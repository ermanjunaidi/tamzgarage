import { useEffect, useState } from 'react'
import { api } from '@/api/client'
import { StatCard } from '@/components/shared/StatCard'
import { Wrench, DollarSign, Package, Clock, CheckCircle, AlertTriangle, TrendingUp, Car } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import type { DashboardStats, RevenueChart } from '@/types'

function formatRupiah(n: number) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n)
}

export function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [chart, setChart] = useState<RevenueChart[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      api.get<DashboardStats>('/dashboard/stats'),
      api.get<RevenueChart[]>('/dashboard/revenue-chart'),
    ])
      .then(([s, c]) => { setStats(s); setChart(c) })
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return (
      <div className="space-y-6">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {Array.from({ length: 8 }).map((_, i) => (
            <div key={i} className="h-28 rounded-xl bg-muted animate-pulse" />
          ))}
        </div>
      </div>
    )
  }

  const maxRev = Math.max(...chart.map((c) => c.amount), 1)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <span className="text-sm text-muted-foreground">{new Date().toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}</span>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard label="Kendaraan Masuk Hari Ini" value={stats?.today_vehicles_in ?? 0} icon={Car} />
        <StatCard label="Revenue Hari Ini" value={formatRupiah(stats?.today_revenue ?? 0)} icon={DollarSign} />
        <StatCard label="Antrian Menunggu" value={stats?.waiting_queue ?? 0} icon={Clock} />
        <StatCard label="Sedang Dikerjakan" value={stats?.in_progress ?? 0} icon={Wrench} />
        <StatCard label="Selesai Hari Ini" value={stats?.completed_today ?? 0} icon={CheckCircle} />
        <StatCard label="Stok Menipis" value={stats?.low_stock_items ?? 0} icon={AlertTriangle} className={stats?.low_stock_items && stats.low_stock_items > 0 ? 'border-amber-300 bg-amber-50' : ''} />
        <StatCard label="Revenue Minggu Ini" value={formatRupiah(stats?.week_revenue ?? 0)} icon={TrendingUp} />
        <StatCard label="Revenue Bulan Ini" value={formatRupiah(stats?.month_revenue ?? 0)} icon={DollarSign} />
      </div>

      {/* Revenue Chart */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Revenue 30 Hari Terakhir</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-48 flex items-end gap-1">
            {chart.map((c) => (
              <div key={c.date} className="flex-1 flex flex-col items-center gap-1">
                <span className="text-[10px] text-muted-foreground">
                  {formatRupiah(c.amount)}
                </span>
                <div
                  className="w-full rounded-t bg-primary/70 hover:bg-primary transition-colors min-h-[2px]"
                  style={{ height: `${(c.amount / maxRev) * 100}%` }}
                  title={`${c.date}: ${formatRupiah(c.amount)}`}
                />
                <span className="text-[10px] text-muted-foreground">
                  {new Date(c.date).getDate()}
                </span>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Summary cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <Card>
          <CardHeader className="pb-2"><CardTitle className="text-sm">Total Pelanggan</CardTitle></CardHeader>
          <CardContent><p className="text-3xl font-bold">{stats?.total_customers ?? 0}</p></CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2"><CardTitle className="text-sm">Total Suku Cadang</CardTitle></CardHeader>
          <CardContent><p className="text-3xl font-bold">{stats?.total_spareparts ?? 0}</p></CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2"><CardTitle className="text-sm">WO Aktif</CardTitle></CardHeader>
          <CardContent><p className="text-3xl font-bold">{stats?.active_work_orders ?? 0}</p></CardContent>
        </Card>
      </div>
    </div>
  )
}

