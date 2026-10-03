'use client';

import React, { useState } from 'react';
import { parseSmartDropdownTemplate } from '@/lib/puck-config';
import {
  CheckCircle,
  XCircle,
  Upload,
  Send,
  HelpCircle,
  FileText,
  Play,
  ArrowUp,
  ArrowDown,
  Sparkles,
  Check,
  Lightbulb,
  ArrowRight,
  Award,
  BookOpen,
  X
} from 'lucide-react';

export interface LessonCompletionPayload {
  score: number;
  essays: Array<{
    question_text: string;
    answer_text: string;
    max_points: number;
  }>;
}

interface PuckLessonViewerProps {
  contentJson: string | object;
  onComplete?: (payload?: LessonCompletionPayload) => void;
  onNavigateBack?: () => void;
}

export default function PuckLessonViewer({
  contentJson,
  onComplete,
  onNavigateBack,
}: PuckLessonViewerProps) {
  let parsedContent: any = { content: [] };
  try {
    if (typeof contentJson === 'string') {
      parsedContent = JSON.parse(contentJson);
    } else if (typeof contentJson === 'object' && contentJson !== null) {
      parsedContent = contentJson;
    }
  } catch (e) {
    console.error('Failed to parse Puck content JSON', e);
  }

  const blocks: any[] = parsedContent.content || [];

  // Local interactive solver state
  const [singleAnswers, setSingleAnswers] = useState<Record<string, number>>({});
  const [multiAnswers, setMultiAnswers] = useState<Record<string, number[]>>({});
  const [dropdownAnswers, setDropdownAnswers] = useState<Record<string, Record<string, string>>>({});
  const [inputAnswers, setInputAnswers] = useState<Record<string, Record<string, string>>>({});
  const [matchAnswers, setMatchAnswers] = useState<Record<string, Record<number, string>>>({});
  const [essayAnswers, setEssayAnswers] = useState<Record<string, string>>({});
  const [sequenceOrders, setSequenceOrders] = useState<Record<string, string[]>>({});
  const [submittedBlocks, setSubmittedBlocks] = useState<Record<string, boolean>>({});
  const [toastMsg, setToastMsg] = useState<string | null>(null);
  const [completionResult, setCompletionResult] = useState<{
    score: number;
    earnedPoints: number;
    totalMaxPoints: number;
    hasQuizzes: boolean;
    essaysCount: number;
  } | null>(null);

  const showToast = (msg: string) => {
    setToastMsg(msg);
    setTimeout(() => setToastMsg(null), 3000);
  };

  const handleSingleSelect = (blockId: string, optIdx: number) => {
    if (submittedBlocks[blockId]) return;
    setSingleAnswers((prev) => ({ ...prev, [blockId]: optIdx }));
  };

  const handleMultiToggle = (blockId: string, optIdx: number) => {
    if (submittedBlocks[blockId]) return;
    setMultiAnswers((prev) => {
      const current = prev[blockId] || [];
      const updated = current.includes(optIdx)
        ? current.filter((i) => i !== optIdx)
        : [...current, optIdx];
      return { ...prev, [blockId]: updated };
    });
  };

  const handleMatchSelect = (blockId: string, pairIdx: number, selectedValue: string) => {
    if (submittedBlocks[blockId]) return;
    setMatchAnswers((prev) => ({
      ...prev,
      [blockId]: {
        ...(prev[blockId] || {}),
        [pairIdx]: selectedValue,
      },
    }));
  };

  const submitTest = (blockId: string) => {
    setSubmittedBlocks((prev) => ({ ...prev, [blockId]: true }));
    showToast('Ответ сохранен!');
  };

  const calculateCompletionPayload = (): {
    payload: LessonCompletionPayload;
    earnedPoints: number;
    totalMaxPoints: number;
  } => {
    let totalMaxPoints = 0;
    let earnedPoints = 0;
    const submittedEssays: LessonCompletionPayload['essays'] = [];

    blocks.forEach((block: any, idx: number) => {
      const { type, props = {} } = block;
      const blockId = props.id || `block-${idx}`;

      switch (type) {
        case 'QuizSingleBlock': {
          const pts = Number(props.points) || 10;
          totalMaxPoints += pts;
          const selected = singleAnswers[blockId];
          const correctIdx = props.options?.findIndex(
            (o: any) => o.isCorrect === true || o.isCorrect === 'true'
          );
          if (selected !== undefined && selected === correctIdx) {
            earnedPoints += pts;
          }
          break;
        }

        case 'QuizMultiBlock': {
          const pts = Number(props.points) || 15;
          totalMaxPoints += pts;
          const selectedList = multiAnswers[blockId] || [];
          const correctIndices = (props.options || [])
            .map((o: any, i: number) => (o.isCorrect === true || o.isCorrect === 'true' ? i : -1))
            .filter((i: number) => i !== -1);

          const isFullyCorrect =
            correctIndices.length > 0 &&
            selectedList.length === correctIndices.length &&
            correctIndices.every((i: number) => selectedList.includes(i));

          if (isFullyCorrect) {
            earnedPoints += pts;
          }
          break;
        }

        case 'QuizMatchBlock': {
          const pts = Number(props.points) || 20;
          totalMaxPoints += pts;
          const userMatches = matchAnswers[blockId] || {};
          const pairs: any[] = props.pairs || [];
          if (pairs.length > 0) {
            let correctPairs = 0;
            pairs.forEach((p, pIdx) => {
              if ((userMatches[pIdx] || '').trim().toLowerCase() === (p.right || '').trim().toLowerCase()) {
                correctPairs++;
              }
            });
            earnedPoints += (correctPairs / pairs.length) * pts;
          }
          break;
        }

        case 'QuizDropdownBlankBlock': {
          const pts = Number(props.points) || 10;
          totalMaxPoints += pts;
          const userVals = dropdownAnswers[blockId] || {};
          const fullText = props.templateText || '';
          const { blanks } = parseSmartDropdownTemplate(fullText, props.blanks);
          if (blanks.length > 0) {
            let correctBlanks = 0;
            blanks.forEach((b) => {
              if ((userVals[b.key] || '').trim().toLowerCase() === (b.correctAnswer || '').trim().toLowerCase()) {
                correctBlanks++;
              }
            });
            earnedPoints += (correctBlanks / blanks.length) * pts;
          }
          break;
        }

        case 'QuizInputBlankBlock': {
          const pts = Number(props.points) || 10;
          totalMaxPoints += pts;
          const userVals = inputAnswers[blockId] || {};
          const fullText = props.templateText || '';
          const { blanks } = parseSmartDropdownTemplate(fullText, props.blanks);
          if (blanks.length > 0) {
            let correctBlanks = 0;
            blanks.forEach((b) => {
              if ((userVals[b.key] || '').trim().toLowerCase() === (b.correctAnswer || '').trim().toLowerCase()) {
                correctBlanks++;
              }
            });
            earnedPoints += (correctBlanks / blanks.length) * pts;
          }
          break;
        }

        case 'QuizEssayBlock': {
          const answer = (essayAnswers[blockId] || '').trim();
          if (answer) {
            submittedEssays.push({
              question_text: props.question || 'Развернутый ответ на задание',
              answer_text: answer,
              max_points: Number(props.points) || 25,
            });
          }
          break;
        }

        case 'FileUploadBlock': {
          submittedEssays.push({
            question_text: props.title || 'Загрузка практической работы',
            answer_text: 'Документ прикреплен к практическому заданию',
            max_points: Number(props.points) || 50,
          });
          break;
        }
      }
    });

    const finalScore = totalMaxPoints > 0 ? Math.round((earnedPoints / totalMaxPoints) * 100) : 100;
    return {
      payload: {
        score: finalScore,
        essays: submittedEssays,
      },
      earnedPoints: Math.round(earnedPoints),
      totalMaxPoints: Math.round(totalMaxPoints),
    };
  };

  const handleFinish = () => {
    const { payload, earnedPoints, totalMaxPoints } = calculateCompletionPayload();
    setCompletionResult({
      score: payload.score,
      earnedPoints,
      totalMaxPoints,
      hasQuizzes: totalMaxPoints > 0,
      essaysCount: payload.essays.length,
    });
    if (onComplete) {
      onComplete(payload);
    }
  };

  const moveSequenceItem = (blockId: string, fromIdx: number, toIdx: number, defaultItems: any[]) => {
    const current = sequenceOrders[blockId] || defaultItems.map((it) => it.text);
    if (toIdx < 0 || toIdx >= current.length) return;
    const copy = [...current];
    const [moved] = copy.splice(fromIdx, 1);
    copy.splice(toIdx, 0, moved);
    setSequenceOrders((prev) => ({ ...prev, [blockId]: copy }));
  };

  if (!blocks || blocks.length === 0) {
    return (
      <div className="p-12 text-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 space-y-2">
        <FileText size={32} className="mx-auto text-slate-400" />
        <h4 className="text-sm font-bold text-slate-700 dark:text-slate-300">
          Содержимое урока пока формируется
        </h4>
        <p className="text-xs text-slate-500">
          Преподаватель еще не добавил учебные материалы в этот урок.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-8 max-w-4xl mx-auto">
      {blocks.map((block: any, idx: number) => {
        const { type, props = {} } = block;
        const blockId = props.id || `block-${idx}`;
        const isSubmitted = submittedBlocks[blockId];

        switch (type) {
          case 'HeaderBlock': {
            const HeadingTag = props.level || 'h1';
            const sizeClasses =
              props.level === 'h1'
                ? 'text-3xl lg:text-4xl font-extrabold text-slate-900 dark:text-white'
                : props.level === 'h2'
                ? 'text-2xl font-bold text-slate-800 dark:text-slate-100'
                : 'text-xl font-bold text-slate-800 dark:text-slate-200';

            return (
              <div key={blockId} className="my-6 space-y-2">
                <HeadingTag className={sizeClasses}>{props.title}</HeadingTag>
                {props.subtitle && (
                  <p className="text-sm text-slate-500 dark:text-slate-400 leading-relaxed max-w-3xl">
                    {props.subtitle}
                  </p>
                )}
              </div>
            );
          }

          case 'TextBlock': {
            if (props.contentHtml) {
              return (
                <div
                  key={blockId}
                  className="my-4 text-sm md:text-base text-slate-700 dark:text-slate-300 leading-relaxed prose dark:prose-invert max-w-none"
                  dangerouslySetInnerHTML={{ __html: props.contentHtml }}
                />
              );
            }
            return (
              <div
                key={blockId}
                className="my-4 text-sm md:text-base text-slate-700 dark:text-slate-300 leading-relaxed whitespace-pre-wrap"
              >
                {props.content}
              </div>
            );
          }

          case 'RichTextBlock': {
            return (
              <div key={blockId} className="my-6 space-y-3">
                {props.title && (
                  <h2 className="text-xl md:text-2xl font-black text-slate-900 dark:text-white border-b border-slate-200 dark:border-slate-800 pb-2">
                    {props.title}
                  </h2>
                )}
                <div
                  className="prose dark:prose-invert max-w-none text-slate-800 dark:text-slate-200 text-sm md:text-base leading-relaxed"
                  dangerouslySetInnerHTML={{ __html: props.contentHtml || '' }}
                />
              </div>
            );
          }

          case 'VideoBlock': {
            return (
              <div key={blockId} className="my-6 space-y-2">
                <div className="relative aspect-video rounded-3xl overflow-hidden bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-md">
                  <iframe
                    src={props.url}
                    title={props.caption || 'Видео'}
                    className="w-full h-full border-0"
                    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                    allowFullScreen
                  />
                </div>
                {props.caption && (
                  <p className="text-xs text-center text-slate-500 dark:text-slate-400 font-medium">
                    {props.caption}
                  </p>
                )}
              </div>
            );
          }

          case 'QuizSingleBlock': {
            const selectedOpt = singleAnswers[blockId];
            return (
              <div
                key={blockId}
                className="p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-4 shadow-sm"
              >
                <div className="flex justify-between items-center">
                  <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-indigo-600 text-white">
                    Один ответ
                  </span>
                  <span className="text-xs font-bold text-indigo-600 dark:text-indigo-400">
                    {props.points || 10} баллов
                  </span>
                </div>
                <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
                  {props.question}
                </h4>
                <div className="space-y-2">
                  {props.options?.map((opt: any, optIdx: number) => {
                    const isSelected = selectedOpt === optIdx;
                    const isCorrect = opt.isCorrect === 'true' || opt.isCorrect === true;

                    let btnClass = 'border-slate-200 dark:border-slate-800 hover:border-indigo-400';
                    if (isSelected && !isSubmitted) {
                      btnClass = 'border-indigo-600 bg-indigo-50/50 dark:bg-indigo-950/40 ring-1 ring-indigo-500';
                    } else if (isSubmitted) {
                      if (isCorrect) {
                        btnClass = 'border-emerald-500 bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300 ring-1 ring-emerald-500/20';
                      } else if (isSelected && !isCorrect) {
                        btnClass = 'border-rose-500 bg-rose-50 dark:bg-rose-950/40 text-rose-700 dark:text-rose-300 ring-1 ring-rose-500/20';
                      }
                    }

                    return (
                      <div key={optIdx} className="space-y-1.5">
                        <div
                          onClick={() => handleSingleSelect(blockId, optIdx)}
                          className={`p-3.5 border rounded-2xl flex items-center justify-between cursor-pointer transition-all text-xs font-medium ${btnClass}`}
                        >
                          <div className="flex items-center gap-3">
                            <div
                              className={`w-4 h-4 rounded-full border flex items-center justify-center ${
                                isSelected ? 'border-indigo-600 bg-indigo-600' : 'border-slate-300 dark:border-slate-700'
                              }`}
                            >
                              {isSelected && <span className="w-1.5 h-1.5 rounded-full bg-white" />}
                            </div>
                            <span>{opt.text}</span>
                          </div>
                          {isSubmitted && (
                            <span>
                              {isCorrect ? (
                                <CheckCircle size={16} className="text-emerald-600" />
                              ) : isSelected ? (
                                <XCircle size={16} className="text-rose-600" />
                              ) : null}
                            </span>
                          )}
                        </div>

                        {/* Hint / Explanation */}
                        {isSubmitted && opt.explain && (isSelected || isCorrect) && (
                          <div className="p-2.5 rounded-xl bg-amber-50 dark:bg-amber-950/40 border border-amber-200 dark:border-amber-900/50 text-[11px] text-amber-800 dark:text-amber-300 flex items-start gap-2">
                            <Lightbulb size={14} className="text-amber-600 flex-shrink-0 mt-0.5" />
                            <span>
                              <strong>Пояснение:</strong> {opt.explain}
                            </span>
                          </div>
                        )}
                      </div>
                    );
                  })}
                </div>

                {!isSubmitted && (
                  <button
                    disabled={selectedOpt === undefined}
                    onClick={() => submitTest(blockId)}
                    className="w-full py-2.5 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-40 text-white font-bold text-xs rounded-xl transition-all shadow-md mt-2 cursor-pointer"
                  >
                    Ответить
                  </button>
                )}
              </div>
            );
          }

          case 'QuizMultiBlock': {
            const selectedList = multiAnswers[blockId] || [];
            return (
              <div
                key={blockId}
                className="p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-4 shadow-sm"
              >
                <div className="flex justify-between items-center">
                  <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-purple-600 text-white">
                    Множественный выбор
                  </span>
                  <span className="text-xs font-bold text-purple-600 dark:text-purple-400">
                    {props.points || 15} баллов
                  </span>
                </div>
                <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
                  {props.question}
                </h4>
                <div className="space-y-2">
                  {props.options?.map((opt: any, optIdx: number) => {
                    const isSelected = selectedList.includes(optIdx);
                    const isCorrect = opt.isCorrect === 'true' || opt.isCorrect === true;

                    let btnClass = 'border-slate-200 dark:border-slate-800 hover:border-purple-400';
                    if (isSelected && !isSubmitted) {
                      btnClass = 'border-purple-600 bg-purple-50/50 dark:bg-purple-950/40 ring-1 ring-purple-500';
                    } else if (isSubmitted) {
                      if (isCorrect) {
                        btnClass = 'border-emerald-500 bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300';
                      } else if (isSelected && !isCorrect) {
                        btnClass = 'border-rose-500 bg-rose-50 dark:bg-rose-950/40 text-rose-700 dark:text-rose-300';
                      }
                    }

                    return (
                      <div
                        key={optIdx}
                        onClick={() => handleMultiToggle(blockId, optIdx)}
                        className={`p-3.5 border rounded-2xl flex items-center justify-between cursor-pointer transition-all text-xs font-medium ${btnClass}`}
                      >
                        <div className="flex items-center gap-3">
                          <div
                            className={`w-4 h-4 rounded-md border flex items-center justify-center ${
                              isSelected ? 'border-purple-600 bg-purple-600' : 'border-slate-300 dark:border-slate-700'
                            }`}
                          >
                            {isSelected && <span className="text-white text-[10px]">✓</span>}
                          </div>
                          <span>{opt.text}</span>
                        </div>
                        {isSubmitted && isCorrect && <CheckCircle size={16} className="text-emerald-600" />}
                      </div>
                    );
                  })}
                </div>

                {!isSubmitted && (
                  <button
                    disabled={selectedList.length === 0}
                    onClick={() => submitTest(blockId)}
                    className="w-full py-2.5 bg-purple-600 hover:bg-purple-700 disabled:opacity-40 text-white font-bold text-xs rounded-xl transition-all shadow-md mt-2 cursor-pointer"
                  >
                    Подтвердить выбор
                  </button>
                )}
              </div>
            );
          }

          case 'QuizMatchBlock': {
            const blockMatches = matchAnswers[blockId] || {};
            const pairs: Array<{ left: string; right: string }> = props.pairs || [];
            const allRightOptions = Array.from(new Set(pairs.map((p) => p.right))).filter(Boolean);

            return (
              <div
                key={blockId}
                className="p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-4 shadow-sm"
              >
                <div className="flex justify-between items-center">
                  <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-cyan-600 text-white">
                    Сопоставление элементов
                  </span>
                  <span className="text-xs font-bold text-cyan-600 dark:text-cyan-400">
                    {props.points || 20} баллов
                  </span>
                </div>
                <h4 className="text-sm font-bold text-slate-900 dark:text-white">{props.question}</h4>

                <div className="space-y-3">
                  {pairs.map((p, pIdx) => {
                    const userSelected = blockMatches[pIdx] || '';
                    const isPairCorrect = userSelected.trim().toLowerCase() === p.right.trim().toLowerCase();

                    let borderClass = 'border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/60';
                    if (isSubmitted) {
                      borderClass = isPairCorrect
                        ? 'border-emerald-500 bg-emerald-50/60 dark:bg-emerald-950/40 text-emerald-800 dark:text-emerald-200'
                        : 'border-rose-500 bg-rose-50/60 dark:bg-rose-950/40 text-rose-800 dark:text-rose-200';
                    }

                    return (
                      <div
                        key={pIdx}
                        className={`p-4 rounded-2xl border flex flex-col md:flex-row md:items-center justify-between gap-3 transition-all ${borderClass}`}
                      >
                        <div className="font-bold text-xs text-slate-900 dark:text-white md:w-1/2 flex items-center gap-2">
                          <span className="w-5 h-5 rounded-full bg-cyan-100 dark:bg-cyan-950 text-cyan-700 dark:text-cyan-300 text-[10px] font-black flex items-center justify-center flex-shrink-0">
                            {pIdx + 1}
                          </span>
                          <span>{p.left}</span>
                        </div>

                        <div className="flex items-center gap-2 md:w-1/2">
                          <ArrowRight size={14} className="text-slate-400 flex-shrink-0" />
                          <div className="flex-1 space-y-1">
                            <select
                              disabled={isSubmitted}
                              value={userSelected}
                              onChange={(e) => handleMatchSelect(blockId, pIdx, e.target.value)}
                              className="w-full px-3 py-2 bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 rounded-xl text-xs font-medium focus:outline-none focus:border-cyan-500 text-slate-800 dark:text-slate-200 cursor-pointer"
                            >
                              <option value="">[ Выберите соответствие ]</option>
                              {allRightOptions.map((opt, oIdx) => (
                                <option key={oIdx} value={opt}>
                                  {opt}
                                </option>
                              ))}
                            </select>

                            {isSubmitted && !isPairCorrect && (
                              <div className="text-[10px] text-rose-500 font-bold">
                                Правильно: {p.right}
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>

                {!isSubmitted && (
                  <button
                    onClick={() => submitTest(blockId)}
                    className="w-full py-2.5 bg-cyan-600 hover:bg-cyan-700 text-white font-bold text-xs rounded-xl shadow-md mt-2 cursor-pointer"
                  >
                    Завершить сопоставление
                  </button>
                )}
              </div>
            );
          }

          case 'QuizDropdownBlankBlock': {
            const blockAnswers = dropdownAnswers[blockId] || {};
            const fullText =
              props.templateText ||
              `${props.textBefore || ''} {${props.options?.map((o: any) => o.text).join('; ') || 'вариант 1; вариант 2'}} ${props.textAfter || ''}`;

            const { tokens } = parseSmartDropdownTemplate(fullText, props.blanks);

            return (
              <div
                key={blockId}
                className="p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-4 shadow-sm"
              >
                <div className="flex justify-between items-center">
                  <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-emerald-600 text-white">
                    Заполните пропуски (Dropdown)
                  </span>
                  <span className="text-xs font-bold text-emerald-600 dark:text-emerald-400">
                    {props.points || 10} баллов
                  </span>
                </div>
                <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
                  {props.question}
                </h4>
                <div className="p-5 bg-slate-50 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700 rounded-2xl text-xs md:text-sm text-slate-800 dark:text-slate-200 leading-loose flex flex-wrap items-center gap-1.5">
                  {tokens.map((token, tIdx) => {
                    if (token.type === 'text') {
                      return <span key={tIdx}>{token.value}</span>;
                    }
                    const b = token.blank!;
                    const userVal = blockAnswers[b.key] || '';
                    const isBlankCorrect =
                      userVal.trim().toLowerCase() === (b.correctAnswer || '').trim().toLowerCase();

                    let selectClass =
                      'border-emerald-500 bg-white dark:bg-slate-900 text-emerald-600 dark:text-emerald-400';
                    if (isSubmitted) {
                      selectClass = isBlankCorrect
                        ? 'border-emerald-500 bg-emerald-100 dark:bg-emerald-950 text-emerald-800 dark:text-emerald-200'
                        : 'border-rose-500 bg-rose-100 dark:bg-rose-950 text-rose-800 dark:text-rose-200';
                    }

                    return (
                      <span key={tIdx} className="inline-flex items-center gap-1">
                        <select
                          value={userVal}
                          disabled={isSubmitted}
                          onChange={(e) =>
                            setDropdownAnswers((prev) => ({
                              ...prev,
                              [blockId]: {
                                ...(prev[blockId] || {}),
                                [b.key]: e.target.value,
                              },
                            }))
                          }
                          className={`px-3 py-1.5 border rounded-xl font-bold text-xs focus:outline-none cursor-pointer transition-all shadow-2xs ${selectClass}`}
                        >
                          <option value="">[ Выберите вариант ]</option>
                          {b.options.map((opt: string, i: number) => (
                            <option key={i} value={opt}>
                              {opt}
                            </option>
                          ))}
                        </select>
                        {isSubmitted && (
                          <span className="text-xs">
                            {isBlankCorrect ? (
                              <Check className="text-emerald-600 inline" size={14} />
                            ) : (
                              <span className="text-[10px] text-rose-500 font-bold">({b.correctAnswer})</span>
                            )}
                          </span>
                        )}
                      </span>
                    );
                  })}
                </div>

                {!isSubmitted && (
                  <button
                    onClick={() => submitTest(blockId)}
                    className="w-full py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs rounded-xl shadow-md transition-all cursor-pointer"
                  >
                    Проверить ответы
                  </button>
                )}
              </div>
            );
          }

          case 'QuizInputBlankBlock': {
            const blockAnswers = inputAnswers[blockId] || {};
            const fullText =
              props.templateText ||
              `${props.prefixText || ''} {${props.correctAnswer || 'ответ'}} ${props.suffixText || ''}`;

            const { tokens } = parseSmartDropdownTemplate(fullText, props.blanks);

            return (
              <div
                key={blockId}
                className="p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-4 shadow-sm"
              >
                <div className="flex justify-between items-center">
                  <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-blue-600 text-white">
                    Заполните пропуски (Ручной ввод)
                  </span>
                  <span className="text-xs font-bold text-blue-600 dark:text-blue-400">
                    {props.points || 10} баллов
                  </span>
                </div>
                <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
                  {props.question}
                </h4>
                <div className="p-5 bg-slate-50 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700 rounded-2xl text-xs md:text-sm text-slate-800 dark:text-slate-200 leading-loose flex flex-wrap items-center gap-1.5">
                  {tokens.map((token, tIdx) => {
                    if (token.type === 'text') {
                      return <span key={tIdx}>{token.value}</span>;
                    }
                    const b = token.blank!;
                    const userVal = blockAnswers[b.key] || '';
                    const isBlankCorrect =
                      userVal.trim().toLowerCase() === (b.correctAnswer || '').trim().toLowerCase();

                    let inputClass = 'border-blue-500 bg-white dark:bg-slate-900 text-blue-600 dark:text-blue-400';
                    if (isSubmitted) {
                      inputClass = isBlankCorrect
                        ? 'border-emerald-500 bg-emerald-100 dark:bg-emerald-950 text-emerald-800 dark:text-emerald-200'
                        : 'border-rose-500 bg-rose-100 dark:bg-rose-950 text-rose-800 dark:text-rose-200';
                    }

                    return (
                      <span key={tIdx} className="inline-flex items-center gap-1">
                        <input
                          type="text"
                          disabled={isSubmitted}
                          value={userVal}
                          onChange={(e) =>
                            setInputAnswers((prev) => ({
                              ...prev,
                              [blockId]: {
                                ...(prev[blockId] || {}),
                                [b.key]: e.target.value,
                              },
                            }))
                          }
                          placeholder={`{${b.key}}...`}
                          className={`px-3 py-1 border rounded-xl font-bold text-xs focus:outline-none transition-all shadow-2xs min-w-[120px] ${inputClass}`}
                        />
                        {isSubmitted && (
                          <span className="text-xs">
                            {isBlankCorrect ? (
                              <Check className="text-emerald-600 inline" size={14} />
                            ) : (
                              <span className="text-[10px] text-rose-500 font-bold">({b.correctAnswer})</span>
                            )}
                          </span>
                        )}
                      </span>
                    );
                  })}
                </div>

                {!isSubmitted && (
                  <button
                    onClick={() => submitTest(blockId)}
                    className="w-full py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-bold text-xs rounded-xl shadow-md transition-all cursor-pointer"
                  >
                    Проверить ответы
                  </button>
                )}
              </div>
            );
          }

          case 'QuizSequenceBlock': {
            const items = sequenceOrders[blockId] || props.items?.map((it: any) => it.text) || [];
            return (
              <div
                key={blockId}
                className="p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-4 shadow-sm"
              >
                <div className="flex justify-between items-center">
                  <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-amber-600 text-white">
                    Расставьте по порядку
                  </span>
                  <span className="text-xs font-bold text-amber-600 dark:text-amber-400">
                    {props.points || 15} баллов
                  </span>
                </div>
                <h4 className="text-sm font-bold text-slate-900 dark:text-white">{props.question}</h4>
                <div className="space-y-2">
                  {items.map((txt: string, i: number) => (
                    <div
                      key={i}
                      className="p-3 bg-slate-50 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700 rounded-2xl flex items-center justify-between text-xs font-bold text-slate-800 dark:text-slate-200"
                    >
                      <div className="flex items-center gap-3">
                        <span className="w-5 h-5 rounded-full bg-amber-100 dark:bg-amber-950 text-amber-700 dark:text-amber-300 text-[10px] font-black flex items-center justify-center">
                          {i + 1}
                        </span>
                        <span>{txt}</span>
                      </div>
                      {!isSubmitted && (
                        <div className="flex gap-1 text-slate-400">
                          <button
                            disabled={i === 0}
                            onClick={() => moveSequenceItem(blockId, i, i - 1, props.items || [])}
                            className="p-1 hover:text-indigo-600 disabled:opacity-20 cursor-pointer"
                          >
                            <ArrowUp size={14} />
                          </button>
                          <button
                            disabled={i === items.length - 1}
                            onClick={() => moveSequenceItem(blockId, i, i + 1, props.items || [])}
                            className="p-1 hover:text-indigo-600 disabled:opacity-20 cursor-pointer"
                          >
                            <ArrowDown size={14} />
                          </button>
                        </div>
                      )}
                    </div>
                  ))}
                </div>

                {!isSubmitted && (
                  <button
                    onClick={() => submitTest(blockId)}
                    className="w-full py-2.5 bg-amber-600 hover:bg-amber-700 text-white font-bold text-xs rounded-xl shadow-md mt-2 cursor-pointer"
                  >
                    Зафиксировать последовательность
                  </button>
                )}
              </div>
            );
          }

          case 'QuizEssayBlock': {
            const answer = essayAnswers[blockId] || '';
            return (
              <div
                key={blockId}
                className="p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-4 shadow-sm"
              >
                <div className="flex justify-between items-center">
                  <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-rose-600 text-white">
                    Развернутый ответ
                  </span>
                  <span className="text-xs font-bold text-rose-600 dark:text-rose-400">
                    {props.points || 25} баллов
                  </span>
                </div>
                <h4 className="text-sm font-bold text-slate-900 dark:text-white">{props.question}</h4>
                <textarea
                  rows={5}
                  disabled={isSubmitted}
                  value={answer}
                  onChange={(e) => setEssayAnswers((prev) => ({ ...prev, [blockId]: e.target.value }))}
                  placeholder="Напишите развернутый ответ своими словами..."
                  className="w-full p-4 bg-slate-50 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700 rounded-2xl text-xs focus:ring-2 focus:ring-rose-500 focus:outline-none"
                />
                {!isSubmitted ? (
                  <button
                    disabled={!answer.trim()}
                    onClick={() => submitTest(blockId)}
                    className="w-full py-2.5 bg-rose-600 hover:bg-rose-700 disabled:opacity-40 text-white font-bold text-xs rounded-xl shadow-md flex items-center justify-center gap-2 cursor-pointer"
                  >
                    <Send size={14} />
                    <span>Отправить на проверку преподавателю</span>
                  </button>
                ) : (
                  <div className="p-3 bg-amber-50 dark:bg-amber-950/40 text-amber-700 dark:text-amber-300 rounded-xl text-xs font-medium">
                    Ответ сохранен и будет отправлен преподавателю на проверку.
                  </div>
                )}
              </div>
            );
          }

          case 'FileUploadBlock': {
            return (
              <div
                key={blockId}
                className="p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-4 shadow-sm"
              >
                <div className="flex justify-between items-center">
                  <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-slate-700 text-white">
                    Загрузка работы
                  </span>
                  <span className="text-xs font-bold text-slate-700 dark:text-slate-300">
                    {props.points || 50} баллов
                  </span>
                </div>
                <h4 className="text-sm font-bold text-slate-900 dark:text-white">{props.title}</h4>
                <p className="text-xs text-slate-600 dark:text-slate-400">{props.instructions}</p>
                <div className="p-8 border-2 border-dashed border-slate-300 dark:border-slate-700 rounded-2xl text-center space-y-2 cursor-pointer hover:border-indigo-500 transition-colors">
                  <Upload className="mx-auto text-slate-400" size={24} />
                  <span className="text-xs font-bold text-indigo-600 dark:text-indigo-400 block">
                    Нажмите для выбора архива или документа
                  </span>
                  <span className="text-[10px] text-slate-400 block">
                    Разрешено: {props.allowedTypes || '.zip, .pdf'} (до {props.maxSizeMB || 25} МБ)
                  </span>
                </div>
              </div>
            );
          }

          default:
            return null;
        }
      })}

      {/* Completion Button */}
      <div className="pt-6 border-t border-slate-200 dark:border-slate-800 flex justify-end">
        <button
          onClick={handleFinish}
          className="px-6 py-3 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold rounded-2xl shadow-lg shadow-emerald-600/20 flex items-center gap-2 cursor-pointer transition-all active:scale-95"
        >
          <Sparkles size={16} />
          <span>Завершить урок и перейти к следующему</span>
        </button>
      </div>

      {/* Completion Results Modal */}
      {completionResult && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4 animate-in fade-in duration-200">
          <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 max-w-md w-full shadow-2xl space-y-5 text-center">
            <div className="w-14 h-14 rounded-2xl bg-emerald-100 dark:bg-emerald-950 text-emerald-600 dark:text-emerald-400 flex items-center justify-center mx-auto shadow-inner">
              <Award size={32} />
            </div>

            <div className="space-y-1.5">
              <h3 className="text-lg font-black text-slate-900 dark:text-white">
                🎉 Урок успешно завершен!
              </h3>
              <p className="text-xs text-slate-500 dark:text-slate-400">
                Ваш прогресс и результаты тестирования сохранены
              </p>
            </div>

            {/* Score Card */}
            {completionResult.hasQuizzes ? (
              <div className="p-4 bg-slate-50 dark:bg-slate-800/60 rounded-2xl border border-slate-200 dark:border-slate-700 space-y-2">
                <div className="text-[11px] uppercase font-extrabold text-slate-400 tracking-wider">
                  Результат за урок
                </div>
                <div className="text-3xl font-black text-indigo-600 dark:text-indigo-400">
                  {completionResult.earnedPoints} / {completionResult.totalMaxPoints}
                  <span className="text-xs text-slate-400 font-bold ml-1.5">
                    ({completionResult.score}%)
                  </span>
                </div>
                <div className="text-[11px] text-slate-500 font-medium">
                  {completionResult.score >= 80
                    ? '🌟 Отличный результат! Материал усвоен на высоком уровне.'
                    : completionResult.score >= 50
                    ? '👍 Хорошая работа! Рекомендуем повторить неверные ответы.'
                    : '📖 Попробуйте повторить материал урока для лучшего понимания.'}
                </div>
              </div>
            ) : (
              <div className="p-4 bg-emerald-50 dark:bg-emerald-950/40 rounded-2xl border border-emerald-200 dark:border-emerald-900/50 text-xs font-bold text-emerald-700 dark:text-emerald-300">
                ✓ Теоретический материал урока успешно пройден
              </div>
            )}

            {completionResult.essaysCount > 0 && (
              <div className="p-3 bg-amber-50 dark:bg-amber-950/40 rounded-xl border border-amber-200 dark:border-amber-900/50 text-xs text-amber-800 dark:text-amber-300 text-left flex items-center gap-2 font-medium">
                <span className="w-2 h-2 rounded-full bg-amber-500 animate-pulse flex-shrink-0" />
                <span>
                  {completionResult.essaysCount} зад. с развернутым ответом отправлено на проверку преподавателю.
                </span>
              </div>
            )}

            <div className="pt-2 flex flex-col gap-2">
              {onNavigateBack && (
                <button
                  onClick={onNavigateBack}
                  className="w-full py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs rounded-xl shadow-lg shadow-indigo-600/20 transition-all cursor-pointer"
                >
                  Вернуться к программе курса
                </button>
              )}
              <button
                onClick={() => setCompletionResult(null)}
                className="w-full py-2 text-xs font-bold text-slate-500 hover:text-slate-700 dark:hover:text-slate-300 transition-colors cursor-pointer"
              >
                Просмотреть материалы урока
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Toast Notification */}
      {toastMsg && (
        <div className="fixed bottom-6 right-6 z-50 bg-indigo-600 text-white text-xs font-bold px-4 py-3 rounded-2xl shadow-xl flex items-center gap-2 animate-bounce">
          <CheckCircle size={16} />
          <span>{toastMsg}</span>
        </div>
      )}
    </div>
  );
}
