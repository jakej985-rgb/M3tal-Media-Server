/* app.js — M3TAL Dashboard v2 */

const socket = io();

// ── Clock ────────────────────────────────────────────────────────
function tick() {
    const el = document.getElementById('live-clock');
    if (!el) return;
    const now = new Date();
    const hh = String(now.getHours()).padStart(2, '0');
    const mm = String(now.getMinutes()).padStart(2, '0');
    const ss = String(now.getSeconds()).padStart(2, '0');
    el.textContent = `${hh}:${mm}:${ss}`;
}
setInterval(tick, 1000);
tick();

// ── Resource Chart ───────────────────────────────────────────────
let chart = null;
const MAX_POINTS = 30;
const cpuData  = Array(MAX_POINTS).fill(null);
const memData  = Array(MAX_POINTS).fill(null);
const timeLabels = Array(MAX_POINTS).fill('');

function initChart() {
    const canvas = document.getElementById('resource-chart');
    if (!canvas) return;

    chart = new Chart(canvas, {
        type: 'line',
        data: {
            labels: timeLabels,
            datasets: [
                {
                    label: 'CPU',
                    data: cpuData,
                    borderColor: '#22c55e',
                    backgroundColor: 'rgba(34,197,94,0.08)',
                    borderWidth: 1.5,
                    tension: 0.4,
                    fill: true,
                    pointRadius: 0,
                },
                {
                    label: 'MEM',
                    data: memData,
                    borderColor: '#a855f7',
                    backgroundColor: 'rgba(168,85,247,0.08)',
                    borderWidth: 1.5,
                    tension: 0.4,
                    fill: true,
                    pointRadius: 0,
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            animation: { duration: 400 },
            interaction: { mode: 'index', intersect: false },
            scales: {
                x: {
                    display: false,
                    grid: { display: false }
                },
                y: {
                    min: 0, max: 100,
                    grid: { color: 'rgba(255,255,255,0.04)', drawBorder: false },
                    ticks: {
                        color: '#4b5e75',
                        font: { family: "'JetBrains Mono', monospace", size: 10 },
                        callback: v => `${v}%`,
                        maxTicksLimit: 5,
                    },
                    border: { display: false }
                }
            },
            plugins: {
                legend: { display: false },
                tooltip: {
                    backgroundColor: 'rgba(13,17,23,0.9)',
                    borderColor: 'rgba(0,212,170,0.2)',
                    borderWidth: 1,
                    titleColor: '#94a3b8',
                    bodyColor: '#e2e8f0',
                    bodyFont: { family: "'JetBrains Mono', monospace", size: 11 },
                }
            }
        }
    });
}

function pushChartPoint(cpu, mem) {
    if (!chart) return;
    const now = new Date();
    const label = `${String(now.getHours()).padStart(2,'0')}:${String(now.getMinutes()).padStart(2,'0')}`;

    cpuData.push(cpu);
    memData.push(mem);
    timeLabels.push(label);

    if (cpuData.length > MAX_POINTS)   { cpuData.shift(); }
    if (memData.length > MAX_POINTS)   { memData.shift(); }
    if (timeLabels.length > MAX_POINTS){ timeLabels.shift(); }

    chart.update('none');
}

// ── Socket – real-time metrics ────────────────────────────────────
socket.on('metrics_update', (data) => {
    const sys = data.system || {};
    const cpu = sys.cpu || 0;
    const mem = sys.mem || 0;

    // Stat cards
    setText('stat-cpu', `${cpu.toFixed(1)}%`);
    setText('stat-mem', `${(sys.mem_gb || 0).toFixed(1)} GB`);

    // Push to chart
    pushChartPoint(cpu, mem);
});

// ── Helpers ───────────────────────────────────────────────────────
function setText(id, val) {
    const el = document.getElementById(id);
    if (el) el.textContent = val;
}

function getStatusClass(status) {
    const s = (status || '').toLowerCase();
    if (s === 'running' || s === 'online') return 'running';
    if (s === 'restarting') return 'restarting';
    if (s === 'offline' || s === 'exited') return 'offline';
    if (s === 'missing') return 'missing';
    return 'unknown';
}

function getCpuClass(cpu) {
    if (cpu >= 80) return 'cpu-crit';
    if (cpu >= 50) return 'cpu-high';
    return '';
}

// ── UI Interactions ───────────────────────────────────────────────
function togglePanel(header) {
    const panel = header.parentElement;
    panel.classList.toggle('collapsed');
}

function toggleRow(rowId) {
    const detailsRow = document.getElementById(rowId);
    if (!detailsRow) return;
    
    const isVisible = detailsRow.style.display !== 'none';
    
    // Hide all other details rows first (optional, for accordion effect)
    // document.querySelectorAll('.details-row').forEach(r => r.style.display = 'none');
    
    detailsRow.style.display = isVisible ? 'none' : 'table-row';
}

// ── Health score ──────────────────────────────────────────────────
async function refreshHealth() {
    try {
        const res  = await fetch('/api/health/report');
        const data = await res.json();
        const score = data.score || 0;
        const verdict = data.verdict || 'Healthy';

        // Main Score
        const scoreEl = document.getElementById('health-score');
        if (scoreEl) scoreEl.textContent = score;
        
        // Mini Card Score (Standardized ID)
        const healthVal = document.getElementById('stat-health-val');
        if (healthVal) {
            healthVal.textContent = `${score}%`;
        }
        
        const ring = document.getElementById('health-ring');
        if (ring) {
            const offset = 220 - (220 * score / 100);
            ring.style.strokeDashoffset = offset;
        }

        // Mini Card Score
        setText('stat-health-val', `${score}%`);
        
        const ringMini = document.getElementById('gsi-ring-mini');
        if (ringMini) {
            const offset = 220 - (220 * score / 100);
            ringMini.style.strokeDashoffset = offset;
        }

        const verdictEl = document.getElementById('system-verdict');
        if (verdictEl) {
            verdictEl.textContent = verdict.toUpperCase();
            verdictEl.className = `badge ${score >= 80 ? 'running' : score >= 50 ? 'restarting' : 'offline'}`;
        }
    } catch (_) {}
}

// ── Hardware Metrics ──────────────────────────────────────────────
async function refreshHardware() {
    try {
        const [tRes, sRes, gRes] = await Promise.all([
            fetch('/api/metrics/temperature'),
            fetch('/api/metrics/storage'),
            fetch('/api/metrics/gpu')
        ]);
        const tData = await tRes.json();
        const sData = await sRes.json();
        const gData = await gRes.json();

        // Update GPU Card
        const gpuLoadEl = document.getElementById('stat-gpu-load');
        if (gpuLoadEl) {
            if (gData.active) {
                setText('stat-gpu-load', gData.load !== undefined ? gData.load : 'OK');
                setText('stat-gpu-temp', Math.round(gData.temp || 0));
                setText('stat-gpu-vram', gData.mem_used || 0);
            } else {
                setText('stat-gpu-load', 'OFF');
            }
            
            const gpuIcon = document.getElementById('stat-gpu-card').querySelector('.stat-icon');
            if (gData.load >= 80 || gData.temp >= 80) {
                gpuIcon.style.background = 'rgba(239, 68, 68, 0.15)'; gpuIcon.style.color = '#ef4444';
            } else if (gData.load >= 50 || gData.temp >= 70) {
                gpuIcon.style.background = 'rgba(245, 158, 11, 0.15)'; gpuIcon.style.color = '#f59e0b';
            } else {
                gpuIcon.style.background = 'rgba(249, 115, 22, 0.15)'; gpuIcon.style.color = '#f97316';
            }
        }

        // Update Temperature Card
        const tempValEl = document.getElementById('stat-temp');
        if (tempValEl && tData.cpu_temp !== undefined) {
            const cpu = tData.cpu_temp != null ? Math.round(tData.cpu_temp) : '--';
            const gpu = tData.gpu_temp != null ? Math.round(tData.gpu_temp) : '--';
            
            tempValEl.innerHTML = `${cpu}<span>°C</span> <small style="font-size: 0.6em; opacity: 0.5;">/</small> ${gpu}<span>°C</span>`;
            
            const maxTemp = Math.max(tData.cpu_temp || 0, tData.gpu_temp || 0);
            const tempIcon = tempValEl.parentElement.parentElement.querySelector('.stat-icon');
            if (maxTemp >= 85) {
                tempIcon.style.background = 'rgba(239, 68, 68, 0.15)'; tempIcon.style.color = '#ef4444';
            } else if (maxTemp >= 75) {
                tempIcon.style.background = 'rgba(245, 158, 11, 0.15)'; tempIcon.style.color = '#f59e0b';
            } else {
                tempIcon.style.background = 'rgba(34, 197, 94, 0.15)'; tempIcon.style.color = '#22c55e';
            }
        }

        // Update Storage Card
        const storageGrid = document.getElementById('stat-storage-grid');
        if (storageGrid && sData.disks) {
            let maxUsage = 0;
            let gridHtml = '';
            const driveKeys = Object.keys(sData.disks).sort(); 
            
            if (driveKeys.length === 0) {
                gridHtml = '<div class="stat-sub">No drives detected</div>';
            } else {
                let rows = { names: '', space: '', temp: '' };
                
                driveKeys.forEach(key => {
                    const disk = sData.disks[key];
                    if (disk.percent > maxUsage) maxUsage = disk.percent;
                    
                    const free = disk.free != null ? `${disk.free}G` : '--';
                    const temp = disk.temp != null ? `${Math.round(disk.temp)}°C` : '--';
                    
                    rows.names += `<div>${key}</div>`;
                    rows.space += `<div>${free}</div>`;
                    rows.temp  += `<div>${temp}</div>`;
                });

                gridHtml = `
                    <div class="storage-table">
                        <div class="row names">${rows.names}</div>
                        <div class="row space">${rows.space}</div>
                        <div class="row temp">${rows.temp}</div>
                    </div>
                `;
            }
            
            storageGrid.innerHTML = gridHtml;
            
            const storageIcon = document.getElementById('stat-storage-card').querySelector('.stat-icon');
            if (maxUsage >= 95) {
                storageIcon.style.background = 'rgba(239, 68, 68, 0.15)'; storageIcon.style.color = '#ef4444';
            } else if (maxUsage >= 85) {
                storageIcon.style.background = 'rgba(245, 158, 11, 0.15)'; storageIcon.style.color = '#f59e0b';
            } else {
                storageIcon.style.background = 'rgba(14, 165, 233, 0.15)'; storageIcon.style.color = '#0ea5e9';
            }
        }
    } catch (_) {}
}

// ── Container table ───────────────────────────────────────────────
async function refreshFleet() {
    try {
        const [hRes, mRes] = await Promise.all([
            fetch('/api/health/report'),
            fetch('/api/metrics')
        ]);
        const hData = await hRes.json();
        const mData = await mRes.json();

        // Build metrics lookup
        const metricsByName = {};
        (mData.containers || []).forEach(c => {
            metricsByName[c.name] = c;
            metricsByName[c.name.replace('m3tal-', '')] = c;
        });

        const containers = hData?.agent_health?.monitor_containers?.containers || {};
        const entries = Object.entries(containers);

        // Stat cards
        const online  = entries.filter(([,v]) => ['online','running'].includes((v.status||'').toLowerCase())).length;
        const total   = entries.length;
        setText('stat-containers-val', `${online} / ${total}`);
        setText('stat-containers-sub', 'Running');

        // Uptime
        const uptimeEl = document.getElementById('stat-uptime');
        if (uptimeEl && hData.uptime) uptimeEl.textContent = hData.uptime;

        // Hardware metrics for expanded view
        const [tRes, sRes] = await Promise.all([
            fetch('/api/metrics/temperature'),
            fetch('/api/metrics/storage')
        ]);
        const tData = await tRes.json();
        const sData = await sRes.json();

        // Table body
        const tbody = document.getElementById('fleet-tbody');
        if (!tbody) return;

        if (entries.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="loading-text">Waiting for agent data…</td></tr>';
            return;
        }

        // Sort: running first
        const order = { running: 0, online: 0, restarting: 1, offline: 2, missing: 3, unknown: 4 };
        entries.sort((a, b) => (order[(a[1].status||'').toLowerCase()] ?? 4) - (order[(b[1].status||'').toLowerCase()] ?? 4));

        let html = '';
        entries.forEach(([name, info]) => {
            const status = info.status || 'unknown';
            const sc     = getStatusClass(status);
            const m      = metricsByName[name] || {};
            const cpu    = m.cpu  != null ? m.cpu.toFixed(1)  + '%' : '—';
            const mem    = m.mem_usage || '—';
            const uptime = info.raw_status || '—';
            const cpuClass = m.cpu != null ? getCpuClass(m.cpu) : '';
            
            // Sub-metrics for details
            const cpuTemp = tData.cpu_temp != null ? Math.round(tData.cpu_temp) : '--';
            const gpuTemp = tData.gpu_temp != null ? Math.round(tData.gpu_temp) : '--';
            const storage = sData.disks?.root?.percent != null ? sData.disks.root.percent + '%' : '—';
            const rowId = `details-${name.replace(/[^a-z0-9]/gi, '-')}`;

            html += `
                <tr class="container-row" onclick="toggleRow('${rowId}')">
                    <td><span class="container-name">${name}</span></td>
                    <td><span class="badge ${sc}">${status.toUpperCase()}</span></td>
                    <td class="metric-cell ${cpuClass}">${cpu}</td>
                    <td class="metric-cell">${mem}</td>
                    <td class="metric-cell">${uptime}</td>
                    <td>
                        <div class="actions-cell">
                            <button class="action-btn logs" title="Logs" onclick="event.stopPropagation(); doAction('logs','${name}')">≡</button>
                        </div>
                    </td>
                </tr>
                <tr id="${rowId}" class="details-row" style="display: none;">
                    <td colspan="6">
                        <div class="details-box">
                            <div class="details-metrics">
                                <div><strong>CPU Usage:</strong> ${cpu}</div>
                                <div><strong>MEM Usage:</strong> ${mem}</div>
                                <div><strong>TEMP:</strong> ${cpuTemp}°C / ${gpuTemp}°C</div>
                                <div><strong>STORAGE:</strong> ${storage}</div>
                                <div><strong>UPTIME:</strong> ${uptime}</div>
                            </div>
                            <div class="details-actions">
                                <button class="big-btn heal action-btn" onclick="event.stopPropagation(); doAction('restart','${name}')">↺ Restart</button>
                                <button class="big-btn reboot action-btn" onclick="event.stopPropagation(); doAction('stop','${name}')">■ Stop</button>
                                <button class="big-btn scan action-btn" onclick="event.stopPropagation(); doAction('logs','${name}')">≡ Logs</button>
                            </div>
                        </div>
                    </td>
                </tr>
            `;
        });
        tbody.innerHTML = html;
    } catch (_) {}
}

// ── Activity feed ─────────────────────────────────────────────────
async function refreshActivity() {
    try {
        const [aRes, hRes] = await Promise.all([
            fetch('/api/anomalies'),
            fetch('/api/health/report')
        ]);
        const aData = await aRes.json();
        const hData = await hRes.json();

        const feed = document.getElementById('activity-feed');
        if (!feed) return;

        const issues = [
            ...(aData.issues || []).map(i => ({
                title: i.target || 'Container',
                sub:   i.reason || i.message || '',
                type:  (i.type === 'critical') ? 'warn' : 'warn',
                time:  formatTime()
            })),
            ...(hData.issues || []).map(msg => ({
                title: 'System',
                sub:   msg,
                type:  'warn',
                time:  formatTime()
            }))
        ];

        const now = formatTime();

        const pinned = [{
            title: issues.length === 0 ? 'All systems operational' : `${issues.length} issue(s) detected`,
            sub:   issues.length === 0 ? 'No issues detected'       : 'Review anomalies below',
            type:  issues.length === 0 ? 'ok' : 'warn',
            time:  now
        }];

        const all = [...pinned, ...issues].slice(0, 8);

        const iconMap = { ok: '✓', warn: '⚠', info: 'ℹ' };

        feed.innerHTML = all.map(item => `
            <div class="activity-item">
                <div class="activity-icon ${item.type}">${iconMap[item.type] || 'ℹ'}</div>
                <div class="activity-text">
                    <div class="activity-title">${item.title}</div>
                    ${item.sub ? `<div class="activity-sub">${item.sub}</div>` : ''}
                </div>
                <div class="activity-time">${item.time}</div>
            </div>
        `).join('');
    } catch (_) {}
}


function formatTime() {
    const n = new Date();
    return `${String(n.getHours()).padStart(2,'0')}:${String(n.getMinutes()).padStart(2,'0')}`;
}

// ── Actions (wired to /api/action) ───────────────────────────
async function doAction(action, container) {
    console.log(`Action: ${action} on ${container}`);
    const btn = event.currentTarget;
    const origHtml = btn.innerHTML;
    btn.innerHTML = '⏳';
    btn.disabled = true;

    try {
        const res = await fetch('/api/action', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ action, container })
        });
        const data = await res.json();
        
        if (data.ok) {
            if (action === 'logs' && data.logs) {
                // Show logs in an alert or custom overlay
                alert(`Logs for ${container}:\n\n${data.logs.substring(0, 1000)}${data.logs.length > 1000 ? '...' : ''}`);
            } else {
                btn.style.background = 'var(--green-dim)';
                btn.style.color = 'var(--green)';
                setTimeout(() => {
                    btn.style.background = '';
                    btn.style.color = '';
                }, 2000);
            }
        } else {
            alert(`Error: ${data.error || 'Failed'}`);
            btn.style.background = 'var(--red-dim)';
            btn.style.color = 'var(--red)';
        }
    } catch (e) {
        alert(`Request failed: ${e.message}`);
    } finally {
        setTimeout(() => {
            btn.innerHTML = origHtml;
            btn.disabled = false;
        }, action === 'logs' ? 0 : 2000);
    }
}

