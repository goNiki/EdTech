'use client';

import React, { useState } from 'react';
import { api } from '@/lib/api';
import { X, PlusCircle, BookOpen, Clock, Sparkles } from 'lucide-react';

interface CreateLessonModalProps {
  courseId: number;
  sectionId: number | null;
  sectionTitle?: string;
  isOpen: boolean;
  orderIndex: number;
  onClose: () => void;
  onCreated: () => void;
}

export default function ModalCreateLesson({
  courseId,
  sectionId,
  sectionTitle,
  isOpen,
  orderIndex,
  onClose,
  onCreated,
}: CreateLessonModalProps) {
  const [title, setTitle] = useState('');
  const [desc, setDesc] = useState('');
  const [type, setType] = useState('lecture');
  const [status, setStatus] = useState<'draft' | 'published'>('draft');
  const [isFree, setIsFree] = useState(false);
  const [duration, setDuration] = useState<number | ''>(20);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  if (!isOpen || !sectionId) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    setIsSubmitting(true);
    setErrorMessage(null);
    try {
      const res = await api.post('/lessons', {
        course_id: Number(courseId),
        section_id: Number(sectionId),
        title: title.trim(),
        description: desc.trim() || undefined,
        content: '{}',
        type: type,
        is_free: isFree,
        duration: duration ? Number(duration) : undefined,
      });

      const data = res.data.data || res.data;
      const lessonId = data.id || data.ID || data.lesson?.id;

      // Update status if published was chosen
      if (status === 'published' && lessonId) {
        await api.patch(`/lessons/${lessonId}/status`, {
          status: 'published',
        });
      }

      setTitle('');
      setDesc('');
      setType('lecture');
      setStatus('draft');
      setIsFree(false);
      setDuration(20);
      setErrorMessage(null);
      onCreated();
      onClose();
    } catch (err: any) {
      console.error('Failed to create lesson', err);
      if (err.response?.status === 403) {
        setErrorMessage('У вас нет прав на редактирование этого курса. Изменения не сохранены.');
      } else {
        setErrorMessage(err.response?.data?.message || err.response?.data?.error || 'Ошибка при создании урока');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4">
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 max-w-lg w-full shadow-2xl space-y-5 animate-in fade-in zoom-in duration-200">
        <div className="flex justify-between items-center pb-3 border-b border-slate-100 dark:border-slate-800">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-xl bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 flex items-center justify-center">
              <BookOpen size={16} />
            </div>
            <div>
              <h3 className="text-base font-extrabold text-slate-900 dark:text-white">
                Новый урок
              </h3>
              <p className="text-[11px] text-slate-400">
                Модуль: {sectionTitle || 'Текущий модуль'} • Урок #{orderIndex}
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-xl text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors"
          >
            <X size={20} />
          </button>
        </div>

        {errorMessage && (
          <div className="p-3.5 rounded-xl bg-rose-50 dark:bg-rose-950/60 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 text-xs font-semibold flex items-center justify-between">
            <span>{errorMessage}</span>
            <button
              type="button"
              onClick={() => setErrorMessage(null)}
              className="text-rose-400 hover:text-rose-600 ml-2"
            >
              <X size={14} />
            </button>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
              Название урока *
            </label>
            <input
              type="text"
              required
              autoFocus
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Например: Урок 1. Конфигурация окружения и первый роут"
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium text-slate-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
              Описание урока (краткий анонс)
            </label>
            <textarea
              rows={2}
              value={desc}
              onChange={(e) => setDesc(e.target.value)}
              placeholder="Кратко расскажите, какие темы будут рассмотрены..."
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium text-slate-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Тип занятия
              </label>
              <select
                value={type}
                onChange={(e) => setType(e.target.value)}
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold text-slate-800 dark:text-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none cursor-pointer"
              >
                <option value="lecture">Лекция (Теория)</option>
                <option value="practice">Практика</option>
                <option value="quiz">Контрольный тест</option>
                <option value="homework">Домашнее задание</option>
              </select>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Длительность (минут)
              </label>
              <input
                type="number"
                min={1}
                value={duration}
                onChange={(e) => setDuration(e.target.value ? Number(e.target.value) : '')}
                placeholder="20"
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold text-slate-800 dark:text-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none"
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
              Статус публикации
            </label>
            <select
              value={status}
              onChange={(e) => setStatus(e.target.value as 'draft' | 'published')}
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold text-slate-800 dark:text-slate-200 focus:ring-2 focus:ring-indigo-500 focus:outline-none cursor-pointer"
            >
              <option value="draft">Черновик (наполняется в редакторе Puck)</option>
              <option value="published">Опубликован (сразу доступен студентам)</option>
            </select>
          </div>

          <div className="flex items-center gap-3 p-3 bg-slate-50 dark:bg-slate-800/80 rounded-xl border border-slate-200 dark:border-slate-700">
            <input
              type="checkbox"
              id="createLessonIsFree"
              checked={isFree}
              onChange={(e) => setIsFree(e.target.checked)}
              className="w-4 h-4 rounded text-indigo-600 focus:ring-indigo-500 cursor-pointer"
            />
            <label
              htmlFor="createLessonIsFree"
              className="text-xs font-semibold text-slate-700 dark:text-slate-300 cursor-pointer"
            >
              Бесплатный демо-урок (доступен до зачисления на курс)
            </label>
          </div>

          <div className="flex justify-end gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={isSubmitting || !title.trim()}
              className="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-40 text-white text-xs font-bold rounded-xl shadow-md transition-all flex items-center gap-2"
            >
              <PlusCircle size={15} />
              <span>{isSubmitting ? 'Создание...' : 'Создать урок'}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
