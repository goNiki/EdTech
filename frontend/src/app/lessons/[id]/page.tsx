'use client';

import React, { useEffect, useState, use } from 'react';
import { api } from '@/lib/api';
import ProtectedRoute from '@/components/ProtectedRoute';
import PuckLessonViewer, { LessonCompletionPayload } from '@/components/player/PuckLessonViewer';
import { useRouter } from 'next/navigation';
import { ChevronLeft, CheckCircle, Sparkles, ArrowRight, RotateCcw, Loader2 } from 'lucide-react';

export default function LessonPlayer({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const [lessonData, setLessonData] = useState<any>(null);
  const [progressData, setProgressData] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isCompleted, setIsCompleted] = useState(false);
  const [toastMsg, setToastMsg] = useState<string | null>(null);

  const showToast = (msg: string) => {
    setToastMsg(msg);
    setTimeout(() => setToastMsg(null), 3000);
  };

  useEffect(() => {
    const fetchLessonAndProgress = async () => {
      try {
        // 1. Notify backend about starting the lesson
        await api.post(`/lessons/${id}/start`).catch(() => {});

        // 2. Fetch lesson content & progress in parallel
        const [lessonRes, progressRes] = await Promise.all([
          api.get(`/lessons/${id}`),
          api.get(`/lessons/${id}/progress`).catch(() => null),
        ]);

        const lData = lessonRes.data?.data?.lesson || lessonRes.data?.lesson || lessonRes.data || {};
        setLessonData(lData);

        const pData = progressRes?.data?.data || progressRes?.data;
        if (pData && (pData.status === 'completed' || pData.Status === 'completed')) {
          setIsCompleted(true);
          setProgressData(pData);
        }
      } catch (err) {
        console.error('Failed to load lesson', err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchLessonAndProgress();
  }, [id]);

  const [isCompleting, setIsCompleting] = useState(false);

  const handleRetake = () => {
    setIsCompleted(false);
    showToast('Режим тренировки: вы можете заново решить задания урока.');
  };

  const handleComplete = async (payload?: LessonCompletionPayload) => {
    setIsCompleting(true);
    try {
      const res = await api.post(`/lessons/${id}/complete`, {
        score: payload?.score ?? 100,
        answers: payload?.answers ?? [],
        essays: payload?.essays ?? [],
      });
      const data = res.data?.data || res.data;
      const verifiedScore = data?.score ?? data?.Score ?? payload?.score ?? 100;

      setIsCompleted(true);
      setProgressData((prev: any) => ({
        ...prev,
        status: 'completed',
        score: verifiedScore,
        completed_at: new Date().toISOString(),
        ...(data || {}),
      }));
      showToast('Урок успешно завершен и проверен на сервере!');
      return data;
    } catch (error: any) {
      console.error('Failed to complete lesson on server', error);
      setIsCompleted(true);
      const msg = error.response?.data?.message || 'Результат зафиксирован локально.';
      showToast(msg);
      return null;
    } finally {
      setIsCompleting(false);
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
            {isCompleted ? (
              <div className="flex items-center gap-2">
                <div className="flex items-center gap-1.5 px-3.5 py-1.5 bg-emerald-50 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-300 rounded-xl border border-emerald-200/80 dark:border-emerald-800/60 text-xs font-bold shadow-2xs">
                  <CheckCircle size={14} className="text-emerald-500" />
                  <span>Пройден • {progressData?.score ?? 100}%</span>
                </div>
                <button
                  type="button"
                  onClick={handleRetake}
                  className="px-3.5 py-1.5 bg-white dark:bg-slate-800 hover:bg-slate-50 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold transition-all flex items-center gap-1 cursor-pointer shadow-2xs"
                  title="Пройти материал заново"
                >
                  <RotateCcw size={13} />
                  <span className="hidden sm:inline">Пройти заново</span>
                </button>
              </div>
            ) : (
              <button
                disabled={isCompleting}
                onClick={() => handleComplete()}
                className="px-5 py-2 rounded-xl text-xs font-bold transition-all shadow-md flex items-center gap-1.5 cursor-pointer bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white"
              >
                {isCompleting ? (
                  <>
                    <Loader2 size={16} className="animate-spin" />
                    <span>Проверка...</span>
                  </>
                ) : (
                  <>
                    <CheckCircle size={16} />
                    <span>Завершить урок</span>
                  </>
                )}
              </button>
            )}
          </div>
        </header>

        {/* Lesson Body */}
        <main className="flex-1 max-w-4xl w-full mx-auto py-8 px-6 space-y-6">
          {/* Result Banner when already completed */}
          {isCompleted && progressData && (
            <div className="bg-gradient-to-r from-emerald-500/10 via-teal-500/10 to-indigo-500/10 dark:from-emerald-950/40 dark:via-teal-950/30 dark:to-indigo-950/40 border border-emerald-200 dark:border-emerald-800/80 rounded-3xl p-6 sm:p-8 shadow-xs space-y-4">
              <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                <div className="flex items-center gap-4">
                  <div className="w-14 h-14 rounded-2xl bg-emerald-500 text-white flex items-center justify-center font-black text-2xl shadow-md flex-shrink-0">
                    <CheckCircle size={28} />
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="px-2.5 py-0.5 rounded-lg text-[10px] font-black uppercase tracking-wider bg-emerald-100 dark:bg-emerald-900/60 text-emerald-800 dark:text-emerald-300">
                        Урок пройден
                      </span>
                      {progressData.completed_at && (
                        <span className="text-xs text-slate-400">
                          {new Date(progressData.completed_at).toLocaleDateString('ru-RU', {
                            day: 'numeric',
                            month: 'long',
                            year: 'numeric',
                          })}
                        </span>
                      )}
                    </div>
                    <h2 className="text-xl font-extrabold text-slate-900 dark:text-white mt-1">
                      Вы успешно завершили этот урок
                    </h2>
                  </div>
                </div>

                <div className="flex items-baseline gap-1.5 bg-white dark:bg-slate-900 px-5 py-3 rounded-2xl border border-emerald-100 dark:border-emerald-900/40 shadow-xs">
                  <span className="text-2xl font-black text-emerald-600 dark:text-emerald-400">
                    {progressData.score ?? 100}
                  </span>
                  <span className="text-xs font-bold text-slate-400">/ 100 баллов</span>
                </div>
              </div>

              {/* Feedback from teacher if available */}
              {progressData.feedback && (
                <div className="p-4 rounded-2xl bg-white/80 dark:bg-slate-900/80 border border-emerald-100 dark:border-emerald-900/40 text-xs space-y-1">
                  <span className="text-[10px] uppercase font-bold text-emerald-600 dark:text-emerald-400 tracking-wider">
                    Комментарий преподавателя
                  </span>
                  <p className="text-slate-700 dark:text-slate-200 italic leading-relaxed">
                    «{progressData.feedback}»
                  </p>
                </div>
              )}

              <div className="flex flex-wrap items-center justify-between gap-3 pt-2 border-t border-emerald-100/60 dark:border-emerald-900/40 text-xs">
                <span className="text-slate-500 dark:text-slate-400">
                  Вы можете свободно просматривать материалы или заново решить тесты для тренировки.
                </span>
                <button
                  type="button"
                  onClick={handleRetake}
                  className="px-4 py-2 bg-white dark:bg-slate-800 hover:bg-slate-50 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 font-bold rounded-xl border border-slate-200 dark:border-slate-700 transition-all flex items-center gap-1.5 shadow-2xs cursor-pointer"
                >
                  <RotateCcw size={13} />
                  <span>Пройти заново</span>
                </button>
              </div>
            </div>
          )}

          <PuckLessonViewer
            contentJson={lessonData?.content || lessonData?.Content || '{}'}
            onComplete={handleComplete}
            initialProgress={progressData}
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
