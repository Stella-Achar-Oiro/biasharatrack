import { createContext, useContext, ReactNode, useEffect, useState } from 'react';
import { useAuthState } from '../utils/auth';
import { User } from '../../types/user';  // Import the User interface

interface RegisterData {
  fullName: string;
  email: string;
  password: string;
  businessName: string;
  telephone: string;
  location: string;
}

interface AuthContextType {
  isAuthenticated: boolean;
  user: User | null;
  loading: boolean;
  error: string | null;
  login: (credentials: { email: string; password: string }) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => void;
  checkAuth: () => void;
  getToken: () => string | null;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [initialized, setInitialized] = useState(false);
  const { 
    user,
    loading = true,
    error,
    login,
    register,
    logout,
    checkAuth
  } = useAuthState();

  const getToken = () => localStorage.getItem('token');

  useEffect(() => {
    const token = getToken();
    if (token) {
      checkAuth();
    }
    setInitialized(true);
  }, [checkAuth]);

  const value = {
    isAuthenticated: !!user,
    user,
    loading,
    error,
    login,
    register,
    logout,
    checkAuth,
    getToken
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};