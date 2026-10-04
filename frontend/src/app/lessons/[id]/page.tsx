'use client';

import React, { useEffect, useState, use } from 'react';
import { api } from '@/lib/api';
import ProtectedRoute from '@/components/ProtectedRoute';
import PuckLessonViewer, { LessonCompletionPayload } from '@/components/player/PuckLessonViewer';
import QuizStepperPlayer from '@/components/player/QuizStepperPlayer';
import QuizPreflightScreen from '@/components/player/QuizPreflightScreen';
import QuizResultScreen from '@/components/player/QuizResultScreen';
import LessonHeaderNav, { LessonNavContext } from '@/components/player/LessonHeaderNav';
import HomeworkFeedbackCard, { HomeworkFeedbackData } from '@/components/player/HomeworkFeedbackCard';
import { useRouter } from 'next/navigation';
import { ChevronLeft, CheckCircle, CheckCircle2, Clock, Sparkles, ArrowRight, RotateCcw, Loader2, AlertTriangle, Zap, LayoutList, Layers } from 'lucide-react';

export default function LessonPlayer({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const [lessonData, setLessonData] = useState<any>(null);
  const [progressData, setProgressData] = useState<any>(null);
  const [navData, setNavData] = useState<LessonNavContext | null>(null);
  const [homeworkFeedback, setHomeworkFeedback] = useState<HomeworkFeedbackData | null>(null);
  const [pendingLessonId, setPendingLessonId] = useState<number | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isCompleted, setIsCompleted] = useState(false);
  const [isAttemptStarted, setIsAttemptStarted] = useState(false);
  const [activeAttempt, setActiveAttempt] = useState<any>(null);
  const [activeAttemptId, setActiveAttemptId] = useState<number | null>(null);
  const [showResultScreen, setShowResultScreen] = useState(false);
  const [lastResult, setLastResult] = useState<{
    score: number;
    earnedPoints: number;
    totalPoints: number;
  } | null>(null);
  const [nextLesson, setNextLesson] = useState<{ id: number; title: string } | null>(null);
  const [toastMsg, setToastMsg] = useState<string | null>(null);
  const [isExitModalOpen, setIsExitModalOpen] = useState(false);
  const [playerMode, setPlayerMode] = useState<'stepper' | 'scroll'>('stepper');

  const showToast = (msg: string) => {
    setToastMsg(msg);
    setTimeout(() => setToastMsg(null), 3000);
  };

  useEffect(() => {
    const fetchLessonAndProgress = async () => {
      try {
        // Fetch lesson content, progress, attempts summary, active attempt, navigation, and homework feedback in parallel
        const [lessonRes, progressRes, attemptsRes, activeAttRes, navRes, hwRes] = await Promise.all([
          api.get(`/lessons/${id}`),
          api.get(`/lessons/${id}/progress`).catch(() => null),
          api.get(`/lessons/${id}/attempts/summary`).catch(() => null),
          api.get(`/lessons/${id}/attempts/active`).catch(() => null),
          api.get(`/lessons/${id}/navigation`).catch(() => null),
          api.get(`/lessons/${id}/homework-feedback`).catch(() => null),
        ]);

        const lData = lessonRes.data?.data?.lesson || lessonRes.data?.lesson || lessonRes.data || {};
        setLessonData(lData);

        const pData = progressRes?.data?.data || progressRes?.data;
        const attData = attemptsRes?.data?.data || attemptsRes?.data;
        const actData = activeAttRes?.data?.data || activeAttRes?.data;
        const nData = navRes?.data?.data || navRes?.data;
        const hwData = hwRes?.data?.data || hwRes?.data;

        if (hwData?.has_submission) {
          setHomeworkFeedback(hwData);
        }

        if (nData) {
          setNavData(nData);
          if (nData.next_lesson) {
            setNextLesson(nData.next_lesson);
          }
        }

        if (actData?.has_active_attempt && actData?.attempt) {
          setActiveAttempt(actData);
          setActiveAttemptId(actData.attempt.id);
          setIsAttemptStarted(true);
        }

        const hasBeenCompleted = Boolean(
          (pData && (pData.status === 'completed' || pData.Status === 'completed' || pData.completed_at || pData.CompletedAt)) ||
          (pData?.score !== null && pData?.score !== undefined && pData?.score > 0) ||
          (attData && (attData.is_passed || attData.total_attempts_made > 0))
        );

        if (hasBeenCompleted) {
          setIsCompleted(true);
          setProgressData(pData || attData);
        } else {
          // Only start if not already completed
          await api.post(`/lessons/${id}/start`).catch(() => {});
        }
      } catch (err) {
        console.error('Failed to load lesson', err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchLessonAndProgress();
  }, [id]);

  const handleStartAttempt = async () => {
    try {
      const res = await api.post(`/lessons/${id}/attempts/start`).catch(() => api.post(`/lessons/${id}/start`));
      const attId = res?.data?.data?.attempt_id || res?.data?.attempt_id;
      if (attId) {
        setActiveAttemptId(attId);
      }
    } catch {}
    setShowResultScreen(false);
    setIsAttemptStarted(true);
    setIsCompleted(false);
    showToast('Попытка начата! Удачи в тестировании.');
  };

  const [isCompleting, setIsCompleting] = useState(false);

  // Determine if lesson contains interactive quizzes/assignments
  const hasQuizzes = React.useMemo(() => {
    try {
      const rawContent = lessonData?.content || lessonData?.Content;
      if (!rawContent) return false;
      const parsed = typeof rawContent === 'string' ? JSON.parse(rawContent) : rawContent;
      const blocks = parsed?.content || [];
      return blocks.some((b: any) => 
        b.type?.startsWith('Quiz') || b.type === 'FileUploadBlock'
      );
    } catch (e) {
      return false;
    }
  }, [lessonData]);

  const quizPointsInfo = React.useMemo(() => {
    try {
      const rawContent = lessonData?.content || lessonData?.Content;
      if (!rawContent) return { totalPoints: 100 };
      const parsed = typeof rawContent === 'string' ? JSON.parse(rawContent) : rawContent;
      const blocks = parsed?.content || [];
      let sum = 0;
      blocks.forEach((b: any) => {
        if (b.type?.startsWith('Quiz') || b.type === 'FileUploadBlock') {
          sum += Number(b.props?.points) || 10;
        }
      });
      return { totalPoints: sum || 100 };
    } catch {
      return { totalPoints: 100 };
    }
  }, [lessonData]);

  const effectiveQuizSettings = React.useMemo(() => {
    try {
      const rawContent = lessonData?.content || lessonData?.Content;
      if (rawContent) {
        const parsed = typeof rawContent === 'string' ? JSON.parse(rawContent) : rawContent;
        if (parsed?.quiz_settings) return parsed.quiz_settings;
      }
      return lessonData?.quiz_settings || lessonData?.QuizSettings;
    } catch {
      return lessonData?.quiz_settings || lessonData?.QuizSettings;
    }
  }, [lessonData]);

  const courseId = lessonData?.course_id || lessonData?.CourseID;

  // Fetch course structure to determine next lesson
  useEffect(() => {
    if (!courseId) return;
    api.get(`/courses/${courseId}/structure`)
      .then((res) => {
        const data = res.data?.data || res.data;
        const secs = data?.Sections || data?.sections || [];
        const allLessons: any[] = [];
        secs.forEach((s: any) => {
          const ls = s.Lessons || s.lessons || [];
          ls.forEach((l: any) => allLessons.push(l));
        });
        const currentIdx = allLessons.findIndex((l: any) => (l.ID || l.id) === Number(id));
        if (currentIdx !== -1 && currentIdx + 1 < allLessons.length) {
          const next = allLessons[currentIdx + 1];
          setNextLesson({
            id: next.ID || next.id,
            title: next.Title || next.title || `Урок ${currentIdx + 2}`,
          });
        }
      })
      .catch((err) => {
        console.warn('Could not fetch course structure for next lesson', err);
      });
  }, [courseId, id]);

  const performExit = () => {
    if (courseId) router.push(`/dashboard/courses/${courseId}`);
    else router.back();
  };

  const handleNavigateToLesson = (targetLessonId: number) => {
    if (hasQuizzes && !isCompleted && !showResultScreen) {
      setPendingLessonId(targetLessonId);
      setIsExitModalOpen(true);
    } else {
      router.push(`/lessons/${targetLessonId}`);
    }
  };

  const handleRequestExit = () => {
    if (hasQuizzes && !isCompleted && !showResultScreen) {
      setIsExitModalOpen(true);
    } else {
      performExit();
    }
  };

  const handleRetake = () => {
    setShowResultScreen(false);
    setIsAttemptStarted(true);
    setIsCompleted(false);
    showToast('Режим тестирования: вы можете заново решить задания урока.');
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

      if (hasQuizzes) {
        const totalPts = quizPointsInfo.totalPoints;
        const earnedPts = Math.round((verifiedScore / 100) * totalPts);
        setLastResult({
          score: verifiedScore,
          earnedPoints: earnedPts,
          totalPoints: totalPts,
        });
        setShowResultScreen(true);
      }

      showToast('Урок успешно завершен и проверен на сервере!');
      return data;
    } catch (error: any) {
      console.error('Failed to complete lesson on server', error);
      setIsCompleted(true);
      const verifiedScore = payload?.score ?? 100;
      if (hasQuizzes) {
        const totalPts = quizPointsInfo.totalPoints;
        const earnedPts = Math.round((verifiedScore / 100) * totalPts);
        setLastResult({
          score: verifiedScore,
          earnedPoints: earnedPts,
          totalPoints: totalPts,
        });
        setShowResultScreen(true);
      }
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

  return (
    <ProtectedRoute allowedRoles={['student', 'teacher', 'author', 'admin']}>
      <div className="min-h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 flex flex-col overflow-x-hidden">
        {/* Sticky Player Header */}
        <header className="sticky top-0 z-40 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200 dark:border-slate-800 px-4 sm:px-6 py-3 flex items-center justify-between gap-3">
          <div className="flex items-center gap-2 sm:gap-4 min-w-0">
            <button
              onClick={handleRequestExit}
              className="p-2 rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors shadow-xs cursor-pointer flex-shrink-0"
              title="Назад к курсу"
            >
              <ChevronLeft size={18} />
            </button>
            <div className="min-w-0">
              <h1 className="font-extrabold text-sm md:text-base text-slate-900 dark:text-white truncate max-w-[130px] sm:max-w-xs md:max-w-sm">
                {lessonData?.title || lessonData?.Title || 'Урок'}
              </h1>
              <p className="text-[10px] text-slate-400 uppercase font-bold tracking-wider truncate">
                {lessonData?.type || lessonData?.Type || 'Интерактивный урок'}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2 sm:gap-3 flex-shrink-0">
            {/* Homework submission status badge */}
            {homeworkFeedback?.has_submission && (
              <span
                className={`hidden md:inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl border text-xs font-bold shadow-2xs ${
                  homeworkFeedback.status === 'pending'
                    ? 'bg-indigo-50 dark:bg-indigo-950/60 text-indigo-700 dark:text-indigo-300 border-indigo-200/80 dark:border-indigo-800/60'
                    : 'bg-emerald-50 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-300 border-emerald-200/80 dark:border-emerald-800/60'
                }`}
                title={
                  homeworkFeedback.status === 'pending'
                    ? 'Письменная работа находится на проверке у преподавателя'
                    : 'Преподаватель проверил вашу работу'
                }
              >
                {homeworkFeedback.status === 'pending' ? (
                  <>
                    <Clock size={13} className="text-indigo-600 dark:text-indigo-400" />
                    <span>ДЗ на проверке</span>
                  </>
                ) : (
                  <>
                    <CheckCircle2 size={13} className="text-emerald-600 dark:text-emerald-400" />
                    <span>ДЗ проверено</span>
                  </>
                )}
              </span>
            )}

            {/* In-Player Navigation and Course Syllabus */}
            <LessonHeaderNav
              navData={navData}
              onNavigateToLesson={handleNavigateToLesson}
              onRequestExit={handleRequestExit}
            />

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
            ) : hasQuizzes ? (
              <div className="flex items-center gap-2">
                {/* Stepper vs Scroll Toggle */}
                <div className="flex items-center bg-slate-100 dark:bg-slate-800 p-0.5 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold">
                  <button
                    type="button"
                    onClick={() => setPlayerMode('stepper')}
                    className={`px-3 py-1.5 rounded-lg transition-all cursor-pointer flex items-center gap-1.5 ${
                      playerMode === 'stepper'
                        ? 'bg-white dark:bg-slate-900 text-indigo-600 dark:text-indigo-400 shadow-xs'
                        : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200'
                    }`}
                    title="Пошаговый режим: 1 вопрос на экран"
                  >
                    <Layers size={13} />
                    <span>По шагам</span>
                  </button>
                  <button
                    type="button"
                    onClick={() => setPlayerMode('scroll')}
                    className={`px-3 py-1.5 rounded-lg transition-all cursor-pointer flex items-center gap-1.5 ${
                      playerMode === 'scroll'
                        ? 'bg-white dark:bg-slate-900 text-indigo-600 dark:text-indigo-400 shadow-xs'
                        : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200'
                    }`}
                    title="Вся лента: все задания сразу"
                  >
                    <LayoutList size={13} />
                    <span>Вся лента</span>
                  </button>
                </div>
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
        <main className="flex-1 max-w-4xl w-full mx-auto py-6 sm:py-8 px-4 sm:px-6 space-y-6">
          {/* Homework Teacher Review & Feedback Card */}
          {homeworkFeedback?.has_submission && (!isAttemptStarted || isCompleted) && (
            <HomeworkFeedbackCard data={homeworkFeedback} />
          )}

          {hasQuizzes && showResultScreen && lastResult ? (
            <QuizResultScreen
              score={lastResult.score}
              earnedPoints={lastResult.earnedPoints}
              totalMaxPoints={lastResult.totalPoints}
              passingScore={effectiveQuizSettings?.passing_score_percent || lessonData?.passing_score || lessonData?.PassingScore || 70}
              bestScore={progressData?.score || lastResult.score}
              remainingAttempts={2}
              nextLesson={nextLesson}
              onNextLesson={nextLesson ? () => router.push(`/lessons/${nextLesson.id}`) : undefined}
              onRetake={handleRetake}
              onNavigateBack={performExit}
            />
          ) : hasQuizzes && isCompleted && !isAttemptStarted ? (
            <QuizPreflightScreen
              lessonTitle={lessonData?.title || lessonData?.Title || 'Контрольное тестирование'}
              bestScore={progressData?.score ?? 100}
              passingScore={effectiveQuizSettings?.passing_score_percent || lessonData?.passing_score || lessonData?.PassingScore || 70}
              attemptsMade={progressData?.attempts_count || 1}
              maxAttempts={effectiveQuizSettings?.max_attempts ?? (lessonData?.max_attempts || 3)}
              attemptsHistory={progressData?.attempts || [
                {
                  attempt_number: 1,
                  score: progressData?.score ?? 100,
                  completed_at: progressData?.completed_at,
                  is_best: true,
                }
              ]}
              onStartAttempt={handleStartAttempt}
              onReviewMaterials={() => {
                setIsAttemptStarted(true);
                setPlayerMode('scroll');
              }}
              onNavigateBack={performExit}
            />
          ) : (
            <>
              {/* Result Banner when already completed (for non-quiz lessons) */}
              {isCompleted && !hasQuizzes && progressData && (
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
                      Вы можете свободно просматривать материалы или повторить урок.
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

              {hasQuizzes && playerMode === 'stepper' ? (
                <QuizStepperPlayer
                  lessonId={id}
                  attemptId={activeAttemptId || activeAttempt?.attempt?.id}
                  initialActiveAttempt={activeAttempt}
                  contentJson={lessonData?.content || lessonData?.Content || '{}'}
                  onComplete={handleComplete}
                  initialProgress={progressData}
                  onNavigateBack={handleRequestExit}
                  lessonTitle={lessonData?.title || lessonData?.Title}
                  quizSettings={effectiveQuizSettings}
                />
              ) : (
                <PuckLessonViewer
                  contentJson={lessonData?.content || lessonData?.Content || '{}'}
                  onComplete={handleComplete}
                  initialProgress={progressData}
                  onNavigateBack={handleRequestExit}
                />
              )}
            </>
          )}
        </main>

        {/* Footer Navigation */}
        <footer className="border-t border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 py-4 px-6 text-center">
          <button
            onClick={handleRequestExit}
            className="text-xs font-bold text-indigo-600 dark:text-indigo-400 hover:underline inline-flex items-center gap-1 cursor-pointer"
          >
            <span>Вернуться к содержанию курса</span>
            <ArrowRight size={14} />
          </button>
        </footer>

        {/* Confirm Exit Quiz Modal */}
        {isExitModalOpen && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4 animate-in fade-in duration-200">
            <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 max-w-md w-full shadow-2xl space-y-5 text-center">
              <div className="w-14 h-14 rounded-2xl bg-amber-100 dark:bg-amber-950 text-amber-600 dark:text-amber-400 flex items-center justify-center mx-auto shadow-inner">
                <AlertTriangle size={30} />
              </div>

              <div className="space-y-1.5">
                <h3 className="text-lg font-black text-slate-900 dark:text-white">
                  Покинуть тестирование?
                </h3>
                <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
                  В уроке содержатся проверочные задания. Если вы покинете урок сейчас, неотвеченные вопросы будут оценены в 0 баллов.
                </p>
              </div>

              <div className="grid grid-cols-2 gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => {
                    setIsExitModalOpen(false);
                    setPendingLessonId(null);
                  }}
                  className="px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 hover:bg-slate-50 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 text-xs font-bold transition-all shadow-2xs cursor-pointer"
                >
                  Продолжить тест
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setIsExitModalOpen(false);
                    if (pendingLessonId) {
                      const target = pendingLessonId;
                      setPendingLessonId(null);
                      router.push(`/lessons/${target}`);
                    } else {
                      performExit();
                    }
                  }}
                  className="px-4 py-2.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold transition-all shadow-md shadow-rose-600/20 cursor-pointer"
                >
                  Завершить досрочно
                </button>
              </div>
            </div>
          </div>
        )}

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
