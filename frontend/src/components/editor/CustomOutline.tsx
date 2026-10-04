'use client';

import React from 'react';
import { usePuck } from '@puckeditor/core';
import {
  Heading,
  AlignLeft,
  FileText,
  Video,
  CheckCircle2,
  ListChecks,
  Shuffle,
  ChevronDownSquare,
  KeyRound,
  ListOrdered,
  PenTool,
  UploadCloud,
  ArrowUp,
  ArrowDown,
  Copy,
  Trash2,
  Layers,
  Presentation
} from 'lucide-react';

function stripHtml(html: string): string {
  if (!html) return '';
  return html.replace(/<[^>]*>?/gm, ' ').replace(/\s+/g, ' ').trim();
}

interface BlockOutlineMeta {
  title: string;
  previewText: string;
  badge?: string;
  icon: React.ReactNode;
}

function getBlockOutlineMeta(item: any): BlockOutlineMeta {
  const type = item.type;
  const props = item.props || {};

  switch (type) {
    case 'HeaderBlock':
      return {
        title: 'Заголовок',
        previewText: props.title || props.subtitle || 'Без текста',
        icon: <Heading size={14} className="text-blue-400" />,
      };

    case 'TextBlock':
      return {
        title: 'Текст (Canvas)',
        previewText: props.content || stripHtml(props.contentHtml || '') || 'Пустой текстовый блок',
        icon: <AlignLeft size={14} className="text-emerald-400" />,
      };

    case 'RichTextBlock':
      return {
        title: 'Лекция / Word',
        previewText: props.title || stripHtml(props.contentHtml || '') || 'Статья лекции',
        badge: 'Word',
        icon: <FileText size={14} className="text-indigo-400" />,
      };

    case 'VideoBlock':
      return {
        title: 'Видео-урок',
        previewText: props.caption || props.url || 'Видео',
        icon: <Video size={14} className="text-red-400" />,
      };

    case 'PresentationBlock':
      return {
        title: 'Презентация / Слайды',
        previewText: props.title || props.embedUrl || props.pdfUrl || 'Слайды презентации',
        badge: props.mode === 'pdf' ? 'PDF' : 'Слайды',
        icon: <Presentation size={14} className="text-amber-400" />,
      };

    case 'QuizSingleBlock':
      return {
        title: 'Тест: один ответ',
        previewText: props.question || 'Вопрос без текста',
        badge: `${props.points || 10} б.`,
        icon: <CheckCircle2 size={14} className="text-indigo-400" />,
      };

    case 'QuizMultiBlock':
      return {
        title: 'Тест: выбор ответов',
        previewText: props.question || 'Вопрос без текста',
        badge: `${props.points || 15} б.`,
        icon: <ListChecks size={14} className="text-purple-400" />,
      };

    case 'QuizMatchBlock':
      return {
        title: 'Сопоставление пар',
        previewText: props.question || 'Задание на сопоставление',
        badge: `${props.points || 20} б.`,
        icon: <Shuffle size={14} className="text-cyan-400" />,
      };

    case 'QuizDropdownBlankBlock':
      return {
        title: 'Выпадающие списки',
        previewText: props.question || props.templateText || 'Задание с пропусками',
        badge: `${props.points || 15} б.`,
        icon: <ChevronDownSquare size={14} className="text-amber-400" />,
      };

    case 'QuizInputBlankBlock':
      return {
        title: 'Ввод слова (пропуск)',
        previewText: props.question || props.templateText || 'Ввод пропущенного слова',
        badge: `${props.points || 15} б.`,
        icon: <KeyRound size={14} className="text-teal-400" />,
      };

    case 'QuizSequenceBlock':
      return {
        title: 'Последовательность',
        previewText: props.question || 'Хронологический порядок',
        badge: `${props.points || 15} б.`,
        icon: <ListOrdered size={14} className="text-amber-400" />,
      };

    case 'QuizEssayBlock':
      return {
        title: 'Развернутый ответ (эссе)',
        previewText: props.question || 'Тема эссе',
        badge: `${props.points || 25} б.`,
        icon: <PenTool size={14} className="text-rose-400" />,
      };

    case 'FileUploadBlock':
      return {
        title: 'Загрузка файла / ДЗ',
        previewText: props.title || props.instructions || 'Практическое задание',
        badge: `${props.points || 50} б.`,
        icon: <UploadCloud size={14} className="text-slate-400" />,
      };

    default:
      return {
        title: type || 'Блок',
        previewText: props.title || props.question || props.content || '',
        icon: <Layers size={14} className="text-slate-400" />,
      };
  }
}

