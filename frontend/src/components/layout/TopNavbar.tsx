'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import { useTheme } from 'next-themes';
import { ArrowLeft, Menu, Moon, Sun } from 'lucide-react';
import { useAuth } from '@/store/useAuth';
import { useLayoutStore } from '@/store/useLayoutStore';
import NotificationsDropdown from './NotificationsDropdown';

interface TopNavbarProps {
  title?: string;
  subtitle?: string;
  showBack?: boolean;
}

export default function TopNavbar({ title, subtitle, showBack = true }: TopNavbarProps) {
  const router = useRouter();
  const { theme, setTheme } = useTheme();
  const { viewMode } = useAuth();
  const { toggleMobileMenu } = useLayoutStore();

  return (
    <header className="h-16 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200 dark:border-slate-800 sticky top-0 z-20 px-4 sm:px-6 flex items-center justify-between">
      <div className="flex items-center gap-3 sm:gap-4 min-w-0">
        {/* Mobile Hamburger Menu Toggle */}
        <button
          type="button"
          onClick={toggleMobileMenu}
          className="p-2 rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors shadow-xs md:hidden cursor-pointer flex-shrink-0"
          title="Открыть меню"
          aria-label="Открыть мобильное меню"
        >
          <Menu size={18} />
        </button>

        {showBack && (
          <button
            onClick={() => router.back()}
            className="p-2 rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors shadow-xs cursor-pointer flex-shrink-0"
            title="Назад"
          >
            <ArrowLeft size={16} />
          </button>
        )}
        {title && (
          <div className="min-w-0">
            <h1 className="text-sm sm:text-base font-extrabold text-slate-900 dark:text-white tracking-tight leading-none truncate">
              {title}
            </h1>
            {subtitle && <p className="text-[11px] sm:text-xs text-slate-500 mt-0.5 truncate">{subtitle}</p>}
          </div>
        )}
      </div>

      <div className="flex items-center gap-2 sm:gap-3 flex-shrink-0">
        {/* Notifications Dropdown */}
        <NotificationsDropdown />

        {/* Theme Toggle Button */}
        <button
          onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
          className="p-2 rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors shadow-xs cursor-pointer"
          title="Сменить тему"
        >
          {theme === 'dark' ? <Sun size={16} className="text-amber-400" /> : <Moon size={16} />}
        </button>
      </div>
    </header>
  );
}
