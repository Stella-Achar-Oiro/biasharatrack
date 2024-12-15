import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from './homepage/Home';
import Login from './pages/Login';
import SignUp from './pages/SignUp';
import NotFound from './pages/error/NotFound';
import Unauthorized from './pages/error/Unauthorized';
import ServerError from './pages/error/ServerError';
import ErrorBoundary from './homepage/components/ErrorBoundary';
import FloatingLanguageSwitcher from './homepage/components/FloatingLanguageSwitcher';
import InventoryManagement from './pages/inventory/InventoryManagement';
import SalesManagement from './pages/sales/SalesManagement';
import CreditsManagement from './pages/credits/CreditsManagement';
import Dashboard from './pages/Dashboard';
import ProtectedRoute from './context/ProtectedRoute';
import ReceiptManagement from './pages/receipts/ReceiptManagement';
import { useAuthState } from './utils/auth';
import { useEffect } from 'react';


function App() {
  const { checkAuth } = useAuthState();

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  return (
    <ErrorBoundary>
      <Router>
        <FloatingLanguageSwitcher />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<SignUp />} />

          <Route path="/unauthorized" element={<Unauthorized />} />
          <Route path="/server-error" element={<ServerError />} />
          <Route path="*" element={<NotFound />} />

          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/dashboard/inventory"
            element={
              <ProtectedRoute>
                <InventoryManagement />
              </ProtectedRoute>
            }
          />
          <Route
            path="/dashboard/sales"
            element={
              <ProtectedRoute>
                <SalesManagement />
              </ProtectedRoute>
            }
          />

          <Route
            path="/dashboard/credits"
            element={
              <ProtectedRoute>
                <CreditsManagement />
              </ProtectedRoute>
            }
          />

          <Route
            path="/dashboard/receipts"
            element={
              <ProtectedRoute>
                <ReceiptManagement />
              </ProtectedRoute>
            }
          />
          <Route path="/unauthorized" element={<Unauthorized />} />
          <Route path="/server-error" element={<ServerError />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </Router>
    </ErrorBoundary>
  );
}

export default App;