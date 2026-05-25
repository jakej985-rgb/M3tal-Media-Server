import React, { useState, useEffect } from 'react';
import { api, subscribeWS } from '../api';

export default function AI() {
  const [models, setModels] = useState([]);
  const [selectedModel, setSelectedModel] = useState('');
  const [prompt, setPrompt] = useState('');
  const [mode, setMode] = useState('');
  
  const [queue, setQueue] = useState([]);
  const [history, setHistory] = useState([]);
  const [executing, setExecuting] = useState(false);
  const [error, setError] = useState(null);

  const fetchData = async () => {
    // Fetch models
    const modelsRes = await api.getAIModels();
    if (modelsRes.status === 'success') {
      setModels(modelsRes.data || []);
      if (modelsRes.data && modelsRes.data.length > 0 && !selectedModel) {
        setSelectedModel(modelsRes.data[0]);
      }
    }

    // Fetch queue
    const queueRes = await api.getAIQueue();
    if (queueRes.status === 'success') {
      setQueue(queueRes.data || []);
    }
  };

  useEffect(() => {
    fetchData();

    // Auto-update queue lists on AI state change events
    const sub = subscribeWS('/api/v2/ws/events', (event) => {
      if (event.type?.startsWith('ai.job.')) {
        fetchData();
      }
    });

    return () => sub.close();
  }, []);

  const handleRun = async (e) => {
    e.preventDefault();
    if (!prompt.trim()) return;

    setExecuting(true);
    setError(null);

    // Save prompt request
    const promptItem = {
      prompt,
      model: selectedModel,
      status: 'pending',
      response: '',
      time: new Date().toLocaleTimeString(),
    };
    setHistory((prev) => [promptItem, ...prev]);

    const res = await api.runAI(prompt, mode);
    
    setHistory((prev) => {
      const updated = [...prev];
      if (res.status === 'success') {
        updated[0].status = 'completed';
        updated[0].response = res.data.response;
        updated[0].model = res.data.model; // Active model used
      } else {
        updated[0].status = 'failed';
        updated[0].response = res.error || 'Execution failed';
      }
      return updated;
    });

    setPrompt('');
    setExecuting(false);
    fetchData();
  };

  const getStatusColor = (status) => {
    const s = status.toLowerCase();
    if (s === 'completed' || s === 'success') return 'text-emerald-400';
    if (s === 'failed') return 'text-rose-400';
    return 'text-amber-400';
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-white">AI Copilot Playground</h1>
        <p className="text-slate-400 text-sm">Submit inference requests and inspect sequential background execution queues.</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Run Prompt Form */}
        <div className="lg:col-span-2 space-y-6">
          <form onSubmit={handleRun} className="bg-slate-900 border border-slate-800 rounded-xl p-6 space-y-4 shadow-sm">
            <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider">Submit Inference Prompt</h2>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              {/* Model */}
              <div>
                <label className="block text-xs font-semibold text-slate-400 mb-1.5 uppercase">Model Option</label>
                <select
                  value={selectedModel}
                  onChange={(e) => setSelectedModel(e.target.value)}
                  className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-emerald-500"
                >
                  {models.map((m) => (
                    <option key={m} value={m}>{m}</option>
                  ))}
                  {models.length === 0 && <option value="">No models available</option>}
                </select>
              </div>

              {/* Mode */}
              <div>
                <label className="block text-xs font-semibold text-slate-400 mb-1.5 uppercase">Execution Mode</label>
                <select
                  value={mode}
                  onChange={(e) => setMode(e.target.value)}
                  className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-emerald-500"
                >
                  <option value="">Default</option>
                  <option value="chat">Chat (General)</option>
                  <option value="code">Code (Programming)</option>
                </select>
              </div>
            </div>

            <div>
              <label className="block text-xs font-semibold text-slate-400 mb-1.5 uppercase">Prompt Context</label>
              <textarea
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                placeholder="Ask M3TAL anything..."
                rows="4"
                className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-sm text-white placeholder-slate-600 focus:outline-none focus:border-emerald-500 font-sans"
              ></textarea>
            </div>

            <div className="flex justify-end">
              <button
                type="submit"
                disabled={executing || !prompt.trim()}
                className="px-5 py-2 text-sm font-semibold rounded-lg bg-emerald-600 hover:bg-emerald-500 disabled:bg-slate-800 disabled:text-slate-600 text-white shadow transition-colors flex items-center gap-2"
              >
                {executing && (
                  <span className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                )}
                {executing ? 'Executing...' : 'Run Generation'}
              </button>
            </div>
          </form>

          {/* History */}
          <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 space-y-4">
            <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider">Prompt Logs</h2>
            <div className="space-y-4 divide-y divide-slate-800/40">
              {history.map((item, idx) => (
                <div key={idx} className={`${idx > 0 ? 'pt-4' : ''} space-y-2`}>
                  <div className="flex justify-between items-center text-xs">
                    <span className="text-slate-500">{item.time} | Model: {item.model}</span>
                    <span className={`font-semibold uppercase ${getStatusColor(item.status)}`}>{item.status}</span>
                  </div>
                  <div className="text-sm font-medium text-slate-200">Q: {item.prompt}</div>
                  {item.response && (
                    <div className="bg-slate-950/60 rounded-lg p-4 text-xs font-mono text-slate-300 whitespace-pre-wrap leading-relaxed border border-slate-950">
                      {item.response}
                    </div>
                  )}
                </div>
              ))}
              {history.length === 0 && (
                <div className="text-center py-6 text-slate-500 italic">No prompt runs recorded yet.</div>
              )}
            </div>
          </div>
        </div>

        {/* AI Queue Panel */}
        <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 h-[400px] flex flex-col shadow-sm">
          <h2 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-4">Sequential Queue (Con=1)</h2>
          <div className="flex-1 overflow-y-auto space-y-3 pr-2">
            {queue.map((job) => (
              <div key={job.id} className="p-3 bg-slate-950/40 border border-slate-800/60 rounded-lg space-y-2">
                <div className="flex justify-between items-center text-xs">
                  <span className="text-emerald-400 font-bold font-mono">{job.id}</span>
                  <span className={`font-semibold uppercase ${getStatusColor(job.status)}`}>{job.status}</span>
                </div>
                <div className="text-xs font-medium text-white truncate">Q: {job.prompt}</div>
                <div className="text-[10px] text-slate-500">Mode: {job.mode || 'default'}</div>
              </div>
            ))}
            {queue.length === 0 && (
              <div className="flex h-full items-center justify-center text-slate-700 italic text-sm">
                Queue is empty (MaxConcurrency=1 worker idle)
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
