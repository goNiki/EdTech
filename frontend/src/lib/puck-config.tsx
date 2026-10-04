import type { Config, CustomField } from '@puckeditor/core';
import React, { useRef, useState } from 'react';
import { InlineText, RichTextCanvasEditor, RichTextWordEditor, usePuckPropUpdater } from '@/components/editor/InlineEditable';
import { Plus, Trash2, Presentation as PresentationIcon, Upload, Loader2 } from 'lucide-react';
import PresentationViewer from '@/components/player/PresentationViewer';
import { api } from './api';

export interface DropdownBlankItem {
  key: string;
  options: string;
  correctAnswer: string;
}

export interface InputBlankItem {
  key: string;
  correctAnswer: string;
  caseSensitive?: string;
}

export type PuckProps = {
  HeaderBlock: {
    title: string;
    subtitle?: string;
    level?: 'h1' | 'h2' | 'h3';
  };
  TextBlock: {
    content?: string;
    contentHtml?: string;
    align?: 'left' | 'center' | 'right' | 'justify';
  };
  RichTextBlock: {
    title?: string;
    contentHtml: string;
  };
  VideoBlock: {
    url: string;
    caption?: string;
  };
  QuizSingleBlock: {
    question: string;
    points: number;
    options: Array<{ text: string; isCorrect: string | boolean; explain?: string }>;
  };
  QuizMultiBlock: {
    question: string;
    points: number;
    options: Array<{ text: string; isCorrect: string | boolean }>;
  };
  QuizMatchBlock: {
    question: string;
    points: number;
    pairs: Array<{ left: string; right: string }>;
  };
  QuizDropdownBlankBlock: {
    question: string;
    points: number;
    templateText: string;
    blanks?: Array<DropdownBlankItem>;
    // Backward compatibility props
    textBefore?: string;
    textAfter?: string;
    options?: Array<{ text: string }>;
    correctIndex?: number;
  };
  QuizInputBlankBlock: {
    question: string;
    points: number;
    templateText: string;
    blanks?: Array<InputBlankItem>;
    // Backward compatibility props
    prefixText?: string;
    suffixText?: string;
    correctAnswer?: string;
    caseSensitive?: string;
  };
  QuizSequenceBlock: {
    question: string;
    points: number;
    items: Array<{ text: string }>;
  };
  QuizEssayBlock: {
    question: string;
    points: number;
    rubric?: string;
    sampleAnswer?: string;
  };
  FileUploadBlock: {
    title: string;
    instructions: string;
    points: number;
    allowedTypes?: string;
    maxSizeMB?: number;
  };
  PresentationBlock: {
    title?: string;
    mode?: 'embed' | 'pdf';
    embedUrl?: string;
    pdfUrl?: string;
    aspectRatio?: '16:9' | '4:3';
    allowDownload?: boolean;
  };
};

export interface ParsedDropdownBlank {
  key: string;
  rawTag: string;
  correctAnswer: string;
  options: string[];
}

export interface ParsedTemplateToken {
  type: 'text' | 'blank';
  value: string;
  blank?: ParsedDropdownBlank;
}

export function parseSmartDropdownTemplate(
  templateText: string,
  existingBlanks?: Array<DropdownBlankItem | InputBlankItem>
): { tokens: ParsedTemplateToken[]; blanks: ParsedDropdownBlank[] } {
  if (!templateText) {
    return { tokens: [], blanks: [] };
  }

  const regex = /\{([^{}]+)\}/g;
  const tokens: ParsedTemplateToken[] = [];
  const blanks: ParsedDropdownBlank[] = [];

  let lastIndex = 0;
  let match: RegExpExecArray | null;
  let blankCounter = 1;

  while ((match = regex.exec(templateText)) !== null) {
    if (match.index > lastIndex) {
      tokens.push({
        type: 'text',
        value: templateText.substring(lastIndex, match.index),
      });
    }

    const rawInner = match[1].trim();
    const rawTag = match[0];

    let correctAnswer = '';
    let distractors: string[] = [];

    if (rawInner.includes(';')) {
      const parts = rawInner.split(';');
      correctAnswer = parts[0]?.trim() || '';
      if (parts[1]) {
        distractors = parts[1]
          .split(',')
          .map((s) => s.trim())
          .filter(Boolean);
      }
    } else {
      correctAnswer = rawInner;
    }

    const key = `blank_${blankCounter++}`;
    const allOptions = Array.from(new Set([correctAnswer, ...distractors])).filter(Boolean);

    const blankObj: ParsedDropdownBlank = {
      key,
      rawTag,
      correctAnswer,
      options: allOptions.length > 0 ? allOptions : [correctAnswer],
    };

    blanks.push(blankObj);
    tokens.push({
      type: 'blank',
      value: rawTag,
      blank: blankObj,
    });

    lastIndex = regex.lastIndex;
  }

  if (lastIndex < templateText.length) {
    tokens.push({
      type: 'text',
      value: templateText.substring(lastIndex),
    });
  }

  return { tokens, blanks };
}

export function parseTemplateParts(text: string): Array<{ type: 'text' | 'blank'; value: string; key?: string }> {
  if (!text) return [];
  const regex = /\{([^}]+)\}/g;
  const parts: Array<{ type: 'text' | 'blank'; value: string; key?: string }> = [];
  let lastIndex = 0;
  let match: RegExpExecArray | null;

  while ((match = regex.exec(text)) !== null) {
    if (match.index > lastIndex) {
      parts.push({
        type: 'text',
        value: text.substring(lastIndex, match.index),
      });
    }

    parts.push({
      type: 'blank',
      value: match[0],
      key: match[1].trim(),
    });
    lastIndex = regex.lastIndex;
  }

  if (lastIndex < text.length) {
    parts.push({
      type: 'text',
      value: text.substring(lastIndex),
    });
  }

  return parts;
}

