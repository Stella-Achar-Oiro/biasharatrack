// src/context/ProtectedRoute.tsx
import { Navigate, Outlet } from 'react-router-dom';
import { useEffect } from 'react';
import { useAuthState } from '../utils/auth';

const ProtectedRoute = () => {
  const { isAuthenticated, checkAuth, loading } = useAuthState();

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  if (loading) {
    return <div>Loading...</div>; // Or your loading component
  }

  return isAuthenticated() ? <Outlet /> : <Navigate to="/login" />;
};

export default ProtectedRoute;