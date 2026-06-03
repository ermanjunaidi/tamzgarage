import { useEffect, useState, useCallback } from 'react'
import { api } from '@/api/client'
import { DataTable } from '@/components/shared/DataTable'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Plus, Truck, Phone, Mail, MapPin, Users } from 'lucide-react'
import type { Supplier, PurchaseOrder } from '@/types'

export function SuppliersPage() {
  const [data, setData] = useState<Supplier[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [selected, setSelected] = useState<Supplier | null>(null)
  const [pos, setPOs] = useState<PurchaseOrder[]>([])
  const [showForm, setShowForm] = useState(false)
  const [showCreatePO, setShowCreatePO] = useState(false)

  const fetchData = useCallback(async () => {
    setLoading(true)
    try { setData(await api.get<Supplier[]>(`/suppliers?search=${encodeURIComponent(search)}`)) }
    finally { setLoading(false) }
  }, [search])

  useEffect(() => { fetchData() }, [fetchData])

  const handleSelect = async (s: Supplier) => {
    setSelected(s)
    const allPOs = await api.get<PurchaseOrder[]>('/purchase-orders')
    setPOs(allPOs.filter((po) => po.supplier_id === s.id))
  }

  const columns = [
    { key: 'code' as const, header: 'Kode', className: 'font-mono text-xs' },
    { key: 'name' as const, header: 'Nama', className: 'font-medium' },
    { key: 'contact_person' as const, header: 'Kontak' },
    { key: 'phone' as const, header: 'Telepon' },
    { key: 'email' as const, header: 'Email' },
    {
      key: 'is_active' as const, header: 'Status',
      render: (s: Supplier) => <Badge variant={s.is_active ? 'success' : 'danger'}>{s.is_active ? 'Aktif' : 'Nonaktif'}</Badge>
    },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Supplier</h1>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => setShowCreatePO(true)}><Truck className="size-4 mr-1" />Buat PO</Button>
          <Button onClick={() => setShowForm(true)}><Plus className="size-4 mr-1" />Tambah Supplier</Button>
        </div>
      </div>

      <DataTable
        columns={columns}
        data={data}
        loading={loading}
        searchValue={search}
        onSearchChange={(v) => setSearch(v)}
        searchPlaceholder="Cari nama, kontak, atau telepon..."
        onRowClick={handleSelect}
        keyExtractor={(s) => s.id}
      />

      {/* Detail */}
      {selected && (
        <Dialog open onOpenChange={() => { setSelected(null); setPOs([]) }}>
          <DialogContent className="max-w-lg">
            <DialogHeader><DialogTitle>{selected.name}</DialogTitle></DialogHeader>
            <div className="space-y-3 text-sm">
              {selected.contact_person && <div className="flex items-center gap-2"><Users className="size-4 text-muted-foreground" />{selected.contact_person}</div>}
              {selected.phone && <div className="flex items-center gap-2"><Phone className="size-4 text-muted-foreground" />{selected.phone}</div>}
              {selected.email && <div className="flex items-center gap-2"><Mail className="size-4 text-muted-foreground" />{selected.email}</div>}
              {selected.address && <div className="flex items-center gap-2"><MapPin className="size-4 text-muted-foreground" />{selected.address}</div>}
              {selected.tax_id && <p><strong>NPWP:</strong> {selected.tax_id}</p>}
            </div>
          </DialogContent>
        </Dialog>
      )}

      {/* Create Supplier */}
      {showForm && (
        <SupplierFormDialog onClose={() => setShowForm(false)} onSaved={() => { setShowForm(false); fetchData() }} />
      )}

      {/* Create PO */}
      {showCreatePO && (
        <CreatePODialog suppliers={data} onClose={() => setShowCreatePO(false)} onSaved={() => setShowCreatePO(false)} />
      )}
    </div>
  )
}

function SupplierFormDialog({ onClose, onSaved }: { onClose: () => void; onSaved: () => void }) {
  const [form, setForm] = useState({ name: '', contact_person: '', phone: '', email: '', address: '', tax_id: '' })
  const [saving, setSaving] = useState(false)

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent className="max-w-lg">
        <DialogHeader><DialogTitle>Tambah Supplier</DialogTitle></DialogHeader>
        <form onSubmit={async (e) => { e.preventDefault(); setSaving(true); try { await api.post('/suppliers', form); onSaved() } finally { setSaving(false) } }} className="space-y-3">
          <Input placeholder="Nama *" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          <Input placeholder="Kontak Person" value={form.contact_person} onChange={(e) => setForm({ ...form, contact_person: e.target.value })} />
          <Input placeholder="Telepon" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
          <Input placeholder="Email" type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
          <Input placeholder="Alamat" value={form.address} onChange={(e) => setForm({ ...form, address: e.target.value })} />
          <Input placeholder="NPWP" value={form.tax_id} onChange={(e) => setForm({ ...form, tax_id: e.target.value })} />
          <div className="flex gap-2 justify-end">
            <Button variant="outline" type="button" onClick={onClose}>Batal</Button>
            <Button type="submit" disabled={saving}>Simpan</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}

function CreatePODialog({ suppliers, onClose, onSaved }: { suppliers: Supplier[]; onClose: () => void; onSaved: () => void }) {
  const [supplierId, setSupplierId] = useState('')
  const [items, setItems] = useState([{ item_name: '', quantity: '', unit_price: '', total_price: '' }])
  const [saving, setSaving] = useState(false)

  const updateItem = (idx: number, field: string, value: string) => {
    const newItems = [...items]
    ;(newItems[idx] as Record<string, string>)[field] = value
    if (field === 'quantity' || field === 'unit_price') {
      const qty = parseFloat(newItems[idx].quantity || '0')
      const price = parseFloat(newItems[idx].unit_price || '0')
      newItems[idx].total_price = (qty * price).toString()
    }
    setItems(newItems)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      await api.post('/purchase-orders', {
        supplier_id: supplierId,
        items: items.map((i) => ({ ...i, quantity: parseFloat(i.quantity), unit_price: parseFloat(i.unit_price), total_price: parseFloat(i.total_price) })),
      })
      onSaved()
    } finally { setSaving(false) }
  }

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent className="max-w-lg">
        <DialogHeader><DialogTitle>Buat Purchase Order</DialogTitle></DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <select className="flex h-9 w-full rounded-lg border border-input bg-background px-3 text-sm" value={supplierId} onChange={(e) => setSupplierId(e.target.value)} required>
            <option value="">Pilih Supplier...</option>
            {suppliers.map((s) => <option key={s.id} value={s.id}>{s.name}</option>)}
          </select>
          {items.map((item, idx) => (
            <div key={idx} className="grid grid-cols-4 gap-2 items-end">
              <Input placeholder="Item" value={item.item_name} onChange={(e) => updateItem(idx, 'item_name', e.target.value)} required />
              <Input type="number" placeholder="Qty" value={item.quantity} onChange={(e) => updateItem(idx, 'quantity', e.target.value)} required />
              <Input type="number" placeholder="Harga" value={item.unit_price} onChange={(e) => updateItem(idx, 'unit_price', e.target.value)} required />
              <Input placeholder="Total" value={item.total_price} readOnly className="bg-muted" />
            </div>
          ))}
          <Button type="button" variant="ghost" size="sm" onClick={() => setItems([...items, { item_name: '', quantity: '', unit_price: '', total_price: '' }])}>+ Tambah Item</Button>
          <div className="flex gap-2 justify-end">
            <Button variant="outline" type="button" onClick={onClose}>Batal</Button>
            <Button type="submit" disabled={saving}>Simpan PO</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
