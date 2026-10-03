'use client';

import { useAuth } from '@/store/useAuth';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

export default function ProtectedRoute({ 
  children, 
  allowedRoles 
}: { 
  children: React.ReactNode;
  allowedRoles?: string[];
}) {
  const { user, isAuthenticated, isLoading, fetchUser } = useAuth();
  const router = useRouter();
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    const checkAuth = async () => {
      if (!isAuthenticated) {
        await fetchUser();
      }
      setIsChecking(false);
    };
    checkAuth();
  }, [isAuthenticated, fetchUser]);

  useEffect(() => {
    if (!isChecking && !isLoading && !isAuthenticated) {
      router.push('/login');
    }
  }, [isChecking, isLoading, isAuthenticated, router]);

  if (isChecking || isLoading) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-50">
        <div className="animate-pulse flex flex-col items-center space-y-4">
          <div className="h-12 w-12 bg-gray-300 rounded-full"></div>
          <div className="h-4 w-32 bg-gray-300 rounded"></div>
        </div>
      </div>
    );
  }

  if (allowedRoles && user) {
    const hasRole = allowedRoles.includes(user.role || 'student');
    if (!hasRole) {
      return (
        <div className="flex h-screen items-center justify-center">
          <div className="text-center space-y-4">
            <h1 className="text-2xl font-bold">Доступ запрещен</h1>
            <p className="text-muted-foreground">У вас нет прав для просмотра этой страницы.</p>
          </div>
        </div>
      );
    }
  }

  return isAuthenticated ? <>{children}</> : null;
}
