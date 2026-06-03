import { useEffect, useState, useCallback } from 'react'
import { api } from '@/api/client'
import { DataTable } from '@/components/shared/DataTable'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Plus, Wrench } from 'lucide-react'
import type { WorkOrder, WOStatus, PaginatedResponse } from '@/types'

const statusColors: Record<WOStatus, string> = {
  menunggu: 'bg-slate-100 text-slate-700',
  diagnosis: 'bg-sky-100 text-sky-700',
  menunggu_persetujuan: 'bg-amber-100 text-amber-700',
  dikerjakan: 'bg-blue-100 text-blue-700',
  qc: 'bg-purple-100 text-purple-700',
  selesai: 'bg-emerald-100 text-emerald-700',
  diambil: 'bg-green-100 text-green-700',
  batal: 'bg-red-100 text-red-700',
}

const statusLabels: Record<WOStatus, string> = {
  menunggu: 'Menunggu',
  diagnosis: 'Diagnosis',
  menunggu_persetujuan: 'Menunggu Persetujuan',
  dikerjakan: 'Dikerjakan',
  qc: 'QC',
  selesai: 'Selesai',
  diambil: 'Diambil',
  batal: 'Batal',
}

const workflowButtons: { status: WOStatus[]; next: WOStatus }[] = [
  { status: ['menunggu'], next: 'diagnosis' },
  { status: ['diagnosis'], next: 'menunggu_persetujuan' },
  { status: ['menunggu_persetujuan'], next: 'dikerjakan' },
  { status: ['dikerjakan'], next: 'qc' },
  { status: ['qc'], next: 'selesai' },
  { status: ['selesai'], next: 'diambil' },
]

