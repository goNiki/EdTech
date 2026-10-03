'use client';

import React from 'react';
import { UserCheck, FileText, Sparkles } from 'lucide-react';

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
  onOpenDrilldown: (studentId: number) => void;
  onAddStudent: () => void;
}

export default function StudentsTable({
  students,
  onOpenDrilldown,
  onAddStudent,
}: StudentsTableProps) {
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
        <button
          onClick={onAddStudent}
          className="px-5 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs rounded-xl shadow-md transition-all"
        >
          Зачислить студента
        </button>
      </div>
    );
  }

  return (
    <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl overflow-hidden shadow-xs">
      <div className="p-5 border-b border-slate-100 dark:border-slate-800 flex justify-between items-center">
        <div>
          <h4 className="text-sm font-bold text-slate-900 dark:text-white">
            Зачисленные студенты ({students.length})
          </h4>
          <p className="text-[10px] text-slate-400">
            Сводная таблица успеваемости и доступа к урокам
          </p>
        </div>
        <button
          onClick={onAddStudent}
          className="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-xl shadow-xs transition-colors"
        >
          + Добавить студента
        </button>
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
                  <button
                    onClick={() => onOpenDrilldown(st.id)}
                    className="px-3.5 py-1.5 bg-indigo-50 dark:bg-indigo-950 hover:bg-indigo-100 dark:hover:bg-indigo-900 text-indigo-600 dark:text-indigo-400 font-bold rounded-xl transition-all shadow-2xs"
                  >
                    Детальный отчет
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
