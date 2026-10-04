'use client';

import React, { useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import { Maximize2, Minimize2 } from 'lucide-react';

interface FocusModeToggleProps {
  isFocusMode: boolean;
  onToggle: () => void;
  className?: string;
}

export default function FocusModeToggle({
  isFocusMode,
  onToggle,
  className = '',
}: FocusModeToggleProps) {
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  // Global keyboard shortcuts: 'F' / 'f' to toggle, 'Escape' to exit
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Don't trigger if user is typing in input or textarea
      const target = e.target as HTMLElement | null;
      const isInput =
        target?.tagName === 'INPUT' ||
        target?.tagName === 'TEXTAREA' ||
        target?.isContentEditable;

      if (isInput) return;

      if (e.key === 'f' || e.key === 'F') {
        e.preventDefault();
        onToggle();
      } else if (e.key === 'Escape' && isFocusMode) {
        e.preventDefault();
        onToggle();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isFocusMode, onToggle]);

  return (
    <>
      {/* Standard Header Button */}
      <button
        type="button"
        onClick={onToggle}
        className={`p-2 rounded-xl border transition-all flex items-center gap-1.5 text-xs font-bold cursor-pointer active:scale-95 shadow-2xs ${
          isFocusMode
            ? 'bg-indigo-600 text-white border-indigo-600 shadow-md ring-2 ring-indigo-500/20'
            : 'border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-indigo-600 dark:hover:text-indigo-400'
        } ${className}`}
        title={isFocusMode ? 'Выйти из режима фокуса (Esc или F)' : 'Режим глубокой концентрации (клавиша F)'}
        aria-label={isFocusMode ? 'Обычный режим' : 'Режим фокуса'}
      >
        {isFocusMode ? <Minimize2 size={16} /> : <Maximize2 size={16} />}
        <span className="hidden lg:inline">
          {isFocusMode ? 'Фокус вкл' : 'Фокус'}
        </span>
      </button>

      {/* Floating Exit Pill in Focus Mode rendered to body */}
      {mounted && isFocusMode && createPortal(
        <div className="fixed top-4 right-4 z-50 flex items-center gap-2 animate-in fade-in slide-in-from-top-2 duration-300">
          <button
            type="button"
            onClick={onToggle}
            className="px-4 py-2 rounded-full bg-slate-900/90 hover:bg-slate-900 dark:bg-slate-100/90 dark:hover:bg-white text-white dark:text-slate-900 backdrop-blur-md shadow-2xl border border-white/20 dark:border-black/20 flex items-center gap-2 text-xs font-extrabold transition-all cursor-pointer opacity-80 hover:opacity-100 active:scale-95 ring-2 ring-indigo-500/30"
            title="Выйти из режима фокуса (Esc или F)"
          >
            <Minimize2 size={14} className="text-indigo-400 dark:text-indigo-600" />
            <span>Выйти из фокуса (Esc)</span>
          </button>
        </div>,
        document.body
      )}
    </>
  );
}
