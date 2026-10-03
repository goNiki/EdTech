'use client';

import React, { useEffect, useState, useTransition } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { useAuth } from '@/store/useAuth';
import TopNavbar from '@/components/layout/TopNavbar';
import Sidebar from '@/components/layout/Sidebar';
import {
  Search,
  SlidersHorizontal,
  ArrowUpDown,
  BookOpen,
  Users,
  PlayCircle,
  Check,
  X,
  Sparkles,
  Layers
} from 'lucide-react';

interface Course {
  id: number;
  title: string;
  slug: string;
  short_description?: string;
  description: string;
  cover_url?: string;
  intro_video_url?: string;
  status: string;
  visibility: string;
  difficulty?: string;
  language?: string;
  total_lessons: number;
  total_sections: number;
  enrolled_count: number;
  created_at: string;
}

export default function CoursesCatalogPage() {
  const router = useRouter();
  const { user, isAuthenticated } = useAuth();

  const [courses, setCourses] = useState<Course[]>([]);
  const [enrolledCourseIds, setEnrolledCourseIds] = useState<number[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Filters State (matches CourseFilter)
  const [search, setSearch] = useState('');
  const [difficulty, setDifficulty] = useState('');
  const [language, setLanguage] = useState('');
  const [sortBy, setSortBy] = useState('created_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  // Video Preview Modal
  const [previewVideoUrl, setPreviewVideoUrl] = useState<string | null>(null);
  const [previewCourseTitle, setPreviewCourseTitle] = useState<string>('');

  // Toast notification state
  const [toastMessage, setToastMessage] = useState<string | null>(null);

  const showToast = (msg: string) => {
    setToastMessage(msg);
    setTimeout(() => setToastMessage(null), 3000);
  };

  const fetchCourses = async () => {
    setIsLoading(true);
    try {
      const params = new URLSearchParams();
      if (search.trim()) params.append('search', search.trim());
      if (difficulty) params.append('difficulty', difficulty);
      if (language) params.append('language', language);
      if (sortBy) params.append('sort_by', sortBy);
      if (sortOrder) params.append('sort_order', sortOrder);

      const res = await api.get(`/courses?${params.toString()}`);
      const data = res.data.data || res.data;
      setCourses(data.courses || data || []);

      if (isAuthenticated) {
        try {
          const myRes = await api.get('/courses/my');
          const myData = myRes.data.data || myRes.data;
          const enrolled = (myData.courses || myData || []).map((c: any) => c.id);
          setEnrolledCourseIds(enrolled);
        } catch {
          // Guest or error
        }
      }
    } catch (err) {
      console.error('Failed to load courses:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      fetchCourses();
    }, 250);
    return () => clearTimeout(timer);
  }, [search, difficulty, language, sortBy, sortOrder, isAuthenticated]);

  const handleEnroll = async (courseId: number, e: React.MouseEvent) => {
    e.stopPropagation();
    if (!isAuthenticated) {
      router.push('/login');
      return;
    }

    try {
      await api.post(`/courses/${courseId}/enroll`);
      setEnrolledCourseIds((prev) => [...prev, courseId]);
      showToast('Вы успешно записались на курс!');
    } catch (err: any) {
      const msg = err.response?.data?.message || err.response?.data?.error || 'Ошибка при записи на курс';
      showToast(msg);
    }
  };

  const toggleSortOrder = () => {
    setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
  };

  const getDifficultyLabel = (diff?: string) => {
    switch (diff) {
      case 'beginner':
        return 'Начальный';
      case 'intermediate':
        return 'Средний';
      case 'advanced':
        return 'Продвинутый';
      default:
        return 'Все уровни';
    }
  };

  return (
    <div className="min-h-screen flex bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100">
      <Sidebar />

      <div className="flex-1 ml-64 min-w-0 flex flex-col">
        <TopNavbar title="Все курсы" subtitle="Каталог образовательных программ и практических модулей" />

        <main className="p-8 max-w-7xl w-full mx-auto space-y-8">
          {/* Header Banner */}
          <div className="bg-gradient-to-r from-indigo-600 via-indigo-500 to-purple-600 rounded-3xl p-8 text-white shadow-xl shadow-indigo-500/10 flex flex-col md:flex-row justify-between md:items-center gap-6">
            <div className="space-y-2 max-w-2xl">
              <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/20 backdrop-blur-md text-xs font-bold uppercase tracking-wider">
                <Sparkles size={14} /> Открытые программы
              </div>
              <h2 className="text-3xl font-extrabold tracking-tight">Исследуйте мир новых знаний</h2>
              <p className="text-indigo-100 text-sm leading-relaxed">
                Выбирайте курсы от ведущих авторов, проходите интерактивные уроки в формате Content-as-Data и получайте мгновенную обратную связь.
              </p>
            </div>
            <div className="bg-white/10 backdrop-blur-md border border-white/20 rounded-2xl p-4 flex items-center gap-4 flex-shrink-0">
              <div className="w-12 h-12 rounded-xl bg-white text-indigo-600 flex items-center justify-center font-black text-xl shadow-sm">
                {courses.length}
              </div>
              <div>
                <p className="text-xs text-indigo-100 font-medium">Доступно курсов</p>
                <p className="text-sm font-extrabold text-white">В каталоге</p>
              </div>
            </div>
          </div>

          {/* Filter and Search Bar */}
          <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-5 shadow-sm space-y-4">
            <div className="flex flex-col md:flex-row gap-4 items-center justify-between">
              {/* Search Input */}
              <div className="relative w-full md:w-96">
                <Search size={18} className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
                <input
                  type="text"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  placeholder="Поиск курсов по названию или теме..."
                  className="w-full pl-11 pr-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all placeholder:text-slate-400"
                />
              </div>

              {/* Filters Group */}
              <div className="flex flex-wrap items-center gap-3 w-full md:w-auto">
                {/* Difficulty Select */}
                <select
                  value={difficulty}
                  onChange={(e) => setDifficulty(e.target.value)}
                  aria-label="Уровень сложности"
                  className="px-3.5 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl text-xs font-semibold text-slate-700 dark:text-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500 cursor-pointer"
                >
                  <option value="">Все уровни</option>
                  <option value="beginner">Начальный</option>
                  <option value="intermediate">Средний</option>
                  <option value="advanced">Продвинутый</option>
                </select>

                {/* Language Select */}
                <select
                  value={language}
                  onChange={(e) => setLanguage(e.target.value)}
                  aria-label="Язык курса"
                  className="px-3.5 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl text-xs font-semibold text-slate-700 dark:text-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500 cursor-pointer"
                >
                  <option value="">Все языки</option>
                  <option value="RU">Русский (RU)</option>
                  <option value="EN">English (EN)</option>
                </select>

                {/* Sort By Select */}
                <select
                  value={sortBy}
                  onChange={(e) => setSortBy(e.target.value)}
                  aria-label="Сортировка"
                  className="px-3.5 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl text-xs font-semibold text-slate-700 dark:text-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500 cursor-pointer"
                >
                  <option value="created_at">По дате добавления</option>
                  <option value="title">По названию</option>
                  <option value="enrolled_count">По популярности</option>
                </select>

                {/* Sort Order Toggle */}
                <button
                  onClick={toggleSortOrder}
                  className="p-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl text-slate-600 dark:text-slate-300 hover:text-indigo-600 hover:border-indigo-300 transition-all flex items-center gap-1.5 text-xs font-bold shadow-xs"
                  title={sortOrder === 'asc' ? 'По возрастанию (А-Я)' : 'По убыванию (Я-А)'}
                >
                  <ArrowUpDown size={15} />
                  <span className="uppercase">{sortOrder}</span>
                </button>
              </div>
            </div>
          </div>

          {/* Courses Grid */}
          {isLoading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {[1, 2, 3, 4, 5, 6].map((i) => (
                <div
                  key={i}
                  className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl h-96 p-5 animate-pulse space-y-4"
                >
                  <div className="bg-slate-200 dark:bg-slate-800 h-44 rounded-2xl w-full" />
                  <div className="space-y-2">
                    <div className="bg-slate-200 dark:bg-slate-800 h-5 rounded-md w-3/4" />
                    <div className="bg-slate-200 dark:bg-slate-800 h-3 rounded-md w-full" />
                    <div className="bg-slate-200 dark:bg-slate-800 h-3 rounded-md w-2/3" />
                  </div>
                </div>
              ))}
            </div>
          ) : courses.length === 0 ? (
            <div className="p-16 text-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 shadow-sm space-y-3">
              <div className="w-12 h-12 rounded-2xl bg-indigo-50 dark:bg-slate-800 text-indigo-600 flex items-center justify-center mx-auto">
                <BookOpen size={24} />
              </div>
              <h3 className="text-lg font-bold text-slate-800 dark:text-slate-200">Курсы не найдены</h3>
              <p className="text-xs text-slate-500 max-w-sm mx-auto">
                Попробуйте изменить параметры поиска или сбросить установленные фильтры.
              </p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {courses.map((course) => {
                const isEnrolled = enrolledCourseIds.includes(course.id);

                return (
                  <div
                    key={course.id}
                    onClick={() => router.push(`/courses/${course.slug || course.id}`)}
                    className="group bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl overflow-hidden shadow-sm hover:shadow-xl hover:border-indigo-300 dark:hover:border-indigo-800 transition-all duration-300 flex flex-col cursor-pointer"
                  >
                    {/* Course Card Media */}
                    <div className="h-48 bg-slate-100 dark:bg-slate-800 relative overflow-hidden">
                      {course.cover_url ? (
                        <img
                          src={course.cover_url}
                          alt={course.title}
                          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                        />
                      ) : (
                        <div className="w-full h-full bg-gradient-to-tr from-indigo-500 to-purple-600 flex items-center justify-center text-white font-extrabold text-2xl">
                          {course.title.slice(0, 2).toUpperCase()}
                        </div>
                      )}

                      <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent" />

                      {/* Top Badges */}
                      <div className="absolute top-3 left-3 right-3 flex justify-between items-center gap-2">
                        <span className="px-2.5 py-1 rounded-xl text-[10px] font-bold bg-white/90 dark:bg-slate-900/90 text-indigo-600 dark:text-indigo-400 backdrop-blur-md shadow-xs">
                          {getDifficultyLabel(course.difficulty)}
                        </span>
                        {course.language && (
                          <span className="px-2 py-0.5 rounded-lg text-[10px] font-extrabold bg-black/40 text-white backdrop-blur-md uppercase tracking-wider">
                            {course.language}
                          </span>
                        )}
                      </div>

                      {/* Intro Video Button */}
                      {course.intro_video_url && (
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setPreviewVideoUrl(course.intro_video_url || null);
                            setPreviewCourseTitle(course.title);
                          }}
                          className="absolute bottom-3 right-3 p-2 bg-white/90 dark:bg-slate-900/90 hover:bg-white text-indigo-600 rounded-xl backdrop-blur-md shadow-md hover:scale-105 transition-all flex items-center gap-1 text-[11px] font-bold"
                          title="Смотреть промо-ролик"
                        >
                          <PlayCircle size={16} />
                          <span>Интро</span>
                        </button>
                      )}

                      {/* Title overlay */}
                      <h3 className="absolute bottom-3 left-3 right-16 text-white font-extrabold text-base line-clamp-1 drop-shadow-md">
                        {course.title}
                      </h3>
                    </div>

                    {/* Card Content */}
                    <div className="p-5 flex-1 flex flex-col justify-between space-y-4">
                      <p className="text-xs text-slate-500 line-clamp-2 leading-relaxed">
                        {course.short_description || course.description || 'Описание курса появится позже.'}
                      </p>

                      <div className="pt-3 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between text-xs text-slate-500 font-medium">
                        <div className="flex items-center gap-1.5">
                          <Layers size={14} className="text-indigo-500" />
                          <span>{course.total_lessons || 0} уроков</span>
                        </div>
                        <div className="flex items-center gap-1.5">
                          <Users size={14} className="text-purple-500" />
                          <span>{course.enrolled_count || 0} студентов</span>
                        </div>
                      </div>

                      {/* Action Button */}
                      <div className="pt-2">
                        {isEnrolled ? (
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              router.push(`/dashboard/courses/${course.id}`);
                            }}
                            className="w-full py-2.5 bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 hover:bg-indigo-100 text-xs font-bold rounded-2xl transition-all flex items-center justify-center gap-1.5 ring-1 ring-indigo-200 dark:ring-indigo-800"
                          >
                            <Check size={16} />
                            <span>Продолжить обучение</span>
                          </button>
                        ) : (
                          <button
                            onClick={(e) => handleEnroll(course.id, e)}
                            className="w-full py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-2xl transition-all shadow-md shadow-indigo-600/20"
                          >
                            Присоединиться к курсу
                          </button>
                        )}
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </main>
      </div>

      {/* Video Preview Modal */}
      {previewVideoUrl && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-xs p-4">
          <div className="bg-white dark:bg-slate-900 rounded-3xl overflow-hidden max-w-2xl w-full border border-slate-200 dark:border-slate-800 shadow-2xl space-y-4">
            <div className="p-4 border-b border-slate-100 dark:border-slate-800 flex justify-between items-center">
              <h4 className="text-sm font-bold text-slate-900 dark:text-white truncate">
                Интро: {previewCourseTitle}
              </h4>
              <button
                onClick={() => setPreviewVideoUrl(null)}
                className="p-1 rounded-lg text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors"
              >
                <X size={20} />
              </button>
            </div>
            <div className="aspect-video w-full bg-black">
              <iframe
                src={previewVideoUrl.replace('watch?v=', 'embed/')}
                className="w-full h-full border-none"
                allowFullScreen
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
              />
            </div>
          </div>
        </div>
      )}

      {/* Toast */}
      {toastMessage && (
        <div className="fixed bottom-6 right-6 z-50 bg-slate-900 text-white text-xs font-bold px-4 py-3 rounded-2xl shadow-xl flex items-center gap-2 animate-bounce">
          <Sparkles size={16} className="text-indigo-400" />
          <span>{toastMessage}</span>
        </div>
      )}
    </div>
  );
}
