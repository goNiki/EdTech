'use client';

import React, { useState } from 'react';
import { api } from '@/lib/api';
import { X, PlusCircle, Layers, Sparkles } from 'lucide-react';

interface CreateModuleModalProps {
  courseId: number;
  isOpen: boolean;
  orderIndex: number;
  onClose: () => void;
  onCreated: () => void;
}

export default function ModalCreateModule({
  courseId,
  isOpen,
  orderIndex,
  onClose,
  onCreated,
}: CreateModuleModalProps) {
  const [title, setTitle] = useState('');
  const [desc, setDesc] = useState('');
  const [status, setStatus] = useState<'draft' | 'published'>('draft');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    setIsSubmitting(true);
    setErrorMessage(null);
    try {
      const res = await api.post('/sections', {
        course_id: Number(courseId),
        title: title.trim(),
        description: desc.trim() || undefined,
        order_index: orderIndex,
      });

      const data = res.data.data || res.data;
      const sectionId = data.id || data.ID || data.section?.id;

      // Update status if published was chosen
      if (status === 'published' && sectionId) {
        await api.patch(`/sections/${sectionId}/status`, {
          status: 'published',
        });
      }

      setTitle('');
      setDesc('');
      setStatus('draft');
      setErrorMessage(null);
      onCreated();
      onClose();
    } catch (err: any) {
      console.error('Failed to create module', err);
      if (err.response?.status === 403) {
        setErrorMessage('У вас нет прав на редактирование этого курса. Изменения не сохранены.');
      } else {
        setErrorMessage(err.response?.data?.message || err.response?.data?.error || 'Ошибка при создании модуля');
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
              <Layers size={16} />
            </div>
            <div>
              <h3 className="text-base font-extrabold text-slate-900 dark:text-white">
                Новый учебный модуль
              </h3>
              <p className="text-[11px] text-slate-400">Порядковый номер модуля: #{orderIndex}</p>
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
              Название модуля *
            </label>
            <input
              type="text"
              required
              autoFocus
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Например: Модуль 1. Основы архитектуры и синтаксиса"
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium text-slate-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
              Описание модуля (цели и результаты)
            </label>
            <textarea
              rows={3}
              value={desc}
              onChange={(e) => setDesc(e.target.value)}
              placeholder="Опишите, какие навыки получит студент после прохождения этого модуля..."
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium text-slate-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
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
              <option value="draft">Черновик (скрыт от студентов до наполнения)</option>
              <option value="published">Опубликован (доступен сразу после сохранения)</option>
            </select>
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
              <span>{isSubmitting ? 'Создание...' : 'Создать модуль'}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
