'use client';

import React, { useState, useRef, useEffect } from 'react';
import {
  Presentation,
  Maximize2,
  Minimize2,
  Download,
  ExternalLink,
  ChevronLeft,
  ChevronRight,
  FileText,
  AlertCircle
} from 'lucide-react';

export interface PresentationViewerProps {
  title?: string;
  mode?: 'embed' | 'pdf';
  embedUrl?: string;
  pdfUrl?: string;
  aspectRatio?: '16:9' | '4:3';
  allowDownload?: boolean;
}

/**
 * Нормализация URL презентаций (Google Slides, Canva, SpeakerDeck)
 */
export function normalizePresentationUrl(url?: string): string {
  if (!url) return '';
  const trimmed = url.trim();

  // Google Slides: заменяем /edit или /pub на /embed
  if (trimmed.includes('docs.google.com/presentation')) {
    if (trimmed.includes('/pub')) {
      return trimmed.replace('/pub', '/embed');
    }
    if (trimmed.includes('/edit')) {
      return trimmed.split('/edit')[0] + '/embed?start=false&loop=false&delayms=3000';
    }
  }

  // Canva: добавляем view?embed если не указано
  if (trimmed.includes('canva.com/design') && !trimmed.includes('view?embed')) {
    const base = trimmed.split('?')[0];
    return `${base}/view?embed`;
  }

  return trimmed;
}

export default function PresentationViewer({
  title,
  mode = 'embed',
  embedUrl,
  pdfUrl,
  aspectRatio = '16:9',
  allowDownload = true,
}: PresentationViewerProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  const normalizedEmbed = normalizePresentationUrl(embedUrl);
  const activePdfUrl = pdfUrl ? `${pdfUrl}#page=${currentPage}&toolbar=0&navpanes=0&view=FitH` : '';

  const handleToggleFullscreen = async () => {
    if (!containerRef.current) return;

    if (!document.fullscreenElement) {
      try {
        await containerRef.current.requestFullscreen();
        setIsFullscreen(true);
      } catch (err) {
        console.warn('Fullscreen request denied', err);
      }
    } else {
      try {
        await document.exitFullscreen();
        setIsFullscreen(false);
      } catch (err) {
        console.warn('Exit fullscreen failed', err);
      }
    }
  };

  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(Boolean(document.fullscreenElement));
    };
    document.addEventListener('fullscreenchange', handleFullscreenChange);
    return () => {
      document.removeEventListener('fullscreenchange', handleFullscreenChange);
    };
  }, []);

  // Keyboard navigation for PDF slide flipping
  useEffect(() => {
    if (mode !== 'pdf') return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'ArrowRight' || e.key === 'PageDown') {
        setCurrentPage((p) => p + 1);
      } else if (e.key === 'ArrowLeft' || e.key === 'PageUp') {
        setCurrentPage((p) => Math.max(1, p - 1));
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [mode]);

  const aspectClass = aspectRatio === '4:3' ? 'aspect-[4/3]' : 'aspect-video';

  return (
    <div
      ref={containerRef}
      className={`relative my-6 rounded-3xl overflow-hidden border border-slate-200 dark:border-slate-800 bg-slate-900 shadow-xl flex flex-col ${
        isFullscreen ? 'h-screen w-screen p-0 rounded-none' : ''
      }`}
    >
      {/* Presentation Top Control Toolbar */}
      <div className="px-5 py-3 bg-slate-950/90 border-b border-slate-800 flex items-center justify-between text-white z-20">
        <div className="flex items-center gap-2.5 min-w-0">
          <span className="p-1.5 rounded-xl bg-amber-500/20 text-amber-400">
            <Presentation size={16} />
          </span>
          <span className="text-xs sm:text-sm font-extrabold truncate text-slate-100">
            {title || 'Интерактивная презентация к лекции'}
          </span>
          <span className="px-2 py-0.5 rounded-md text-[10px] font-bold uppercase tracking-wider bg-slate-800 text-slate-300 hidden sm:inline">
            {mode === 'pdf' ? 'PDF Слайды' : 'Embed'}
          </span>
        </div>

        <div className="flex items-center gap-2">
          {mode === 'pdf' && (
            <div className="flex items-center gap-1 bg-slate-800/80 rounded-xl p-1 text-xs">
              <button
                type="button"
                onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                disabled={currentPage <= 1}
                className="p-1 rounded-lg hover:bg-slate-700 disabled:opacity-30 transition-colors"
                title="Предыдущий слайд (←)"
              >
                <ChevronLeft size={16} />
              </button>
              <span className="px-2 font-mono font-bold text-[11px] text-amber-400">
                Слайд {currentPage}
              </span>
              <button
                type="button"
                onClick={() => setCurrentPage((p) => p + 1)}
                className="p-1 rounded-lg hover:bg-slate-700 transition-colors"
                title="Следующий слайд (→)"
              >
                <ChevronRight size={16} />
              </button>
            </div>
          )}

          {allowDownload && (pdfUrl || embedUrl) && (
            <a
              href={pdfUrl || embedUrl}
              target="_blank"
              rel="noopener noreferrer"
              download={Boolean(pdfUrl)}
              className="p-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 transition-colors flex items-center gap-1.5 text-xs font-semibold"
              title="Открыть в новой вкладке / Скачать"
            >
              {pdfUrl ? <Download size={15} /> : <ExternalLink size={15} />}
              <span className="hidden md:inline">{pdfUrl ? 'Скачать PDF' : 'Открыть'}</span>
            </a>
          )}

          <button
            type="button"
            onClick={handleToggleFullscreen}
            className="p-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 transition-colors"
            title={isFullscreen ? 'Выйти из полноэкранного режима' : 'Во весь экран'}
          >
            {isFullscreen ? <Minimize2 size={16} /> : <Maximize2 size={16} />}
          </button>
        </div>
      </div>

      {/* Main Slides Display Area */}
      <div className={`w-full relative flex-1 flex items-center justify-center bg-slate-950 ${aspectClass}`}>
        {mode === 'pdf' ? (
          pdfUrl ? (
            <iframe
              src={activePdfUrl}
              title={title || 'Презентация PDF'}
              className="w-full h-full border-0"
            />
          ) : (
            <div className="p-8 text-center text-slate-400 space-y-2">
              <FileText size={32} className="mx-auto text-amber-500 opacity-60" />
              <p className="text-xs font-bold text-slate-300">Файл PDF еще не загружен</p>
              <p className="text-[11px] text-slate-500">
                Преподаватель может прикрепить файл презентации в редакторе курса.
              </p>
            </div>
          )
        ) : normalizedEmbed ? (
          <iframe
            src={normalizedEmbed}
            title={title || 'Встроенная презентация'}
            className="w-full h-full border-0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            allowFullScreen
          />
        ) : (
          <div className="p-8 text-center text-slate-400 space-y-2">
            <Presentation size={32} className="mx-auto text-amber-500 opacity-60" />
            <p className="text-xs font-bold text-slate-300">Ссылка на презентацию не указана</p>
            <p className="text-[11px] text-slate-500">
              Укажите ссылку на Google Slides, Canva или SpeakerDeck.
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
