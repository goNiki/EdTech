'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import TopNavbar from '@/components/layout/TopNavbar';
import { PlusCircle, Sparkles, Save, ArrowRight } from 'lucide-react';

export default function CreateCoursePage() {
  const router = useRouter();

  const [title, setTitle] = useState('');
  const [slug, setSlug] = useState('');
  const [category, setCategory] = useState('1');
  const [shortDesc, setShortDesc] = useState('');
  const [desc, setDesc] = useState('');
  const [coverUrl, setCoverUrl] = useState('');
  const [videoUrl, setVideoUrl] = useState('');
  const [difficulty, setDifficulty] = useState('beginner');
  const [language, setLanguage] = useState('RU');
  const [visibility, setVisibility] = useState('public');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleTitleChange = (val: string) => {
    setTitle(val);
    // Auto generate translit slug
    const generated = val
      .toLowerCase()
      .trim()
      .replace(/[^\w\s-]/g, '')
      .replace(/[\s_-]+/g, '-')
      .replace(/^-+|-+$/g, '');
    setSlug(generated || `course-${Date.now()}`);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) {
      alert('Пожалуйста, укажите название курса');
      return;
    }

    setIsSubmitting(true);
    try {
      const payload = {
        title: title.trim(),
        slug: slug.trim() || `course-${Date.now()}`,
        short_description: shortDesc.trim() || undefined,
        description: desc.trim() || title.trim(),
        cover_url: coverUrl.trim() || 'https://placehold.co/600x340/6366f1/ffffff?text=Course+Cover',
        intro_video_url: videoUrl.trim() || undefined,
        difficulty: difficulty,
        language: language,
        visibility: visibility,
      };

      const res = await api.post('/courses', payload);
      const data = res.data.data || res.data;
      const newId = data.id || data.Id || data.course?.id;

      if (newId) {
        router.push(`/teacher/courses/${newId}/curriculum`);
      } else {
        router.push('/teacher/courses');
      }
    } catch (err: any) {
      console.error('Failed to create course', err);
      alert(err.response?.data?.message || err.response?.data?.error || 'Ошибка при создании курса');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex flex-col min-h-screen">
      <TopNavbar title="Создать новый курс" subtitle="Конструктор параметров образовательной программы" />

      <main className="p-8 max-w-4xl w-full mx-auto space-y-8 flex-1">
        <form onSubmit={handleSubmit} className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-8 shadow-sm space-y-6">
          <div className="space-y-1 pb-4 border-b border-slate-100 dark:border-slate-800">
            <h2 className="text-xl font-extrabold text-slate-900 dark:text-white">
              Основные параметры курса
            </h2>
            <p className="text-xs text-slate-500">
              Курс будет создан со статусом «Черновик». Вы сможете опубликовать его после добавления модулей и уроков.
            </p>
          </div>

          {/* Title & Slug */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Название курса *
              </label>
              <input
                type="text"
                required
                value={title}
                onChange={(e) => handleTitleChange(e.target.value)}
                placeholder="Например: Полный курс по Golang и микросервисам"
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
              />
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                URL Slug (для адресной строки) *
              </label>
              <input
                type="text"
                required
                value={slug}
                onChange={(e) => setSlug(e.target.value)}
                placeholder="golang-microservices"
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
              />
            </div>
          </div>

          {/* Short Description */}
          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
              Краткое описание (для карточки курса)
            </label>
            <input
              type="text"
              value={shortDesc}
              onChange={(e) => setShortDesc(e.target.value)}
              placeholder="Кратко опишите, о чем данный курс и для кого он предназначен..."
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
          </div>

          {/* Full Description */}
          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
              Полное развернутое описание курса
            </label>
            <textarea
              rows={4}
              value={desc}
              onChange={(e) => setDesc(e.target.value)}
              placeholder="Подробная программа, требования к ученикам и результаты обучения..."
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
          </div>

          {/* Media Links */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Ссылка на обложку (Cover Image URL)
              </label>
              <input
                type="url"
                value={coverUrl}
                onChange={(e) => setCoverUrl(e.target.value)}
                placeholder="https://example.com/cover.jpg"
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
              />
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Ссылка на интро-видео (YouTube URL)
              </label>
              <input
                type="url"
                value={videoUrl}
                onChange={(e) => setVideoUrl(e.target.value)}
                placeholder="https://www.youtube.com/watch?v=..."
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
              />
            </div>
          </div>

          {/* Selects: Difficulty, Language, Visibility */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Уровень сложности
              </label>
              <select
                value={difficulty}
                onChange={(e) => setDifficulty(e.target.value)}
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold focus:ring-2 focus:ring-indigo-500 focus:outline-none cursor-pointer"
              >
                <option value="beginner">Начальный (Beginner)</option>
                <option value="intermediate">Средний (Intermediate)</option>
                <option value="advanced">Продвинутый (Advanced)</option>
              </select>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Язык курса
              </label>
              <select
                value={language}
                onChange={(e) => setLanguage(e.target.value)}
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold focus:ring-2 focus:ring-indigo-500 focus:outline-none cursor-pointer"
              >
                <option value="RU">Русский (RU)</option>
                <option value="EN">English (EN)</option>
              </select>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Видимость в каталоге
              </label>
              <select
                value={visibility}
                onChange={(e) => setVisibility(e.target.value)}
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold focus:ring-2 focus:ring-indigo-500 focus:outline-none cursor-pointer"
              >
                <option value="public">Публичный (Виден всем)</option>
                <option value="private">Приватный (По ссылке)</option>
              </select>
            </div>
          </div>

          {/* Submit Action */}
          <div className="flex justify-end gap-3 pt-6 border-t border-slate-100 dark:border-slate-800">
            <button
              type="button"
              onClick={() => router.back()}
              className="px-5 py-3 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-600 dark:text-slate-400 hover:bg-slate-50"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="px-8 py-3 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-xl shadow-lg shadow-indigo-600/20 transition-all flex items-center gap-2"
            >
              <Sparkles size={16} />
              <span>{isSubmitting ? 'Создание курса...' : 'Создать и перейти к программе'}</span>
              <ArrowRight size={16} />
            </button>
          </div>
        </form>
      </main>
    </div>
  );
}
