'use client';

import React from 'react';
import {
  Award,
  CheckCircle2,
  RotateCcw,
  BookOpen,
  ArrowRight,
  ShieldCheck,
  ChevronLeft,
  Sparkles,
  HelpCircle,
  Clock,
  AlertCircle
} from 'lucide-react';

export interface QuizAttemptItem {
  attempt_number?: number;
  score: number;
  completed_at?: string;
  is_best?: boolean;
}

export interface QuizPreflightScreenProps {
  lessonTitle: string;
  bestScore: number;
  passingScore?: number;
  attemptsMade?: number;
  maxAttempts?: number;
  attemptsHistory?: QuizAttemptItem[];
  onStartAttempt: () => void;
  onReviewMaterials: () => void;
  onNavigateBack: () => void;
}

export default function QuizPreflightScreen({
  lessonTitle,
  bestScore,
  passingScore = 70,
  attemptsMade = 1,
  maxAttempts = 3,
  attemptsHistory = [],
  onStartAttempt,
  onReviewMaterials,
  onNavigateBack,
}: QuizPreflightScreenProps) {
  const isPassed = bestScore >= passingScore;
  const remainingAttempts = Math.max(0, maxAttempts - attemptsMade);
  const canAttempt = maxAttempts === 0 || remainingAttempts > 0;

  return (
    <div className="max-w-2xl mx-auto w-full py-6 space-y-6 animate-in fade-in duration-300">
      {/* Top Card */}
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 sm:p-8 shadow-xl space-y-6 text-center">
        {/* Status Badge & Icon */}
        <div className="w-16 h-16 rounded-3xl bg-indigo-50 dark:bg-indigo-950/60 border border-indigo-200/80 dark:border-indigo-800/60 text-indigo-600 dark:text-indigo-400 flex items-center justify-center mx-auto shadow-inner">
          <Award size={36} />
        </div>

        <div className="space-y-1.5">
          <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-indigo-50 dark:bg-indigo-950/60 text-indigo-700 dark:text-indigo-300 text-[11px] font-black uppercase tracking-wider border border-indigo-200/80 dark:border-indigo-800/60">
            <ShieldCheck size={13} />
            <span>Контрольное тестирование</span>
          </div>
          <h2 className="text-xl sm:text-2xl font-black text-slate-900 dark:text-white">
            {lessonTitle}
          </h2>
          <p className="text-xs text-slate-500 dark:text-slate-400">
            Вы уже проходили этот тест ранее. Ваш лучший результат надёжно сохранён.
          </p>
        </div>

        {/* Highlight Score Box */}
        <div className="grid grid-cols-2 sm:grid-cols-3 gap-3 p-4 bg-slate-50 dark:bg-slate-800/50 rounded-2xl border border-slate-200 dark:border-slate-800 text-left">
          <div className="space-y-1">
            <span className="text-[10px] font-extrabold uppercase text-slate-400 tracking-wider">
              Статус теста
            </span>
            <div className="flex items-center gap-1.5">
              {isPassed ? (
                <>
                  <CheckCircle2 size={16} className="text-emerald-500" />
                  <span className="text-xs font-black text-emerald-600 dark:text-emerald-400">
                    Зачтено
                  </span>
                </>
              ) : (
                <>
                  <AlertCircle size={16} className="text-amber-500" />
                  <span className="text-xs font-black text-amber-600 dark:text-amber-400">
                    Не сдан
                  </span>
                </>
              )}
            </div>
          </div>

          <div className="space-y-1">
            <span className="text-[10px] font-extrabold uppercase text-slate-400 tracking-wider">
              Высший балл
            </span>
            <div className="text-sm font-black text-indigo-600 dark:text-indigo-400">
              {bestScore}%
            </div>
          </div>

          <div className="space-y-1 col-span-2 sm:col-span-1">
            <span className="text-[10px] font-extrabold uppercase text-slate-400 tracking-wider">
              Попытки
            </span>
            <div className="text-xs font-bold text-slate-700 dark:text-slate-200">
              {maxAttempts === 0
                ? `${attemptsMade} (без лимита)`
                : `${attemptsMade} из ${maxAttempts}`}
            </div>
          </div>
        </div>

        {/* History List if available */}
        {attemptsHistory.length > 0 && (
          <div className="text-left space-y-2 pt-2 border-t border-slate-100 dark:border-slate-800">
            <h4 className="text-xs font-black uppercase text-slate-400 tracking-wider">
              История попыток
            </h4>
            <div className="space-y-1.5">
              {attemptsHistory.map((att, idx) => (
                <div
                  key={idx}
                  className="flex items-center justify-between p-2.5 rounded-xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 text-xs font-semibold"
                >
                  <div className="flex items-center gap-2 text-slate-700 dark:text-slate-300">
                    <span>Попытка #{att.attempt_number || idx + 1}</span>
                    {att.is_best && (
                      <span className="px-1.5 py-0.5 rounded text-[10px] font-bold bg-amber-100 dark:bg-amber-950/60 text-amber-700 dark:text-amber-300">
                        ★ Лучшая
                      </span>
                    )}
                  </div>
                  <span className="font-black text-indigo-600 dark:text-indigo-400">
                    {att.score}%
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Action Buttons */}
        <div className="space-y-3 pt-2">
          {canAttempt ? (
            <button
              type="button"
              onClick={onStartAttempt}
              className="w-full py-3.5 px-6 rounded-2xl bg-indigo-600 hover:bg-indigo-700 text-white font-black text-xs sm:text-sm shadow-lg shadow-indigo-600/25 flex items-center justify-center gap-2 transition-all active:scale-98 cursor-pointer"
            >
              <RotateCcw size={16} />
              <span>
                {maxAttempts > 0
                  ? `Начать попытку ${attemptsMade + 1} из ${maxAttempts} →`
                  : 'Пройти тестирование заново →'}
              </span>
            </button>
          ) : (
            <div className="p-3 rounded-2xl bg-amber-50 dark:bg-amber-950/40 border border-amber-200 dark:border-amber-900/40 text-amber-800 dark:text-amber-200 text-xs font-bold">
              Все доступные попытки тестирования исчерпаны.
            </div>
          )}

          <div className="flex items-center justify-center gap-4 text-xs font-bold pt-1">
            <button
              type="button"
              onClick={onReviewMaterials}
              className="text-slate-600 dark:text-slate-400 hover:text-indigo-600 dark:hover:text-indigo-400 flex items-center gap-1.5 cursor-pointer"
            >
              <BookOpen size={14} />
              <span>Посмотреть материалы урока</span>
            </button>
            <span className="text-slate-300 dark:text-slate-700">•</span>
            <button
              type="button"
              onClick={onNavigateBack}
              className="text-slate-600 dark:text-slate-400 hover:text-indigo-600 dark:hover:text-indigo-400 flex items-center gap-1.5 cursor-pointer"
            >
              <ChevronLeft size={14} />
              <span>Назад к программе курса</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
