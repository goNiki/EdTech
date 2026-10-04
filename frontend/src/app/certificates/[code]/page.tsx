'use client';

import React, { useEffect, useState, use } from 'react';
import Link from 'next/link';
import { CertificateData, verifyCertificate } from '@/lib/certificates';
import TopNavbar from '@/components/layout/TopNavbar';
import CertificateModal from '@/components/certificate/CertificateModal';
import {
  ShieldCheck,
  AlertTriangle,
  Award,
  Calendar,
  User,
  BookOpen,
  ArrowRight,
  Printer,
  Search,
  ExternalLink,
  Sparkles
} from 'lucide-react';

export default function CertificateVerificationPage({
  params,
}: {
  params: Promise<{ code: string }>;
}) {
  const { code } = use(params);
  const [cert, setCert] = useState<CertificateData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [searchCode, setSearchCode] = useState('');

  useEffect(() => {
    const loadCert = async () => {
      setIsLoading(true);
      try {
        const data = await verifyCertificate(code);
        setCert(data);
      } catch (err) {
        console.error('Error verifying certificate:', err);
      } finally {
        setIsLoading(false);
      }
    };

    loadCert();
  }, [code]);

  const formattedDate = cert
    ? new Date(cert.issued_at).toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
      })
    : '';

  return (
    <div className="flex flex-col min-h-screen bg-slate-50 dark:bg-slate-950">
      <TopNavbar
        title="Верификация сертификата"
        subtitle="Официальный реестр подлинности документов ED.Learn"
      />

      <main className="p-6 sm:p-10 max-w-4xl w-full mx-auto flex-1 flex flex-col justify-center">
        {isLoading ? (
          <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-12 text-center animate-pulse space-y-4 shadow-sm">
            <div className="w-16 h-16 bg-slate-200 dark:bg-slate-800 rounded-2xl mx-auto" />
            <div className="h-6 w-64 bg-slate-200 dark:bg-slate-800 rounded mx-auto" />
            <div className="h-4 w-48 bg-slate-200 dark:bg-slate-800 rounded mx-auto" />
          </div>
        ) : cert ? (
          <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl overflow-hidden shadow-xl space-y-8 p-6 sm:p-10">
            {/* Status Header */}
            <div className="flex flex-col sm:flex-row items-center gap-4 p-5 rounded-2xl bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-900/50">
              <div className="w-14 h-14 rounded-2xl bg-emerald-500 text-white flex items-center justify-center flex-shrink-0 shadow-lg shadow-emerald-500/20">
                <ShieldCheck size={32} />
              </div>
              <div className="text-center sm:text-left space-y-1">
                <div className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-emerald-100 dark:bg-emerald-900/80 text-emerald-800 dark:text-emerald-300 text-[10px] font-black uppercase tracking-wider">
                  Подлинность подтверждена
                </div>
                <h2 className="text-lg sm:text-xl font-black text-slate-900 dark:text-white">
                  Официальный сертификат платформы ED.Learn
                </h2>
                <p className="text-xs text-slate-600 dark:text-slate-400">
                  Документ с номером <span className="font-mono font-bold text-slate-900 dark:text-slate-100">{cert.certificate_code}</span> является действительным и внесен в единый реестр.
                </p>
              </div>
            </div>

            {/* Certificate Details Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="p-5 rounded-2xl bg-slate-50 dark:bg-slate-800/60 border border-slate-200/80 dark:border-slate-700/80 space-y-2">
                <div className="flex items-center gap-2 text-slate-400 text-xs font-semibold uppercase">
                  <User size={15} />
                  <span>Выпускник курса</span>
                </div>
                <div className="text-xl font-black text-slate-900 dark:text-white">
                  {cert.student_name}
                </div>
                <p className="text-xs text-slate-500">
                  Успешно завершил(а) все модули программы и сдал(а) итоговые тесты.
                </p>
              </div>

              <div className="p-5 rounded-2xl bg-slate-50 dark:bg-slate-800/60 border border-slate-200/80 dark:border-slate-700/80 space-y-2">
                <div className="flex items-center gap-2 text-indigo-500 text-xs font-semibold uppercase">
                  <BookOpen size={15} />
                  <span>Образовательная программа</span>
                </div>
                <div className="text-lg font-black text-slate-900 dark:text-white leading-snug">
                  {cert.course_title}
                </div>
                <div className="flex items-center gap-3 pt-1 text-xs font-semibold">
                  <span className="text-emerald-600 dark:text-emerald-400">
                    Итоговый балл: {cert.final_score}%
                  </span>
                  <span className="text-slate-400">•</span>
                  <span className="text-slate-500 flex items-center gap-1">
                    <Calendar size={13} />
                    {formattedDate}
                  </span>
                </div>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex flex-wrap items-center justify-between gap-4 pt-4 border-t border-slate-100 dark:border-slate-800">
              <button
                type="button"
                onClick={() => setIsModalOpen(true)}
                className="px-6 py-3 rounded-2xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-600 hover:to-amber-700 text-white font-bold text-xs shadow-lg shadow-amber-500/20 flex items-center gap-2 cursor-pointer transition-all active:scale-95"
              >
                <Award size={18} />
                <span>Полноразмерный диплом</span>
              </button>

              <div className="flex items-center gap-3">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(true)}
                  className="px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-semibold text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 flex items-center gap-2 transition-colors cursor-pointer"
                >
                  <Printer size={15} />
                  <span>Распечатать</span>
                </button>

                <Link
                  href="/courses"
                  className="px-4 py-2.5 rounded-xl bg-slate-900 dark:bg-white text-white dark:text-slate-900 text-xs font-bold hover:opacity-90 flex items-center gap-1.5 transition-opacity"
                >
                  <span>Каталог курсов</span>
                  <ArrowRight size={14} />
                </Link>
              </div>
            </div>

            {/* Full Diploma Modal */}
            <CertificateModal
              isOpen={isModalOpen}
              onClose={() => setIsModalOpen(false)}
              certificate={cert}
            />
          </div>
        ) : (
          <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-10 text-center space-y-6 shadow-sm">
            <div className="w-16 h-16 rounded-2xl bg-rose-50 dark:bg-rose-950/60 text-rose-600 dark:text-rose-400 flex items-center justify-center mx-auto">
              <AlertTriangle size={32} />
            </div>

            <div className="space-y-2 max-w-md mx-auto">
              <h2 className="text-xl font-black text-slate-900 dark:text-white">
                Сертификат не найден
              </h2>
              <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
                Документ с номером <span className="font-mono font-bold text-rose-600">{code}</span> не зарегистрирован в единой базе сертификатов ED.Learn или был аннулирован.
              </p>
            </div>

            {/* Search alternative */}
            <form
              onSubmit={(e) => {
                e.preventDefault();
                if (searchCode.trim()) {
                  window.location.href = `/certificates/${encodeURIComponent(searchCode.trim())}`;
                }
              }}
              className="max-w-md mx-auto flex gap-2"
            >
              <input
                type="text"
                placeholder="Введите номер: EDL-2026-..."
                value={searchCode}
                onChange={(e) => setSearchCode(e.target.value)}
                className="flex-1 px-4 py-2.5 text-xs bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
              <button
                type="submit"
                className="px-4 py-2.5 bg-indigo-600 text-white rounded-xl text-xs font-bold hover:bg-indigo-700 transition-colors flex items-center gap-1.5 cursor-pointer"
              >
                <Search size={14} />
                <span>Проверить</span>
              </button>
            </form>

            <div className="pt-4">
              <Link
                href="/courses"
                className="text-xs font-semibold text-indigo-600 dark:text-indigo-400 hover:underline inline-flex items-center gap-1"
              >
                <span>Перейти к каталогу образовательных программ</span>
                <ArrowRight size={13} />
              </Link>
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
