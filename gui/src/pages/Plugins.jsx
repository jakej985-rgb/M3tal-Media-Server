import React, { useState, useEffect } from 'react';
import { api } from '../api';

export default function Plugins() {
  const [plugins, setPlugins] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchPlugins = async () => {
    setLoading(true);
    const res = await api.getPlugins();
    if (res.status === 'success') {
      setPlugins(res.data);
      setError(null);
    } else {
      setError(res.error || 'Failed to fetch plugins');
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchPlugins();
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white">Plugin Registry</h1>
          <p className="text-slate-400 text-sm">Inspect loaded Routes, Stacks, and Middleware plugins.</p>
        </div>
        <button
          onClick={fetchPlugins}
          className="px-4 py-2 text-sm font-semibold rounded-lg bg-slate-800 hover:bg-slate-700 text-white border border-slate-700/60 transition-colors"
        >
          Refresh Plugins
        </button>
      </div>

      {error && (
        <div className="p-4 rounded-xl bg-rose-500/10 text-rose-400 border border-rose-500/20 text-sm">
          {error}
        </div>
      )}

      {loading && !plugins ? (
        <div className="flex justify-center items-center h-48 text-slate-500">
          Loading plugin catalog...
        </div>
      ) : plugins ? (
        <div className="space-y-8">
          <div className="text-sm font-semibold text-slate-400 bg-slate-900 border border-slate-800 rounded-lg px-4 py-3">
            Summary: <span className="text-emerald-400 font-bold">{plugins.summary || '0 plugins loaded'}</span>
          </div>

          {/* Stacks Plugins */}
          <div className="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden shadow-sm">
            <div className="px-6 py-4 bg-slate-950/20 border-b border-slate-800/80">
              <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider">📁 Stack Plugins</h2>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="border-b border-slate-800/80 bg-slate-950/10 text-xs font-semibold text-slate-400 uppercase tracking-wider">
                    <th className="p-4">Name</th>
                    <th className="p-4">Version</th>
                    <th className="p-4">Path / Category</th>
                    <th className="p-4">Description</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800/40 text-sm">
                  {plugins.stacks?.map((item, idx) => (
                    <tr key={idx} className="hover:bg-slate-950/10 transition-colors">
                      <td className="p-4 font-semibold text-white">{item.metadata?.name || item.name}</td>
                      <td className="p-4 font-mono text-xs text-slate-400">{item.metadata?.version || '1.0.0'}</td>
                      <td className="p-4 text-xs font-mono text-slate-400">{item.composePath || 'custom-compose.yml'}</td>
                      <td className="p-4 text-slate-400">{item.metadata?.description || 'No description provided.'}</td>
                    </tr>
                  ))}
                  {(!plugins.stacks || plugins.stacks.length === 0) && (
                    <tr>
                      <td colSpan="4" className="text-center p-6 text-slate-500 italic">No stack plugins registered.</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>

          {/* Routes Plugins */}
          <div className="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden shadow-sm">
            <div className="px-6 py-4 bg-slate-950/20 border-b border-slate-800/80">
              <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider">🌐 Route Plugins</h2>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="border-b border-slate-800/80 bg-slate-950/10 text-xs font-semibold text-slate-400 uppercase tracking-wider">
                    <th className="p-4">Name</th>
                    <th className="p-4">Domain / Port</th>
                    <th className="p-4">Entrypoints</th>
                    <th className="p-4">Middlewares</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800/40 text-sm">
                  {plugins.routes?.map((item, idx) => (
                    <tr key={idx} className="hover:bg-slate-950/10 transition-colors">
                      <td className="p-4 font-semibold text-white">{item.metadata?.name || item.service}</td>
                      <td className="p-4 font-mono text-xs text-emerald-400">
                        {item.domain}:{item.port}
                      </td>
                      <td className="p-4 text-xs text-slate-400">{item.entrypoints || 'web'}</td>
                      <td className="p-4">
                        <div className="flex flex-wrap gap-1">
                          {item.middlewares?.map((m) => (
                            <span key={m} className="px-1.5 py-0.5 text-[10px] font-mono rounded bg-slate-800 text-slate-300 border border-slate-700">
                              {m}
                            </span>
                          )) || <span className="text-slate-500 text-xs">none</span>}
                        </div>
                      </td>
                    </tr>
                  ))}
                  {(!plugins.routes || plugins.routes.length === 0) && (
                    <tr>
                      <td colSpan="4" className="text-center p-6 text-slate-500 italic">No route plugins registered.</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>

          {/* Middleware Plugins */}
          <div className="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden shadow-sm">
            <div className="px-6 py-4 bg-slate-950/20 border-b border-slate-800/80">
              <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider">🔒 Middleware Plugins</h2>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="border-b border-slate-800/80 bg-slate-950/10 text-xs font-semibold text-slate-400 uppercase tracking-wider">
                    <th className="p-4">Name</th>
                    <th className="p-4">Type</th>
                    <th className="p-4">Author</th>
                    <th className="p-4">Config Summary</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800/40 text-sm">
                  {plugins.middleware?.map((item, idx) => (
                    <tr key={idx} className="hover:bg-slate-950/10 transition-colors">
                      <td className="p-4 font-semibold text-white">{item.metadata?.name || item.name}</td>
                      <td className="p-4 font-mono text-xs text-slate-400">{item.type}</td>
                      <td className="p-4 text-slate-400">{item.metadata?.author || 'System'}</td>
                      <td className="p-4 text-xs font-mono text-slate-400 max-w-[200px] truncate" title={JSON.stringify(item.config)}>
                        {JSON.stringify(item.config)}
                      </td>
                    </tr>
                  ))}
                  {(!plugins.middleware || plugins.middleware.length === 0) && (
                    <tr>
                      <td colSpan="4" className="text-center p-6 text-slate-500 italic">No middleware plugins registered.</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      ) : (
        <div className="text-center p-8 bg-slate-900 border border-slate-800 rounded-xl text-slate-500">
          No plugin data available.
        </div>
      )}
    </div>
  );
}
