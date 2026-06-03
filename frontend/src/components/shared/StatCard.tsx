import { cn } from '@/lib/utils'
import type { LucideIcon } from 'lucide-react'

export function StatCard({
  label, value, icon: Icon, trend, trendLabel, className
}: {
  label: string
  value: string | number
  icon?: LucideIcon
  trend?: 'up' | 'down'
  trendLabel?: string
  className?: string
}) {
  return (
    <div className={cn('rounded-xl border bg-card p-4 lg:p-5 shadow-sm', className)}>
      <div className="flex items-start justify-between">
        <div>
          <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{label}</p>
          <p className="mt-1 text-2xl font-bold">{value}</p>
          {trendLabel && (
            <p className={cn('mt-1 text-xs font-medium', trend === 'up' ? 'text-emerald-600' : 'text-red-600')}>
              {trend === 'up' ? '↑' : '↓'} {trendLabel}
            </p>
          )}
        </div>
        {Icon && (
          <div className="rounded-lg bg-primary/10 p-2">
            <Icon className="size-5 text-primary" />
          </div>
        )}
      </div>
    </div>
  )
}
