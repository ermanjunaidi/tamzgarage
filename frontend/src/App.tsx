import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from '@/context/AuthContext'
import { Layout } from '@/components/layout/Layout'
import { LoginPage } from '@/pages/auth/LoginPage'
import { DashboardPage } from '@/pages/dashboard/DashboardPage'
import { CustomersPage } from '@/pages/customers/CustomersPage'
import { VehiclesPage } from '@/pages/vehicles/VehiclesPage'
import { WorkOrdersPage } from '@/pages/workorders/WorkOrdersPage'
import { InventoryPage } from '@/pages/inventory/InventoryPage'
import { EmployeesPage } from '@/pages/employees/EmployeesPage'
import { POSPage } from '@/pages/pos/POSPage'
import { SuppliersPage } from '@/pages/suppliers/SuppliersPage'
import { ReportsPage } from '@/pages/reports/ReportsPage'
import { Loader2 } from 'lucide-react'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { token, loading } = useAuth()
  if (loading) return <div className="flex min-h-screen items-center justify-center"><Loader2 className="size-8 animate-spin text-muted-foreground" /></div>
  if (!token) return <Navigate to="/login" replace />
  return <>{children}</>
}

function PublicRoute({ children }: { children: React.ReactNode }) {
  const { token, loading } = useAuth()
  if (loading) return <div className="flex min-h-screen items-center justify-center"><Loader2 className="size-8 animate-spin text-muted-foreground" /></div>
  if (token) return <Navigate to="/" replace />
  return <>{children}</>
}

function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<PublicRoute><LoginPage /></PublicRoute>} />
      <Route path="/" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
        <Route index element={<DashboardPage />} />
        <Route path="customers" element={<CustomersPage />} />
        <Route path="vehicles" element={<VehiclesPage />} />
        <Route path="work-orders" element={<WorkOrdersPage />} />
        <Route path="inventory" element={<InventoryPage />} />
        <Route path="employees" element={<EmployeesPage />} />
        <Route path="pos" element={<POSPage />} />
        <Route path="suppliers" element={<SuppliersPage />} />
        <Route path="reports" element={<ReportsPage />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

function App() {
  return (
    <AuthProvider>
      <AppRoutes />
    </AuthProvider>
  )
}

export default App
