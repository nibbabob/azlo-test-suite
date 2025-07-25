let ws;

function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//${window.location.host}/ws`);

    ws.onopen = function() {
        console.log('WebSocket connected');
        document.getElementById('results').innerHTML = '<div class="loading">Ready for testing...</div>';
    };

    ws.onmessage = function(event) {
        const data = JSON.parse(event.data);
        updateDashboard(data);
    };

    ws.onclose = function() {
        console.log('WebSocket disconnected');
        document.getElementById('results').innerHTML = '<div class="loading">Connection lost. Retrying...</div>';
        setTimeout(connectWebSocket, 3000);
    };

    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
    };
}

function updateDashboard(data) {
    // Update project info
    updateProjectInfo(data);

    // Update stats
    const overallCoverage = data.overall_coverage || 0;
    document.getElementById('overall-coverage').textContent = `${overallCoverage.toFixed(1)}%`;
    document.getElementById('total-tests').textContent = data.total_tests || 0;
    document.getElementById('passed-tests').textContent = data.passed_tests || 0;
    document.getElementById('failed-tests').textContent = (data.total_tests || 0) - (data.passed_tests || 0);

    // Update coverage color
    const coverageEl = document.getElementById('overall-coverage');
    coverageEl.className = `stat-number ${getCoverageClass(overallCoverage)}`;

    // Update results
    const resultsEl = document.getElementById('results');

    // This new logic handles the live updates gracefully.
    if (data.results && data.results.length > 0) {
        // If we have results, display them. This will happen for every package update.
        resultsEl.innerHTML = data.results.map(result => createPackageHTML(result)).join('');
        delete resultsEl.dataset.isRunning; // A run with results is no longer in the "initial" state
    } else if (data.results && resultsEl.dataset.isRunning === "true") {
        // This handles the very first message of a test run, which has 0 results.
        // We clear the "Running tests..." message to make way for the incoming package results.
        resultsEl.innerHTML = '';
    } else if (data.results) {
        // This handles the case where there are no results and a test is NOT running.
        resultsEl.innerHTML = '<div class="loading">No test results yet. Click "Run Tests".</div>';
    }


    // Update last run time
    if (data.last_run) {
        const lastRunEl = document.getElementById('last-run');
        const lastRun = new Date(data.last_run);
        lastRunEl.textContent = `Last run: ${lastRun.toLocaleTimeString()}`;
    }
}

function updateProjectInfo(data) {
    if (data.project_name) {
        document.getElementById('project-name').textContent = data.project_name;
    }
    if (data.project_path) {
        document.getElementById('project-path-text').textContent = data.project_path;
    }
}

function getCoverageClass(coverage) {
    if (coverage >= 80) return 'coverage-high';
    if (coverage >= 60) return 'coverage-medium';
    return 'coverage-low';
}

function escapeHtml(text) {
    if (text === null || typeof text === 'undefined') {
        return '';
    }
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function createPackageHTML(result) {
    const statusClass = result.passed ? 'passed' : 'failed';
    const statusText = result.passed ? 'PASSED' : 'FAILED';
    const coverageClass = getCoverageClass(result.coverage || 0);
    const duration = result.duration ? (result.duration / 1000000).toFixed(0) : '0';

    let coverageButtons = '';
    if (result.files && result.files.length > 0) {
        coverageButtons += `<button class="coverage-link" onclick="event.stopPropagation(); showCoverage('${escapeHtml(result.package || '')}')">📊 View Coverage</button>`;
    }
    if (result.html_coverage_file) {
        coverageButtons += `<button class="coverage-link html-coverage-link" onclick="event.stopPropagation(); openHTMLCoverage('${escapeHtml(result.html_coverage_file || '')}')">📋 HTML Report</button>`;
    }

    return `
        <div class="package-result ${statusClass}">
            <div class="package-header" onclick="togglePackage(this)">
                <div class="package-name">${escapeHtml(result.package || '')}</div>
                <div class="package-stats">
                    <div class="coverage-badge ${coverageClass}">${(result.coverage || 0).toFixed(1)}%</div>
                    <div class="status-badge ${statusClass}">${statusText}</div>
                    <div class="duration">${duration}ms</div>
                </div>
            </div>
            <div class="package-details">
                <div class="test-output">${escapeHtml(result.output || '')}</div>
                <div class="coverage-buttons">${coverageButtons}</div>
            </div>
        </div>`;
}

// --- PROJECT PATH SELECTION FUNCTIONS ---

function showProjectModal() {
    document.getElementById('project-modal').classList.add('show');
    document.body.style.overflow = 'hidden';
    loadCurrentProjectInfo();
}

function closeProjectModal() {
    document.getElementById('project-modal').classList.remove('show');
    document.body.style.overflow = '';
}

async function setProjectPath() {
    const pathInput = document.getElementById('manual-path-input');
    const path = pathInput.value.trim();
    
    if (!path) {
        alert('Please enter a project path');
        return;
    }

    try {
        const response = await fetch('/set-project-path', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ path: path })
        });

        const result = await response.json();
        
        if (result.success) {
            closeProjectModal();
            // The WebSocket will automatically update the UI with new project info
            setTimeout(() => {
                alert('Project path updated successfully!');
            }, 100);
        } else {
            alert(`Error: ${result.message}`);
        }
    } catch (error) {
        console.error('Error setting project path:', error);
        alert('Error setting project path. Please check the console for details.');
    }
}

