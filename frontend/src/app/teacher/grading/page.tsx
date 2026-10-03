'use client';

import React, { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import TopNavbar from '@/components/layout/TopNavbar';
import PendingHomeworksQueue, { PendingHWItem } from '@/components/teacher/PendingHomeworksQueue';
import ModalGradeHW from '@/components/teacher/ModalGradeHW';
import { CheckSquare, Sparkles, CheckCircle2 } from 'lucide-react';

export default function TeacherGlobalGradingPage() {
  const [courses, setCourses] = useState<any[]>([]);
  const [selectedCourseId, setSelectedCourseId] = useState<number | null>(null);
  const [pendingHWs, setPendingHWs] = useState<PendingHWItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [gradingHW, setGradingHW] = useState<PendingHWItem | null>(null);
  const [toastMsg, setToastMsg] = useState<string | null>(null);

  const showToast = (msg: string) => {
    setToastMsg(msg);
    setTimeout(() => setToastMsg(null), 3000);
  };

  const fetchGradingQueue = async () => {
    try {
      // 1. Get author courses
      const cRes = await api.get('/courses/my');
      const cData = cRes.data.data || cRes.data;
      const list = cData.courses || cData || [];
      setCourses(list);

      const targetId = selectedCourseId || list[0]?.id;
      if (targetId) {
        setSelectedCourseId(targetId);
        const hwRes = await api.get(`/courses/${targetId}/grading/pending`).catch(() => null);
        if (hwRes) {
          const hwData = hwRes.data.data || hwRes.data;
          setPendingHWs(hwData.items || hwData || []);
        }
      }
    } catch (err) {
      console.error('Failed to load grading queue', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchGradingQueue();
  }, [selectedCourseId]);

  return (
    <div className="flex flex-col min-h-screen">
      <TopNavbar
        title="Проверка заданий"
        subtitle="Очередь ответов на открытые вопросы и практические работы"
      />

      <main className="p-8 max-w-7xl w-full mx-auto space-y-8 flex-1">
        {/* Header Controls */}
        <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-sm flex flex-col sm:flex-row justify-between sm:items-center gap-4">
          <div className="space-y-1">
            <h2 className="text-xl font-extrabold text-slate-900 dark:text-white">
              Очередь на проверку
            </h2>
            <p className="text-xs text-slate-500">
              Выберите курс для фильтрации сданных практических работ
            </p>
          </div>

          {courses.length > 0 && (
            <div className="flex items-center gap-2">
              <label className="text-xs font-bold text-slate-600 dark:text-slate-400">Курс:</label>
              <select
                value={selectedCourseId || ''}
                onChange={(e) => setSelectedCourseId(Number(e.target.value))}
                className="px-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-bold text-slate-800 dark:text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 cursor-pointer"
              >
                {courses.map((c) => (
                  <option key={c.id} value={c.id}>
                    {c.title}
                  </option>
                ))}
              </select>
            </div>
          )}
        </div>

        {/* Pending Queue List */}
        {isLoading ? (
          <div className="p-12 text-center text-slate-400 animate-pulse space-y-3">
            <Sparkles className="mx-auto text-indigo-500" size={24} />
            <p className="text-xs font-semibold">Загрузка очереди домашних заданий...</p>
          </div>
        ) : (
          <PendingHomeworksQueue
            items={pendingHWs}
            onOpenGradeModal={(hw) => setGradingHW(hw)}
          />
        )}
      </main>

      <ModalGradeHW
        hw={gradingHW}
        onClose={() => setGradingHW(null)}
        onGraded={() => {
          showToast('Оценка успешно выставлена и сохранена!');
          fetchGradingQueue();
        }}
      />

      {/* Toast Notification */}
      {toastMsg && (
        <div className="fixed bottom-6 right-6 z-50 bg-emerald-600 text-white text-xs font-bold px-4 py-3 rounded-2xl shadow-xl flex items-center gap-2 animate-bounce">
          <CheckCircle2 size={16} />
          <span>{toastMsg}</span>
        </div>
      )}
    </div>
  );
}
