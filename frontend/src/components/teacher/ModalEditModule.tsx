'use client';

import React, { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { X, ArrowUp, ArrowDown, Save, Trash2 } from 'lucide-react';

interface EditModuleModalProps {
  moduleData: any | null;
  isOpen: boolean;
  onClose: () => void;
  onSaved: () => void;
}

export default function ModalEditModule({
  moduleData,
  isOpen,
  onClose,
  onSaved,
}: EditModuleModalProps) {
  const [title, setTitle] = useState('');
  const [desc, setDesc] = useState('');
  const [status, setStatus] = useState('draft');
  const [lessons, setLessons] = useState<any[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    if (moduleData) {
      const sec = moduleData.Section || moduleData.section || moduleData;
      setTitle(sec.Title || sec.title || '');
      setDesc(sec.Description || sec.description || '');
      setStatus(sec.Status || sec.status || 'draft');
      setLessons(moduleData.Lessons || moduleData.lessons || []);
      setErrorMessage(null);
    }
  }, [moduleData]);

  if (!isOpen || !moduleData) return null;

  const sec = moduleData.Section || moduleData.section || moduleData;
  const secId = sec.ID || sec.id;

  const moveLesson = (fromIdx: number, toIdx: number) => {
    if (toIdx < 0 || toIdx >= lessons.length) return;
    const copy = [...lessons];
    const [moved] = copy.splice(fromIdx, 1);
    copy.splice(toIdx, 0, moved);
    setLessons(copy);
  };

  const handleDelete = async () => {
    if (!confirm('Вы уверены, что хотите удалить этот модуль? Все связанные уроки будут удалены.')) return;
    setIsSubmitting(true);
    setErrorMessage(null);
    try {
      await api.delete(`/sections/${secId}`);
      onSaved();
      onClose();
    } catch (err: any) {
      console.error('Failed to delete module', err);
      if (err.response?.status === 403) {
        setErrorMessage('У вас нет прав на редактирование этого курса. Изменения не сохранены.');
      } else {
        setErrorMessage(err.response?.data?.message || err.response?.data?.error || 'Ошибка при удалении модуля');
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
      // 1. Update title/description
      await api.patch(`/sections/${secId}`, {
        title: title.trim(),
        description: desc.trim(),
      });

      // 2. Update status
      await api.patch(`/sections/${secId}/status`, {
        status: status,
      });

      // 3. Reorder lessons if changed
      const lessonIds = lessons.map((l: any) => Number(l.ID || l.id));
      if (lessonIds.length > 0) {
        await api.put(`/sections/${secId}/reorder-lessons`, {
          item_ids: lessonIds,
        });
      }

      onSaved();
      onClose();
    } catch (err: any) {
      console.error('Failed to update module', err);
      if (err.response?.status === 403) {
        setErrorMessage('У вас нет прав на редактирование этого курса. Изменения не сохранены.');
      } else {
        setErrorMessage(err.response?.data?.message || err.response?.data?.error || 'Ошибка при сохранении модуля');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4">
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 max-w-lg w-full max-h-[90vh] flex flex-col shadow-2xl space-y-5 animate-in fade-in zoom-in duration-200">
        <div className="flex justify-between items-center pb-3 border-b border-slate-100 dark:border-slate-800 flex-shrink-0">
          <h3 className="text-base font-extrabold text-slate-900 dark:text-white">
            Настройки модуля
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

        <form onSubmit={handleSubmit} className="flex-1 overflow-y-auto space-y-4 pr-1">
          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">Название модуля *</label>
            <input
              type="text"
              required
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">Описание модуля</label>
            <textarea
              rows={2}
              value={desc}
              onChange={(e) => setDesc(e.target.value)}
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">Статус публикации</label>
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

          {/* Lessons list order */}
          {lessons.length > 0 && (
            <div className="space-y-2 pt-2 border-t border-slate-100 dark:border-slate-800">
              <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                Порядок уроков в модуле
              </label>
              <div className="space-y-1.5">
                {lessons.map((les: any, lIdx: number) => (
                  <div
                    key={les.ID || les.id || lIdx}
                    className="p-2.5 bg-slate-50 dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 flex items-center justify-between text-xs"
                  >
                    <span className="font-semibold text-slate-800 dark:text-slate-200 truncate">
                      {lIdx + 1}. {les.Title || les.title}
                    </span>
                    <div className="flex items-center gap-1">
                      <button
                        type="button"
                        disabled={lIdx === 0}
                        onClick={() => moveLesson(lIdx, lIdx - 1)}
                        className="p-1 rounded hover:bg-slate-200 dark:hover:bg-slate-700 disabled:opacity-20"
                      >
                        <ArrowUp size={14} />
                      </button>
                      <button
                        type="button"
                        disabled={lIdx === lessons.length - 1}
                        onClick={() => moveLesson(lIdx, lIdx + 1)}
                        className="p-1 rounded hover:bg-slate-200 dark:hover:bg-slate-700 disabled:opacity-20"
                      >
                        <ArrowDown size={14} />
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          <div className="flex justify-between items-center pt-4 border-t border-slate-100 dark:border-slate-800">
            <button
              type="button"
              onClick={handleDelete}
              disabled={isSubmitting}
              className="px-3.5 py-2.5 rounded-xl text-xs font-bold text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/40 flex items-center gap-1.5 transition-colors disabled:opacity-40"
              title="Удалить модуль"
            >
              <Trash2 size={15} />
              <span>Удалить модуль</span>
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
                <span>{isSubmitting ? 'Сохранение...' : 'Сохранить модуль'}</span>
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}
