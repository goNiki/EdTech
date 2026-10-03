'use client';

import React, { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { X, Save, Trash2 } from 'lucide-react';

interface EditLessonModalProps {
  lessonData: any | null;
  isOpen: boolean;
  onClose: () => void;
  onSaved: () => void;
}

export default function ModalEditLesson({
  lessonData,
  isOpen,
  onClose,
  onSaved,
}: EditLessonModalProps) {
  const [title, setTitle] = useState('');
  const [desc, setDesc] = useState('');
  const [type, setType] = useState('lecture');
  const [status, setStatus] = useState('draft');
  const [isFree, setIsFree] = useState(false);
  const [duration, setDuration] = useState<number | ''>('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    if (lessonData) {
      setTitle(lessonData.Title || lessonData.title || '');
      setDesc(lessonData.Description || lessonData.description || '');
      setType(lessonData.Type || lessonData.type || 'lecture');
      setStatus(lessonData.Status || lessonData.status || 'draft');
      setIsFree(lessonData.IsFree || lessonData.is_free || false);
      setDuration(lessonData.Duration || lessonData.duration || '');
      setErrorMessage(null);
    }
  }, [lessonData]);

  if (!isOpen || !lessonData) return null;

  const lessonId = lessonData.ID || lessonData.id;

  const handleDelete = async () => {
    if (!confirm('Вы уверены, что хотите удалить этот урок?')) return;
    setIsSubmitting(true);
    setErrorMessage(null);
    try {
      await api.delete(`/lessons/${lessonId}`);
      onSaved();
      onClose();
    } catch (err: any) {
      console.error('Failed to delete lesson', err);
      if (err.response?.status === 403) {
        setErrorMessage('У вас нет прав на редактирование этого курса. Изменения не сохранены.');
      } else {
        setErrorMessage(err.response?.data?.message || err.response?.data?.error || 'Ошибка при удалении урока');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setErrorMessage(null);
    try {
      // 1. Update basic fields
      await api.patch(`/lessons/${lessonId}`, {
        title: title.trim(),
        description: desc.trim(),
        type: type,
        is_free: isFree,
        duration: duration ? Number(duration) : undefined,
      });

      // 2. Update status
      await api.patch(`/lessons/${lessonId}/status`, {
        status: status,
      });

      onSaved();
      onClose();
    } catch (err: any) {
      console.error('Failed to update lesson metadata', err);
      if (err.response?.status === 403) {
        setErrorMessage('У вас нет прав на редактирование этого курса. Изменения не сохранены.');
      } else {
        setErrorMessage(err.response?.data?.message || err.response?.data?.error || 'Ошибка при сохранении параметров урока');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4">
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 max-w-lg w-full shadow-2xl space-y-5 animate-in fade-in zoom-in duration-200">
        <div className="flex justify-between items-center pb-3 border-b border-slate-100 dark:border-slate-800">
          <h3 className="text-base font-extrabold text-slate-900 dark:text-white">
            Изменить параметры урока
          </h3>
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
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">Название урока *</label>
            <input
              type="text"
              required
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">Описание урока</label>
            <textarea
              rows={2}
              value={desc}
              onChange={(e) => setDesc(e.target.value)}
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">Тип урока</label>
              <select
                value={type}
                onChange={(e) => setType(e.target.value)}
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold focus:ring-2 focus:ring-indigo-500 focus:outline-none cursor-pointer"
              >
                <option value="lecture">Лекция (Теория)</option>
                <option value="practice">Практика</option>
                <option value="quiz">Контрольный тест</option>
                <option value="homework">Домашнее задание</option>
              </select>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">Статус</label>
              <select
                value={status}
                onChange={(e) => setStatus(e.target.value)}
                className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-semibold focus:ring-2 focus:ring-indigo-500 focus:outline-none cursor-pointer"
              >
                <option value="draft">Черновик (draft)</option>
                <option value="published">Опубликован (published)</option>
                <option value="archived">В архиве (archived)</option>
              </select>
            </div>
          </div>

          <div className="flex items-center gap-3 p-3 bg-slate-50 dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700">
            <input
              type="checkbox"
              id="isFreeCheckbox"
              checked={isFree}
              onChange={(e) => setIsFree(e.target.checked)}
              className="w-4 h-4 rounded text-indigo-600 focus:ring-indigo-500 cursor-pointer"
            />
            <label htmlFor="isFreeCheckbox" className="text-xs font-semibold text-slate-700 dark:text-slate-300 cursor-pointer">
              Бесплатный демо-урок (доступен до зачисления на курс)
            </label>
          </div>

          <div className="flex justify-between items-center pt-4 border-t border-slate-100 dark:border-slate-800">
            <button
              type="button"
              onClick={handleDelete}
              disabled={isSubmitting}
              className="px-3.5 py-2.5 rounded-xl text-xs font-bold text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/40 flex items-center gap-1.5 transition-colors disabled:opacity-40"
              title="Удалить урок"
            >
              <Trash2 size={15} />
              <span>Удалить урок</span>
            </button>

            <div className="flex items-center gap-3">
              <button
                type="button"
                onClick={onClose}
                className="px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-600 dark:text-slate-400 hover:bg-slate-50"
              >
                Отмена
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-xl shadow-md transition-all flex items-center gap-2"
              >
                <Save size={14} />
                <span>{isSubmitting ? 'Сохранение...' : 'Сохранить урок'}</span>
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}
