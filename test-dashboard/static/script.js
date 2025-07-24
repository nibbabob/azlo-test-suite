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
    if (data.results && data.results.length > 0) {
        resultsEl.innerHTML = data.results.map(result => createPackageHTML(result)).join('');
    } else if (data.results) {
        resultsEl.innerHTML = '<div class="loading">No test results yet. Click "Run Tests".</div>';
    }

    // Update last run time
    if (data.last_run) {
        const lastRunEl = document.getElementById('last-run');
        const lastRun = new Date(data.last_run);
        lastRunEl.textContent = `Last run: ${lastRun.toLocaleTimeString()}`;
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

    const coverageButton = result.files && result.files.length > 0 ?
        `<button class="coverage-link" onclick="event.stopPropagation(); showCoverage('${escapeHtml(result.package || '')}')">View Coverage</button>` : '';

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
                ${coverageButton}
            </div>
        </div>`;
}

function showCoverage(packageName) {
    document.getElementById('coverage-package-name').textContent = `${packageName} Coverage`;
    document.getElementById('coverage-modal').classList.add('show');
    document.body.style.overflow = 'hidden'; // Prevent background scrolling

    // Fetch coverage data
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

    // Select first file by default
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
    // Note: escapeHtml should be called *before* this function.
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
    document.body.style.overflow = ''; // Re-enable scrolling
}

function togglePackage(header) {
    const details = header.nextElementSibling;
    details.classList.toggle('expanded');
}

function runTests() {
    fetch('/run-tests', { method: 'POST' })
        .then(() => {
            document.getElementById('results').innerHTML = '<div class="loading">Running tests...</div>';
        })
        .catch(error => {
            console.error('Error running tests:', error);
            document.getElementById('results').innerHTML = '<div class="loading">Failed to start tests.</div>';
        });
}

// --- INITIALIZATION ---

// Connect WebSocket on page load
connectWebSocket();

// Set up event listeners once the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    const runButton = document.getElementById('run-button');
    const modal = document.getElementById('coverage-modal');

    // Run tests button
    runButton.addEventListener('click', runTests);

    // Close modal when clicking outside the content
    modal.addEventListener('click', function(event) {
        if (event.target === modal) {
            closeCoverage();
        }
    });

    // Close modal with Escape key
    document.addEventListener('keydown', function(event) {
        if (event.key === 'Escape') {
            closeCoverage();
        }
    });

    // Auto-run tests shortly after page load
    setTimeout(runTests, 1000);
});