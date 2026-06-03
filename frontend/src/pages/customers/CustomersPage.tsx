import { useEffect, useState, useCallback } from 'react'
import { api } from '@/api/client'
import { DataTable } from '@/components/shared/DataTable'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Plus, Phone, Mail, MapPin, Star } from 'lucide-react'
import type { Customer, PaginatedResponse } from '@/types'

const categoryColors: Record<string, string> = {
  regular: 'bg-slate-100 text-slate-700',
  member: 'bg-blue-100 text-blue-700',
  fleet: 'bg-purple-100 text-purple-700',
}

export function CustomersPage() {
  const [data, setData] = useState<PaginatedResponse<Customer>>({ data: [], page: 1, limit: 20 })
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const [selected, setSelected] = useState<Customer | null>(null)
  const [showForm, setShowForm] = useState(false)

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const res = await api.get<PaginatedResponse<Customer>>(`/customers?page=${page}&limit=20&search=${encodeURIComponent(search)}`)
      setData(res)
    } finally {
      setLoading(false)
    }
  }, [page, search])

  useEffect(() => { fetchData() }, [fetchData])

  const columns = [
    { key: 'code' as const, header: 'Kode', className: 'w-24' },
    { key: 'full_name' as const, header: 'Nama', className: 'font-medium' },
    { key: 'phone' as const, header: 'Telepon' },
    {
      key: 'category' as const, header: 'Kategori',
      render: (c: Customer) => <Badge className={categoryColors[c.category] || ''}>{c.category}</Badge>
    },
    { key: 'vehicle_count' as const, header: 'Kendaraan', className: 'text-center' },
    { key: 'created_at' as const, header: 'Terdaftar', render: (c: Customer) => new Date(c.created_at).toLocaleDateString('id-ID') },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Pelanggan</h1>
        <Button onClick={() => setShowForm(true)}><Plus className="size-4 mr-1" />Tambah Pelanggan</Button>
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
        searchPlaceholder="Cari nama, telepon, atau kode..."
        onRowClick={setSelected}
        keyExtractor={(c) => c.id}
      />

      {/* Detail */}
      {selected && (
        <Dialog open onOpenChange={() => setSelected(null)}>
          <DialogContent className="max-w-lg">
            <DialogHeader><DialogTitle>{selected.full_name}</DialogTitle></DialogHeader>
            <div className="space-y-3 text-sm">
              <div className="flex items-center gap-2"><Phone className="size-4 text-muted-foreground" />{selected.phone}</div>
              {selected.email && <div className="flex items-center gap-2"><Mail className="size-4 text-muted-foreground" />{selected.email}</div>}
              {selected.address && <div className="flex items-center gap-2"><MapPin className="size-4 text-muted-foreground" />{selected.address}</div>}
              <div className="flex items-center gap-2">
                <Badge className={categoryColors[selected.category]}>{selected.category}</Badge>
                {selected.loyalty_points > 0 && <span className="flex items-center gap-1 text-amber-600"><Star className="size-3" />{selected.loyalty_points} poin</span>}
              </div>
              {selected.notes && <p className="text-muted-foreground border-t pt-2 mt-2">{selected.notes}</p>}
            </div>
          </DialogContent>
        </Dialog>
      )}

      {/* Create Form */}
      {showForm && (
        <CustomerFormDialog
          onClose={() => setShowForm(false)}
          onSaved={() => { setShowForm(false); fetchData() }}
        />
      )}
    </div>
  )
}

function CustomerFormDialog({ onClose, onSaved }: { onClose: () => void; onSaved: () => void }) {
  const [form, setForm] = useState({ full_name: '', phone: '', email: '', address: '', category: 'regular', notes: '' })
  const [saving, setSaving] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      await api.post('/customers', form)
      onSaved()
    } finally { setSaving(false) }
  }

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent className="max-w-lg">
        <DialogHeader><DialogTitle>Tambah Pelanggan</DialogTitle></DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input placeholder="Nama Lengkap" value={form.full_name} onChange={(e) => setForm({ ...form, full_name: e.target.value })} required />
          <Input placeholder="No. Telepon" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} required />
          <Input placeholder="Email" type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
          <Input placeholder="Alamat" value={form.address} onChange={(e) => setForm({ ...form, address: e.target.value })} />
          <select className="flex h-9 w-full rounded-lg border border-input bg-background px-3 text-sm" value={form.category} onChange={(e) => setForm({ ...form, category: e.target.value })}>
            <option value="regular">Regular</option>
            <option value="member">Member</option>
            <option value="fleet">Fleet</option>
          </select>
          <Input placeholder="Catatan" value={form.notes} onChange={(e) => setForm({ ...form, notes: e.target.value })} />
          <div className="flex gap-2 justify-end">
            <Button variant="outline" type="button" onClick={onClose}>Batal</Button>
            <Button type="submit" disabled={saving}>Simpan</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
