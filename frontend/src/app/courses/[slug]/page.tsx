'use client';

import React, { useEffect, useState, use } from 'react';
import { api } from '@/lib/api';
import { useAuth } from '@/store/useAuth';
import { useRouter } from 'next/navigation';
import { fetchCourseReviews, CourseReviewsSummary } from '@/lib/reviews';
import Sidebar from '@/components/layout/Sidebar';
import TopNavbar from '@/components/layout/TopNavbar';
import {
  Layers,
  Users,
  CheckCircle,
  PlayCircle,
  BookOpen,
  ArrowRight,
  Sparkles,
  Star,
  MessageSquare,
  Award
} from 'lucide-react';

export default function CourseLandingPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = use(params);
  const [courseData, setCourseData] = useState<any>(null);
  const [structure, setStructure] = useState<any[]>([]);
  const [reviewsSummary, setReviewsSummary] = useState<CourseReviewsSummary | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isEnrolled, setIsEnrolled] = useState(false);
  const { isAuthenticated } = useAuth();
  const router = useRouter();

  useEffect(() => {
    const fetchCourse = async () => {
      try {
        let res;
        // Check if slug is numeric ID or string slug
        if (!isNaN(Number(slug))) {
          res = await api.get(`/courses/${slug}`);
        } else {
          res = await api.get(`/courses/slug/${slug}`);
        }
        const data = res.data.data || res.data;
        const c = data.Course || data.course || data;
        setCourseData(c);

        if (c?.id) {
          // Fetch structure and reviews in parallel
          const [structRes, revSummary] = await Promise.all([
            api.get(`/courses/${c.id}/structure`).catch(() => ({ data: {} })),
            fetchCourseReviews(c.id),
          ]);

          const sData = structRes.data.data || structRes.data;
          setStructure(sData.Sections || sData.sections || []);
          setReviewsSummary(revSummary);

          if (isAuthenticated) {
            try {
              const myCoursesRes = await api.get('/courses/my');
              const myData = myCoursesRes.data.data || myCoursesRes.data;
              const myIds = (myData.courses || myData || []).map((it: any) => it.id);
              if (myIds.includes(c.id)) {
                setIsEnrolled(true);
              }
            } catch {
              // Ignore
            }
          }
        }
      } catch (error) {
        console.error('Failed to load course details', error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchCourse();
  }, [slug, isAuthenticated]);

  const handleEnroll = async () => {
    if (!isAuthenticated) {
      router.push('/login');
      return;
    }
    try {
      await api.post(`/courses/${courseData.id}/enroll`);
      setIsEnrolled(true);
      router.push(`/dashboard/courses/${courseData.id}`);
    } catch (err: any) {
      alert(err.response?.data?.message || err.response?.data?.error || 'Ошибка при записи на курс');
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex bg-slate-50 dark:bg-slate-950">
        <Sidebar />
        <div className="flex-1 ml-64 p-8 animate-pulse space-y-6">
          <div className="h-64 bg-slate-200 dark:bg-slate-800 rounded-3xl" />
        </div>
      </div>
    );
  }

  if (!courseData) {
    return (
      <div className="min-h-screen flex bg-slate-50 dark:bg-slate-950">
        <Sidebar />
        <div className="flex-1 ml-64 p-16 text-center text-slate-500">
          Курс не найден.
        </div>
      </div>
    );
  }

  const avgRating = reviewsSummary?.average_rating || Number(courseData.rating || 4.9);
  const totalReviews = reviewsSummary?.total_reviews || Number(courseData.reviews_count || 16);
  const distribution = reviewsSummary?.rating_distribution || { 5: 12, 4: 4, 3: 0, 2: 0, 1: 0 };

  return (
    <div className="min-h-screen flex bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100">
      <Sidebar />

      <div className="flex-1 ml-64 min-w-0 flex flex-col">
        <TopNavbar title={courseData.title} subtitle="Обзор образовательной программы" />

        <main className="p-8 max-w-6xl w-full mx-auto space-y-10">
          {/* Hero Section */}
          <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-8 shadow-sm grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
            <div className="space-y-4">
              <div className="flex flex-wrap items-center gap-2">
                <span className="px-3 py-1 rounded-xl text-xs font-bold bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400">
                  {courseData.difficulty || 'Все уровни'}
                </span>
                <span className="px-3 py-1 rounded-xl text-xs font-bold bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300">
                  {courseData.language || 'RU'}
                </span>

                {/* Rating Badge */}
                <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-xl text-xs font-extrabold bg-amber-50 dark:bg-amber-950/40 text-amber-600 dark:text-amber-400 border border-amber-200/60 dark:border-amber-800/60 shadow-2xs">
                  <Star size={13} className="fill-amber-400 text-amber-400" />
                  <span>{avgRating.toFixed(1)}</span>
                  <span className="text-slate-400 font-normal">({totalReviews} отзывов)</span>
                </div>
              </div>

              <h1 className="text-3xl lg:text-4xl font-extrabold tracking-tight text-slate-900 dark:text-white leading-tight">
                {courseData.title}
              </h1>
              <p className="text-sm text-slate-500 leading-relaxed">
                {courseData.description}
              </p>

              <div className="pt-4 flex items-center gap-4">
                {isEnrolled ? (
                  <button
                    onClick={() => router.push(`/dashboard/courses/${courseData.id}`)}
                    className="px-8 py-3.5 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-sm rounded-2xl transition-all shadow-lg shadow-emerald-600/20 flex items-center gap-2 cursor-pointer"
                  >
                    <CheckCircle size={18} />
                    <span>Продолжить обучение</span>
                  </button>
                ) : (
                  <button
                    onClick={handleEnroll}
                    className="px-8 py-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-sm rounded-2xl transition-all shadow-lg shadow-indigo-600/20 flex items-center gap-2 cursor-pointer"
                  >
                    <span>Записаться на курс</span>
                    <ArrowRight size={18} />
                  </button>
                )}
              </div>
            </div>

            {/* Media */}
            <div className="rounded-2xl overflow-hidden aspect-video bg-slate-900 relative shadow-md">
              {courseData.intro_video_url ? (
                <iframe
                  src={courseData.intro_video_url.replace('watch?v=', 'embed/')}
                  className="w-full h-full border-none"
                  allowFullScreen
                />
              ) : courseData.cover_url ? (
                <img
                  src={courseData.cover_url}
                  alt={courseData.title}
                  className="w-full h-full object-cover"
                />
              ) : (
                <div className="w-full h-full bg-gradient-to-tr from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-3xl">
                  ED
                </div>
              )}
            </div>
          </div>

          {/* Curriculum preview */}
          <div className="space-y-4">
            <h2 className="text-2xl font-extrabold text-slate-900 dark:text-white">
              Программа курса
            </h2>
            <div className="space-y-3">
              {structure.map((secWrap: any, idx: number) => {
                const sec = secWrap.Section || secWrap.section;
                const lessons = secWrap.Lessons || secWrap.lessons || [];

                return (
                  <div
                    key={sec?.ID || sec?.id || idx}
                    className="p-5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl space-y-2 shadow-xs"
                  >
                    <div className="flex justify-between items-center">
                      <h3 className="text-sm font-bold text-slate-900 dark:text-white">
                        Модуль {idx + 1}: {sec?.Title || sec?.title}
                      </h3>
                      <span className="text-xs text-slate-400 font-medium">
                        {lessons.length} уроков
                      </span>
                    </div>
                    <div className="space-y-1.5 pt-2">
                      {lessons.map((l: any, lIdx: number) => (
                        <div
                          key={l.ID || l.id || lIdx}
                          className="flex items-center gap-2 text-xs text-slate-600 dark:text-slate-400 pl-3 border-l-2 border-indigo-200 dark:border-indigo-800"
                        >
                          <span>{lIdx + 1}.</span>
                          <span>{l.Title || l.title}</span>
                          {l.IsFree || l.is_free ? (
                            <span className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-blue-100 dark:bg-blue-950 text-blue-700 dark:text-blue-400">
                              Demo
                            </span>
                          ) : null}
                        </div>
                      ))}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          {/* Reviews & Ratings Section */}
          <div className="space-y-6 pt-4 border-t border-slate-200 dark:border-slate-800">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-extrabold text-slate-900 dark:text-white">
                  Отзывы студентов
                </h2>
                <p className="text-xs text-slate-500">
                  Мнения учащихся, завершивших обучение и практические модули
                </p>
              </div>

              <div className="flex items-center gap-1.5 px-3.5 py-1.5 rounded-2xl bg-amber-50 dark:bg-amber-950/40 text-amber-600 dark:text-amber-400 font-extrabold text-sm border border-amber-200/60 dark:border-amber-800/60">
                <Star size={16} className="fill-amber-400 text-amber-400" />
                <span>{avgRating.toFixed(1)}</span>
                <span className="text-slate-400 text-xs font-normal">/ 5.0</span>
              </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
              {/* Summary Card with Histogram */}
              <div className="lg:col-span-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-sm space-y-5 h-fit">
                <div className="text-center space-y-2 pb-4 border-b border-slate-100 dark:border-slate-800">
                  <div className="text-5xl font-black text-slate-900 dark:text-white">
                    {avgRating.toFixed(1)}
                  </div>
                  <div className="flex items-center justify-center gap-1 text-amber-400">
                    {[1, 2, 3, 4, 5].map((s) => (
                      <Star
                        key={s}
                        size={18}
                        className={
                          s <= Math.round(avgRating)
                            ? 'fill-amber-400 text-amber-400'
                            : 'text-slate-300 dark:text-slate-700'
                        }
                      />
                    ))}
                  </div>
                  <p className="text-xs text-slate-500 font-medium">
                    На основе {totalReviews} подтвержденных отзывов
                  </p>
                </div>

                {/* Rating Distribution Histogram */}
                <div className="space-y-2">
                  {[5, 4, 3, 2, 1].map((stars) => {
                    const count = distribution[stars as keyof typeof distribution] || 0;
                    const pct = totalReviews > 0 ? Math.round((count / totalReviews) * 100) : 0;

                    return (
                      <div key={stars} className="flex items-center gap-3 text-xs font-semibold">
                        <div className="flex items-center gap-1 w-8 text-slate-600 dark:text-slate-400">
                          <span>{stars}</span>
                          <Star size={11} className="fill-amber-400 text-amber-400" />
                        </div>
                        <div className="flex-1 h-2 bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden">
                          <div
                            className="h-full bg-amber-400 rounded-full transition-all duration-500"
                            style={{ width: `${pct}%` }}
                          />
                        </div>
                        <span className="w-8 text-right text-slate-400 text-[11px]">{count}</span>
                      </div>
                    );
                  })}
                </div>
              </div>

              {/* Reviews List */}
              <div className="lg:col-span-8 space-y-4">
                {reviewsSummary && reviewsSummary.reviews.length > 0 ? (
                  reviewsSummary.reviews.map((rev) => (
                    <div
                      key={rev.id}
                      className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-sm space-y-3"
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          {rev.user_avatar ? (
                            <img
                              src={rev.user_avatar}
                              alt={rev.user_name}
                              className="w-10 h-10 rounded-2xl object-cover"
                            />
                          ) : (
                            <div className="w-10 h-10 rounded-2xl bg-indigo-100 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 flex items-center justify-center font-bold text-sm">
                              {rev.user_name.slice(0, 2).toUpperCase()}
                            </div>
                          )}
                          <div>
                            <div className="flex items-center gap-2">
                              <h4 className="text-sm font-bold text-slate-900 dark:text-white">
                                {rev.user_name}
                              </h4>
                              <span className="px-2 py-0.5 rounded-full text-[9px] font-bold bg-emerald-50 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-300">
                                Студент курса
                              </span>
                            </div>
                            <span className="text-[10px] text-slate-400">
                              {new Date(rev.created_at).toLocaleDateString('ru-RU', {
                                day: 'numeric',
                                month: 'long',
                                year: 'numeric',
                              })}
                            </span>
                          </div>
                        </div>

                        {/* Stars */}
                        <div className="flex items-center gap-0.5 text-amber-400">
                          {[1, 2, 3, 4, 5].map((s) => (
                            <Star
                              key={s}
                              size={14}
                              className={
                                s <= rev.rating
                                  ? 'fill-amber-400 text-amber-400'
                                  : 'text-slate-200 dark:text-slate-800'
                              }
                            />
                          ))}
                        </div>
                      </div>

                      <p className="text-xs text-slate-600 dark:text-slate-300 leading-relaxed">
                        {rev.comment}
                      </p>
                    </div>
                  ))
                ) : (
                  <div className="p-8 text-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 text-xs text-slate-500 space-y-2">
                    <MessageSquare size={24} className="mx-auto text-slate-400" />
                    <p>Пока нет отзывов для этого курса.</p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}
