package dashboard

import (
	"bufio"
	"bytes" // Added this import
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// --- Data Structures (No Changes) ---

type CoverageBlock struct {
	StartLine int  `json:"start_line"`
	EndLine   int  `json:"end_line"`
	Count     int  `json:"count"`
	Covered   bool `json:"covered"`
}

type FileCoverage struct {
	Filename string          `json:"filename"`
	Content  string          `json:"content"`
	Blocks   []CoverageBlock `json:"blocks"`
	Coverage float64         `json:"coverage"`
}

type TestResult struct {
	Package          string         `json:"package"`
	Passed           bool           `json:"passed"`
	Output           string         `json:"output"`
	Duration         time.Duration  `json:"duration"`
	Coverage         float64        `json:"coverage"`
	Files            []FileCoverage `json:"files"`
	Timestamp        time.Time      `json:"timestamp"`
	HTMLCoverageFile string         `json:"html_coverage_file,omitempty"`
}

type DashboardData struct {
	Results         []TestResult `json:"results"`
	OverallCoverage float64      `json:"overall_coverage"`
	TotalTests      int          `json:"total_tests"`
	PassedTests     int          `json:"passed_tests"`
	LastRun         time.Time    `json:"last_run"`
	ProjectPath     string       `json:"project_path"`
	ProjectName     string       `json:"project_name"`
}

// --- Core Dashboard Component (No Changes) ---

type TestDashboard struct {
	Data        DashboardData
	Clients     map[*websocket.Conn]bool
	Broadcast   chan DashboardData
	Upgrader    websocket.Upgrader
	ProjectPath string
	HTMLFiles   map[string]time.Time
}

func NewTestDashboard() *TestDashboard {
	currentDir, _ := os.Getwd()
	td := &TestDashboard{
		Clients:     make(map[*websocket.Conn]bool),
		Broadcast:   make(chan DashboardData),
		ProjectPath: currentDir,
		HTMLFiles:   make(map[string]time.Time),
		Data: DashboardData{
			ProjectPath: currentDir,
			ProjectName: filepath.Base(currentDir),
		},
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	go td.cleanupHTMLFiles()
	return td
}

// --- Exported Methods (No Changes) ---

func (td *TestDashboard) SetProjectPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path does not exist: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}
	if !td.isGoProject(path) {
		return fmt.Errorf("directory does not appear to be a Go project (no go.mod or *.go files found)")
	}
	td.ProjectPath = path
	td.Data.ProjectPath = path
	td.Data.ProjectName = filepath.Base(path)
	log.Printf("Project path changed to: %s", path)
	td.Broadcast <- td.Data
	return nil
}

func (td *TestDashboard) GetProjectInfo() map[string]interface{} {
	packages, _ := td.findGoPackages(td.ProjectPath)
	return map[string]interface{}{
		"project_path":   td.ProjectPath,
		"project_name":   td.Data.ProjectName,
		"packages_found": len(packages),
		"packages":       packages,
	}
}

func (td *TestDashboard) BroadcastUpdates() {
	for data := range td.Broadcast {
		td.Data = data
		for client := range td.Clients {
			err := client.WriteJSON(data)
			if err != nil {
				client.Close()
				delete(td.Clients, client)
			}
		}
	}
}

func (td *TestDashboard) RunTests() {
	log.Printf("Running tests in project: %s", td.ProjectPath)

	packages, err := td.findGoPackages(td.ProjectPath)
	if err != nil {
		log.Printf("Error finding packages: %v", err)
		return
	}

	// Broadcast a "starting" state to clear the UI and show the total package count.
	startingData := DashboardData{
		Results:         []TestResult{},
		OverallCoverage: 0.0,
		TotalTests:      len(packages),
		PassedTests:     0,
		LastRun:         time.Now(),
		ProjectPath:     td.ProjectPath,
		ProjectName:     td.Data.ProjectName,
	}
	td.Broadcast <- startingData

	if len(packages) == 0 {
		log.Printf("No test packages found in %s", td.ProjectPath)
		return
	}

	var results []TestResult // This will hold the results as they come in.

	for _, pkg := range packages {
		result := td.runPackageTests(pkg)
		results = append(results, result) // Add the new result to our list

		// Recalculate stats based on the results we have so far
		var totalCoverage float64
		var passedTests int
		for _, r := range results {
			totalCoverage += r.Coverage
			if r.Passed {
				passedTests++
			}
		}

		// The denominator is the number of packages completed so far, giving a running average.
		overallCoverage := 0.0
		if len(results) > 0 {
			overallCoverage = totalCoverage / float64(len(results))
		}

		// Create and broadcast the intermediate data object
		intermediateData := DashboardData{
			Results:         results, // Send the list of packages completed so far
			OverallCoverage: overallCoverage,
			TotalTests:      len(packages), // Total expected packages
			PassedTests:     passedTests,
			LastRun:         time.Now(),
			ProjectPath:     td.ProjectPath,
			ProjectName:     td.Data.ProjectName,
		}
		td.Broadcast <- intermediateData // Send the live update
	}

	log.Printf("Test run complete. Found %d packages, %d passed", len(packages), len(results))
}

