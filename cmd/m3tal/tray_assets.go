package main

// TrayHTML is the complete HTML page for the system tray popup
const TrayHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>M3TAL System Tray</title>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&family=Rajdhani:wght@400;500;600;700&display=swap');
        
        :root {
            --bg: #080b10;
            --bg-card: #0d1117;
            --border: rgba(0, 212, 170, 0.12);
            --border-dim: rgba(255,255,255,0.06);
            --teal: #00d4aa;
            --teal-dim: rgba(0, 212, 170, 0.15);
            --green: #22c55e;
            --green-dim: rgba(34, 197, 94, 0.15);
            --amber: #f59e0b;
            --amber-dim: rgba(245, 158, 11, 0.15);
            --red: #ef4444;
            --red-dim: rgba(239, 68, 68, 0.15);
            --text-1: #e2e8f0;
            --text-2: #94a3b8;
            --text-3: #4b5e75;
            --font-ui: 'Rajdhani', -apple-system, sans-serif;
            --font-mono: 'JetBrains Mono', monospace;
        }

        body {
            font-family: var(--font-ui);
            background: var(--bg);
            color: var(--text-1);
            margin: 0;
            padding: 1.25rem;
            min-height: 100vh;
            box-sizing: border-box;
            display: flex;
            flex-direction: column;
            gap: 1.25rem;
        }

        .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            border-bottom: 1px solid var(--border);
            padding-bottom: 0.5rem;
        }

        .title {
            font-size: 1.3rem;
            font-weight: 700;
            letter-spacing: 2px;
            text-transform: uppercase;
            color: var(--teal);
            margin: 0;
        }

        .grid {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 1rem;
        }

        @media (max-width: 600px) {
            .grid {
                grid-template-columns: 1fr;
            }
        }

        .card {
            background: var(--bg-card);
            border: 1px solid var(--border-dim);
            border-radius: 8px;
            padding: 1rem;
            display: flex;
            flex-direction: column;
            gap: 0.5rem;
            position: relative;
            overflow: hidden;
            transition: border-color 0.2s;
        }

        .card:hover {
            border-color: var(--border);
        }

        .card::after {
            content: '';
            position: absolute;
            bottom: 0; left: 0; right: 0;
            height: 2px;
            background: var(--teal);
            opacity: 0.3;
        }

        .card-header {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-weight: 600;
            font-size: 0.85rem;
            letter-spacing: 1px;
            text-transform: uppercase;
            color: var(--text-2);
        }

        .card-icon {
            font-size: 1.2rem;
        }

        .metric-value {
            font-family: var(--font-mono);
            font-size: 2.2rem;
            font-weight: 700;
            color: var(--text-1);
        }

        .metric-sub {
            font-size: 0.85rem;
            color: var(--text-2);
        }

        .progress-track {
            width: 100%;
            height: 6px;
            background: rgba(255,255,255,0.05);
            border-radius: 3px;
            overflow: hidden;
            margin-top: 4px;
        }

        .progress-fill {
            height: 100%;
            background: var(--teal);
            width: 0%;
            transition: width 0.3s ease;
        }

        .progress-fill.low { background: var(--green); }
        .progress-fill.mid { background: var(--amber); }
        .progress-fill.high { background: var(--red); }

        .container-list {
            max-height: 180px;
            overflow-y: auto;
            display: flex;
            flex-direction: column;
            gap: 0.5rem;
            padding-right: 0.25rem;
        }

        /* Scrollbar styles */
        .container-list::-webkit-scrollbar {
            width: 4px;
        }
        .container-list::-webkit-scrollbar-track {
            background: rgba(255,255,255,0.02);
        }
        .container-list::-webkit-scrollbar-thumb {
            background: var(--border);
            border-radius: 2px;
        }

        .container-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-family: var(--font-mono);
            font-size: 0.85rem;
            padding: 0.35rem 0.6rem;
            background: rgba(0,0,0,0.15);
            border: 1px solid var(--border-dim);
            border-radius: 4px;
        }

        .badge {
            font-size: 0.7rem;
            font-weight: 700;
            padding: 2px 6px;
            border-radius: 3px;
            text-transform: uppercase;
        }
        .badge.running { background: var(--green-dim); color: var(--green); border: 1px solid rgba(34, 197, 94, 0.2); }
        .badge.stopped { background: var(--red-dim); color: var(--red); border: 1px solid rgba(239, 68, 68, 0.2); }

        .storage-item {
            display: flex;
            justify-content: space-between;
            font-size: 0.9rem;
            font-family: var(--font-mono);
        }
    </style>
