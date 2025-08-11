import React, { useEffect, useState } from 'react';

export const HardwarePage: React.FC = () => {
  const [spaceId, setSpaceId] = useState<number | null>(null);
  const [list, setList] = useState<any[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState<any>({ id: 0, device_id: '', name: '', app_id: null, status: 'online' });
  const [showModal, setShowModal] = useState(false);

  useEffect(() => {
    // try parse space id from url
    const segs = window.location.pathname.split('/');
    const idx = segs.findIndex(s => s === 'space');
    const idStr = idx >= 0 ? segs[idx + 1] : '';
    const sid = Number(idStr);
    if (!isNaN(sid)) setSpaceId(sid);
  }, []);

  useEffect(() => {
    if (!spaceId) return;
    void fetchList();
  }, [spaceId]);

  async function fetchList() {
    setLoading(true);
    try {
      const resp = await fetch('/api/iot/devices/list', {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ space_id: spaceId, page: 1, page_size: 50 })
      });
      const data = await resp.json();
      setList(data.list || []);
      setTotal(data.total || 0);
    } finally {
      setLoading(false);
    }
  }

  async function saveDevice() {
    if (!spaceId) return;
    const payload = { id: form.id, space_id: spaceId, device_id: form.device_id, name: form.name, app_id: form.app_id, status: form.status, description: form.description };
    await fetch('/api/iot/devices/upsert', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) });
    setShowModal(false);
    await fetchList();
  }

  return (
    <div className="p-4">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold">AI 硬件设备（{total}）</h2>
        <button className="px-3 py-1 bg-blue-600 text-white rounded" onClick={() => { setForm({ id: 0, device_id: '', name: '', app_id: null, status: 'online' }); setShowModal(true); }}>新增设备</button>
      </div>

      {loading ? <div>Loading...</div> : (
        <table className="w-full border text-sm">
          <thead>
            <tr className="bg-gray-50">
              <th className="p-2 text-left">设备ID</th>
              <th className="p-2 text-left">名称</th>
              <th className="p-2 text-left">绑定应用</th>
              <th className="p-2 text-left">状态</th>
              <th className="p-2 text-left">操作</th>
            </tr>
          </thead>
          <tbody>
            {list.map((it, idx) => (
              <tr key={idx} className="border-t">
                <td className="p-2">{it.device_id}</td>
                <td className="p-2">{it.name}</td>
                <td className="p-2">{it.app_id ?? '-'}</td>
                <td className="p-2">{it.status}</td>
                <td className="p-2">
                  <button className="px-2 py-1 text-blue-600" onClick={() => { setForm({ id: it.id, device_id: it.device_id, name: it.name, app_id: it.app_id, status: it.status, description: it.description }); setShowModal(true); }}>编辑</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {showModal && (
        <div className="fixed inset-0 bg-black/20 flex items-center justify-center">
          <div className="bg-white p-4 rounded shadow w-[520px]">
            <h3 className="text-md font-semibold mb-3">{form.id ? '编辑设备' : '新增设备'}</h3>
            <div className="space-y-3">
              <div>
                <label className="block text-sm mb-1">设备ID</label>
                <input className="w-full border p-2" value={form.device_id} onChange={e => setForm({ ...form, device_id: e.target.value })} />
              </div>
              <div>
                <label className="block text-sm mb-1">名称</label>
                <input className="w-full border p-2" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
              </div>
              <div>
                <label className="block text-sm mb-1">绑定应用 AppID（可选）</label>
                <input className="w-full border p-2" value={form.app_id ?? ''} onChange={e => setForm({ ...form, app_id: Number(e.target.value) || null })} />
              </div>
              <div>
                <label className="block text-sm mb-1">状态</label>
                <select className="w-full border p-2" value={form.status} onChange={e => setForm({ ...form, status: e.target.value })}>
                  <option value="online">online</option>
                  <option value="offline">offline</option>
                  <option value="pairing">pairing</option>
                  <option value="blocked">blocked</option>
                </select>
              </div>
            </div>
            <div className="mt-4 flex justify-end gap-2">
              <button className="px-3 py-1" onClick={() => setShowModal(false)}>取消</button>
              <button className="px-3 py-1 bg-blue-600 text-white" onClick={saveDevice}>保存</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};