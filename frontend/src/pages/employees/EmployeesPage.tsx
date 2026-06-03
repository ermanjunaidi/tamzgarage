import { useEffect, useState } from 'react'
import { api } from '@/api/client'
import { DataTable } from '@/components/shared/DataTable'
import { Badge } from '@/components/ui/badge'
import { Wrench } from 'lucide-react'
import type { Mechanic } from '@/types'

export function EmployeesPage() {
  const [mechanics, setMechanics] = useState<Mechanic[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get<Mechanic[]>('/employees/mechanics').then(setMechanics).finally(() => setLoading(false))
  }, [])

  const columns = [
    { key: 'full_name' as const, header: 'Nama', className: 'font-medium' },
    { key: 'phone' as const, header: 'Telepon' },
    {
      key: 'skills' as const, header: 'Keahlian',
      render: (m: Mechanic) => (
        <div className="flex gap-1 flex-wrap">
          {(m.skills || []).map((s: string) => <Badge key={s} className="bg-slate-100 text-slate-700">{s}</Badge>)}
          {(!m.skills || m.skills.length === 0) && <span className="text-muted-foreground">-</span>}
        </div>
      )
    },
    {
      key: 'active_jobs' as const, header: 'WO Aktif',
      render: (m: Mechanic) => <span className={m.active_jobs > 3 ? 'text-amber-600 font-bold' : ''}>{m.active_jobs}</span>
    },
    {
      key: 'is_active' as const, header: 'Status',
      render: (m: Mechanic) => <Badge variant={m.is_active ? 'success' : 'danger'}>{m.is_active ? 'Aktif' : 'Nonaktif'}</Badge>
    },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Wrench className="size-6" />
        <h1 className="text-2xl font-bold">Karyawan & Mekanik</h1>
      </div>

      <DataTable
        columns={columns}
        data={mechanics}
        loading={loading}
        keyExtractor={(m) => m.id}
        emptyMessage="Belum ada mekanik terdaftar"
      />
    </div>
  )
}