export const config: Config<PuckProps> = {
  components: {
    HeaderBlock: {
      label: 'Заголовок',
      fields: {
        title: { type: 'text', label: 'Основной заголовок' },
        subtitle: { type: 'textarea', label: 'Подзаголовок (опционально)' },
        level: {
          type: 'select',
          label: 'Уровень заголовка',
          options: [
            { label: 'H1 (Главный)', value: 'h1' },
            { label: 'H2 (Раздел)', value: 'h2' },
            { label: 'H3 (Подраздел)', value: 'h3' },
          ],
        },
      },
      defaultProps: {
        title: 'Архитектура сервисов и проектирование',
        subtitle: 'Интерактивный конспект и практические задания к уроку',
        level: 'h1',
      },
      render: (props: any) => {
        const { title, subtitle, level, id } = props;
        const { updateProp, isEditing, selectThisBlock } = usePuckPropUpdater(id);
        const HeadingTag = level || 'h1';
        const sizeClasses =
          level === 'h1'
            ? 'text-3xl lg:text-4xl font-extrabold text-slate-900 dark:text-white'
            : level === 'h2'
            ? 'text-2xl font-bold text-slate-800 dark:text-slate-100'
            : 'text-xl font-bold text-slate-800 dark:text-slate-200';

        return (
          <div className="my-6 space-y-2" id={id} onClick={() => selectThisBlock()}>
            <HeadingTag className={sizeClasses}>
              <InlineText
                value={title}
                onChange={(val) => updateProp('title', val)}
                onFocusBlock={selectThisBlock}
                isEditing={isEditing}
                placeholder="Введите основной заголовок..."
              />
            </HeadingTag>
            <div className="text-sm text-slate-500 dark:text-slate-400 leading-relaxed max-w-3xl">
              <InlineText
                value={subtitle || ''}
                onChange={(val) => updateProp('subtitle', val)}
                onFocusBlock={selectThisBlock}
                isEditing={isEditing}
                placeholder="Введите подзаголовок (опционально)..."
                multiline
              />
            </div>
          </div>
        );
      },
    },

    TextBlock: {
      label: 'Текстовый блок (Canvas)',
      fields: {
        contentHtml: { type: 'textarea', label: 'HTML-содержимое блока' },
        align: {
          type: 'select',
          label: 'Выравнивание',
          options: [
            { label: 'По левому краю', value: 'left' },
            { label: 'По центру', value: 'center' },
            { label: 'По правому краю', value: 'right' },
            { label: 'По ширине', value: 'justify' },
          ],
        },
      },
      defaultProps: {
        contentHtml:
          '<p>В парадигме <strong>Content-as-Data</strong> весь контент описывается строго типизированным JSON-объектом, что исключает XSS и обеспечивает переиспользуемость компонентов.</p>',
        content:
          'В парадигме Content-as-Data весь контент описывается строго типизированным JSON-объектом, что исключает XSS и обеспечивает переиспользуемость компонентов.',
        align: 'left',
      },
      render: (props: any) => {
        const { contentHtml, content, align, id, puck } = props;
        const { updateProps, isEditing, selectThisBlock } = usePuckPropUpdater(id);
        const effectiveEditing = puck?.isEditing ?? isEditing ?? true;
        const initialHtml = contentHtml || (content ? `<p>${content}</p>` : '<p>Начните вводить текст лекции...</p>');

        return (
          <div className="my-3 w-full" id={id} onClick={() => selectThisBlock()}>
            <RichTextCanvasEditor
              htmlContent={initialHtml}
              onChange={(val) => {
                updateProps({
                  contentHtml: val,
                  content: val.replace(/<[^>]+>/g, ' '),
                });
              }}
              isEditing={effectiveEditing}
              onFocusBlock={selectThisBlock}
              defaultAlign={align || 'left'}
            />
          </div>
        );
      },
    },

    RichTextBlock: {
      label: 'Лекция / Статья (Word)',
      fields: {
        title: { type: 'text', label: 'Заголовок раздела лекции' },
        contentHtml: { type: 'textarea', label: 'HTML-содержимое статьи' },
      },
      defaultProps: {
        title: 'Теоретический материал к уроку',
        contentHtml: `
          <p>Добро пожаловать в интерактивный конспект лекции! Вы можете писать и форматировать этот текст прямо здесь, как в Microsoft Word или Google Docs.</p>
          <p><strong>Ключевые преимущества Clean Architecture:</strong></p>
          <ul>
            <li>Независимость от фреймворков и библиотек</li>
            <li>Легкая тестируемость бизнес-логики без поднятия БД</li>
            <li>Единый доменный язык (Ubiquitous Language)</li>
          </ul>
          <div class="my-5 p-4 rounded-2xl bg-indigo-50/80 dark:bg-indigo-950/40 border-l-4 border-indigo-500 flex items-start gap-3">
            <span class="text-xl">💡</span>
            <div class="text-sm text-slate-800 dark:text-slate-200">
              <strong>Совет преподавателя:</strong> нажимайте кнопки на панели Word сверху, чтобы вставлять картинки, списки, цитаты и цветные врезки.
            </div>
          </div>
        `,
      },
      render: (props: any) => {
        const { contentHtml, title, id, puck } = props;
        const { updateProp, isEditing, selectThisBlock } = usePuckPropUpdater(id);
        const effectiveEditing = puck?.isEditing ?? isEditing ?? true;

        return (
          <div className="my-5 w-full" id={id} onClick={() => selectThisBlock()}>
            <RichTextWordEditor
              htmlContent={contentHtml || ''}
              onChange={(val) => updateProp('contentHtml', val)}
              title={title}
              onTitleChange={(val) => updateProp('title', val)}
              isEditing={effectiveEditing}
              onFocusBlock={selectThisBlock}
            />
          </div>
        );
      },
    },

    VideoBlock: {
      label: 'Видео-урок',
      fields: {
        url: { type: 'text', label: 'Ссылка на видео (YouTube / MP4 URL)' },
        caption: { type: 'text', label: 'Подпись к видео (опционально)' },
      },
      defaultProps: {
        url: 'https://www.youtube.com/embed/dQw4w9WgXcQ',
        caption: 'Видео-объяснение темы урока',
      },
      render: (props: any) => {
        const { url, caption, id } = props;
        const { updateProp, isEditing, selectThisBlock } = usePuckPropUpdater(id);

        return (
          <div className="my-6 space-y-2" id={id} onClick={() => selectThisBlock()}>
            <div className="relative aspect-video rounded-3xl overflow-hidden bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-md">
              <iframe
                src={url}
                title={caption || 'Видео'}
                className="w-full h-full border-0"
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                allowFullScreen
              />
            </div>
            <div className="text-xs text-center text-slate-500 dark:text-slate-400 font-medium">
              <InlineText
                value={caption || ''}
                onChange={(val) => updateProp('caption', val)}
                onFocusBlock={selectThisBlock}
                isEditing={isEditing}
                placeholder="Добавьте подпись к видео..."
              />
            </div>
          </div>
        );
      },
    },

    PresentationBlock: {
      label: 'Презентация / Слайды',
      fields: {
        title: { type: 'text', label: 'Название презентации' },
        mode: {
          type: 'select',
          label: 'Режим отображения',
          options: [
            { label: 'Встраивание по ссылке (Google Slides / Canva)', value: 'embed' },
            { label: 'Загрузка PDF-файла', value: 'pdf' },
          ],
        },
        embedUrl: { type: 'text', label: 'Ссылка на слайды (Google Slides, Canva, SpeakerDeck)' },
        pdfUrl: { type: 'text', label: 'URL загруженного PDF документа' },
        aspectRatio: {
          type: 'select',
          label: 'Соотношение сторон',
          options: [
            { label: '16:9 (Широкоформатный)', value: '16:9' },
            { label: '4:3 (Классический)', value: '4:3' },
          ],
        },
        allowDownload: {
          type: 'select',
          label: 'Разрешить скачивание',
          options: [
            { label: 'Да', value: true },
            { label: 'Нет', value: false },
          ],
        },
      },
      defaultProps: {
        title: 'Презентация к лекции',
        mode: 'embed',
        embedUrl: 'https://docs.google.com/presentation/d/e/2PACX-1vT_DEMO_SLIDES/embed',
        pdfUrl: '',
        aspectRatio: '16:9',
        allowDownload: true,
      },
      render: (props: any) => {
        const { title, mode, embedUrl, pdfUrl, aspectRatio, allowDownload, id } = props;
        const { updateProp, isEditing, selectThisBlock } = usePuckPropUpdater(id);

        return (
          <div className="my-6 space-y-3" id={id} onClick={() => selectThisBlock()}>
            {isEditing && (
              <div
                data-puck-overlay-portal="true"
                onPointerDown={(e) => e.stopPropagation()}
                onMouseDown={(e) => e.stopPropagation()}
                className="p-3 bg-amber-50/80 dark:bg-amber-950/40 border border-amber-200 dark:border-amber-800 rounded-2xl flex flex-wrap items-center justify-between gap-3 text-xs"
              >
                <div className="flex items-center gap-2 font-bold text-amber-900 dark:text-amber-200">
                  <PresentationIcon size={16} />
                  <span>Настройки презентации:</span>
                </div>

                <div className="flex items-center gap-2">
                  <select
                    value={mode || 'embed'}
                    onChange={(e) => updateProp('mode', e.target.value)}
                    className="px-2.5 py-1 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl font-semibold cursor-pointer"
                  >
                    <option value="embed">Ссылка (Google Slides / Canva)</option>
                    <option value="pdf">Файл PDF</option>
                  </select>

                  {mode === 'pdf' && (
                    <label className="px-3 py-1 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl flex items-center gap-1.5 cursor-pointer transition-colors shadow-xs">
                      <Upload size={13} />
                      <span>{pdfUrl ? 'Заменить PDF' : 'Загрузить PDF'}</span>
                      <input
                        type="file"
                        accept="application/pdf"
                        className="hidden"
                        onChange={async (e) => {
                          const file = e.target.files?.[0];
                          if (!file) return;
                          const formData = new FormData();
                          formData.append('file', file);
                          formData.append('category', 'presentation');
                          try {
                            const res = await api.post('/upload', formData, {
                              headers: { 'Content-Type': 'multipart/form-data' },
                            });
                            const fileUrl = res.data?.data?.file_url || res.data?.file_url;
                            if (fileUrl) {
                              updateProp('pdfUrl', fileUrl);
                            }
                          } catch {
                            const localUrl = URL.createObjectURL(file);
                            updateProp('pdfUrl', localUrl);
                          }
                          e.target.value = '';
                        }}
                      />
                    </label>
                  )}
                </div>
              </div>
            )}

            <PresentationViewer
              title={title}
              mode={mode}
              embedUrl={embedUrl}
              pdfUrl={pdfUrl}
              aspectRatio={aspectRatio}
              allowDownload={allowDownload}
            />
          </div>
        );
      },
    },

    QuizSingleBlock: {
      label: 'Тест: один ответ',
      fields: {
        question: { type: 'textarea', label: 'Текст вопроса' },
        points: { type: 'number', label: 'Максимальный балл' },
        options: {
          type: 'array',
          label: 'Варианты ответов',
          getItemSummary: (item: any) => item.text || 'Новый вариант',
          arrayFields: {
            text: { type: 'text', label: 'Текст варианта' },
            isCorrect: {
              type: 'select',
              label: 'Правильный ответ?',
              options: [
                { label: 'Нет (Неверный)', value: 'false' },
                { label: 'Да (Правильный)', value: 'true' },
              ],
            },
            explain: { type: 'text', label: 'Пояснение к ответу (опционально)' },
          },
        },
      },
      defaultProps: {
        question: '1. Какой паттерн используется для изоляции бизнес-логики от внешних фреймворков?',
        points: 10,
        options: [
          { text: 'Clean Architecture (Чистая Архитектура)', isCorrect: 'true', explain: 'Верно! Бизнес-логика живет в чистом домене.' },
          { text: 'Smart UI Pattern', isCorrect: 'false', explain: 'Неверно, логика смешивается с представлением.' },
          { text: 'Active Record Pattern', isCorrect: 'false', explain: 'Привязывает домен напрямую к структуре БД.' },
        ],
      },
      render: (props: any) => {
        const { question, points, options, id } = props;
        const { updateProp, updateNestedArrayItem, addArrayItem, removeArrayItem, selectThisBlock, isEditing } = usePuckPropUpdater(id);

        return (
          <div
            className="p-6 bg-transparent border-2 border-indigo-500/40 dark:border-indigo-500/30 rounded-3xl my-6 space-y-4 shadow-sm transition-all"
            id={id}
            onClick={() => selectThisBlock()}
          >
            <div className="flex justify-between items-center">
              <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-indigo-600 text-white">
                Один правильный ответ
              </span>
              <span className="text-xs font-bold text-indigo-600 dark:text-indigo-400">
                {points || 10} баллов
              </span>
            </div>

            <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
              <InlineText
                value={question}
                onChange={(val) => updateProp('question', val)}
                onFocusBlock={selectThisBlock}
                isEditing={isEditing}
                placeholder="Введите текст вопроса..."
                multiline
              />
            </h4>

            <div className="space-y-2">
              {options?.map((opt: any, idx: number) => {
                const isCorrect = opt.isCorrect === 'true' || opt.isCorrect === true;
                return (
                  <div
                    key={idx}
                    className={`group/opt p-3 bg-white dark:bg-slate-900 border rounded-2xl flex items-center justify-between text-xs font-medium transition-all ${
                      isCorrect
                        ? 'border-emerald-500 text-emerald-700 dark:text-emerald-300 ring-1 ring-emerald-500/20 shadow-sm'
                        : 'border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300'
                    }`}
                  >
                    <div className="flex-1 mr-2">
                      <InlineText
                        value={opt.text}
                        onChange={(val) => updateNestedArrayItem('options', idx, 'text', val)}
                        onFocusBlock={selectThisBlock}
                        isEditing={isEditing}
                        placeholder="Введите вариант ответа..."
                      />
                    </div>

                    <div className="flex items-center gap-1.5 flex-shrink-0" data-puck-overlay-portal="true">
                      {isEditing ? (
                        <button
                          type="button"
                          data-puck-overlay-portal="true"
                          onPointerDown={(e) => e.stopPropagation()}
                          onMouseDown={(e) => e.stopPropagation()}
                          onClick={(e) => {
                            e.stopPropagation();
                            selectThisBlock();
                            // В одиночном выборе ровно 1 вариант правильный (Radio)
                            const updatedOptions = (options || []).map((optItem: any, i: number) => ({
                              ...optItem,
                              isCorrect: i === idx ? 'true' : 'false',
                            }));
                            updateProp('options', updatedOptions);
                          }}
                          className={`edtech-inline-btn px-2.5 py-1 text-[10px] font-black rounded-lg transition-all cursor-pointer ${
                            isCorrect
                              ? 'bg-emerald-600 text-white shadow-sm'
                              : 'bg-slate-100 dark:bg-slate-800 text-slate-400 hover:text-emerald-500 hover:bg-emerald-50 dark:hover:bg-emerald-950/40'
                          }`}
                          title="Нажмите, чтобы сделать этот вариант единственным верным (Radio)"
                        >
                          {isCorrect ? '✓ Верный' : 'Сделать верным'}
                        </button>
                      ) : (
                        isCorrect && (
                          <span className="px-2 py-0.5 bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300 text-[10px] font-black rounded-lg">
                            Верный
                          </span>
                        )
                      )}

                      {isEditing && options.length > 2 && (
                        <button
                          type="button"
                          data-puck-overlay-portal="true"
                          onPointerDown={(e) => e.stopPropagation()}
                          onMouseDown={(e) => e.stopPropagation()}
                          onClick={(e) => {
                            e.stopPropagation();
                            selectThisBlock();
                            removeArrayItem('options', idx);
                          }}
                          className="edtech-inline-btn p-1.5 text-slate-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-colors rounded-lg cursor-pointer"
                          title="Удалить вариант"
                        >
                          <Trash2 size={13} />
                        </button>
                      )}
                    </div>
                  </div>
                );
              })}

              {isEditing && (
                <button
                  type="button"
                  data-puck-overlay-portal="true"
                  onPointerDown={(e) => e.stopPropagation()}
                  onMouseDown={(e) => e.stopPropagation()}
                  onClick={(e) => {
                    e.stopPropagation();
                    selectThisBlock();
                    addArrayItem('options', {
                      text: `Вариант ответа ${(options?.length || 0) + 1}`,
                      isCorrect: 'false',
                      explain: '',
                    });
                  }}
                  className="edtech-inline-btn edtech-inline-editable w-full py-2.5 border-2 border-dashed border-indigo-200 dark:border-indigo-800/60 rounded-xl text-xs font-bold text-indigo-600 dark:text-indigo-400 hover:bg-indigo-50 dark:hover:bg-indigo-950/40 flex items-center justify-center gap-1.5 transition-all mt-2 cursor-pointer shadow-sm"
                >
                  <Plus size={15} />
                  <span>Добавить вариант ответа</span>
                </button>
              )}
            </div>
          </div>
        );
      },
    },

    QuizMultiBlock: {
      label: 'Тест: выбор ответов',
      fields: {
        question: { type: 'textarea', label: 'Текст вопроса' },
        points: { type: 'number', label: 'Максимальный балл' },
        options: {
          type: 'array',
          label: 'Варианты ответов (несколько правильных)',
          getItemSummary: (item: any) => item.text || 'Новый вариант',
          arrayFields: {
            text: { type: 'text', label: 'Текст варианта' },
            isCorrect: {
              type: 'select',
              label: 'Правильный ответ?',
              options: [
                { label: 'Нет', value: 'false' },
                { label: 'Да', value: 'true' },
              ],
            },
          },
        },
      },
      defaultProps: {
        question: '2. Выберите ВСЕ преимущества использования Domain-Driven Design (DDD):',
        points: 15,
        options: [
          { text: 'Единый язык (Ubiquitous Language) между разработчиками и бизнесом', isCorrect: 'true' },
          { text: 'Четкие границы контекстов (Bounded Contexts)', isCorrect: 'true' },
          { text: 'Мгновенная компиляция любого кода без тестов', isCorrect: 'false' },
          { text: 'Богатая доменная модель без протекания деталей БД', isCorrect: 'true' },
        ],
      },
      render: (props: any) => {
        const { question, points, options, id } = props;
        const { updateProp, updateNestedArrayItem, addArrayItem, removeArrayItem, selectThisBlock, isEditing } = usePuckPropUpdater(id);

        return (
          <div
            className="p-6 bg-transparent border-2 border-purple-500/40 dark:border-purple-500/30 rounded-3xl my-6 space-y-4 shadow-sm transition-all"
            id={id}
            onClick={() => selectThisBlock()}
          >
            <div className="flex justify-between items-center">
              <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-purple-600 text-white">
                Несколько правильных ответов
              </span>
              <span className="text-xs font-bold text-purple-600 dark:text-purple-400">
                {points || 15} баллов
              </span>
            </div>

            <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
              <InlineText
                value={question}
                onChange={(val) => updateProp('question', val)}
                onFocusBlock={selectThisBlock}
                isEditing={isEditing}
                placeholder="Введите текст вопроса..."
                multiline
              />
            </h4>

            <div className="space-y-2">
              {options?.map((opt: any, idx: number) => {
                const isCorrect = opt.isCorrect === 'true' || opt.isCorrect === true;
                return (
                  <div
                    key={idx}
                    className={`group/opt p-3 bg-white dark:bg-slate-900 border rounded-2xl flex items-center justify-between text-xs font-medium transition-all ${
                      isCorrect
                        ? 'border-purple-500 text-purple-700 dark:text-purple-300 ring-1 ring-purple-500/20 shadow-sm'
                        : 'border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300'
                    }`}
                  >
                    <div className="flex-1 mr-2">
                      <InlineText
                        value={opt.text}
                        onChange={(val) => updateNestedArrayItem('options', idx, 'text', val)}
                        onFocusBlock={selectThisBlock}
                        isEditing={isEditing}
                        placeholder="Введите вариант ответа..."
                      />
                    </div>

                    <div className="flex items-center gap-1.5 flex-shrink-0" data-puck-overlay-portal="true">
                      {isEditing ? (
                        <button
                          type="button"
                          data-puck-overlay-portal="true"
                          onPointerDown={(e) => e.stopPropagation()}
                          onMouseDown={(e) => e.stopPropagation()}
                          onClick={(e) => {
                            e.stopPropagation();
                            selectThisBlock();
                            updateNestedArrayItem('options', idx, 'isCorrect', isCorrect ? 'false' : 'true');
                          }}
                          className={`edtech-inline-btn px-2.5 py-1 text-[10px] font-black rounded-lg transition-all cursor-pointer ${
                            isCorrect
                              ? 'bg-purple-600 text-white shadow-sm'
                              : 'bg-slate-100 dark:bg-slate-800 text-slate-400 hover:text-purple-500 hover:bg-purple-50 dark:hover:bg-purple-950/40'
                          }`}
                          title="Нажмите, чтобы переключить правильность ответа"
                        >
                          {isCorrect ? '✓ Верный' : 'Сделать верным'}
                        </button>
                      ) : (
                        isCorrect && (
                          <span className="px-2 py-0.5 bg-purple-100 dark:bg-purple-950 text-purple-700 dark:text-purple-300 text-[10px] font-black rounded-lg">
                            Верный
                          </span>
                        )
                      )}

                      {isEditing && options.length > 2 && (
                        <button
                          type="button"
                          data-puck-overlay-portal="true"
                          onPointerDown={(e) => e.stopPropagation()}
                          onMouseDown={(e) => e.stopPropagation()}
                          onClick={(e) => {
                            e.stopPropagation();
                            selectThisBlock();
                            removeArrayItem('options', idx);
                          }}
                          className="edtech-inline-btn p-1.5 text-slate-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-colors rounded-lg cursor-pointer"
                          title="Удалить вариант"
                        >
                          <Trash2 size={13} />
                        </button>
                      )}
                    </div>
                  </div>
                );
              })}

              {isEditing && (
                <button
                  type="button"
                  data-puck-overlay-portal="true"
                  onPointerDown={(e) => e.stopPropagation()}
                  onMouseDown={(e) => e.stopPropagation()}
                  onClick={(e) => {
                    e.stopPropagation();
                    selectThisBlock();
                    addArrayItem('options', {
                      text: `Вариант ответа ${(options?.length || 0) + 1}`,
                      isCorrect: 'false',
                    });
                  }}
                  className="edtech-inline-btn edtech-inline-editable w-full py-2.5 border-2 border-dashed border-purple-200 dark:border-purple-800/60 rounded-xl text-xs font-bold text-purple-600 dark:text-purple-400 hover:bg-purple-50 dark:hover:bg-purple-950/40 flex items-center justify-center gap-1.5 transition-all mt-2 cursor-pointer shadow-sm"
                >
                  <Plus size={15} />
                  <span>Добавить вариант ответа</span>
                </button>
              )}
            </div>
          </div>
        );
      },
    },

    QuizMatchBlock: {
      label: 'Сопоставление пар',
      fields: {
        question: { type: 'textarea', label: 'Текст вопроса' },
        points: { type: 'number', label: 'Максимальный балл' },
        pairs: {
          type: 'array',
          label: 'Пары для сопоставления',
          getItemSummary: (item: any) =>
            item.left ? `${item.left} ➔ ${item.right || '?'}` : 'Новая пара',
          arrayFields: {
            left: { type: 'text', label: 'Левая колонка (Термин)' },
            right: { type: 'text', label: 'Правая колонка (Определение/Пара)' },
          },
        },
      },
      defaultProps: {
        question: '3. Сопоставьте архитектурные понятия и их определения:',
        points: 20,
        pairs: [
          { left: 'AST Tree', right: 'Иерархическая структура блоков страницы' },
          { left: 'Block Props', right: 'Параметры и настройки конкретного компонента' },
          { left: 'Puck Canvas', right: 'Интерактивная область предпросмотра редактора' },
        ],
      },
      render: (props: any) => {
        const { question, points, pairs, id } = props;
        const { updateProp, updateNestedArrayItem, addArrayItem, removeArrayItem, selectThisBlock, isEditing } = usePuckPropUpdater(id);

        return (
          <div
            className="p-6 bg-transparent border-2 border-cyan-500/40 dark:border-cyan-500/30 rounded-3xl my-6 space-y-4 shadow-sm"
            id={id}
            onClick={() => selectThisBlock()}
          >
            <div className="flex justify-between items-center">
              <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-cyan-600 text-white">
                Сопоставление 2 колонок
              </span>
              <span className="text-xs font-bold text-cyan-600 dark:text-cyan-400">
                {points || 20} баллов
              </span>
            </div>

            <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
              <InlineText
                value={question}
                onChange={(val) => updateProp('question', val)}
                onFocusBlock={selectThisBlock}
                isEditing={isEditing}
                placeholder="Введите текст вопроса..."
                multiline
              />
            </h4>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <div className="text-[11px] font-extrabold text-cyan-700 dark:text-cyan-400 uppercase tracking-wider">Левая колонка:</div>
                {pairs?.map((p: any, i: number) => (
                  <div key={i} className="p-3 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl text-xs font-bold text-slate-800 dark:text-slate-200">
                    <InlineText
                      value={p.left}
                      onChange={(val) => updateNestedArrayItem('pairs', i, 'left', val)}
                      onFocusBlock={selectThisBlock}
                      isEditing={isEditing}
                      placeholder="Термин..."
                    />
                  </div>
                ))}
              </div>

              <div className="space-y-2">
                <div className="text-[11px] font-extrabold text-cyan-700 dark:text-cyan-400 uppercase tracking-wider">Правая колонка:</div>
                {pairs?.map((p: any, i: number) => (
                  <div key={i} className="group/pair p-3 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl text-xs font-medium text-slate-600 dark:text-slate-400 flex items-center justify-between">
                    <div className="flex-1 mr-2">
                      <InlineText
                        value={p.right}
                        onChange={(val) => updateNestedArrayItem('pairs', i, 'right', val)}
                        onFocusBlock={selectThisBlock}
                        isEditing={isEditing}
                        placeholder="Определение..."
                      />
                    </div>
                    {isEditing && pairs.length > 2 && (
                      <button
                        type="button"
                        data-puck-overlay-portal="true"
                        onPointerDown={(e) => e.stopPropagation()}
                        onMouseDown={(e) => e.stopPropagation()}
                        onClick={(e) => {
                          e.stopPropagation();
                          selectThisBlock();
                          removeArrayItem('pairs', i);
                        }}
                        className="edtech-inline-btn p-1 text-slate-400 hover:text-rose-500 rounded-lg cursor-pointer"
                        title="Удалить пару"
                      >
                        <Trash2 size={13} />
                      </button>
                    )}
                  </div>
                ))}
              </div>
            </div>

            {isEditing && (
              <button
                type="button"
                data-puck-overlay-portal="true"
                onPointerDown={(e) => e.stopPropagation()}
                onMouseDown={(e) => e.stopPropagation()}
                onClick={(e) => {
                  e.stopPropagation();
                  selectThisBlock();
                  addArrayItem('pairs', { left: 'Новый термин', right: 'Новое определение' });
                }}
                className="edtech-inline-btn edtech-inline-editable w-full py-2.5 border-2 border-dashed border-cyan-200 dark:border-cyan-800/60 rounded-xl text-xs font-bold text-cyan-600 dark:text-cyan-400 hover:bg-cyan-50/50 dark:hover:bg-cyan-950/30 flex items-center justify-center gap-1.5 transition-colors mt-2 cursor-pointer shadow-sm"
              >
                <Plus size={14} />
                <span>Добавить пару для сопоставления</span>
              </button>
            )}
          </div>
        );
      },
    },

    QuizDropdownBlankBlock: {
      label: 'Выпадающие списки (пропуски)',
      fields: {
        question: { type: 'textarea', label: 'Текст вопроса / задание' },
        points: { type: 'number', label: 'Максимальный балл' },
        templateText: {
          type: 'custom',
          label: 'Текст с вариантами {Правильный; Неверный 1, Неверный 2}',
          render: ({ value, onChange, readOnly }: { value: string; onChange: (v: string) => void; readOnly?: boolean }) => {
            const { blanks } = parseSmartDropdownTemplate(value || '');

            return (
              <div className="space-y-3">
                <textarea
                  value={value || ''}
                  disabled={readOnly}
                  onChange={(e) => onChange(e.target.value)}
                  rows={4}
                  placeholder="Введите текст с вариантами: {OpenAPI; Buf, Contract}"
                  className="w-full p-3 text-xs bg-slate-900 border border-slate-700 rounded-xl text-slate-100 focus:outline-none focus:border-indigo-500 font-sans leading-relaxed resize-y"
                />

                {blanks.length > 0 ? (
                  <div className="p-3 bg-slate-900/90 border border-emerald-500/30 rounded-xl space-y-2.5">
                    <div className="flex items-center justify-between text-[11px] font-bold text-emerald-400">
                      <span>✨ Распознано полей (аргументов): {blanks.length}</span>
                    </div>

                    <div className="space-y-2">
                      {blanks.map((b, idx) => {
                        const distractors = b.options.slice(1);
                        return (
                          <div
                            key={idx}
                            className="p-2.5 bg-slate-950/80 border border-slate-800 rounded-lg space-y-1 text-xs"
                          >
                            <div className="flex items-center justify-between">
                              <span className="font-extrabold text-indigo-400 text-[11px]">
                                {b.key}
                              </span>
                              <span className="text-[10px] bg-emerald-950 text-emerald-300 border border-emerald-800 px-1.5 py-0.2 rounded font-bold">
                                ✓ {b.correctAnswer}
                              </span>
                            </div>

                            {distractors.length > 0 ? (
                              <div className="text-[10px] text-slate-400 flex flex-wrap gap-1 items-center">
                                <span className="text-slate-500">Варианты ошибок:</span>
                                {distractors.map((d, dIdx) => (
                                   <span
                                    key={dIdx}
                                    className="px-1.5 py-0.2 bg-slate-900 border border-slate-800 rounded text-slate-300"
                                  >
                                    ✕ {d}
                                  </span>
                                ))}
                              </div>
                            ) : (
                              <div className="text-[10px] text-amber-400/80 italic">
                                Укажите другие варианты через запятую: {`{${b.correctAnswer}; вариант2, вариант3}`}
                              </div>
                            )}
                          </div>
                        );
                      })}
                    </div>
                  </div>
                ) : (
                  <div className="p-2.5 bg-slate-900/60 border border-slate-800 rounded-xl text-[11px] text-slate-400 space-y-1">
                    <p className="font-bold text-slate-300">💡 Формат записи:</p>
                    <p>
                      Пишите варианты в фигурных скобках:{' '}
                      <code className="text-emerald-400 font-mono text-[10px]">
                        {'{Правильный; Ошибка1, Ошибка2}'}
                      </code>
                    </p>
                    <p className="text-[10px] text-slate-500">
                      Первый вариант всегда считается правильным ответом.
                    </p>
                  </div>
                )}
              </div>
            );
          },
        },
      },
      defaultProps: {
        question: '4. Заполните пропуски, выбрав правильные варианты из списков:',
        points: 10,
        templateText:
          'Спецификация для REST API называется {OpenAPI; Buf, Contract}, а для RPC систем часто используется {gRPC; Buf, OpenAPI}.',
      },
      render: (props: any) => {
        const { question, points, templateText, blanks, textBefore, textAfter, options, id } = props;
        const { updateProp, selectThisBlock, isEditing } = usePuckPropUpdater(id);

        const fullText =
          templateText ||
          `${textBefore || ''} {${options?.map((o: any) => o.text).join('; ') || 'вариант 1; вариант 2'}} ${textAfter || ''}`;

        const { tokens } = parseSmartDropdownTemplate(fullText, blanks);

        return (
          <div
            className="p-6 bg-transparent border-2 border-emerald-500/40 dark:border-emerald-500/30 rounded-3xl my-6 space-y-4 shadow-sm"
            id={id}
            onClick={() => selectThisBlock()}
          >
            <div className="flex justify-between items-center">
              <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-emerald-600 text-white">
                Пропуски со списком (Dropdown)
              </span>
              <span className="text-xs font-bold text-emerald-600 dark:text-emerald-400">
                {points || 10} баллов
              </span>
            </div>

            <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
              <InlineText
                value={question}
                onChange={(val) => updateProp('question', val)}
                onFocusBlock={selectThisBlock}
                isEditing={isEditing}
                placeholder="Введите текст вопроса..."
                multiline
              />
            </h4>

            <div className="p-5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl text-xs md:text-sm text-slate-800 dark:text-slate-200 leading-loose flex flex-wrap items-center gap-1.5">
              {tokens.map((token, tIdx) => {
                if (token.type === 'text') {
                  return <span key={tIdx}>{token.value}</span>;
                }
                const b = token.blank!;
                const distractors = b.options.slice(1);

                return (
                  <span
                    key={tIdx}
                    className="inline-flex items-center gap-1 px-3 py-1 bg-emerald-100 dark:bg-emerald-950/80 text-emerald-800 dark:text-emerald-200 rounded-xl border border-emerald-300 dark:border-emerald-700 font-bold mx-1 text-xs shadow-2xs"
                    title={`Поле {${b.key}}: правильный ответ "${b.correctAnswer}" | Варианты: ${b.options.join(', ')}`}
                  >
                    <span className="text-[10px] text-emerald-600 font-black">[{b.key}]</span>
                    <span className="flex items-center gap-1">
                      <span className="bg-emerald-600 text-white px-1.5 py-0.2 rounded text-[10px]">
                        ✓ {b.correctAnswer}
                      </span>
                      {distractors.length > 0 && (
                        <span className="text-[10px] text-slate-400">
                          (варианты: {distractors.join(', ')})
                        </span>
                      )}
                    </span>
                    <span className="text-[10px] text-emerald-600">▼</span>
                  </span>
                );
              })}
            </div>
          </div>
        );
      },
    },

    QuizInputBlankBlock: {
      label: 'Ввод слова (пропуск)',
      fields: {
        question: { type: 'textarea', label: 'Текст вопроса / задание' },
        points: { type: 'number', label: 'Максимальный балл' },
        templateText: {
          type: 'custom',
          label: 'Текст с пропущенными словами {ПравильныйОтвет}',
          render: ({ value, onChange, readOnly }: { value: string; onChange: (v: string) => void; readOnly?: boolean }) => {
            const { blanks } = parseSmartDropdownTemplate(value || '');

            return (
              <div className="space-y-3">
                <textarea
                  value={value || ''}
                  disabled={readOnly}
                  onChange={(e) => onChange(e.target.value)}
                  rows={4}
                  placeholder="Введите текст с ответом: {OpenAPI}"
                  className="w-full p-3 text-xs bg-slate-900 border border-slate-700 rounded-xl text-slate-100 focus:outline-none focus:border-indigo-500 font-sans leading-relaxed resize-y"
                />

                {blanks.length > 0 ? (
                  <div className="p-3 bg-slate-900/90 border border-blue-500/30 rounded-xl space-y-2.5">
                    <div className="flex items-center justify-between text-[11px] font-bold text-blue-400">
                      <span>✨ Распознано полей ввода (аргументов): {blanks.length}</span>
                    </div>

                    <div className="space-y-2">
                      {blanks.map((b, idx) => (
                        <div
                          key={idx}
                          className="p-2.5 bg-slate-950/80 border border-slate-800 rounded-lg flex items-center justify-between text-xs"
                        >
                          <span className="font-extrabold text-blue-400 text-[11px]">
                            {b.key}
                          </span>
                          <span className="text-[10px] bg-blue-950 text-blue-300 border border-blue-800 px-2 py-0.5 rounded font-bold">
                            Правильный ответ: «{b.correctAnswer}»
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>
                ) : (
                  <div className="p-2.5 bg-slate-900/60 border border-slate-800 rounded-xl text-[11px] text-slate-400 space-y-1">
                    <p className="font-bold text-slate-300">💡 Формат записи:</p>
                    <p>
                      Укажите правильный ответ в скобках:{' '}
                      <code className="text-blue-400 font-mono text-[10px]">{'{OpenAPI}'}</code>
                    </p>
                  </div>
                )}
              </div>
            );
          },
        },
      },
      defaultProps: {
        question: '5. Введите точные термины для заполнения пропусков:',
        points: 10,
        templateText:
          'Спецификация для описания контрактов REST API называется {OpenAPI}, а для RPC систем часто используется {gRPC}.',
      },
      render: (props: any) => {
        const { question, points, templateText, blanks, prefixText, suffixText, correctAnswer, id } = props;
        const { updateProp, selectThisBlock, isEditing } = usePuckPropUpdater(id);

        const fullText =
          templateText ||
          `${prefixText || ''} {${correctAnswer || 'ответ'}} ${suffixText || ''}`;

        const { tokens } = parseSmartDropdownTemplate(fullText, blanks);

        return (
          <div
            className="p-6 bg-transparent border-2 border-blue-500/40 dark:border-blue-500/30 rounded-3xl my-6 space-y-4 shadow-sm"
            id={id}
            onClick={() => selectThisBlock()}
          >
            <div className="flex justify-between items-center">
              <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-blue-600 text-white">
                Пропуски с ручным вводом
              </span>
              <span className="text-xs font-bold text-blue-600 dark:text-blue-400">
                {points || 10} баллов
              </span>
            </div>

            <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
              <InlineText
                value={question}
                onChange={(val) => updateProp('question', val)}
                onFocusBlock={selectThisBlock}
                isEditing={isEditing}
                placeholder="Введите текст вопроса..."
                multiline
              />
            </h4>

            <div className="p-5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl text-xs md:text-sm text-slate-800 dark:text-slate-200 leading-loose flex flex-wrap items-center gap-1.5">
              {tokens.map((token, tIdx) => {
                if (token.type === 'text') {
                  return <span key={tIdx}>{token.value}</span>;
                }
                const b = token.blank!;

                return (
                  <span
                    key={tIdx}
                    className="inline-flex items-center gap-1.5 px-3 py-1 bg-blue-50 dark:bg-blue-950/80 text-blue-700 dark:text-blue-300 font-bold rounded-xl border border-dashed border-blue-400 dark:border-blue-700 shadow-2xs mx-1"
                    title={`Поле {${b.key}}: правильный ответ "${b.correctAnswer}"`}
                  >
                    <span className="text-[10px] text-blue-500 font-black">[{b.key}]</span>
                    <span className="italic font-normal">«{b.correctAnswer}»</span>
                  </span>
                );
              })}
            </div>
          </div>
        );
      },
    },

    QuizSequenceBlock: {
      label: 'Последовательность (порядок)',
      fields: {
        question: { type: 'textarea', label: 'Текст вопроса' },
        points: { type: 'number', label: 'Максимальный балл' },
        items: {
          type: 'array',
          label: 'Элементы последовательности (в правильном порядке)',
          getItemSummary: (item: any) => item.text || 'Новый этап',
          arrayFields: {
            text: { type: 'text', label: 'Этап / Шаг' },
          },
        },
      },
      defaultProps: {
        question: '6. Расставьте этапы обработки запроса в правильном хронологическом порядке:',
        points: 15,
        items: [
          { text: '1. Прием HTTP запроса маршрутизатором (Router)' },
          { text: '2. Проверка JWT-токена в Middleware' },
          { text: '3. Валидация входного DTO в хендлере' },
          { text: '4. Выполнение бизнес-логики в сервисе' },
          { text: '5. Запись в БД через репозиторий и ответ клиенту' },
        ],
      },
      render: (props: any) => {
        const { question, points, items, id } = props;
        const { updateProp, updateNestedArrayItem, addArrayItem, removeArrayItem, selectThisBlock, isEditing } = usePuckPropUpdater(id);

        return (
          <div
            className="p-6 bg-transparent border-2 border-amber-500/40 dark:border-amber-500/30 rounded-3xl my-6 space-y-4 shadow-sm"
            id={id}
            onClick={() => selectThisBlock()}
          >
            <div className="flex justify-between items-center">
              <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-amber-600 text-white">
                Хронологическая последовательность
              </span>
              <span className="text-xs font-bold text-amber-600 dark:text-amber-400">
                {points || 15} баллов
              </span>
            </div>

            <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
              <InlineText
                value={question}
                onChange={(val) => updateProp('question', val)}
                onFocusBlock={selectThisBlock}
                isEditing={isEditing}
                placeholder="Введите текст вопроса..."
                multiline
              />
            </h4>

            <div className="space-y-2">
              {items?.map((item: any, idx: number) => (
                <div
                  key={idx}
                  className="group/seq p-3 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl text-xs font-bold text-slate-800 dark:text-slate-200 flex items-center justify-between gap-3"
                >
                  <div className="flex items-center gap-2.5 flex-1 mr-2">
                    <span className="w-5 h-5 rounded-full bg-amber-100 dark:bg-amber-950 text-amber-700 dark:text-amber-300 text-[10px] font-black flex items-center justify-center flex-shrink-0">
                      {idx + 1}
                    </span>
                    <div className="flex-1">
                      <InlineText
                        value={item.text}
                        onChange={(val) => updateNestedArrayItem('items', idx, 'text', val)}
                        onFocusBlock={selectThisBlock}
                        isEditing={isEditing}
                        placeholder="Введите название шага..."
                      />
                    </div>
                  </div>

                  {isEditing && items.length > 2 && (
                    <button
                      type="button"
                      data-puck-overlay-portal="true"
                      onPointerDown={(e) => e.stopPropagation()}
                      onMouseDown={(e) => e.stopPropagation()}
                      onClick={(e) => {
                        e.stopPropagation();
                        selectThisBlock();
                        removeArrayItem('items', idx);
                      }}
                      className="edtech-inline-btn p-1 text-slate-400 hover:text-rose-500 rounded-lg cursor-pointer"
                      title="Удалить шаг"
                    >
                      <Trash2 size={13} />
                    </button>
                  )}
                </div>
              ))}

              {isEditing && (
                <button
                  type="button"
                  data-puck-overlay-portal="true"
                  onPointerDown={(e) => e.stopPropagation()}
                  onMouseDown={(e) => e.stopPropagation()}
                  onClick={(e) => {
                    e.stopPropagation();
                    selectThisBlock();
                    addArrayItem('items', { text: `Шаг ${(items?.length || 0) + 1}` });
                  }}
                  className="edtech-inline-btn edtech-inline-editable w-full py-2 border-2 border-dashed border-amber-200 dark:border-amber-800/60 rounded-xl text-xs font-bold text-amber-600 dark:text-amber-400 hover:bg-amber-50/50 dark:hover:bg-amber-950/30 flex items-center justify-center gap-1.5 transition-colors mt-2 cursor-pointer shadow-sm"
                >
                  <Plus size={14} />
                  <span>Добавить этап</span>
                </button>
              )}
            </div>
          </div>
        );
      },
    },

    QuizEssayBlock: {
      label: 'Развернутый ответ (эссе)',
      fields: {
        question: { type: 'textarea', label: 'Тема эссе / Развернутый вопрос' },
        points: { type: 'number', label: 'Максимальный балл' },
        rubric: { type: 'textarea', label: 'Критерии оценивания для учителя' },
        sampleAnswer: { type: 'textarea', label: 'Эталонный ответ / Ключевые тезисы' },
      },
      defaultProps: {
        question: '7. Опишите разницу между синхронным REST и асинхронным Event-Driven взаимодействием:',
        points: 25,
        rubric: 'Оценивается понимание блокировок, брокеров сообщений (Kafka/RabbitMQ) и eventual consistency.',
        sampleAnswer: 'REST — request-reply, блокирует клиент. Event-Driven — слабая связность через события.',
      },
      render: (props: any) => {
        const { question, points, rubric, id } = props;
        const { updateProp, selectThisBlock, isEditing } = usePuckPropUpdater(id);

        return (
          <div
            className="p-6 bg-transparent border-2 border-rose-500/40 dark:border-rose-500/30 rounded-3xl my-6 space-y-4 shadow-sm"
            id={id}
            onClick={() => selectThisBlock()}
          >
            <div className="flex justify-between items-center">
              <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-rose-600 text-white">
                Развернутый ответ (Эссе)
              </span>
              <span className="text-xs font-bold text-rose-600 dark:text-rose-400">
                {points || 25} баллов
              </span>
            </div>

            <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
              <InlineText
                value={question}
                onChange={(val) => updateProp('question', val)}
                onFocusBlock={selectThisBlock}
                isEditing={isEditing}
                placeholder="Введите тему эссе или развернутый вопрос..."
                multiline
              />
            </h4>

            {rubric && (
              <p className="text-xs text-slate-500 dark:text-slate-400 italic">
                Критерии: {rubric}
              </p>
            )}

            <div className="p-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl text-xs text-slate-400 italic">
              [Студент напишет здесь свой развернутый ответ в текстовом поле]
            </div>
          </div>
        );
      },
    },

    FileUploadBlock: {
      label: 'Загрузка файла / ДЗ',
      fields: {
        title: { type: 'text', label: 'Заголовок практического задания' },
        instructions: { type: 'textarea', label: 'Инструкции по выполнению и загрузке' },
        points: { type: 'number', label: 'Максимальный балл' },
        allowedTypes: { type: 'text', label: 'Разрешенные форматы (напр. .zip, .pdf, .go)' },
        maxSizeMB: { type: 'number', label: 'Максимальный размер (МБ)' },
      },
      defaultProps: {
        title: '8. Практическая работа: Проектирование микросервиса',
        instructions: 'Прикрепите архив с исходным кодом решения или ссылку на GitHub репозиторий.',
        points: 50,
        allowedTypes: '.zip, .tar.gz, .go, .pdf',
        maxSizeMB: 25,
      },
      render: (props: any) => {
        const { title, instructions, points, allowedTypes, maxSizeMB, id } = props;
        const { updateProp, selectThisBlock, isEditing } = usePuckPropUpdater(id);

        return (
          <div
            className="p-6 bg-transparent border-2 border-slate-400/40 dark:border-slate-600/40 rounded-3xl my-6 space-y-4 shadow-sm"
            id={id}
            onClick={() => selectThisBlock()}
          >
            <div className="flex justify-between items-center">
              <span className="px-2.5 py-1 rounded-xl text-[10px] font-extrabold uppercase tracking-wider bg-slate-700 text-white">
                Загрузка файла / ДЗ
              </span>
              <span className="text-xs font-bold text-slate-700 dark:text-slate-300">
                {points || 50} баллов
              </span>
            </div>

            <h4 className="text-sm md:text-base font-bold text-slate-900 dark:text-white">
              <InlineText
                value={title}
                onChange={(val) => updateProp('title', val)}
                onFocusBlock={selectThisBlock}
                isEditing={isEditing}
                placeholder="Заголовок практической работы..."
              />
            </h4>

            <div className="text-xs text-slate-600 dark:text-slate-400">
              <InlineText
                value={instructions}
                onChange={(val) => updateProp('instructions', val)}
                onFocusBlock={selectThisBlock}
                isEditing={isEditing}
                placeholder="Инструкции по выполнению..."
                multiline
              />
            </div>

            <div className="p-8 border-2 border-dashed border-slate-300 dark:border-slate-600 rounded-2xl text-center space-y-2 bg-white/60 dark:bg-slate-900/60">
              <span className="text-xs font-bold text-indigo-600 dark:text-indigo-400 block">
                Перетащите файл или нажмите для выбора
              </span>
              <span className="text-[11px] text-slate-400 block">
                Форматы: {allowedTypes || 'Любые'} (до {maxSizeMB || 25} МБ)
              </span>
            </div>
          </div>
        );
      },
    },
  },
};