</head>
<body>
    <div class="header">
        <h1 class="title">🛠️ M3TAL System Monitor</h1>
    </div>

    <div class="grid">
        <!-- CPU Card -->
        <div class="card">
            <div class="card-header">
                <span class="card-icon">⚙️</span>
                <span>CPU Engine</span>
            </div>
            <div class="metric-value" id="cpu-val">0.0%</div>
            <div class="metric-sub" id="cpu-sub">Temp: --°C</div>
            <div class="progress-track">
                <div class="progress-fill" id="cpu-bar"></div>
            </div>
        </div>

        <!-- GPU Card -->
        <div class="card">
            <div class="card-header">
                <span class="card-icon">🎮</span>
                <span>GPU Engine</span>
            </div>
            <div class="metric-value" id="gpu-val">--%</div>
            <div class="metric-sub" id="gpu-sub">Temp: --°C | VRAM: -- / -- MB</div>
            <div class="progress-track">
                <div class="progress-fill" id="gpu-bar"></div>
            </div>
        </div>

        <!-- Storage Card -->
        <div class="card" style="grid-column: span 2;">
            <div class="card-header">
                <span class="card-icon">💽</span>
                <span>Storage Cards</span>
            </div>
            <div class="storage-item">
                <span style="color: var(--text-2);">Root Drive (/)</span>
                <span id="storage-text">-- / -- GB</span>
            </div>
            <div class="progress-track">
                <div class="progress-fill" id="storage-bar"></div>
            </div>
        </div>

        <!-- Containers Card -->
        <div class="card" style="grid-column: span 2;">
            <div class="card-header">
                <span class="card-icon">📦</span>
                <span id="container-title">Containers: 0 UP</span>
            </div>
            <div class="container-list" id="container-list">
                <div class="container-item">
                    <span>Loading containers...</span>
                </div>
            </div>
        </div>
    </div>

    <script>
        function getProgressClass(val) {
            if (val < 50) return 'low';
            if (val < 85) return 'mid';
            return 'high';
        }

        async function updateStats() {
            try {
                const response = await fetch('/tray/api/stats');
                const stats = await response.json();

                // CPU
                document.getElementById('cpu-val').innerText = stats.cpu_usage.toFixed(1) + '%';
                document.getElementById('cpu-sub').innerText = 'Temp: ' + (stats.cpu_temp > 0 ? stats.cpu_temp.toFixed(0) + '°C' : '--°C') + ' | RAM: ' + stats.memory_used.toFixed(1) + ' / ' + stats.memory_total.toFixed(0) + ' GB';
                const cpuBar = document.getElementById('cpu-bar');
                cpuBar.style.width = stats.cpu_usage + '%';
                cpuBar.className = 'progress-fill ' + getProgressClass(stats.cpu_usage);

                // GPU
                if (stats.gpu_model && stats.gpu_model !== "No GPU Detected") {
                    document.getElementById('gpu-val').innerText = stats.gpu_usage.toFixed(1) + '%';
                    document.getElementById('gpu-sub').innerText = 'Temp: ' + stats.gpu_temp.toFixed(0) + '°C | VRAM: ' + stats.gpu_mem_used.toFixed(0) + ' / ' + stats.gpu_mem_total.toFixed(0) + ' MB (' + stats.gpu_model + ')';
                    const gpuBar = document.getElementById('gpu-bar');
                    gpuBar.style.width = stats.gpu_usage + '%';
                    gpuBar.className = 'progress-fill ' + getProgressClass(stats.gpu_usage);
                } else {
                    document.getElementById('gpu-val').innerText = 'N/A';
                    document.getElementById('gpu-sub').innerText = 'No GPU detected';
                    document.getElementById('gpu-bar').style.width = '0%';
                }

                // Storage
                document.getElementById('storage-text').innerText = stats.disk_used.toFixed(1) + ' / ' + stats.disk_total.toFixed(0) + ' GB (' + stats.disk_usage.toFixed(1) + '%)';
                const storageBar = document.getElementById('storage-bar');
                storageBar.style.width = stats.disk_usage + '%';
                storageBar.className = 'progress-fill ' + getProgressClass(stats.disk_usage);

            } catch (e) {
                console.error("Failed to fetch stats:", e);
            }
        }

        async function updateContainers() {
            try {
                const response = await fetch('/tray/api/containers');
                const list = await response.json();

                const running = list.filter(c => c.state === 'running').length;
                document.getElementById('container-title').innerText = 'Containers: ' + running + ' / ' + list.length + ' UP';

                const listEl = document.getElementById('container-list');
                listEl.innerHTML = '';
                if (list.length === 0) {
                    listEl.innerHTML = '<div class="container-item"><span>No containers found.</span></div>';
                    return;
                }

                list.forEach(c => {
                    const item = document.createElement('div');
                    item.className = 'container-item';
                    
                    const nameSpan = document.createElement('span');
                    nameSpan.innerText = c.names && c.names.length > 0 ? c.names[0].replace('/', '') : c.image;
                    
                    const badge = document.createElement('span');
                    badge.className = 'badge ' + (c.state === 'running' ? 'running' : 'stopped');
                    badge.innerText = c.state;

                    item.appendChild(nameSpan);
                    item.appendChild(badge);
                    listEl.appendChild(item);
                });
            } catch (e) {
                console.error("Failed to fetch containers:", e);
            }
        }

        // Poll every 2 seconds
        setInterval(updateStats, 2000);
        setInterval(updateContainers, 2000);

        // Initial call
        updateStats();
        updateContainers();
    </script>
</body>
</html>`
