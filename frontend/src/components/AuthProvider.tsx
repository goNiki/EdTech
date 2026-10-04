'use client';

import React, { useEffect } from 'react';
import { useAuth } from '@/store/useAuth';

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const fetchUser = useAuth((state) => state.fetchUser);

  useEffect(() => {
    if (typeof window !== 'undefined' && localStorage.getItem('access_token')) {
      fetchUser();
    }
  }, [fetchUser]);

  return <>{children}</>;
}