async function doGlobalAction(action) {
    console.log(`Global action: ${action}`);
    
    // Add confirmation for reboot
    if (action === 'reboot' && !confirm('Are you sure you want to reboot the entire host system?')) {
        return;
    }

    const btn = event.currentTarget;
    const origHtml = btn.innerHTML;
    btn.innerHTML = '⏳ Processing...';
    btn.disabled = true;

    try {
        const res = await fetch('/api/action', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ action })
        });
        const data = await res.json();
        
        if (data.ok) {
            if (action === 'status') {
                 alert(`Status:\nScore: ${data.score}%\nVerdict: ${data.verdict}\nSystem: ${data.system}`);
            } else {
                 btn.innerHTML = `✅ ${data.message || 'Success'}`;
            }
        } else {
            alert(`Error: ${data.error || 'Failed'}`);
            btn.innerHTML = '❌ Error';
        }
    } catch (e) {
        alert(`Request failed: ${e.message}`);
        btn.innerHTML = '❌ Failed';
    } finally {
        setTimeout(() => {
            btn.innerHTML = origHtml;
            btn.disabled = false;
        }, 3000);
    }
}

// ── Boot ──────────────────────────────────────────────────────────
document.addEventListener('DOMContentLoaded', () => {
    // Telegram Web App Initialization
    if (window.Telegram && window.Telegram.WebApp) {
        window.Telegram.WebApp.ready();
        window.Telegram.WebApp.expand();
    }

    initChart();
    refreshHealth();
    refreshHardware();
    refreshFleet();
    refreshActivity();

    setInterval(refreshHealth,   5000);
    setInterval(refreshHardware, 10000);
    setInterval(refreshFleet,    8000);
    setInterval(refreshActivity, 12000);
});

function setText(id, text) {
    const el = document.getElementById(id);
    if (el) el.textContent = text;
}
