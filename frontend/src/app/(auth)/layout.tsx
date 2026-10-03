import React from 'react';
import Link from 'next/link';
import { ArrowLeft } from 'lucide-react';

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen flex flex-col justify-between text-slate-800 dark:text-slate-100 bg-slate-50 dark:bg-slate-950 antialiased relative overflow-x-hidden selection:bg-indigo-500 selection:text-white font-sans transition-colors duration-200">
      {/* Top Bar */}
      <header className="relative z-10 w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 flex items-center justify-between">
        <Link href="/" className="flex items-center gap-3 group">
          <div className="w-11 h-11 rounded-2xl bg-gradient-to-tr from-indigo-600 via-indigo-500 to-cyan-400 flex items-center justify-center text-white font-extrabold text-xl shadow-lg shadow-indigo-500/25 group-hover:scale-105 transition-transform">
            ED
          </div>
          <div className="flex flex-col">
            <span className="font-black text-xl tracking-tight text-slate-900 dark:text-white group-hover:text-indigo-600 transition-colors">
              ED.Learn
            </span>
            <span className="text-[10px] uppercase font-bold tracking-widest text-indigo-500">
              Курсы & Тестирование
            </span>
          </div>
        </Link>

        <Link
          href="/"
          className="inline-flex items-center gap-2 text-xs sm:text-sm font-semibold text-slate-600 dark:text-slate-300 hover:text-indigo-600 bg-white dark:bg-slate-900 hover:bg-slate-100 dark:hover:bg-slate-800 border border-slate-200/80 dark:border-slate-800 px-4 py-2.5 rounded-xl shadow-xs transition-all"
        >
          <ArrowLeft size={16} />
          <span>На главную</span>
        </Link>
      </header>

      {/* Main Content Container */}
      <main className="relative z-10 flex-1 flex items-center justify-center px-4 sm:px-6 lg:px-8 py-4 sm:py-8">
        {children}
      </main>

      {/* Minimal Footer */}
      <footer className="relative z-10 py-6 text-center text-xs text-slate-400 dark:text-slate-500">
        &copy; 2026 ED.Learn Platform. Все права защищены.
      </footer>
    </div>
  );
}
