# 🧪 Azlo Test Suite

> **A beautiful, real-time Go testing dashboard with comprehensive coverage visualization**  
> *Part of the professional development tools from [Azlo.pro](https://www.azlo.pro)*

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version" />
  <img src="https://img.shields.io/badge/WebSocket-Real--time-4CAF50?style=for-the-badge&logo=websocket&logoColor=white" alt="WebSocket" />
  <img src="https://img.shields.io/badge/Coverage-HTML%20%2B%20Interactive-2196F3?style=for-the-badge&logo=go&logoColor=white" alt="Coverage" />
  <img src="https://img.shields.io/badge/Dark%20Theme-Professional-6c2c91?style=for-the-badge&logo=visualstudiocode&logoColor=white" alt="Dark Theme" />
</p>


Hey! So, check this out. I made a thing.

## The Problem I Got Tired Of

**Before:** Jumping between the terminal, a browser tab for coverage, and my code editor just to see if my tests passed. It's clunky and kills the flow.

**After:** One clean dashboard. You see test results pop up in real-time. You can explore your code coverage interactively. It just works.

---

## ✨ What It Actually Does

### **🚀 Real-Time Everything**
Powered by WebSockets, you see tests pass or fail the second they're done. No more waiting. It finds all your test packages automatically and tells you how long they took to run.

### **📊 See Your Coverage Two Ways**
* An **Interactive Explorer** to click through files and see covered lines.
* The standard Go **HTML Report**, but styled with a proper dark theme so it doesn't burn your eyes out.
* You get a clear coverage score, color-coded so you know if you're doing good (green's the goal!).

### **🎨 A UI That's Actually Nice to Look At**
Clean, dark theme. It's designed for people who spend hours looking at code.

### **📁 It Just Works**
Pick your Go project folder, and you're off. Doesn't matter how your project is structured.

---

## 🚀 Getting Started

**Prerequisites:**
* Go 1.24+
* A Go project with some `*_test.go` files.

**Let's Go:**

1. **Clone the thing:**
   ```bash
   git clone <repository-url> azlo-test-suite
   cd azlo-test-suite
   ```

2. **Install the stuff it needs:**
   ```bash
   go mod tidy
   ```

3. **Run it:**
   ```bash
   go run main.go
   ```

4. **Open `http://localhost:8484` in your browser.**

Stick the dashboard on a second monitor while you code. Write a test, save the file, and watch the coverage percentage go up. It's a great feeling.

---

## 📊 What You'll See

### **Project Info Bar**
Shows your current project name and path - always know which project you're testing.

### **Stats Overview**
* **Overall Coverage**: Project-wide coverage with color coding (green ≥80%, yellow 60-79%, red <60%)
* **Total Packages**: Number of packages with tests
* **Passed/Failed**: Quick status overview

### **Package Results**
Each package shows:
* **Status**: Green border (passed) or red border (failed)
* **Coverage Badge**: Color-coded percentage
* **Duration**: How long tests took
* **Expandable Details**: Click to see full test output

### **Coverage Options**
For packages with coverage:
* **📊 View Coverage**: Interactive file browser with line-by-line highlighting
* **📋 HTML Report**: Go's native coverage report (opens in new tab)

---

## 🛠️ How It's Built

I love to orchestrate systems and see them run smoothly. This is no different.

### **Backend (Go)**
Just pure Go with Gorilla WebSocket for the real-time magic and Mux for handling requests. It calls the `go test` and `go tool cover` commands you already use.

### **Frontend (Vanilla JS)**
No heavy frameworks. Just simple, clean JavaScript that talks to the backend and makes everything look good.

### **Privacy-First**
This runs entirely on your machine. No data goes anywhere. No tracking. That's a big deal for me.

---

## 🎯 Usage Tips

### **Daily Workflow**
1. **Keep it visible** on a second monitor or split screen
2. **Write tests incrementally** and watch coverage improve
3. **Use HTML reports** for deep coverage analysis
4. **Set coverage goals** - aim for those green badges
5. **Monitor test speed** to catch slow tests early

### **Team Stuff**
* **Code Reviews**: Share HTML coverage reports
* **CI/CD**: Validate coverage before deployment
* **Onboarding**: Help new people understand test coverage

---

## 🔧 Configuration

### **Custom Port**
```bash
# Different port if 8484 is busy
PORT=3000 go run main.go
```

### **Project Structure**
Works with any Go layout. Whether you've got:
```
your-project/
├── go.mod
├── main.go
├── pkg/
│   └── something/
│       ├── code.go
│       └── code_test.go
└── internal/
    └── stuff/
        ├── more.go
        └── more_test.go
```

It'll find your tests and run them.

---

## 🚨 If Things Break

**No packages found?**
* Make sure you have `*_test.go` files
* Check that `go test` works from the command line

**Coverage not showing?**
* Verify your tests actually run code
* Try `go test -cover` manually first

**Can't connect?**
* Port 8484 might be busy - try `PORT=9090 go run main.go`
* Check the console for error messages

**Project path issues?**
* Use the full path to your project
* Make sure there's a `go.mod` file or `.go` files in the directory

---

## 💡 Why I Made This

Because I wanted it. It's a tool built to solve a problem that was bugging me, designed with the same problem-solving energy I put into **[Azlo.pro](https://www.azlo.pro)**. 

At Azlo.pro, I build custom automation and high-performance backends that eliminate manual work. This dashboard scratches the same itch - taking something tedious (checking test results) and making it smooth and enjoyable.

It's about creating things that are not only functional but also a pleasure to use. I'm the "glass half full" type, and I believe our tools should reflect that.

Built with ADHD-powered focus and a love for elegant systems. I hope you find it useful!

---

## 🤝 Questions? Issues?

Fork it, mess with it, make it your own. If you get stuck, open an issue on GitHub or shoot me an email: [christian.nielsen@azlo.pro](mailto:christian.nielsen@azlo.pro)

---

<p align="center">
  <strong>Stop context-switching between terminal and browser.</strong><br>
  <em>Let the dashboard do the work.</em>
</p>