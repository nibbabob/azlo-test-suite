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
}

// --- Core Dashboard Component ---

type TestDashboard struct {
	Data      DashboardData
	Clients   map[*websocket.Conn]bool
	Broadcast chan DashboardData
	Upgrader  websocket.Upgrader
}

func NewTestDashboard() *TestDashboard {
	return &TestDashboard{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan DashboardData),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}
}

// --- Exported Methods (for use by handlers/main) ---

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
	log.Println("Running tests...")

	packages, err := findGoPackages(".")
	if err != nil {
		log.Printf("Error finding packages: %v", err)
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
	}

	td.Broadcast <- dashboardData
	log.Println("Test run complete.")
}

// --- Internal Helper Functions (unexported) ---

// runPackageTests executes tests for a single package.
func (td *TestDashboard) runPackageTests(pkg string) TestResult {
	start := time.Now()
	coverProfile := fmt.Sprintf("coverage_%s.out", strings.ReplaceAll(pkg, "/", "_"))
	defer os.Remove(coverProfile)

	cmd := exec.Command("go", "test", "-coverprofile="+coverProfile, "./"+pkg)
	output, err := cmd.CombinedOutput()

	result := TestResult{
		Package:   pkg,
		Passed:    err == nil,
		Output:    string(output),
		Duration:  time.Since(start),
		Timestamp: time.Now(),
	}

	if fileExists(coverProfile) {
		result.Coverage = extractCoverage(string(output))
		result.Files = td.parseCoverageProfile(coverProfile, pkg)
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
		content, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("Error reading file %s: %v", filename, err)
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

func findGoPackages(root string) ([]string, error) {
	var packages []string
	uniquePackages := make(map[string]bool)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "vendor" || d.Name() == ".git" || d.Name() == "static" {
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
