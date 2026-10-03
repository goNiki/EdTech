'use client';

import ProtectedRoute from '@/components/ProtectedRoute';
import Sidebar from '@/components/layout/Sidebar';

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <ProtectedRoute allowedRoles={['student', 'teacher', 'admin', 'author']}>
      <div className="min-h-screen flex bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100">
        <Sidebar />
        <div className="flex-1 ml-64 min-w-0 transition-all duration-300">
          {children}
        </div>
      </div>
    </ProtectedRoute>
  );
}
