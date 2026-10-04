'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { useAuth } from '@/store/useAuth';
import { CourseCard } from '@/components/CourseCard';
import { BookOpen, GraduationCap, ArrowRight, Loader2 } from 'lucide-react';

export default function StudentDashboard() {
  const router = useRouter();
  const { user, viewMode, setViewMode, isLoading: isAuthLoading } = useAuth();
  const [courses, setCourses] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  // Role Guard / Redirect teachers to Author Studio
  const isTeacherRole = user && ['teacher', 'author', 'admin'].includes(user.role);

  useEffect(() => {
    if (!isAuthLoading && isTeacherRole && viewMode === 'teacher') {
      router.replace('/teacher/courses');
    }
  }, [user, viewMode, isAuthLoading, isTeacherRole, router]);

  useEffect(() => {
    const fetchMyCourses = async () => {
      try {
        const res = await api.get('/courses/my');
        const payload = res.data?.data || res.data;
        const courseList = payload?.courses || (Array.isArray(payload) ? payload : []);
        setCourses(courseList);
      } catch (err) {
        console.error('Failed to load user courses', err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchMyCourses();
  }, []);

  if (isAuthLoading || (isTeacherRole && viewMode === 'teacher')) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] space-y-4">
        <Loader2 className="w-8 h-8 animate-spin text-indigo-600" />
        <p className="text-sm font-semibold text-slate-600 dark:text-slate-400">
          Перенаправление в Авторскую студию...
        </p>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-extrabold text-gray-900 dark:text-white">Мое обучение</h1>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {[1, 2, 3].map(i => <div key={i} className="h-72 bg-gray-200 dark:bg-slate-800 animate-pulse rounded-2xl"></div>)}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Teacher View As Student Banner */}
      {isTeacherRole && viewMode === 'student' && (
        <div className="p-4 bg-gradient-to-r from-indigo-50 to-purple-50 dark:from-indigo-950/40 dark:to-purple-950/40 border border-indigo-200 dark:border-indigo-800/80 rounded-2xl flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 shadow-xs">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-xl bg-indigo-600 text-white flex items-center justify-center flex-shrink-0 shadow-xs">
              <GraduationCap size={20} />
            </div>
            <div>
              <p className="text-xs font-bold text-slate-800 dark:text-slate-200">
                Режим предварительного просмотра студента
              </p>
              <p className="text-[11px] text-slate-500 dark:text-slate-400">
                Вы вошли как преподаватель ({user?.role}), но сейчас просматриваете кабинет глазами ученика.
              </p>
            </div>
          </div>
          <button
            type="button"
            onClick={() => {
              setViewMode('teacher');
              router.push('/teacher/courses');
            }}
            className="px-3.5 py-1.5 bg-indigo-600 hover:bg-indigo-700 active:scale-95 text-white text-xs font-bold rounded-xl flex items-center gap-1.5 transition-all shadow-xs cursor-pointer"
          >
            <span>В Авторскую студию</span>
            <ArrowRight size={14} />
          </button>
        </div>
      )}

      <div>
        <h1 className="text-3xl font-extrabold text-gray-900 dark:text-white">Мое обучение</h1>
        <p className="text-gray-500 dark:text-slate-400 mt-2">Продолжайте изучать начатые курсы</p>
      </div>

      {courses.length === 0 ? (
        <div className="flex flex-col items-center justify-center p-12 bg-white dark:bg-slate-900 rounded-3xl border border-dashed border-gray-200 dark:border-slate-800">
          <BookOpen size={48} className="text-gray-300 dark:text-slate-700 mb-4" />
          <h2 className="text-xl font-bold text-gray-700 dark:text-slate-300">У вас пока нет курсов</h2>
          <p className="text-gray-500 dark:text-slate-500 mt-2 mb-6">Перейдите в каталог, чтобы найти что-то интересное.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {courses.map((course: any) => (
            <CourseCard 
              key={course.id} 
              course={course} 
              href={`/dashboard/courses/${course.id}`} // Направляем внутрь дашборда курса
            />
          ))}
        </div>
      )}
    </div>
  );
}
