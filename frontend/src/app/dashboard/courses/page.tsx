'use client';

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import TopNavbar from '@/components/layout/TopNavbar';
import { BookOpen, Layers, Play, CheckCircle2, ArrowRight } from 'lucide-react';

interface EnrolledCourse {
  id: number;
  title: string;
  slug: string;
  description: string;
  cover_url?: string;
  progress_percent?: number;
  total_lessons?: number;
  completed_lessons?: number;
  average_score?: number;
}

export default function MyCoursesPage() {
  const router = useRouter();
  const [courses, setCourses] = useState<EnrolledCourse[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchMyCourses = async () => {
      try {
        const res = await api.get('/courses/my');
        const data = res.data.data || res.data;
        const list = data.courses || data || [];
        setCourses(list);
      } catch (err) {
        console.error('Failed to fetch my courses', err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchMyCourses();
  }, []);

  return (
    <div className="flex flex-col min-h-screen">
      <TopNavbar title="Мои курсы" subtitle="Программы, на которые вы зачислены" />

      <main className="p-8 max-w-7xl w-full mx-auto space-y-8 flex-1">
        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl h-80 p-5 animate-pulse space-y-4"
              >
                <div className="bg-slate-200 dark:bg-slate-800 h-40 rounded-2xl w-full" />
                <div className="h-4 bg-slate-200 dark:bg-slate-800 rounded-md w-3/4" />
                <div className="h-3 bg-slate-200 dark:bg-slate-800 rounded-md w-1/2" />
              </div>
            ))}
          </div>
        ) : courses.length === 0 ? (
          <div className="p-16 text-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 shadow-sm space-y-4 max-w-lg mx-auto mt-12">
            <div className="w-16 h-16 rounded-3xl bg-indigo-50 dark:bg-slate-800 text-indigo-600 flex items-center justify-center mx-auto shadow-sm">
              <BookOpen size={32} />
            </div>
            <h3 className="text-xl font-extrabold text-slate-900 dark:text-white">
              Вы пока не записаны ни на один курс
            </h3>
            <p className="text-xs text-slate-500 leading-relaxed">
              Откройте наш каталог, выберите подходящую программу и начните обучение прямо сейчас!
            </p>
            <button
              onClick={() => router.push('/courses')}
              className="px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs rounded-2xl transition-all shadow-md flex items-center gap-2 mx-auto"
            >
              <span>Перейти в каталог курсов</span>
              <ArrowRight size={16} />
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {courses.map((course) => {
              const progress = course.progress_percent || 0;

              return (
                <div
                  key={course.id}
                  onClick={() => router.push(`/dashboard/courses/${course.id}`)}
                  className="group bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl overflow-hidden shadow-sm hover:shadow-xl hover:border-indigo-300 dark:hover:border-indigo-800 transition-all duration-300 flex flex-col cursor-pointer"
                >
                  <div className="h-44 bg-slate-100 dark:bg-slate-800 relative overflow-hidden">
                    {course.cover_url ? (
                      <img
                        src={course.cover_url}
                        alt={course.title}
                        className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                      />
                    ) : (
                      <div className="w-full h-full bg-gradient-to-tr from-indigo-600 to-purple-600 flex items-center justify-center text-white font-bold text-2xl">
                        {course.title.slice(0, 2).toUpperCase()}
                      </div>
                    )}
                    <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
                    <h3 className="absolute bottom-3 left-4 right-4 text-white font-extrabold text-base line-clamp-1 drop-shadow-md">
                      {course.title}
                    </h3>
                  </div>

                  <div className="p-5 flex-1 flex flex-col justify-between space-y-4">
                    {/* Progress Bar */}
                    <div className="space-y-1.5">
                      <div className="flex justify-between items-center text-xs font-semibold">
                        <span className="text-slate-500">Прогресс обучения</span>
                        <span className="font-extrabold text-indigo-600 dark:text-indigo-400">
                          {progress}%
                        </span>
                      </div>
                      <div className="w-full bg-slate-100 dark:bg-slate-800 h-2 rounded-full overflow-hidden">
                        <div
                          className="bg-indigo-600 h-full rounded-full transition-all duration-500"
                          style={{ width: `${progress}%` }}
                        />
                      </div>
                    </div>

                    <div className="pt-3 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between text-xs text-slate-500">
                      <div className="flex items-center gap-1.5">
                        <Layers size={14} className="text-indigo-500" />
                        <span>{course.total_lessons || 0} уроков</span>
                      </div>
                      {course.average_score !== undefined && (
                        <div className="flex items-center gap-1 font-bold text-emerald-600">
                          <CheckCircle2 size={14} />
                          <span>Балл: {course.average_score} / 100</span>
                        </div>
                      )}
                    </div>

                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        router.push(`/dashboard/courses/${course.id}`);
                      }}
                      className="w-full py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-2xl transition-all shadow-md flex items-center justify-center gap-2"
                    >
                      <Play size={14} fill="currentColor" />
                      <span>Продолжить обучение</span>
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </main>
    </div>
  );
}
