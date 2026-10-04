'use client';

import React, { useState, useEffect, useRef } from 'react';
import {
  ParsedQuiz,
  QuizBlockType,
  parseQuizzesFromText,
  readTextFromFile,
  convertQuizzesToPuckBlocks,
  SAMPLE_QUIZ_TEMPLATES,
} from '@/lib/quiz-parser';
import {
  Zap,
  X,
  Upload,
  CheckCircle2,
  AlertTriangle,
  HelpCircle,
  Sparkles,
  Circle,
  Layers,
  ArrowRight,
  Split,
  ListOrdered,
  FileInput,
  FileCheck,
  Check,
  Plus
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
  const [selectedTemplateIndex, setSelectedTemplateIndex] = useState(0);
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

  const validQuizzes = quizzes.filter((q) => q.warnings.length === 0);
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

  const handleApplyTemplate = (append: boolean = false) => {
    const template = SAMPLE_QUIZ_TEMPLATES[selectedTemplateIndex];
    if (!template) return;

    if (append) {
      setInputText((prev) => (prev.trim() ? `${prev.trim()}\n\n${template.text}` : template.text));
    } else {
      setInputText(template.text);
    }
  };

  const handleApplyImport = () => {
    if (quizzes.length === 0) return;
    const blocks = convertQuizzesToPuckBlocks(quizzes);
    onImport(blocks);
    onClose();
  };

  const getTypeBadge = (type: QuizBlockType) => {
    switch (type) {
      case 'single':
        return { label: 'Single Choice', color: 'bg-indigo-600 text-white' };
      case 'multi':
        return { label: 'Multiple Choice', color: 'bg-purple-600 text-white' };
      case 'match':
        return { label: 'Matching Pairs', color: 'bg-cyan-600 text-white' };
      case 'dropdown':
        return { label: 'Dropdown Blanks', color: 'bg-emerald-600 text-white' };
      case 'input':
        return { label: 'Text Input', color: 'bg-blue-600 text-white' };
      case 'sequence':
        return { label: 'Sequence', color: 'bg-amber-600 text-white' };
      case 'essay':
        return { label: 'Essay', color: 'bg-rose-600 text-white' };
      case 'file':
        return { label: 'File Upload', color: 'bg-slate-700 text-white' };
    }
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto bg-black/75 backdrop-blur-xs flex items-center justify-center p-3 sm:p-6">
      <div className="relative w-full max-w-6xl h-[92vh] bg-white dark:bg-slate-900 rounded-3xl shadow-2xl border border-slate-200 dark:border-slate-800 flex flex-col overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        
        {/* Header */}
        <div className="p-4 sm:px-6 bg-slate-50 dark:bg-slate-800/80 border-b border-slate-200 dark:border-slate-700 flex items-center justify-between gap-3">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-xl bg-amber-500 text-white flex items-center justify-center shadow-md shadow-amber-500/20">
              <Zap size={20} />
            </div>
            <div>
              <h2 className="text-sm sm:text-base font-black text-slate-900 dark:text-white flex items-center gap-2">
                <span>⚡ Быстрый импорт тестов: Все 8 типов блоков Puck</span>
                <span className="px-2 py-0.5 rounded-full text-[10px] font-bold bg-amber-100 dark:bg-amber-950/80 text-amber-700 dark:text-amber-400 border border-amber-300 dark:border-amber-800">
                  Full-Spectrum v2
                </span>
              </h2>
              <p className="text-[11px] text-slate-500 dark:text-slate-400">
                Вставьте текст с вопросами из Word, PDF или блокнота — они автоматически распознаются в интерактивные блоки Puck.
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

        {/* Template Quick Selection Bar */}
        <div className="px-4 py-2 bg-slate-100/70 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700/80 flex flex-wrap items-center justify-between gap-2">
          <div className="flex items-center gap-1.5 overflow-x-auto py-1 max-w-full">
            <span className="text-[11px] font-bold text-slate-500 mr-1 flex-shrink-0">
              Шаблоны:
            </span>
            {SAMPLE_QUIZ_TEMPLATES.map((tmpl, idx) => (
              <button
                key={tmpl.name}
                type="button"
                onClick={() => {
                  setSelectedTemplateIndex(idx);
                  setInputText(tmpl.text);
                }}
                className={`px-2.5 py-1 text-[11px] font-bold rounded-lg transition-all flex-shrink-0 cursor-pointer ${
                  selectedTemplateIndex === idx
                    ? 'bg-amber-500 text-white shadow-xs'
                    : 'bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700 border border-slate-200/80 dark:border-slate-700'
                }`}
              >
                {tmpl.name}
              </button>
            ))}
          </div>

          <div className="flex items-center gap-1.5 flex-shrink-0">
            <button
              type="button"
              onClick={() => handleApplyTemplate(true)}
              className="px-2.5 py-1 text-[11px] font-bold bg-white dark:bg-slate-800 text-amber-600 dark:text-amber-400 border border-amber-300 dark:border-amber-800/60 rounded-lg hover:bg-amber-50 dark:hover:bg-amber-950/40 flex items-center gap-1 cursor-pointer"
              title="Добавить выбранный шаблон к уже введенному тексту"
            >
              <Plus size={12} />
              <span>Добавить в конец</span>
            </button>
          </div>
        </div>

        {/* Format Help Drawer (Expandable) */}
        {showHelpTooltip && (
          <div className="p-4 bg-indigo-50/90 dark:bg-indigo-950/40 border-b border-indigo-200 dark:border-indigo-900 text-xs text-slate-700 dark:text-slate-300 grid grid-cols-1 md:grid-cols-4 gap-3 animate-in slide-in-from-top-2 overflow-y-auto max-h-48">
            <div className="p-2.5 bg-white dark:bg-slate-900 rounded-xl border border-indigo-100 dark:border-indigo-800 space-y-1">
              <span className="font-bold text-indigo-600 dark:text-indigo-400 block text-[11px]">1 & 2. Выбор (Single/Multi)</span>
              <pre className="text-[10px] font-mono text-slate-500 dark:text-slate-400 whitespace-pre-wrap">
{`[SINGLE]
Вопрос?
A. Вариант 1
B. Вариант 2 *
[MULTI]
[x] Ответ 1
[ ] Ответ 2`}
              </pre>
            </div>

            <div className="p-2.5 bg-white dark:bg-slate-900 rounded-xl border border-indigo-100 dark:border-indigo-800 space-y-1">
              <span className="font-bold text-cyan-600 dark:text-cyan-400 block text-[11px]">3. Сопоставление пар</span>
              <pre className="text-[10px] font-mono text-slate-500 dark:text-slate-400 whitespace-pre-wrap">
{`[MATCH]
Сопоставьте:
Ключ 1 :: Определение 1
Ключ 2 :: Определение 2
POINTS: 15`}
              </pre>
            </div>

            <div className="p-2.5 bg-white dark:bg-slate-900 rounded-xl border border-indigo-100 dark:border-indigo-800 space-y-1">
              <span className="font-bold text-emerald-600 dark:text-emerald-400 block text-[11px]">4 & 5. Пропуски в тексте</span>
              <pre className="text-[10px] font-mono text-slate-500 dark:text-slate-400 whitespace-pre-wrap">
{`[DROPDOWN]
Текст {Ответ; Вариант1, Вариант2}
[INPUT]
Текст с полем [слово]`}
              </pre>
            </div>

            <div className="p-2.5 bg-white dark:bg-slate-900 rounded-xl border border-indigo-100 dark:border-indigo-800 space-y-1">
              <span className="font-bold text-amber-600 dark:text-amber-400 block text-[11px]">6, 7 & 8. Порядок, Эссе, Файл</span>
              <pre className="text-[10px] font-mono text-slate-500 dark:text-slate-400 whitespace-pre-wrap">
{`[SEQUENCE]
-> Шаг 1
-> Шаг 2
[ESSAY] РУБРИКА: ...
[FILE] ФОРМАТЫ: .zip`}
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
                Исходный текст тестов и заданий
              </span>

              <div className="flex items-center gap-2">
                {/* Upload from file button */}
                <button
                  type="button"
                  onClick={() => fileInputRef.current?.click()}
                  disabled={isLoadingFile}
                  className="px-2.5 py-1 text-xs font-semibold rounded-xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 flex items-center gap-1.5 transition-colors cursor-pointer"
                  title="Загрузить из файла Word (.docx), Markdown (.md) или блокнота (.txt)"
                >
                  <Upload size={13} />
                  <span>{isLoadingFile ? 'Чтение...' : 'Загрузить файл'}</span>
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
                  className="px-2 py-1 text-xs text-rose-500 hover:text-rose-700 font-medium cursor-pointer"
                >
                  Очистить
                </button>
              </div>
            </div>

            <textarea
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              placeholder="Вставьте сюда текст ваших заданий (поддерживаются все 8 типов блоков)..."
              className="flex-1 w-full p-4 font-mono text-xs leading-relaxed rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 text-slate-800 dark:text-slate-200 focus:outline-none focus:ring-2 focus:ring-amber-500 resize-none shadow-inner"
            />
          </div>

          {/* Right Column: Live Interactive Preview */}
          <div className="flex flex-col h-full bg-white dark:bg-slate-900 p-4 sm:p-5 space-y-3 overflow-hidden">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <span className="text-xs font-black uppercase tracking-wider text-slate-500">
                  Интерактивный предпросмотр блоков
                </span>
                <span className="px-2 py-0.5 rounded-full text-[10px] font-extrabold bg-indigo-100 dark:bg-indigo-950/80 text-indigo-700 dark:text-indigo-400">
                  {quizzes.length} заданий ({totalPoints} баллов)
                </span>
              </div>

              {validQuizzes.length < quizzes.length && (
                <span className="text-[11px] text-amber-500 font-bold flex items-center gap-1">
                  <AlertTriangle size={13} />
                  {quizzes.length - validQuizzes.length} требуют внимания
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
                    Ожидание текста заданий
                  </h4>
                  <p className="text-xs max-w-sm">
                    Начните вводить или вставьте текст в окне слева. Парсер на лету распознает все 8 типов блоков Puck, определит параметры и подсчитает баллы.
                  </p>
                </div>
              ) : (
                quizzes.map((quiz, qIdx) => {
                  const badge = getTypeBadge(quiz.type);

                  return (
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
                          <span className="w-6 h-6 rounded-lg bg-slate-800 text-white font-bold text-xs flex items-center justify-center">
                            {qIdx + 1}
                          </span>
                          <span
                            className={`px-2 py-0.5 rounded-lg text-[10px] font-extrabold uppercase shadow-xs ${badge.color}`}
                          >
                            {badge.label}
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

                      {/* Specific Block Previews */}
                      {/* 1 & 2. Single and Multi choice options */}
                      {(quiz.type === 'single' || quiz.type === 'multi') && quiz.options && (
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
                      )}

                      {/* 3. Matching pairs */}
                      {quiz.type === 'match' && quiz.pairs && (
                        <div className="space-y-1.5">
                          <div className="text-[11px] font-semibold text-cyan-600 dark:text-cyan-400 mb-1">
                            Пар для сопоставления: {quiz.pairs.length}
                          </div>
                          {quiz.pairs.map((p, pIdx) => (
                            <div
                              key={pIdx}
                              className="p-2 rounded-xl text-xs bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex items-center justify-between gap-2"
                            >
                              <span className="font-semibold text-slate-800 dark:text-slate-200 flex-1">
                                {p.left}
                              </span>
                              <ArrowRight size={13} className="text-slate-400 flex-shrink-0" />
                              <span className="text-cyan-700 dark:text-cyan-300 font-medium flex-1 text-right">
                                {p.right}
                              </span>
                            </div>
                          ))}
                        </div>
                      )}

                      {/* 4. Dropdown blanks */}
                      {quiz.type === 'dropdown' && (
                        <div className="p-3 bg-emerald-50/50 dark:bg-emerald-950/30 rounded-xl border border-emerald-200 dark:border-emerald-800/60 space-y-1.5">
                          <div className="text-[11px] font-bold text-emerald-700 dark:text-emerald-300">
                            Пропусков с селектом: {quiz.blanksCount || 0}
                          </div>
                          <div className="text-xs font-mono text-slate-700 dark:text-slate-300 whitespace-pre-wrap leading-relaxed">
                            {quiz.templateText}
                          </div>
                        </div>
                      )}

                      {/* 5. Input blanks */}
                      {quiz.type === 'input' && (
                        <div className="p-3 bg-blue-50/50 dark:bg-blue-950/30 rounded-xl border border-blue-200 dark:border-blue-800/60 space-y-1.5">
                          <div className="text-[11px] font-bold text-blue-700 dark:text-blue-300">
                            Полей ручного ввода: {quiz.blanksCount || 0}
                          </div>
                          <div className="text-xs font-mono text-slate-700 dark:text-slate-300 whitespace-pre-wrap leading-relaxed">
                            {quiz.templateText}
                          </div>
                        </div>
                      )}

                      {/* 6. Sequence items */}
                      {quiz.type === 'sequence' && quiz.items && (
                        <div className="space-y-1.5">
                          <div className="text-[11px] font-semibold text-amber-600 dark:text-amber-400 mb-1">
                            Шагов последовательности: {quiz.items.length}
                          </div>
                          {quiz.items.map((it, itIdx) => (
                            <div
                              key={itIdx}
                              className="p-2 rounded-xl text-xs bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex items-center gap-2"
                            >
                              <span className="w-5 h-5 rounded-full bg-amber-100 dark:bg-amber-950 text-amber-700 dark:text-amber-300 text-[10px] font-bold flex items-center justify-center flex-shrink-0">
                                {itIdx + 1}
                              </span>
                              <span className="text-slate-800 dark:text-slate-200">
                                {it.text}
                              </span>
                            </div>
                          ))}
                        </div>
                      )}

                      {/* 7. Essay */}
                      {quiz.type === 'essay' && (
                        <div className="space-y-2 text-xs">
                          {quiz.rubric && (
                            <div className="p-2.5 bg-rose-50/50 dark:bg-rose-950/30 rounded-xl border border-rose-200 dark:border-rose-900/60">
                              <span className="font-bold text-rose-700 dark:text-rose-300 block mb-0.5 text-[11px]">
                                Рубрика оценивания:
                              </span>
                              <span className="text-slate-700 dark:text-slate-300">
                                {quiz.rubric}
                              </span>
                            </div>
                          )}
                          {quiz.sampleAnswer && (
                            <div className="p-2.5 bg-slate-100/70 dark:bg-slate-800/60 rounded-xl">
                              <span className="font-bold text-slate-600 dark:text-slate-400 block mb-0.5 text-[11px]">
                                Эталонный ответ:
                              </span>
                              <span className="text-slate-700 dark:text-slate-300 italic">
                                {quiz.sampleAnswer}
                              </span>
                            </div>
                          )}
                        </div>
                      )}

                      {/* 8. File Upload */}
                      {quiz.type === 'file' && (
                        <div className="p-3 bg-slate-100/70 dark:bg-slate-800/60 rounded-xl space-y-2 text-xs">
                          {quiz.instructions && (
                            <div>
                              <span className="font-bold text-slate-600 dark:text-slate-400 block text-[11px]">
                                Инструкция студенту:
                              </span>
                              <span className="text-slate-700 dark:text-slate-300">
                                {quiz.instructions}
                              </span>
                            </div>
                          )}
                          <div className="flex items-center gap-4 text-[11px] font-medium text-slate-500 dark:text-slate-400 pt-1 border-t border-slate-200 dark:border-slate-700">
                            <span>Форматы: <strong>{quiz.allowedTypes}</strong></span>
                            <span>Лимит: <strong>{quiz.maxSizeMB} МБ</strong></span>
                          </div>
                        </div>
                      )}

                      {/* Warnings badge */}
                      {quiz.warnings.length > 0 && (
                        <div className="mt-3 p-2 rounded-xl bg-amber-100/70 dark:bg-amber-950/50 text-amber-800 dark:text-amber-300 text-[11px] font-medium flex items-center gap-1.5">
                          <AlertTriangle size={14} className="flex-shrink-0" />
                          <span>{quiz.warnings.join(' ')}</span>
                        </div>
                      )}
                    </div>
                  );
                })
              )}
            </div>
          </div>
        </div>

        {/* Footer Actions */}
        <div className="p-4 sm:px-6 bg-slate-50 dark:bg-slate-800/80 border-t border-slate-200 dark:border-slate-700 flex items-center justify-between gap-3">
          <div className="text-xs text-slate-500">
            {quizzes.length > 0 ? (
              <span>
                Будет создано <strong>{quizzes.length}</strong> интерактивных блоков тестов на сумму <strong>{totalPoints}</strong> баллов
              </span>
            ) : (
              <span>Введите текст тестов или выберите шаблон для импорта</span>
            )}
          </div>

          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
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
              <span>Добавить {quizzes.length} блоков в урок</span>
            </button>
          </div>
        </div>

      </div>
    </div>
  );
}
