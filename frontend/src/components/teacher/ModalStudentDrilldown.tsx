'use client';

import React, { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { X, CheckCircle, XCircle, Award, BookOpen, Clock, Sparkles } from 'lucide-react';

interface DrilldownModalProps {
  courseId: number;
  studentId: number | null;
  onClose: () => void;
}

export default function ModalStudentDrilldown({
  courseId,
  studentId,
  onClose,
}: DrilldownModalProps) {
  const [report, setReport] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (!studentId || !courseId) return;

    const fetchDrilldown = async () => {
      setIsLoading(true);
      try {
        const res = await api.get(`/courses/${courseId}/students/${studentId}/drilldown`);
        const data = res.data.data || res.data;
        setReport(data);
      } catch (err) {
        console.error('Failed to load student drilldown', err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchDrilldown();
  }, [courseId, studentId]);

  if (!studentId) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4">
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 max-w-3xl w-full max-h-[90vh] flex flex-col shadow-2xl space-y-5 animate-in fade-in zoom-in duration-200">
        {/* Modal Header */}
        <div className="flex justify-between items-center pb-4 border-b border-slate-100 dark:border-slate-800 flex-shrink-0">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-2xl bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 flex items-center justify-center font-black text-sm">
              {report?.student_name?.[0] || 'S'}
            </div>
            <div>
              <h3 className="text-base font-extrabold text-slate-900 dark:text-white">
                {report?.student_name || 'Детальный отчет студента'}
              </h3>
              <p className="text-xs text-slate-400">
                {report?.student_email} • @{report?.student_username}
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

        {/* Modal Body */}
        {isLoading ? (
          <div className="py-16 text-center text-slate-400 space-y-3">
            <Sparkles className="mx-auto text-indigo-500 animate-spin" size={24} />
            <p className="text-xs font-semibold">Загрузка подробного лога успеваемости...</p>
          </div>
        ) : !report ? (
          <div className="py-12 text-center text-slate-500 text-xs">
            Не удалось загрузить данные по студенту.
          </div>
        ) : (
          <div className="flex-1 overflow-y-auto space-y-6 pr-1 custom-scrollbar">
            {/* Quick Metrics Bar */}
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
              <div className="p-3 bg-slate-50 dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700">
                <span className="text-[10px] uppercase font-bold text-slate-400">Прогресс курса</span>
                <p className="text-lg font-black text-indigo-600 dark:text-indigo-400">
                  {report.progress_percent || 0}%
                </p>
              </div>

              <div className="p-3 bg-slate-50 dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700">
                <span className="text-[10px] uppercase font-bold text-slate-400">Средний балл</span>
                <p className="text-lg font-black text-emerald-600">
                  {report.average_score || 0} / 100
                </p>
              </div>

              <div className="p-3 bg-slate-50 dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 col-span-2 sm:col-span-1">
                <span className="text-[10px] uppercase font-bold text-slate-400">Пройдено уроков</span>
                <p className="text-lg font-black text-purple-600">
                  {report.completed_lessons_count || 0}
                </p>
              </div>
            </div>

            {/* Test Attempts & Quiz Drilldowns */}
            {report.test_attempts && report.test_attempts.length > 0 && (
              <div className="space-y-3">
                <h4 className="text-xs font-extrabold uppercase tracking-wider text-slate-400">
                  Срезы знаний и контрольные тесты
                </h4>

                <div className="space-y-3">
                  {report.test_attempts.map((att: any, aIdx: number) => (
                    <div
                      key={att.id || aIdx}
                      className="p-4 bg-slate-50 dark:bg-slate-800/70 rounded-2xl border border-slate-200 dark:border-slate-700 space-y-3"
                    >
                      <div className="flex justify-between items-center text-xs">
                        <div className="flex items-center gap-2">
                          <Award size={16} className="text-indigo-500" />
                          <span className="font-bold text-slate-900 dark:text-white">
                            {att.quiz_title || `Тестирование #${att.attempt_number || aIdx + 1}`}
                          </span>
                        </div>
                        <span className="font-black px-2.5 py-0.5 rounded-lg bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300">
                          {att.score || 0} / 100
                        </span>
                      </div>

                      {/* Answers detail */}
                      {att.answers && att.answers.length > 0 && (
                        <div className="space-y-2 pt-2 border-t border-slate-200 dark:border-slate-700">
                          {att.answers.map((ans: any, ansIdx: number) => (
                            <div
                              key={ans.id || ansIdx}
                              className="p-3 bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 space-y-1 text-xs"
                            >
                              <p className="font-semibold text-slate-800 dark:text-slate-200">
                                {ans.question_text || `Вопрос #${ansIdx + 1}`}
                              </p>
                              <div className="flex items-center gap-1.5 font-medium">
                                {ans.is_correct ? (
                                  <span className="text-emerald-600 flex items-center gap-1">
                                    <CheckCircle size={14} /> Ответ студента: «{ans.chosen_answer}» (Верно)
                                  </span>
                                ) : (
                                  <span className="text-rose-600 flex items-center gap-1">
                                    <XCircle size={14} /> Ответ студента: «{ans.chosen_answer}» (Неверно)
                                  </span>
                                )}
                              </div>
                              {ans.feedback && (
                                <p className="text-[11px] text-slate-500 bg-slate-50 dark:bg-slate-800 p-2 rounded-lg mt-1">
                                  Рецензия учителя: {ans.feedback}
                                </p>
                              )}
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Lesson history logs */}
            {report.lesson_logs && report.lesson_logs.length > 0 && (
              <div className="space-y-3">
                <h4 className="text-xs font-extrabold uppercase tracking-wider text-slate-400">
                  История прохождения уроков
                </h4>
                <div className="space-y-2">
                  {report.lesson_logs.map((log: any, idx: number) => (
                    <div
                      key={log.lesson_id || idx}
                      className="p-3 bg-slate-50 dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 flex justify-between items-center text-xs"
                    >
                      <div className="flex items-center gap-2">
                        <BookOpen size={14} className="text-slate-400" />
                        <span className="font-bold text-slate-800 dark:text-slate-200">
                          {log.lesson_title}
                        </span>
                      </div>
                      <span className="font-semibold text-emerald-600 dark:text-emerald-400">
                        {log.status || 'Пройден'}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
