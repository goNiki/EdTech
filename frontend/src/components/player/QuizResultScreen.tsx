'use client';

import React from 'react';
import {
  Award,
  CheckCircle2,
  XCircle,
  RotateCcw,
  ArrowRight,
  ShieldCheck,
  ChevronLeft,
  Sparkles,
  TrendingUp,
  Clock,
  ListOrdered
} from 'lucide-react';

export interface QuizResultScreenProps {
  score: number;
  earnedPoints: number;
  totalMaxPoints: number;
  passingScore?: number;
  bestScore?: number;
  remainingAttempts?: number;
  nextLesson?: { id: number; title: string } | null;
  onRetake?: () => void;
  onNextLesson?: () => void;
  onNavigateBack: () => void;
}

export default function QuizResultScreen({
  score,
  earnedPoints,
  totalMaxPoints,
  passingScore = 70,
  bestScore,
  remainingAttempts = 2,
  nextLesson,
  onRetake,
  onNextLesson,
  onNavigateBack,
}: QuizResultScreenProps) {
  const isPassed = score >= passingScore;
  const highestScore = Math.max(score, bestScore || score);

  return (
    <div className="max-w-2xl mx-auto w-full py-8 space-y-6 animate-in fade-in duration-300">
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 sm:p-8 shadow-2xl space-y-6 text-center">
        {/* Animated Icon */}
        <div
          className={`w-20 h-20 rounded-3xl flex items-center justify-center mx-auto shadow-inner transition-transform ${
            isPassed
              ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-600 dark:text-emerald-400'
              : 'bg-amber-100 dark:bg-amber-950 text-amber-600 dark:text-amber-400'
          }`}
        >
          {isPassed ? <Award size={44} /> : <RotateCcw size={40} />}
        </div>

        {/* Title & Status */}
        <div className="space-y-2">
          <div
            className={`inline-flex items-center gap-1.5 px-3.5 py-1 rounded-full text-xs font-black uppercase tracking-wider border ${
              isPassed
                ? 'bg-emerald-50 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-300 border-emerald-200 dark:border-emerald-800/60'
                : 'bg-amber-50 dark:bg-amber-950/60 text-amber-700 dark:text-amber-300 border-amber-200 dark:border-amber-800/60'
            }`}
          >
            <ShieldCheck size={14} />
            <span>{isPassed ? 'Тест успешно сдан!' : 'Порог не пройден'}</span>
          </div>

          <h2 className="text-2xl sm:text-3xl font-black text-slate-900 dark:text-white">
            {isPassed ? 'Отличная работа! 🎉' : 'Попробуйте еще раз! 💪'}
          </h2>
          <p className="text-xs sm:text-sm text-slate-500 dark:text-slate-400">
            {isPassed
              ? 'Ваш результат зафиксирован и учтен в общем прогрессе курса.'
              : `Для успешной сдачи необходимо набрать от ${passingScore}% баллов.`}
          </p>
        </div>

        {/* Big Score Card */}
        <div className="p-6 bg-slate-50 dark:bg-slate-800/50 rounded-3xl border border-slate-200 dark:border-slate-800 space-y-4">
          <div className="flex flex-col items-center justify-center">
            <div className="text-5xl sm:text-6xl font-black text-indigo-600 dark:text-indigo-400 tracking-tight">
              {score}%
            </div>
            <div className="text-xs sm:text-sm font-bold text-slate-500 dark:text-slate-400 mt-1">
              {earnedPoints} из {totalMaxPoints} баллов набрано
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3 pt-3 border-t border-slate-200/60 dark:border-slate-700/60 text-left text-xs">
            <div className="space-y-0.5">
              <span className="text-[10px] uppercase font-bold text-slate-400">
                Проходной порог
              </span>
              <p className="font-bold text-slate-700 dark:text-slate-300">
                {passingScore}%{' '}
                <span className={isPassed ? 'text-emerald-500 font-extrabold' : 'text-rose-500 font-extrabold'}>
                  ({score >= passingScore ? `+${score - passingScore}%` : `-${passingScore - score}%`})
                </span>
              </p>
            </div>

            <div className="space-y-0.5">
              <span className="text-[10px] uppercase font-bold text-slate-400">
                Высший результат
              </span>
              <p className="font-bold text-slate-700 dark:text-slate-300">
                {highestScore}%
              </p>
            </div>
          </div>
        </div>

        {/* Primary Action: Next Lesson */}
        <div className="space-y-3 pt-2">
          {nextLesson && onNextLesson && (
            <button
              type="button"
              onClick={onNextLesson}
              className="w-full py-4 px-6 rounded-2xl bg-indigo-600 hover:bg-indigo-700 text-white font-black text-sm shadow-xl shadow-indigo-600/30 flex items-center justify-center gap-2 transition-all active:scale-98 cursor-pointer"
            >
              <span>Следующий урок: {nextLesson.title}</span>
              <ArrowRight size={18} />
            </button>
          )}

          {/* Secondary Actions */}
          <div className="flex flex-col sm:flex-row items-center gap-3">
            {onRetake && remainingAttempts > 0 && (
              <button
                type="button"
                onClick={onRetake}
                className="w-full sm:flex-1 py-3 px-4 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 hover:bg-slate-50 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 text-xs font-bold transition-all flex items-center justify-center gap-1.5 cursor-pointer shadow-xs"
              >
                <RotateCcw size={14} />
                <span>Пройти заново ({remainingAttempts} ост.)</span>
              </button>
            )}

            <button
              type="button"
              onClick={onNavigateBack}
              className="w-full sm:flex-1 py-3 px-4 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 hover:bg-slate-50 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 text-xs font-bold transition-all flex items-center justify-center gap-1.5 cursor-pointer shadow-xs"
            >
              <ChevronLeft size={14} />
              <span>К содержанию курса</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
