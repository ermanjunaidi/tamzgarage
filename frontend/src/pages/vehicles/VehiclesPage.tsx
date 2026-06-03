import { useEffect, useState, useCallback } from 'react'
import { api } from '@/api/client'
import { DataTable } from '@/components/shared/DataTable'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Plus } from 'lucide-react'
import type { Vehicle } from '@/types'

export function VehiclesPage() {
  const [data, setData] = useState<Vehicle[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [selected, setSelected] = useState<Vehicle | null>(null)

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const res = await api.get<Vehicle[]>(`/vehicles?search=${encodeURIComponent(search)}`)
      setData(res)
    } finally { setLoading(false) }
  }, [search])

  useEffect(() => { fetchData() }, [fetchData])

  const columns = [
    {
      key: 'plate_number' as const, header: 'No. Plat', className: 'font-medium',
      render: (v: Vehicle) => (
        <div className="flex items-center gap-2">
          <div className="rounded-md bg-slate-100 px-2 py-0.5 font-mono text-sm font-bold">{v.plate_number}</div>
        </div>
      )
    },
    { key: 'brand' as const, header: 'Merk' },
    { key: 'model' as const, header: 'Model' },
    { key: 'year' as const, header: 'Tahun' },
    { key: 'color' as const, header: 'Warna' },
    { key: 'customer_name' as const, header: 'Pemilik' },
    {
      key: 'last_km' as const, header: 'KM Terakhir',
      render: (v: Vehicle) => v.last_km ? `${v.last_km.toLocaleString()} km` : '-'
    },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Kendaraan</h1>
        <Button onClick={() => setShowForm(true)}><Plus className="size-4 mr-1" />Tambah Kendaraan</Button>
      </div>

      <DataTable
        columns={columns}
        data={data}
        loading={loading}
        searchValue={search}
        onSearchChange={(v) => setSearch(v)}
        searchPlaceholder="Cari plat, merk, atau pemilik..."
        onRowClick={setSelected}
        keyExtractor={(v) => v.id}
      />

      {selected && (
        <Dialog open onOpenChange={() => setSelected(null)}>
          <DialogContent className="max-w-lg">
            <DialogHeader><DialogTitle>{selected.brand} {selected.model}</DialogTitle></DialogHeader>
            <div className="space-y-2 text-sm">
              <p><strong>Plat:</strong> <span className="font-mono">{selected.plate_number}</span></p>
              <p><strong>Pemilik:</strong> {selected.customer_name || '-'}</p>
              {selected.year && <p><strong>Tahun:</strong> {selected.year}</p>}
              {selected.color && <p><strong>Warna:</strong> {selected.color}</p>}
              {selected.vin && <p><strong>VIN:</strong> {selected.vin}</p>}
              {selected.engine_number && <p><strong>No. Mesin:</strong> {selected.engine_number}</p>}
              {selected.last_km && <p><strong>KM Terakhir:</strong> {selected.last_km.toLocaleString()} km</p>}
              {selected.next_service_km && <p><strong>Servis Berikutnya:</strong> {selected.next_service_km.toLocaleString()} km</p>}
              {selected.notes && <p className="text-muted-foreground border-t pt-2 mt-2">{selected.notes}</p>}
            </div>
          </DialogContent>
        </Dialog>
      )}

      {showForm && (
        <VehicleFormDialog onClose={() => setShowForm(false)} onSaved={() => { setShowForm(false); fetchData() }} />
      )}
    </div>
  )
}

function VehicleFormDialog({ onClose, onSaved }: { onClose: () => void; onSaved: () => void }) {
  const [form, setForm] = useState({
    plate_number: '', brand: '', model: '', variant: '', year: '', color: '',
    vin: '', engine_number: '', last_km: '', next_service_km: '', notes: '',
  })
  const [saving, setSaving] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      await api.post('/vehicles', {
        ...form,
        year: form.year ? parseInt(form.year) : null,
        last_km: form.last_km ? parseInt(form.last_km) : null,
        next_service_km: form.next_service_km ? parseInt(form.next_service_km) : null,
      })
      onSaved()
    } finally { setSaving(false) }
  }

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent className="max-w-lg">
        <DialogHeader><DialogTitle>Tambah Kendaraan</DialogTitle></DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-3">
          <Input placeholder="No. Plat *" value={form.plate_number} onChange={(e) => setForm({ ...form, plate_number: e.target.value })} required />
          <div className="grid grid-cols-2 gap-3">
            <Input placeholder="Merk *" value={form.brand} onChange={(e) => setForm({ ...form, brand: e.target.value })} required />
            <Input placeholder="Model" value={form.model} onChange={(e) => setForm({ ...form, model: e.target.value })} />
          </div>
          <div className="grid grid-cols-3 gap-3">
            <Input placeholder="Tahun" type="number" value={form.year} onChange={(e) => setForm({ ...form, year: e.target.value })} />
            <Input placeholder="Warna" value={form.color} onChange={(e) => setForm({ ...form, color: e.target.value })} />
            <Input placeholder="KM Terakhir" type="number" value={form.last_km} onChange={(e) => setForm({ ...form, last_km: e.target.value })} />
          </div>
          <Input placeholder="VIN / No. Rangka" value={form.vin} onChange={(e) => setForm({ ...form, vin: e.target.value })} />
          <Input placeholder="No. Mesin" value={form.engine_number} onChange={(e) => setForm({ ...form, engine_number: e.target.value })} />
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
