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
    updateProjectInfo(data);

    const overallCoverage = data.overall_coverage || 0;
    document.getElementById('overall-coverage').textContent = `${overallCoverage.toFixed(1)}%`;
    document.getElementById('total-tests').textContent = data.total_tests || 0;
    document.getElementById('passed-tests').textContent = data.passed_tests || 0;
    document.getElementById('failed-tests').textContent = (data.total_tests || 0) - (data.passed_tests || 0);

    const coverageEl = document.getElementById('overall-coverage');
    coverageEl.className = `stat-number ${getCoverageClass(overallCoverage)}`;

    const resultsEl = document.getElementById('results');
    if (data.results && data.results.length > 0) {
        resultsEl.innerHTML = data.results.map(result => createPackageHTML(result)).join('');
        delete resultsEl.dataset.isRunning;
    } else if (data.results && resultsEl.dataset.isRunning === "true") {
        resultsEl.innerHTML = '';
    } else if (data.results) {
        resultsEl.innerHTML = '<div class="loading">No test results yet. Click "Run Tests".</div>';
    }

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
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ path: path })
        });
        const result = await response.json();
        if (result.success) {
            closeProjectModal();
            setTimeout(() => alert('Project path updated successfully!'), 100);
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
                </div>`;
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
            ${packagesHtml}`;
    } catch (error) {
        console.error('Error loading project info:', error);
        document.getElementById('current-project-info').innerHTML = '<div class="loading">Error loading project information</div>';
    }
}

function openHTMLCoverage(filename) {
    if (!filename) {
        alert('HTML coverage report not available');
        return;
    }
    const url = `/html-coverage/${encodeURIComponent(filename)}`;
    window.open(url, '_blank', 'width=1200,height=800,scrollbars=yes,resizable=yes');
}

