import { useState } from 'react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Search, Plus, ChevronLeft, ChevronRight, Loader2 } from 'lucide-react'

interface Column<T> {
  key: keyof T | string
  header: string
  render?: (item: T) => React.ReactNode
  className?: string
}

interface DataTableProps<T> {
  columns: Column<T>[]
  data: T[]
  loading?: boolean
  page?: number
  totalPages?: number
  onPageChange?: (page: number) => void
  searchValue?: string
  onSearchChange?: (value: string) => void
  searchPlaceholder?: string
  onCreateClick?: () => void
  createLabel?: string
  emptyMessage?: string
  onRowClick?: (item: T) => void
  keyExtractor?: (item: T) => string
}

export function DataTable<T>({
  columns,
  data,
  loading = false,
  page = 1,
  totalPages = 1,
  onPageChange,
  searchValue = '',
  onSearchChange,
  searchPlaceholder = 'Cari...',
  onCreateClick,
  createLabel = 'Tambah',
  emptyMessage = 'Tidak ada data',
  onRowClick,
  keyExtractor,
}: DataTableProps<T>) {
  const [searchTimeout, setSearchTimeout] = useState<number | null>(null)

  const handleSearch = (value: string) => {
    if (searchTimeout) clearTimeout(searchTimeout)
    const t = window.setTimeout(() => onSearchChange?.(value), 300)
    setSearchTimeout(t)
  }

  return (
    <div className="space-y-4">
      {/* Toolbar */}
      <div className="flex flex-col sm:flex-row gap-3 items-start sm:items-center justify-between">
        {onSearchChange && (
          <div className="relative w-full sm:w-72">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
            <Input
              defaultValue={searchValue}
              placeholder={searchPlaceholder}
              onChange={(e) => handleSearch(e.target.value)}
              className="pl-9 h-9"
            />
          </div>
        )}
        {onCreateClick && (
          <Button size="sm" onClick={onCreateClick} className="shrink-0">
            <Plus className="size-4 mr-1" />
            {createLabel}
          </Button>
        )}
      </div>

      {/* Table */}
      <div className="rounded-lg border bg-card overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b bg-muted/50">
                {columns.map((col) => (
                  <th
                    key={String(col.key)}
                    className={col.className + ' text-left text-xs font-medium text-muted-foreground uppercase tracking-wider py-3 px-4'}
                  >
                    {col.header}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={columns.length} className="py-12 text-center text-muted-foreground">
                    <Loader2 className="size-6 mx-auto mb-2 animate-spin" />
                    <span className="text-sm">Memuat data...</span>
                  </td>
                </tr>
              ) : data.length === 0 ? (
                <tr>
                  <td colSpan={columns.length} className="py-12 text-center text-muted-foreground text-sm">
                    {emptyMessage}
                  </td>
                </tr>
              ) : (
                data.map((item, idx) => (
                  <tr
                    key={keyExtractor ? keyExtractor(item) : idx}
                    onClick={() => onRowClick?.(item)}
                    className={onRowClick ? 'cursor-pointer hover:bg-muted/50 border-b last:border-0 transition-colors' : 'border-b last:border-0'}
                  >
                    {columns.map((col) => (
                      <td key={String(col.key)} className="py-3 px-4 text-sm">
                        {col.render ? col.render(item) : String((item as Record<string, unknown>)[String(col.key)] ?? '')}
                      </td>
                    ))}
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Pagination */}
      {onPageChange && totalPages > 1 && (
        <div className="flex items-center justify-between text-sm text-muted-foreground">
          <span>Halaman {page} dari {totalPages}</span>
          <div className="flex gap-1">
            <Button variant="outline" size="xs" disabled={page <= 1} onClick={() => onPageChange(page - 1)}>
              <ChevronLeft className="size-3" />
            </Button>
            <Button variant="outline" size="xs" disabled={page >= totalPages} onClick={() => onPageChange(page + 1)}>
              <ChevronRight className="size-3" />
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}
