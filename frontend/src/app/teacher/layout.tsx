'use client';

import { usePathname } from 'next/navigation';
import ProtectedRoute from '@/components/ProtectedRoute';
import Sidebar from '@/components/layout/Sidebar';
import { useLayoutStore } from '@/store/useLayoutStore';

export default function TeacherLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const { isCollapsed } = useLayoutStore();

  // If inside full-screen editor (Puck edit mode), don't render sidebar
  if (pathname.includes('/edit')) {
    return <ProtectedRoute allowedRoles={['teacher', 'author', 'admin']}>{children}</ProtectedRoute>;
  }

  return (
    <ProtectedRoute allowedRoles={['teacher', 'author', 'admin']}>
      <div className="min-h-screen flex bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100">
        <Sidebar />
        <div className={`flex-1 ml-0 ${isCollapsed ? 'md:ml-20' : 'md:ml-64'} min-w-0 transition-all duration-300 overflow-x-hidden`}>
          {children}
        </div>
      </div>
    </ProtectedRoute>
  );
}
