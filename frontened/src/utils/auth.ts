import { useState, useCallback } from 'react';
import { API_URL } from './api';

export interface User {
  id: string;
  full_name: string;
  email: string;
  imageUrl?: string;
  role: 'owner' | 'manager' | 'staff';
  business_name: string;
  telephone: string;
  location: string;
  token: string;
}
interface LoginCredentials {
  email: string;
  password: string;
}

interface RegisterCredentials {
  fullName: string;
  email: string;
  password: string;
}

interface AuthResponse {
  message: string;
  token: string;
  user: User;
}

export const authFetch = async (url: string, options: RequestInit = {}) => {
  const token = localStorage.getItem('token');

  console.log('Token:', token);
  console.log('URL:', url);

  const headers = {
    ...options.headers,
    'Authorization': `Bearer ${token}`,
  };



  try {
    const response = await fetch(`${API_URL}${url}`, {
      ...options,
      headers,
    });
  

  
    if (response.status === 401) {
      console.log('Unauthorized access detected');
      localStorage.removeItem('token');
      window.location.href = '/login';
      throw new Error('Unauthorized');
    }


    return response;
  } catch (error) {
    console.error('Fetch Error:', error);
    throw error;
  }
};
export function useAuthState() {
  const [user, setUser] = useState<User | null>(null);
  // Changed initial loading to false since we're not doing any initial fetch
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  


  const login = useCallback(async (credentials: LoginCredentials) => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await authFetch('/login', {
        method: 'POST',
        body: JSON.stringify(credentials),
      });
  
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Login failed');
      }
  
      const data: AuthResponse = await response.json();
      console.log("User data",data);
      setUser(data.user);
      localStorage.setItem('token', data.token);
      localStorage.setItem('user', JSON.stringify(data.user));
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Login failed. Please try again.';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const register = useCallback(async (credentials: RegisterCredentials): Promise<void> => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await authFetch('/register', {
        method: 'POST',
        body: JSON.stringify(credentials),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Registration failed');
      }

      const data: AuthResponse = await response.json();
      console.log("User data",data);
      setUser(data.user);
      localStorage.setItem('token', data.token);
      localStorage.setItem('user', JSON.stringify(data.user));
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Registration failed. Please try again.';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const logout = useCallback(() => {
    setUser(null);
    setError(null);
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }, []);

  const checkAuth = useCallback(async () => {
    const token = localStorage.getItem('token');
    if (!token) {
      setUser(null);
      return;
    }

    try {
      setLoading(true);
      const response = await authFetch('/verify-token');

      if (!response.ok) {
        throw new Error('Token verification failed');
      }

      const { user: userData } = await response.json();
      setUser(userData);
      localStorage.setItem('user', JSON.stringify(userData));
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Authentication check failed';
      setError(errorMessage);
      // Clear everything if token verification fails
      logout();
    } finally {
      setLoading(false);
    }
  }, [logout]);

  return {
    user,
    loading,
    error,
    login,
    register,
    logout,
    checkAuth
  };
}