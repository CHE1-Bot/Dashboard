// dashboard.js

document.addEventListener('DOMContentLoaded', function() {
    const searchBtn = document.getElementById('search-btn');
    const ticketIdInput = document.getElementById('ticket-id');
    const transcriptsList = document.getElementById('transcripts-list');
    const dropdownButton = document.querySelector('.dropdown-button');
    const dropdownContent = document.querySelector('.dropdown-content');

    // Load all transcripts on page load
    loadAllTranscripts();

    // Dropdown toggle
    dropdownButton.addEventListener('click', function() {
        dropdownContent.parentElement.classList.toggle('active');
    });

    // Sorting options
    dropdownContent.addEventListener('click', function(e) {
        if (e.target.tagName === 'A') {
            e.preventDefault();
            const sortType = e.target.getAttribute('data-sort');
            sortTranscripts(sortType);
            dropdownContent.parentElement.classList.remove('active');
        }
    });

    // Search functionality
    searchBtn.addEventListener('click', function() {
        const ticketId = ticketIdInput.value.trim();
        if (ticketId) {
            searchTranscript(ticketId);
        } else {
            loadAllTranscripts();
        }
    });

    ticketIdInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            searchBtn.click();
        }
    });

    function loadAllTranscripts() {
        fetch('/api/transcripts')
            .then(response => response.json())
            .then(data => {
                displayTranscripts(data);
            })
            .catch(error => {
                console.error('Error loading transcripts:', error);
                transcriptsList.innerHTML = '<p>Error loading transcripts.</p>';
            });
    }

    function searchTranscript(id) {
        fetch(`/api/transcripts/${id}`)
            .then(response => {
                if (response.ok) {
                    return response.json();
                } else {
                    throw new Error('Not found');
                }
            })
            .then(data => {
                displayTranscripts([data]);
            })
            .catch(error => {
                console.error('Error fetching transcript:', error);
                transcriptsList.innerHTML = '<p>No transcript found for this ID.</p>';
            });
    }

    function displayTranscripts(transcripts) {
        if (transcripts && transcripts.length > 0) {
            transcriptsList.innerHTML = transcripts.map(t => `
                <div class="transcript-item">
                    <h3>Ticket ID: ${t.id}</h3>
                    <div class="transcript-content">${t.content}</div>
                </div>
            `).join('');
        } else {
            transcriptsList.innerHTML = '<p>No transcripts found.</p>';
        }
    }

    function sortTranscripts(sortType) {
        // Implement sorting logic
        // For now, just reload
        loadAllTranscripts();
    }
});