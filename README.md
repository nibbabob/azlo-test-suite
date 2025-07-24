# 🧪 Azlo Test Suite

> **A beautiful, real-time Go testing dashboard with comprehensive coverage visualization**  
> *Part of the professional development tools from [Azlo.pro](https://www.azlo.pro)*

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version" />
  <img src="https://img.shields.io/badge/WebSocket-Real--time-4CAF50?style=for-the-badge&logo=websocket&logoColor=white" alt="WebSocket" />
  <img src="https://img.shields.io/badge/Coverage-HTML%20%2B%20Interactive-2196F3?style=for-the-badge&logo=go&logoColor=white" alt="Coverage" />
  <img src="https://img.shields.io/badge/Dark%20Theme-Professional-6c2c91?style=for-the-badge&logo=visualstudiocode&logoColor=white" alt="Dark Theme" />
</p>

The **Azlo Test Suite** transforms the tedious process of running Go tests into a delightful, productive experience. Built with the problem-solving philosophy of [Azlo.pro](https://www.azlo.pro), this tool eliminates the friction between writing code and understanding test coverage.

## 🎯 The Problem We Solve

**Before:** Developers waste time switching between terminal, coverage files, and IDEs to understand test results and coverage gaps.

**After:** One elegant dashboard that provides real-time test feedback, interactive coverage exploration, and professional HTML reports - all styled for modern development workflows.

---

## ✨ Features That Drive Results

### **🚀 Real-Time Test Execution**
- **WebSocket-powered updates**: See test results instantly as they complete
- **Multi-package discovery**: Automatically finds and tests all packages with test files
- **One-click testing**: Run comprehensive test suites with a single button
- **Execution timing**: Monitor test performance and identify slow tests

### **📊 Dual Coverage Visualization**
- **Interactive Coverage Explorer**: Browse coverage file-by-file with line-level highlighting
- **Native Go HTML Reports**: Professional coverage reports with custom dark theme styling
- **Coverage Analytics**: Color-coded badges and metrics (green ≥80%, yellow 60-79%, red <60%)
- **Real-time Updates**: Watch coverage improve as you write tests

### **🎨 Professional Dark Theme UI**
- **Modern Design**: Clean, distraction-free interface optimized for extended use
- **Syntax Highlighting**: Proper Go code highlighting in coverage views
- **Responsive Layout**: Works seamlessly on different screen sizes
- **Developer-Focused**: Built by developers, for developers

### **📁 Flexible Project Management**
- **Any Go Project**: Works with any Go project structure
- **Easy Project Switching**: Browse folders or enter paths manually
- **Project Validation**: Ensures selected directories are valid Go projects
- **Cross-Platform**: Windows, macOS, and Linux support

---

## 🚀 Quick Start

### Prerequisites
- **Go 1.24+** installed on your system
- A Go project with test files (`*_test.go`)

### Installation & Setup

1. **Clone the Azlo Test Suite**:
   ```bash
   git clone <repository-url> azlo-test-suite
   cd azlo-test-suite
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Start the dashboard**:
   ```bash
   go run main.go
   ```

4. **Open your browser** to `http://localhost:8484`

### First Run Experience

1. **Select Your Project**: Click "📁 Select Project" to choose your Go project folder
2. **Run Tests**: Click "Run Tests" to execute your test suite
3. **Explore Coverage**: Use both the interactive coverage viewer and HTML reports
4. **Keep It Open**: Perfect for a second monitor while you develop

---

## 📊 Understanding Your Dashboard

### **Project Info Bar**
Shows your current project name and path - always know which project you're testing.

### **Stats Overview**
- **Overall Coverage**: Project-wide coverage percentage with color coding
- **Total Packages**: Number of packages containing tests
- **Passed/Failed**: Quick status overview across all packages

### **Package Results**
Each package displays:
- **Status Border**: Green (passed) or red (failed)
- **Coverage Badge**: Color-coded coverage percentage
- **Duration**: Test execution time in milliseconds
- **Expandable Details**: Click to see full test output

### **Coverage Options**
For each package with coverage data:
- **📊 View Coverage**: Interactive file browser with line-by-line highlighting
- **📋 HTML Report**: Go's native coverage report (opens in new tab)

---

## 🛠️ Advanced Configuration

### **Environment Variables**
```bash
# Custom port (default: 8484)
PORT=3000 go run main.go

# Example with different port
PORT=9090 go run main.go
```

### **Project Structure Support**
Works with any Go project layout:
```
your-project/
├── go.mod
├── main.go
├── pkg/
│   ├── auth/
│   │   ├── auth.go
│   │   └── auth_test.go
│   └── handlers/
│       ├── handlers.go
│       └── handlers_test.go
├── internal/
│   └── service/
│       ├── service.go
│       └── service_test.go
└── cmd/
    └── app/
        └── main.go
```

---

## 🎯 Best Practices for Maximum Impact

### **Development Workflow**
1. **Keep the dashboard visible** on a second monitor or split screen
2. **Write tests incrementally** and watch coverage improve in real-time
3. **Use HTML reports** for comprehensive coverage analysis
4. **Set coverage goals** - aim for green badges (≥80% coverage)
5. **Monitor test duration** to identify performance bottlenecks

### **Team Usage**
- **Code Reviews**: Use HTML coverage reports to identify untested code paths
- **CI/CD Integration**: Validate coverage thresholds before deployment
- **Onboarding**: Help new team members understand project test coverage
- **Technical Debt**: Track coverage improvements over time

---

## 🔧 Technical Architecture

### **Backend (Go)**
- **Gorilla WebSocket**: Real-time client communication
- **Gorilla Mux**: HTTP routing and static file serving
- **Go Toolchain Integration**: Native `go test` and `go tool cover` usage
- **File System Monitoring**: Automatic cleanup of temporary coverage files

### **Frontend (Vanilla JS)**
- **WebSocket Client**: Real-time dashboard updates
- **Modern CSS**: Professional dark theme with responsive design
- **No Dependencies**: Pure JavaScript for maximum compatibility
- **File System Access API**: Modern folder selection (with fallback)

### **Security & Privacy**
- **Local Development Only**: No data leaves your machine
- **Temporary Files**: Automatic cleanup of coverage reports
- **No Tracking**: Built with privacy-first principles from [Azlo.pro](https://www.azlo.pro)

---

## 🚀 API Reference

### **WebSocket Endpoints**
- `GET /ws` - Real-time dashboard updates

### **HTTP Endpoints**
- `POST /run-tests` - Trigger test execution
- `GET /coverage/{package}` - Get package coverage data
- `GET /html-coverage/{filename}` - Serve styled HTML coverage reports
- `POST /set-project-path` - Change project directory
- `GET /project-info` - Get current project information

---

## 🛠️ Troubleshooting

### **Common Issues**

**No packages found**
- Ensure your project has `*_test.go` files
- Check that tests are in the same package as your code

**Coverage not showing**
- Verify your tests actually execute code paths
- Check that `go test -cover` works from command line

**WebSocket connection failed**
- Ensure port 8484 is available
- Try setting a different port: `PORT=9090 go run main.go`

**Project path issues**
- Use absolute paths when possible
- Ensure the directory contains a `go.mod` file or `.go` files

### **Debug Mode**
Run with verbose logging:
```bash
go run main.go 2>&1 | tee dashboard.log
```

---

## 💡 Why Choose Azlo Test Suite?

### **Problem-Focused Design**
Built to solve real developer pain points - not just another dashboard. Every feature addresses specific testing workflow inefficiencies.

### **Production-Ready Architecture**
Enterprise-grade patterns that handle real projects, complex structures, and team environments.

### **Developer Experience First**
Created by developers who understand the daily grind of test-driven development and coverage optimization.

---

## 🤝 About Azlo.pro

The Azlo Test Suite represents the problem-solving philosophy of [**Azlo.pro**](https://www.azlo.pro) - custom development solutions that eliminate manual work and streamline business processes.

**Azlo.pro specializes in:**
- 🤖 **Custom Automation**: Robust workflows beyond Zapier limitations
- ⚡ **High-Performance Backends**: Blazing-fast APIs with Go & Rust
- 🚀 **Rapid MVP Development**: Ideas to working prototypes at incredible speed
- 🧠 **AI Integration**: LLM-powered workflow automation
- 📊 **Data Unification**: Transform messy data into actionable insights

*Built with ADHD-powered focus and privacy-first principles.*

---

## 📞 Support & Contributions

**Need Help?**
- Check the troubleshooting section above
- Review the [GitHub Issues](repository-issues-url)
- Contact: [christian.nielsen@azlo.pro](mailto:christian.nielsen@azlo.pro)

**Want to Contribute?**
This tool is designed to be easily customizable for specific team needs. Fork it, enhance it, make it yours.

---

## 📄 License

This project is provided as-is for educational and development use.

---

<p align="center">
  <strong>Stop drowning in manual testing workflows.</strong><br>
  <em>Let the Azlo Test Suite automate your Go testing experience.</em>
</p>

<p align="center">
  <a href="https://www.azlo.pro">🌐 Learn more about Azlo.pro custom development services</a>
</p>