'use client';

import ProtectedRoute from '@/components/ProtectedRoute';
import Sidebar from '@/components/layout/Sidebar';
import { useLayoutStore } from '@/store/useLayoutStore';

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const { isCollapsed } = useLayoutStore();

  return (
    <ProtectedRoute allowedRoles={['student', 'teacher', 'admin', 'author']}>
      <div className="min-h-screen flex bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100">
        <Sidebar />
        <div className={`flex-1 ml-0 ${isCollapsed ? 'md:ml-20' : 'md:ml-64'} min-w-0 transition-all duration-300 overflow-x-hidden`}>
          {children}
        </div>
      </div>
    </ProtectedRoute>
  );
}
