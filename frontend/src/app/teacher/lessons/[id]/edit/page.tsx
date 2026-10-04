'use client';

import { Puck, blocksPlugin, outlinePlugin } from '@puckeditor/core';
import '@puckeditor/core/dist/index.css';
import { config } from '@/lib/puck-config';
import CustomOutline from '@/components/editor/CustomOutline';
import { useEffect, useState, use } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { useAuth } from '@/store/useAuth';
import { ArrowLeft, Save, Sparkles, CheckCircle2, Lock, ShieldAlert, FileUp, Loader2, Zap } from 'lucide-react';
import { useRef } from 'react';
import { convertDocumentToHtml } from '@/lib/document-importer';
import BulkQuizImportModal from '@/components/editor/BulkQuizImportModal';

export default function LessonEditor({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const { user } = useAuth();
  const [initialData, setInitialData] = useState<any>(null);
  const [lessonMeta, setLessonMeta] = useState<any>(null);
  const [isReadOnly, setIsReadOnly] = useState(false);
  const [forbiddenAlert, setForbiddenAlert] = useState<string | null>(null);
  const [toastMsg, setToastMsg] = useState<string | null>(null);
  const [isImportingHeader, setIsImportingHeader] = useState(false);
  const [isBulkQuizModalOpen, setIsBulkQuizModalOpen] = useState(false);
  const headerFileInputRef = useRef<HTMLInputElement>(null);

  const showToast = (msg: string) => {
    setToastMsg(msg);
    setTimeout(() => setToastMsg(null), 3000);
  };

  const handleBulkQuizzesImport = (blocks: any[]) => {
    setInitialData((prev: any) => ({
      ...prev,
      content: [...(prev?.content || []), ...blocks],
    }));
    showToast(`Успешно добавлено ${blocks.length} блоков тестов в урок!`);
  };

  const handleHeaderFileImport = async (file: File) => {
    setIsImportingHeader(true);
    try {
      const res = await convertDocumentToHtml(file);
      const newBlock = {
        type: 'RichTextBlock',
        props: {
          id: `RichTextBlock-${Date.now()}`,
          title: res.extractedTitle || 'Импортированная лекция',
          contentHtml: res.html,
        },
      };
      setInitialData((prev: any) => ({
        ...prev,
        content: [...(prev?.content || []), newBlock],
      }));
      showToast(`Документ успешно добавлен в урок (${res.wordCount} слов, ${res.tablesCount} таблиц)`);
    } catch (err: any) {
      alert(err.message || 'Ошибка импорта документа');
    } finally {
      setIsImportingHeader(false);
    }
  };

  useEffect(() => {
    const fetchLesson = async () => {
      try {
        const response = await api.get(`/lessons/${id}`);
        const lesson = response.data?.data?.lesson || response.data?.lesson || response.data || {};
        setLessonMeta(lesson);

        if (lesson.can_edit === false || lesson.permissions?.can_edit === false) {
          setIsReadOnly(true);
        }
        if (user && user.role !== 'admin' && lesson.course_creator_id && user.id !== lesson.course_creator_id && user.role !== 'teacher') {
          setIsReadOnly(true);
        }

        if (lesson.content && lesson.content !== '{}' && lesson.content !== '') {
          try {
            setInitialData(JSON.parse(lesson.content));
          } catch (e) {
            setInitialData({});
          }
        } else {
          setInitialData({});
        }
      } catch (err) {
        console.warn('Could not fetch lesson content, initializing empty editor', err);
        setInitialData({});
        setLessonMeta({ title: 'Новый урок', description: 'Описание появится позже...' });
      }
    };
    fetchLesson();
  }, [id, user]);

  const handleSave = async (data: any) => {
    if (isReadOnly) {
      setForbiddenAlert('У вас нет прав на редактирование этого курса. Изменения не сохранены.');
      return;
    }
    try {
      const contentString = JSON.stringify(data);
      await api.patch(`/lessons/${id}`, {
        content: contentString,
      });
      showToast('Контент и интерактивные тесты успешно сохранены!');
    } catch (err: any) {
      console.error('Ошибка сохранения:', err);
      if (err.response?.status === 403) {
        setIsReadOnly(true);
        setForbiddenAlert('У вас нет прав на редактирование этого курса. Изменения не сохранены.');
      } else {
        alert(err.response?.data?.message || err.response?.data?.error || 'Ошибка при сохранении урока');
      }
    }
  };

  if (!initialData) {
    return (
      <div className="flex h-screen items-center justify-center bg-slate-900 text-white">
        <div className="animate-pulse text-sm font-bold flex items-center gap-3">
          <Sparkles className="text-indigo-400" />
          <span>Загрузка визуального редактора Puck...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="h-screen w-full flex flex-col bg-slate-950">
      {/* Top Custom Bar */}
      <header className="h-14 bg-slate-900 border-b border-slate-800 px-6 flex items-center justify-between text-white flex-shrink-0 z-50">
        <div className="flex items-center gap-4">
          <button
            onClick={() => router.back()}
            className="p-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 transition-colors flex items-center gap-1.5 text-xs font-bold"
          >
            <ArrowLeft size={16} />
            <span>Вернуться к курсу</span>
          </button>
          <div className="flex items-center gap-2">
            <span className="font-extrabold text-sm text-slate-200">
              {lessonMeta?.title || lessonMeta?.Title || 'Редактор урока'}
            </span>
            {isReadOnly ? (
              <span className="px-2.5 py-1 rounded text-[10px] font-extrabold uppercase tracking-wider bg-rose-950 text-rose-300 border border-rose-800 flex items-center gap-1.5 shadow-xs">
                <Lock size={12} />
                Режим просмотра (Только чтение)
              </span>
            ) : (
              <span className="px-2 py-0.5 rounded text-[9px] font-extrabold uppercase tracking-wider bg-indigo-950 text-indigo-400 border border-indigo-800">
                {lessonMeta?.status || lessonMeta?.Status || 'Черновик'}
              </span>
            )}
          </div>
        </div>

        <div className="flex items-center gap-3">
          <button
            type="button"
            disabled={isReadOnly || isImportingHeader}
            onClick={() => headerFileInputRef.current?.click()}
            className="px-3 py-1.5 rounded-xl bg-purple-600 hover:bg-purple-700 disabled:opacity-50 text-white text-xs font-bold transition-all flex items-center gap-1.5 shadow-md shadow-purple-600/20 cursor-pointer"
            title="Импортировать готовый документ Word (.docx) или Markdown (.md) как новый раздел лекции"
          >
            {isImportingHeader ? (
              <>
                <Loader2 size={14} className="animate-spin" />
                <span>Импорт документа...</span>
              </>
            ) : (
              <>
                <FileUp size={14} />
                <span>Вставить документ (.docx / .md)</span>
              </>
            )}
          </button>
          <input
            ref={headerFileInputRef}
            type="file"
            accept=".docx,.md,.markdown"
            className="hidden"
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (file) handleHeaderFileImport(file);
              e.target.value = '';
            }}
          />

          {/* Bulk Quiz Import Button */}
          <button
            type="button"
            disabled={isReadOnly}
            onClick={() => setIsBulkQuizModalOpen(true)}
            className="px-3.5 py-1.5 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-600 hover:to-amber-700 disabled:opacity-50 text-white text-xs font-black transition-all flex items-center gap-1.5 shadow-md shadow-amber-500/20 cursor-pointer active:scale-95"
            title="Быстрый импорт тестов из текста (Aiken, GIFT, звездочки, чекбоксы)"
          >
            <Zap size={14} className="fill-white" />
            <span>⚡ Быстрый импорт тестов</span>
          </button>

          <span className="text-xs text-slate-400 font-medium hidden sm:inline">
            Режим: Content-as-Data Visual Builder
          </span>
        </div>
      </header>

      {/* Puck Editor Container */}
      <div className="flex-1 w-full overflow-hidden">
        <Puck
          config={config}
          data={initialData}
          onPublish={handleSave}
          plugins={[
            blocksPlugin({ label: 'Блоки' }),
            outlinePlugin({ label: 'Структура' }),
          ]}
          overrides={{
            outline: CustomOutline,
          }}
        />
      </div>

      {/* Forbidden 403 Alert Modal */}
      {forbiddenAlert && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-xs p-4 animate-in fade-in duration-200">
          <div className="bg-slate-900 border border-rose-500/40 rounded-3xl p-6 max-w-md w-full shadow-2xl space-y-4 text-center">
            <div className="w-12 h-12 rounded-2xl bg-rose-500/10 text-rose-400 flex items-center justify-center mx-auto">
              <ShieldAlert size={24} />
            </div>
            <h3 className="text-base font-extrabold text-white">Доступ ограничен (403 Forbidden)</h3>
            <p className="text-xs text-slate-300 leading-relaxed">
              {forbiddenAlert}
            </p>
            <div className="flex gap-3 pt-2">
              <button
                onClick={() => setForbiddenAlert(null)}
                className="flex-1 py-2.5 bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-bold rounded-xl transition-colors"
              >
                Закрыть
              </button>
              <button
                onClick={() => router.back()}
                className="flex-1 py-2.5 bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold rounded-xl transition-colors"
              >
                Вернуться к курсу
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Toast Notification */}
      {toastMsg && (
        <div className="fixed bottom-6 right-6 z-50 bg-indigo-600 text-white text-xs font-bold px-5 py-3.5 rounded-2xl shadow-2xl flex items-center gap-2 animate-bounce">
          <CheckCircle2 size={16} />
          <span>{toastMsg}</span>
        </div>
      )}

      {/* Bulk Quiz Import Modal */}
      <BulkQuizImportModal
        isOpen={isBulkQuizModalOpen}
        onClose={() => setIsBulkQuizModalOpen(false)}
        onImport={handleBulkQuizzesImport}
      />
    </div>
  );
}
