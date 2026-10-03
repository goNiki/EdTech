'use client';

import React from 'react';
import { CheckCircle2, Clock, FileText, Paperclip } from 'lucide-react';

export interface PendingHWItem {
  id: number;
  attempt_id: number;
  answer_id: number;
  student_id: number;
  student_name: string;
  student_email: string;
  lesson_title: string;
  question_text: string;
  student_answer: string;
  attachment_url?: string;
  max_points?: number;
  submitted_at?: string;
}

interface PendingQueueProps {
  items: PendingHWItem[];
  onOpenGradeModal: (item: PendingHWItem) => void;
}

export default function PendingHomeworksQueue({
  items,
  onOpenGradeModal,
}: PendingQueueProps) {
  if (items.length === 0) {
    return (
      <div className="p-12 text-center bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-3">
        <div className="w-12 h-12 rounded-2xl bg-emerald-50 dark:bg-emerald-950 text-emerald-600 dark:text-emerald-400 flex items-center justify-center mx-auto shadow-xs">
          <CheckCircle2 size={24} />
        </div>
        <h4 className="text-sm font-bold text-slate-800 dark:text-slate-200">
          Все домашние задания проверены!
        </h4>
        <p className="text-xs text-slate-500 max-w-sm mx-auto">
          Новые отправленные работы студентов сразу отобразятся в этой очереди.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {items.map((hw) => (
        <div
          key={hw.id || hw.answer_id}
          className="p-5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl flex flex-col md:flex-row justify-between md:items-center gap-4 hover:border-indigo-300 dark:hover:border-indigo-700 transition-all shadow-xs"
        >
          <div className="space-y-2 flex-1 min-w-0">
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-sm font-bold text-slate-900 dark:text-white">
                {hw.student_name}
              </span>
              <span className="text-xs text-slate-400">({hw.student_email})</span>
              <span className="px-2 py-0.5 rounded-lg text-[10px] font-bold bg-amber-50 dark:bg-amber-950 text-amber-700 dark:text-amber-300 flex items-center gap-1">
                <Clock size={12} /> {hw.submitted_at || 'Недавно'}
              </span>
            </div>

            <p className="text-xs font-semibold text-indigo-600 dark:text-indigo-400">
              Урок: {hw.lesson_title}
            </p>

            <div className="p-3 bg-slate-50 dark:bg-slate-800/60 rounded-2xl border border-slate-200 dark:border-slate-700/60 space-y-1">
              <p className="text-[11px] font-bold text-slate-700 dark:text-slate-300">
                {hw.question_text}
              </p>
              <p className="text-xs text-slate-600 dark:text-slate-300 line-clamp-2 italic">
                «{hw.student_answer}»
              </p>
            </div>

            {hw.attachment_url && (
              <a
                href={hw.attachment_url}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center gap-1 text-[11px] font-bold text-indigo-600 dark:text-indigo-400 hover:underline"
              >
                <Paperclip size={14} />
                <span>Прикрепленный файл решения</span>
              </a>
            )}
          </div>

          <button
            onClick={() => onOpenGradeModal(hw)}
            className="px-5 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs rounded-xl shadow-md transition-all flex-shrink-0"
          >
            Проверить работу
          </button>
        </div>
      ))}
    </div>
  );
}
