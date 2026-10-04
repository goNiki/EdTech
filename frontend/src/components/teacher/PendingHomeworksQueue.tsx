'use client';

import React from 'react';
import { CheckCircle2, Clock, Paperclip, BookOpen, AlertCircle, User } from 'lucide-react';

export interface PendingHWItem {
  id: number;
  attempt_id: number;
  answer_id: number;
  student_id: number;
  student_name: string;
  student_email: string;
  student_username?: string;
  course_id?: number;
  course_title?: string;
  lesson_id?: number;
  lesson_title: string;
  question_id?: number;
  question_text: string;
  student_answer: string;
  attachment_url?: string;
  max_points?: number;
  rubric?: string;
  submitted_at?: string;
}

interface PendingQueueProps {
  items: PendingHWItem[];
  onOpenGradeModal: (item: PendingHWItem) => void;
  onSelectCourse?: (courseId: number) => void;
}

function formatRelativeTime(dateStr?: string): { text: string; isOverdue: boolean } {
  if (!dateStr) return { text: 'Недавно', isOverdue: false };
  const date = new Date(dateStr);
  if (isNaN(date.getTime())) return { text: dateStr, isOverdue: false };

  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffHours / 24);

  const isOverdue = diffHours >= 48;

  if (diffHours < 1) return { text: 'Менее часа назад', isOverdue: false };
  if (diffHours === 1) return { text: '1 час назад', isOverdue: false };
  if (diffHours < 24) return { text: `${diffHours} ч. назад`, isOverdue: false };
  if (diffDays === 1) return { text: 'Вчера', isOverdue: false };
  if (diffDays === 2) return { text: '2 дня назад', isOverdue: true };
  return { text: `${diffDays} дн. назад`, isOverdue: true };
}

export default function PendingHomeworksQueue({
  items,
  onOpenGradeModal,
  onSelectCourse,
}: PendingQueueProps) {
  if (items.length === 0) {
    return (
      <div className="p-14 text-center bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-3 shadow-xs">
        <div className="w-14 h-14 rounded-2xl bg-emerald-50 dark:bg-emerald-950 text-emerald-600 dark:text-emerald-400 flex items-center justify-center mx-auto shadow-xs">
          <CheckCircle2 size={28} />
        </div>
        <h4 className="text-base font-extrabold text-slate-800 dark:text-slate-200">
          Все домашние задания проверены!
        </h4>
        <p className="text-xs text-slate-500 max-w-sm mx-auto leading-relaxed">
          В очереди нет ожидающих проверки работ. Новые отправленные ответы студентов автоматически появятся здесь.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {items.map((hw) => {
        const timeInfo = formatRelativeTime(hw.submitted_at);

        return (
          <div
            key={hw.id || hw.answer_id}
            className="p-5 sm:p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl flex flex-col md:flex-row justify-between md:items-center gap-5 hover:border-indigo-300 dark:hover:border-indigo-700 transition-all shadow-xs"
          >
            <div className="space-y-3 flex-1 min-w-0">
              {/* Top metadata tags */}
              <div className="flex flex-wrap items-center gap-2">
                {/* Course badge */}
                {hw.course_title && (
                  <button
                    type="button"
                    onClick={() => hw.course_id && onSelectCourse?.(hw.course_id)}
                    className="px-2.5 py-1 rounded-xl text-xs font-bold bg-indigo-50 hover:bg-indigo-100 dark:bg-indigo-950/80 dark:hover:bg-indigo-900 text-indigo-700 dark:text-indigo-300 border border-indigo-200 dark:border-indigo-800 flex items-center gap-1.5 transition-colors cursor-pointer"
                    title="Фильтровать только по этому курсу"
                  >
                    <BookOpen size={13} />
                    <span>{hw.course_title}</span>
                  </button>
                )}

                {/* Lesson title */}
                <span className="px-2.5 py-1 rounded-xl text-xs font-semibold bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 border border-slate-200 dark:border-slate-700">
                  Урок: {hw.lesson_title}
                </span>

                {/* Submission time with overdue indicator */}
                <span
                  className={`px-2.5 py-1 rounded-xl text-xs font-bold flex items-center gap-1.5 ${
                    timeInfo.isOverdue
                      ? 'bg-rose-50 dark:bg-rose-950/80 text-rose-700 dark:text-rose-300 border border-rose-300 dark:border-rose-800'
                      : 'bg-amber-50 dark:bg-amber-950/80 text-amber-700 dark:text-amber-300 border border-amber-200 dark:border-amber-800'
                  }`}
                  title={timeInfo.isOverdue ? 'Работа ожидает проверки более 48 часов' : 'Время сдачи'}
                >
                  {timeInfo.isOverdue ? <AlertCircle size={13} /> : <Clock size={13} />}
                  <span>{timeInfo.text}</span>
                </span>
              </div>

              {/* Student info */}
              <div className="flex items-center gap-2.5">
                <div className="w-7 h-7 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 flex items-center justify-center text-xs font-bold flex-shrink-0">
                  <User size={14} />
                </div>
                <div className="flex items-baseline gap-2">
                  <span className="text-sm font-bold text-slate-900 dark:text-white">
                    {hw.student_name}
                  </span>
                  <span className="text-xs text-slate-400">({hw.student_email})</span>
                </div>
              </div>

              {/* Question and answer preview card */}
              <div className="p-3.5 bg-slate-50 dark:bg-slate-800/60 rounded-2xl border border-slate-200 dark:border-slate-700/60 space-y-1.5">
                <p className="text-xs font-bold text-slate-800 dark:text-slate-200">
                  {hw.question_text}
                </p>
                <p className="text-xs text-slate-600 dark:text-slate-300 line-clamp-3 italic whitespace-pre-wrap leading-relaxed">
                  «{hw.student_answer}»
                </p>
              </div>

              {/* Attachment if present */}
              {hw.attachment_url && (
                <div>
                  <a
                    href={hw.attachment_url}
                    target="_blank"
                    rel="noreferrer"
                    className="inline-flex items-center gap-1.5 text-xs font-bold text-indigo-600 dark:text-indigo-400 hover:underline bg-indigo-50/50 dark:bg-indigo-950/30 px-3 py-1.5 rounded-xl border border-indigo-200/60 dark:border-indigo-800/40"
                  >
                    <Paperclip size={14} />
                    <span>Прикрепленный файл решения</span>
                  </a>
                </div>
              )}
            </div>

            {/* Action button */}
            <div className="flex sm:flex-col items-center justify-end gap-2 flex-shrink-0 pt-2 sm:pt-0">
              <button
                type="button"
                onClick={() => onOpenGradeModal(hw)}
                className="w-full sm:w-auto px-6 py-3 bg-emerald-600 hover:bg-emerald-700 active:scale-95 text-white font-bold text-xs rounded-xl shadow-md shadow-emerald-600/20 transition-all flex items-center justify-center gap-2 cursor-pointer"
              >
                <span>Проверить работу</span>
              </button>
            </div>
          </div>
        );
      })}
    </div>
  );
}
