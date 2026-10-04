'use client';

import React, { useState } from 'react';
import { CertificateData } from '@/lib/certificates';
import {
  Award,
  Printer,
  Copy,
  Check,
  ExternalLink,
  X,
  ShieldCheck,
  Sparkles
} from 'lucide-react';
import Link from 'next/link';

interface CertificateModalProps {
  isOpen: boolean;
  onClose: () => void;
  certificate: CertificateData;
}

export default function CertificateModal({
  isOpen,
  onClose,
  certificate,
}: CertificateModalProps) {
  const [copied, setCopied] = useState(false);

  if (!isOpen) return null;

  const formattedDate = new Date(certificate.issued_at).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });

  const verifyUrl = typeof window !== 'undefined'
    ? `${window.location.origin}/certificates/${certificate.certificate_code}`
    : `/certificates/${certificate.certificate_code}`;

  const handleCopyLink = async () => {
    try {
      await navigator.clipboard.writeText(verifyUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2500);
    } catch {
      // ignore
    }
  };

  const handlePrint = () => {
    window.print();
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto bg-black/70 backdrop-blur-xs flex items-center justify-center p-4 print:p-0 print:bg-white print:static">
      <div className="relative w-full max-w-4xl bg-white dark:bg-slate-900 rounded-3xl shadow-2xl border border-slate-200 dark:border-slate-800 overflow-hidden flex flex-col my-8 print:border-none print:shadow-none print:m-0 print:max-w-none">
        
        {/* Top Control Bar (Hidden on print) */}
        <div className="p-4 sm:px-6 bg-slate-50 dark:bg-slate-800/80 border-b border-slate-200 dark:border-slate-700 flex flex-wrap items-center justify-between gap-3 print:hidden">
          <div className="flex items-center gap-2">
            <span className="p-1.5 rounded-lg bg-amber-500/10 text-amber-500">
              <Award size={18} />
            </span>
            <span className="text-xs sm:text-sm font-extrabold text-slate-800 dark:text-slate-200">
              Официальный сертификат ED.Learn
            </span>
          </div>

          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={handleCopyLink}
              className="px-3 py-1.5 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-semibold text-slate-700 dark:text-slate-300 hover:bg-white dark:hover:bg-slate-800 flex items-center gap-1.5 transition-colors cursor-pointer"
            >
              {copied ? <Check size={14} className="text-emerald-500" /> : <Copy size={14} />}
              <span>{copied ? 'Ссылка скопирована' : 'Копировать ссылку'}</span>
            </button>

            <Link
              href={`/certificates/${certificate.certificate_code}`}
              target="_blank"
              className="px-3 py-1.5 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-semibold text-slate-700 dark:text-slate-300 hover:bg-white dark:hover:bg-slate-800 flex items-center gap-1.5 transition-colors"
            >
              <ExternalLink size={14} />
              <span className="hidden sm:inline">Верификация</span>
            </Link>

            <button
              type="button"
              onClick={handlePrint}
              className="px-3.5 py-1.5 rounded-xl bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold shadow-md shadow-indigo-600/20 flex items-center gap-1.5 transition-all cursor-pointer"
            >
              <Printer size={14} />
              <span>Печать / PDF</span>
            </button>

            <button
              type="button"
              onClick={onClose}
              className="p-1.5 rounded-xl text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 hover:bg-slate-200/50 dark:hover:bg-slate-800 transition-colors ml-2 cursor-pointer"
              aria-label="Закрыть"
            >
              <X size={18} />
            </button>
          </div>
        </div>

        {/* Certificate Printable Body */}
        <div className="p-6 sm:p-10 md:p-12 print:p-8 bg-gradient-to-br from-amber-50/40 via-white to-slate-50 dark:from-slate-900 dark:via-slate-950 dark:to-slate-900 flex justify-center items-center">
          {/* Certificate Container with Decorative Golden Borders */}
          <div className="w-full relative border-8 border-double border-amber-300 dark:border-amber-600/70 p-6 sm:p-10 rounded-2xl bg-white dark:bg-slate-900/95 shadow-inner space-y-8 text-center print:border-4 print:border-amber-500">
            
            {/* Top Ornamental Corners */}
            <div className="absolute top-2 left-2 text-amber-400 opacity-60 print:opacity-100">
              <Sparkles size={24} />
            </div>
            <div className="absolute top-2 right-2 text-amber-400 opacity-60 print:opacity-100">
              <Sparkles size={24} />
            </div>
            <div className="absolute bottom-2 left-2 text-amber-400 opacity-60 print:opacity-100">
              <Sparkles size={24} />
            </div>
            <div className="absolute bottom-2 right-2 text-amber-400 opacity-60 print:opacity-100">
              <Sparkles size={24} />
            </div>

            {/* Platform Brand Header */}
            <div className="space-y-2">
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-amber-50 dark:bg-amber-950/60 border border-amber-200 dark:border-amber-800 text-amber-700 dark:text-amber-400 text-xs font-black tracking-widest uppercase">
                <Award size={14} />
                <span>ED.Learn Academy</span>
              </div>
              <h1 className="text-2xl sm:text-3xl md:text-4xl font-black text-slate-900 dark:text-white uppercase tracking-wider font-serif">
                Сертификат об окончании
              </h1>
              <p className="text-xs sm:text-sm font-semibold text-slate-400 uppercase tracking-widest">
                Certificate of Course Completion
              </p>
            </div>

            {/* Recipient Notice */}
            <div className="space-y-3 max-w-xl mx-auto">
              <p className="text-xs sm:text-sm text-slate-500 dark:text-slate-400 italic">
                Настоящим подтверждается, что
              </p>
              <div className="text-2xl sm:text-4xl font-extrabold text-slate-900 dark:text-white tracking-tight pb-2 border-b-2 border-amber-300 dark:border-amber-600/60 inline-block px-6">
                {certificate.student_name}
              </div>
              <p className="text-xs sm:text-sm text-slate-500 dark:text-slate-400 italic pt-1">
                успешно освоил(а) образовательную программу курса:
              </p>
              <h3 className="text-lg sm:text-2xl font-black text-indigo-700 dark:text-indigo-400 leading-snug">
                «{certificate.course_title}»
              </h3>
            </div>

            {/* Performance and Details Grid */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-4 max-w-2xl mx-auto border-t border-slate-100 dark:border-slate-800">
              <div className="p-3 rounded-xl bg-slate-50 dark:bg-slate-800/50">
                <span className="text-[10px] uppercase font-bold text-slate-400 block mb-0.5">
                  Успеваемость
                </span>
                <span className="text-sm sm:text-base font-extrabold text-emerald-600 dark:text-emerald-400">
                  {certificate.final_score}% баллов
                </span>
              </div>

              <div className="p-3 rounded-xl bg-slate-50 dark:bg-slate-800/50">
                <span className="text-[10px] uppercase font-bold text-slate-400 block mb-0.5">
                  Дата выпуска
                </span>
                <span className="text-xs sm:text-sm font-bold text-slate-800 dark:text-slate-200">
                  {formattedDate}
                </span>
              </div>

              <div className="p-3 rounded-xl bg-slate-50 dark:bg-slate-800/50">
                <span className="text-[10px] uppercase font-bold text-slate-400 block mb-0.5">
                  ID Документа
                </span>
                <span className="text-xs font-mono font-bold text-slate-800 dark:text-slate-200">
                  {certificate.certificate_code}
                </span>
              </div>
            </div>

            {/* Signatures & Seal Block */}
            <div className="pt-6 flex flex-col sm:flex-row items-center justify-between gap-6 max-w-2xl mx-auto border-t border-slate-100 dark:border-slate-800">
              {/* Instructor Signature */}
              <div className="text-center sm:text-left space-y-1">
                <div className="h-8 flex items-end justify-center sm:justify-start">
                  <span className="font-serif italic text-base text-slate-600 dark:text-slate-300 font-bold border-b border-slate-400 pb-0.5 px-2">
                    ED.Learn Academic Council
                  </span>
                </div>
                <p className="text-[10px] text-slate-400 uppercase font-semibold">
                  Аттестационная комиссия
                </p>
              </div>

              {/* Gold Verification Seal */}
              <div className="w-20 h-20 rounded-full border-4 border-dashed border-amber-400 dark:border-amber-500 bg-amber-50 dark:bg-amber-950/60 flex flex-col items-center justify-center text-amber-600 dark:text-amber-400 shadow-sm">
                <ShieldCheck size={28} />
                <span className="text-[8px] font-black uppercase tracking-tighter mt-0.5">
                  VERIFIED
                </span>
              </div>

              {/* Verification Info */}
              <div className="text-center sm:text-right space-y-1">
                <p className="text-[10px] text-slate-400 font-semibold uppercase">
                  Проверка подлинности онлайн
                </p>
                <p className="text-[11px] font-mono text-indigo-600 dark:text-indigo-400 truncate max-w-xs">
                  {verifyUrl}
                </p>
              </div>
            </div>

          </div>
        </div>

        {/* Print Styles */}
        <style jsx global>{`
          @media print {
            body {
              background: white !important;
              color: black !important;
            }
            nav, header, footer {
              display: none !important;
            }
            @page {
              size: landscape;
              margin: 10mm;
            }
          }
        `}</style>
      </div>
    </div>
  );
}
