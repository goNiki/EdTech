'use client';

import React, { useState, useMemo, useEffect } from 'react';
import {
  CheckCircle,
  CheckCircle2,
  HelpCircle,
  Sparkles,
  ArrowRight,
  ArrowLeft,
  Award,
  BookOpen,
  Loader2,
  Check,
  RotateCcw,
  Upload,
  FileText,
  Clock,
  ShieldCheck,
  ChevronRight,
  ChevronLeft,
  Zap,
  EyeOff
} from 'lucide-react';
import { parseSmartDropdownTemplate } from '@/lib/puck-config';
import { api } from '@/lib/api';
import { LessonCompletionPayload, LessonAnswerItem } from './PuckLessonViewer';
import { QuizSettings, defaultQuizSettings } from '@/components/teacher/ModalQuizSettings';

export interface QuizStepperPlayerProps {
  contentJson: string | object;
  onComplete?: (payload: LessonCompletionPayload) => Promise<any> | void;
  onNavigateBack?: () => void;
  initialProgress?: any;
  lessonTitle?: string;
  quizSettings?: QuizSettings;
}

export default function QuizStepperPlayer({
  contentJson,
  onComplete,
  onNavigateBack,
  initialProgress,
  lessonTitle,
  quizSettings,
}: QuizStepperPlayerProps) {
  let parsedContent: any = { content: [] };
  try {
    if (typeof contentJson === 'string') {
      parsedContent = JSON.parse(contentJson);
    } else if (typeof contentJson === 'object' && contentJson !== null) {
      parsedContent = contentJson;
    }
  } catch (e) {
    console.error('Failed to parse Puck content in QuizStepperPlayer', e);
  }

  const allBlocks: any[] = parsedContent.content || [];

  // Separate theoretical intro blocks and question blocks
  const theoryBlocks = useMemo(() => {
    return allBlocks.filter((b: any) =>
      ['HeaderBlock', 'TextBlock', 'RichTextBlock', 'VideoBlock', 'PresentationBlock'].includes(b.type)
    );
  }, [allBlocks]);

  const questionBlocks = useMemo(() => {
    return allBlocks.filter((b: any) =>
      b.type?.startsWith('Quiz') || b.type === 'FileUploadBlock'
    );
  }, [allBlocks]);

  // Steps definition: Intro step (if theory exists) + Question steps
  const steps = useMemo(() => {
    const list: Array<{ type: 'theory' | 'question'; data: any; questionIndex?: number }> = [];
    if (theoryBlocks.length > 0) {
      list.push({ type: 'theory', data: theoryBlocks });
    }
    questionBlocks.forEach((qb, idx) => {
      list.push({ type: 'question', data: qb, questionIndex: idx + 1 });
    });
    return list;
  }, [theoryBlocks, questionBlocks]);

  // Current active step index (0 .. steps.length - 1)
  const [currentStepIdx, setCurrentStepIdx] = useState<number>(0);

  // Answers State
  const [singleAnswers, setSingleAnswers] = useState<Record<string, number>>({});
  const [multiAnswers, setMultiAnswers] = useState<Record<string, number[]>>({});
  const [matchAnswers, setMatchAnswers] = useState<Record<string, Record<number, string>>>({});
  const [dropdownAnswers, setDropdownAnswers] = useState<Record<string, Record<string, string>>>({});
  const [inputAnswers, setInputAnswers] = useState<Record<string, Record<string, string>>>({});
  const [sequenceOrders, setSequenceOrders] = useState<Record<string, string[]>>({});
  const [essayAnswers, setEssayAnswers] = useState<Record<string, string>>({});
  const [uploadedFiles, setUploadedFiles] = useState<
    Record<string, { name: string; size: number; url: string }>
  >({});
  const [isUploading, setIsUploading] = useState<Record<string, boolean>>({});

  // Verification & Completion state
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [completionResult, setCompletionResult] = useState<{
    score: number;
    earnedPoints: number;
    totalMaxPoints: number;
    hasQuizzes: boolean;
  } | null>(null);

  // Check if a question step is answered
  const isQuestionAnswered = (block: any): boolean => {
    if (!block || !block.props) return false;
    const bid = block.props.id || 'default';
    switch (block.type) {
      case 'QuizSingleBlock':
        return singleAnswers[bid] !== undefined;
      case 'QuizMultiBlock':
        return (multiAnswers[bid] || []).length > 0;
      case 'QuizMatchBlock':
        return Object.keys(matchAnswers[bid] || {}).length > 0;
      case 'QuizDropdownBlankBlock':
        return Object.keys(dropdownAnswers[bid] || {}).length > 0;
      case 'QuizInputBlankBlock':
        return Object.keys(inputAnswers[bid] || {}).length > 0;
      case 'QuizSequenceBlock':
        return true;
      case 'QuizEssayBlock':
        return (essayAnswers[bid] || '').trim().length > 0;
      case 'FileUploadBlock':
        return uploadedFiles[bid] !== undefined;
      default:
        return false;
    }
  };

  const answeredQuestionsCount = useMemo(() => {
    return questionBlocks.filter(isQuestionAnswered).length;
  }, [questionBlocks, singleAnswers, multiAnswers, matchAnswers, dropdownAnswers, inputAnswers, sequenceOrders, essayAnswers, uploadedFiles]);

  const progressPercent = questionBlocks.length > 0
    ? Math.round((answeredQuestionsCount / questionBlocks.length) * 100)
    : 100;

  // File upload handler
  const handleFileUpload = async (blockId: string, file: File) => {
    setIsUploading((prev) => ({ ...prev, [blockId]: true }));
    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('category', 'homework');

      const res = await api.post('/upload', formData, {
        params: { category: 'homework' },
        headers: { 'Content-Type': 'multipart/form-data' },
      });

      const data = res.data?.data || res.data;
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
    } catch (e) {
      console.error('File upload failed', e);
      alert('Ошибка при загрузке файла');
    } finally {
      setIsUploading((prev) => ({ ...prev, [blockId]: false }));
    }
  };

  // Calculate final score & payload
  const calculatePayload = (): { payload: LessonCompletionPayload; earnedPoints: number; totalMaxPoints: number } => {
    let totalMaxPoints = 0;
    let earnedPoints = 0;
    const submittedEssays: LessonCompletionPayload['essays'] = [];
    const answersList: LessonAnswerItem[] = [];

    questionBlocks.forEach((block: any, idx: number) => {
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
          earnedPoints += pts;
          break;
        }

        case 'QuizEssayBlock': {
          const answer = (essayAnswers[blockId] || '').trim();
          submittedEssays.push({
            question_text: props.question || 'Эссе',
            answer_text: answer || 'Ответ не был дан',
            max_points: Number(props.points) || 25,
          });
          break;
        }

        case 'FileUploadBlock': {
          const uploaded = uploadedFiles[blockId];
          submittedEssays.push({
            question_text: props.title || 'Практическое задание',
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
    setIsSubmitting(true);
    try {
      const { payload, earnedPoints, totalMaxPoints } = calculatePayload();
      let serverResponse: any = null;
      if (onComplete) {
        serverResponse = await onComplete(payload);
      }

      setCompletionResult({
        score: serverResponse?.score ?? payload.score,
        earnedPoints: serverResponse?.earned_points ?? earnedPoints,
        totalMaxPoints: serverResponse?.total_max_points ?? totalMaxPoints,
        hasQuizzes: totalMaxPoints > 0,
      });
    } catch (e) {
      console.error('Failed to submit in Stepper', e);
      alert('Ошибка при проверке результатов');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Quiz Settings & Timers
  const effectiveQuizSettings: QuizSettings = useMemo(() => {
    return quizSettings || parsedContent.quiz_settings || defaultQuizSettings;
  }, [quizSettings, parsedContent]);

  const [remainingOverallSeconds, setRemainingOverallSeconds] = useState<number>(() => {
    return (effectiveQuizSettings.time_limit_minutes || 0) * 60;
  });

  const [questionRemainingSeconds, setQuestionRemainingSeconds] = useState<number>(() => {
    return effectiveQuizSettings.question_time_limit_seconds || 0;
  });

  // Overall quiz timer countdown
  useEffect(() => {
    if (!effectiveQuizSettings.time_limit_minutes || effectiveQuizSettings.time_limit_minutes <= 0) return;
    
    const interval = setInterval(() => {
      setRemainingOverallSeconds((prev) => {
        if (prev <= 1) {
          clearInterval(interval);
          handleFinish();
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [effectiveQuizSettings.time_limit_minutes]);

  // Per-question blitz timer countdown
  useEffect(() => {
    const blitzSecs = effectiveQuizSettings.question_time_limit_seconds;
    if (!blitzSecs || blitzSecs <= 0) return;
    if (steps[currentStepIdx]?.type !== 'question') return;

    setQuestionRemainingSeconds(blitzSecs);

    const interval = setInterval(() => {
      setQuestionRemainingSeconds((prev) => {
        if (prev <= 1) {
          clearInterval(interval);
          if (currentStepIdx < steps.length - 1) {
            setCurrentStepIdx((idx) => idx + 1);
          } else {
            handleFinish();
          }
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [currentStepIdx, effectiveQuizSettings.question_time_limit_seconds, steps.length]);

  const formatMMSS = (totalSeconds: number) => {
    const m = Math.floor(Math.max(0, totalSeconds) / 60);
    const s = Math.max(0, totalSeconds) % 60;
    return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
  };

  const isBlindMode = effectiveQuizSettings.feedback_mode === 'exam_blind';

  const currentStep = steps[currentStepIdx];
  const isLastStep = currentStepIdx === steps.length - 1;
  const isFirstStep = currentStepIdx === 0;

  if (steps.length === 0) {
    return (
      <div className="p-12 text-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 space-y-3">
        <FileText size={36} className="mx-auto text-slate-400" />
        <h4 className="text-base font-bold text-slate-700 dark:text-slate-300">
          Содержимое теста пока формируется
        </h4>
        <p className="text-xs text-slate-500">
          Преподаватель еще не добавил вопросы в этот урок.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6 max-w-3xl mx-auto w-full">
      {/* Top Wizard Status & Numbers Bar */}
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-5 shadow-sm space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-3 text-xs">
          <div className="flex items-center gap-2 font-bold text-slate-700 dark:text-slate-200">
            <span className="w-2 h-2 rounded-full bg-indigo-500 animate-pulse" />
            <span>
              {currentStep.type === 'theory'
                ? 'Вводная теория к уроку'
                : `Вопрос ${currentStep.questionIndex} из ${questionBlocks.length}`}
            </span>
            {isBlindMode && (
              <span className="px-2 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider bg-purple-100 dark:bg-purple-950/70 text-purple-700 dark:text-purple-300 border border-purple-200 dark:border-purple-800 flex items-center gap-1">
                <EyeOff size={11} />
                <span>Blind Mode</span>
              </span>
            )}
          </div>

          <div className="flex items-center gap-3 text-slate-500 text-[11px] font-semibold">
            {effectiveQuizSettings.time_limit_minutes > 0 && (
              <div
                className={`flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-black tracking-wider transition-all ${
                  remainingOverallSeconds < 120
                    ? 'bg-rose-100 text-rose-700 dark:bg-rose-950 dark:text-rose-300 animate-pulse border border-rose-300 dark:border-rose-800'
                    : 'bg-indigo-50 text-indigo-700 dark:bg-indigo-950/80 dark:text-indigo-300 border border-indigo-200/60 dark:border-indigo-800/60'
                }`}
                title="Оставшееся время на весь тест"
              >
                <Clock size={13} />
                <span>{formatMMSS(remainingOverallSeconds)}</span>
              </div>
            )}
            <span>
              Отвечено: <strong className="text-emerald-600 dark:text-emerald-400">{answeredQuestionsCount}</strong> из {questionBlocks.length}
            </span>
            <span className="text-slate-300 dark:text-slate-700">•</span>
            <span>{progressPercent}%</span>
          </div>
        </div>

        {/* Progress bar */}
        <div className="w-full bg-slate-100 dark:bg-slate-800 h-2 rounded-full overflow-hidden">
          <div
            className="h-full bg-gradient-to-r from-indigo-500 to-emerald-500 transition-all duration-300 rounded-full"
            style={{ width: `${progressPercent}%` }}
          />
        </div>

        {/* Questions Number Grid / Tiles */}
        <div className="flex items-center gap-2 overflow-x-auto pb-1 pt-1 scrollbar-none">
          {theoryBlocks.length > 0 && (
            <button
              type="button"
              onClick={() => setCurrentStepIdx(0)}
              className={`px-3 py-1.5 rounded-xl text-xs font-bold transition-all cursor-pointer flex-shrink-0 flex items-center gap-1.5 ${
                currentStepIdx === 0
                  ? 'bg-indigo-600 text-white shadow-md ring-2 ring-indigo-500/20'
                  : 'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700'
              }`}
            >
              <BookOpen size={13} />
              <span>Теория</span>
            </button>
          )}

          {questionBlocks.map((qBlock, qIdx) => {
            const stepNumber = theoryBlocks.length > 0 ? qIdx + 1 : qIdx;
            const isCurrent = currentStepIdx === stepNumber;
            const isAnswered = isQuestionAnswered(qBlock);

            return (
              <button
                key={qIdx}
                type="button"
                onClick={() => setCurrentStepIdx(stepNumber)}
                className={`w-9 h-9 rounded-xl text-xs font-bold transition-all flex items-center justify-center flex-shrink-0 cursor-pointer ${
                  isCurrent
                    ? 'bg-indigo-600 text-white ring-4 ring-indigo-500/20 shadow-md scale-105'
                    : isAnswered
                    ? 'bg-emerald-500 text-white shadow-xs'
                    : 'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 hover:bg-slate-200 dark:hover:bg-slate-700'
                }`}
                title={`Перейти к вопросу ${qIdx + 1}${isAnswered ? ' (Ответ дан)' : ''}`}
              >
                {isAnswered && !isCurrent ? (
                  <Check size={14} className="stroke-[3]" />
                ) : (
                  <span>{qIdx + 1}</span>
                )}
              </button>
            );
          })}
        </div>
      </div>

      {/* Main Step Card */}
      <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 sm:p-8 shadow-xl min-h-[380px] flex flex-col justify-between transition-all">
        {/* Case 1: Theory Step */}
        {currentStep.type === 'theory' ? (
          <div className="space-y-6">
            <div className="flex items-center gap-2">
              <span className="px-3 py-1 rounded-xl text-[10px] font-black uppercase tracking-wider bg-indigo-50 dark:bg-indigo-950/60 text-indigo-700 dark:text-indigo-300 border border-indigo-200 dark:border-indigo-800">
                Теоретическая часть
              </span>
            </div>

            <div className="space-y-4">
              {currentStep.data.map((block: any, idx: number) => {
                if (block.type === 'HeaderBlock') {
                  return (
                    <div key={idx} className="space-y-1">
                      <h2 className="text-xl sm:text-2xl font-black text-slate-900 dark:text-white">
                        {block.props?.title}
                      </h2>
                      {block.props?.subtitle && (
                        <p className="text-sm text-slate-500 dark:text-slate-400">
                          {block.props.subtitle}
                        </p>
                      )}
                    </div>
                  );
                }
                if (block.type === 'TextBlock' || block.type === 'RichTextBlock') {
                  return (
                    <div
                      key={idx}
                      className="prose dark:prose-invert max-w-none text-slate-700 dark:text-slate-300 leading-relaxed text-sm sm:text-base"
                      dangerouslySetInnerHTML={{
                        __html: block.props?.contentHtml || block.props?.content || '',
                      }}
                    />
                  );
                }
                return null;
              })}
            </div>
          </div>
        ) : (
          /* Case 2: Interactive Question Step */
          (() => {
            const block = currentStep.data;
            const { type, props = {} } = block;
            const bid = props.id || `q-${currentStep.questionIndex}`;

            return (
              <div className="space-y-6">
                {/* Blitz Timer Progress Bar */}
                {effectiveQuizSettings.question_time_limit_seconds > 0 && (
                  <div className="space-y-1 p-3 rounded-2xl bg-amber-50/60 dark:bg-amber-950/30 border border-amber-200/80 dark:border-amber-900/60 animate-in fade-in">
                    <div className="flex items-center justify-between text-[11px] font-bold">
                      <span className="text-amber-700 dark:text-amber-300 flex items-center gap-1.5">
                        <Zap size={13} className="fill-amber-500 text-amber-500" />
                        <span>Блиц-таймер на вопрос</span>
                      </span>
                      <span className="font-mono text-amber-800 dark:text-amber-200 font-extrabold">
                        {questionRemainingSeconds} сек.
                      </span>
                    </div>
                    <div className="w-full bg-amber-200/60 dark:bg-slate-800 h-1.5 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-gradient-to-r from-amber-500 to-rose-500 transition-all duration-1000 ease-linear rounded-full"
                        style={{
                          width: `${Math.max(
                            0,
                            (questionRemainingSeconds / (effectiveQuizSettings.question_time_limit_seconds || 1)) * 100
                          )}%`,
                        }}
                      />
                    </div>
                  </div>
                )}

                {/* Header: Question badge, points */}
                <div className="flex items-center justify-between gap-3">
                  <div className="flex items-center gap-2">
                    <span className="w-6 h-6 rounded-full bg-indigo-600 text-white font-black text-xs flex items-center justify-center shadow-xs">
                      {currentStep.questionIndex}
                    </span>
                    <span className="px-2.5 py-0.5 rounded-lg text-[10px] font-extrabold uppercase tracking-wider bg-indigo-50 dark:bg-indigo-950/60 text-indigo-700 dark:text-indigo-300 border border-indigo-200/80 dark:border-indigo-800/60">
                      {type === 'QuizSingleBlock'
                        ? 'Одиночный выбор'
                        : type === 'QuizMultiBlock'
                        ? 'Множественный выбор'
                        : type === 'QuizMatchBlock'
                        ? 'Сопоставление пар'
                        : type === 'QuizDropdownBlankBlock'
                        ? 'Заполнение пропусков'
                        : type === 'QuizInputBlankBlock'
                        ? 'Ввод слова'
                        : type === 'QuizSequenceBlock'
                        ? 'Последовательность'
                        : type === 'QuizEssayBlock'
                        ? 'Развернутый ответ'
                        : 'Загрузка файла'}
                    </span>
                    {isBlindMode && (
                      <span className="px-2 py-0.5 rounded-lg text-[10px] font-bold bg-purple-50 dark:bg-purple-950/60 text-purple-700 dark:text-purple-300 border border-purple-200/80 dark:border-purple-800/60 flex items-center gap-1">
                        <EyeOff size={11} />
                        <span>Blind Mode</span>
                      </span>
                    )}
                  </div>

                  <span className="text-xs font-bold text-slate-500 dark:text-slate-400">
                    {props.points || 10} баллов
                  </span>
                </div>

                {/* Question Prompt */}
                <h3 className="text-base sm:text-lg font-black text-slate-900 dark:text-white leading-snug">
                  {props.question || props.title || 'Ответьте на вопрос:'}
                </h3>

                {/* Subtitle / Hint */}
                {type === 'QuizMultiBlock' && (
                  <p className="text-xs text-indigo-600 dark:text-indigo-400 font-semibold">
                    💡 Выберите один или несколько правильных вариантов
                  </p>
                )}

                {/* Options / Inputs by block type */}
                <div className="space-y-3 pt-2">
                  {/* 1. QuizSingleBlock */}
                  {type === 'QuizSingleBlock' && (
                    <div className="space-y-2.5">
                      {(props.options || []).map((opt: any, optIdx: number) => {
                        const isSelected = singleAnswers[bid] === optIdx;
                        return (
                          <div
                            key={optIdx}
                            onClick={() =>
                              setSingleAnswers((prev) => ({ ...prev, [bid]: optIdx }))
                            }
                            className={`p-4 rounded-2xl border-2 flex items-center gap-3.5 cursor-pointer transition-all select-none ${
                              isSelected
                                ? 'border-emerald-500 bg-emerald-50/50 dark:bg-emerald-950/30 text-emerald-950 dark:text-emerald-100 shadow-sm ring-1 ring-emerald-500/20'
                                : 'border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-950 hover:border-slate-300 dark:hover:border-slate-700 text-slate-800 dark:text-slate-200'
                            }`}
                          >
                            <div
                              className={`w-5 h-5 rounded-full border-2 flex items-center justify-center flex-shrink-0 transition-colors ${
                                isSelected
                                  ? 'border-emerald-500 bg-emerald-500 text-white'
                                  : 'border-slate-300 dark:border-slate-700'
                              }`}
                            >
                              {isSelected && <Check size={12} className="stroke-[3]" />}
                            </div>
                            <span className="text-xs sm:text-sm font-bold flex-1">
                              {opt.text}
                            </span>
                          </div>
                        );
                      })}
                    </div>
                  )}

                  {/* 2. QuizMultiBlock */}
                  {type === 'QuizMultiBlock' && (
                    <div className="space-y-2.5">
                      {(props.options || []).map((opt: any, optIdx: number) => {
                        const currentList = multiAnswers[bid] || [];
                        const isSelected = currentList.includes(optIdx);
                        return (
                          <div
                            key={optIdx}
                            onClick={() => {
                              setMultiAnswers((prev) => {
                                const exist = prev[bid] || [];
                                const next = exist.includes(optIdx)
                                  ? exist.filter((i) => i !== optIdx)
                                  : [...exist, optIdx];
                                return { ...prev, [bid]: next };
                              });
                            }}
                            className={`p-4 rounded-2xl border-2 flex items-center gap-3.5 cursor-pointer transition-all select-none ${
                              isSelected
                                ? 'border-purple-500 bg-purple-50/50 dark:bg-purple-950/30 text-purple-950 dark:text-purple-100 shadow-sm ring-1 ring-purple-500/20'
                                : 'border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-950 hover:border-slate-300 dark:hover:border-slate-700 text-slate-800 dark:text-slate-200'
                            }`}
                          >
                            <div
                              className={`w-5 h-5 rounded-lg border-2 flex items-center justify-center flex-shrink-0 transition-colors ${
                                isSelected
                                  ? 'border-purple-600 bg-purple-600 text-white'
                                  : 'border-slate-300 dark:border-slate-700'
                              }`}
                            >
                              {isSelected && <Check size={12} className="stroke-[3]" />}
                            </div>
                            <span className="text-xs sm:text-sm font-bold flex-1">
                              {opt.text}
                            </span>
                          </div>
                        );
                      })}
                    </div>
                  )}

                  {/* 3. QuizMatchBlock */}
                  {type === 'QuizMatchBlock' && (
                    <div className="space-y-3">
                      {(props.pairs || []).map((pair: any, pIdx: number) => {
                        const currentVal = matchAnswers[bid]?.[pIdx] || '';
                        const allRightOptions = (props.pairs || []).map((p: any) => p.right).filter(Boolean);
                        return (
                          <div
                            key={pIdx}
                            className="p-3.5 rounded-2xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-950 flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-xs"
                          >
                            <span className="font-bold text-slate-800 dark:text-slate-200 sm:w-1/2">
                              {pair.left}
                            </span>
                            <select
                              value={currentVal}
                              onChange={(e) => {
                                const val = e.target.value;
                                setMatchAnswers((prev) => ({
                                  ...prev,
                                  [bid]: { ...(prev[bid] || {}), [pIdx]: val },
                                }));
                              }}
                              className="px-3 py-2 rounded-xl border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900 font-semibold text-slate-800 dark:text-slate-200 sm:w-1/2 focus:outline-indigo-500"
                            >
                              <option value="">Выберите соответствие...</option>
                              {allRightOptions.map((opt: string, oIdx: number) => (
                                <option key={oIdx} value={opt}>
                                  {opt}
                                </option>
                              ))}
                            </select>
                          </div>
                        );
                      })}
                    </div>
                  )}

                  {/* 4. QuizEssayBlock */}
                  {type === 'QuizEssayBlock' && (
                    <div className="space-y-2">
                      <textarea
                        rows={6}
                        value={essayAnswers[bid] || ''}
                        onChange={(e) =>
                          setEssayAnswers((prev) => ({ ...prev, [bid]: e.target.value }))
                        }
                        placeholder="Напишите ваш ответ здесь..."
                        className="w-full p-4 rounded-2xl border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-950 text-slate-900 dark:text-white text-xs sm:text-sm font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                      />
                      <div className="text-right text-[11px] text-slate-400 font-semibold">
                        Символов: {(essayAnswers[bid] || '').length}
                      </div>
                    </div>
                  )}

                  {/* 5. FileUploadBlock */}
                  {type === 'FileUploadBlock' && (
                    <div className="space-y-3">
                      {uploadedFiles[bid] ? (
                        <div className="p-4 rounded-2xl bg-emerald-50 dark:bg-emerald-950/30 border border-emerald-300 dark:border-emerald-800 flex items-center justify-between">
                          <div className="flex items-center gap-2.5">
                            <CheckCircle2 size={18} className="text-emerald-500" />
                            <div>
                              <p className="text-xs font-bold text-slate-800 dark:text-slate-200">
                                {uploadedFiles[bid].name}
                              </p>
                              <p className="text-[10px] text-slate-400">
                                {(uploadedFiles[bid].size / 1024).toFixed(1)} KB
                              </p>
                            </div>
                          </div>
                          <button
                            type="button"
                            onClick={() => {
                              const copy = { ...uploadedFiles };
                              delete copy[bid];
                              setUploadedFiles(copy);
                            }}
                            className="text-xs text-rose-500 hover:underline font-bold"
                          >
                            Заменить
                          </button>
                        </div>
                      ) : (
                        <label className="p-8 rounded-2xl border-2 border-dashed border-indigo-300 dark:border-indigo-800 hover:border-indigo-500 bg-indigo-50/20 dark:bg-indigo-950/20 flex flex-col items-center justify-center gap-2 cursor-pointer transition-colors text-center">
                          <Upload size={28} className="text-indigo-600 dark:text-indigo-400" />
                          <span className="text-xs font-bold text-slate-800 dark:text-slate-200">
                            {isUploading[bid] ? 'Загрузка файла...' : 'Нажмите для выбора файла решения'}
                          </span>
                          <span className="text-[10px] text-slate-400">
                            PDF, DOCX, ZIP, PNG (до 25 МБ)
                          </span>
                          <input
                            type="file"
                            className="hidden"
                            disabled={isUploading[bid]}
                            onChange={(e) => {
                              const f = e.target.files?.[0];
                              if (f) handleFileUpload(bid, f);
                            }}
                          />
                        </label>
                      )}
                    </div>
                  )}
                </div>
              </div>
            );
          })()
        )}

        {/* Bottom Wizard Stepper Navigation Buttons */}
        <div className="pt-8 border-t border-slate-100 dark:border-slate-800/80 flex items-center justify-between gap-3 mt-6">
          <button
            type="button"
            disabled={isFirstStep}
            onClick={() => setCurrentStepIdx((prev) => Math.max(0, prev - 1))}
            className="px-4 py-2.5 rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 text-slate-700 dark:text-slate-200 text-xs font-bold transition-all hover:bg-slate-50 dark:hover:bg-slate-800 disabled:opacity-30 disabled:cursor-not-allowed flex items-center gap-1.5 cursor-pointer shadow-2xs"
          >
            <ChevronLeft size={16} />
            <span>Назад</span>
          </button>

          {isLastStep ? (
            <button
              type="button"
              disabled={isSubmitting}
              onClick={handleFinish}
              className="px-6 py-2.5 rounded-2xl bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-black transition-all shadow-md shadow-emerald-600/25 flex items-center gap-2 cursor-pointer active:scale-95 disabled:opacity-60"
            >
              {isSubmitting ? (
                <>
                  <Loader2 size={16} className="animate-spin" />
                  <span>Проверка ответов...</span>
                </>
              ) : (
                <>
                  <CheckCircle size={16} />
                  <span>Завершить тест ✓</span>
                </>
              )}
            </button>
          ) : (
            <button
              type="button"
              onClick={() => setCurrentStepIdx((prev) => Math.min(steps.length - 1, prev + 1))}
              className="px-6 py-2.5 rounded-2xl bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-black transition-all shadow-md shadow-indigo-600/25 flex items-center gap-1.5 cursor-pointer active:scale-95"
            >
              <span>Следующий вопрос</span>
              <ChevronRight size={16} />
            </button>
          )}
        </div>
      </div>

      {/* Results Modal */}
      {completionResult && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4 animate-in fade-in duration-200">
          <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 max-w-md w-full shadow-2xl space-y-5 text-center">
            <div className="w-14 h-14 rounded-2xl bg-emerald-100 dark:bg-emerald-950 text-emerald-600 dark:text-emerald-400 flex items-center justify-center mx-auto shadow-inner">
              <Award size={32} />
            </div>

            <div className="space-y-1.5">
              <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-emerald-50 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-300 text-[11px] font-bold border border-emerald-200 dark:border-emerald-800/60 mx-auto">
                <ShieldCheck size={13} />
                <span>Тестирование завершено</span>
              </div>
              <h3 className="text-lg font-black text-slate-900 dark:text-white">
                🎉 Результат зафиксирован!
              </h3>
            </div>

            {/* Score Card */}
            <div className="p-4 bg-slate-50 dark:bg-slate-800/60 rounded-2xl border border-slate-200 dark:border-slate-700 space-y-2">
              <div className="text-[11px] uppercase font-extrabold text-slate-400 tracking-wider">
                Набранные баллы
              </div>
              <div className="text-3xl font-black text-indigo-600 dark:text-indigo-400">
                {completionResult.earnedPoints} / {completionResult.totalMaxPoints}
                <span className="text-xs text-slate-400 font-bold ml-1.5">
                  ({completionResult.score}%)
                </span>
              </div>
              <div className="text-[11px] text-slate-500 font-medium">
                {completionResult.score >= 80
                  ? '🌟 Превосходный результат! Тема отлично усвоена.'
                  : completionResult.score >= 50
                  ? '👍 Хорошая работа! Рекомендуется повторить сложные вопросы.'
                  : '📖 Рекомендуем вернуться к теории и пройти тест еще раз.'}
              </div>
            </div>

            <button
              type="button"
              onClick={() => {
                setCompletionResult(null);
                if (onNavigateBack) onNavigateBack();
              }}
              className="w-full py-3 rounded-2xl bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs shadow-md transition-all cursor-pointer"
            >
              Вернуться к курсу
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
