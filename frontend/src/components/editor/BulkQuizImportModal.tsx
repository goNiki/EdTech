'use client';

import React, { useState, useEffect, useRef } from 'react';
import {
  ParsedQuiz,
  parseQuizzesFromText,
  readTextFromFile,
  convertQuizzesToPuckBlocks,
  SAMPLE_QUIZ_TEMPLATES,
} from '@/lib/quiz-parser';
import {
  Zap,
  X,
  FileText,
  Upload,
  CheckCircle2,
  AlertTriangle,
  HelpCircle,
  Sparkles,
  Check,
  Circle,
  Layers,
  FileCode
} from 'lucide-react';

interface BulkQuizImportModalProps {
  isOpen: boolean;
  onClose: () => void;
  onImport: (blocks: any[]) => void;
}

export default function BulkQuizImportModal({
  isOpen,
  onClose,
  onImport,
}: BulkQuizImportModalProps) {
  const [inputText, setInputText] = useState(SAMPLE_QUIZ_TEMPLATES[0].text);
  const [quizzes, setQuizzes] = useState<ParsedQuiz[]>([]);
  const [isLoadingFile, setIsLoadingFile] = useState(false);
  const [showHelpTooltip, setShowHelpTooltip] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Debounced parsing
  useEffect(() => {
    const handler = setTimeout(() => {
      const parsed = parseQuizzesFromText(inputText);
      setQuizzes(parsed);
    }, 150);

    return () => clearTimeout(handler);
  }, [inputText]);

  if (!isOpen) return null;

  const validQuizzes = quizzes.filter((q) => q.warnings.length === 0 && q.options.length >= 2);
  const totalPoints = quizzes.reduce((sum, q) => sum + (q.points || 0), 0);

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setIsLoadingFile(true);
    try {
      const text = await readTextFromFile(file);
      setInputText(text);
    } catch (err: any) {
      alert('Ошибка чтения файла: ' + (err.message || 'Не удалось прочитать текст'));
    } finally {
      setIsLoadingFile(false);
      e.target.value = '';
    }
  };

  const handleApplyImport = () => {
    if (quizzes.length === 0) return;
    const blocks = convertQuizzesToPuckBlocks(quizzes);
    onImport(blocks);
    onClose();
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto bg-black/75 backdrop-blur-xs flex items-center justify-center p-3 sm:p-6">
      <div className="relative w-full max-w-6xl h-[90vh] bg-white dark:bg-slate-900 rounded-3xl shadow-2xl border border-slate-200 dark:border-slate-800 flex flex-col overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        
        {/* Header */}
        <div className="p-4 sm:px-6 bg-slate-50 dark:bg-slate-800/80 border-b border-slate-200 dark:border-slate-700 flex items-center justify-between gap-3">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-xl bg-amber-500 text-white flex items-center justify-center shadow-md shadow-amber-500/20">
              <Zap size={20} />
            </div>
            <div>
              <h2 className="text-sm sm:text-base font-black text-slate-900 dark:text-white flex items-center gap-2">
                <span>⚡ Быстрый импорт тестов</span>
                <span className="px-2 py-0.5 rounded-full text-[10px] font-bold bg-amber-100 dark:bg-amber-950/80 text-amber-700 dark:text-amber-400 border border-amber-300 dark:border-amber-800">
                  Bulk Parser
                </span>
              </h2>
              <p className="text-[11px] text-slate-500 dark:text-slate-400">
                Вставьте текст с вопросами из Word, PDF или блокнота — они мгновенно превратятся в интерактивные блоки урока.
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={() => setShowHelpTooltip(!showHelpTooltip)}
              className="p-2 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-200/50 dark:hover:bg-slate-800 transition-colors"
              title="Справка по форматам"
            >
              <HelpCircle size={18} />
            </button>
            <button
              type="button"
              onClick={onClose}
              className="p-2 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-200/50 dark:hover:bg-slate-800 transition-colors cursor-pointer"
            >
              <X size={18} />
            </button>
          </div>
        </div>

        {/* Format Help Drawer (Expandable) */}
        {showHelpTooltip && (
          <div className="p-4 bg-indigo-50/90 dark:bg-indigo-950/40 border-b border-indigo-200 dark:border-indigo-900 text-xs text-slate-700 dark:text-slate-300 grid grid-cols-1 md:grid-cols-3 gap-4 animate-in slide-in-from-top-2">
            <div className="p-3 bg-white dark:bg-slate-900 rounded-xl border border-indigo-100 dark:border-indigo-800 space-y-1">
              <span className="font-bold text-indigo-600 dark:text-indigo-400 block">Формат 1: Звездочка (*)</span>
              <pre className="text-[10px] font-mono text-slate-500 dark:text-slate-400 whitespace-pre-wrap">
{`1. Вопрос?
a) Вариант 1
b) Вариант 2 *
c) Вариант 3`}
              </pre>
            </div>

            <div className="p-3 bg-white dark:bg-slate-900 rounded-xl border border-indigo-100 dark:border-indigo-800 space-y-1">
              <span className="font-bold text-indigo-600 dark:text-indigo-400 block">Формат 2: Чекбоксы [x]</span>
              <pre className="text-[10px] font-mono text-slate-500 dark:text-slate-400 whitespace-pre-wrap">
{`2. Выберите верное:
[x] Вариант А
[ ] Вариант Б
[x] Вариант В`}
              </pre>
            </div>

            <div className="p-3 bg-white dark:bg-slate-900 rounded-xl border border-indigo-100 dark:border-indigo-800 space-y-1">
              <span className="font-bold text-indigo-600 dark:text-indigo-400 block">Формат 3: Aiken / Moodle</span>
              <pre className="text-[10px] font-mono text-slate-500 dark:text-slate-400 whitespace-pre-wrap">
{`3. Вопрос?
A. Вариант 1
B. Вариант 2
ANSWER: B`}
              </pre>
            </div>
          </div>
        )}

        {/* Split View Content Area */}
        <div className="flex-1 grid grid-cols-1 lg:grid-cols-2 divide-y lg:divide-y-0 lg:divide-x divide-slate-200 dark:divide-slate-800 overflow-hidden">
          
          {/* Left Column: Raw Text Input */}
          <div className="flex flex-col h-full bg-slate-50/50 dark:bg-slate-950/40 p-4 sm:p-5 space-y-3 overflow-hidden">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <span className="text-xs font-black uppercase tracking-wider text-slate-500">
                Исходный текст вопросов
              </span>

              <div className="flex items-center gap-2">
                {/* Insert sample templates */}
                <select
                  onChange={(e) => {
                    const found = SAMPLE_QUIZ_TEMPLATES.find((t) => t.name === e.target.value);
                    if (found) setInputText(found.text);
                  }}
                  className="px-2.5 py-1 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300 font-semibold cursor-pointer"
                  defaultValue=""
                >
                  <option value="" disabled>
                    Примеры шаблонов...
                  </option>
                  {SAMPLE_QUIZ_TEMPLATES.map((tmpl) => (
                    <option key={tmpl.name} value={tmpl.name}>
                      {tmpl.name}
                    </option>
                  ))}
                </select>

                {/* Upload from file button */}
                <button
                  type="button"
                  onClick={() => fileInputRef.current?.click()}
                  disabled={isLoadingFile}
                  className="px-2.5 py-1 text-xs font-semibold rounded-xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 flex items-center gap-1.5 transition-colors cursor-pointer"
                  title="Загрузить из файла Word (.docx), Markdown (.md) или блокнота (.txt)"
                >
                  <Upload size={13} />
                  <span>{isLoadingFile ? 'Чтение...' : 'Файл'}</span>
                </button>
                <input
                  ref={fileInputRef}
                  type="file"
                  accept=".docx,.txt,.md"
                  className="hidden"
                  onChange={handleFileUpload}
                />

                <button
                  type="button"
                  onClick={() => setInputText('')}
                  className="px-2 py-1 text-xs text-rose-500 hover:text-rose-700 font-medium"
                >
                  Очистить
                </button>
              </div>
            </div>

            <textarea
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              placeholder="Вставьте сюда текст ваших вопросов..."
              className="flex-1 w-full p-4 font-mono text-xs leading-relaxed rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 text-slate-800 dark:text-slate-200 focus:outline-none focus:ring-2 focus:ring-amber-500 resize-none shadow-inner"
            />
          </div>

          {/* Right Column: Live Interactive Preview */}
          <div className="flex flex-col h-full bg-white dark:bg-slate-900 p-4 sm:p-5 space-y-3 overflow-hidden">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <span className="text-xs font-black uppercase tracking-wider text-slate-500">
                  Интерактивный предпросмотр
                </span>
                <span className="px-2 py-0.5 rounded-full text-[10px] font-extrabold bg-indigo-100 dark:bg-indigo-950/80 text-indigo-700 dark:text-indigo-400">
                  {quizzes.length} вопросов ({totalPoints} баллов)
                </span>
              </div>

              {validQuizzes.length < quizzes.length && (
                <span className="text-[11px] text-amber-500 font-bold flex items-center gap-1">
                  <AlertTriangle size={13} />
                  {quizzes.length - validQuizzes.length} требуют проверки
                </span>
              )}
            </div>

            {/* Quizzes Scroll Area */}
            <div className="flex-1 overflow-y-auto space-y-3 pr-1">
              {quizzes.length === 0 ? (
                <div className="h-full flex flex-col items-center justify-center text-center p-8 text-slate-400 space-y-2">
                  <div className="w-12 h-12 rounded-2xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center">
                    <Sparkles size={22} className="text-amber-500" />
                  </div>
                  <h4 className="text-sm font-bold text-slate-700 dark:text-slate-300">
                    Ожидание текста вопросов
                  </h4>
                  <p className="text-xs max-w-sm">
                    Начните вводить или вставьте текст в окне слева. Парсер на лету распознает вопросы, определит правильные ответы и подсчитает баллы.
                  </p>
                </div>
              ) : (
                quizzes.map((quiz, qIdx) => (
                  <div
                    key={quiz.id}
                    className={`p-4 rounded-2xl border transition-all ${
                      quiz.warnings.length > 0
                        ? 'border-amber-300 bg-amber-50/20 dark:border-amber-800 dark:bg-amber-950/10'
                        : 'border-slate-200 dark:border-slate-800 bg-slate-50/40 dark:bg-slate-900/40'
                    }`}
                  >
                    {/* Header */}
                    <div className="flex items-center justify-between gap-2 mb-2">
                      <div className="flex items-center gap-2">
                        <span className="w-6 h-6 rounded-lg bg-indigo-600 text-white font-bold text-xs flex items-center justify-center">
                          {qIdx + 1}
                        </span>
                        <span
                          className={`px-2 py-0.5 rounded-lg text-[10px] font-extrabold uppercase ${
                            quiz.type === 'multi'
                              ? 'bg-purple-100 dark:bg-purple-950 text-purple-700 dark:text-purple-300'
                              : 'bg-blue-100 dark:bg-blue-950 text-blue-700 dark:text-blue-300'
                          }`}
                        >
                          {quiz.type === 'multi' ? 'Множественный выбор' : 'Одиночный выбор'}
                        </span>
                      </div>
                      <span className="text-[11px] font-bold text-slate-400">
                        {quiz.points} баллов
                      </span>
                    </div>

                    {/* Question text */}
                    <h4 className="text-xs font-bold text-slate-900 dark:text-white leading-relaxed mb-3">
                      {quiz.question}
                    </h4>

                    {/* Options list */}
                    <div className="space-y-1.5">
                      {quiz.options.map((opt, oIdx) => (
                        <div
                          key={oIdx}
                          className={`p-2 rounded-xl text-xs flex items-center gap-2.5 transition-colors ${
                            opt.isCorrect
                              ? 'bg-emerald-50 dark:bg-emerald-950/50 text-emerald-900 dark:text-emerald-200 font-bold border border-emerald-300 dark:border-emerald-800'
                              : 'bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300 border border-slate-200 dark:border-slate-700/60'
                          }`}
                        >
                          {opt.isCorrect ? (
                            <CheckCircle2 size={15} className="text-emerald-600 flex-shrink-0" />
                          ) : (
                            <Circle size={15} className="text-slate-300 dark:text-slate-600 flex-shrink-0" />
                          )}
                          <span className="truncate">{opt.text}</span>
                        </div>
                      ))}
                    </div>

                    {/* Warnings badge */}
                    {quiz.warnings.length > 0 && (
                      <div className="mt-3 p-2 rounded-xl bg-amber-100/70 dark:bg-amber-950/50 text-amber-800 dark:text-amber-300 text-[11px] font-medium flex items-center gap-1.5">
                        <AlertTriangle size={14} className="flex-shrink-0" />
                        <span>{quiz.warnings.join(' ')}</span>
                      </div>
                    )}
                  </div>
                ))
              )}
            </div>
          </div>
        </div>

        {/* Footer Actions */}
        <div className="p-4 sm:px-6 bg-slate-50 dark:bg-slate-800/80 border-t border-slate-200 dark:border-slate-700 flex items-center justify-between gap-3">
          <div className="text-xs text-slate-500">
            {quizzes.length > 0 ? (
              <span>
                Будет создано <strong>{quizzes.length}</strong> интерактивных блоков тестов
              </span>
            ) : (
              <span>Введите текст тестов для импорта</span>
            )}
          </div>

          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
            >
              Отмена
            </button>
            <button
              type="button"
              disabled={quizzes.length === 0}
              onClick={handleApplyImport}
              className="px-6 py-2.5 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-600 hover:to-amber-700 disabled:opacity-40 text-white text-xs font-black shadow-lg shadow-amber-500/20 flex items-center gap-2 transition-all cursor-pointer active:scale-95"
            >
              <Zap size={15} />
              <span>Добавить {quizzes.length} тестов в урок</span>
            </button>
          </div>
        </div>

      </div>
    </div>
  );
}
