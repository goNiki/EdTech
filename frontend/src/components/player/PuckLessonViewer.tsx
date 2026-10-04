'use client';

import React, { useState, useRef } from 'react';
import { parseSmartDropdownTemplate } from '@/lib/puck-config';
import { api } from '@/lib/api';
import {
  CheckCircle,
  CheckCircle2,
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
  X,
  Loader2,
  Trash2,
  Paperclip,
  ShieldCheck,
  RotateCcw
} from 'lucide-react';
import PresentationViewer from './PresentationViewer';

export interface LessonAnswerItem {
  block_id: string;
  answer: any;
}

export interface LessonCompletionPayload {
  score: number;
  answers: LessonAnswerItem[];
  essays: Array<{
    question_text: string;
    answer_text: string;
    max_points: number;
  }>;
}

interface PuckLessonViewerProps {
  contentJson: string | object;
  onComplete?: (payload: LessonCompletionPayload) => Promise<any> | void;
  onNavigateBack?: () => void;
  initialProgress?: any;
}

export default function PuckLessonViewer({
  contentJson,
  onComplete,
  onNavigateBack,
  initialProgress,
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

  // Server Anti-cheat Verification states
  const [isVerifying, setIsVerifying] = useState(false);
  const [serverValidationResults, setServerValidationResults] = useState<
    Record<string, { isCorrect: boolean; feedback?: string; correctAnswer?: any }>
  >({});

  // File Upload states
  const [uploadedFiles, setUploadedFiles] = useState<
    Record<string, { name: string; size: number; url: string }>
  >({});
  const [uploadProgress, setUploadProgress] = useState<Record<string, number>>({});
  const [isUploading, setIsUploading] = useState<Record<string, boolean>>({});
  const [dragOverBlocks, setDragOverBlocks] = useState<Record<string, boolean>>({});
  const [fileErrors, setFileErrors] = useState<Record<string, string>>({});

  const [completionResult, setCompletionResult] = useState<{
    score: number;
    earnedPoints: number;
    totalMaxPoints: number;
    hasQuizzes: boolean;
    essaysCount: number;
    isServerVerified?: boolean;
  } | null>(null);

  const showToast = (msg: string) => {
    setToastMsg(msg);
    setTimeout(() => setToastMsg(null), 3500);
  };

  const handleFileUpload = async (blockId: string, file: File, maxSizeMB: number = 25) => {
    if (file.size > maxSizeMB * 1024 * 1024) {
      const errText = `Файл превышает допустимый размер ${maxSizeMB} МБ`;
      setFileErrors((prev) => ({ ...prev, [blockId]: errText }));
      showToast(errText);
      return;
    }

    setIsUploading((prev) => ({ ...prev, [blockId]: true }));
    setUploadProgress((prev) => ({ ...prev, [blockId]: 10 }));
    setFileErrors((prev) => ({ ...prev, [blockId]: '' }));

    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('category', 'homework');

      const res = await api.post('/upload', formData, {
        params: { category: 'homework' },
        headers: { 'Content-Type': 'multipart/form-data' },
        onUploadProgress: (progressEvent) => {
          if (progressEvent.total) {
            const pct = Math.round((progressEvent.loaded * 90) / progressEvent.total);
            setUploadProgress((prev) => ({ ...prev, [blockId]: Math.max(10, pct) }));
          }
        },
      });

      const data = res.data.data || res.data;
      const fileUrl = data.file_url || data.url || data.FileUrl || '';

      const resolveUrl = (url: string) => {
        if (!url) return '';
        if (url.startsWith('http://') || url.startsWith('https://')) return url;
        const baseUrl = (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8082/api/v1').replace(/\/api\/v1\/?$/, '');
        return `${baseUrl}${url.startsWith('/') ? '' : '/'}${url}`;
      };

      setUploadedFiles((prev) => ({
        ...prev,
        [blockId]: {
          name: file.name,
          size: file.size,
          url: resolveUrl(fileUrl),
        },
      }));
      setUploadProgress((prev) => ({ ...prev, [blockId]: 100 }));
      showToast(`Файл «${file.name}» успешно загружен`);
    } catch (err: any) {
      console.error('Failed to upload homework file', err);
      const msg = err.response?.data?.message || err.response?.data?.error || 'Ошибка при загрузке файла';
      setFileErrors((prev) => ({ ...prev, [blockId]: msg }));
      showToast(msg);
    } finally {
      setIsUploading((prev) => ({ ...prev, [blockId]: false }));
    }
  };

  const handleRemoveFile = (blockId: string) => {
    setUploadedFiles((prev) => {
      const copy = { ...prev };
      delete copy[blockId];
      return copy;
    });
    setUploadProgress((prev) => {
      const copy = { ...prev };
      delete copy[blockId];
      return copy;
    });
    setFileErrors((prev) => {
      const copy = { ...prev };
      delete copy[blockId];
      return copy;
    });
  };

  const handleSingleSelect = (blockId: string, optIdx: number) => {
    if (submittedBlocks[blockId] && !isVerifying) return;
    setSingleAnswers((prev) => ({ ...prev, [blockId]: optIdx }));
  };

  const handleMultiToggle = (blockId: string, optIdx: number) => {
    if (submittedBlocks[blockId] && !isVerifying) return;
    setMultiAnswers((prev) => {
      const current = prev[blockId] || [];
      const updated = current.includes(optIdx)
        ? current.filter((i) => i !== optIdx)
        : [...current, optIdx];
      return { ...prev, [blockId]: updated };
    });
  };

  const handleMatchSelect = (blockId: string, pairIdx: number, selectedValue: string) => {
    if (submittedBlocks[blockId] && !isVerifying) return;
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
    const answersList: LessonAnswerItem[] = [];

    blocks.forEach((block: any, idx: number) => {
      const { type, props = {} } = block;
      const blockId = props.id || `block-${idx}`;

      switch (type) {
        case 'QuizSingleBlock': {
          const pts = Number(props.points) || 10;
          totalMaxPoints += pts;
          const selected = singleAnswers[blockId];
          answersList.push({
            block_id: blockId,
            answer: { type: 'single', selected_option: selected },
          });

          const correctIdx = props.options?.findIndex(
            (o: any) => o.isCorrect === true || o.isCorrect === 'true'
          );
          if (selected !== undefined && correctIdx !== -1 && selected === correctIdx) {
            earnedPoints += pts;
          }
          break;
        }

        case 'QuizMultiBlock': {
          const pts = Number(props.points) || 15;
          totalMaxPoints += pts;
          const selectedList = multiAnswers[blockId] || [];
          answersList.push({
            block_id: blockId,
            answer: { type: 'multi', selected_options: selectedList },
          });

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
          answersList.push({
            block_id: blockId,
            answer: { type: 'match', pairs: userMatches },
          });

          const pairs: any[] = props.pairs || [];
          if (pairs.length > 0) {
            let correctPairs = 0;
            pairs.forEach((p, pIdx) => {
              if (p.right && (userMatches[pIdx] || '').trim().toLowerCase() === (p.right || '').trim().toLowerCase()) {
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
          answersList.push({
            block_id: blockId,
            answer: { type: 'dropdown', blanks: userVals },
          });

          const fullText = props.templateText || '';
          const { blanks } = parseSmartDropdownTemplate(fullText, props.blanks);
          if (blanks.length > 0) {
            let correctBlanks = 0;
            blanks.forEach((b) => {
              if (b.correctAnswer && (userVals[b.key] || '').trim().toLowerCase() === (b.correctAnswer || '').trim().toLowerCase()) {
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
          answersList.push({
            block_id: blockId,
            answer: { type: 'input', blanks: userVals },
          });

          const fullText = props.templateText || '';
          const { blanks } = parseSmartDropdownTemplate(fullText, props.blanks);
          if (blanks.length > 0) {
            let correctBlanks = 0;
            blanks.forEach((b) => {
              if (b.correctAnswer && (userVals[b.key] || '').trim().toLowerCase() === (b.correctAnswer || '').trim().toLowerCase()) {
                correctBlanks++;
              }
            });
            earnedPoints += (correctBlanks / blanks.length) * pts;
          }
          break;
        }

        case 'QuizSequenceBlock': {
          const pts = Number(props.points) || 15;
          totalMaxPoints += pts;
          const order = sequenceOrders[blockId] || props.items?.map((it: any) => it.text) || [];
          answersList.push({
            block_id: blockId,
            answer: { type: 'sequence', order },
          });
          // Sequence points calculated on server or default
          earnedPoints += pts;
          break;
        }

        case 'QuizEssayBlock': {
          const answer = (essayAnswers[blockId] || '').trim();
          answersList.push({
            block_id: blockId,
            answer: { type: 'essay', text: answer },
          });
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
          const uploaded = uploadedFiles[blockId];
          answersList.push({
            block_id: blockId,
            answer: { type: 'file', file: uploaded || null },
          });
          submittedEssays.push({
            question_text: props.title || 'Загрузка практической работы',
            answer_text: uploaded
              ? `Файл решения: ${uploaded.name} (скачать: ${uploaded.url})`
              : 'Файл решения не был прикреплен',
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
        answers: answersList,
        essays: submittedEssays,
      },
      earnedPoints: Math.round(earnedPoints),
      totalMaxPoints: Math.round(totalMaxPoints),
    };
  };

  const handleFinish = async () => {
    setIsVerifying(true);
    try {
      const { payload, earnedPoints, totalMaxPoints } = calculateCompletionPayload();
      let serverResponse: any = null;

      if (onComplete) {
        serverResponse = await onComplete(payload);
      }

      // Parse server validation response (anti-cheat verification)
      const verifiedScore = serverResponse?.score ?? serverResponse?.Score ?? payload.score;
      const validatedMap: Record<string, { isCorrect: boolean; feedback?: string; correctAnswer?: any }> = {};

      const rawResults = serverResponse?.results || serverResponse?.validated_blocks || serverResponse?.answers;
      if (rawResults && typeof rawResults === 'object') {
        Object.entries(rawResults).forEach(([bid, item]: [string, any]) => {
          validatedMap[bid] = {
            isCorrect: Boolean(item.is_correct ?? item.isCorrect ?? item.correct),
            feedback: item.feedback || item.message,
            correctAnswer: item.correct_answer ?? item.correctAnswer,
          };
        });
      }

      setServerValidationResults(validatedMap);

      // Mark all blocks as submitted to display highlights
      const updatedSubmitted: Record<string, boolean> = { ...submittedBlocks };
      blocks.forEach((b: any, idx: number) => {
        const bid = b.props?.id || `block-${idx}`;
        updatedSubmitted[bid] = true;
      });
      setSubmittedBlocks(updatedSubmitted);

      setCompletionResult({
        score: verifiedScore,
        earnedPoints: serverResponse?.earned_points ?? earnedPoints,
        totalMaxPoints: serverResponse?.total_max_points ?? totalMaxPoints,
        hasQuizzes: totalMaxPoints > 0,
        essaysCount: payload.essays.length,
        isServerVerified: true,
      });
    } catch (err) {
      console.error('Failed to submit lesson for validation', err);
      showToast('Ошибка при проверке ответов на сервере');
    } finally {
      setIsVerifying(false);
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
        const serverRes = serverValidationResults[blockId];

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
                  className="lesson-content-area my-4 leading-relaxed prose prose-slate dark:prose-invert max-w-none prose-headings:font-black prose-p:leading-relaxed prose-table:my-0 text-slate-700 dark:text-slate-300"
                  dangerouslySetInnerHTML={{ __html: props.contentHtml }}
                />
              );
            }
            return (
              <div
                key={blockId}
                className="lesson-content-area my-4 leading-relaxed whitespace-pre-wrap text-slate-700 dark:text-slate-300"
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
                  className="lesson-content-area prose prose-slate dark:prose-invert max-w-none prose-headings:font-black prose-p:leading-relaxed prose-table:my-0 text-slate-800 dark:text-slate-200 leading-relaxed"
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

          case 'PresentationBlock': {
            return (
              <div key={blockId} className="my-6">
                <PresentationViewer {...props} />
              </div>
            );
          }

          case 'QuizSingleBlock': {
            const selectedOpt = singleAnswers[blockId];
            const hasServerResult = serverRes !== undefined;
            const isCorrectServer = serverRes?.isCorrect;

            return (
              <div
                key={blockId}
                className={`p-6 bg-white dark:bg-slate-900 border rounded-3xl space-y-4 shadow-sm transition-all ${
                  hasServerResult
                    ? isCorrectServer
                      ? 'border-emerald-400/80 dark:border-emerald-800 ring-1 ring-emerald-400/30'
                      : 'border-rose-400/80 dark:border-rose-800 ring-1 ring-rose-400/30'
                    : 'border-slate-200 dark:border-slate-800'
                }`}
              >
                <div className="flex justify-between items-center">
                  <div className="flex items-center gap-2">
                    <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-indigo-600 text-white">
                      Один ответ
                    </span>
                    {hasServerResult && (
                      <span
                        className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-lg text-[10px] font-bold ${
                          isCorrectServer
                            ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300'
                            : 'bg-rose-100 dark:bg-rose-950 text-rose-700 dark:text-rose-300'
                        }`}
                      >
                        <ShieldCheck size={12} />
                        <span>{isCorrectServer ? 'Сервер: Верно' : 'Сервер: Неверно'}</span>
                      </span>
                    )}
                  </div>
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
                    const hasCorrectFlag = opt.isCorrect !== undefined;
                    const isCorrectLocal = opt.isCorrect === 'true' || opt.isCorrect === true;

                    let btnClass = 'border-slate-200 dark:border-slate-800 hover:border-indigo-400';
                    if (isSelected && !isSubmitted) {
                      btnClass = 'border-indigo-600 bg-indigo-50/50 dark:bg-indigo-950/40 ring-1 ring-indigo-500';
                    } else if (isSubmitted) {
                      if (hasServerResult) {
                        if (isSelected) {
                          btnClass = isCorrectServer
                            ? 'border-emerald-500 bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300 ring-1 ring-emerald-500/20'
                            : 'border-rose-500 bg-rose-50 dark:bg-rose-950/40 text-rose-700 dark:text-rose-300 ring-1 ring-rose-500/20';
                        }
                      } else if (hasCorrectFlag) {
                        if (isCorrectLocal) {
                          btnClass = 'border-emerald-500 bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300 ring-1 ring-emerald-500/20';
                        } else if (isSelected && !isCorrectLocal) {
                          btnClass = 'border-rose-500 bg-rose-50 dark:bg-rose-950/40 text-rose-700 dark:text-rose-300 ring-1 ring-rose-500/20';
                        }
                      } else if (isSelected) {
                        btnClass = 'border-indigo-500 bg-indigo-50/40 dark:bg-indigo-950/30 text-indigo-700 dark:text-indigo-300';
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
                              {hasServerResult && isSelected ? (
                                isCorrectServer ? (
                                  <CheckCircle size={16} className="text-emerald-600" />
                                ) : (
                                  <XCircle size={16} className="text-rose-600" />
                                )
                              ) : hasCorrectFlag ? (
                                isCorrectLocal ? (
                                  <CheckCircle size={16} className="text-emerald-600" />
                                ) : isSelected ? (
                                  <XCircle size={16} className="text-rose-600" />
                                ) : null
                              ) : isSelected ? (
                                <Check size={16} className="text-indigo-600" />
                              ) : null}
                            </span>
                          )}
                        </div>

                        {/* Hint / Explanation */}
                        {isSubmitted && opt.explain && (isSelected || isCorrectLocal) && (
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
                    Зафиксировать выбор
                  </button>
                )}
              </div>
            );
          }

          case 'QuizMultiBlock': {
            const selectedList = multiAnswers[blockId] || [];
            const hasServerResult = serverRes !== undefined;
            const isCorrectServer = serverRes?.isCorrect;

            return (
              <div
                key={blockId}
                className={`p-6 bg-white dark:bg-slate-900 border rounded-3xl space-y-4 shadow-sm transition-all ${
                  hasServerResult
                    ? isCorrectServer
                      ? 'border-emerald-400/80 dark:border-emerald-800 ring-1 ring-emerald-400/30'
                      : 'border-rose-400/80 dark:border-rose-800 ring-1 ring-rose-400/30'
                    : 'border-slate-200 dark:border-slate-800'
                }`}
              >
                <div className="flex justify-between items-center">
                  <div className="flex items-center gap-2">
                    <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-purple-600 text-white">
                      Множественный выбор
                    </span>
                    {hasServerResult && (
                      <span
                        className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-lg text-[10px] font-bold ${
                          isCorrectServer
                            ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300'
                            : 'bg-rose-100 dark:bg-rose-950 text-rose-700 dark:text-rose-300'
                        }`}
                      >
                        <ShieldCheck size={12} />
                        <span>{isCorrectServer ? 'Сервер: Верно' : 'Сервер: Неверно'}</span>
                      </span>
                    )}
                  </div>
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
                    const hasCorrectFlag = opt.isCorrect !== undefined;
                    const isCorrectLocal = opt.isCorrect === 'true' || opt.isCorrect === true;

                    let btnClass = 'border-slate-200 dark:border-slate-800 hover:border-purple-400';
                    if (isSelected && !isSubmitted) {
                      btnClass = 'border-purple-600 bg-purple-50/50 dark:bg-purple-950/40 ring-1 ring-purple-500';
                    } else if (isSubmitted) {
                      if (hasServerResult) {
                        if (isSelected) {
                          btnClass = isCorrectServer
                            ? 'border-emerald-500 bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300'
                            : 'border-rose-500 bg-rose-50 dark:bg-rose-950/40 text-rose-700 dark:text-rose-300';
                        }
                      } else if (hasCorrectFlag) {
                        if (isCorrectLocal) {
                          btnClass = 'border-emerald-500 bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300';
                        } else if (isSelected && !isCorrectLocal) {
                          btnClass = 'border-rose-500 bg-rose-50 dark:bg-rose-950/40 text-rose-700 dark:text-rose-300';
                        }
                      } else if (isSelected) {
                        btnClass = 'border-purple-500 bg-purple-50/40 dark:bg-purple-950/30 text-purple-700 dark:text-purple-300';
                      }
                    }

                    return (
                      <div key={optIdx} className="space-y-1.5">
                        <div
                          onClick={() => handleMultiToggle(blockId, optIdx)}
                          className={`p-3.5 border rounded-2xl flex items-center justify-between cursor-pointer transition-all text-xs font-medium ${btnClass}`}
                        >
                          <div className="flex items-center gap-3">
                            <div
                              className={`w-4 h-4 rounded border flex items-center justify-center ${
                                isSelected ? 'border-purple-600 bg-purple-600 text-white' : 'border-slate-300 dark:border-slate-700'
                              }`}
                            >
                              {isSelected && <Check size={12} strokeWidth={3} />}
                            </div>
                            <span>{opt.text}</span>
                          </div>
                          {isSubmitted && (
                            <span>
                              {hasServerResult && isSelected ? (
                                isCorrectServer ? (
                                  <CheckCircle size={16} className="text-emerald-600" />
                                ) : (
                                  <XCircle size={16} className="text-rose-600" />
                                )
                              ) : hasCorrectFlag ? (
                                isCorrectLocal ? (
                                  <CheckCircle size={16} className="text-emerald-600" />
                                ) : isSelected ? (
                                  <XCircle size={16} className="text-rose-600" />
                                ) : null
                              ) : isSelected ? (
                                <Check size={16} className="text-purple-600" />
                              ) : null}
                            </span>
                          )}
                        </div>

                        {isSubmitted && opt.explain && (isSelected || isCorrectLocal) && (
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
                    disabled={selectedList.length === 0}
                    onClick={() => submitTest(blockId)}
                    className="w-full py-2.5 bg-purple-600 hover:bg-purple-700 disabled:opacity-40 text-white font-bold text-xs rounded-xl transition-all shadow-md mt-2 cursor-pointer"
                  >
                    Зафиксировать выбор
                  </button>
                )}
              </div>
            );
          }

          case 'QuizMatchBlock': {
            const blockMatches = matchAnswers[blockId] || {};
            const pairs: Array<{ left: string; right: string }> = props.pairs || [];
            const allRightOptions = Array.from(new Set(pairs.map((p) => p.right))).filter(Boolean);
            const hasServerResult = serverRes !== undefined;
            const isCorrectServer = serverRes?.isCorrect;

            return (
              <div
                key={blockId}
                className={`p-6 bg-white dark:bg-slate-900 border rounded-3xl space-y-4 shadow-sm transition-all ${
                  hasServerResult
                    ? isCorrectServer
                      ? 'border-emerald-400/80 dark:border-emerald-800 ring-1 ring-emerald-400/30'
                      : 'border-rose-400/80 dark:border-rose-800 ring-1 ring-rose-400/30'
                    : 'border-slate-200 dark:border-slate-800'
                }`}
              >
                <div className="flex justify-between items-center">
                  <div className="flex items-center gap-2">
                    <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-cyan-600 text-white">
                      Сопоставление элементов
                    </span>
                    {hasServerResult && (
                      <span
                        className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-lg text-[10px] font-bold ${
                          isCorrectServer
                            ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300'
                            : 'bg-rose-100 dark:bg-rose-950 text-rose-700 dark:text-rose-300'
                        }`}
                      >
                        <ShieldCheck size={12} />
                        <span>{isCorrectServer ? 'Сервер: Верно' : 'Сервер: Неверно'}</span>
                      </span>
                    )}
                  </div>
                  <span className="text-xs font-bold text-cyan-600 dark:text-cyan-400">
                    {props.points || 20} баллов
                  </span>
                </div>
                <h4 className="text-sm font-bold text-slate-900 dark:text-white">{props.question}</h4>

                <div className="space-y-3">
                  {pairs.map((p, pIdx) => {
                    const userSelected = blockMatches[pIdx] || '';
                    const hasRightAnswer = Boolean(p.right);
                    const isPairCorrect = hasRightAnswer && userSelected.trim().toLowerCase() === p.right.trim().toLowerCase();

                    let borderClass = 'border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/60';
                    if (isSubmitted) {
                      if (hasServerResult) {
                        borderClass = isCorrectServer
                          ? 'border-emerald-500 bg-emerald-50/60 dark:bg-emerald-950/40 text-emerald-800 dark:text-emerald-200'
                          : 'border-rose-500 bg-rose-50/60 dark:bg-rose-950/40 text-rose-800 dark:text-rose-200';
                      } else if (hasRightAnswer) {
                        borderClass = isPairCorrect
                          ? 'border-emerald-500 bg-emerald-50/60 dark:bg-emerald-950/40 text-emerald-800 dark:text-emerald-200'
                          : 'border-rose-500 bg-rose-50/60 dark:bg-rose-950/40 text-rose-800 dark:text-rose-200';
                      }
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

                            {isSubmitted && hasRightAnswer && !isPairCorrect && !hasServerResult && (
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
                    Зафиксировать сопоставление
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
            const hasServerResult = serverRes !== undefined;
            const isCorrectServer = serverRes?.isCorrect;

            return (
              <div
                key={blockId}
                className={`p-6 bg-white dark:bg-slate-900 border rounded-3xl space-y-4 shadow-sm transition-all ${
                  hasServerResult
                    ? isCorrectServer
                      ? 'border-emerald-400/80 dark:border-emerald-800 ring-1 ring-emerald-400/30'
                      : 'border-rose-400/80 dark:border-rose-800 ring-1 ring-rose-400/30'
                    : 'border-slate-200 dark:border-slate-800'
                }`}
              >
                <div className="flex justify-between items-center">
                  <div className="flex items-center gap-2">
                    <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-emerald-600 text-white">
                      Заполните пропуски (Dropdown)
                    </span>
                    {hasServerResult && (
                      <span
                        className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-lg text-[10px] font-bold ${
                          isCorrectServer
                            ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300'
                            : 'bg-rose-100 dark:bg-rose-950 text-rose-700 dark:text-rose-300'
                        }`}
                      >
                        <ShieldCheck size={12} />
                        <span>{isCorrectServer ? 'Сервер: Верно' : 'Сервер: Неверно'}</span>
                      </span>
                    )}
                  </div>
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
                    const hasCorrectAnswer = Boolean(b.correctAnswer);
                    const isBlankCorrect =
                      hasCorrectAnswer &&
                      userVal.trim().toLowerCase() === (b.correctAnswer || '').trim().toLowerCase();

                    let selectClass =
                      'border-emerald-500 bg-white dark:bg-slate-900 text-emerald-600 dark:text-emerald-400';
                    if (isSubmitted) {
                      if (hasServerResult) {
                        selectClass = isCorrectServer
                          ? 'border-emerald-500 bg-emerald-100 dark:bg-emerald-950 text-emerald-800 dark:text-emerald-200'
                          : 'border-rose-500 bg-rose-100 dark:bg-rose-950 text-rose-800 dark:text-rose-200';
                      } else if (hasCorrectAnswer) {
                        selectClass = isBlankCorrect
                          ? 'border-emerald-500 bg-emerald-100 dark:bg-emerald-950 text-emerald-800 dark:text-emerald-200'
                          : 'border-rose-500 bg-rose-100 dark:bg-rose-950 text-rose-800 dark:text-rose-200';
                      }
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
                        {isSubmitted && hasCorrectAnswer && !hasServerResult && (
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
                    Зафиксировать ответы
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
            const hasServerResult = serverRes !== undefined;
            const isCorrectServer = serverRes?.isCorrect;

            return (
              <div
                key={blockId}
                className={`p-6 bg-white dark:bg-slate-900 border rounded-3xl space-y-4 shadow-sm transition-all ${
                  hasServerResult
                    ? isCorrectServer
                      ? 'border-emerald-400/80 dark:border-emerald-800 ring-1 ring-emerald-400/30'
                      : 'border-rose-400/80 dark:border-rose-800 ring-1 ring-rose-400/30'
                    : 'border-slate-200 dark:border-slate-800'
                }`}
              >
                <div className="flex justify-between items-center">
                  <div className="flex items-center gap-2">
                    <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-blue-600 text-white">
                      Заполните пропуски (Ручной ввод)
                    </span>
                    {hasServerResult && (
                      <span
                        className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-lg text-[10px] font-bold ${
                          isCorrectServer
                            ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300'
                            : 'bg-rose-100 dark:bg-rose-950 text-rose-700 dark:text-rose-300'
                        }`}
                      >
                        <ShieldCheck size={12} />
                        <span>{isCorrectServer ? 'Сервер: Верно' : 'Сервер: Неверно'}</span>
                      </span>
                    )}
                  </div>
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
                    const hasCorrectAnswer = Boolean(b.correctAnswer);
                    const isBlankCorrect =
                      hasCorrectAnswer &&
                      userVal.trim().toLowerCase() === (b.correctAnswer || '').trim().toLowerCase();

                    let inputClass = 'border-blue-500 bg-white dark:bg-slate-900 text-blue-600 dark:text-blue-400';
                    if (isSubmitted) {
                      if (hasServerResult) {
                        inputClass = isCorrectServer
                          ? 'border-emerald-500 bg-emerald-100 dark:bg-emerald-950 text-emerald-800 dark:text-emerald-200'
                          : 'border-rose-500 bg-rose-100 dark:bg-rose-950 text-rose-800 dark:text-rose-200';
                      } else if (hasCorrectAnswer) {
                        inputClass = isBlankCorrect
                          ? 'border-emerald-500 bg-emerald-100 dark:bg-emerald-950 text-emerald-800 dark:text-emerald-200'
                          : 'border-rose-500 bg-rose-100 dark:bg-rose-950 text-rose-800 dark:text-rose-200';
                      }
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
                        {isSubmitted && hasCorrectAnswer && !hasServerResult && (
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
                    Зафиксировать ответы
                  </button>
                )}
              </div>
            );
          }

          case 'QuizSequenceBlock': {
            const items = sequenceOrders[blockId] || props.items?.map((it: any) => it.text) || [];
            const hasServerResult = serverRes !== undefined;
            const isCorrectServer = serverRes?.isCorrect;

            return (
              <div
                key={blockId}
                className={`p-6 bg-white dark:bg-slate-900 border rounded-3xl space-y-4 shadow-sm transition-all ${
                  hasServerResult
                    ? isCorrectServer
                      ? 'border-emerald-400/80 dark:border-emerald-800 ring-1 ring-emerald-400/30'
                      : 'border-rose-400/80 dark:border-rose-800 ring-1 ring-rose-400/30'
                    : 'border-slate-200 dark:border-slate-800'
                }`}
              >
                <div className="flex justify-between items-center">
                  <div className="flex items-center gap-2">
                    <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-amber-600 text-white">
                      Расставьте по порядку
                    </span>
                    {hasServerResult && (
                      <span
                        className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-lg text-[10px] font-bold ${
                          isCorrectServer
                            ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300'
                            : 'bg-rose-100 dark:bg-rose-950 text-rose-700 dark:text-rose-300'
                        }`}
                      >
                        <ShieldCheck size={12} />
                        <span>{isCorrectServer ? 'Сервер: Верно' : 'Сервер: Неверно'}</span>
                      </span>
                    )}
                  </div>
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
                        <div className="flex items-center gap-1">
                          <button
                            disabled={i === 0}
                            onClick={() => moveSequenceItem(blockId, i, i - 1, props.items || [])}
                            className="p-1.5 hover:bg-slate-200 dark:hover:bg-slate-700 rounded-lg disabled:opacity-30 cursor-pointer"
                          >
                            <ArrowUp size={14} />
                          </button>
                          <button
                            disabled={i === items.length - 1}
                            onClick={() => moveSequenceItem(blockId, i, i + 1, props.items || [])}
                            className="p-1.5 hover:bg-slate-200 dark:hover:bg-slate-700 rounded-lg disabled:opacity-30 cursor-pointer"
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
                    className="w-full py-2.5 bg-amber-600 hover:bg-amber-700 text-white font-bold text-xs rounded-xl shadow-md transition-all cursor-pointer"
                  >
                    Зафиксировать порядок
                  </button>
                )}
              </div>
            );
          }

          case 'QuizEssayBlock': {
            const val = essayAnswers[blockId] || '';
            return (
              <div
                key={blockId}
                className="p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-4 shadow-sm"
              >
                <div className="flex justify-between items-center">
                  <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-rose-600 text-white">
                    Развернутый ответ (Эссе)
                  </span>
                  <span className="text-xs font-bold text-rose-600 dark:text-rose-400">
                    до {props.points || 25} баллов
                  </span>
                </div>
                <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
                  {props.question}
                </h4>
                {props.hint && (
                  <p className="text-xs text-slate-500 dark:text-slate-400 italic">
                    Подсказка: {props.hint}
                  </p>
                )}
                <textarea
                  disabled={isSubmitted}
                  rows={props.minLines || 4}
                  value={val}
                  onChange={(e) => setEssayAnswers((prev) => ({ ...prev, [blockId]: e.target.value }))}
                  placeholder="Введите ваш ответ здесь..."
                  className="w-full p-4 bg-slate-50 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700 rounded-2xl text-xs md:text-sm font-medium focus:ring-2 focus:ring-rose-500 focus:outline-none disabled:opacity-70 leading-relaxed"
                />
                {!isSubmitted && (
                  <button
                    disabled={!val.trim()}
                    onClick={() => submitTest(blockId)}
                    className="w-full py-2.5 bg-rose-600 hover:bg-rose-700 disabled:opacity-40 text-white font-bold text-xs rounded-xl transition-all shadow-md cursor-pointer"
                  >
                    Сохранить черновик ответа
                  </button>
                )}
              </div>
            );
          }

          case 'FileUploadBlock': {
            const uploaded = uploadedFiles[blockId];
            const uploading = isUploading[blockId];
            const pct = uploadProgress[blockId] || 0;
            const err = fileErrors[blockId];
            const isDragOver = dragOverBlocks[blockId];
            const allowed = props.allowedTypes || 'Любые файлы';
            const maxMB = Number(props.maxSizeMB) || 25;

            return (
              <div
                key={blockId}
                className="p-6 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl space-y-4 shadow-sm"
              >
                <div className="flex justify-between items-center">
                  <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-emerald-600 text-white flex items-center gap-1.5">
                    <Upload size={12} />
                    <span>Практическое задание (Файл)</span>
                  </span>
                  <span className="text-xs font-bold text-emerald-600 dark:text-emerald-400">
                    до {props.points || 50} баллов
                  </span>
                </div>

                <div className="space-y-1">
                  <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
                    {props.title || 'Загрузите решение практической работы'}
                  </h4>
                  {props.description && (
                    <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
                      {props.description}
                    </p>
                  )}
                </div>

                {uploaded ? (
                  <div className="p-4 bg-emerald-50/60 dark:bg-emerald-950/30 border border-emerald-200 dark:border-emerald-800/80 rounded-2xl flex items-center justify-between gap-4">
                    <div className="flex items-center gap-3 overflow-hidden">
                      <div className="w-10 h-10 rounded-xl bg-emerald-500 text-white flex items-center justify-center flex-shrink-0">
                        <Paperclip size={18} />
                      </div>
                      <div className="overflow-hidden">
                        <p className="text-xs font-bold text-slate-900 dark:text-white truncate">
                          {uploaded.name}
                        </p>
                        <p className="text-[10px] text-slate-500">
                          {(uploaded.size / (1024 * 1024)).toFixed(2)} МБ •{' '}
                          <a
                            href={uploaded.url}
                            target="_blank"
                            rel="noreferrer"
                            className="text-emerald-600 hover:underline font-semibold"
                          >
                            Просмотреть
                          </a>
                        </p>
                      </div>
                    </div>

                    {!isSubmitted && (
                      <button
                        type="button"
                        onClick={() => handleRemoveFile(blockId)}
                        className="p-2 rounded-xl text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-colors cursor-pointer"
                        title="Удалить файл"
                      >
                        <Trash2 size={16} />
                      </button>
                    )}
                  </div>
                ) : (
                  <div
                    onDragOver={(e) => {
                      e.preventDefault();
                      setDragOverBlocks((prev) => ({ ...prev, [blockId]: true }));
                    }}
                    onDragLeave={(e) => {
                      e.preventDefault();
                      setDragOverBlocks((prev) => ({ ...prev, [blockId]: false }));
                    }}
                    onDrop={(e) => {
                      e.preventDefault();
                      setDragOverBlocks((prev) => ({ ...prev, [blockId]: false }));
                      const file = e.dataTransfer.files?.[0];
                      if (file) handleFileUpload(blockId, file, maxMB);
                    }}
                    className={`border-2 border-dashed rounded-2xl p-6 text-center transition-all ${
                      isDragOver
                        ? 'border-indigo-500 bg-indigo-50/50 dark:bg-indigo-950/30 scale-[1.01]'
                        : 'border-slate-200 dark:border-slate-700 hover:border-slate-300 dark:hover:border-slate-600 bg-slate-50/50 dark:bg-slate-800/30'
                    }`}
                  >
                    <input
                      type="file"
                      id={`file-${blockId}`}
                      className="hidden"
                      disabled={uploading || isSubmitted}
                      onChange={(e) => {
                        const file = e.target.files?.[0];
                        if (file) handleFileUpload(blockId, file, maxMB);
                      }}
                    />
                    <label
                      htmlFor={`file-${blockId}`}
                      className="cursor-pointer flex flex-col items-center justify-center space-y-2"
                    >
                      <div className="w-12 h-12 rounded-2xl bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 flex items-center justify-center shadow-xs">
                        {uploading ? (
                          <Loader2 size={22} className="animate-spin" />
                        ) : (
                          <Upload size={22} />
                        )}
                      </div>

                      <div>
                        <span className="text-xs font-bold text-indigo-600 dark:text-indigo-400 hover:underline">
                          {uploading ? 'Загрузка файла...' : 'Нажмите для выбора файла'}
                        </span>
                        <span className="text-xs text-slate-500"> или перетащите его сюда</span>
                      </div>

                      {uploading && (
                        <div className="w-full max-w-xs mt-2 space-y-1">
                          <div className="h-1.5 w-full bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
                            <div
                              className="h-full bg-indigo-600 transition-all duration-300"
                              style={{ width: `${pct}%` }}
                            />
                          </div>
                          <span className="text-[10px] text-slate-400 font-bold">{pct}%</span>
                        </div>
                      )}

                      {err && (
                        <p className="text-[11px] text-rose-500 font-semibold">{err}</p>
                      )}

                      <span className="text-[10px] text-slate-400 block">
                        Разрешено: {allowed} (до {maxMB} МБ)
                      </span>
                    </label>
                  </div>
                )}
              </div>
            );
          }

          default:
            return null;
        }
      })}

      {/* Completion Button */}
      <div className="pt-6 border-t border-slate-200 dark:border-slate-800 flex flex-col sm:flex-row items-center justify-between gap-4">
        {initialProgress && (initialProgress.status === 'completed' || initialProgress.Status === 'completed') ? (
          <div className="flex items-center gap-2 text-xs font-semibold text-emerald-700 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/40 px-4 py-2.5 rounded-2xl border border-emerald-200/80 dark:border-emerald-800/60 shadow-2xs">
            <CheckCircle2 size={16} className="text-emerald-500" />
            <span>Урок уже сдан • Подтвержденный балл: {initialProgress.score ?? 100}%</span>
          </div>
        ) : (
          <div className="text-xs text-slate-400">
            Заполните ответы на задания. Итоговая проверка будет выполнена на сервере.
          </div>
        )}

        <button
          onClick={handleFinish}
          disabled={isVerifying}
          className="px-6 py-3 bg-emerald-600 hover:bg-emerald-700 disabled:opacity-60 text-white text-xs font-bold rounded-2xl shadow-lg shadow-emerald-600/20 flex items-center gap-2 cursor-pointer transition-all active:scale-95 disabled:cursor-not-allowed"
        >
          {isVerifying ? (
            <>
              <Loader2 size={16} className="animate-spin" />
              <span>Проверка ответов на сервере...</span>
            </>
          ) : (
            <>
              <Sparkles size={16} />
              <span>
                {initialProgress && (initialProgress.status === 'completed' || initialProgress.Status === 'completed')
                  ? 'Обновить результат и проверить на сервере'
                  : 'Завершить урок и отправить на проверку'}
              </span>
            </>
          )}
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
              <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-emerald-50 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-300 text-[11px] font-bold border border-emerald-200 dark:border-emerald-800/60 mx-auto">
                <ShieldCheck size={13} />
                <span>Серверная валидация пройдена</span>
              </div>
              <h3 className="text-lg font-black text-slate-900 dark:text-white">
                🎉 Урок успешно завершен!
              </h3>
              <p className="text-xs text-slate-500 dark:text-slate-400">
                Ваш прогресс и результаты тестирования проверены и зафиксированы на сервере
              </p>
            </div>

            {/* Score Card */}
            {completionResult.hasQuizzes ? (
              <div className="p-4 bg-slate-50 dark:bg-slate-800/60 rounded-2xl border border-slate-200 dark:border-slate-700 space-y-2">
                <div className="text-[11px] uppercase font-extrabold text-slate-400 tracking-wider">
                  Подтвержденный результат
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
                    ? '👍 Хорошая работа! Рекомендуем повторить темы с неверными ответами.'
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
                Просмотреть результаты на странице
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