function showCoverage(packageName) {
    document.getElementById('coverage-package-name').textContent = `${packageName} Coverage`;
    document.getElementById('coverage-modal').classList.add('show');
    document.body.style.overflow = 'hidden';

    fetch(`/coverage/${encodeURIComponent(packageName)}`)
        .then(response => response.ok ? response.json() : Promise.reject('Failed to fetch coverage data'))
        .then(files => displayCoverageFiles(files))
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

    files.forEach((file) => {
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

// ===================================================================================
// START OF THE NEW, ROBUST SYNTAX HIGHLIGHTING LOGIC
// ===================================================================================

/**
 * Tokenizes a line of Go code into an array of objects.
 * This is more robust than using multiple regex replacements.
 * @param {string} rawCode - The raw line of code.
 * @returns {Array<{type: string, content: string}>} - An array of tokens.
 */
function tokenizeGoCode(rawCode) {
    const tokens = [];
    let remainingCode = rawCode;

    // The order of these rules is critical for correct parsing.
    const tokenRules = [
        { type: 'string', regex: /^"([^"\\\\]|\\\\.)*"/ },
        { type: 'string', regex: /^`[^`]*`/ },
        { type: 'comment', regex: /^\/\/.*$/ },
        { type: 'comment', regex: /^\/\*[\s\S]*?\*\// },
        { type: 'keyword', regex: /^\b(package|import|func|var|const|type|struct|interface|if|else|for|range|switch|case|default|return|break|continue|go|defer|select|chan|map)\b/ },
        { type: 'number', regex: /^\b\d+(\.\d+)?\b/ },
        { type: 'function', regex: /^\b([A-Z][a-zA-Z0-9_]*)(?=\s*\()/ },
        { type: 'whitespace', regex: /^\s+/ },
        { type: 'other', regex: /^./ } // Fallback for any other single character
    ];

    while (remainingCode.length > 0) {
        let matched = false;
        for (const rule of tokenRules) {
            const match = remainingCode.match(rule.regex);
            if (match) {
                tokens.push({ type: rule.type, content: match[0] });
                remainingCode = remainingCode.substring(match[0].length);
                matched = true;
                break;
            }
        }
        if (!matched) { // Safety break to prevent infinite loops
            tokens.push({ type: 'other', content: remainingCode[0] });
            remainingCode = remainingCode.substring(1);
        }
    }
    return tokens;
}

/**
 * Renders the source code view by manually creating DOM elements.
 * This avoids innerHTML parsing issues and is much more reliable.
 * @param {object} file - The file object containing the code.
 */
function displaySourceCode(file) {
    const sourceCode = document.getElementById('source-code');
    sourceCode.innerHTML = ''; // Clear previous content to prevent memory leaks.

    const lines = file.content.split('\n');
    const lineCoverage = {};

    if (file.blocks) {
        file.blocks.forEach(block => {
            for (let i = block.start_line; i <= block.end_line; i++) {
                lineCoverage[i] = block.covered;
            }
        });
    }

    lines.forEach((line, index) => {
        const lineNumber = index + 1;
        const covered = lineCoverage[lineNumber];

        const lineDiv = document.createElement('div');
        lineDiv.className = 'code-line';
        if (covered === true) lineDiv.classList.add('covered');
        else if (covered === false && line.trim() !== '') lineDiv.classList.add('uncovered');

        const numberDiv = document.createElement('div');
        numberDiv.className = 'line-number';
        numberDiv.textContent = lineNumber;

        const contentDiv = document.createElement('div');
        contentDiv.className = 'line-content';

        if (line.trim() === '') {
            contentDiv.innerHTML = '&nbsp;'; // Use innerHTML for non-breaking space
        } else {
            const tokens = tokenizeGoCode(line);
            tokens.forEach(token => {
                // For plain text, create a text node to ensure it's not parsed as HTML.
                if (token.type === 'other' || token.type === 'whitespace') {
                    contentDiv.appendChild(document.createTextNode(token.content));
                } else {
                    // For highlighted tokens, create a styled span element.
                    const span = document.createElement('span');
                    span.className = `go-${token.type}`; // e.g., 'go-string', 'go-keyword'
                    span.textContent = token.content;
                    contentDiv.appendChild(span);
                }
            });
        }

        lineDiv.appendChild(numberDiv);
        lineDiv.appendChild(contentDiv);
        sourceCode.appendChild(lineDiv);
    });
}

// ===================================================================================
// END OF THE NEW LOGIC
// ===================================================================================

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
    resultsEl.innerHTML = '<div class="loading">Running tests...</div>';
    resultsEl.dataset.isRunning = "true";

    fetch('/run-tests', { method: 'POST' })
        .catch(error => {
            console.error('Error running tests:', error);
            resultsEl.innerHTML = '<div class="loading">Failed to start tests.</div>';
            delete resultsEl.dataset.isRunning;
        });
}

connectWebSocket();

document.addEventListener('DOMContentLoaded', () => {
    const runButton = document.getElementById('run-button');
    const projectButton = document.getElementById('project-button');
    const setPathButton = document.getElementById('set-path-button');
    const manualPathInput = document.getElementById('manual-path-input');
    const projectModal = document.getElementById('project-modal');
    const coverageModal = document.getElementById('coverage-modal');

    runButton.addEventListener('click', runTests);
    projectButton.addEventListener('click', showProjectModal);
    setPathButton.addEventListener('click', setProjectPath);

    manualPathInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') setProjectPath();
    });

    projectModal.addEventListener('click', (event) => {
        if (event.target === projectModal) closeProjectModal();
    });

    coverageModal.addEventListener('click', (event) => {
        if (event.target === coverageModal) closeCoverage();
    });

    document.addEventListener('keydown', (event) => {
        if (event.key === 'Escape') {
            if (projectModal.classList.contains('show')) closeProjectModal();
            if (coverageModal.classList.contains('show')) closeCoverage();
        }
    });

    setTimeout(runTests, 1000);
});