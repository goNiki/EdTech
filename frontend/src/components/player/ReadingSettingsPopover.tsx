'use client';

import React, { useState, useEffect, useRef } from 'react';
import {
  Sun,
  Moon,
  Coffee,
  Check,
  AlignLeft,
  Maximize2,
  Minimize2,
  Sliders
} from 'lucide-react';
import { api } from '@/lib/api';

export interface ReadingSettings {
  fontScale: 'compact' | 'medium' | 'large' | 'xlarge';
  contentWidth: 'standard' | 'wide';
  lineHeight: 'normal' | 'relaxed';
  readingTheme: 'system' | 'light' | 'dark' | 'sepia';
}

export const DEFAULT_READING_SETTINGS: ReadingSettings = {
  fontScale: 'medium',
  contentWidth: 'standard',
  lineHeight: 'normal',
  readingTheme: 'system',
};

const STORAGE_KEY = 'edtech_reading_preferences';

export function loadReadingSettings(): ReadingSettings {
  if (typeof window === 'undefined') return DEFAULT_READING_SETTINGS;
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) {
      const parsed = JSON.parse(raw);
      return { ...DEFAULT_READING_SETTINGS, ...parsed };
    }
  } catch (e) {
    console.error('Failed to parse reading preferences from localStorage', e);
  }
  return DEFAULT_READING_SETTINGS;
}

export function saveReadingSettings(settings: ReadingSettings) {
  if (typeof window === 'undefined') return;
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
    // Background sync to server if authenticated
    api.patch('/auth/preferences', {
      font_scale: settings.fontScale,
      content_width: settings.contentWidth,
      line_height: settings.lineHeight,
      reading_theme: settings.readingTheme,
    }).catch(() => {
      // Ignore background sync errors gracefully
    });
  } catch (e) {
    console.error('Failed to save reading preferences', e);
  }
}

interface ReadingSettingsPopoverProps {
  settings: ReadingSettings;
  onChange: (newSettings: ReadingSettings) => void;
  className?: string;
}