async function loadCurrentProjectInfo() {
    try {
        const response = await fetch('/project-info');
        const info = await response.json();
        
        const infoContainer = document.getElementById('current-project-info');
        
        let packagesHtml = '';
        if (info.packages && info.packages.length > 0) {
            packagesHtml = `
                <div class="project-info-item">
                    <span class="project-info-label">Test Packages:</span>
                    <span class="project-info-value">${info.packages_found}</span>
                </div>
                <div class="project-packages-list">
                    ${info.packages.map(pkg => `<div class="project-package-item">${pkg}</div>`).join('')}
                </div>
            `;
        }
        
        infoContainer.innerHTML = `
            <div class="project-info-item">
                <span class="project-info-label">Current Path:</span>
                <span class="project-info-value">${info.project_path || 'Not set'}</span>
            </div>
            <div class="project-info-item">
                <span class="project-info-label">Project Name:</span>
                <span class="project-info-value">${info.project_name || 'Unknown'}</span>
            </div>
            ${packagesHtml}
        `;
    } catch (error) {
        console.error('Error loading project info:', error);
        document.getElementById('current-project-info').innerHTML = 
            '<div class="loading">Error loading project information</div>';
    }
}

// --- COVERAGE FUNCTIONS (unchanged) ---

function openHTMLCoverage(filename) {
    if (!filename) {
        alert('HTML coverage report not available');
        return;
    }
    
    // Open the HTML coverage report in a new tab/window
    const url = `/html-coverage/${encodeURIComponent(filename)}`;
    window.open(url, '_blank', 'width=1200,height=800,scrollbars=yes,resizable=yes');
}

