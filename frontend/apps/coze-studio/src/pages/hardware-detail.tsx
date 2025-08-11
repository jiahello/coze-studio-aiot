import React, { useEffect, useMemo, useState } from 'react';

export const HardwareDetailPage: React.FC = () => {
  const segs = window.location.pathname.split('/');
  const spaceId = Number(segs[segs.findIndex(s => s === 'space') + 1]);
  const deviceId = decodeURIComponent(segs[segs.findIndex(s => s === 'hardware') + 1] || '');

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [effective, setEffective] = useState<{provider:string;model:string;voice:string;source:string} | null>(null);
  const [form, setForm] = useState<{provider:string;model:string;voice:string}>({provider:'doubao',model:'speech-1',voice:'doubao-standard'});
  const [voices, setVoices] = useState<any[]>([]);
  const [provider, setProvider] = useState('doubao');
  const [previewUrl, setPreviewUrl] = useState('');
  const [previewLoading, setPreviewLoading] = useState(false);

  const crumb = useMemo(() => (
    <div className="text-sm mb-3">
      <a className="text-blue-600" href={`/space/${spaceId}/hardware`}>AI 硬件</a>
      <span className="mx-2">/</span>
      <span>{deviceId}</span>
    </div>
  ), [spaceId, deviceId]);

  useEffect(() => { void fetchEffective(); }, [deviceId]);
  useEffect(() => { void fetchVoices(); }, [provider, spaceId]);

  async function fetchEffective() {
    if (!deviceId) return;
    setLoading(true); setError(null);
    try {
      const resp = await fetch(`/api/iot/devices/tts/get?device_id=${encodeURIComponent(deviceId)}`);
      if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
      const data = await resp.json();
      setEffective(data);
      setForm({ provider: data.provider, model: data.model, voice: data.voice });
      setProvider(data.provider);
    } catch (e:any) {
      setError(e.message || '加载失败');
    } finally { setLoading(false); }
  }

  async function fetchVoices() {
    const body = { space_id: spaceId ?? null, provider, page: 1, page_size: 100 };
    const resp = await fetch('/api/tts/voices/list', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) });
    const data = await resp.json();
    setVoices(data.list || []);
  }

  async function saveTTS() {
    setError(null);
    const payload = { device_id: deviceId, provider: form.provider, model: form.model, voice: form.voice };
    const resp = await fetch('/api/iot/devices/tts/set', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) });
    if (!resp.ok) { setError(`保存失败: HTTP ${resp.status}`); return; }
    await fetchEffective();
  }

  let previewTimer: any;
  async function doPreview() {
    clearTimeout(previewTimer);
    setPreviewLoading(true); setPreviewUrl(''); setError(null);
    try {
      const resp = await fetch('/api/tts/preview', { method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({ provider: form.provider, model: form.model, voice: form.voice, text: 'Hello', space_id: spaceId }) });
      const data = await resp.json();
      setPreviewUrl(data.sample_url || '');
    } catch (e:any) {
      setError(e.message || '预听失败');
    } finally { setPreviewLoading(false); }
  }
  function previewDebounced() { clearTimeout(previewTimer); previewTimer = setTimeout(doPreview, 300); }

  return (
    <div className="p-4">
      {crumb}
      <h2 className="text-lg font-semibold mb-4">设备详情</h2>
      {loading ? <div>加载中...</div> : error ? <div className="text-red-600">{error}</div> : (
        <div className="space-y-6">
          <div className="border p-4 rounded">
            <h3 className="font-semibold mb-2">TTS 设置</h3>
            <p className="text-xs text-gray-500 mb-2">生效来源：{effective?.source}</p>
            <div className="grid grid-cols-3 gap-3">
              <div>
                <label className="block text-sm mb-1">Provider</label>
                <select className="w-full border p-2" value={form.provider} onChange={e => { setForm({ ...form, provider: e.target.value }); setProvider(e.target.value); previewDebounced(); }}>
                  <option value="doubao">Doubao</option>
                  <option value="aliyun">Aliyun</option>
                  <option value="tencent">Tencent</option>
                  <option value="minimax">MiniMax</option>
                </select>
              </div>
              <div>
                <label className="block text-sm mb-1">Model</label>
                <input className="w-full border p-2" value={form.model} onChange={e => { setForm({ ...form, model: e.target.value }); previewDebounced(); }} />
              </div>
              <div>
                <label className="block text-sm mb-1">Voice</label>
                <select className="w-full border p-2" value={form.voice} onChange={e => { setForm({ ...form, voice: e.target.value }); previewDebounced(); }}>
                  <option value="">请选择音色</option>
                  {voices.map((v:any, i:number) => <option key={i} value={v.voice_code}>{v.name} ({v.voice_code})</option>)}
                </select>
              </div>
            </div>
            <div className="mt-3 flex gap-2 items-center">
              <button className="px-3 py-1 bg-blue-600 text-white rounded" onClick={saveTTS}>保存</button>
              <button className="px-3 py-1 border rounded" onClick={doPreview} disabled={previewLoading}>{previewLoading ? '预听中...' : '预听'}</button>
              {previewUrl && <a className="text-blue-600" href={previewUrl} target="_blank">打开样音</a>}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};