'use client';

import React, { useState } from 'react';
import {
  X,
  Clock,
  Zap,
  RotateCcw,
  Award,
  EyeOff,
  Eye,
  Sliders,
  CheckCircle2,
  ShieldCheck,
  AlertCircle
} from 'lucide-react';

export interface QuizSettings {
  time_limit_minutes: number;
  question_time_limit_seconds: number;
  max_attempts: number;
  passing_score_percent: number;
  feedback_mode: 'immediate' | 'exam_blind';
  shuffle_questions?: boolean;
}

export const defaultQuizSettings: QuizSettings = {
  time_limit_minutes: 0,
  question_time_limit_seconds: 0,
  max_attempts: 3,
  passing_score_percent: 70,
  feedback_mode: 'immediate',
  shuffle_questions: false,
};

interface ModalQuizSettingsProps {
  isOpen: boolean;
  onClose: () => void;
  settings: QuizSettings;
  onSave: (settings: QuizSettings) => void;
}

export default function ModalQuizSettings({
  isOpen,
  onClose,
  settings,
  onSave,
}: ModalQuizSettingsProps) {
  const [form, setForm] = useState<QuizSettings>({
    ...defaultQuizSettings,
    ...settings,
  });

  React.useEffect(() => {
    if (isOpen) {
      setForm({
        ...defaultQuizSettings,
        ...settings,
      });
    }
  }, [isOpen, settings]);

  if (!isOpen) return null;

  const handleApply = (e: React.FormEvent) => {
    e.preventDefault();
    onSave(form);
    onClose();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/75 backdrop-blur-xs p-4 animate-in fade-in duration-200">
      <div className="bg-slate-900 border border-slate-800 rounded-3xl max-w-lg w-full shadow-2xl flex flex-col max-h-[92vh] overflow-hidden text-slate-100">
        {/* Header */}
        <div className="px-6 py-4 border-b border-slate-800 flex items-center justify-between">
          <div className="flex items-center gap-2.5">
            <div className="w-9 h-9 rounded-xl bg-indigo-500/20 text-indigo-400 flex items-center justify-center border border-indigo-500/30">
              <Sliders size={18} />
            </div>
            <div>
              <h3 className="font-extrabold text-sm sm:text-base text-white">
                Параметры тестирования
              </h3>
              <p className="text-[11px] text-slate-400">
                Таймеры, лимит попыток и правила проверки теста
              </p>
            </div>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
          >
            <X size={18} />
          </button>
        </div>

        {/* Form Body */}
        <form onSubmit={handleApply} className="p-6 space-y-6 overflow-y-auto flex-1 text-xs">
          {/* 1. Time Limit on Entire Quiz */}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <label className="font-bold text-slate-200 flex items-center gap-1.5">
                <Clock size={14} className="text-indigo-400" />
                <span>Общий таймер на тест (минуты)</span>
              </label>
              <span className="text-slate-400 font-mono text-[11px]">
                {form.time_limit_minutes > 0 ? `${form.time_limit_minutes} мин.` : 'Без лимита'}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <input
                type="number"
                min="0"
                max="300"
                value={form.time_limit_minutes}
                onChange={(e) => setForm({ ...form, time_limit_minutes: Math.max(0, parseInt(e.target.value) || 0) })}
                className="w-24 px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-white font-mono focus:border-indigo-500 focus:outline-none"
              />
              <div className="flex flex-wrap gap-1.5">
                {[0, 10, 15, 30, 45].map((mins) => (
                  <button
                    key={mins}
                    type="button"
                    onClick={() => setForm({ ...form, time_limit_minutes: mins })}
                    className={`px-2.5 py-1.5 rounded-lg font-bold text-[11px] transition-colors ${
                      form.time_limit_minutes === mins
                        ? 'bg-indigo-600 text-white'
                        : 'bg-slate-800 text-slate-300 hover:bg-slate-700'
                    }`}
                  >
                    {mins === 0 ? 'Выкл' : `${mins}м`}
                  </button>
                ))}
              </div>
            </div>
            <p className="text-[11px] text-slate-400">
              По истечении таймера ответы будут принудительно отправлены на сервер.
            </p>
          </div>

          {/* 2. Blitz Timer on each question */}
          <div className="space-y-2 pt-2 border-t border-slate-800/80">
            <div className="flex items-center justify-between">
              <label className="font-bold text-slate-200 flex items-center gap-1.5">
                <Zap size={14} className="text-amber-400" />
                <span>Блиц-таймер на вопрос (секунды)</span>
              </label>
              <span className="text-slate-400 font-mono text-[11px]">
                {form.question_time_limit_seconds > 0 ? `${form.question_time_limit_seconds} сек.` : 'Выключен'}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <input
                type="number"
                min="0"
                max="300"
                value={form.question_time_limit_seconds}
                onChange={(e) => setForm({ ...form, question_time_limit_seconds: Math.max(0, parseInt(e.target.value) || 0) })}
                className="w-24 px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-white font-mono focus:border-indigo-500 focus:outline-none"
              />
              <div className="flex flex-wrap gap-1.5">
                {[0, 20, 30, 45, 60].map((secs) => (
                  <button
                    key={secs}
                    type="button"
                    onClick={() => setForm({ ...form, question_time_limit_seconds: secs })}
                    className={`px-2.5 py-1.5 rounded-lg font-bold text-[11px] transition-colors ${
                      form.question_time_limit_seconds === secs
                        ? 'bg-amber-600 text-white'
                        : 'bg-slate-800 text-slate-300 hover:bg-slate-700'
                    }`}
                  >
                    {secs === 0 ? 'Выкл' : `${secs}с`}
                  </button>
                ))}
              </div>
            </div>
            <p className="text-[11px] text-slate-400">
              Полоса времени тает над вопросом. По истечении лимита плеер переходит к следующему вопросу.
            </p>
          </div>

          {/* 3. Max Attempts */}
          <div className="space-y-2 pt-2 border-t border-slate-800/80">
            <div className="flex items-center justify-between">
              <label className="font-bold text-slate-200 flex items-center gap-1.5">
                <RotateCcw size={14} className="text-emerald-400" />
                <span>Лимит попыток сдачи</span>
              </label>
              <span className="text-slate-400 font-mono text-[11px]">
                {form.max_attempts === 0 ? 'Без ограничений' : `${form.max_attempts} попытки`}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <input
                type="number"
                min="0"
                max="10"
                value={form.max_attempts}
                onChange={(e) => setForm({ ...form, max_attempts: Math.max(0, parseInt(e.target.value) || 0) })}
                className="w-24 px-3 py-2 bg-slate-950 border border-slate-800 rounded-xl text-white font-mono focus:border-indigo-500 focus:outline-none"
              />
              <div className="flex flex-wrap gap-1.5">
                {[0, 1, 2, 3, 5].map((att) => (
                  <button
                    key={att}
                    type="button"
                    onClick={() => setForm({ ...form, max_attempts: att })}
                    className={`px-2.5 py-1.5 rounded-lg font-bold text-[11px] transition-colors ${
                      form.max_attempts === att
                        ? 'bg-emerald-600 text-white'
                        : 'bg-slate-800 text-slate-300 hover:bg-slate-700'
                    }`}
                  >
                    {att === 0 ? 'Без лимита' : `${att} ${att === 1 ? 'экзамен' : 'поп.'}`}
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* 4. Passing Score */}
          <div className="space-y-2 pt-2 border-t border-slate-800/80">
            <div className="flex items-center justify-between">
              <label className="font-bold text-slate-200 flex items-center gap-1.5">
                <Award size={14} className="text-purple-400" />
                <span>Проходной порог для зачета</span>
              </label>
              <span className="text-purple-400 font-black text-sm">
                {form.passing_score_percent}%
              </span>
            </div>
            <input
              type="range"
              min="10"
              max="100"
              step="5"
              value={form.passing_score_percent}
              onChange={(e) => setForm({ ...form, passing_score_percent: parseInt(e.target.value) || 70 })}
              className="w-full accent-purple-500 cursor-pointer"
            />
            <div className="flex justify-between text-[10px] text-slate-400">
              <span>10% (мягкий)</span>
              <span>70% (стандарт)</span>
              <span>100% (строгий)</span>
            </div>
          </div>

          {/* 5. Feedback Mode / Blind Exam Mode */}
          <div className="space-y-2 pt-2 border-t border-slate-800/80">
            <label className="font-bold text-slate-200">
              Режим показа правильных ответов
            </label>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
              <button
                type="button"
                onClick={() => setForm({ ...form, feedback_mode: 'immediate' })}
                className={`p-3 rounded-2xl border text-left transition-all ${
                  form.feedback_mode === 'immediate'
                    ? 'border-indigo-500 bg-indigo-950/50 text-white shadow-xs'
                    : 'border-slate-800 bg-slate-950/50 text-slate-400 hover:border-slate-700'
                }`}
              >
                <div className="flex items-center gap-1.5 font-black text-xs text-indigo-400 mb-1">
                  <Eye size={14} />
                  <span>Обучающий</span>
                </div>
                <p className="text-[11px] leading-relaxed">
                  Показывает студенту правильные ответы и пояснения после завершения.
                </p>
              </button>

              <button
                type="button"
                onClick={() => setForm({ ...form, feedback_mode: 'exam_blind' })}
                className={`p-3 rounded-2xl border text-left transition-all ${
                  form.feedback_mode === 'exam_blind'
                    ? 'border-rose-500 bg-rose-950/50 text-white shadow-xs'
                    : 'border-slate-800 bg-slate-950/50 text-slate-400 hover:border-slate-700'
                }`}
              >
                <div className="flex items-center gap-1.5 font-black text-xs text-rose-400 mb-1">
                  <EyeOff size={14} />
                  <span>Экзамен (Blind)</span>
                </div>
                <p className="text-[11px] leading-relaxed">
                  Скрывает правильные ответы от студента, показывая только итоговый балл.
                </p>
              </button>
            </div>
          </div>

          {/* Footer Actions */}
          <div className="pt-4 border-t border-slate-800 flex items-center justify-end gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2.5 rounded-xl border border-slate-800 bg-slate-800 hover:bg-slate-700 text-slate-300 font-bold text-xs transition-colors"
            >
              Отмена
            </button>
            <button
              type="submit"
              className="px-5 py-2.5 rounded-xl bg-indigo-600 hover:bg-indigo-700 text-white font-black text-xs transition-colors shadow-lg shadow-indigo-600/25 flex items-center gap-1.5"
            >
              <CheckCircle2 size={15} />
              <span>Применить параметры</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