export function CustomOutline({ children }: { children?: React.ReactNode } = {}) {
  const { appState, dispatch, selectedItem } = usePuck();
  const items = appState.data.content || [];

  return (
    <div className="flex flex-col h-full bg-slate-900 border-r border-slate-800 text-white select-none">
      {/* Outline Header */}
      <div className="p-3.5 border-b border-slate-800 flex items-center justify-between bg-slate-950/60">
        <div className="flex items-center gap-2">
          <Layers size={15} className="text-indigo-400" />
          <span className="text-xs font-black uppercase tracking-wider text-slate-200">
            Слои и блоки
          </span>
        </div>
        <span className="px-2 py-0.5 rounded-full text-[10px] font-extrabold bg-slate-800 text-slate-300 border border-slate-700">
          {items.length} {items.length === 1 ? 'слой' : items.length < 5 ? 'слоя' : 'слоев'}
        </span>
      </div>

      {/* Layers List */}
      <div className="flex-1 overflow-y-auto p-2.5 space-y-2">
        {items.length === 0 ? (
          <div className="p-6 text-center text-slate-500 text-xs font-medium">
            Урок пока пуст. Перетащите блоки из вкладки «Блоки» на холст.
          </div>
        ) : (
          items.map((item, index) => {
            const isSelected = selectedItem?.props?.id === item.props?.id;
            const meta = getBlockOutlineMeta(item);

            return (
              <div
                key={item.props?.id || index}
                onClick={() => {
                  dispatch({
                    type: 'setUi',
                    ui: { itemSelector: { index, zone: 'default' } },
                  });
                  // Smoothly scroll to block in preview canvas
                  const el = document.getElementById(item.props?.id) || document.querySelector(`[data-puck-component-id="${item.props?.id}"]`);
                  el?.scrollIntoView({ behavior: 'smooth', block: 'center' });
                }}
                className={`group p-3 rounded-2xl border transition-all cursor-pointer flex flex-col gap-1.5 ${
                  isSelected
                    ? 'bg-indigo-950/80 border-indigo-500 shadow-md ring-1 ring-indigo-500/40'
                    : 'bg-slate-950/40 border-slate-800 hover:border-slate-700 hover:bg-slate-800/50'
                }`}
              >
                {/* Header row: Index, Icon, Type Name, Badge */}
                <div className="flex items-center justify-between gap-2">
                  <div className="flex items-center gap-1.5 min-w-0">
                    <span className="text-[10px] font-black text-slate-500">#{index + 1}</span>
                    <span className="p-1 rounded-md bg-slate-800/80">{meta.icon}</span>
                    <span className="text-xs font-bold text-slate-200 truncate">{meta.title}</span>
                  </div>

                  {meta.badge && (
                    <span className="px-1.5 py-0.5 rounded text-[9px] font-black bg-indigo-950 text-indigo-300 border border-indigo-800/60 flex-shrink-0">
                      {meta.badge}
                    </span>
                  )}
                </div>

                {/* Text Snippet / Content Preview */}
                <div className="text-[11px] text-slate-400 line-clamp-2 leading-relaxed pl-5 font-medium">
                  {meta.previewText ? (
                    <span>«{meta.previewText}»</span>
                  ) : (
                    <span className="opacity-40 italic">Пустой блок</span>
                  )}
                </div>

                {/* Quick actions on hover */}
                <div className="flex items-center justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity pt-1 border-t border-slate-800/50">
                  <button
                    type="button"
                    title="Поднять выше"
                    disabled={index === 0}
                    onClick={(e) => {
                      e.stopPropagation();
                      dispatch({
                        type: 'reorder',
                        sourceIndex: index,
                        destinationIndex: index - 1,
                        destinationZone: 'default',
                      });
                    }}
                    className="p-1 rounded-lg hover:bg-slate-700 text-slate-400 hover:text-white disabled:opacity-20 transition-colors"
                  >
                    <ArrowUp size={12} />
                  </button>

                  <button
                    type="button"
                    title="Опустить ниже"
                    disabled={index === items.length - 1}
                    onClick={(e) => {
                      e.stopPropagation();
                      dispatch({
                        type: 'reorder',
                        sourceIndex: index,
                        destinationIndex: index + 1,
                        destinationZone: 'default',
                      });
                    }}
                    className="p-1 rounded-lg hover:bg-slate-700 text-slate-400 hover:text-white disabled:opacity-20 transition-colors"
                  >
                    <ArrowDown size={12} />
                  </button>

                  <button
                    type="button"
                    title="Дублировать блок"
                    onClick={(e) => {
                      e.stopPropagation();
                      dispatch({
                        type: 'duplicate',
                        sourceIndex: index,
                        sourceZone: 'default',
                      });
                    }}
                    className="p-1 rounded-lg hover:bg-slate-700 text-slate-400 hover:text-white transition-colors"
                  >
                    <Copy size={12} />
                  </button>

                  <button
                    type="button"
                    title="Удалить блок"
                    onClick={(e) => {
                      e.stopPropagation();
                      dispatch({
                        type: 'remove',
                        index,
                        zone: 'default',
                      });
                    }}
                    className="p-1 rounded-lg hover:bg-rose-950/60 text-slate-400 hover:text-rose-400 transition-colors"
                  >
                    <Trash2 size={12} />
                  </button>
                </div>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}

export default CustomOutline;
