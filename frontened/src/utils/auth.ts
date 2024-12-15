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

interface FetchOptions extends RequestInit {
  requiresAuth?: boolean;
}

export const authFetch = async (url: string, options: FetchOptions = {}) => {
  const { requiresAuth = true, ...fetchOptions } = options;
  const token = localStorage.getItem('token');

  // Create base headers
  const defaultHeaders: HeadersInit = {
    'Content-Type': 'application/json',
  };

  // Only add Authorization header if token exists and route requires auth
  if (token && requiresAuth) {
    defaultHeaders['Authorization'] = `Bearer ${token}`;
  }

  // Don't set Content-Type if we're sending FormData
  if (options.body instanceof FormData) {
    delete defaultHeaders['Content-Type'];
  }

  const headers = {
    ...defaultHeaders,
    ...options.headers,
  };

  try {
    console.log('Making request to:', `${API_URL}${url}`);
    console.log('With headers:', headers);

    const response = await fetch(`${API_URL}${url}`, {
      ...fetchOptions,
      headers,
      credentials: 'include',
      mode: 'cors',
    });

    // Only redirect to login if the route requires auth and we get a 401
    if (response.status === 401 && requiresAuth) {
      console.log('Unauthorized access detected - clearing auth state');
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
      throw new Error('Unauthorized');
    }

    return response;
  } catch (error) {
    console.error('Fetch Error:', error);
    if (error instanceof Error && error.message === 'Failed to fetch') {
      console.error('Network error - API might be down or CORS issue');
    }
    throw error;
  }
};

export function useAuthState() {
  const [user, setUser] = useState<User | null>(() => {
    const savedUser = localStorage.getItem('user');
    return savedUser ? JSON.parse(savedUser) : null;
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const login = useCallback(async (credentials: LoginCredentials) => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await authFetch('/login', {
        method: 'POST',
        body: JSON.stringify(credentials),
        requiresAuth: false // Login endpoint is public
      });
  
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Login failed');
      }
  
      const data: AuthResponse = await response.json();
      console.log("Login successful:", data);
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
        requiresAuth: false // Register endpoint is public
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Registration failed');
      }

      const data: AuthResponse = await response.json();
      console.log("Registration successful:", data);
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
    window.location.href = '/';
  }, []);

  const checkAuth = useCallback(async () => {
    const token = localStorage.getItem('token');
    if (!token) {
      setUser(null);
      return false;
    }

    try {
      setLoading(true);
      const response = await authFetch('/verify-token', { requiresAuth: true });

      if (!response.ok) {
        throw new Error('Token verification failed');
      }

      const { user: userData } = await response.json();
      setUser(userData);
      localStorage.setItem('user', JSON.stringify(userData));
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Authentication check failed';
      setError(errorMessage);
      logout();
      return false;
    } finally {
      setLoading(false);
    }
  }, [logout]);

  const isAuthenticated = useCallback(() => {
    return !!user && !!localStorage.getItem('token');
  }, [user]);

  return {
    user,
    loading,
    error,
    login,
    register,
    logout,
    checkAuth,
    isAuthenticated
  };
}