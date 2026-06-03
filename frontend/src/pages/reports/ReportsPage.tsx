import { useEffect, useState } from 'react'
import { api } from '@/api/client'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { BarChart3, TrendingUp, Package, Wrench, Users } from 'lucide-react'
import { useAuth } from '@/context/AuthContext'

interface RevenueRow { date: string; count: number; total: number; paid: number }
interface StockRow { code: string; name: string; category: string; brand: string; current_stock: number; min_stock: number; selling_price: number; stock_value: number; low_stock: boolean }
interface WOStatusRow { status: string; count: number; value: number }
interface CustomerRow { full_name: string; phone: string; wo_count: number; total_spent: number }

function formatRupiah(n: number) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n)
}

const statusLabels: Record<string, string> = {
  menunggu: 'Menunggu', diagnosis: 'Diagnosis', menunggu_persetujuan: 'Menunggu Persetujuan',
  dikerjakan: 'Dikerjakan', qc: 'QC', selesai: 'Selesai', diambil: 'Diambil', batal: 'Batal',
}

export function ReportsPage() {
  const [revenue, setRevenue] = useState<RevenueRow[]>([])
  const [stock, setStock] = useState<StockRow[]>([])
  const [woStatus, setWOStatus] = useState<WOStatusRow[]>([])
  const [topCust, setTopCust] = useState<CustomerRow[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      api.get<RevenueRow[]>('/reports/revenue'),
      api.get<{ data: StockRow[] }>('/reports/stock'),
      api.get<WOStatusRow[]>('/reports/work-orders'),
      api.get<CustomerRow[]>('/reports/top-customers'),
    ]).then(([r, s, w, c]) => {
      setRevenue(r)
      setStock(s.data)
      setWOStatus(w)
      setTopCust(c)
    }).finally(() => setLoading(false))
  }, [])

  if (loading) {
    return (
      <div className="space-y-6">
        <h1 className="text-2xl font-bold">Laporan</h1>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {Array.from({ length: 4 }).map((_, i) => <div key={i} className="h-64 rounded-xl bg-muted animate-pulse" />)}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <BarChart3 className="size-6" />
        <h1 className="text-2xl font-bold">Laporan & Analytics</h1>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Revenue Report */}
        <Card>
          <CardHeader><CardTitle className="text-base flex items-center gap-2"><TrendingUp className="size-4" />Pendapatan</CardTitle></CardHeader>
          <CardContent>
            <div className="max-h-72 overflow-y-auto">
              <table className="w-full text-sm">
                <thead className="bg-muted/50 sticky top-0">
                  <tr><th className="py-2 px-3 text-left">Tanggal</th><th className="py-2 px-3 text-right">WO</th><th className="py-2 px-3 text-right">Total</th></tr>
                </thead>
                <tbody>
                  {revenue.map((r) => (
                    <tr key={r.date} className="border-t hover:bg-muted/30">
                      <td className="py-1.5 px-3">{new Date(r.date).toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })}</td>
                      <td className="py-1.5 px-3 text-right">{r.count}</td>
                      <td className="py-1.5 px-3 text-right font-medium">{formatRupiah(r.total)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>

        {/* Work Order Status */}
        <Card>
          <CardHeader><CardTitle className="text-base flex items-center gap-2"><Wrench className="size-4" />Status Work Order</CardTitle></CardHeader>
          <CardContent>
            <table className="w-full text-sm">
              <thead className="bg-muted/50"><tr><th className="py-2 px-3 text-left">Status</th><th className="py-2 px-3 text-right">Jumlah</th><th className="py-2 px-3 text-right">Nilai</th></tr></thead>
              <tbody>
                {woStatus.map((w) => (
                  <tr key={w.status} className="border-t">
                    <td className="py-1.5 px-3">{statusLabels[w.status] || w.status}</td>
                    <td className="py-1.5 px-3 text-right">{w.count}</td>
                    <td className="py-1.5 px-3 text-right">{formatRupiah(w.value)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>

        {/* Stock Report */}
        <Card>
          <CardHeader><CardTitle className="text-base flex items-center gap-2"><Package className="size-4" />Stok Inventori</CardTitle></CardHeader>
          <CardContent>
            <div className="max-h-72 overflow-y-auto">
              <table className="w-full text-sm">
                <thead className="bg-muted/50 sticky top-0">
                  <tr><th className="py-2 px-3 text-left">Item</th><th className="py-2 px-3 text-right">Stok</th><th className="py-2 px-3 text-right">Nilai</th></tr>
                </thead>
                <tbody>
                  {stock.map((s) => (
                    <tr key={s.code} className={s.low_stock ? 'border-t bg-red-50' : 'border-t hover:bg-muted/30'}>
                      <td className="py-1.5 px-3">
                        <span className="font-medium">{s.name}</span>
                        {s.low_stock && <Badge variant="warning" className="ml-2">Stok Menipis</Badge>}
                      </td>
                      <td className="py-1.5 px-3 text-right font-mono">{s.current_stock}</td>
                      <td className="py-1.5 px-3 text-right">{formatRupiah(s.stock_value)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>

        {/* Top Customers */}
        <Card>
          <CardHeader><CardTitle className="text-base flex items-center gap-2"><Users className="size-4" />Top 10 Pelanggan</CardTitle></CardHeader>
          <CardContent>
            <table className="w-full text-sm">
              <thead className="bg-muted/50"><tr><th className="py-2 px-3 text-left">Nama</th><th className="py-2 px-3 text-right">WO</th><th className="py-2 px-3 text-right">Total</th></tr></thead>
              <tbody>
                {topCust.map((c) => (
                  <tr key={c.phone} className="border-t hover:bg-muted/30">
                    <td className="py-1.5 px-3">{c.full_name}</td>
                    <td className="py-1.5 px-3 text-right">{c.wo_count}</td>
                    <td className="py-1.5 px-3 text-right font-medium">{formatRupiah(c.total_spent)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
