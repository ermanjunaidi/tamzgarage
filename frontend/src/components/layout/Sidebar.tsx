import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'
import { cn } from '@/lib/utils'
import {
  LayoutDashboard, Users, Car, Wrench, Package, UserCog, ShoppingCart,
  Truck, FileText, BarChart3, LogOut, ChevronDown, ChevronRight, Menu, X
} from 'lucide-react'
import { useState } from 'react'
import { Button } from '@/components/ui/button'

const menuItems = [
  {
    label: 'Dashboard',
    icon: LayoutDashboard,
    path: '/',
    roles: ['super_admin', 'admin_cabang', 'service_advisor', 'mekanik', 'kasir', 'gudang'],
  },
  { label: 'Pelanggan', icon: Users, path: '/customers', roles: ['super_admin', 'admin_cabang', 'service_advisor'] },
  { label: 'Kendaraan', icon: Car, path: '/vehicles', roles: ['super_admin', 'admin_cabang', 'service_advisor'] },
  { label: 'Work Order', icon: Wrench, path: '/work-orders', roles: ['super_admin', 'admin_cabang', 'service_advisor', 'mekanik'] },
  { label: 'Inventori', icon: Package, path: '/inventory', roles: ['super_admin', 'admin_cabang', 'gudang'] },
  { label: 'Karyawan', icon: UserCog, path: '/employees', roles: ['super_admin', 'admin_cabang'] },
  { label: 'POS / Kasir', icon: ShoppingCart, path: '/pos', roles: ['super_admin', 'admin_cabang', 'kasir', 'service_advisor'] },
  { label: 'Supplier', icon: Truck, path: '/suppliers', roles: ['super_admin', 'admin_cabang', 'gudang'] },
  { label: 'Laporan', icon: BarChart3, path: '/reports', roles: ['super_admin', 'admin_cabang'] },
]

const statusLabels: Record<string, string> = {
  super_admin: 'Super Admin',
  admin_cabang: 'Admin Cabang',
  service_advisor: 'Service Advisor',
  mekanik: 'Mekanik',
  kasir: 'Kasir',
  gudang: 'Gudang',
  pelanggan: 'Pelanggan',
}

export function Sidebar() {
  const { user, hasRole, logout } = useAuth()
  const location = useLocation()
  const navigate = useNavigate()
  const [collapsed, setCollapsed] = useState(false)
  const [mobileOpen, setMobileOpen] = useState(false)

  const visibleItems = menuItems.filter((item) => hasRole(...item.roles))

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const sidebarContent = (
    <div className={cn('flex h-full flex-col bg-sidebar text-sidebar-foreground border-r border-border transition-all', collapsed ? 'w-16' : 'w-60')}>
      {/* Logo */}
      <div className="flex h-14 items-center gap-2 px-3 border-b border-sidebar-border">
        <div className="flex size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground font-bold text-sm shrink-0">
          BP
        </div>
        {!collapsed && <span className="font-semibold text-sm">BengkelPro</span>}
      </div>

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto py-2 px-2">
        {visibleItems.map((item) => {
          const Icon = item.icon
          const isActive = location.pathname === item.path || (item.path !== '/' && location.pathname.startsWith(item.path))
          return (
            <Link
              key={item.path}
              to={item.path}
              onClick={() => setMobileOpen(false)}
              className={cn(
                'flex items-center gap-3 rounded-lg px-3 py-2 text-sm mb-0.5 transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground',
                isActive && 'bg-sidebar-accent text-sidebar-accent-foreground font-medium',
                collapsed && 'justify-center px-2'
              )}
              title={collapsed ? item.label : undefined}
            >
              <Icon className="size-4 shrink-0" />
              {!collapsed && <span>{item.label}</span>}
            </Link>
          )
        })}
      </nav>

      {/* User & Logout */}
      <div className="border-t border-sidebar-border p-3">
        {!collapsed && user && (
          <div className="mb-2 text-xs text-sidebar-foreground/60">
            <div className="truncate font-medium text-sidebar-foreground">{user.full_name}</div>
            <div>{statusLabels[user.role] || user.role}</div>
          </div>
        )}
        <Button
          variant="ghost"
          size="sm"
          className={cn('w-full justify-start text-sidebar-foreground/60 hover:text-sidebar-foreground', collapsed && 'justify-center')}
          onClick={handleLogout}
        >
          <LogOut className={cn('size-4', !collapsed && 'mr-2')} />
          {!collapsed && 'Logout'}
        </Button>
      </div>
    </div>
  )

  return (
    <>
      {/* Mobile toggle */}
      <button
        className="fixed top-3 left-3 z-50 lg:hidden bg-background border rounded-md p-1.5"
        onClick={() => setMobileOpen(!mobileOpen)}
      >
        {mobileOpen ? <X className="size-5" /> : <Menu className="size-5" />}
      </button>

      {/* Desktop */}
      <aside className="hidden lg:block h-screen sticky top-0">
        <button
          className="absolute -right-3 top-6 z-10 bg-background border rounded-full p-0.5 size-6 flex items-center justify-center hover:bg-muted"
          onClick={() => setCollapsed(!collapsed)}
        >
          {collapsed ? <ChevronRight className="size-3" /> : <ChevronDown className="size-3" />}
        </button>
        {sidebarContent}
      </aside>

      {/* Mobile */}
      {mobileOpen && (
        <div className="fixed inset-0 z-40 lg:hidden">
          <div className="absolute inset-0 bg-black/40" onClick={() => setMobileOpen(false)} />
          <div className="relative z-50 h-screen">
            {sidebarContent}
          </div>
        </div>
      )}
    </>
  )
}
