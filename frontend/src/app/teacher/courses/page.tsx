'use client';

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import TopNavbar from '@/components/layout/TopNavbar';
import {
  PlusCircle,
  BookOpen,
  Layers,
  Users,
  Settings,
  Sliders,
  Sparkles,
  ArrowRight
} from 'lucide-react';

export default function TeacherCoursesPage() {
  const router = useRouter();
  const [courses, setCourses] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchCourses = async () => {
      try {
        const res = await api.get('/courses/my');
        const data = res.data.data || res.data;
        const list = data.courses || data || [];
        setCourses(list);
      } catch (err) {
        console.error('Failed to load teacher courses', err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchCourses();
  }, []);

  return (
    <div className="flex flex-col min-h-screen">
      <TopNavbar title="Мои курсы (Автор)" subtitle="Управление вашими образовательными программами" />

      <main className="p-8 max-w-7xl w-full mx-auto space-y-8 flex-1">
        {/* Header Action Banner */}
        <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 p-6 rounded-3xl shadow-sm">
          <div>
            <h2 className="text-xl font-extrabold text-slate-900 dark:text-white">
              Авторская студия
            </h2>
            <p className="text-xs text-slate-500 mt-0.5">
              Создавайте модули, настраивайте тесты и отслеживайте прогресс студентов
            </p>
          </div>

          <button
            onClick={() => router.push('/teacher/courses/new')}
            className="px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs rounded-2xl transition-all shadow-md shadow-indigo-600/20 flex items-center justify-center gap-2"
          >
            <PlusCircle size={18} />
            <span>Создать новый курс</span>
          </button>
        </div>

        {/* Courses Grid */}
        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl h-80 p-5 animate-pulse space-y-4"
              >
                <div className="bg-slate-200 dark:bg-slate-800 h-36 rounded-2xl w-full" />
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
              У вас пока нет созданных курсов
            </h3>
            <p className="text-xs text-slate-500 leading-relaxed">
              Нажмите кнопку ниже, чтобы сконфигурировать ваш первый курс и собрать программу обучения.
            </p>
            <button
              onClick={() => router.push('/teacher/courses/new')}
              className="px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs rounded-2xl transition-all shadow-md flex items-center gap-2 mx-auto"
            >
              <PlusCircle size={16} />
              <span>Создать первый курс</span>
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {courses.map((course) => {
              const isPublished = course.status === 'published';
              const statusBadge = isPublished
                ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300'
                : 'bg-amber-100 dark:bg-amber-950 text-amber-700 dark:text-amber-300';
              const statusText = isPublished ? 'Опубликован' : 'Черновик';
              const visText = course.visibility === 'public' ? 'Публичный' : 'Приватный';

              return (
                <div
                  key={course.id}
                  className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl overflow-hidden shadow-sm hover:shadow-xl hover:border-indigo-300 dark:hover:border-indigo-800 transition-all duration-300 flex flex-col justify-between"
                >
                  <div>
                    {/* Media Header */}
                    <div className="h-40 bg-slate-100 dark:bg-slate-800 relative overflow-hidden">
                      {course.cover_url ? (
                        <img
                          src={course.cover_url}
                          alt={course.title}
                          className="w-full h-full object-cover"
                        />
                      ) : (
                        <div className="w-full h-full bg-gradient-to-tr from-indigo-600 to-purple-600 flex items-center justify-center text-white font-bold text-2xl">
                          {course.title?.slice(0, 2).toUpperCase()}
                        </div>
                      )}
                      <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent" />
                      <h3 className="absolute bottom-3 left-4 right-4 text-white font-extrabold text-base line-clamp-1 drop-shadow-md">
                        {course.title}
                      </h3>
                    </div>

                    {/* Card Body */}
                    <div className="p-5 space-y-4">
                      <div className="flex flex-wrap gap-2">
                        <span className={`px-2.5 py-0.5 rounded-lg text-[10px] font-bold ${statusBadge}`}>
                          {statusText}
                        </span>
                        <span className="text-[10px] text-slate-600 dark:text-slate-400 bg-slate-100 dark:bg-slate-800 px-2.5 py-0.5 rounded-lg font-medium">
                          {visText}
                        </span>
                        <span className="text-[10px] text-indigo-600 bg-indigo-50 dark:bg-indigo-950 px-2.5 py-0.5 rounded-lg font-bold">
                          Автор
                        </span>
                      </div>

                      <p className="text-xs text-slate-500 line-clamp-2">
                        {course.short_description || course.description || 'Описание курса...'}
                      </p>

                      <div className="pt-3 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between text-xs text-slate-500">
                        <div className="flex items-center gap-1.5 font-medium">
                          <Layers size={14} className="text-indigo-500" />
                          <span>{course.total_lessons || 0} уроков</span>
                        </div>
                        <div className="flex items-center gap-1.5 font-medium">
                          <Users size={14} className="text-purple-500" />
                          <span>{course.enrolled_count || 0} студентов</span>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="p-5 pt-0 flex gap-2">
                    <button
                      onClick={() => router.push(`/teacher/courses/${course.id}/curriculum`)}
                      className="w-full py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-2xl transition-all shadow-md flex items-center justify-center gap-2"
                    >
                      <Sliders size={14} />
                      <span>Управление курсом</span>
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
