import { useEffect, useState, useCallback } from 'react'
import { api } from '@/api/client'
import { DataTable } from '@/components/shared/DataTable'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Plus, Package, ArrowDownToLine } from 'lucide-react'
import type { Sparepart, PaginatedResponse, StockMutation } from '@/types'

export function InventoryPage() {
  const [data, setData] = useState<PaginatedResponse<Sparepart>>({ data: [], page: 1, limit: 20 })
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const [selected, setSelected] = useState<Sparepart | null>(null)
  const [mutations, setMutations] = useState<StockMutation[]>([])
  const [showForm, setShowForm] = useState(false)
  const [showStockIn, setShowStockIn] = useState<Sparepart | null>(null)

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const res = await api.get<PaginatedResponse<Sparepart>>(`/spareparts?page=${page}&limit=20&search=${encodeURIComponent(search)}`)
      setData(res)
    } finally { setLoading(false) }
  }, [page, search])

  useEffect(() => { fetchData() }, [fetchData])

  const handleRowClick = async (sp: Sparepart) => {
    setSelected(sp)
    const m = await api.get<StockMutation[]>(`/stock-mutations?sparepart_id=${sp.id}`)
    setMutations(m)
  }

  const columns = [
    { key: 'code' as const, header: 'Kode', className: 'font-mono text-xs' },
    { key: 'name' as const, header: 'Nama', className: 'font-medium' },
    { key: 'category' as const, header: 'Kategori' },
    { key: 'brand' as const, header: 'Merk' },
    {
      key: 'current_stock' as const, header: 'Stok',
      render: (s: Sparepart) => (
        <span className={s.current_stock <= s.min_stock ? 'text-red-600 font-bold' : ''}>
          {s.current_stock} {s.unit}
        </span>
      )
    },
    {
      key: 'min_stock' as const, header: 'Stok Min',
    },
    {
      key: 'selling_price' as const, header: 'Harga Jual',
      render: (s: Sparepart) => `Rp ${s.selling_price.toLocaleString('id-ID')}`
    },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Inventori & Suku Cadang</h1>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => setShowForm(true)}><Plus className="size-4 mr-1" />Tambah</Button>
        </div>
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
        searchPlaceholder="Cari kode, nama, atau merk..."
        onRowClick={handleRowClick}
        keyExtractor={(s) => s.id}
      />

      {/* Detail Dialog */}
      {selected && (
        <Dialog open onOpenChange={() => { setSelected(null); setMutations([]) }}>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-3">
                {selected.name}
                {selected.current_stock <= selected.min_stock && <Badge variant="warning">Stok Menipis</Badge>}
              </DialogTitle>
            </DialogHeader>
            <div className="space-y-4 text-sm">
              <div className="grid grid-cols-3 gap-3 bg-muted/50 rounded-lg p-4">
                <div><strong>Kode:</strong> <span className="font-mono">{selected.code}</span></div>
                <div><strong>Kategori:</strong> {selected.category}</div>
                <div><strong>Merk:</strong> {selected.brand || '-'}</div>
                <div><strong>Stok:</strong> {selected.current_stock} {selected.unit}</div>
                <div><strong>Stok Min:</strong> {selected.min_stock} {selected.unit}</div>
                <div><strong>Harga Jual:</strong> Rp {selected.selling_price.toLocaleString('id-ID')}</div>
              </div>

              <div>
                <div className="flex items-center justify-between mb-2">
                  <strong>Riwayat Stok</strong>
                  <Button size="xs" variant="outline" onClick={() => setShowStockIn(selected)}>
                    <ArrowDownToLine className="size-3 mr-1" />Stok Masuk
                  </Button>
                </div>
                <div className="rounded-lg border overflow-hidden">
                  <table className="w-full text-sm">
                    <thead className="bg-muted/50">
                      <tr>
                        <th className="text-left py-2 px-3">Tanggal</th>
                        <th className="text-left py-2 px-3">Tipe</th>
                        <th className="text-right py-2 px-3">Qty</th>
                        <th className="text-left py-2 px-3">Catatan</th>
                      </tr>
                    </thead>
                    <tbody>
                      {mutations.map((m) => (
                        <tr key={m.id} className="border-t">
                          <td className="py-2 px-3">{new Date(m.created_at).toLocaleDateString('id-ID')}</td>
                          <td className="py-2 px-3">
                            <Badge className={m.mutation_type === 'masuk' ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'}>{m.mutation_type}</Badge>
                          </td>
                          <td className="py-2 px-3 text-right">{m.quantity}</td>
                          <td className="py-2 px-3 text-muted-foreground">{m.notes || '-'}</td>
                        </tr>
                      ))}
                      {mutations.length === 0 && (
                        <tr><td colSpan={4} className="py-4 text-center text-muted-foreground">Belum ada mutasi stok</td></tr>
                      )}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </DialogContent>
        </Dialog>
      )}

      {/* Stock In Dialog */}
      {showStockIn && (
        <StockInDialog
          sparepart={showStockIn}
          onClose={() => setShowStockIn(null)}
          onSaved={() => { setShowStockIn(null); fetchData(); }}
        />
      )}

      {/* Create Dialog */}
      {showForm && (
        <SparepartFormDialog onClose={() => setShowForm(false)} onSaved={() => { setShowForm(false); fetchData() }} />
      )}
    </div>
  )
}

