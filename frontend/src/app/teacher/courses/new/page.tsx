'use client';

import React, { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { slugify, isValidSlug } from '@/lib/utils';
import { Category, fetchCategories, DEFAULT_CATEGORIES } from '@/lib/categories';
import TopNavbar from '@/components/layout/TopNavbar';
import {
  PlusCircle,
  Sparkles,
  Save,
  ArrowRight,
  Upload,
  Image as ImageIcon,
  Loader2,
  Trash2,
  RefreshCw,
  CheckCircle2,
  AlertCircle,
  Tag
} from 'lucide-react';

export default function CreateCoursePage() {
  const router = useRouter();

  const [title, setTitle] = useState('');
  const [slug, setSlug] = useState('');
  const [categories, setCategories] = useState<Category[]>(DEFAULT_CATEGORIES);
  const [category, setCategory] = useState('1');
  const [shortDesc, setShortDesc] = useState('');
  const [desc, setDesc] = useState('');
  const [coverUrl, setCoverUrl] = useState('');
  const [videoUrl, setVideoUrl] = useState('');
  const [difficulty, setDifficulty] = useState('beginner');
  const [language, setLanguage] = useState('RU');
  const [visibility, setVisibility] = useState('public');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSlugManuallyEdited, setIsSlugManuallyEdited] = useState(false);
  const [isUploadingCover, setIsUploadingCover] = useState(false);
  const [coverError, setCoverError] = useState<string | null>(null);
  const coverFileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    let isMounted = true;
    fetchCategories().then((list) => {
      if (isMounted && list.length > 0) {
        setCategories(list);
        if (!list.some((c) => String(c.id) === category)) {
          setCategory(String(list[0].id));
        }
      }
    });
    return () => {
      isMounted = false;
    };
  }, []);

  const resolveUrl = (url: string) => {
    if (!url) return '';
    if (url.startsWith('http://') || url.startsWith('https://')) return url;
    const baseUrl = (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8082/api/v1').replace(/\/api\/v1\/?$/, '');
    return `${baseUrl}${url.startsWith('/') ? '' : '/'}${url}`;
  };

  const handleCoverUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (file.size > 10 * 1024 * 1024) {
      setCoverError('Размер обложки не должен превышать 10 МБ');
      return;
    }

    setIsUploadingCover(true);
    setCoverError(null);

    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('category', 'course_cover');

      const res = await api.post('/upload', formData, {
        params: { category: 'course_cover' },
        headers: { 'Content-Type': 'multipart/form-data' },
      });

      const data = res.data.data || res.data;
      const uploadedUrl = data.file_url || data.url || data.FileUrl;
      if (uploadedUrl) {
        setCoverUrl(resolveUrl(uploadedUrl));
      }
    } catch (err: any) {
      console.error('Course cover upload failed', err);
      setCoverError(err.response?.data?.message || err.response?.data?.error || 'Не удалось загрузить обложку');
    } finally {
      setIsUploadingCover(false);
      if (coverFileInputRef.current) {
        coverFileInputRef.current.value = '';
      }
    }
  };

  const handleTitleChange = (val: string) => {
    setTitle(val);
    if (!isSlugManuallyEdited) {
      const generated = slugify(val);
      setSlug(generated);
    }
  };

  const handleSlugChange = (val: string) => {
    setIsSlugManuallyEdited(true);
    setSlug(val.toLowerCase().replace(/[^a-z0-9-]/g, ''));
  };

  const handleRegenerateSlug = () => {
    const generated = slugify(title);
    setSlug(generated);
    setIsSlugManuallyEdited(false);
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
        category_id: Number(category),
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

  const selectedCategoryObj = categories.find((c) => String(c.id) === String(category));

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
              <div className="flex items-center justify-between">
                <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                  URL Slug (для адресной строки) *
                </label>
                {isSlugManuallyEdited && (
                  <button
                    type="button"
                    onClick={handleRegenerateSlug}
                    className="text-[11px] font-bold text-indigo-600 hover:text-indigo-700 dark:text-indigo-400 flex items-center gap-1 transition-colors cursor-pointer"
                    title="Сгенерировать слаг заново из названия курса"
                  >
                    <RefreshCw size={11} />
                    <span>Синхронизировать</span>
                  </button>
                )}
              </div>
              <div className="relative">
                <input
                  type="text"
                  required
                  value={slug}
                  onChange={(e) => handleSlugChange(e.target.value)}
                  placeholder="arhitektura-mikroservisov-i-go"
                  className={`w-full p-3 pr-8 bg-slate-50 dark:bg-slate-800 border rounded-xl text-xs font-medium focus:ring-2 focus:outline-none transition-all ${
                    slug && !isValidSlug(slug)
                      ? 'border-rose-400 focus:ring-rose-400 text-rose-700 dark:text-rose-300'
                      : slug && isValidSlug(slug)
                      ? 'border-emerald-400/80 focus:ring-emerald-400'
                      : 'border-slate-200 dark:border-slate-700 focus:ring-indigo-500'
                  }`}
                />
                {slug && (
                  <div className="absolute right-3 top-1/2 -translate-y-1/2 flex items-center pointer-events-none">
                    {isValidSlug(slug) ? (
                      <span className="text-emerald-500 text-[10px] font-bold flex items-center gap-0.5">
                        <CheckCircle2 size={14} />
                      </span>
                    ) : (
                      <span className="text-rose-500 text-[10px] font-bold flex items-center gap-0.5">
                        <AlertCircle size={14} />
                      </span>
                    )}
                  </div>
                )}
              </div>
              <p className="text-[11px] text-slate-400">
                {slug && !isValidSlug(slug) ? (
                  <span className="text-rose-500 font-semibold">
                    Допустимы только строчные латинские буквы, цифры и одиночные дефисы
                  </span>
                ) : (
                  <span>
                    Ссылка на курс:{' '}
                    <code className="text-indigo-600 dark:text-indigo-400 font-mono">
                      /courses/{slug || '...'}
                    </code>
                  </span>
                )}
              </p>
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
              <div className="flex items-center justify-between">
                <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                  Обложка курса (Cover Image)
                </label>
                <button
                  type="button"
                  onClick={() => coverFileInputRef.current?.click()}
                  disabled={isUploadingCover}
                  className="text-xs font-bold text-indigo-600 hover:text-indigo-700 dark:text-indigo-400 flex items-center gap-1.5"
                >
                  {isUploadingCover ? (
                    <>
                      <Loader2 size={13} className="animate-spin" />
                      <span>Загрузка...</span>
                    </>
                  ) : (
                    <>
                      <Upload size={13} />
                      <span>Загрузить с диска</span>
                    </>
                  )}
                </button>
              </div>

              <input
                type="file"
                ref={coverFileInputRef}
                onChange={handleCoverUpload}
                accept="image/png,image/jpeg,image/webp"
                className="hidden"
              />

              <div className="flex gap-2">
                <input
                  type="url"
                  value={coverUrl}
                  onChange={(e) => setCoverUrl(e.target.value)}
                  placeholder="https://example.com/cover.jpg или выберите файл"
                  className="flex-1 p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                />
                {coverUrl && (
                  <button
                    type="button"
                    onClick={() => setCoverUrl('')}
                    className="px-3 py-2 rounded-xl border border-rose-200 dark:border-rose-900/50 text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors flex items-center gap-1 text-xs font-bold"
                    title="Удалить обложку"
                  >
                    <Trash2 size={14} />
                  </button>
                )}
              </div>

              {coverError && (
                <p className="text-[11px] text-rose-600 dark:text-rose-400 font-semibold">{coverError}</p>
              )}

              {coverUrl && (
                <div className="relative mt-2 w-full h-32 rounded-xl overflow-hidden border border-slate-200 dark:border-slate-700 bg-slate-100 dark:bg-slate-800">
                  <img
                    src={coverUrl}
                    alt="Предпросмотр обложки"
                    className="w-full h-full object-cover"
                    onError={(e) => {
                      (e.target as HTMLElement).style.display = 'none';
                    }}
                  />
                  <span className="absolute bottom-1.5 right-1.5 bg-black/60 text-white text-[10px] font-bold px-2 py-0.5 rounded-md backdrop-blur-xs">
                    Превью
                  </span>
                </div>
              )}
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

          {/* Selects: Category, Difficulty, Language, Visibility */}
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300 flex items-center gap-1.5">
                <Tag size={13} className="text-indigo-600 dark:text-indigo-400" />
                <span>Категория курса *</span>
              </label>
              <select
                value={category}
                onChange={(e) => setCategory(e.target.value)}
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold focus:ring-2 focus:ring-indigo-500 focus:outline-none cursor-pointer"
              >
                {categories.map((cat) => (
                  <option key={cat.id} value={cat.id}>
                    {cat.name}
                  </option>
                ))}
              </select>
            </div>

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

          {selectedCategoryObj?.description && (
            <p className="text-[11px] text-slate-500 dark:text-slate-400 italic">
              Выбрана категория «{selectedCategoryObj.name}»: {selectedCategoryObj.description}
            </p>
          )}

          {/* Submit Action */}
          <div className="flex justify-end gap-3 pt-6 border-t border-slate-100 dark:border-slate-800">
            <button
              type="button"
              onClick={() => router.back()}
              className="px-5 py-3 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-600 dark:text-slate-400 hover:bg-slate-50 cursor-pointer"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="px-8 py-3 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-xl shadow-lg shadow-indigo-600/20 transition-all flex items-center gap-2 cursor-pointer disabled:opacity-50"
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
