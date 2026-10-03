'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useAuth } from '@/store/useAuth';
import {
  BookOpen,
  GraduationCap,
  PlusCircle,
  Settings,
  User as UserIcon,
  LogOut,
  ChevronLeft,
  ChevronRight,
  CheckSquare,
  Compass,
  ArrowLeft
} from 'lucide-react';

export default function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const { user, viewMode, setViewMode, logout } = useAuth();
  const [collapsed, setCollapsed] = useState(false);
  const [isLogoutModalOpen, setIsLogoutModalOpen] = useState(false);

  const canSwitchRole = ['teacher', 'author', 'admin'].includes(user?.role || '');

  const handleRoleChange = (mode: 'student' | 'teacher') => {
    setViewMode(mode);
    if (mode === 'teacher') {
      router.push('/teacher/courses');
    } else {
      router.push('/courses');
    }
  };

  const isTeacherView = viewMode === 'teacher' && canSwitchRole;

  const studentNavItems = [
    { label: 'Все курсы', href: '/courses', icon: Compass },
    { label: 'Мои курсы', href: '/dashboard/courses', icon: BookOpen },
    { label: 'Личный кабинет', href: '/dashboard/settings', icon: UserIcon },
    { label: 'Настройки', href: '/dashboard/settings?tab=preferences', icon: Settings },
  ];

  const teacherNavItems = [
    { label: 'Мои курсы (Автор)', href: '/teacher/courses', icon: BookOpen },
    { label: 'Создать новый курс', href: '/teacher/courses/new', icon: PlusCircle },
    { label: 'Проверка заданий', href: '/teacher/grading', icon: CheckSquare },
    { label: 'Личный кабинет', href: '/dashboard/settings', icon: UserIcon },
    { label: 'Настройки', href: '/dashboard/settings?tab=preferences', icon: Settings },
  ];

  const navItems = isTeacherView ? teacherNavItems : studentNavItems;

  return (
    <>
      <aside
        className={`bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 flex flex-col justify-between transition-all duration-300 z-30 fixed h-full shadow-sm ${
          collapsed ? 'w-20' : 'w-64'
        }`}
      >
        <div>
          {/* Brand Header */}
          <div className="h-16 px-4 flex items-center justify-between border-b border-slate-100 dark:border-slate-800">
            <Link href="/" className="flex items-center gap-3 overflow-hidden">
              <div className="w-9 h-9 rounded-xl bg-gradient-to-tr from-indigo-600 via-indigo-500 to-cyan-400 flex items-center justify-center text-white font-extrabold text-lg flex-shrink-0 shadow-sm">
                ED
              </div>
              {!collapsed && (
                <div className="flex flex-col whitespace-nowrap">
                  <span className="font-extrabold text-base tracking-tight text-slate-900 dark:text-white">
                    ED.Learn
                  </span>
                  <span className="text-[9px] uppercase font-bold tracking-widest text-indigo-500">
                    {isTeacherView ? 'Преподаватель' : 'Студент'}
                  </span>
                </div>
              )}
            </Link>
            <button
              onClick={() => setCollapsed(!collapsed)}
              className="p-1.5 rounded-lg text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
              title={collapsed ? 'Развернуть' : 'Свернуть'}
            >
              {collapsed ? <ChevronRight size={18} /> : <ChevronLeft size={18} />}
            </button>
          </div>

          {/* Role Switcher Banner (for teachers/authors/admins) */}
          {canSwitchRole && !collapsed && (
            <div className="p-3">
              <div className="bg-indigo-50/80 dark:bg-slate-800/80 border border-indigo-100 dark:border-slate-700/80 rounded-2xl p-2.5">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-[10px] font-bold uppercase tracking-wider text-indigo-600 dark:text-indigo-400">
                    Режим просмотра
                  </span>
                  <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
                </div>
                <div className="grid grid-cols-2 gap-1 bg-slate-200/60 dark:bg-slate-900 p-1 rounded-xl">
                  <button
                    onClick={() => handleRoleChange('student')}
                    className={`py-1.5 px-2 text-xs font-semibold rounded-lg transition-all text-center ${
                      !isTeacherView
                        ? 'bg-white dark:bg-indigo-600 text-indigo-600 dark:text-white shadow-sm'
                        : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'
                    }`}
                  >
                    Студент
                  </button>
                  <button
                    onClick={() => handleRoleChange('teacher')}
                    className={`py-1.5 px-2 text-xs font-semibold rounded-lg transition-all text-center ${
                      isTeacherView
                        ? 'bg-white dark:bg-indigo-600 text-indigo-600 dark:text-white shadow-sm'
                        : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'
                    }`}
                  >
                    Учитель
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* Navigation Links */}
          <nav className="px-3 space-y-1.5 mt-2">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = pathname === item.href || (item.href !== '/courses' && item.href !== '/teacher/courses' && pathname.startsWith(item.href));

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`flex items-center gap-3 px-3.5 py-2.5 rounded-2xl text-sm font-semibold transition-all ${
                    isActive
                      ? 'bg-indigo-50 dark:bg-indigo-950/60 text-indigo-600 dark:text-indigo-400 shadow-sm ring-1 ring-indigo-100 dark:ring-indigo-900/40'
                      : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-slate-900 dark:hover:text-white'
                  }`}
                  title={collapsed ? item.label : undefined}
                >
                  <Icon size={18} className="flex-shrink-0" />
                  {!collapsed && <span className="truncate">{item.label}</span>}
                </Link>
              );
            })}
          </nav>
        </div>

        {/* Footer Profile & Logout */}
        <div className="p-3 border-t border-slate-100 dark:border-slate-800 mt-auto space-y-2">
          {!collapsed ? (
            <div className="flex items-center gap-3 p-2 rounded-2xl bg-slate-50 dark:bg-slate-800/50">
              <div className="w-9 h-9 rounded-full bg-gradient-to-tr from-indigo-500 to-purple-600 text-white flex items-center justify-center font-bold text-xs flex-shrink-0 shadow-sm">
                {user?.first_name?.[0] || user?.username?.[0] || 'U'}
                {user?.last_name?.[0] || ''}
              </div>
              <div className="flex flex-col min-w-0 flex-1">
                <span className="text-xs font-bold text-slate-800 dark:text-slate-200 truncate">
                  {user?.first_name ? `${user.first_name} ${user.last_name || ''}` : user?.username}
                </span>
                <span className="text-[10px] text-slate-400 truncate">{user?.email}</span>
              </div>
            </div>
          ) : (
            <div className="w-9 h-9 mx-auto rounded-full bg-indigo-100 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 flex items-center justify-center font-bold text-xs shadow-sm">
              {user?.first_name?.[0] || user?.username?.[0] || 'U'}
            </div>
          )}

          <button
            onClick={() => setIsLogoutModalOpen(true)}
            className={`w-full flex items-center gap-3 px-3.5 py-2 rounded-xl text-xs font-bold text-rose-600 dark:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/50 transition-colors ${
              collapsed ? 'justify-center' : ''
            }`}
            title={collapsed ? 'Выйти' : undefined}
          >
            <LogOut size={16} className="flex-shrink-0" />
            {!collapsed && <span>Выйти</span>}
          </button>
        </div>
      </aside>

      {/* Logout Confirmation Modal */}
      {isLogoutModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4">
          <div className="bg-white dark:bg-slate-900 rounded-3xl p-6 max-w-sm w-full border border-slate-200 dark:border-slate-800 shadow-2xl space-y-4">
            <h3 className="text-lg font-bold text-slate-900 dark:text-white">Выход из аккаунта</h3>
            <p className="text-xs text-slate-500">Вы уверены, что хотите завершить сессию в ED.Learn?</p>
            <div className="flex gap-2 pt-2">
              <button
                onClick={() => setIsLogoutModalOpen(false)}
                className="flex-1 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors"
              >
                Отмена
              </button>
              <button
                onClick={() => {
                  setIsLogoutModalOpen(false);
                  logout();
                  router.push('/login');
                }}
                className="flex-1 py-2.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold transition-colors shadow-md"
              >
                Выйти
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
