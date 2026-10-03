import { create } from 'zustand';
import { api } from '@/lib/api';

export interface User {
  id: number;
  email: string;
  first_name?: string;
  last_name?: string;
  username: string;
  avatar_url?: string;
  bio?: string;
  role: string;
  email_verified?: boolean;
  is_active?: boolean;
  created_at?: string;
  last_login_at?: string;
}

export type ViewMode = 'student' | 'teacher';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  viewMode: ViewMode;
  login: (accessToken: string, refreshToken: string) => void;
  logout: () => void;
  fetchUser: () => Promise<void>;
  setViewMode: (mode: ViewMode) => void;
  setUser: (user: User) => void;
}

export const useAuth = create<AuthState>((set, get) => ({
  user: null,
  isAuthenticated: false,
  isLoading: true,
  viewMode: (typeof window !== 'undefined' && (localStorage.getItem('view_mode') as ViewMode)) || 'student',

  login: (accessToken, refreshToken) => {
    localStorage.setItem('access_token', accessToken);
    localStorage.setItem('refresh_token', refreshToken);
    set({ isAuthenticated: true });
    get().fetchUser();
  },

  logout: () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    set({ user: null, isAuthenticated: false, isLoading: false, viewMode: 'student' });
  },

  setViewMode: (mode: ViewMode) => {
    localStorage.setItem('view_mode', mode);
    set({ viewMode: mode });
  },

  setUser: (user: User) => {
    set({ user });
  },

  fetchUser: async () => {
    try {
      set({ isLoading: true });
      const { data } = await api.get('/auth/me');
      const userData = data.data || data;
      
      let initialViewMode = get().viewMode;
      if (typeof window !== 'undefined') {
        const savedMode = localStorage.getItem('view_mode') as ViewMode;
        if (savedMode) {
          initialViewMode = savedMode;
        } else if (['teacher', 'author', 'admin'].includes(userData.role)) {
          initialViewMode = 'teacher';
          localStorage.setItem('view_mode', 'teacher');
        }
      }

      set({ 
        user: userData, 
        isAuthenticated: true, 
        isLoading: false,
        viewMode: initialViewMode
      });
    } catch (error) {
      set({ user: null, isAuthenticated: false, isLoading: false });
    }
  },
}));