export function WorkOrdersPage() {
  const [data, setData] = useState<PaginatedResponse<WorkOrder>>({ data: [], page: 1, limit: 20 })
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)
  const [selected, setSelected] = useState<WorkOrder | null>(null)
  const [showCreate, setShowCreate] = useState(false)

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const res = await api.get<PaginatedResponse<WorkOrder>>(`/work-orders?page=${page}&limit=20&status=${statusFilter}&search=${encodeURIComponent(search)}`)
      setData(res)
    } finally { setLoading(false) }
  }, [page, search, statusFilter])

  useEffect(() => { fetchData() }, [fetchData])

  const updateStatus = async (id: string, status: string) => {
    await api.patch(`/work-orders/${id}/status`, { status })
    fetchData()
  }

  const columns = [
    { key: 'wo_number' as const, header: 'WO#', className: 'font-mono text-sm font-medium' },
    { key: 'customer_name' as const, header: 'Pelanggan' },
    { key: 'vehicle_plate' as const, header: 'Kendaraan', render: (w: WorkOrder) => <span className="font-mono">{w.vehicle_plate || '-'}</span> },
    {
      key: 'status' as const, header: 'Status',
      render: (w: WorkOrder) => <Badge className={statusColors[w.status]}>{statusLabels[w.status]}</Badge>
    },
    {
      key: 'grand_total' as const, header: 'Total',
      render: (w: WorkOrder) => w.grand_total ? `Rp ${w.grand_total.toLocaleString('id-ID')}` : '-'
    },
    { key: 'created_at' as const, header: 'Tanggal', render: (w: WorkOrder) => new Date(w.created_at).toLocaleDateString('id-ID') },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Work Order</h1>
        <Button onClick={() => setShowCreate(true)}><Plus className="size-4 mr-1" />Buat WO</Button>
      </div>

      {/* Status filter tabs */}
      <div className="flex gap-2 flex-wrap">
        <Button variant={statusFilter === '' ? 'default' : 'outline'} size="xs" onClick={() => setStatusFilter('')}>Semua</Button>
        {Object.entries(statusLabels).map(([k, v]) => (
          <Button key={k} variant={statusFilter === k ? 'default' : 'outline'} size="xs" onClick={() => setStatusFilter(k)}>{v}</Button>
        ))}
      </div>

      <DataTable
        columns={columns}
        data={data.data}
        loading={loading}
        page={page}
        totalPages={Math.ceil((data.total ?? 0) / 20)}
        onPageChange={setPage}
        searchValue={search}
        onSearchChange={(v) => { setSearch(v); setPage(1) }}
        searchPlaceholder="Cari WO, pelanggan, atau plat..."
        onRowClick={async (w) => {
          const full = await api.get<WorkOrder>(`/work-orders/${w.id}`)
          setSelected(full)
        }}
        keyExtractor={(w) => w.id}
      />

      {/* Detail Dialog */}
      {selected && (
        <Dialog open onOpenChange={() => setSelected(null)}>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-3">
                <span className="font-mono">{selected.wo_number}</span>
                <Badge className={statusColors[selected.status]}>{statusLabels[selected.status]}</Badge>
              </DialogTitle>
            </DialogHeader>
            <div className="space-y-4 text-sm">
              <div className="grid grid-cols-2 gap-3 bg-muted/50 rounded-lg p-4">
                <div><strong>Pelanggan:</strong> {selected.customer_name || '-'}</div>
                <div><strong>Kendaraan:</strong> <span className="font-mono">{selected.vehicle_plate || '-'}</span></div>
                <div><strong>Keluhan:</strong> {selected.complaint}</div>
                <div><strong>Estimasi:</strong> {selected.estimated_cost ? `Rp ${selected.estimated_cost.toLocaleString('id-ID')}` : '-'}</div>
              </div>

              {selected.diagnosis && (
                <div><strong>Diagnosis:</strong> <p className="text-muted-foreground mt-1">{selected.diagnosis}</p></div>
              )}

              {/* Services */}
              {selected.services && selected.services.length > 0 && (
                <div>
                  <strong className="block mb-2">Jasa:</strong>
                  <div className="rounded-lg border overflow-hidden">
                    <table className="w-full text-sm">
                      <thead className="bg-muted/50"><tr><th className="text-left py-2 px-3">Jasa</th><th className="text-right py-2 px-3">Qty</th><th className="text-right py-2 px-3">Harga</th><th className="text-right py-2 px-3">Total</th></tr></thead>
                      <tbody>
                        {selected.services.map((s) => (
                          <tr key={s.id} className="border-t">
                            <td className="py-2 px-3">{s.service_name}</td>
                            <td className="text-right py-2 px-3">{s.quantity}</td>
                            <td className="text-right py-2 px-3">{s.unit_price.toLocaleString('id-ID')}</td>
                            <td className="text-right py-2 px-3 font-medium">{s.total_price.toLocaleString('id-ID')}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}

              {/* Spareparts */}
              {selected.spareparts && selected.spareparts.length > 0 && (
                <div>
                  <strong className="block mb-2">Suku Cadang:</strong>
                  <div className="rounded-lg border overflow-hidden">
                    <table className="w-full text-sm">
                      <thead className="bg-muted/50"><tr><th className="text-left py-2 px-3">Item</th><th className="text-right py-2 px-3">Qty</th><th className="text-right py-2 px-3">Harga</th><th className="text-right py-2 px-3">Total</th></tr></thead>
                      <tbody>
                        {selected.spareparts.map((sp) => (
                          <tr key={sp.id} className="border-t">
                            <td className="py-2 px-3">{sp.sparepart_name}</td>
                            <td className="text-right py-2 px-3">{sp.quantity}</td>
                            <td className="text-right py-2 px-3">{sp.unit_price.toLocaleString('id-ID')}</td>
                            <td className="text-right py-2 px-3 font-medium">{sp.total_price.toLocaleString('id-ID')}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}

              {/* Mechanics */}
              {selected.mechanics && selected.mechanics.length > 0 && (
                <div><strong>Mekanik:</strong> {selected.mechanics.map((m) => m.full_name).join(', ')}</div>
              )}

              {/* Workflow buttons */}
              <div className="flex gap-2 flex-wrap pt-2 border-t">
                {workflowButtons.filter((w) => w.status.includes(selected.status)).map((w) => (
                  <Button key={w.next} size="sm" variant={selected.status === 'menunggu_persetujuan' ? 'default' : 'default'} onClick={() => updateStatus(selected.id, w.next)}>
                    {w.next === 'menunggu_persetujuan' ? 'Kirim ke Persetujuan' : `→ ${statusLabels[w.next]}`}
                  </Button>
                ))}
                {selected.status !== 'batal' && selected.status !== 'diambil' && (
                  <Button size="sm" variant="destructive" onClick={() => updateStatus(selected.id, 'batal')}>Batalkan</Button>
                )}
              </div>
            </div>
          </DialogContent>
        </Dialog>
      )}

      {/* Create WO Dialog */}
      {showCreate && (
        <CreateWODialog onClose={() => setShowCreate(false)} onSaved={() => { setShowCreate(false); fetchData() }} />
      )}
    </div>
  )
}

function CreateWODialog({ onClose, onSaved }: { onClose: () => void; onSaved: () => void }) {
  const [form, setForm] = useState({ customer_id: '', vehicle_id: '', complaint: '', estimated_cost: '', estimated_hours: '' })
  const [saving, setSaving] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      await api.post('/work-orders', {
        customer_id: form.customer_id || null,
        vehicle_id: form.vehicle_id || null,
        complaint: form.complaint,
        estimated_cost: parseFloat(form.estimated_cost) || 0,
        estimated_hours: parseFloat(form.estimated_hours) || null,
        services: [], spareparts: [], mechanics: [],
      })
      onSaved()
    } finally { setSaving(false) }
  }

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent className="max-w-lg">
        <DialogHeader><DialogTitle>Buat Work Order Baru</DialogTitle></DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-3">
          <Input placeholder="Customer ID (opsional)" value={form.customer_id} onChange={(e) => setForm({ ...form, customer_id: e.target.value })} />
          <Input placeholder="Vehicle ID (opsional)" value={form.vehicle_id} onChange={(e) => setForm({ ...form, vehicle_id: e.target.value })} />
          <textarea className="flex w-full rounded-lg border border-input bg-background px-3 py-2 text-sm min-h-[80px]" placeholder="Keluhan *" value={form.complaint} onChange={(e) => setForm({ ...form, complaint: e.target.value })} required />
          <div className="grid grid-cols-2 gap-3">
            <Input type="number" placeholder="Estimasi Biaya" value={form.estimated_cost} onChange={(e) => setForm({ ...form, estimated_cost: e.target.value })} />
            <Input type="number" placeholder="Estimasi Jam" value={form.estimated_hours} onChange={(e) => setForm({ ...form, estimated_hours: e.target.value })} />
          </div>
          <div className="flex gap-2 justify-end">
            <Button variant="outline" type="button" onClick={onClose}>Batal</Button>
            <Button type="submit" disabled={saving}>Simpan</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
