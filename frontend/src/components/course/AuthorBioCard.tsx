'use client';

import React from 'react';
import {
  ShieldCheck,
  GraduationCap,
  BookOpen,
  Star,
  Award,
  Sparkles
} from 'lucide-react';

export interface CourseAuthor {
  id: number;
  name: string;
  avatar_url?: string;
  headline?: string;
  bio?: string;
  courses_count?: number;
  total_students?: number;
}

interface AuthorBioCardProps {
  author?: CourseAuthor | null;
  className?: string;
}

export default function AuthorBioCard({ author, className = '' }: AuthorBioCardProps) {
  if (!author) return null;

  const initials = (author.name || 'П')
    .split(' ')
    .filter(Boolean)
    .map((n) => n[0])
    .join('')
    .slice(0, 2)
    .toUpperCase();

  const headline = author.headline?.trim() || 'Преподаватель и ведущий автор курса';
  const bio = author.bio?.trim() || 'Практикующий эксперт с глубоким опытом в предметной области. Обучает студентов структурированному подходу, практическим навыкам и подготовке к экзаменам на высший балл.';
  const studentsCount = author.total_students || 0;
  const coursesCount = author.courses_count || 1;

  return (
    <div
      className={`bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 sm:p-8 shadow-sm space-y-6 ${className}`}
    >
      <div className="flex items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <Award size={18} className="text-indigo-600 dark:text-indigo-400" />
          <h3 className="text-xs font-black uppercase tracking-wider text-slate-500 dark:text-slate-400">
            Преподаватель курса
          </h3>
        </div>

        <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-xl bg-indigo-50 dark:bg-indigo-950/60 border border-indigo-200/80 dark:border-indigo-800/60 text-indigo-700 dark:text-indigo-300 text-xs font-extrabold shadow-2xs">
          <ShieldCheck size={14} className="text-indigo-600 dark:text-indigo-400" />
          <span>Верифицирован</span>
        </span>
      </div>

      <div className="flex flex-col sm:flex-row items-center sm:items-start gap-6">
        {/* Avatar */}
        <div className="relative flex-shrink-0">
          {author.avatar_url ? (
            <img
              src={author.avatar_url}
              alt={author.name}
              className="w-20 h-20 sm:w-24 sm:h-24 rounded-full object-cover ring-4 ring-indigo-500/20 shadow-md"
            />
          ) : (
            <div className="w-20 h-20 sm:w-24 sm:h-24 rounded-full bg-gradient-to-tr from-indigo-600 via-purple-600 to-indigo-800 text-white font-black text-2xl sm:text-3xl flex items-center justify-center ring-4 ring-indigo-500/20 shadow-md">
              {initials}
            </div>
          )}
          <div
            className="absolute -bottom-1 -right-1 w-7 h-7 rounded-full bg-emerald-500 text-white border-2 border-white dark:border-slate-900 flex items-center justify-center shadow-xs"
            title="Сертифицированный преподаватель платформы"
          >
            <ShieldCheck size={14} />
          </div>
        </div>

        {/* Info */}
        <div className="flex-1 text-center sm:text-left space-y-2">
          <div>
            <h4 className="text-xl sm:text-2xl font-black text-slate-900 dark:text-white leading-tight">
              {author.name}
            </h4>
            <p className="text-xs sm:text-sm font-bold text-indigo-600 dark:text-indigo-400 mt-1">
              {headline}
            </p>
          </div>

          <p className="text-xs sm:text-sm text-slate-600 dark:text-slate-300 leading-relaxed pt-1">
            {bio}
          </p>

          {/* Metrics */}
          <div className="flex flex-wrap items-center justify-center sm:justify-start gap-3 pt-3">
            <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200/80 dark:border-slate-700 text-xs font-bold text-slate-700 dark:text-slate-200 shadow-2xs">
              <GraduationCap size={15} className="text-indigo-600 dark:text-indigo-400" />
              <span>
                <strong>{studentsCount}</strong> {studentsCount === 1 ? 'студент' : 'студентов'}
              </span>
            </div>

            <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200/80 dark:border-slate-700 text-xs font-bold text-slate-700 dark:text-slate-200 shadow-2xs">
              <BookOpen size={14} className="text-indigo-600 dark:text-indigo-400" />
              <span>
                <strong>{coursesCount}</strong> {coursesCount === 1 ? 'курс' : 'курса'}
              </span>
            </div>

            <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-amber-50 dark:bg-amber-950/40 border border-amber-200/80 dark:border-amber-800/60 text-xs font-bold text-amber-700 dark:text-amber-300 shadow-2xs">
              <Star size={14} className="fill-amber-400 text-amber-400" />
              <span>
                <strong>4.9</strong> рейтинг
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