function StockInDialog({ sparepart, onClose, onSaved }: { sparepart: Sparepart; onClose: () => void; onSaved: () => void }) {
  const [qty, setQty] = useState('')
  const [notes, setNotes] = useState('')
  const [saving, setSaving] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      await api.post('/spareparts/stock-in', { sparepart_id: sparepart.id, quantity: parseFloat(qty), notes })
      onSaved()
    } finally { setSaving(false) }
  }

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent className="max-w-sm">
        <DialogHeader><DialogTitle>Stok Masuk: {sparepart.name}</DialogTitle></DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-3">
          <p className="text-sm text-muted-foreground">Stok saat ini: {sparepart.current_stock} {sparepart.unit}</p>
          <Input type="number" step="any" placeholder="Jumlah" value={qty} onChange={(e) => setQty(e.target.value)} required />
          <Input placeholder="Catatan" value={notes} onChange={(e) => setNotes(e.target.value)} />
          <div className="flex gap-2 justify-end">
            <Button variant="outline" type="button" onClick={onClose}>Batal</Button>
            <Button type="submit" disabled={saving}>Simpan</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}

function SparepartFormDialog({ onClose, onSaved }: { onClose: () => void; onSaved: () => void }) {
  const [form, setForm] = useState({ code: '', name: '', category: '', brand: '', unit: 'pcs', purchase_price: '', selling_price: '', current_stock: '0', min_stock: '5' })
  const [saving, setSaving] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      await api.post('/spareparts', {
        ...form,
        purchase_price: parseFloat(form.purchase_price),
        selling_price: parseFloat(form.selling_price),
        current_stock: parseFloat(form.current_stock),
        min_stock: parseFloat(form.min_stock),
      })
      onSaved()
    } finally { setSaving(false) }
  }

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent className="max-w-lg">
        <DialogHeader><DialogTitle>Tambah Suku Cadang</DialogTitle></DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-3">
          <div className="grid grid-cols-2 gap-3">
            <Input placeholder="Kode *" value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} required />
            <Input placeholder="Nama *" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          </div>
          <div className="grid grid-cols-3 gap-3">
            <Input placeholder="Kategori *" value={form.category} onChange={(e) => setForm({ ...form, category: e.target.value })} required />
            <Input placeholder="Merk" value={form.brand} onChange={(e) => setForm({ ...form, brand: e.target.value })} />
            <select className="flex h-9 w-full rounded-lg border border-input bg-background px-3 text-sm" value={form.unit} onChange={(e) => setForm({ ...form, unit: e.target.value })}>
              <option value="pcs">pcs</option>
              <option value="liter">liter</option>
              <option value="set">set</option>
              <option value="box">box</option>
            </select>
          </div>
          <div className="grid grid-cols-2 gap-3">
            <Input type="number" placeholder="Harga Beli *" value={form.purchase_price} onChange={(e) => setForm({ ...form, purchase_price: e.target.value })} required />
            <Input type="number" placeholder="Harga Jual *" value={form.selling_price} onChange={(e) => setForm({ ...form, selling_price: e.target.value })} required />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <Input type="number" placeholder="Stok Awal" value={form.current_stock} onChange={(e) => setForm({ ...form, current_stock: e.target.value })} />
            <Input type="number" placeholder="Stok Minimum" value={form.min_stock} onChange={(e) => setForm({ ...form, min_stock: e.target.value })} />
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
