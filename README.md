# 🧪 Go Test Dashboard

A beautiful, real-time web-based dashboard for Go testing with coverage visualization. Perfect for keeping on your second monitor while writing unit tests!

## ✨ Features

- **Real-time Updates**: WebSocket-based live updates as tests run
- **Coverage Visualization**: Red/green HTML coverage reports with drill-down capabilities
- **Multi-Package Support**: Automatically discovers and tests all packages with tests
- **Beautiful UI**: Dark theme with color-coded coverage and test status
- **One-Click Testing**: Run all tests with a single button click
- **Detailed Output**: View test output and execution times for each package
- **HTML Coverage Reports**: Click through to detailed Go coverage HTML reports

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or later
- A Go project with test files

### Installation

1. **Clone or create the dashboard files** in your Go project root:
   ```
   your-go-project/
   ├── main.go              # The dashboard server
   ├── go.mod               # Module definition
   ├── calculator/          # Example package
   │   ├── calculator.go
   │   └── calculator_test.go
   └── utils/               # Another example package
       ├── strings.go
       └── strings_test.go
   ```

2. **Initialize the Go module and install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the dashboard**:
   ```bash
   go run main.go
   ```

4. **Open your browser** to `http://localhost:8080`

## 📊 Usage

### Running Tests

- **Auto-run**: Tests automatically run when you first open the dashboard
- **Manual run**: Click the "Run Tests" button to run tests on demand
- **Real-time updates**: The dashboard updates automatically when tests complete

### Understanding the Interface

#### Stats Overview
- **Overall Coverage**: Average coverage across all packages
- **Total Packages**: Number of packages with tests
- **Passed/Failed**: Quick overview of test results

#### Package Results
- **Green border**: All tests passed
- **Red border**: Some tests failed
- **Coverage badge**: Color-coded coverage percentage
  - Green: ≥80% coverage
  - Yellow: 60-79% coverage
  - Red: <60% coverage

#### Detailed Views
- **Click package headers** to expand and see test output
- **View Coverage button** opens detailed HTML coverage reports in new tab
- **Test output** shows the raw Go test output for debugging

## 🛠️ Customization

### Environment Variables

- `PORT`: Set the server port (default: 8080)
  ```bash
  PORT=3000 go run main.go
  ```

### Project Structure

The dashboard automatically discovers packages by looking for `*_test.go` files. It works with any Go project structure:

```
your-project/
├── main.go
├── go.mod
├── pkg1/
│   ├── code.go
│   └── code_test.go
├── pkg2/
│   ├── code.go
│   └── code_test.go
└── internal/
    └── pkg3/
        ├── code.go
        └── code_test.go
```

## 🎯 Best Practices

### For Maximum Effectiveness

1. **Keep it visible**: Run on a second monitor or split screen
2. **Write tests incrementally**: Watch coverage improve in real-time
3. **Use the coverage drill-down**: Click "View Coverage" to see exactly what lines need tests
4. **Set coverage goals**: Aim for the green coverage badge (≥80%)

### Writing Better Tests

The dashboard helps you identify:
- **Uncovered code**: Red areas in coverage reports
- **Failing tests**: Red package borders and detailed error output
- **Slow tests**: Duration shown for each package

## 🔧 Advanced Features

### WebSocket API

The dashboard exposes a WebSocket endpoint at `/ws` for programmatic access:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('Test results:', data);
};
```

### REST API

- `POST /run-tests`: Trigger test execution
- `GET /coverage/{filename}`: Serve coverage HTML files

## 🐛 Troubleshooting

### Common Issues

1. **No packages found**: Ensure you have `*_test.go` files in your project
2. **Coverage not showing**: Check that your tests are in the same package as your code
3. **WebSocket connection failed**: Ensure port 8080 is available or set a different PORT
4. **Tests not running**: Verify your Go installation and that `go test` works from command line

### Debug Mode

Run with verbose output:
```bash
go run main.go 2>&1 | tee dashboard.log
```

## 📝 Example Output

```
🧪 Go Test Dashboard starting on http://localhost:8080
📊 Open in your browser to see live test results and coverage
```

## 🤝 Contributing

This dashboard is designed to be easily customizable. Feel free to:

- Modify the CSS for different themes
- Add new metrics to the dashboard
- Integrate with CI/CD pipelines
- Add notification systems

## 📄 License

This project is provided as-is for educational and development use.

---

**Happy Testing!** 🎉 Keep this dashboard open while you code and watch your test coverage grow in real-time.