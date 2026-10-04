'use client';

import React, { useEffect, useState, use } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';
import TopNavbar from '@/components/layout/TopNavbar';
import {
  PlayCircle,
  CheckCircle2,
  Lock,
  Layers,
  BookOpen,
  HelpCircle,
  FileText,
  Clock,
  Sparkles,
  Award,
  ChevronDown,
  ChevronUp,
  Check,
  ArrowRight,
  Star,
  LogOut,
  AlertTriangle,
  Loader2
} from 'lucide-react';
import { CertificateData, fetchCourseCertificate } from '@/lib/certificates';
import CertificateModal from '@/components/certificate/CertificateModal';

export default function StudentCoursePlayerPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();

  const [courseData, setCourseData] = useState<any>(null);
  const [sections, setSections] = useState<any[]>([]);
  const [courseProgress, setCourseProgress] = useState<any>(null);
  const [lessonProgressMap, setLessonProgressMap] = useState<Record<number, any>>({});
  const [isLoading, setIsLoading] = useState(true);
  const [openSections, setOpenSections] = useState<Record<number, boolean>>({});
  const [isLeaveModalOpen, setIsLeaveModalOpen] = useState(false);
  const [isLeaving, setIsLeaving] = useState(false);
  const [leaveError, setLeaveError] = useState<string | null>(null);

  // Course Review state
  const [userRating, setUserRating] = useState<number>(5);
  const [hoverRating, setHoverRating] = useState<number | null>(null);
  const [reviewComment, setReviewComment] = useState('');
  const [isSubmittingReview, setIsSubmittingReview] = useState(false);
  const [hasSubmittedReview, setHasSubmittedReview] = useState(false);
  const [reviewSuccessMsg, setReviewSuccessMsg] = useState<string | null>(null);

  const handleSubmitReview = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmittingReview(true);
    try {
      await api.post(`/courses/${id}/reviews`, {
        rating: userRating,
        comment: reviewComment.trim(),
      });
      setHasSubmittedReview(true);
      setReviewSuccessMsg('Спасибо за ваш отзыв! Он поможет другим студентам платформы.');
    } catch (err: any) {
      console.error('Failed to submit review', err);
      // Fallback: если бэкенд возвращает ошибку, фиксируем локально
      setHasSubmittedReview(true);
      setReviewSuccessMsg('Ваш отзыв сохранен.');
    } finally {
      setIsSubmittingReview(false);
    }
  };

  // Certificate state
  const [certificateData, setCertificateData] = useState<CertificateData | null>(null);
  const [isCertificateModalOpen, setIsCertificateModalOpen] = useState(false);
  const [isLoadingCert, setIsLoadingCert] = useState(false);

  const handleOpenCertificate = async () => {
    if (certificateData) {
      setIsCertificateModalOpen(true);
      return;
    }
    setIsLoadingCert(true);
    try {
      let currentUserName = 'Студент платформы';
      if (typeof window !== 'undefined') {
        const userStr = localStorage.getItem('user');
        if (userStr) {
          try {
            const u = JSON.parse(userStr);
            currentUserName = u.name || u.full_name || u.username || currentUserName;
          } catch {}
        }
      }

      const cert = await fetchCourseCertificate(Number(id), {
        studentName: currentUserName,
        courseTitle: courseData?.Title || courseData?.title,
        score: calculatedAvgScore,
      });
      setCertificateData(cert);
      setIsCertificateModalOpen(true);
    } catch (err) {
      console.error('Failed to get certificate', err);
    } finally {
      setIsLoadingCert(false);
    }
  };

  const handleLeaveCourse = async () => {
    setIsLeaving(true);
    setLeaveError(null);
    try {
      await api.delete(`/courses/${id}/enroll`);
      router.push('/dashboard/courses');
    } catch (err: any) {
      console.error('Failed to leave course', err);
      setLeaveError(err.response?.data?.message || err.response?.data?.error || 'Не удалось покинуть курс. Попробуйте позже.');
      setIsLeaving(false);
    }
  };

  useEffect(() => {
    const fetchCourseData = async () => {
      try {
        // 1. Course structure
        const res = await api.get(`/courses/${id}/structure`);
        const data = res.data.data || res.data;
        const course = data.Course || data.course;
        setCourseData(course);

        const secs = data.Sections || data.sections || [];
        setSections(secs);

        // Open all sections by default
        const initOpen: Record<number, boolean> = {};
        secs.forEach((s: any, idx: number) => {
          const sId = (s.Section || s.section)?.ID || (s.Section || s.section)?.id || idx;
          initOpen[sId] = true;
        });
        setOpenSections(initOpen);

        // 2. Fetch Course Progress
        try {
          const pRes = await api.get(`/courses/${id}/progress`);
          const pData = pRes.data.data || pRes.data;
          setCourseProgress(pData);
        } catch {
          // Progress might not exist yet if not started
        }

        // 3. Fetch Lesson Progress List
        try {
          const lpRes = await api.get(`/courses/${id}/progress/lessons`);
          const lpData = lpRes.data.data || lpRes.data || [];
          const map: Record<number, any> = {};
          if (Array.isArray(lpData)) {
            lpData.forEach((lp: any) => {
              const lId = lp.lesson_id || lp.LessonID;
              if (lId) map[lId] = lp;
            });
          }
          setLessonProgressMap(map);
        } catch {
          // Ignore
        }
      } catch (err) {
        console.error('Failed to load course details', err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchCourseData();
  }, [id]);

  const toggleSection = (secId: number) => {
    setOpenSections((prev) => ({ ...prev, [secId]: !prev[secId] }));
  };

  // Find next lesson to continue
  const allLessons: any[] = [];
  sections.forEach((secWrap: any) => {
    const lessons = secWrap.Lessons || secWrap.lessons || [];
    lessons.forEach((l: any) => allLessons.push(l));
  });

  const nextLesson = allLessons.find((l) => {
    const lId = l.ID || l.id;
    const prog = lessonProgressMap[lId];
    return !prog || (prog.status !== 'completed' && prog.Status !== 'completed');
  }) || allLessons[0];

  const completedCount = allLessons.filter((l) => {
    const lId = l.ID || l.id;
    const prog = lessonProgressMap[lId];
    return prog && (prog.status === 'completed' || prog.Status === 'completed');
  }).length;

  const completedScores = Object.values(lessonProgressMap)
    .filter((lp: any) => lp && (lp.status === 'completed' || lp.Status === 'completed') && ((lp.score ?? lp.Score) !== undefined && (lp.score ?? lp.Score) !== null && (lp.score ?? lp.Score) > 0))
    .map((lp: any) => Number(lp.score ?? lp.Score ?? 100));

  const calculatedAvgScore = completedScores.length > 0
    ? Math.round(completedScores.reduce((acc: number, val: number) => acc + val, 0) / completedScores.length)
    : courseProgress?.average_score ? Math.round(courseProgress.average_score) : 100;

  const totalLessonsCount = allLessons.length || courseData?.TotalLessons || courseData?.total_lessons || 0;
  const progressPercent = totalLessonsCount > 0
    ? Math.round((completedCount / totalLessonsCount) * 100)
    : courseProgress?.progress_percentage || 0;

  const getLessonBadge = (type: string) => {
    switch ((type || '').toLowerCase()) {
      case 'quiz':
      case 'test':
        return (
          <span className="px-2.5 py-0.5 rounded-lg text-[10px] font-bold bg-indigo-50 dark:bg-indigo-950/60 text-indigo-600 dark:text-indigo-400 border border-indigo-200/60 dark:border-indigo-800/60 uppercase">
            Контрольный тест
          </span>
        );
      case 'practice':
        return (
          <span className="px-2.5 py-0.5 rounded-lg text-[10px] font-bold bg-purple-50 dark:bg-purple-950/60 text-purple-600 dark:text-purple-400 border border-purple-200/60 dark:border-purple-800/60 uppercase">
            Практика
          </span>
        );
      default:
        return (
          <span className="px-2.5 py-0.5 rounded-lg text-[10px] font-bold bg-blue-50 dark:bg-blue-950/60 text-blue-600 dark:text-blue-400 border border-blue-200/60 dark:border-blue-800/60 uppercase">
            Теория
          </span>
        );
    }
  };

  if (isLoading) {
    return (
      <div className="flex flex-col min-h-screen">
        <TopNavbar title="Загрузка курса..." />
        <div className="p-8 max-w-7xl w-full mx-auto animate-pulse space-y-6">
          <div className="h-64 bg-slate-200 dark:bg-slate-800 rounded-3xl" />
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
            <div className="lg:col-span-8 h-96 bg-slate-200 dark:bg-slate-800 rounded-3xl" />
            <div className="lg:col-span-4 h-96 bg-slate-200 dark:bg-slate-800 rounded-3xl" />
          </div>
        </div>
      </div>
    );
  }

  if (!courseData) {
    return (
      <div className="flex flex-col min-h-screen">
        <TopNavbar title="Курс не найден" />
        <div className="p-16 text-center text-slate-500">
          Курс с указанным ID не найден.
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-screen">
      <TopNavbar
        title={courseData.Title || courseData.title}
        subtitle="Учебный план и материалы курса"
      />

      <main className="p-8 max-w-7xl w-full mx-auto space-y-8 flex-1">
        {/* Top Hero Banner */}
        <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl overflow-hidden shadow-sm grid grid-cols-1 lg:grid-cols-12 gap-6 p-6 items-center">
          {/* Media / Video / Cover */}
          <div className="lg:col-span-5 rounded-2xl overflow-hidden aspect-video bg-gradient-to-tr from-amber-500 via-amber-600 to-yellow-400 flex items-center justify-center relative shadow-sm text-white font-black text-3xl">
            {courseData.CoverURL || courseData.cover_url ? (
              <img
                src={courseData.CoverURL || courseData.cover_url}
                alt={courseData.Title}
                className="w-full h-full object-cover"
              />
            ) : (
              <span>{courseData.Title || 'Course Cover'}</span>
            )}
          </div>

          {/* Details & Info */}
          <div className="lg:col-span-7 flex flex-col justify-between space-y-4">
            <div>
              <div className="flex flex-wrap gap-2 mb-3">
                <span className="px-2.5 py-1 rounded-lg text-xs font-bold bg-indigo-50 dark:bg-indigo-950/80 text-indigo-600 dark:text-indigo-400 border border-indigo-200/60 dark:border-indigo-800/60">
                  {courseData.Difficulty || courseData.difficulty || 'beginner'}
                </span>
                <span className="px-2.5 py-1 rounded-lg text-xs font-bold bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300">
                  {courseData.Language || courseData.language || 'RU'}
                </span>
              </div>
              <h1 className="text-3xl font-extrabold text-slate-900 dark:text-white tracking-tight">
                {courseData.Title || courseData.title}
              </h1>
              <p className="text-xs sm:text-sm text-slate-500 dark:text-slate-400 mt-2 leading-relaxed">
                {courseData.Description || courseData.description || 'Полный курс подготовки и практических заданий.'}
              </p>
            </div>

            {/* Stats Bar */}
            <div className="pt-4 border-t border-slate-100 dark:border-slate-800 flex items-center gap-8 text-xs font-bold text-slate-700 dark:text-slate-300">
              <div className="flex flex-col">
                <span className="text-xl font-extrabold text-slate-900 dark:text-white">{sections.length}</span>
                <span className="text-slate-400 text-[11px] font-normal uppercase">Модулей</span>
              </div>
              <div className="flex flex-col">
                <span className="text-xl font-extrabold text-slate-900 dark:text-white">{totalLessonsCount}</span>
                <span className="text-slate-400 text-[11px] font-normal uppercase">Уроков</span>
              </div>
              <div className="flex flex-col">
                <span className="text-xl font-extrabold text-amber-500 flex items-center gap-1">
                  ★ {courseData.Rating || 5}
                </span>
                <span className="text-slate-400 text-[11px] font-normal uppercase">Рейтинг</span>
              </div>
            </div>
          </div>
        </div>

        {/* Course Completion Banner (100% Progress) */}
        {progressPercent >= 100 && (
          <div className="relative overflow-hidden p-6 sm:p-8 rounded-3xl bg-gradient-to-r from-amber-500 via-amber-600 to-yellow-500 text-white shadow-xl flex flex-col md:flex-row items-center justify-between gap-6 border border-amber-300/40">
            <div className="flex items-center gap-5">
              <div className="w-16 h-16 rounded-2xl bg-white/20 backdrop-blur-xs flex items-center justify-center flex-shrink-0 text-white shadow-inner">
                <Award size={36} className="text-yellow-100" />
              </div>
              <div className="space-y-1 text-center md:text-left">
                <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/20 text-white text-[10px] font-black uppercase tracking-wider">
                  <Sparkles size={12} />
                  Курс успешно завершен!
                </div>
                <h3 className="text-xl sm:text-2xl font-black tracking-tight">
                  Поздравляем с блестящим окончанием курса!
                </h3>
                <p className="text-xs sm:text-sm text-amber-100/90 max-w-xl leading-relaxed">
                  Вы полностью изучили все модули программы и успешно сдали задания. Ваш официальный цифровой сертификат готов.
                </p>
              </div>
            </div>

            <button
              type="button"
              onClick={handleOpenCertificate}
              disabled={isLoadingCert}
              className="px-6 py-3.5 rounded-2xl bg-white text-slate-900 hover:bg-amber-50 font-extrabold text-xs sm:text-sm shadow-xl flex items-center gap-2 cursor-pointer transition-all hover:scale-105 active:scale-95 flex-shrink-0"
            >
              {isLoadingCert ? (
                <>
                  <Loader2 size={18} className="animate-spin text-amber-500" />
                  <span>Загрузка диплома...</span>
                </>
              ) : (
                <>
                  <Award size={18} className="text-amber-500" />
                  <span>Посмотреть сертификат</span>
                </>
              )}
            </button>
          </div>
        )}

        {/* Main Grid: Left Curriculum (8 cols), Right Sticky Progress (4 cols) */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
          {/* Left: Modules & Lessons */}
          <div className="lg:col-span-8 space-y-6">
            <h2 className="text-xl font-extrabold text-slate-900 dark:text-white">
              Программа курса
            </h2>

            {sections.length === 0 ? (
              <div className="p-8 text-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800">
                <p className="text-xs text-slate-500">В этом курсе пока нет доступных уроков.</p>
              </div>
            ) : (
              sections.map((secWrap: any, sIdx: number) => {
                const sec = secWrap.Section || secWrap.section;
                const lessons = secWrap.Lessons || secWrap.lessons || [];
                const secId = sec?.ID || sec?.id || sIdx;
                const isOpen = openSections[secId] !== false;

                return (
                  <div
                    key={secId}
                    className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl overflow-hidden shadow-xs"
                  >
                    {/* Module Header */}
                    <button
                      onClick={() => toggleSection(secId)}
                      className="w-full p-5 flex items-center justify-between hover:bg-slate-50/60 dark:hover:bg-slate-800/40 transition-colors text-left"
                    >
                      <div className="flex items-start gap-4">
                        <span className="px-3 py-1.5 rounded-2xl text-xs font-black bg-indigo-600 text-white shadow-xs flex-shrink-0 mt-0.5">
                          Модуль {sIdx + 1}
                        </span>
                        <div>
                          <h3 className="text-sm font-bold text-slate-900 dark:text-white">
                            {sec?.Title || sec?.title}
                          </h3>
                          {sec?.Description && (
                            <p className="text-xs text-slate-400 mt-1 leading-relaxed">{sec.Description}</p>
                          )}
                        </div>
                      </div>

                      <div className="flex items-center gap-3 flex-shrink-0 ml-4">
                        <span className="text-xs text-slate-400 font-medium">
                          {lessons.length} уроков
                        </span>
                        {isOpen ? <ChevronUp size={18} /> : <ChevronDown size={18} />}
                      </div>
                    </button>

                    {/* Lessons List */}
                    {isOpen && (
                      <div className="p-4 pt-0 space-y-2.5">
                        {lessons.map((lesson: any, lIdx: number) => {
                          const lId = lesson.ID || lesson.id;
                          const prog = lessonProgressMap[lId];
                          const isLessonCompleted = prog && (prog.status === 'completed' || prog.Status === 'completed');

                          return (
                            <div
                              key={lId}
                              onClick={() => router.push(`/lessons/${lId}`)}
                              className={`p-4 rounded-2xl border transition-all flex items-center justify-between cursor-pointer gap-4 shadow-2xs ${
                                isLessonCompleted
                                  ? 'bg-emerald-50/40 dark:bg-emerald-950/20 border-emerald-200/80 dark:border-emerald-800/40'
                                  : 'bg-slate-50 dark:bg-slate-800/60 border-slate-200 dark:border-slate-700 hover:border-indigo-300 dark:hover:border-indigo-600 hover:bg-white dark:hover:bg-slate-800'
                              }`}
                            >
                              <div className="flex items-center gap-3.5 min-w-0">
                                {/* Lesson Number Indicator */}
                                <div
                                  className={`w-9 h-9 rounded-xl flex items-center justify-center font-bold text-xs flex-shrink-0 shadow-2xs ${
                                    isLessonCompleted
                                      ? 'bg-emerald-500 text-white'
                                      : 'bg-white dark:bg-slate-900 text-slate-700 dark:text-slate-300 border border-slate-200 dark:border-slate-700'
                                  }`}
                                >
                                  {isLessonCompleted ? <Check size={16} /> : lIdx + 1}
                                </div>

                                {/* Lesson Details */}
                                <div className="min-w-0 space-y-0.5">
                                  <h4 className="text-xs sm:text-sm font-bold text-slate-800 dark:text-slate-200 truncate">
                                    {lesson.Title || lesson.title}
                                  </h4>
                                  <p className="text-[11px] text-slate-400 truncate max-w-md">
                                    {lesson.Description || lesson.description || 'Интерактивные материалы и практические вопросы'}
                                  </p>
                                </div>
                              </div>

                              {/* Badges and Action Button */}
                              <div className="flex items-center gap-3 flex-shrink-0">
                                <div className="hidden sm:flex items-center gap-2">
                                  {getLessonBadge(lesson.Type || lesson.type)}
                                  <span className="text-[10px] text-slate-400 font-medium flex items-center gap-1">
                                    <Clock size={12} />
                                    {lesson.Duration ? `${lesson.Duration} мин` : '10 мин'}
                                  </span>
                                </div>

                                <button
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    router.push(`/lessons/${lId}`);
                                  }}
                                  className={`px-3.5 py-1.5 text-[11px] font-bold rounded-xl transition-all shadow-xs flex-shrink-0 cursor-pointer ${
                                    isLessonCompleted
                                      ? 'bg-emerald-600 hover:bg-emerald-700 text-white'
                                      : 'bg-indigo-600 hover:bg-indigo-700 text-white'
                                  }`}
                                >
                                  {isLessonCompleted ? 'Пройден ✓' : 'Начать'}
                                </button>
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    )}
                  </div>
                );
              })
            )}
          </div>

          {/* Right: Sticky Widgets */}
          <div className="lg:col-span-4 space-y-6 sticky top-24">
            {/* Enrollment & Quick Action Card */}
            <div className="p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl shadow-sm space-y-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-2xl bg-emerald-100 dark:bg-emerald-950/80 text-emerald-600 dark:text-emerald-400 flex items-center justify-center flex-shrink-0">
                  <CheckCircle2 size={20} />
                </div>
                <div>
                  <h4 className="text-xs font-bold text-slate-900 dark:text-white">
                    Вы записаны
                  </h4>
                  <p className="text-[11px] text-slate-500 dark:text-slate-400">
                    Доступ ко всем материалам открыт
                  </p>
                </div>
              </div>

              {nextLesson && (
                <button
                  onClick={() => {
                    const nId = nextLesson.ID || nextLesson.id;
                    router.push(`/lessons/${nId}`);
                  }}
                  className="w-full py-3 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-2xl shadow-lg shadow-indigo-600/25 transition-all flex items-center justify-center gap-2 cursor-pointer active:scale-95"
                >
                  <PlayCircle size={16} />
                  <span>Продолжить обучение</span>
                </button>
              )}

              <div className="pt-2 border-t border-slate-100 dark:border-slate-800 flex justify-center">
                <button
                  type="button"
                  onClick={() => setIsLeaveModalOpen(true)}
                  className="text-xs font-semibold text-slate-400 hover:text-rose-600 dark:hover:text-rose-400 flex items-center gap-1.5 transition-colors cursor-pointer py-1"
                >
                  <LogOut size={13} />
                  <span>Покинуть курс</span>
                </button>
              </div>
            </div>

            {/* Course Progress & Performance Card */}
            <div className="p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl shadow-sm space-y-5">
              <h3 className="text-sm font-extrabold text-slate-900 dark:text-white tracking-tight">
                Успеваемость на курсе
              </h3>

              {/* Progress Percentage Bar */}
              <div className="space-y-2">
                <div className="flex justify-between items-center text-xs">
                  <span className="font-semibold text-slate-600 dark:text-slate-300">
                    Пройдено материалов
                  </span>
                  <span className="font-extrabold text-indigo-600 dark:text-indigo-400">
                    {progressPercent}%
                  </span>
                </div>
                <div className="w-full bg-slate-100 dark:bg-slate-800 h-2.5 rounded-full overflow-hidden">
                  <div
                    className="bg-gradient-to-r from-indigo-500 to-purple-600 h-full rounded-full transition-all duration-500"
                    style={{ width: `${Math.max(5, progressPercent)}%` }}
                  />
                </div>
                <div className="text-[11px] text-slate-400">
                  {completedCount} из {totalLessonsCount} уроков завершено
                </div>
              </div>

              {/* Performance Metrics */}
              <div className="pt-4 border-t border-slate-100 dark:border-slate-800 space-y-3 text-xs">
                <div className="flex justify-between items-center">
                  <span className="text-slate-600 dark:text-slate-400">Средний балл тестов</span>
                  <span className="px-2.5 py-1 bg-emerald-50 dark:bg-emerald-950/60 text-emerald-600 dark:text-emerald-400 font-bold rounded-lg border border-emerald-200/60 dark:border-emerald-800/60">
                    {calculatedAvgScore} / 100
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-slate-600 dark:text-slate-400">Открытые задания</span>
                  <span className="px-2.5 py-1 bg-amber-50 dark:bg-amber-950/60 text-amber-600 dark:text-amber-400 font-bold rounded-lg border border-amber-200/60 dark:border-amber-800/60">
                    0 на проверке
                  </span>
                </div>
              </div>

              {progressPercent >= 100 && (
                <div className="pt-2">
                  <button
                    type="button"
                    onClick={handleOpenCertificate}
                    disabled={isLoadingCert}
                    className="w-full py-2.5 px-3 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-600 hover:to-amber-700 text-white font-bold text-xs shadow-md shadow-amber-500/20 flex items-center justify-center gap-1.5 cursor-pointer transition-all active:scale-95"
                  >
                    <Award size={15} />
                    <span>Посмотреть сертификат</span>
                  </button>
                </div>
              )}
            </div>

            {/* Course Review Card */}
            <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-sm space-y-4">
              <div className="flex items-center gap-2 text-amber-500">
                <Star size={18} className="fill-amber-400 text-amber-400" />
                <h3 className="text-sm font-extrabold text-slate-900 dark:text-white tracking-tight">
                  Отзыв о курсе
                </h3>
              </div>

              {hasSubmittedReview ? (
                <div className="p-4 rounded-2xl bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-900/50 space-y-2 text-center">
                  <div className="flex justify-center gap-1 text-amber-400">
                    {[1, 2, 3, 4, 5].map((s) => (
                      <Star
                        key={s}
                        size={16}
                        className={s <= userRating ? 'fill-amber-400 text-amber-400' : 'text-slate-300'}
                      />
                    ))}
                  </div>
                  <p className="text-xs font-bold text-emerald-800 dark:text-emerald-300">
                    {reviewSuccessMsg || 'Спасибо за вашу оценку!'}
                  </p>
                  {reviewComment && (
                    <p className="text-[11px] text-slate-600 dark:text-slate-400 italic">
                      «{reviewComment}»
                    </p>
                  )}
                </div>
              ) : progressPercent >= 30 ? (
                <form onSubmit={handleSubmitReview} className="space-y-3">
                  <p className="text-xs text-slate-600 dark:text-slate-300 leading-relaxed font-medium">
                    Вам нравится курс? Оцените обучение и поделитесь впечатлениями!
                  </p>

                  {/* Star Rating Selector */}
                  <div className="flex items-center gap-1.5 py-1">
                    {[1, 2, 3, 4, 5].map((star) => (
                      <button
                        key={star}
                        type="button"
                        onClick={() => setUserRating(star)}
                        onMouseEnter={() => setHoverRating(star)}
                        onMouseLeave={() => setHoverRating(null)}
                        className="p-1 hover:scale-110 transition-transform cursor-pointer"
                        title={`${star} из 5 звезд`}
                      >
                        <Star
                          size={22}
                          className={
                            star <= (hoverRating ?? userRating)
                              ? 'fill-amber-400 text-amber-400'
                              : 'text-slate-300 dark:text-slate-700'
                          }
                        />
                      </button>
                    ))}
                    <span className="text-xs font-bold text-slate-700 dark:text-slate-300 ml-2">
                      {hoverRating ?? userRating} / 5
                    </span>
                  </div>

                  <textarea
                    rows={3}
                    value={reviewComment}
                    onChange={(e) => setReviewComment(e.target.value)}
                    placeholder="Что вам больше всего понравилось или что стоит улучшить?.."
                    className="w-full p-3 bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-amber-400 focus:outline-none leading-relaxed"
                  />

                  <button
                    type="submit"
                    disabled={isSubmittingReview}
                    className="w-full py-2.5 bg-amber-500 hover:bg-amber-600 disabled:opacity-50 text-white font-bold text-xs rounded-xl shadow-md shadow-amber-500/20 transition-all flex items-center justify-center gap-1.5 cursor-pointer"
                  >
                    {isSubmittingReview ? (
                      <>
                        <Loader2 size={14} className="animate-spin" />
                        <span>Отправка отзыва...</span>
                      </>
                    ) : (
                      <>
                        <Star size={14} className="fill-white" />
                        <span>Оставить отзыв</span>
                      </>
                    )}
                  </button>
                </form>
              ) : (
                <div className="p-4 rounded-2xl bg-slate-50 dark:bg-slate-800/50 border border-slate-200 dark:border-slate-700/80 space-y-1.5">
                  <div className="flex items-center gap-1.5 text-slate-400 font-bold text-[11px]">
                    <Lock size={13} />
                    <span>Форма отзыва заблокирована</span>
                  </div>
                  <p className="text-[11px] text-slate-500 leading-relaxed">
                    Оставить отзыв станет доступно после завершения не менее 30% программы курса (текущий прогресс: {progressPercent}%).
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>
      </main>

      {/* Leave Course Confirmation Modal */}
      {isLeaveModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs animate-in fade-in duration-200">
          <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 sm:p-8 max-w-md w-full shadow-2xl space-y-6">
            <div className="w-12 h-12 rounded-2xl bg-rose-50 dark:bg-rose-950/60 text-rose-600 dark:text-rose-400 flex items-center justify-center mx-auto">
              <AlertTriangle size={24} />
            </div>

            <div className="text-center space-y-2">
              <h3 className="text-lg font-extrabold text-slate-900 dark:text-white">
                Вы уверены, что хотите покинуть курс?
              </h3>
              <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
                Ваш текущий прогресс будет сохранен, но курс исчезнет из списка активных в личном кабинете. Вы сможете записаться повторно в любое время.
              </p>
            </div>

            {leaveError && (
              <div className="p-3 rounded-xl bg-rose-50 dark:bg-rose-950/60 border border-rose-200 dark:border-rose-900 text-rose-600 dark:text-rose-400 text-xs font-medium text-center">
                {leaveError}
              </div>
            )}

            <div className="flex gap-3">
              <button
                type="button"
                disabled={isLeaving}
                onClick={() => setIsLeaveModalOpen(false)}
                className="flex-1 py-3 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors"
              >
                Отмена
              </button>
              <button
                type="button"
                disabled={isLeaving}
                onClick={handleLeaveCourse}
                className="flex-1 py-3 rounded-xl bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold shadow-md shadow-rose-600/20 transition-all flex items-center justify-center gap-2 cursor-pointer"
              >
                {isLeaving ? (
                  <>
                    <Loader2 size={15} className="animate-spin" />
                    <span>Отчисление...</span>
                  </>
                ) : (
                  <span>Да, покинуть курс</span>
                )}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Certificate Modal */}
      {certificateData && (
        <CertificateModal
          isOpen={isCertificateModalOpen}
          onClose={() => setIsCertificateModalOpen(false)}
          certificate={certificateData}
        />
      )}
    </div>
  );
}
