import { useEffect, useState, useCallback } from 'react'
import { api } from '@/api/client'
import { DataTable } from '@/components/shared/DataTable'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { ShoppingCart } from 'lucide-react'
import type { Invoice, PaginatedResponse, Payment } from '@/types'

const statusColors: Record<string, string> = {
  belum_bayar: 'bg-red-100 text-red-700',
  dp: 'bg-amber-100 text-amber-700',
  lunas: 'bg-emerald-100 text-emerald-700',
  batal: 'bg-slate-100 text-slate-700',
}

const statusLabels: Record<string, string> = {
  belum_bayar: 'Belum Bayar',
  dp: 'DP',
  lunas: 'Lunas',
  batal: 'Batal',
}

export function POSPage() {
  const [data, setData] = useState<PaginatedResponse<Invoice>>({ data: [], page: 1, limit: 20 })
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState('')
  const [selected, setSelected] = useState<Invoice | null>(null)
  const [showPayment, setShowPayment] = useState<Invoice | null>(null)

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const res = await api.get<PaginatedResponse<Invoice>>(`/invoices?page=${page}&limit=20&status=${statusFilter}`)
      setData(res)
    } finally { setLoading(false) }
  }, [page, statusFilter])

  useEffect(() => { fetchData() }, [fetchData])

  const columns = [
    { key: 'invoice_number' as const, header: 'Invoice#', className: 'font-mono text-sm font-medium' },
    { key: 'customer_name' as const, header: 'Pelanggan' },
    {
      key: 'grand_total' as const, header: 'Total',
      render: (i: Invoice) => `Rp ${i.grand_total.toLocaleString('id-ID')}`
    },
    {
      key: 'amount_paid' as const, header: 'Dibayar',
      render: (i: Invoice) => `Rp ${i.amount_paid.toLocaleString('id-ID')}`
    },
    {
      key: 'balance_due' as const, header: 'Sisa',
      render: (i: Invoice) => (
        <span className={i.balance_due > 0 ? 'text-red-600 font-bold' : 'text-emerald-600'}>
          Rp {i.balance_due.toLocaleString('id-ID')}
        </span>
      )
    },
    {
      key: 'status' as const, header: 'Status',
      render: (i: Invoice) => <Badge className={statusColors[i.status]}>{statusLabels[i.status]}</Badge>
    },
    {
      key: 'created_at' as const, header: 'Tanggal',
      render: (i: Invoice) => new Date(i.created_at).toLocaleDateString('id-ID')
    },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <ShoppingCart className="size-6" />
          <h1 className="text-2xl font-bold">POS / Kasir</h1>
        </div>
      </div>

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
        onRowClick={async (i) => {
          const full = await api.get<Invoice>(`/invoices/${i.id}`)
          setSelected(full)
        }}
        keyExtractor={(i) => i.id}
      />

      {/* Invoice Detail */}
      {selected && (
        <Dialog open onOpenChange={() => setSelected(null)}>
          <DialogContent className="max-w-xl">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-3">
                <span className="font-mono">{selected.invoice_number}</span>
                <Badge className={statusColors[selected.status]}>{statusLabels[selected.status]}</Badge>
              </DialogTitle>
            </DialogHeader>
            <div className="space-y-4 text-sm">
              <div className="grid grid-cols-2 gap-3 bg-muted/50 rounded-lg p-4">
                <div><strong>Pelanggan:</strong> {selected.customer_name || '-'}</div>
                <div><strong>Tanggal:</strong> {new Date(selected.created_at).toLocaleDateString('id-ID')}</div>
                <div><strong>Subtotal:</strong> Rp {selected.subtotal.toLocaleString('id-ID')}</div>
                <div><strong>Diskon:</strong> Rp {selected.discount.toLocaleString('id-ID')}</div>
                <div><strong>Pajak:</strong> Rp {selected.tax_amount.toLocaleString('id-ID')}</div>
                <div><strong className="text-base">Grand Total:</strong> <span className="font-bold text-base">Rp {selected.grand_total.toLocaleString('id-ID')}</span></div>
              </div>

              {/* Payments */}
              {selected.payments && selected.payments.length > 0 && (
                <div>
                  <strong className="block mb-2">Pembayaran:</strong>
                  <div className="rounded-lg border overflow-hidden">
                    <table className="w-full text-sm">
                      <thead className="bg-muted/50">
                        <tr><th className="py-2 px-3 text-left">Metode</th><th className="py-2 px-3 text-right">Jumlah</th><th className="py-2 px-3 text-left">Tanggal</th></tr>
                      </thead>
                      <tbody>
                        {selected.payments.map((p) => (
                          <tr key={p.id} className="border-t">
                            <td className="py-2 px-3">{p.method}</td>
                            <td className="py-2 px-3 text-right">Rp {p.amount.toLocaleString('id-ID')}</td>
                            <td className="py-2 px-3">{new Date(p.created_at).toLocaleDateString('id-ID')}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}

              {selected.status !== 'lunas' && selected.status !== 'batal' && (
                <Button className="w-full" onClick={() => setShowPayment(selected)}>Proses Pembayaran</Button>
              )}
            </div>
          </DialogContent>
        </Dialog>
      )}

      {/* Payment Dialog */}
      {showPayment && (
        <PaymentDialog
          invoice={showPayment}
          onClose={() => setShowPayment(null)}
          onSaved={() => { setShowPayment(null); fetchData(); setSelected(null) }}
        />
      )}
    </div>
  )
}

function PaymentDialog({ invoice, onClose, onSaved }: { invoice: Invoice; onClose: () => void; onSaved: () => void }) {
  const [amount, setAmount] = useState(invoice.balance_due.toString())
  const [method, setMethod] = useState('tunai')
  const [ref, setRef] = useState('')
  const [saving, setSaving] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      await api.post('/payments', { invoice_id: invoice.id, amount: parseFloat(amount), method, reference: ref })
      onSaved()
    } finally { setSaving(false) }
  }

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent className="max-w-sm">
        <DialogHeader><DialogTitle>Pembayaran</DialogTitle></DialogHeader>
        <div className="text-sm mb-3">
          <p><strong>Invoice:</strong> {invoice.invoice_number}</p>
          <p><strong>Total:</strong> Rp {invoice.grand_total.toLocaleString('id-ID')}</p>
          <p><strong>Sisa:</strong> <span className="text-red-600">Rp {invoice.balance_due.toLocaleString('id-ID')}</span></p>
        </div>
        <form onSubmit={handleSubmit} className="space-y-3">
          <div>
            <label className="text-sm font-medium">Jumlah</label>
            <Input type="number" step="any" value={amount} onChange={(e) => setAmount(e.target.value)} required />
          </div>
          <div>
            <label className="text-sm font-medium">Metode</label>
            <select className="flex h-9 w-full rounded-lg border border-input bg-background px-3 text-sm" value={method} onChange={(e) => setMethod(e.target.value)}>
              <option value="tunai">Tunai</option>
              <option value="debit">Debit</option>
              <option value="kredit">Kredit</option>
              <option value="qris">QRIS</option>
              <option value="transfer">Transfer</option>
            </select>
          </div>
          <Input placeholder="Referensi (opsional)" value={ref} onChange={(e) => setRef(e.target.value)} />
          <div className="flex gap-2 justify-end">
            <Button variant="outline" type="button" onClick={onClose}>Batal</Button>
            <Button type="submit" disabled={saving}>Bayar</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
