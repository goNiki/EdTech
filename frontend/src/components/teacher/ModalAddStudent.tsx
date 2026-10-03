'use client';

import React, { useState } from 'react';
import { api } from '@/lib/api';
import { X, UserPlus, Mail, Sparkles } from 'lucide-react';

interface AddStudentModalProps {
  courseId: number;
  isOpen: boolean;
  onClose: () => void;
  onAdded: () => void;
}

export default function ModalAddStudent({
  courseId,
  isOpen,
  onClose,
  onAdded,
}: AddStudentModalProps) {
  const [emailOrId, setEmailOrId] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!emailOrId.trim()) return;

    setIsSubmitting(true);
    try {
      const isNumeric = /^\d+$/.test(emailOrId.trim());
      const payload: any = {};
      if (isNumeric) {
        payload.user_id = Number(emailOrId.trim());
      } else {
        payload.email = emailOrId.trim();
      }

      await api.post(`/courses/${courseId}/students`, payload);
      setEmailOrId('');
      onAdded();
      onClose();
    } catch (err: any) {
      console.error('Failed to enroll student', err);
      alert(err.response?.data?.message || err.response?.data?.error || 'Не удалось зачислить студента. Проверьте Email/ID.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4">
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 max-w-md w-full shadow-2xl space-y-5 animate-in fade-in zoom-in duration-200">
        <div className="flex justify-between items-center pb-3 border-b border-slate-100 dark:border-slate-800">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-xl bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 flex items-center justify-center">
              <UserPlus size={16} />
            </div>
            <h3 className="text-base font-extrabold text-slate-900 dark:text-white">
              Зачислить студента на курс
            </h3>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-xl text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors"
          >
            <X size={20} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-1.5">
            <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
              Email или ID пользователя *
            </label>
            <input
              type="text"
              required
              value={emailOrId}
              onChange={(e) => setEmailOrId(e.target.value)}
              placeholder="student@example.com или 42"
              className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
            />
            <p className="text-[10px] text-slate-400">
              Студент получит доступ ко всем материалам и появится в списке группы.
            </p>
          </div>

          <div className="flex justify-end gap-3 pt-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-600 dark:text-slate-400 hover:bg-slate-50"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={isSubmitting || !emailOrId.trim()}
              className="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-40 text-white text-xs font-bold rounded-xl shadow-md transition-all flex items-center gap-2"
            >
              <UserPlus size={14} />
              <span>{isSubmitting ? 'Зачисление...' : 'Зачислить'}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
