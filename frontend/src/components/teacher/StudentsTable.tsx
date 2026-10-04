'use client';

import React, { useState } from 'react';
import { UserCheck, FileText, Sparkles, UserX, AlertTriangle, Loader2, Download, FileSpreadsheet } from 'lucide-react';
import { api } from '@/lib/api';

export interface StudentItem {
  id: number;
  name: string;
  email: string;
  username: string;
  progress: number;
  score: number;
  completed_lessons?: number;
  has_pending?: boolean;
}

interface StudentsTableProps {
  students: StudentItem[];
  courseId?: string | number;
  onOpenDrilldown: (studentId: number) => void;
  onAddStudent: () => void;
  onRemoveStudent?: (student: StudentItem) => Promise<void> | void;
  showToast?: (msg: string, type?: 'success' | 'error') => void;
}

export default function StudentsTable({
  students,
  courseId,
  onOpenDrilldown,
  onAddStudent,
  onRemoveStudent,
  showToast,
}: StudentsTableProps) {
  const [studentToRemove, setStudentToRemove] = useState<StudentItem | null>(null);
  const [isRemoving, setIsRemoving] = useState(false);
  const [isExporting, setIsExporting] = useState(false);

  const handleExportGradebook = async () => {
    if (!courseId) return;
    setIsExporting(true);
    try {
      const res = await api.get(`/courses/${courseId}/analytics/export?format=csv`, {
        responseType: 'blob',
      });

      let filename = `gradebook_course_${courseId}_${new Date().toISOString().split('T')[0]}.csv`;
      const disposition = res.headers ? res.headers['content-disposition'] : null;
      if (disposition && disposition.includes('filename=')) {
        const match = disposition.match(/filename=["']?([^"';]+)["']?/);
        if (match && match[1]) {
          filename = match[1];
        }
      }

      const blob = new Blob([res.data], { type: 'text/csv;charset=utf-8;' });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', filename);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);

      if (showToast) {
        showToast('Ведомость курса успешно выгружена');
      }
    } catch (err: any) {
      console.error('Failed to export gradebook', err);
      if (showToast) {
        showToast(
          err.response?.data?.message || err.response?.data?.error || 'Ошибка при выгрузке ведомости курса',
          'error'
        );
      }
    } finally {
      setIsExporting(false);
    }
  };

  const handleConfirmRemove = async () => {
    if (!studentToRemove || !onRemoveStudent) return;
    setIsRemoving(true);
    try {
      await onRemoveStudent(studentToRemove);
      setStudentToRemove(null);
    } finally {
      setIsRemoving(false);
    }
  };
  if (students.length === 0) {
    return (
      <div className="p-12 text-center bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-4">
        <UserCheck size={32} className="mx-auto text-slate-400" />
        <h4 className="text-sm font-bold text-slate-700 dark:text-slate-300">
          На данный курс еще никто не записан
        </h4>
        <p className="text-xs text-slate-500 max-w-sm mx-auto">
          Вы можете зачислить учеников вручную по Email или ID, либо опубликовать курс в открытом каталоге.
        </p>
        <div className="flex items-center justify-center gap-3 flex-wrap">
          {courseId && (
            <button
              type="button"
              disabled={isExporting}
              onClick={handleExportGradebook}
              className="px-4 py-2.5 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 font-bold text-xs rounded-xl shadow-2xs transition-all flex items-center gap-1.5 cursor-pointer disabled:opacity-60"
              title="Выгрузить ведомость успеваемости в формате CSV"
            >
              {isExporting ? (
                <>
                  <Loader2 size={14} className="animate-spin text-indigo-600" />
                  <span>Формирование отчета...</span>
                </>
              ) : (
                <>
                  <Download size={14} />
                  <span>Выгрузить ведомость (CSV)</span>
                </>
              )}
            </button>
          )}
          <button
            onClick={onAddStudent}
            className="px-5 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs rounded-xl shadow-md transition-all cursor-pointer"
          >
            Зачислить студента
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl overflow-hidden shadow-xs">
      <div className="p-5 border-b border-slate-100 dark:border-slate-800 flex flex-col sm:flex-row justify-between sm:items-center gap-4">
        <div>
          <h4 className="text-sm font-bold text-slate-900 dark:text-white">
            Зачисленные студенты ({students.length})
          </h4>
          <p className="text-[10px] text-slate-400">
            Сводная таблица успеваемости и доступа к урокам
          </p>
        </div>
        <div className="flex items-center gap-2.5 flex-wrap">
          {courseId && (
            <button
              type="button"
              disabled={isExporting}
              onClick={handleExportGradebook}
              className="px-3.5 py-2 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 text-xs font-bold rounded-xl transition-all shadow-2xs flex items-center gap-1.5 cursor-pointer disabled:opacity-60"
              title="Выгрузить ведомость успеваемости в формате CSV / Excel"
            >
              {isExporting ? (
                <>
                  <Loader2 size={14} className="animate-spin text-indigo-600" />
                  <span>Формирование отчета...</span>
                </>
              ) : (
                <>
                  <Download size={14} />
                  <span>Выгрузить ведомость (CSV)</span>
                </>
              )}
            </button>
          )}

          <button
            onClick={onAddStudent}
            className="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-xl shadow-xs transition-colors flex items-center gap-1.5 cursor-pointer"
          >
            <span>+ Добавить студента</span>
          </button>
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-slate-50 dark:bg-slate-800/60 border-b border-slate-100 dark:border-slate-800 text-[10px] font-extrabold uppercase tracking-wider text-slate-400">
              <th className="px-5 py-3">Студент</th>
              <th className="px-5 py-3">Прогресс</th>
              <th className="px-5 py-3">Средний балл</th>
              <th className="px-5 py-3">Статус ДЗ</th>
              <th className="px-5 py-3 text-right">Действия</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 dark:divide-slate-800 text-xs">
            {students.map((st) => (
              <tr
                key={st.id}
                className="hover:bg-slate-50/70 dark:hover:bg-slate-800/40 transition-colors"
              >
                <td className="px-5 py-4">
                  <p className="font-bold text-slate-900 dark:text-white">{st.name}</p>
                  <p className="text-[10px] text-slate-400 mt-0.5">
                    {st.email} • @{st.username}
                  </p>
                </td>
                <td className="px-5 py-4">
                  <div className="flex items-center gap-2">
                    <div className="bg-slate-100 dark:bg-slate-800 h-2 rounded-full w-24 overflow-hidden">
                      <div
                        className="bg-indigo-600 h-full rounded-full"
                        style={{ width: `${st.progress}%` }}
                      />
                    </div>
                    <span className="font-bold text-slate-700 dark:text-slate-300">
                      {st.progress}%
                    </span>
                  </div>
                </td>
                <td className="px-5 py-4">
                  <span
                    className={`font-black ${
                      st.score >= 80 ? 'text-emerald-600' : 'text-amber-600'
                    }`}
                  >
                    {st.score} / 100
                  </span>
                </td>
                <td className="px-5 py-4">
                  {st.has_pending ? (
                    <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-[10px] font-bold bg-rose-50 dark:bg-rose-950/80 text-rose-700 dark:text-rose-300">
                      <span className="w-1.5 h-1.5 rounded-full bg-rose-500 animate-pulse" />
                      Есть ДЗ на проверку
                    </span>
                  ) : (
                    <span className="inline-flex items-center px-2.5 py-1 rounded-lg text-[10px] font-medium bg-emerald-50 dark:bg-emerald-950/80 text-emerald-700 dark:text-emerald-300">
                      Все проверено
                    </span>
                  )}
                </td>
                <td className="px-5 py-4 text-right">
                  <div className="flex items-center justify-end gap-2">
                    <button
                      onClick={() => onOpenDrilldown(st.id)}
                      className="px-3.5 py-1.5 bg-indigo-50 dark:bg-indigo-950 hover:bg-indigo-100 dark:hover:bg-indigo-900 text-indigo-600 dark:text-indigo-400 font-bold rounded-xl transition-all shadow-2xs text-xs"
                    >
                      Детальный отчет
                    </button>
                    {onRemoveStudent && (
                      <button
                        onClick={() => setStudentToRemove(st)}
                        className="p-1.5 rounded-xl border border-rose-200/80 dark:border-rose-900/60 text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/40 hover:text-rose-600 transition-colors"
                        title="Исключить студента из курса"
                      >
                        <UserX size={15} />
                      </button>
                    )}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Remove Student Confirmation Modal */}
      {studentToRemove && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs animate-in fade-in duration-200">
          <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 sm:p-8 max-w-md w-full shadow-2xl space-y-6">
            <div className="w-12 h-12 rounded-2xl bg-rose-50 dark:bg-rose-950/60 text-rose-600 dark:text-rose-400 flex items-center justify-center mx-auto">
              <AlertTriangle size={24} />
            </div>

            <div className="text-center space-y-2">
              <h3 className="text-lg font-extrabold text-slate-900 dark:text-white">
                Исключить студента из курса?
              </h3>
              <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
                Вы собираетесь отчислить учащегося{' '}
                <strong className="text-slate-800 dark:text-slate-200">{studentToRemove.name}</strong> ({studentToRemove.email}). 
                Он потеряет доступ к материалам курса, а счетчик учащихся будет уменьшен.
              </p>
            </div>

            <div className="flex gap-3">
              <button
                type="button"
                disabled={isRemoving}
                onClick={() => setStudentToRemove(null)}
                className="flex-1 py-3 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors"
              >
                Отмена
              </button>
              <button
                type="button"
                disabled={isRemoving}
                onClick={handleConfirmRemove}
                className="flex-1 py-3 rounded-xl bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold shadow-md shadow-rose-600/20 transition-all flex items-center justify-center gap-2 cursor-pointer"
              >
                {isRemoving ? (
                  <>
                    <Loader2 size={15} className="animate-spin" />
                    <span>Исключение...</span>
                  </>
                ) : (
                  <span>Исключить</span>
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
