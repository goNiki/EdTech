'use client';

import React, { useEffect, useState, useMemo, useCallback } from 'react';
import { api } from '@/lib/api';
import TopNavbar from '@/components/layout/TopNavbar';
import PendingHomeworksQueue, { PendingHWItem } from '@/components/teacher/PendingHomeworksQueue';
import ModalGradeHW from '@/components/teacher/ModalGradeHW';
import { CheckSquare, Sparkles, CheckCircle2, Filter, Layers, BookOpen, AlertCircle } from 'lucide-react';

interface CourseItem {
  id: number;
  title: string;
  slug?: string;
}

export default function TeacherGlobalGradingPage() {
  const [courses, setCourses] = useState<CourseItem[]>([]);
  const [selectedCourseId, setSelectedCourseId] = useState<number | null>(null); // null = "Все курсы"
  const [pendingHWs, setPendingHWs] = useState<PendingHWItem[]>([]);
  const [courseCounts, setCourseCounts] = useState<Record<number, number>>({});
  const [isLoading, setIsLoading] = useState(true);
  const [gradingHW, setGradingHW] = useState<PendingHWItem | null>(null);
  const [toastMsg, setToastMsg] = useState<string | null>(null);

  const showToast = (msg: string) => {
    setToastMsg(msg);
    setTimeout(() => setToastMsg(null), 3500);
  };

  const fetchGradingData = useCallback(async () => {
    setIsLoading(true);
    try {
      // 1. Fetch teacher courses
      const cRes = await api.get('/courses/my');
      const cData = cRes.data.data || cRes.data;
      const courseList: CourseItem[] = cData.courses || (Array.isArray(cData) ? cData : []);
      setCourses(courseList);

      // 2. Try global backend endpoint BE-027
      const globalQuery = selectedCourseId ? `?course_id=${selectedCourseId}` : '';
      const globalRes = await api.get(`/teacher/grading/pending${globalQuery}`).catch(() => null);

      if (globalRes && globalRes.data) {
        const payload = globalRes.data.data || globalRes.data;
        const items: PendingHWItem[] = payload.items || (Array.isArray(payload) ? payload : []);
        setPendingHWs(items);

        // Parse course counts if summary provided
        if (payload.courses_summary && Array.isArray(payload.courses_summary)) {
          const counts: Record<number, number> = {};
          payload.courses_summary.forEach((s: any) => {
            counts[s.course_id] = s.pending_count || 0;
          });
          setCourseCounts(counts);
        } else {
          // Count from items
          const counts: Record<number, number> = {};
          items.forEach((it) => {
            if (it.course_id) {
              counts[it.course_id] = (counts[it.course_id] || 0) + 1;
            }
          });
          setCourseCounts(counts);
        }
      } else {
        // Resilient Fallback: parallel query across all author courses
        const counts: Record<number, number> = {};
        const allItems: PendingHWItem[] = [];

        await Promise.all(
          courseList.map(async (c) => {
            try {
              const res = await api.get(`/courses/${c.id}/grading/pending`);
              const pData = res.data.data || res.data;
              const courseItems: PendingHWItem[] = pData.items || (Array.isArray(pData) ? pData : []);

              counts[c.id] = courseItems.length;
              courseItems.forEach((item) => {
                allItems.push({
                  ...item,
                  course_id: item.course_id || c.id,
                  course_title: item.course_title || c.title,
                });
              });
            } catch {
              counts[c.id] = 0;
            }
          })
        );

        setCourseCounts(counts);

        if (selectedCourseId) {
          setPendingHWs(allItems.filter((it) => it.course_id === selectedCourseId));
        } else {
          setPendingHWs(allItems);
        }
      }
    } catch (err) {
      console.error('Failed to load grading queue', err);
    } finally {
      setIsLoading(false);
    }
  }, [selectedCourseId]);

  useEffect(() => {
    fetchGradingData();
  }, [fetchGradingData]);

  // Total count of pending items across all courses
  const totalPendingCount = useMemo(() => {
    return Object.values(courseCounts).reduce((acc, cnt) => acc + cnt, 0);
  }, [courseCounts]);

  // Sort courses: courses with active homeworks first
  const sortedCourses = useMemo(() => {
    return [...courses].sort((a, b) => {
      const countA = courseCounts[a.id] || 0;
      const countB = courseCounts[b.id] || 0;
      if (countB !== countA) return countB - countA;
      return a.title.localeCompare(b.title);
    });
  }, [courses, courseCounts]);

  // Courses with at least 1 pending homework for quick filter pills
  const activeCourses = useMemo(() => {
    return sortedCourses.filter((c) => (courseCounts[c.id] || 0) > 0);
  }, [sortedCourses, courseCounts]);

  const handleGradedSuccess = () => {
    if (!gradingHW) return;
    const gradedId = gradingHW.id || gradingHW.answer_id;
    const courseId = gradingHW.course_id;

    // Optimistic removal from queue
    setPendingHWs((prev) => prev.filter((it) => (it.id || it.answer_id) !== gradedId));

    // Decrement counts
    if (courseId) {
      setCourseCounts((prev) => ({
        ...prev,
        [courseId]: Math.max(0, (prev[courseId] || 0) - 1),
      }));
    }

    showToast('Оценка успешно выставлена и сохранена!');
    setGradingHW(null);
  };

  return (
    <div className="flex flex-col min-h-screen">
      <TopNavbar
        title="Проверка заданий"
        subtitle="Глобальная очередь ответов на открытые вопросы и практические работы"
      />

      {/* Floating Success Toast */}
      {toastMsg && (
        <div className="fixed bottom-6 right-6 z-50 bg-emerald-600 text-white px-5 py-3 rounded-2xl shadow-xl flex items-center gap-2.5 text-xs font-bold animate-in slide-in-from-bottom-3 duration-200">
          <CheckCircle2 size={16} />
          <span>{toastMsg}</span>
        </div>
      )}

      <main className="p-4 sm:p-8 max-w-7xl w-full mx-auto space-y-6 flex-1">
        {/* Header Controls & Filter Section */}
        <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-sm space-y-4">
          <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-4">
            <div className="space-y-1">
              <div className="flex items-center gap-2.5">
                <h2 className="text-xl font-extrabold text-slate-900 dark:text-white">
                  Очередь на проверку
                </h2>
                <span className="px-2.5 py-0.5 rounded-full text-xs font-black bg-indigo-100 dark:bg-indigo-950/80 text-indigo-700 dark:text-indigo-400">
                  Всего: {totalPendingCount}
                </span>
              </div>
              <p className="text-xs text-slate-500">
                Фильтрация по всем авторским курсам и оперативная оценка студенческих работ
              </p>
            </div>

            {/* Course Selector Dropdown */}
            <div className="flex items-center gap-2">
              <label className="text-xs font-bold text-slate-600 dark:text-slate-400 whitespace-nowrap">
                Курс:
              </label>
              <select
                value={selectedCourseId ?? ''}
                onChange={(e) => {
                  const val = e.target.value;
                  setSelectedCourseId(val === '' ? null : Number(val));
                }}
                className="px-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-bold text-slate-800 dark:text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 cursor-pointer min-w-[200px]"
              >
                <option value="">
                  🌟 Все курсы ({totalPendingCount})
                </option>
                {sortedCourses.map((c) => {
                  const cnt = courseCounts[c.id] || 0;
                  return (
                    <option key={c.id} value={c.id}>
                      {c.title} ({cnt})
                    </option>
                  );
                })}
              </select>
            </div>
          </div>

          {/* Quick Filter Pills */}
          {activeCourses.length > 0 && (
            <div className="pt-3 border-t border-slate-100 dark:border-slate-800 flex flex-wrap items-center gap-2">
              <span className="text-[11px] font-bold text-slate-400 mr-1 flex items-center gap-1">
                <Filter size={12} />
                <span>Быстрый фильтр:</span>
              </span>

              {/* "All" pill */}
              <button
                type="button"
                onClick={() => setSelectedCourseId(null)}
                className={`px-3 py-1 rounded-xl text-xs font-bold transition-all cursor-pointer ${
                  selectedCourseId === null
                    ? 'bg-indigo-600 text-white shadow-xs'
                    : 'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700'
                }`}
              >
                Все курсы ({totalPendingCount})
              </button>

              {/* Active course pills */}
              {activeCourses.map((c) => {
                const cnt = courseCounts[c.id] || 0;
                const isSelected = selectedCourseId === c.id;

                return (
                  <button
                    key={c.id}
                    type="button"
                    onClick={() => setSelectedCourseId(c.id)}
                    className={`px-3 py-1 rounded-xl text-xs font-bold transition-all flex items-center gap-1.5 cursor-pointer ${
                      isSelected
                        ? 'bg-indigo-600 text-white shadow-xs'
                        : 'bg-amber-50 dark:bg-amber-950/60 text-amber-800 dark:text-amber-300 border border-amber-200/80 dark:border-amber-800/80 hover:bg-amber-100 dark:hover:bg-amber-900/60'
                    }`}
                  >
                    <span>{c.title}</span>
                    <span
                      className={`px-1.5 py-0.2 rounded-full text-[10px] font-black ${
                        isSelected
                          ? 'bg-white/20 text-white'
                          : 'bg-amber-200/80 dark:bg-amber-800 text-amber-900 dark:text-amber-200'
                      }`}
                    >
                      {cnt}
                    </span>
                  </button>
                );
              })}
            </div>
          )}
        </div>

        {/* Pending Queue List */}
        {isLoading ? (
          <div className="p-16 text-center text-slate-400 animate-pulse space-y-3 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl">
            <Sparkles className="mx-auto text-indigo-500" size={26} />
            <p className="text-xs font-semibold">Загрузка очереди домашних заданий...</p>
          </div>
        ) : (
          <PendingHomeworksQueue
            items={pendingHWs}
            onOpenGradeModal={(hw) => setGradingHW(hw)}
            onSelectCourse={(courseId) => setSelectedCourseId(courseId)}
          />
        )}
      </main>

      {/* Grade HW Modal */}
      {gradingHW && (
        <ModalGradeHW
          hw={gradingHW}
          onClose={() => setGradingHW(null)}
          onGraded={handleGradedSuccess}
        />
      )}
    </div>
  );
}
