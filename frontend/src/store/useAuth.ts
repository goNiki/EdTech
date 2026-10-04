import { create } from 'zustand';
import { api } from '@/lib/api';

export interface User {
  id: number;
  email: string;
  first_name?: string;
  last_name?: string;
  username: string;
  avatar_url?: string;
  headline?: string;
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
  login: (accessToken: string, refreshToken: string) => Promise<void>;
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

  login: async (accessToken, refreshToken) => {
    localStorage.setItem('access_token', accessToken);
    localStorage.setItem('refresh_token', refreshToken);
    set({ isAuthenticated: true });
    await get().fetchUser();
    const currentUser = get().user;
    if (currentUser && ['teacher', 'author', 'admin'].includes(currentUser.role)) {
      get().setViewMode('teacher');
    } else {
      get().setViewMode('student');
    }
  },

  logout: () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('view_mode');
    set({ user: null, isAuthenticated: false, isLoading: false, viewMode: 'student' });
  },

  setViewMode: (mode: ViewMode) => {
    const currentUser = get().user;
    const isStaff = currentUser && ['teacher', 'author', 'admin'].includes(currentUser.role);
    const finalMode = isStaff ? mode : 'student';
    localStorage.setItem('view_mode', finalMode);
    set({ viewMode: finalMode });
  },

  setUser: (user: User) => {
    set({ user });
  },

  fetchUser: async () => {
    try {
      set({ isLoading: true });
      const { data } = await api.get('/auth/me');
      const userData = data.data || data;
      
      let initialViewMode: ViewMode = 'student';
      if (typeof window !== 'undefined') {
        const isStaff = ['teacher', 'author', 'admin'].includes(userData.role);
        if (isStaff) {
          const savedMode = localStorage.getItem('view_mode') as ViewMode;
          initialViewMode = (savedMode === 'student' || savedMode === 'teacher') ? savedMode : 'teacher';
        } else {
          initialViewMode = 'student';
          localStorage.setItem('view_mode', 'student');
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