function showCoverage(packageName) {
    document.getElementById('coverage-package-name').textContent = `${packageName} Coverage`;
    document.getElementById('coverage-modal').classList.add('show');
    document.body.style.overflow = 'hidden';

    fetch(`/coverage/${encodeURIComponent(packageName)}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch coverage data');
            }
            return response.json();
        })
        .then(files => {
            displayCoverageFiles(files);
        })
        .catch(error => {
            console.error('Error fetching coverage:', error);
            document.getElementById('file-list').innerHTML = '<div class="file-item">Error loading coverage data</div>';
            document.getElementById('source-code').innerHTML = '<div class="loading">Could not load coverage.</div>';
        });
}

function displayCoverageFiles(files) {
    const fileList = document.getElementById('file-list');
    fileList.innerHTML = '';
    document.getElementById('source-code').innerHTML = '<div class="loading">Select a file to view coverage...</div>';

    if (!files || files.length === 0) {
        fileList.innerHTML = '<div class="file-item">No coverage data available</div>';
        return;
    }

    files.forEach((file, index) => {
        const fileItem = document.createElement('div');
        fileItem.className = 'file-item';
        fileItem.onclick = () => selectFile(file, fileItem);

        fileItem.innerHTML = `
            <div class="file-name">${escapeHtml(getFileName(file.filename))}</div>
            <div class="file-coverage">${file.coverage.toFixed(1)}% coverage</div>`;

        fileList.appendChild(fileItem);
    });

    if (files.length > 0) {
        fileList.firstChild.click();
    }
}

function selectFile(file, fileItem) {
    document.querySelectorAll('.file-item.active').forEach(item => item.classList.remove('active'));
    fileItem.classList.add('active');
    displaySourceCode(file);
}

function highlightGoSyntax(code) {
    return code
        .replace(/\b(package|import|func|var|const|type|struct|interface|if|else|for|range|switch|case|default|return|break|continue|go|defer|select|chan|map)\b/g, '<span class="go-keyword">$1</span>')
        .replace(/"([^"\\\\]|\\\\.)*"/g, '<span class="go-string">$&</span>')
        .replace(/'([^'\\\\]|\\\\.)*'/g, '<span class="go-string">$&</span>')
        .replace(/`([^`]*)`/g, '<span class="go-string">$&</span>')
        .replace(/\/\/.*$/gm, '<span class="go-comment">$&</span>')
        .replace(/\/\*[\s\S]*?\*\//g, '<span class="go-comment">$&</span>')
        .replace(/\b\d+(\.\d+)?\b/g, '<span class="go-number">$&</span>')
        .replace(/\b([A-Z][a-zA-Z0-9_]*)(?=\s*\()/g, '<span class="go-function">$1</span>');
}

function displaySourceCode(file) {
    const sourceCode = document.getElementById('source-code');
    const lines = file.content.split('\n');
    const lineCoverage = {};

    if (file.blocks) {
        file.blocks.forEach(block => {
            for (let i = block.start_line; i <= block.end_line; i++) {
                lineCoverage[i] = block.covered;
            }
        });
    }

    sourceCode.innerHTML = lines.map((line, index) => {
        const lineNumber = index + 1;
        const covered = lineCoverage[lineNumber];
        let cssClass = '';

        if (covered === true) {
            cssClass = 'covered';
        } else if (covered === false && line.trim() !== '') {
            cssClass = 'uncovered';
        }

        const highlightedLine = highlightGoSyntax(escapeHtml(line));
        return `
            <div class="code-line ${cssClass}">
                <div class="line-number">${lineNumber}</div>
                <div class="line-content">${highlightedLine || '&nbsp;'}</div>
            </div>`;
    }).join('');
}

function getFileName(fullPath) {
    return fullPath.split('/').pop();
}

function closeCoverage() {
    document.getElementById('coverage-modal').classList.remove('show');
    document.body.style.overflow = '';
}

function togglePackage(header) {
    const details = header.nextElementSibling;
    details.classList.toggle('expanded');
}

function runTests() {
    const resultsEl = document.getElementById('results');
    // Set a "running" state on the results element itself.
    resultsEl.innerHTML = '<div class="loading">Running tests...</div>';
    resultsEl.dataset.isRunning = "true";

    fetch('/run-tests', { method: 'POST' })
        .catch(error => {
            console.error('Error running tests:', error);
            resultsEl.innerHTML = '<div class="loading">Failed to start tests.</div>';
            // Clean up the state attribute on failure.
            delete resultsEl.dataset.isRunning;
        });
}

// --- INITIALIZATION ---

connectWebSocket();

document.addEventListener('DOMContentLoaded', () => {
    const runButton = document.getElementById('run-button');
    const projectButton = document.getElementById('project-button');
    const setPathButton = document.getElementById('set-path-button');
    const manualPathInput = document.getElementById('manual-path-input');
    const projectModal = document.getElementById('project-modal');
    const coverageModal = document.getElementById('coverage-modal');

    // Button event listeners
    runButton.addEventListener('click', runTests);
    projectButton.addEventListener('click', showProjectModal);
    setPathButton.addEventListener('click', setProjectPath);

    // Enter key support for manual path input
    manualPathInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            setProjectPath();
        }
    });

    // Modal close handlers
    projectModal.addEventListener('click', (event) => {
        if (event.target === projectModal) {
            closeProjectModal();
        }
    });

    coverageModal.addEventListener('click', (event) => {
        if (event.target === coverageModal) {
            closeCoverage();
        }
    });

    // Escape key support
    document.addEventListener('keydown', (event) => {
        if (event.key === 'Escape') {
            if (projectModal.classList.contains('show')) {
                closeProjectModal();
            }
            if (coverageModal.classList.contains('show')) {
                closeCoverage();
            }
        }
    });

    // Auto-run tests shortly after page load
    setTimeout(runTests, 1000);
});