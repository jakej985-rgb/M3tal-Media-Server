import React, { useState, useEffect } from 'react';
import { api, subscribeWS } from '../api';

export default function Dashboard() {
  const [metrics, setMetrics] = useState(null);
  const [events, setEvents] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    // 1. Fetch initial metrics and poll
    const fetchMetrics = async () => {
      const res = await api.getMetrics();
      if (res.status === 'success') {
        setMetrics(res.data);
        setError(null);
      } else {
        setError(res.error || 'Failed to fetch metrics');
      }
    };

    fetchMetrics();
    const interval = setInterval(fetchMetrics, 2000);

    // 2. Subscribe to WebSocket events for live notification updates
    const sub = subscribeWS('/api/v2/ws/events', (event) => {
      const timeStr = new Date(event.timestamp * 1000).toLocaleTimeString();
      const newEvent = {
        time: timeStr,
        type: event.type,
        detail: JSON.stringify(event.payload),
      };
      setEvents((prev) => [newEvent, ...prev].slice(0, 15)); // Keep last 15
    });

    return () => {
      clearInterval(interval);
      sub.close();
    };
  }, []);

  const formatUptime = (seconds) => {
    if (!seconds) return '0s';
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    if (days > 0) return `${days}d ${hours}h ${minutes}m`;
    return `${hours}h ${minutes}m`;
  };

  const getUsageColor = (val) => {
    if (val >= 85) return 'bg-rose-500';
    if (val >= 60) return 'bg-amber-500';
    return 'bg-emerald-500';
  };

  const getBorderColor = (val) => {
    if (val >= 85) return 'border-rose-500';
    if (val >= 60) return 'border-amber-500';
    return 'border-emerald-500';
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white">System Dashboard</h1>
          <p className="text-slate-400 text-sm">Real-time resource utilization and event streaming.</p>
        </div>
        {error && (
          <div className="px-3 py-1 text-xs rounded-full bg-rose-500/10 text-rose-400 border border-rose-500/20">
            {error}
          </div>
        )}
      </div>

      {metrics && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* CPU Card */}
          <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 relative overflow-hidden">
            <div className="flex justify-between items-center mb-4">
              <span className="text-slate-400 font-semibold text-sm">Processor (CPU)</span>
              <span className="text-lg font-bold text-white">{metrics.cpu_usage?.toFixed(1)}%</span>
            </div>
            <div className="w-full bg-slate-800 rounded-full h-3">
              <div
                className={`h-3 rounded-full transition-all duration-500 ${getUsageColor(metrics.cpu_usage)}`}
                style={{ width: `${Math.min(100, metrics.cpu_usage)}%` }}
              ></div>
            </div>
            <p className="mt-3 text-xs text-slate-500">Average system workload across all active cores.</p>
          </div>

          {/* Memory Card */}
          <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 relative overflow-hidden">
            <div className="flex justify-between items-center mb-4">
              <span className="text-slate-400 font-semibold text-sm">Memory (RAM)</span>
              <span className="text-lg font-bold text-white">{metrics.memory_usage?.toFixed(1)}%</span>
            </div>
            <div className="w-full bg-slate-800 rounded-full h-3">
              <div
                className={`h-3 rounded-full transition-all duration-500 ${getUsageColor(metrics.memory_usage)}`}
                style={{ width: `${Math.min(100, metrics.memory_usage)}%` }}
              ></div>
            </div>
            <p className="mt-3 text-xs text-slate-500">Total physical RAM usage currently allocated.</p>
          </div>

          {/* Disk Card */}
          <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 relative overflow-hidden">
            <div className="flex justify-between items-center mb-4">
              <span className="text-slate-400 font-semibold text-sm">Disk Storage</span>
              <span className="text-lg font-bold text-white">{metrics.disk_usage?.toFixed(1)}%</span>
            </div>
            <div className="w-full bg-slate-800 rounded-full h-3">
              <div
                className={`h-3 rounded-full transition-all duration-500 ${getUsageColor(metrics.disk_usage)}`}
                style={{ width: `${Math.min(100, metrics.disk_usage)}%` }}
              ></div>
            </div>
            <p className="mt-3 text-xs text-slate-500">Root storage usage on primary storage partition.</p>
          </div>
        </div>
      )}

      {metrics && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Host Info */}
          <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 lg:col-span-1 space-y-4">
            <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider">Host Information</h2>
            <div className="border-t border-slate-800/60 pt-3">
              <div className="flex justify-between text-sm py-1.5">
                <span className="text-slate-500">Hostname</span>
                <span className="text-white font-medium">{metrics.hostname}</span>
              </div>
              <div className="flex justify-between text-sm py-1.5">
                <span className="text-slate-500">OS Platform</span>
                <span className="text-white font-medium">Linux</span>
              </div>
              <div className="flex justify-between text-sm py-1.5">
                <span className="text-slate-500">System Uptime</span>
                <span className="text-white font-medium">{formatUptime(metrics.uptime)}</span>
              </div>
              <div className="flex justify-between text-sm py-1.5">
                <span className="text-slate-500">API Status</span>
                <span className="text-emerald-400 font-medium flex items-center gap-1.5">
                  <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
                  Online
                </span>
              </div>
            </div>
          </div>

          {/* Event Stream */}
          <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 lg:col-span-2 flex flex-col h-[320px]">
            <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-4">Live System Event Log</h2>
            <div className="flex-1 overflow-y-auto space-y-2 pr-2 text-xs font-mono">
              {events.length === 0 ? (
                <div className="flex h-full items-center justify-center text-slate-600 italic">
                  Waiting for system events (container states, stack updates, AI jobs)...
                </div>
              ) : (
                events.map((ev, i) => (
                  <div key={i} className="flex gap-4 p-2 bg-slate-950/40 rounded border border-slate-800/40">
                    <span className="text-slate-500">{ev.time}</span>
                    <span className="text-amber-400 font-semibold">{ev.type}</span>
                    <span className="text-slate-300 break-all">{ev.detail}</span>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
