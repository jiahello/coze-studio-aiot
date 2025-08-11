import React, { useEffect, useState } from 'react';

export const TTSPage: React.FC = () => {
  const [spaceId, setSpaceId] = useState<number | null>(null);
  const [provider, setProvider] = useState('doubao');
  const [voices, setVoices] = useState<any[]>([]);
  const [appId, setAppId] = useState<number | ''>('');
  const [model, setModel] = useState('speech-1');
  const [voice, setVoice] = useState('');

  useEffect(() => {
    const segs = window.location.pathname.split('/');
    const idx = segs.findIndex(s => s === 'space');
    const idStr = idx >= 0 ? segs[idx + 1] : '';
    const sid = Number(idStr);
    if (!isNaN(sid)) setSpaceId(sid);
  }, []);

  useEffect(() => {
    void fetchVoices();
  }, [provider, spaceId]);

  async function fetchVoices() {
    const body = { space_id: spaceId ?? null, provider, page: 1, page_size: 100 };
    const resp = await fetch('/api/tts/voices/list', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) });
    const data = await resp.json();
    setVoices(data.list || []);
  }

  async function saveAppTTS() {
    if (!appId) return;
    const payload = { app_id: Number(appId), provider, model, voice };
    await fetch('/api/apps/tts/set', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) });
    alert('已保存');
  }

  return (
    <div className="p-4 space-y-4">
      <div className="flex items-center gap-3">
        <h2 className="text-lg font-semibold">TTS 音色</h2>
        <select className="border p-2" value={provider} onChange={e => setProvider(e.target.value)}>
          <option value="doubao">Doubao</option>
          <option value="aliyun">Aliyun</option>
          <option value="tencent">Tencent</option>
          <option value="minimax">MiniMax</option>
        </select>
      </div>

      <table className="w-full border text-sm">
        <thead>
          <tr className="bg-gray-50">
            <th className="p-2 text-left">名称</th>
            <th className="p-2 text-left">编码</th>
            <th className="p-2 text-left">语言</th>
            <th className="p-2 text-left">示例</th>
            <th className="p-2 text-left">操作</th>
          </tr>
        </thead>
        <tbody>
          {voices.map((v, i) => (
            <tr key={i} className="border-t">
              <td className="p-2">{v.name}</td>
              <td className="p-2">{v.voice_code}</td>
              <td className="p-2">{v.language || '-'}</td>
              <td className="p-2">{v.sample_url ? <a href={v.sample_url} target="_blank">试听</a> : '-'}</td>
              <td className="p-2"><button className="text-blue-600" onClick={() => setVoice(v.voice_code)}>选择</button></td>
            </tr>
          ))}
        </tbody>
      </table>

      <div className="mt-6 border p-4 rounded w-[560px]">
        <h3 className="font-semibold mb-3">应用 TTS 设置</h3>
        <div className="space-y-3">
          <div>
            <label className="block text-sm mb-1">AppID</label>
            <input className="w-full border p-2" value={appId} onChange={e => setAppId(e.target.value as any)} />
          </div>
          <div>
            <label className="block text-sm mb-1">模型</label>
            <input className="w-full border p-2" value={model} onChange={e => setModel(e.target.value)} />
          </div>
          <div>
            <label className="block text-sm mb-1">音色</label>
            <input className="w-full border p-2" value={voice} onChange={e => setVoice(e.target.value)} placeholder="从上方列表选择或手填" />
          </div>
        </div>
        <div className="mt-3 flex gap-2">
          <button className="px-3 py-1 bg-blue-600 text-white rounded" onClick={saveAppTTS}>保存</button>
        </div>
      </div>
    </div>
  );
};