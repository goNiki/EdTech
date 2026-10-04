'use client';

import React, { useState, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { useRouter } from 'next/navigation';
import {
  ChevronLeft,
  ChevronRight,
  BookOpen,
  CheckCircle2,
  X,
  PlayCircle,
  Award,
  Layers,
  ArrowRight
} from 'lucide-react';

export interface LessonNavContext {
  current_lesson: {
    id: number;
    title: string;
    position?: number;
    section_id?: number;
  };
  course: {
    id: number;
    title: string;
    slug?: string;
  };
  prev_lesson: {
    id: number;
    title: string;
  } | null;
  next_lesson: {
    id: number;
    title: string;
  } | null;
  syllabus: Array<{
    section_id: number;
    section_title: string;
    position?: number;
    lessons: Array<{
      id: number;
      title: string;
      position?: number;
      is_completed: boolean;
      score: number;
    }>;
  }>;
}

interface LessonHeaderNavProps {
  navData: LessonNavContext | null;
  onNavigateToLesson: (lessonId: number) => void;
  onRequestExit?: () => void;
}

export default function LessonHeaderNav({
  navData,
  onNavigateToLesson,
  onRequestExit,
}: LessonHeaderNavProps) {
  const router = useRouter();
  const [isSyllabusOpen, setIsSyllabusOpen] = useState(false);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    if (!isSyllabusOpen) return;
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setIsSyllabusOpen(false);
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isSyllabusOpen]);

  if (!navData) return null;

  const { course, current_lesson, prev_lesson, next_lesson, syllabus } = navData;

  // Calculate course completion stats
  let totalLessons = 0;
  let completedLessons = 0;
  syllabus.forEach((sec) => {
    (sec.lessons || []).forEach((l) => {
      totalLessons += 1;
      if (l.is_completed) completedLessons += 1;
    });
  });
  const progressPercent = totalLessons > 0 ? Math.round((completedLessons / totalLessons) * 100) : 0;

  return (
    <>
      <div className="flex items-center gap-1.5 sm:gap-2">
        {/* Previous Lesson Button */}
        <button
          type="button"
          disabled={!prev_lesson}
          onClick={() => prev_lesson && onNavigateToLesson(prev_lesson.id)}
          className={`p-2 rounded-xl border transition-all flex items-center gap-1.5 text-xs font-bold ${
            prev_lesson
              ? 'border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 shadow-2xs cursor-pointer active:scale-95'
              : 'border-slate-100 dark:border-slate-800/50 bg-slate-50 dark:bg-slate-900/40 text-slate-300 dark:text-slate-700 cursor-not-allowed opacity-50'
          }`}
          title={prev_lesson ? `Предыдущий урок: ${prev_lesson.title}` : 'Это первый урок курса'}
          aria-label="Предыдущий урок"
        >
          <ChevronLeft size={16} />
          <span className="hidden xl:inline truncate max-w-[120px]">
            {prev_lesson ? prev_lesson.title : 'Назад'}
          </span>
        </button>

        {/* Syllabus (Course Content) Sheet Trigger */}
        <button
          type="button"
          onClick={() => setIsSyllabusOpen(true)}
          className="px-2.5 sm:px-3 py-2 rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 hover:border-indigo-300 dark:hover:border-indigo-700 transition-all flex items-center gap-2 text-xs font-bold shadow-2xs cursor-pointer active:scale-95"
          title="Открыть оглавление курса"
          aria-label="Оглавление курса"
        >
          <BookOpen size={15} className="text-indigo-600 dark:text-indigo-400 flex-shrink-0" />
          <span className="hidden md:inline truncate max-w-[140px] text-slate-900 dark:text-white font-extrabold">
            {course.title}
          </span>
          <span className="hidden sm:inline-flex items-center px-1.5 py-0.5 rounded-md bg-indigo-50 dark:bg-indigo-950/70 text-indigo-600 dark:text-indigo-400 text-[10px] font-black">
            {completedLessons}/{totalLessons}
          </span>
        </button>

        {/* Next Lesson Button */}
        <button
          type="button"
          disabled={!next_lesson}
          onClick={() => next_lesson && onNavigateToLesson(next_lesson.id)}
          className={`p-2 rounded-xl border transition-all flex items-center gap-1.5 text-xs font-bold ${
            next_lesson
              ? 'border-indigo-200 dark:border-indigo-800/80 bg-indigo-50/70 dark:bg-indigo-950/50 text-indigo-700 dark:text-indigo-300 hover:bg-indigo-100 dark:hover:bg-indigo-900/60 shadow-2xs cursor-pointer active:scale-95'
              : 'border-slate-100 dark:border-slate-800/50 bg-slate-50 dark:bg-slate-900/40 text-slate-300 dark:text-slate-700 cursor-not-allowed opacity-50'
          }`}
          title={next_lesson ? `Следующий урок: ${next_lesson.title}` : 'Это последний урок курса'}
          aria-label="Следующий урок"
        >
          <span className="hidden xl:inline truncate max-w-[120px]">
            {next_lesson ? next_lesson.title : 'Далее'}
          </span>
          <ChevronRight size={16} />
        </button>
      </div>

      {/* Syllabus Modal Drawer (Right Sheet) */}
      {mounted && isSyllabusOpen && typeof document !== 'undefined' && createPortal(
        <div className="fixed inset-0 z-50 flex justify-end animate-in fade-in duration-200">
          {/* Backdrop */}
          <div
            className="fixed inset-0 bg-black/60 backdrop-blur-xs cursor-pointer"
            onClick={() => setIsSyllabusOpen(false)}
            aria-hidden="true"
          />

          {/* Drawer Content */}
          <div className="relative z-10 w-full max-w-md bg-white dark:bg-slate-900 h-full border-l border-slate-200 dark:border-slate-800 shadow-2xl flex flex-col justify-between animate-in slide-in-from-right duration-300">
            {/* Drawer Header */}
            <div className="p-5 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between gap-3 bg-slate-50/50 dark:bg-slate-950/50">
              <div className="min-w-0">
                <span className="text-[10px] uppercase font-black tracking-wider text-indigo-600 dark:text-indigo-400">
                  Оглавление курса
                </span>
                <h3 className="text-base font-extrabold text-slate-900 dark:text-white truncate">
                  {course.title}
                </h3>
              </div>

              <button
                type="button"
                onClick={() => setIsSyllabusOpen(false)}
                className="p-2 rounded-xl text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 hover:bg-slate-200/60 dark:hover:bg-slate-800 transition-colors cursor-pointer"
                title="Закрыть оглавление"
              >
                <X size={18} />
              </button>
            </div>

            {/* Overall Course Progress Bar */}
            <div className="px-5 py-3.5 bg-indigo-50/60 dark:bg-indigo-950/30 border-b border-indigo-100 dark:border-indigo-900/40">
              <div className="flex items-center justify-between text-xs font-bold text-slate-700 dark:text-slate-200 mb-1.5">
                <span className="flex items-center gap-1.5">
                  <Award size={14} className="text-indigo-600 dark:text-indigo-400" />
                  <span>Общий прогресс курса</span>
                </span>
                <span className="font-mono text-indigo-600 dark:text-indigo-400">
                  {completedLessons} из {totalLessons} ({progressPercent}%)
                </span>
              </div>
              <div className="w-full bg-slate-200 dark:bg-slate-800 h-2 rounded-full overflow-hidden">
                <div
                  className="bg-gradient-to-r from-indigo-500 to-emerald-500 h-full transition-all duration-500 rounded-full"
                  style={{ width: `${progressPercent}%` }}
                />
              </div>
            </div>

            {/* Syllabus Sections & Lessons List */}
            <div className="flex-1 overflow-y-auto p-5 space-y-6">
              {syllabus.map((sec, secIdx) => (
                <div key={sec.section_id || secIdx} className="space-y-2.5">
                  <div className="flex items-center gap-2">
                    <span className="w-5 h-5 rounded-lg bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400 text-[10px] font-black flex items-center justify-center">
                      {secIdx + 1}
                    </span>
                    <h4 className="text-xs font-black uppercase tracking-wider text-slate-500 dark:text-slate-400 truncate">
                      {sec.section_title}
                    </h4>
                  </div>

                  <div className="space-y-1.5 pl-2">
                    {(sec.lessons || []).map((l, lIdx) => {
                      const isCurrent = l.id === current_lesson.id;

                      return (
                        <div
                          key={l.id}
                          onClick={() => {
                            setIsSyllabusOpen(false);
                            if (!isCurrent) {
                              onNavigateToLesson(l.id);
                            }
                          }}
                          className={`p-3 rounded-2xl border transition-all flex items-center justify-between gap-3 cursor-pointer ${
                            isCurrent
                              ? 'bg-indigo-600 text-white border-indigo-600 shadow-md ring-2 ring-indigo-500/20'
                              : 'bg-white dark:bg-slate-900/80 border-slate-200 dark:border-slate-800/80 hover:bg-slate-50 dark:hover:bg-slate-800/60 text-slate-800 dark:text-slate-200'
                          }`}
                        >
                          <div className="flex items-center gap-3 min-w-0">
                            <div className="flex-shrink-0">
                              {l.is_completed ? (
                                <CheckCircle2
                                  size={16}
                                  className={isCurrent ? 'text-white' : 'text-emerald-500'}
                                />
                              ) : (
                                <PlayCircle
                                  size={16}
                                  className={isCurrent ? 'text-white' : 'text-slate-400'}
                                />
                              )}
                            </div>

                            <div className="min-w-0">
                              <p className={`text-xs font-bold truncate leading-tight ${isCurrent ? 'text-white' : ''}`}>
                                {l.title}
                              </p>
                              {isCurrent && (
                                <span className="text-[10px] font-extrabold uppercase tracking-wider text-indigo-200">
                                  Текущий урок
                                </span>
                              )}
                            </div>
                          </div>

                          {/* Completion Badge / Score */}
                          {l.is_completed && (
                            <span
                              className={`text-[11px] font-black font-mono px-2 py-0.5 rounded-lg flex-shrink-0 ${
                                isCurrent
                                  ? 'bg-indigo-700 text-white'
                                  : 'bg-emerald-50 dark:bg-emerald-950/60 text-emerald-600 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800/60'
                              }`}
                            >
                              {l.score || 100}%
                            </span>
                          )}
                        </div>
                      );
                    })}
                  </div>
                </div>
              ))}
            </div>

            {/* Footer Exit to Course */}
            <div className="p-4 border-t border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-950/50">
              <button
                type="button"
                onClick={() => {
                  setIsSyllabusOpen(false);
                  if (onRequestExit) {
                    onRequestExit();
                  } else if (course.slug) {
                    router.push(`/courses/${course.slug}`);
                  } else {
                    router.push('/courses');
                  }
                }}
                className="w-full py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 text-xs font-bold transition-all flex items-center justify-center gap-1.5 shadow-2xs cursor-pointer"
              >
                <span>Вернуться к описанию курса</span>
                <ArrowRight size={14} />
              </button>
            </div>
          </div>
        </div>,
        document.body
      )}
    </>
  );
}
