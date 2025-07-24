package dashboard

import (
	"bufio"
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

// --- Data Structures ---

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
	Package   string         `json:"package"`
	Passed    bool           `json:"passed"`
	Output    string         `json:"output"`
	Duration  time.Duration  `json:"duration"`
	Coverage  float64        `json:"coverage"`
	Files     []FileCoverage `json:"files"`
	Timestamp time.Time      `json:"timestamp"`
}

type DashboardData struct {
	Results         []TestResult `json:"results"`
	OverallCoverage float64      `json:"overall_coverage"`
	TotalTests      int          `json:"total_tests"`
	PassedTests     int          `json:"passed_tests"`
	LastRun         time.Time    `json:"last_run"`
	ProjectPath     string       `json:"project_path"` // New field
	ProjectName     string       `json:"project_name"` // New field
}

// --- Core Dashboard Component ---

type TestDashboard struct {
	Data        DashboardData
	Clients     map[*websocket.Conn]bool
	Broadcast   chan DashboardData
	Upgrader    websocket.Upgrader
	ProjectPath string // Current project path
}

func NewTestDashboard() *TestDashboard {
	// Default to current directory
	currentDir, _ := os.Getwd()

	return &TestDashboard{
		Clients:     make(map[*websocket.Conn]bool),
		Broadcast:   make(chan DashboardData),
		ProjectPath: currentDir,
		Data: DashboardData{
			ProjectPath: currentDir,
			ProjectName: filepath.Base(currentDir),
		},
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}
}

// --- Exported Methods (for use by handlers/main) ---

// SetProjectPath changes the project root directory
func (td *TestDashboard) SetProjectPath(path string) error {
	// Validate that the path exists and is a directory
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path does not exist: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	// Check if it's a Go project (has go.mod or *.go files)
	if !td.isGoProject(path) {
		return fmt.Errorf("directory does not appear to be a Go project (no go.mod or *.go files found)")
	}

	td.ProjectPath = path
	td.Data.ProjectPath = path
	td.Data.ProjectName = filepath.Base(path)

	log.Printf("Project path changed to: %s", path)

	// Broadcast updated project info to all clients
	td.Broadcast <- td.Data

	return nil
}

// GetProjectInfo returns current project information
func (td *TestDashboard) GetProjectInfo() map[string]interface{} {
	packages, _ := td.findGoPackages(td.ProjectPath)

	return map[string]interface{}{
		"project_path":   td.ProjectPath,
		"project_name":   td.Data.ProjectName,
		"packages_found": len(packages),
		"packages":       packages,
	}
}

// BroadcastUpdates sends dashboard data to all connected WebSocket clients.
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

// RunTests discovers all Go packages with tests and executes them.
func (td *TestDashboard) RunTests() {
	log.Printf("Running tests in project: %s", td.ProjectPath)

	packages, err := td.findGoPackages(td.ProjectPath)
	if err != nil {
		log.Printf("Error finding packages: %v", err)
		return
	}

	if len(packages) == 0 {
		log.Printf("No test packages found in %s", td.ProjectPath)
		// Still broadcast an update with zero results
		dashboardData := DashboardData{
			Results:         []TestResult{},
			OverallCoverage: 0.0,
			TotalTests:      0,
			PassedTests:     0,
			LastRun:         time.Now(),
			ProjectPath:     td.ProjectPath,
			ProjectName:     td.Data.ProjectName,
		}
		td.Broadcast <- dashboardData
		return
	}

	var results []TestResult
	var totalCoverage float64
	var passedTests int

	for _, pkg := range packages {
		result := td.runPackageTests(pkg)
		results = append(results, result)
		totalCoverage += result.Coverage
		if result.Passed {
			passedTests++
		}
	}

	overallCoverage := 0.0
	if len(results) > 0 {
		overallCoverage = totalCoverage / float64(len(results))
	}

	dashboardData := DashboardData{
		Results:         results,
		OverallCoverage: overallCoverage,
		TotalTests:      len(packages),
		PassedTests:     passedTests,
		LastRun:         time.Now(),
		ProjectPath:     td.ProjectPath,
		ProjectName:     td.Data.ProjectName,
	}

	td.Broadcast <- dashboardData
	log.Printf("Test run complete. Found %d packages, %d passed", len(packages), passedTests)
}

// --- Internal Helper Functions (unexported) ---

// isGoProject checks if a directory contains Go project files
func (td *TestDashboard) isGoProject(path string) bool {
	// Check for go.mod file
	if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
		return true
	}

	// Check for any .go files in the directory or subdirectories
	hasGoFiles := false
	filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Skip common non-source directories
			if d.Name() == "vendor" || d.Name() == ".git" || d.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(filePath, ".go") {
			hasGoFiles = true
			return fmt.Errorf("found go files") // Use error to break out of walk
		}
		return nil
	})

	return hasGoFiles
}

// runPackageTests executes tests for a single package.
func (td *TestDashboard) runPackageTests(pkg string) TestResult {
	start := time.Now()

	// Create a unique coverage profile name
	coverProfile := fmt.Sprintf("coverage_%s_%d.out",
		strings.ReplaceAll(strings.ReplaceAll(pkg, "/", "_"), string(filepath.Separator), "_"),
		time.Now().UnixNano())

	defer os.Remove(coverProfile)

	// Change to project directory and run tests
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

	// Convert absolute package path to relative path from project root
	relPkg, err := filepath.Rel(td.ProjectPath, pkg)
	if err != nil {
		relPkg = pkg
	}

	// Ensure we use forward slashes for go command (works on Windows too)
	relPkg = filepath.ToSlash(relPkg)
	if relPkg == "." {
		relPkg = "./"
	} else if !strings.HasPrefix(relPkg, "./") {
		relPkg = "./" + relPkg
	}

	cmd := exec.Command("go", "test", "-coverprofile="+coverProfile, relPkg)
	cmd.Dir = td.ProjectPath // Ensure command runs in project directory
	output, testErr := cmd.CombinedOutput()

	result := TestResult{
		Package:   relPkg,
		Passed:    testErr == nil,
		Output:    string(output),
		Duration:  time.Since(start),
		Timestamp: time.Now(),
	}

	// Check if coverage file was created and parse it
	coveragePath := filepath.Join(td.ProjectPath, coverProfile)
	if fileExists(coveragePath) {
		result.Coverage = extractCoverage(string(output))
		result.Files = td.parseCoverageProfile(coveragePath, pkg)
	}

	return result
}

// parseCoverageProfile reads a coverage.out file and parses the data.
func (td *TestDashboard) parseCoverageProfile(profilePath string, _ string) []FileCoverage {
	file, err := os.Open(profilePath)
	if err != nil {
		log.Printf("Error opening coverage profile %s: %v", profilePath, err)
		return nil
	}
	defer file.Close()

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
		// Try to read file relative to project path
		fullPath := filename
		if !filepath.IsAbs(filename) {
			fullPath = filepath.Join(td.ProjectPath, filename)
		}

		content, err := os.ReadFile(fullPath)
		if err != nil {
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

// findGoPackages discovers all directories containing Go test files
func (td *TestDashboard) findGoPackages(root string) ([]string, error) {
	var packages []string
	uniquePackages := make(map[string]bool)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Skip common directories that shouldn't contain tests
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

// Helper functions remain the same
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
