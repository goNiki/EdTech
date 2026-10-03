'use client';

import React, { useEffect, useState, use, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { slugify, isValidSlug } from '@/lib/utils';
import TopNavbar from '@/components/layout/TopNavbar';
import { Save, ArrowLeft, Sparkles, CheckCircle2, Sliders, ExternalLink, Upload, Loader2, Trash2, RefreshCw, AlertCircle } from 'lucide-react';

export default function CourseSettingsPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();

  const [formData, setFormData] = useState({
    title: '',
    slug: '',
    short_description: '',
    description: '',
    cover_url: '',
    intro_video_url: '',
    difficulty: 'beginner',
    language: 'RU',
    visibility: 'public',
    status: 'draft',
    category_id: 1,
  });

  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [toastMsg, setToastMsg] = useState<string | null>(null);
  const [isUploadingCover, setIsUploadingCover] = useState(false);
  const [coverError, setCoverError] = useState<string | null>(null);
  const coverFileInputRef = useRef<HTMLInputElement>(null);

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
        setFormData((prev) => ({ ...prev, cover_url: resolveUrl(uploadedUrl) }));
        showToast('Обложка успешно загружена');
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

  const handleGenerateSlugFromTitle = () => {
    if (!formData.title.trim()) {
      alert('Сначала укажите название курса');
      return;
    }
    const generated = slugify(formData.title);
    setFormData((prev) => ({ ...prev, slug: generated }));
    showToast('Слаг успешно обновлен из названия');
  };

  const handleSlugChange = (val: string) => {
    const clean = val.toLowerCase().replace(/[^a-z0-9-]/g, '');
    setFormData((prev) => ({ ...prev, slug: clean }));
  };

  const showToast = (msg: string) => {
    setToastMsg(msg);
    setTimeout(() => setToastMsg(null), 3000);
  };

  useEffect(() => {
    const fetchCourse = async () => {
      try {
        const { data } = await api.get(`/courses/${id}`);
        const c = data.data?.course || data.data?.Course || data.course || data.Course || data.data || data;

        if (c) {
          setFormData({
            title: c.title || c.Title || '',
            slug: c.slug || c.Slug || '',
            short_description: c.short_description || c.ShortDescription || '',
            description: c.description || c.Description || '',
            cover_url: c.cover_url || c.CoverUrl || '',
            intro_video_url: c.intro_video_url || c.IntroVideoUrl || '',
            difficulty: (c.difficulty || c.Difficulty || 'beginner').toLowerCase(),
            language: c.language || c.Language || 'RU',
            visibility: (c.visibility || c.Visibility || 'public').toLowerCase(),
            status: (c.status || c.Status || 'draft').toLowerCase(),
            category_id: Number(c.category_id || c.CategoryID || 1),
          });
        }
      } catch (err) {
        console.error('Failed to fetch course settings', err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchCourse();
  }, [id]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.title.trim()) {
      alert('Пожалуйста, укажите название курса');
      return;
    }

    setIsSaving(true);
    try {
      // 1. Update general course parameters
      await api.patch(`/courses/${id}`, {
        title: formData.title.trim(),
        slug: formData.slug.trim(),
        short_description: formData.short_description.trim() || undefined,
        description: formData.description.trim() || undefined,
        cover_url: formData.cover_url.trim() || undefined,
        intro_video_url: formData.intro_video_url.trim() || undefined,
        difficulty: formData.difficulty,
        language: formData.language,
        visibility: formData.visibility,
        category_id: Number(formData.category_id),
      });

      // 2. Update status if available
      try {
        await api.patch(`/courses/${id}/status`, {
          status: formData.status,
        });
      } catch {
        // Fallback if status endpoint is optional
      }

      showToast('Настройки курса успешно сохранены!');
    } catch (err: any) {
      console.error('Failed to save course settings', err);
      alert(err.response?.data?.message || err.response?.data?.error || 'Ошибка при сохранении');
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex flex-col min-h-screen">
        <TopNavbar title="Загрузка настроек курса..." />
        <main className="p-8 max-w-4xl w-full mx-auto space-y-6 animate-pulse">
          <div className="h-8 bg-slate-200 dark:bg-slate-800 rounded-xl w-1/3" />
          <div className="h-96 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl" />
        </main>
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-screen">
      <TopNavbar
        title="Настройки курса"
        subtitle="Управление метаданными, видимостью и параметрами курса"
      />

      <main className="p-8 max-w-4xl w-full mx-auto space-y-8 flex-1">
        <form
          onSubmit={handleSubmit}
          className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-8 shadow-sm space-y-6"
        >
          <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-4 pb-4 border-b border-slate-100 dark:border-slate-800">
            <div>
              <h2 className="text-xl font-extrabold text-slate-900 dark:text-white">
                Основные параметры курса
              </h2>
              <p className="text-xs text-slate-500">
                ID курса: #{id} • Изменения сразу вступают в силу после сохранения
              </p>
            </div>

            <button
              type="button"
              onClick={() => router.push(`/teacher/courses/${id}/curriculum`)}
              className="px-4 py-2 bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 hover:bg-indigo-100 text-xs font-bold rounded-xl transition-colors flex items-center gap-1.5 w-fit"
            >
              <Sliders size={14} />
              <span>Перейти к программе курса</span>
            </button>
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
                value={formData.title}
                onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                placeholder="Название курса..."
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium text-slate-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:outline-none"
              />
            </div>

            <div className="space-y-1.5">
              <div className="flex items-center justify-between">
                <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                  URL Slug (для адресной строки) *
                </label>
                <button
                  type="button"
                  onClick={handleGenerateSlugFromTitle}
                  className="text-[11px] font-bold text-indigo-600 hover:text-indigo-700 dark:text-indigo-400 flex items-center gap-1 transition-colors cursor-pointer"
                  title="Сгенерировать слаг из названия курса"
                >
                  <RefreshCw size={11} />
                  <span>Сгенерировать из названия</span>
                </button>
              </div>
              <div className="relative">
                <input
                  type="text"
                  required
                  value={formData.slug}
                  onChange={(e) => handleSlugChange(e.target.value)}
                  placeholder="arhitektura-mikroservisov-i-go"
                  className={`w-full p-3 pr-8 bg-slate-50 dark:bg-slate-800 border rounded-xl text-xs font-medium text-slate-900 dark:text-white focus:ring-2 focus:outline-none transition-all ${
                    formData.slug && !isValidSlug(formData.slug)
                      ? 'border-rose-400 focus:ring-rose-400 text-rose-700 dark:text-rose-300'
                      : formData.slug && isValidSlug(formData.slug)
                      ? 'border-emerald-400/80 focus:ring-emerald-400'
                      : 'border-slate-200 dark:border-slate-700 focus:ring-indigo-500'
                  }`}
                />
                {formData.slug && (
                  <div className="absolute right-3 top-1/2 -translate-y-1/2 flex items-center pointer-events-none">
                    {isValidSlug(formData.slug) ? (
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
                {formData.slug && !isValidSlug(formData.slug) ? (
                  <span className="text-rose-500 font-semibold">
                    Допустимы только строчные латинские буквы, цифры и одиночные дефисы
                  </span>
                ) : (
                  <span>
                    Ссылка на курс:{' '}
                    <code className="text-indigo-600 dark:text-indigo-400 font-mono">
                      /courses/{formData.slug || '...'}
                    </code>
                  </span>
                )}
              </p>
            </div>
          </div>

          {/* Short Description */}
          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
              Краткое описание (для карточки в каталоге)
            </label>
            <input
              type="text"
              value={formData.short_description}
              onChange={(e) => setFormData({ ...formData, short_description: e.target.value })}
              placeholder="О чем этот курс в двух словах..."
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium text-slate-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
          </div>

          {/* Full Description */}
          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
              Подробное описание курса
            </label>
            <textarea
              rows={4}
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              placeholder="Расскажите студентам, что они изучат на курсе, какие требования и цели..."
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium text-slate-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:outline-none"
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
                  value={formData.cover_url}
                  onChange={(e) => setFormData({ ...formData, cover_url: e.target.value })}
                  placeholder="https://example.com/cover.jpg или выберите файл"
                  className="flex-1 p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium text-slate-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                />
                {formData.cover_url && (
                  <button
                    type="button"
                    onClick={() => setFormData({ ...formData, cover_url: '' })}
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

              {formData.cover_url && (
                <div className="relative mt-2 w-full h-32 rounded-xl overflow-hidden border border-slate-200 dark:border-slate-700 bg-slate-100 dark:bg-slate-800">
                  <img
                    src={formData.cover_url}
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
                value={formData.intro_video_url}
                onChange={(e) => setFormData({ ...formData, intro_video_url: e.target.value })}
                placeholder="https://www.youtube.com/watch?v=..."
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium text-slate-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:outline-none"
              />
            </div>
          </div>

          {/* Selects: Difficulty, Language, Visibility, Status */}
          <div className="grid grid-cols-1 sm:grid-cols-4 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Сложность
              </label>
              <select
                value={formData.difficulty}
                onChange={(e) => setFormData({ ...formData, difficulty: e.target.value })}
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold text-slate-800 dark:text-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none cursor-pointer"
              >
                <option value="beginner">Начальный (Beginner)</option>
                <option value="intermediate">Средний (Intermediate)</option>
                <option value="advanced">Продвинутый (Advanced)</option>
              </select>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Язык
              </label>
              <select
                value={formData.language}
                onChange={(e) => setFormData({ ...formData, language: e.target.value })}
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold text-slate-800 dark:text-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none cursor-pointer"
              >
                <option value="RU">Русский (RU)</option>
                <option value="EN">English (EN)</option>
              </select>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Видимость
              </label>
              <select
                value={formData.visibility}
                onChange={(e) => setFormData({ ...formData, visibility: e.target.value })}
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold text-slate-800 dark:text-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none cursor-pointer"
              >
                <option value="public">Публичный (в каталоге)</option>
                <option value="private">Приватный (по ссылке)</option>
              </select>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Статус курса
              </label>
              <select
                value={formData.status}
                onChange={(e) => setFormData({ ...formData, status: e.target.value })}
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold text-slate-800 dark:text-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none cursor-pointer"
              >
                <option value="draft">Черновик (draft)</option>
                <option value="published">Опубликован (published)</option>
                <option value="archived">В архиве (archived)</option>
              </select>
            </div>
          </div>

          {/* Form Actions */}
          <div className="flex justify-between items-center pt-6 border-t border-slate-100 dark:border-slate-800">
            <button
              type="button"
              onClick={() => router.push(`/teacher/courses/${id}/curriculum`)}
              className="text-xs font-bold text-slate-500 hover:text-indigo-600 transition-colors"
            >
              ← Вернуться к программе курса
            </button>

            <button
              type="submit"
              disabled={isSaving}
              className="px-8 py-3 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-xl shadow-lg shadow-indigo-600/20 transition-all flex items-center gap-2"
            >
              <Save size={16} />
              <span>{isSaving ? 'Сохранение...' : 'Сохранить настройки'}</span>
            </button>
          </div>
        </form>
      </main>

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
