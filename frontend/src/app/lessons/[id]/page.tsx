'use client';

import React, { useEffect, useState, use } from 'react';
import { api } from '@/lib/api';
import ProtectedRoute from '@/components/ProtectedRoute';
import PuckLessonViewer, { LessonCompletionPayload } from '@/components/player/PuckLessonViewer';
import { useRouter } from 'next/navigation';
import { ChevronLeft, CheckCircle, Sparkles, ArrowRight } from 'lucide-react';

export default function LessonPlayer({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const [lessonData, setLessonData] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isCompleted, setIsCompleted] = useState(false);
  const [toastMsg, setToastMsg] = useState<string | null>(null);

  const showToast = (msg: string) => {
    setToastMsg(msg);
    setTimeout(() => setToastMsg(null), 3000);
  };

  useEffect(() => {
    const fetchAndStartLesson = async () => {
      try {
        // 1. Notify backend about starting the lesson
        await api.post(`/lessons/${id}/start`).catch(() => {});

        // 2. Fetch lesson content
        const res = await api.get(`/lessons/${id}`);
        const data = res.data?.data?.lesson || res.data?.lesson || res.data || {};
        setLessonData(data);
      } catch (err) {
        console.error('Failed to load lesson', err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchAndStartLesson();
  }, [id]);

  const handleComplete = async (payload?: LessonCompletionPayload) => {
    try {
      await api.post(`/lessons/${id}/complete`, {
        score: payload?.score ?? 100,
        essays: payload?.essays ?? [],
      });
      setIsCompleted(true);
      showToast('Урок успешно завершен! Прогресс и задания сохранены.');
    } catch (error) {
      console.error('Failed to complete lesson', error);
      setIsCompleted(true);
    }
  };

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center bg-slate-50 dark:bg-slate-950 text-slate-500">
        <div className="animate-pulse text-sm font-bold flex items-center gap-3">
          <Sparkles className="text-indigo-600 animate-spin" />
          <span>Подготовка материалов урока...</span>
        </div>
      </div>
    );
  }

  const courseId = lessonData?.course_id || lessonData?.CourseID;

  return (
    <ProtectedRoute allowedRoles={['student', 'teacher', 'author', 'admin']}>
      <div className="min-h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 flex flex-col">
        {/* Sticky Player Header */}
        <header className="sticky top-0 z-40 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200 dark:border-slate-800 px-6 py-3.5 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <button
              onClick={() => {
                if (courseId) router.push(`/dashboard/courses/${courseId}`);
                else router.back();
              }}
              className="p-2 rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors shadow-xs cursor-pointer"
              title="Назад к курсу"
            >
              <ChevronLeft size={18} />
            </button>
            <div>
              <h1 className="font-extrabold text-sm md:text-base text-slate-900 dark:text-white truncate max-w-lg">
                {lessonData?.title || lessonData?.Title || 'Урок'}
              </h1>
              <p className="text-[10px] text-slate-400 uppercase font-bold tracking-wider">
                {lessonData?.type || lessonData?.Type || 'Интерактивный урок'}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <button
              onClick={() => handleComplete()}
              disabled={isCompleted}
              className={`px-5 py-2 rounded-xl text-xs font-bold transition-all shadow-md flex items-center gap-1.5 cursor-pointer ${
                isCompleted
                  ? 'bg-emerald-600 text-white cursor-default'
                  : 'bg-indigo-600 hover:bg-indigo-700 text-white'
              }`}
            >
              <CheckCircle size={16} />
              <span>{isCompleted ? 'Урок завершен' : 'Завершить урок'}</span>
            </button>
          </div>
        </header>

        {/* Lesson Body */}
        <main className="flex-1 max-w-4xl w-full mx-auto py-10 px-6">
          <PuckLessonViewer
            contentJson={lessonData?.content || lessonData?.Content || '{}'}
            onComplete={handleComplete}
            onNavigateBack={() => {
              if (courseId) router.push(`/dashboard/courses/${courseId}`);
              else router.back();
            }}
          />
        </main>

        {/* Footer Navigation */}
        <footer className="border-t border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 py-4 px-6 text-center">
          <button
            onClick={() => {
              if (courseId) router.push(`/dashboard/courses/${courseId}`);
              else router.back();
            }}
            className="text-xs font-bold text-indigo-600 dark:text-indigo-400 hover:underline inline-flex items-center gap-1 cursor-pointer"
          >
            <span>Вернуться к содержанию курса</span>
            <ArrowRight size={14} />
          </button>
        </footer>

        {/* Toast */}
        {toastMsg && (
          <div className="fixed bottom-6 right-6 z-50 bg-emerald-600 text-white text-xs font-bold px-4 py-3 rounded-2xl shadow-xl flex items-center gap-2 animate-bounce">
            <CheckCircle size={16} />
            <span>{toastMsg}</span>
          </div>
        )}
      </div>
    </ProtectedRoute>
  );
}