export default function ReadingSettingsPopover({
  settings,
  onChange,
  className = '',
}: ReadingSettingsPopoverProps) {
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  // Close on outside click
  useEffect(() => {
    if (!isOpen) return;
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    document.addEventListener('keydown', handleKeyDown);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, [isOpen]);

  const updateSetting = <K extends keyof ReadingSettings>(key: K, value: ReadingSettings[K]) => {
    const next = { ...settings, [key]: value };
    onChange(next);
    saveReadingSettings(next);
  };

  const fontScales: Array<{ id: ReadingSettings['fontScale']; label: string; px: string }> = [
    { id: 'compact', label: '14', px: '14px' },
    { id: 'medium', label: '16', px: '16px' },
    { id: 'large', label: '18', px: '18px' },
    { id: 'xlarge', label: '20', px: '20px' },
  ];

  return (
    <div className={`relative ${className}`} ref={containerRef}>
      {/* Trigger Button "Aa" */}
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className={`px-3 py-1.5 rounded-xl border transition-all flex items-center gap-1.5 text-xs font-bold cursor-pointer active:scale-95 shadow-2xs ${
          isOpen
            ? 'bg-indigo-600 text-white border-indigo-600 ring-2 ring-indigo-500/20'
            : 'border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-indigo-600 dark:hover:text-indigo-400'
        }`}
        title="Настройки чтения (шрифт, ширина полотна, интервалы)"
        aria-label="Настройки чтения"
        aria-expanded={isOpen}
      >
        <span className="font-serif font-black text-sm tracking-tighter leading-none">Aa</span>
        <span className="hidden sm:inline text-[11px] font-semibold">Шрифт</span>
      </button>

      {/* Popover Dropdown */}
      {isOpen && (
        <div className="absolute right-0 top-full mt-2 w-72 sm:w-80 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-5 shadow-2xl z-50 animate-in fade-in zoom-in-95 duration-150 space-y-5 text-slate-900 dark:text-slate-100">
          <div className="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-3">
            <h4 className="text-xs font-black uppercase tracking-wider text-slate-500 dark:text-slate-400 flex items-center gap-1.5">
              <Sliders size={14} className="text-indigo-500" />
              <span>Параметры чтения</span>
            </h4>
            <span className="text-[10px] bg-indigo-50 dark:bg-indigo-950/60 text-indigo-600 dark:text-indigo-400 font-bold px-2 py-0.5 rounded-md">
              Автосохранение
            </span>
          </div>

          {/* 1. Font Scale Segmented Selector */}
          <div className="space-y-2">
            <div className="flex justify-between items-center text-xs">
              <span className="font-bold text-slate-700 dark:text-slate-300">Размер текста</span>
              <span className="font-mono text-[11px] text-slate-400">
                {fontScales.find((f) => f.id === settings.fontScale)?.px}
              </span>
            </div>

            <div className="grid grid-cols-4 gap-1.5 p-1 bg-slate-100 dark:bg-slate-800 rounded-2xl">
              {fontScales.map((item) => {
                const isActive = settings.fontScale === item.id;
                return (
                  <button
                    key={item.id}
                    type="button"
                    onClick={() => updateSetting('fontScale', item.id)}
                    className={`py-2 rounded-xl text-xs font-bold transition-all cursor-pointer flex flex-col items-center justify-center gap-0.5 ${
                      isActive
                        ? 'bg-white dark:bg-slate-900 text-indigo-600 dark:text-indigo-400 shadow-xs'
                        : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'
                    }`}
                  >
                    <span className="leading-none text-xs">{item.label}</span>
                    <span className="text-[9px] font-normal opacity-70">
                      {item.id === 'compact'
                        ? 'Мелкий'
                        : item.id === 'medium'
                        ? 'Стандарт'
                        : item.id === 'large'
                        ? 'Крупный'
                        : 'Макс'}
                    </span>
                  </button>
                );
              })}
            </div>
          </div>

          {/* 2. Container Width Selector */}
          <div className="space-y-2">
            <span className="font-bold text-xs text-slate-700 dark:text-slate-300 block">
              Ширина страницы
            </span>
            <div className="grid grid-cols-2 gap-2">
              <button
                type="button"
                onClick={() => updateSetting('contentWidth', 'standard')}
                className={`px-3 py-2.5 rounded-xl border text-xs font-bold transition-all cursor-pointer flex items-center justify-center gap-1.5 ${
                  settings.contentWidth === 'standard'
                    ? 'border-indigo-600 bg-indigo-50 dark:bg-indigo-950/40 text-indigo-700 dark:text-indigo-300 ring-1 ring-indigo-500/20'
                    : 'border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400'
                }`}
              >
                <Minimize2 size={13} />
                <span>Обычная (1024px)</span>
              </button>

              <button
                type="button"
                onClick={() => updateSetting('contentWidth', 'wide')}
                className={`px-3 py-2.5 rounded-xl border text-xs font-bold transition-all cursor-pointer flex items-center justify-center gap-1.5 ${
                  settings.contentWidth === 'wide'
                    ? 'border-indigo-600 bg-indigo-50 dark:bg-indigo-950/40 text-indigo-700 dark:text-indigo-300 ring-1 ring-indigo-500/20'
                    : 'border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400'
                }`}
              >
                <Maximize2 size={13} />
                <span>Широкая (1280px+)</span>
              </button>
            </div>
          </div>

          {/* 3. Line Height (Leading) Selector */}
          <div className="space-y-2">
            <span className="font-bold text-xs text-slate-700 dark:text-slate-300 block">
              Межстрочный интервал
            </span>
            <div className="grid grid-cols-2 gap-2">
              <button
                type="button"
                onClick={() => updateSetting('lineHeight', 'normal')}
                className={`px-3 py-2 rounded-xl border text-xs font-bold transition-all cursor-pointer flex items-center justify-center gap-1.5 ${
                  settings.lineHeight === 'normal'
                    ? 'border-indigo-600 bg-indigo-50 dark:bg-indigo-950/40 text-indigo-700 dark:text-indigo-300 ring-1 ring-indigo-500/20'
                    : 'border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400'
                }`}
              >
                <AlignLeft size={13} />
                <span>Стандарт (1.5)</span>
              </button>

              <button
                type="button"
                onClick={() => updateSetting('lineHeight', 'relaxed')}
                className={`px-3 py-2 rounded-xl border text-xs font-bold transition-all cursor-pointer flex items-center justify-center gap-1.5 ${
                  settings.lineHeight === 'relaxed'
                    ? 'border-indigo-600 bg-indigo-50 dark:bg-indigo-950/40 text-indigo-700 dark:text-indigo-300 ring-1 ring-indigo-500/20'
                    : 'border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400'
                }`}
              >
                <AlignLeft size={13} className="scale-y-125" />
                <span>Просторный (1.75)</span>
              </button>
            </div>
          </div>

          {/* 4. Reading Tone (Color theme) */}
          <div className="space-y-2 pt-1 border-t border-slate-100 dark:border-slate-800">
            <span className="font-bold text-xs text-slate-700 dark:text-slate-300 block">
              Фон чтения
            </span>
            <div className="grid grid-cols-3 gap-2">
              <button
                type="button"
                onClick={() => updateSetting('readingTheme', 'system')}
                className={`p-2 rounded-xl border text-xs font-bold transition-all cursor-pointer flex flex-col items-center gap-1 ${
                  settings.readingTheme === 'system'
                    ? 'border-indigo-600 ring-2 ring-indigo-500/20'
                    : 'border-slate-200 dark:border-slate-800'
                } bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-200`}
              >
                <Sun size={14} />
                <span className="text-[10px]">Системный</span>
              </button>

              <button
                type="button"
                onClick={() => updateSetting('readingTheme', 'sepia')}
                className={`p-2 rounded-xl border text-xs font-bold transition-all cursor-pointer flex flex-col items-center gap-1 ${
                  settings.readingTheme === 'sepia'
                    ? 'border-amber-600 ring-2 ring-amber-500/20'
                    : 'border-amber-200'
                } bg-[#FBF7EE] text-[#5B4636]`}
              >
                <Coffee size={14} />
                <span className="text-[10px]">Сепия</span>
              </button>

              <button
                type="button"
                onClick={() => updateSetting('readingTheme', 'dark')}
                className={`p-2 rounded-xl border text-xs font-bold transition-all cursor-pointer flex flex-col items-center gap-1 ${
                  settings.readingTheme === 'dark'
                    ? 'border-indigo-500 ring-2 ring-indigo-500/20'
                    : 'border-slate-700'
                } bg-slate-950 text-slate-200`}
              >
                <Moon size={14} />
                <span className="text-[10px]">Тёмный</span>
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