// ---- NEW HELPER FUNCTION ----
// This function reads the go.mod file to find the project's module name.
func getModuleName(projectPath string) (string, error) {
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("could not read go.mod: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning go.mod: %w", err)
	}

	return "", fmt.Errorf("module directive not found in go.mod")
}

// --- Internal Helper Functions (with one modification) ---

func (td *TestDashboard) isGoProject(path string) bool {
	if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
		return true
	}
	hasGoFiles := false
	filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "vendor" || d.Name() == ".git" || d.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(filePath, ".go") {
			hasGoFiles = true
			return fmt.Errorf("found go files")
		}
		return nil
	})
	return hasGoFiles
}

func (td *TestDashboard) runPackageTests(pkg string) TestResult {
	start := time.Now()
	coverProfile := fmt.Sprintf("coverage_%s_%d.out",
		strings.ReplaceAll(strings.ReplaceAll(pkg, "/", "_"), string(filepath.Separator), "_"),
		time.Now().UnixNano())
	htmlCoverageFile := fmt.Sprintf("coverage_%s_%d.html",
		strings.ReplaceAll(strings.ReplaceAll(pkg, "/", "_"), string(filepath.Separator), "_"),
		time.Now().UnixNano())

	defer func() {
		os.Remove(coverProfile)
	}()

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(td.ProjectPath)
	if err != nil {
		return TestResult{
			Package:   pkg,
			Passed:    false,
			Output:    fmt.Sprintf("Error changing to project directory: %v", err),
			Duration:  time.Since(start),
			Timestamp: time.Now(),
		}
	}

	relPkg, err := filepath.Rel(td.ProjectPath, pkg)
	if err != nil {
		relPkg = pkg
	}
	relPkg = filepath.ToSlash(relPkg)
	if relPkg == "." {
		relPkg = "./"
	} else if !strings.HasPrefix(relPkg, "./") {
		relPkg = "./" + relPkg
	}

	cmd := exec.Command("go", "test", "-coverprofile="+coverProfile, relPkg)
	cmd.Dir = td.ProjectPath
	output, testErr := cmd.CombinedOutput()

	result := TestResult{
		Package:   relPkg,
		Passed:    testErr == nil,
		Output:    string(output),
		Duration:  time.Since(start),
		Timestamp: time.Now(),
	}

	coveragePath := filepath.Join(td.ProjectPath, coverProfile)
	if fileExists(coveragePath) {
		result.Coverage = extractCoverage(string(output))
		result.Files = td.parseCoverageProfile(coveragePath, pkg) // This function is now fixed
		htmlPath := filepath.Join(td.ProjectPath, htmlCoverageFile)
		if td.generateHTMLCoverage(coveragePath, htmlPath) {
			result.HTMLCoverageFile = htmlCoverageFile
			td.HTMLFiles[htmlCoverageFile] = time.Now().Add(1 * time.Hour)
		}
	}

	return result
}

// ---- MODIFIED FUNCTION ----
func (td *TestDashboard) parseCoverageProfile(profilePath string, _ string) []FileCoverage {
	file, err := os.Open(profilePath)
	if err != nil {
		log.Printf("Error opening coverage profile %s: %v", profilePath, err)
		return nil
	}
	defer file.Close()

	// Get the module name to correctly resolve file paths
	moduleName, modErr := getModuleName(td.ProjectPath)

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() { // Skip the "mode: set" line
		return nil
	}

	fileMap := make(map[string][]CoverageBlock)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}
		locationPart := parts[0]
		count, _ := strconv.Atoi(parts[2])
		colonIdx := strings.LastIndex(locationPart, ":")
		if colonIdx == -1 {
			continue
		}
		filename := locationPart[:colonIdx]
		rangeParts := strings.Split(locationPart[colonIdx+1:], ",")
		if len(rangeParts) != 2 {
			continue
		}
		startParts := strings.Split(rangeParts[0], ".")
		endParts := strings.Split(rangeParts[1], ".")
		if len(startParts) < 1 || len(endParts) < 1 {
			continue
		}
		startLine, _ := strconv.Atoi(startParts[0])
		endLine, _ := strconv.Atoi(endParts[0])
		block := CoverageBlock{
			StartLine: startLine,
			EndLine:   endLine,
			Count:     count,
			Covered:   count > 0,
		}
		fileMap[filename] = append(fileMap[filename], block)
	}

	var files []FileCoverage
	for filename, blocks := range fileMap {
		fullPath := filename
		if !filepath.IsAbs(filename) {
			// This is the new logic to construct the correct path
			relativePath := filename
			// If we found a module name and the path from the coverage file starts with it...
			if modErr == nil && strings.HasPrefix(filename, moduleName+"/") {
				// ...then we strip that module name prefix to get a true relative path.
				relativePath = strings.TrimPrefix(filename, moduleName+"/")
			}
			fullPath = filepath.Join(td.ProjectPath, relativePath)
		}

		content, err := os.ReadFile(fullPath)
		if err != nil {
			// This is the error message that should now disappear
			log.Printf("Error reading file %s: %v", fullPath, err)
			continue
		}

		files = append(files, FileCoverage{
			Filename: filename,
			Content:  string(content),
			Blocks:   blocks,
			Coverage: calculateFileCoverage(blocks),
		})
	}
	return files
}

