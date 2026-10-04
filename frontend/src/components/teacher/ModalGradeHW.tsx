'use client';

import React, { useState } from 'react';
import { api } from '@/lib/api';
import { X, CheckCircle, Paperclip, Send, AlertCircle } from 'lucide-react';
import { PendingHWItem } from './PendingHomeworksQueue';

interface ModalGradeHWProps {
  hw: PendingHWItem | null;
  onClose: () => void;
  onGraded: () => void;
}

export default function ModalGradeHW({ hw, onClose, onGraded }: ModalGradeHWProps) {
  const [points, setPoints] = useState<number | string>(hw?.max_points || 10);
  const [feedback, setFeedback] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  if (!hw) return null;

  const maxPoints = hw.max_points || 20;
  const numPoints = Number(points);
  const isInvalidNegative = numPoints < 0;
  const isInvalidOverMax = numPoints > maxPoints;
  const isPointsInvalid = isInvalidNegative || isInvalidOverMax || isNaN(numPoints) || points === '';

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isPointsInvalid) return;

    setIsSubmitting(true);
    try {
      const attId = hw.attempt_id || hw.id;
      const ansId = hw.answer_id || hw.id;

      await api.post(`/quizzes/attempts/${attId}/answers/${ansId}/grade`, {
        points: Number(points),
        feedback: feedback.trim(),
      });

      onGraded();
      onClose();
    } catch (err: any) {
      console.error('Failed to grade homework', err);
      alert(err.response?.data?.message || err.response?.data?.error || 'Ошибка при выставлении оценки');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4">
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 max-w-xl w-full shadow-2xl space-y-5 animate-in fade-in zoom-in duration-200">
        <div className="flex justify-between items-center pb-3 border-b border-slate-100 dark:border-slate-800">
          <div>
            <h3 className="text-base font-extrabold text-slate-900 dark:text-white">
              Проверка домашнего задания
            </h3>
            <p className="text-xs text-slate-500">Студент: {hw.student_name} ({hw.student_email})</p>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-xl text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors"
          >
            <X size={20} />
          </button>
        </div>

        {/* Question and student answer preview */}
        <div className="space-y-3">
          <div className="p-3 bg-slate-50 dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700">
            <span className="text-[10px] font-bold uppercase tracking-wider text-slate-400">Вопрос задания</span>
            <p className="text-xs font-semibold text-slate-800 dark:text-slate-200 mt-0.5">
              {hw.question_text}
            </p>
          </div>

          <div className="p-3 bg-slate-50 dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 space-y-1">
            <span className="text-[10px] font-bold uppercase tracking-wider text-indigo-600 dark:text-indigo-400">
              Ответ студента
            </span>
            <p className="text-xs text-slate-700 dark:text-slate-300 whitespace-pre-wrap leading-relaxed">
              {hw.student_answer}
            </p>

            {hw.attachment_url && (
              <div className="pt-2 border-t border-slate-200 dark:border-slate-700">
                <a
                  href={hw.attachment_url}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1.5 text-xs font-bold text-indigo-600 dark:text-indigo-400 hover:underline"
                >
                  <Paperclip size={14} />
                  <span>Скачать файл решения</span>
                </a>
              </div>
            )}
          </div>
        </div>

        {/* Grading form */}
        <form onSubmit={handleSubmit} className="space-y-4 pt-2 border-t border-slate-100 dark:border-slate-800">
          <div className="space-y-1.5">
            <div className="flex justify-between items-center">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Выставить баллы (из {maxPoints}) *
              </label>
              {isPointsInvalid && (
                <span className="text-[11px] font-bold text-rose-500 flex items-center gap-1">
                  <AlertCircle size={12} />
                  <span>{isInvalidOverMax ? `Максимум ${maxPoints} б.` : isInvalidNegative ? 'Не может быть < 0' : 'Укажите число'}</span>
                </span>
              )}
            </div>
            <input
              type="number"
              min={0}
              max={maxPoints}
              required
              value={points}
              onChange={(e) => setPoints(e.target.value === '' ? '' : Number(e.target.value))}
              className={`w-full p-3 bg-slate-50 dark:bg-slate-800 border rounded-xl text-xs font-bold transition-all focus:outline-none ${
                isPointsInvalid
                  ? 'border-rose-500 text-rose-600 dark:text-rose-400 bg-rose-50/40 dark:bg-rose-950/20 focus:ring-2 focus:ring-rose-500'
                  : 'border-slate-200 dark:border-slate-700 text-indigo-600 dark:text-indigo-400 focus:ring-2 focus:ring-indigo-500'
              }`}
            />
            {isInvalidOverMax && (
              <p className="text-[11px] font-semibold text-rose-500">
                Балл не может превышать {maxPoints}
              </p>
            )}
            {isInvalidNegative && (
              <p className="text-[11px] font-semibold text-rose-500">
                Балл не может быть отрицательным
              </p>
            )}
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
              Рецензия и обратная связь студенту
            </label>
            <textarea
              rows={3}
              value={feedback}
              onChange={(e) => setFeedback(e.target.value)}
              placeholder="Напишите комментарий, похвалите за верное решение или укажите на ошибки..."
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
          </div>

          <div className="flex justify-end gap-3 pt-3">
            <button
              type="button"
              onClick={onClose}
              className="px-5 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800 cursor-pointer"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={isSubmitting || isPointsInvalid}
              className={`px-6 py-2.5 text-white text-xs font-bold rounded-xl shadow-md transition-all flex items-center gap-2 ${
                isPointsInvalid || isSubmitting
                  ? 'bg-slate-400 dark:bg-slate-700 cursor-not-allowed opacity-60'
                  : 'bg-emerald-600 hover:bg-emerald-700 cursor-pointer active:scale-95'
              }`}
            >
              <Send size={14} />
              <span>{isSubmitting ? 'Сохранение...' : 'Утвердить оценку'}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
