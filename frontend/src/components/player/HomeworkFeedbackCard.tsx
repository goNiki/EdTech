'use client';

import React from 'react';
import {
  Clock,
  CheckCircle2,
  AlertCircle,
  MessageSquareQuote,
  UserCheck,
  Award,
  FileText,
  ChevronDown,
  ChevronUp
} from 'lucide-react';

export interface HomeworkAnswerFeedback {
  answer_id: number;
  question_text: string;
  student_answer: string;
  points: number;
  max_points: number;
  is_correct: boolean | null;
  feedback?: string;
}

export interface HomeworkFeedbackData {
  has_submission: boolean;
  status: 'pending' | 'graded';
  attempt_id: number;
  submitted_at: string;
  graded_at?: string;
  teacher?: {
    id: number;
    name: string;
    avatar_url?: string;
  };
  answers: HomeworkAnswerFeedback[];
}

interface HomeworkFeedbackCardProps {
  data: HomeworkFeedbackData | null;
  className?: string;
}

export default function HomeworkFeedbackCard({ data, className = '' }: HomeworkFeedbackCardProps) {
  const [isExpanded, setIsExpanded] = React.useState(true);

  if (!data || !data.has_submission) return null;

  const isPending = data.status === 'pending';
  const totalEarned = data.answers.reduce((acc, a) => acc + (a.points || 0), 0);
  const totalMax = data.answers.reduce((acc, a) => acc + (a.max_points || 0), 0);
  const percentage = totalMax > 0 ? Math.round((totalEarned / totalMax) * 100) : 100;
  const isPassed = !isPending && percentage >= 60;

  // Format date helper
  const formatDate = (isoString?: string) => {
    if (!isoString) return '';
    try {
      const d = new Date(isoString);
      return d.toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch {
      return isoString;
    }
  };

  return (
    <div
      className={`rounded-3xl border transition-all duration-300 shadow-sm overflow-hidden ${
        isPending
          ? 'bg-gradient-to-br from-indigo-50/80 via-white to-blue-50/50 dark:from-indigo-950/40 dark:via-slate-900 dark:to-slate-900 border-indigo-200 dark:border-indigo-800/60'
          : isPassed
          ? 'bg-gradient-to-br from-emerald-50/80 via-white to-teal-50/50 dark:from-emerald-950/40 dark:via-slate-900 dark:to-slate-900 border-emerald-200 dark:border-emerald-800/60'
          : 'bg-gradient-to-br from-amber-50/80 via-white to-orange-50/50 dark:from-amber-950/40 dark:via-slate-900 dark:to-slate-900 border-amber-200 dark:border-amber-800/60'
      } ${className}`}
    >
      {/* Header Bar */}
      <div className="p-5 sm:p-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/60 dark:border-slate-800/60">
        <div className="flex items-start sm:items-center gap-3.5">
          <div
            className={`w-12 h-12 rounded-2xl flex items-center justify-center flex-shrink-0 shadow-xs ${
              isPending
                ? 'bg-indigo-600 text-white'
                : isPassed
                ? 'bg-emerald-600 text-white'
                : 'bg-amber-600 text-white'
            }`}
          >
            {isPending ? (
              <Clock size={24} className="animate-spin-slow" />
            ) : isPassed ? (
              <CheckCircle2 size={24} />
            ) : (
              <AlertCircle size={24} />
            )}
          </div>

          <div>
            <div className="flex flex-wrap items-center gap-2">
              <span
                className={`text-[11px] font-black uppercase tracking-wider px-2.5 py-0.5 rounded-lg ${
                  isPending
                    ? 'bg-indigo-100 dark:bg-indigo-900/60 text-indigo-700 dark:text-indigo-300'
                    : isPassed
                    ? 'bg-emerald-100 dark:bg-emerald-900/60 text-emerald-700 dark:text-emerald-300'
                    : 'bg-amber-100 dark:bg-amber-900/60 text-amber-700 dark:text-amber-300'
                }`}
              >
                {isPending
                  ? 'На проверке у преподавателя'
                  : isPassed
                  ? '✓ Задание проверено'
                  : 'Требуется доработка'}
              </span>

              {data.submitted_at && (
                <span className="text-xs text-slate-500 dark:text-slate-400">
                  Сдано: {formatDate(data.submitted_at)}
                </span>
              )}
            </div>

            <h3 className="text-base sm:text-lg font-black text-slate-900 dark:text-white mt-1">
              {isPending
                ? 'Письменная работа принята на проверку'
                : 'Рецензия и оценка преподавателя'}
            </h3>
          </div>
        </div>

        {/* Score badge / toggle */}
        <div className="flex items-center gap-3 self-end sm:self-auto">
          {!isPending && totalMax > 0 && (
            <div className="flex items-baseline gap-1.5 px-4 py-2 rounded-2xl bg-white dark:bg-slate-800/80 border border-slate-200/80 dark:border-slate-700 shadow-2xs">
              <span
                className={`text-xl font-black font-mono ${
                  isPassed ? 'text-emerald-600 dark:text-emerald-400' : 'text-amber-600 dark:text-amber-400'
                }`}
              >
                {totalEarned}
              </span>
              <span className="text-xs font-bold text-slate-400">/ {totalMax} б.</span>
              <span className="text-[10px] font-extrabold text-slate-400 ml-1">({percentage}%)</span>
            </div>
          )}

          <button
            type="button"
            onClick={() => setIsExpanded(!isExpanded)}
            className="p-2 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
            title={isExpanded ? 'Свернуть рецензию' : 'Развернуть рецензию'}
          >
            {isExpanded ? <ChevronUp size={20} /> : <ChevronDown size={20} />}
          </button>
        </div>
      </div>

      {/* Expandable Body */}
      {isExpanded && (
        <div className="p-5 sm:p-6 space-y-6">
          {/* Status info bar */}
          {isPending ? (
            <div className="p-4 rounded-2xl bg-indigo-50/70 dark:bg-indigo-950/30 border border-indigo-100 dark:border-indigo-900/50 flex items-start gap-3 text-xs text-indigo-900 dark:text-indigo-200 leading-relaxed">
              <Clock size={18} className="text-indigo-600 dark:text-indigo-400 flex-shrink-0 mt-0.5" />
              <div>
                <p className="font-bold">
                  Ваши ответы на открытые задания и эссе успешно отправлены.
                </p>
                <p className="text-indigo-700/80 dark:text-indigo-300/80 mt-1">
                  Преподаватель проверит полноту ответа, аргументацию и выставит итоговый балл. Рецензия отобразится прямо на этой странице, а вы получите системное уведомление.
                </p>
              </div>
            </div>
          ) : (
            <div className="flex flex-wrap items-center justify-between gap-3 p-3.5 rounded-2xl bg-slate-100/70 dark:bg-slate-800/50 border border-slate-200/60 dark:border-slate-700/50 text-xs">
              <div className="flex items-center gap-2 text-slate-700 dark:text-slate-300">
                <UserCheck size={16} className="text-emerald-500" />
                <span>
                  Проверил(а):{' '}
                  <strong className="text-slate-900 dark:text-white font-extrabold">
                    {data.teacher?.name || 'Преподаватель курса'}
                  </strong>
                </span>
              </div>
              {data.graded_at && (
                <span className="text-slate-400">
                  Дата проверки: {formatDate(data.graded_at)}
                </span>
              )}
            </div>
          )}

          {/* Answers and Feedback List */}
          <div className="space-y-4">
            <h4 className="text-xs font-black uppercase tracking-wider text-slate-500 dark:text-slate-400 flex items-center gap-1.5">
              <FileText size={14} />
              <span>Проверенные вопросы и задания ({data.answers.length})</span>
            </h4>

            {data.answers.map((ans, idx) => (
              <div
                key={ans.answer_id || idx}
                className="p-4 sm:p-5 rounded-2xl bg-white dark:bg-slate-800/90 border border-slate-200/80 dark:border-slate-700 shadow-2xs space-y-3.5"
              >
                {/* Question and Score */}
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <span className="text-[10px] font-black uppercase tracking-wider text-indigo-600 dark:text-indigo-400">
                      Вопрос #{idx + 1}
                    </span>
                    <p className="text-xs sm:text-sm font-extrabold text-slate-900 dark:text-white mt-0.5 leading-snug">
                      {ans.question_text}
                    </p>
                  </div>

                  {!isPending && (
                    <div className="flex items-center gap-1.5 px-3 py-1 rounded-xl bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-800 text-xs font-black flex-shrink-0">
                      <Award size={13} className={ans.points > 0 ? 'text-amber-500' : 'text-slate-400'} />
                      <span className={ans.points > 0 ? 'text-slate-900 dark:text-white' : 'text-slate-400'}>
                        {ans.points}
                      </span>
                      <span className="text-slate-400 text-[10px]">/ {ans.max_points} б.</span>
                    </div>
                  )}
                </div>

                {/* Student's Answer */}
                <div className="p-3.5 rounded-xl bg-slate-50 dark:bg-slate-900/60 border border-slate-100 dark:border-slate-800 text-xs space-y-1">
                  <span className="text-[10px] uppercase font-bold text-slate-400 tracking-wider">
                    Ваш ответ:
                  </span>
                  <p className="text-slate-700 dark:text-slate-300 whitespace-pre-wrap leading-relaxed font-sans">
                    {ans.student_answer || '— (ответ пуст)'}
                  </p>
                </div>

                {/* Teacher's Feedback */}
                {!isPending && (
                  <div
                    className={`p-4 rounded-xl border text-xs space-y-1.5 ${
                      ans.feedback
                        ? 'bg-indigo-50/60 dark:bg-indigo-950/40 border-indigo-100 dark:border-indigo-900/50'
                        : 'bg-slate-50 dark:bg-slate-900/40 border-slate-100 dark:border-slate-800/60'
                    }`}
                  >
                    <div className="flex items-center gap-1.5 text-indigo-600 dark:text-indigo-400">
                      <MessageSquareQuote size={15} />
                      <span className="text-[11px] font-black uppercase tracking-wider">
                        Рецензия преподавателя
                      </span>
                    </div>

                    {ans.feedback ? (
                      <p className="text-slate-800 dark:text-slate-200 whitespace-pre-wrap italic leading-relaxed pl-1">
                        «{ans.feedback}»
                      </p>
                    ) : (
                      <p className="text-slate-400 italic text-[11px] pl-1">
                        Преподаватель выставил оценку без развернутого текстового комментария.
                      </p>
                    )}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