// --- Other helpers (No Changes) ---

func (td *TestDashboard) findGoPackages(root string) ([]string, error) {
	var packages []string
	uniquePackages := make(map[string]bool)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "vendor" || d.Name() == ".git" || d.Name() == "node_modules" || d.Name() == ".vscode" || d.Name() == ".idea" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			dir := filepath.Dir(path)
			if !uniquePackages[dir] {
				uniquePackages[dir] = true
				packages = append(packages, dir)
			}
		}
		return nil
	})
	return packages, err
}

func (td *TestDashboard) generateHTMLCoverage(profilePath, htmlPath string) bool {
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(td.ProjectPath)
	if err != nil {
		log.Printf("Error changing to project directory for HTML coverage: %v", err)
		return false
	}
	cmd := exec.Command("go", "tool", "cover", "-html="+filepath.Base(profilePath), "-o", filepath.Base(htmlPath))
	cmd.Dir = td.ProjectPath
	if err := cmd.Run(); err != nil {
		log.Printf("Error generating HTML coverage: %v", err)
		return false
	}
	return fileExists(htmlPath)
}

func (td *TestDashboard) GetHTMLCoverage(filename string) (string, error) {
	htmlPath := filepath.Join(td.ProjectPath, filename)
	if !fileExists(htmlPath) {
		return "", fmt.Errorf("HTML coverage file not found: %s", filename)
	}
	content, err := os.ReadFile(htmlPath)
	if err != nil {
		return "", fmt.Errorf("error reading HTML coverage file: %v", err)
	}
	htmlContent := string(content)
	customCSS := td.getCustomCoverageCSS()
	if headEndIndex := strings.Index(htmlContent, "</head>"); headEndIndex != -1 {
		htmlContent = htmlContent[:headEndIndex] +
			"\n<style>\n" + customCSS + "\n</style>\n" +
			htmlContent[headEndIndex:]
	}
	return htmlContent, nil
}

func (td *TestDashboard) getCustomCoverageCSS() string {
	return `
/* Custom dark theme for Go coverage reports */
body {
    background-color: #1a1a1a !important; color: #e0e0e0 !important;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif !important;
    margin: 0 !important; padding: 20px !important;
}
pre {
    background: #1a1a1a !important; color: #e0e0e0 !important;
}
.cov0 { background-color: rgba(244, 67, 54, 0.3) !important; color: #ffffff !important; }
.cov1, .cov2, .cov3, .cov4, .cov5, .cov6, .cov7, .cov8, .cov9, .cov10 {
    background-color: rgba(76, 175, 80, 0.3) !important; color: #ffffff !important;
}
a { color: #4CAF50 !important; }
* { border-color: #404040 !important; }
`
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func calculateFileCoverage(blocks []CoverageBlock) float64 {
	if len(blocks) == 0 {
		return 0
	}
	coveredBlocks := 0
	for _, block := range blocks {
		if block.Covered {
			coveredBlocks++
		}
	}
	return float64(coveredBlocks) / float64(len(blocks)) * 100
}

func extractCoverage(output string) float64 {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "coverage:") {
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "coverage:" && i+1 < len(parts) {
					coverageStr := strings.TrimSuffix(parts[i+1], "%")
					if coverage, err := strconv.ParseFloat(coverageStr, 64); err == nil {
						return coverage
					}
				}
			}
		}
	}
	return 0.0
}

func (td *TestDashboard) cleanupHTMLFiles() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		for filename, expiry := range td.HTMLFiles {
			if now.After(expiry) {
				filePath := filepath.Join(td.ProjectPath, filename)
				if err := os.Remove(filePath); err == nil {
					log.Printf("Cleaned up expired HTML coverage file: %s", filename)
				}
				delete(td.HTMLFiles, filename)
			}
		}
	}
}
