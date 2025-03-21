document.addEventListener('DOMContentLoaded', function() {
    const analyzeForm = document.getElementById('analyzeForm');
    const hashForm = document.getElementById('hashForm');
    const resultsDiv = document.getElementById('results');

    // Show loading indicator
    function showLoading() {
        resultsDiv.innerHTML = '<div class="loading">Processing...</div>';
    }

    // Show error message
    function showError(message) {
        resultsDiv.innerHTML = `<div class="error-message">${message}</div>`;
    }

    // Show results
    function showResults(data) {
        let html = '';
        if (Array.isArray(data)) {
            data.forEach(item => {
                html += `
                    <div class="result-item">
                        <h6>${item.file}</h6>
                        <p>Lines: ${item.lines} | Words: ${item.words} | Bytes: ${item.bytes}</p>
                        <small>Duration: ${item.duration}</small>
                    </div>
                `;
            });
        } else {
            html = `
                <div class="result-item">
                    <h6>Hash Result</h6>
                    <p>${data.hash}</p>
                </div>
            `;
        }
        resultsDiv.innerHTML = html;
    }

    // Handle analyze form submission
    analyzeForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        const path = document.getElementById('path').value;
        
        showLoading();
        try {
            const response = await fetch('/api/v1/analyze', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ path }),
            });

            if (!response.ok) {
                throw new Error('Analysis failed');
            }

            const data = await response.json();
            showResults(data);
        } catch (error) {
            showError(error.message);
        }
    });

    // Handle hash form submission
    hashForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        const file = document.getElementById('file').value;
        
        showLoading();
        try {
            const response = await fetch('/api/v1/hash', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ file }),
            });

            if (!response.ok) {
                throw new Error('Hash calculation failed');
            }

            const data = await response.json();
            showResults(data);
        } catch (error) {
            showError(error.message);
        }
    });
}); 