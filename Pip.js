(function() {
  
    const menu = document.createElement('div');
    menu.innerHTML = `
        <div class="control-panel">
            <div class="header">
                <h3 class="title">🎬 Video Control</h3>
                <div class="controls">
                    <button class="pip-toggle">
                        <span class="pip-icon">
                            <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
                                <path d="M8 5v14l11-7z"/>
                            </svg>
                        </span>
                        <span class="pip-text">Picture-in-Picture</span>
                    </button>
                    <button class="close">
                        <svg class="close-icon" viewBox="0 0 24 24" width="18" height="18">
                            <path fill="currentColor" d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
                        </svg>
                    </button>
                </div>
            </div>
        </div>
    `;

    menu.innerHTML += `
        <style>
            .control-panel {
                position: fixed;
                bottom: 25px;
                left: 25px;
                background: rgba(30, 30, 30, 0.95);
                border-radius: 16px;
                box-shadow: 0 8px 32px rgba(0,0,0,0.3);
                width: 280px;
                font-family: 'Segoe UI', system-ui;
                z-index: 100000;
                color: #fff;
                backdrop-filter: blur(12px);
                border: 1px solid rgba(255,255,255,0.1);
                transform: translateY(20px);
                animation: slideUp 0.3s ease-out forwards;
            }
            
            @keyframes slideUp {
                to { transform: translateY(0); opacity: 1; }
                from { transform: translateY(20px); opacity: 0; }
            }
            
            .header {
                padding: 18px;
                display: flex;
                justify-content: space-between;
                align-items: center;
            }
            
            .title {
                margin: 0;
                font-size: 1.1em;
                font-weight: 600;
                letter-spacing: -0.5px;
                background: linear-gradient(45deg, #00b4d8, #90e0ef);
                -webkit-background-clip: text;
                -webkit-text-fill-color: transparent;
            }
            
            .controls {
                display: flex;
                gap: 8px;
                align-items: center;
            }
            
            .pip-toggle {
                padding: 10px 16px;
                border: none;
                border-radius: 10px;
                cursor: pointer;
                background: linear-gradient(135deg, #4361ee, #3a0ca3);
                color: white;
                display: flex;
                align-items: center;
                gap: 8px;
                transition: all 0.2s ease;
                box-shadow: 0 4px 15px rgba(67, 97, 238, 0.3);
            }
            
            .pip-toggle:hover {
                transform: translateY(-1px);
                box-shadow: 0 6px 20px rgba(67, 97, 238, 0.4);
                background: linear-gradient(135deg, #3a0ca3, #4361ee);
            }
            
            .pip-toggle:active {
                transform: translateY(0);
                filter: brightness(0.9);
            }
            
            .pip-icon {
                display: flex;
                align-items: center;
                justify-content: center;
                width: 18px;
                height: 18px;
            }
            
            .pip-text {
                font-weight: 500;
                font-size: 0.95em;
            }
            
            .close {
                padding: 8px;
                border: none;
                border-radius: 50%;
                background: rgba(255, 68, 68, 0.15);
                color: #ff4444;
                cursor: pointer;
                transition: all 0.2s ease;
                display: flex;
                align-items: center;
                justify-content: center;
            }
            
            .close:hover {
                background: rgba(255, 68, 68, 0.25);
                transform: rotate(90deg);
            }
            
            .close-icon {
                transition: transform 0.3s ease;
            }
        </style>
    `;

    const pipToggle = menu.querySelector('.pip-toggle');
    let videoElement = null;

    function togglePictureInPicture() {
        if (!videoElement) {
            videoElement = document.querySelector('video');
            if (!videoElement) {
                alert('No video found on the page 🎥');
                return;
            }
        }

        if (document.pictureInPictureElement) {
            document.exitPictureInPicture()
                .catch(error => console.error("PiP Error:", error));
        } else {
            videoElement.requestPictureInPicture()
                .catch(error => console.error("PiP Error:", error));
        }
    }

    pipToggle.addEventListener('click', togglePictureInPicture);
    menu.querySelector('.close').addEventListener('click', () => document.body.removeChild(menu));
    document.body.appendChild(menu);
})();
