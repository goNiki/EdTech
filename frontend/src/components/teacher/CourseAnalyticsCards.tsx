'use client';

import React from 'react';
import { Users, TrendingUp, Award, Clock } from 'lucide-react';

interface AnalyticsProps {
  totalStudents: number;
  avgProgress: number;
  avgScore: number;
  pendingCount: number;
}

export default function CourseAnalyticsCards({
  totalStudents,
  avgProgress,
  avgScore,
  pendingCount,
}: AnalyticsProps) {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {/* 1. Total Students */}
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 p-5 rounded-3xl shadow-xs space-y-2">
        <div className="flex justify-between items-center text-slate-500 text-xs font-semibold">
          <span>Студентов на курсе</span>
          <div className="w-8 h-8 rounded-xl bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 flex items-center justify-center">
            <Users size={16} />
          </div>
        </div>
        <p className="text-2xl font-black text-slate-900 dark:text-white">{totalStudents}</p>
        <p className="text-[10px] text-slate-400">Активно обучаются</p>
      </div>

      {/* 2. Average Progress */}
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 p-5 rounded-3xl shadow-xs space-y-2">
        <div className="flex justify-between items-center text-slate-500 text-xs font-semibold">
          <span>Средний прогресс</span>
          <div className="w-8 h-8 rounded-xl bg-blue-50 dark:bg-blue-950 text-blue-600 dark:text-blue-400 flex items-center justify-center">
            <TrendingUp size={16} />
          </div>
        </div>
        <p className="text-2xl font-black text-slate-900 dark:text-white">{avgProgress}%</p>
        <div className="w-full bg-slate-100 dark:bg-slate-800 h-1.5 rounded-full overflow-hidden">
          <div className="bg-blue-500 h-full rounded-full" style={{ width: `${avgProgress}%` }} />
        </div>
      </div>

      {/* 3. Average Score */}
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 p-5 rounded-3xl shadow-xs space-y-2">
        <div className="flex justify-between items-center text-slate-500 text-xs font-semibold">
          <span>Средний балл тестов</span>
          <div className="w-8 h-8 rounded-xl bg-emerald-50 dark:bg-emerald-950 text-emerald-600 dark:text-emerald-400 flex items-center justify-center">
            <Award size={16} />
          </div>
        </div>
        <p className="text-2xl font-black text-emerald-600 dark:text-emerald-400">{avgScore} / 100</p>
        <p className="text-[10px] text-slate-400">По всем срезам знаний</p>
      </div>

      {/* 4. Pending Homeworks */}
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 p-5 rounded-3xl shadow-xs space-y-2">
        <div className="flex justify-between items-center text-slate-500 text-xs font-semibold">
          <span>Непроверенные ДЗ</span>
          <div className="w-8 h-8 rounded-xl bg-rose-50 dark:bg-rose-950 text-rose-600 dark:text-rose-400 flex items-center justify-center">
            <Clock size={16} />
          </div>
        </div>
        <div className="flex items-center gap-2">
          <p className="text-2xl font-black text-rose-600 dark:text-rose-400">{pendingCount}</p>
          {pendingCount > 0 && (
            <span className="w-2.5 h-2.5 rounded-full bg-rose-500 animate-pulse" />
          )}
        </div>
        <p className="text-[10px] text-slate-400">Требуют ручной проверки</p>
      </div>
    </div>
  );
}
