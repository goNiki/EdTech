'use client';

import React, { useRef, useState, useCallback, useEffect } from 'react';
import { usePuck, registerOverlayPortal } from '@puckeditor/core';
import {
  Bold,
  Italic,
  Underline,
  Strikethrough,
  Code,
  List,
  ListOrdered,
  AlignLeft,
  AlignCenter,
  AlignRight,
  AlignJustify,
  Heading1,
  Heading2,
  Heading3,
  Quote,
  Image as ImageIcon,
  Link as LinkIcon,
  Minus,
  Sparkles,
  Info,
  Highlighter,
  RotateCcw,
  RotateCw,
  Undo,
  Redo,
  Plus,
  Trash2,
  X,
  Check,
  GripVertical,
  Indent,
  Outdent,
  RemoveFormatting,
  Type,
  Palette,
  FileUp,
  FileText,
  Loader2,
  Upload,
  AlertTriangle,
  Target
} from 'lucide-react';
import { convertDocumentToHtml } from '@/lib/document-importer';
import { api } from '@/lib/api';

/**
 * Hook to safely interact with Puck's dispatch for updating block props.
 */
export function usePuckPropUpdater(blockId?: string) {
  let puckApi: any = null;
  try {
    puckApi = usePuck();
  } catch (e) {
    puckApi = null;
  }

  const selectThisBlock = useCallback(() => {
    if (!puckApi || !blockId) return;

    let selector = puckApi.getSelectorForId?.(blockId);

    if (!selector) {
      const items = puckApi.appState?.data?.content || [];
      const index = items.findIndex((it: any) => (it.props?.id || it.id) === blockId);
      if (index !== -1) {
        selector = { zone: 'root:default-zone', index };
      } else if (puckApi.appState?.data?.zones) {
        for (const [zoneCompound, zoneItems] of Object.entries(puckApi.appState.data.zones)) {
          const zIdx = (zoneItems as any[]).findIndex((it: any) => (it.props?.id || it.id) === blockId);
          if (zIdx !== -1) {
            selector = { zone: zoneCompound, index: zIdx };
            break;
          }
        }
      }
    }

    if (selector) {
      const current = puckApi.appState?.ui?.itemSelector;
      if (current?.zone !== selector.zone || current?.index !== selector.index) {
        puckApi.dispatch({
          type: 'setUi',
          ui: {
            itemSelector: selector,
          },
        });
      }
    }
  }, [puckApi, blockId]);

  const updateProp = useCallback(
    (field: string, value: any) => {
      if (!puckApi || !blockId) return;
      puckApi.dispatch({
        type: 'setData',
        data: (prev: any) => {
          const updateList = (items: any[]): any[] => {
            return (items || []).map((item) => {
              const currentId = item.props?.id || item.id;
              if (currentId === blockId) {
                return {
                  ...item,
                  props: {
                    ...item.props,
                    [field]: value,
                  },
                };
              }
              if (item.props?.content && Array.isArray(item.props.content)) {
                return {
                  ...item,
                  props: {
                    ...item.props,
                    content: updateList(item.props.content),
                  },
                };
              }
              return item;
            });
          };

          const newContent = updateList(prev?.content || []);
          const newZones: Record<string, any[]> = {};
          if (prev?.zones) {
            for (const [k, v] of Object.entries(prev.zones)) {
              newZones[k] = updateList(v as any[]);
            }
          }

          return {
            ...prev,
            content: newContent,
            ...(prev?.zones ? { zones: newZones } : {}),
          };
        },
      });
    },
    [puckApi, blockId]
  );

  const updateProps = useCallback(
    (fields: Record<string, any>) => {
      if (!puckApi || !blockId) return;
      puckApi.dispatch({
        type: 'setData',
        data: (prev: any) => {
          const updateList = (items: any[]): any[] => {
            return (items || []).map((item) => {
              const currentId = item.props?.id || item.id;
              if (currentId === blockId) {
                return {
                  ...item,
                  props: {
                    ...item.props,
                    ...fields,
                  },
                };
              }
              if (item.props?.content && Array.isArray(item.props.content)) {
                return {
                  ...item,
                  props: {
                    ...item.props,
                    content: updateList(item.props.content),
                  },
                };
              }
              return item;
            });
          };

          const newContent = updateList(prev?.content || []);
          const newZones: Record<string, any[]> = {};
          if (prev?.zones) {
            for (const [k, v] of Object.entries(prev.zones)) {
              newZones[k] = updateList(v as any[]);
            }
          }

          return {
            ...prev,
            content: newContent,
            ...(prev?.zones ? { zones: newZones } : {}),
          };
        },
      });
    },
    [puckApi, blockId]
  );

  const updateNestedArrayItem = useCallback(
    (arrayField: string, index: number, subField: string, value: any) => {
      if (!puckApi || !blockId) return;
      puckApi.dispatch({
        type: 'setData',
        data: (prev: any) => {
          const updateList = (items: any[]): any[] => {
            return (items || []).map((item) => {
              const currentId = item.props?.id || item.id;
              if (currentId === blockId) {
                const arr = [...(item.props?.[arrayField] || [])];
                if (arr[index]) {
                  arr[index] = { ...arr[index], [subField]: value };
                }
                return {
                  ...item,
                  props: {
                    ...item.props,
                    [arrayField]: arr,
                  },
                };
              }
              if (item.props?.content && Array.isArray(item.props.content)) {
                return {
                  ...item,
                  props: {
                    ...item.props,
                    content: updateList(item.props.content),
                  },
                };
              }
              return item;
            });
          };

          const newContent = updateList(prev?.content || []);
          const newZones: Record<string, any[]> = {};
          if (prev?.zones) {
            for (const [k, v] of Object.entries(prev.zones)) {
              newZones[k] = updateList(v as any[]);
            }
          }

          return {
            ...prev,
            content: newContent,
            ...(prev?.zones ? { zones: newZones } : {}),
          };
        },
      });
    },
    [puckApi, blockId]
  );

  const addArrayItem = useCallback(
    (arrayField: string, newItem: any) => {
      if (!puckApi || !blockId) return;
      selectThisBlock();
      puckApi.dispatch({
        type: 'setData',
        data: (prev: any) => {
          const updateList = (items: any[]): any[] => {
            return (items || []).map((item) => {
              const currentId = item.props?.id || item.id;
              if (currentId === blockId) {
                const currentArr = Array.isArray(item.props?.[arrayField]) ? item.props[arrayField] : [];
                return {
                  ...item,
                  props: {
                    ...item.props,
                    [arrayField]: [...currentArr, newItem],
                  },
                };
              }
              if (item.props?.content && Array.isArray(item.props.content)) {
                return {
                  ...item,
                  props: {
                    ...item.props,
                    content: updateList(item.props.content),
                  },
                };
              }
              return item;
            });
          };

          const newContent = updateList(prev?.content || []);
          const newZones: Record<string, any[]> = {};
          if (prev?.zones) {
            for (const [k, v] of Object.entries(prev.zones)) {
              newZones[k] = updateList(v as any[]);
            }
          }

          return {
            ...prev,
            content: newContent,
            ...(prev?.zones ? { zones: newZones } : {}),
          };
        },
      });

      setTimeout(() => {
        selectThisBlock();
      }, 50);
    },
    [puckApi, blockId, selectThisBlock]
  );

  const removeArrayItem = useCallback(
    (arrayField: string, index: number) => {
      if (!puckApi || !blockId) return;
      selectThisBlock();
      puckApi.dispatch({
        type: 'setData',
        data: (prev: any) => {
          const updateList = (items: any[]): any[] => {
            return (items || []).map((item) => {
              const currentId = item.props?.id || item.id;
              if (currentId === blockId) {
                const arr = [...(item.props?.[arrayField] || [])];
                arr.splice(index, 1);
                return {
                  ...item,
                  props: {
                    ...item.props,
                    [arrayField]: arr,
                  },
                };
              }
              if (item.props?.content && Array.isArray(item.props.content)) {
                return {
                  ...item,
                  props: {
                    ...item.props,
                    content: updateList(item.props.content),
                  },
                };
              }
              return item;
            });
          };

          const newContent = updateList(prev?.content || []);
          const newZones: Record<string, any[]> = {};
          if (prev?.zones) {
            for (const [k, v] of Object.entries(prev.zones)) {
              newZones[k] = updateList(v as any[]);
            }
          }

          return {
            ...prev,
            content: newContent,
            ...(prev?.zones ? { zones: newZones } : {}),
          };
        },
      });

      setTimeout(() => {
        selectThisBlock();
      }, 50);
    },
    [puckApi, blockId, selectThisBlock]
  );

  return { updateProp, updateProps, updateNestedArrayItem, addArrayItem, removeArrayItem, selectThisBlock, isEditing: !!puckApi };
}

/**
 * Universal InlineText Component:
 * When isEditing=true, allows clicking and editing text directly on the canvas with smooth blur sync.
 * Also notifies Puck to select the parent block so the sidebar opens.
 */
interface InlineTextProps {
  value: string;
  onChange?: (val: string) => void;
  className?: string;
  placeholder?: string;
  multiline?: boolean;
  as?: any;
  isEditing?: boolean;
  onFocusBlock?: () => void;
}

export function InlineText({
  value,
  onChange,
  className = '',
  placeholder = 'Нажмите, чтобы ввести текст...',
  multiline = false,
  as: Component = 'span',
  isEditing = false,
  onFocusBlock,
}: InlineTextProps) {
  const [isFocused, setIsFocused] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (ref.current && !isFocused) {
      if (ref.current.innerText !== (value || '')) {
        ref.current.innerText = value || '';
      }
    }
  }, [value, isFocused]);

  // Native keyboard event isolation to stop Puck from blocking Backspace / Delete
  useEffect(() => {
    const el = ref.current;
    if (!el || !isEditing) return;

    const stopPropagation = (e: KeyboardEvent) => {
      e.stopPropagation();
    };

    el.addEventListener('keydown', stopPropagation, false);
    el.addEventListener('keyup', stopPropagation, false);
    el.addEventListener('keypress', stopPropagation, false);

    return () => {
      el.removeEventListener('keydown', stopPropagation, false);
      el.removeEventListener('keyup', stopPropagation, false);
      el.removeEventListener('keypress', stopPropagation, false);
    };
  }, [isEditing]);

  if (!isEditing || !onChange) {
    return <Component className={className}>{value || placeholder}</Component>;
  }

  return (
    <Component
      ref={ref}
      contentEditable
      suppressContentEditableWarning
      data-puck-overlay-portal="true"
      onPointerDownCapture={(e: React.PointerEvent) => {
        e.stopPropagation();
      }}
      onMouseDownCapture={(e: React.MouseEvent) => {
        e.stopPropagation();
      }}
      onKeyDownCapture={(e: React.KeyboardEvent) => {
        e.stopPropagation();
      }}
      onKeyUpCapture={(e: React.KeyboardEvent) => {
        e.stopPropagation();
      }}
      onFocus={() => {
        setIsFocused(true);
        onFocusBlock?.();
      }}
      onClick={(e: React.MouseEvent) => {
        e.stopPropagation();
        onFocusBlock?.();
      }}
      onBlur={(e: React.FocusEvent<HTMLDivElement>) => {
        setIsFocused(false);
        const text = e.currentTarget.innerText.trim();
        if (text !== value) {
          onChange(text);
        }
      }}
      onKeyDown={(e: React.KeyboardEvent<HTMLDivElement>) => {
        e.stopPropagation();
        if (!multiline && e.key === 'Enter') {
          e.preventDefault();
          e.currentTarget.blur();
        }
      }}
      className={`edtech-inline-editable ${className} cursor-text transition-all ${
        isFocused
          ? 'outline-2 outline-indigo-500 rounded px-1 min-w-[20px] inline-block'
          : 'hover:outline-1 hover:outline-dashed hover:outline-indigo-400/50 rounded px-1'
      }`}
      title="✏️ Нажмите, чтобы редактировать текст прямо в блоке"
    >
      {value || (isFocused ? '' : <span className="opacity-40 italic">{placeholder}</span>)}
    </Component>
  );
}

export const EDITOR_SCOPED_STYLES = `
  .edtech-inline-editable h1, [contenteditable] h1 {
    font-size: 2rem !important;
    line-height: 2.35rem !important;
    font-weight: 800 !important;
    margin-top: 1.25rem !important;
    margin-bottom: 0.5rem !important;
    color: inherit !important;
  }
  .edtech-inline-editable h2, [contenteditable] h2 {
    font-size: 1.5rem !important;
    line-height: 1.85rem !important;
    font-weight: 700 !important;
    margin-top: 1rem !important;
    margin-bottom: 0.4rem !important;
    color: inherit !important;
  }
  .edtech-inline-editable h3, [contenteditable] h3 {
    font-size: 1.25rem !important;
    line-height: 1.65rem !important;
    font-weight: 600 !important;
    margin-top: 0.85rem !important;
    margin-bottom: 0.35rem !important;
    color: inherit !important;
  }
  .edtech-inline-editable p, [contenteditable] p {
    font-size: 1rem !important;
    line-height: 1.65 !important;
    margin-top: 0.5rem !important;
    margin-bottom: 0.5rem !important;
  }
  .edtech-inline-editable ul, [contenteditable] ul {
    list-style-type: disc !important;
    padding-left: 1.75rem !important;
    margin-top: 0.5rem !important;
    margin-bottom: 0.5rem !important;
  }
  .edtech-inline-editable ol, [contenteditable] ol {
    list-style-type: decimal !important;
    padding-left: 1.75rem !important;
    margin-top: 0.5rem !important;
    margin-bottom: 0.5rem !important;
  }
  .edtech-inline-editable li, [contenteditable] li {
    display: list-item !important;
    margin-top: 0.25rem !important;
    margin-bottom: 0.25rem !important;
  }
  .edtech-inline-editable blockquote, [contenteditable] blockquote {
    border-left: 4px solid #6366f1 !important;
    padding: 0.6rem 1rem !important;
    margin: 1rem 0 !important;
    background-color: rgba(99, 102, 241, 0.08) !important;
    border-radius: 0 0.75rem 0.75rem 0 !important;
    font-style: italic !important;
  }
  .edtech-inline-editable u, [contenteditable] u {
    text-decoration: underline !important;
  }
`;

export async function uploadOrConvertImage(file: File): Promise<string> {
  try {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('category', 'lesson_media');
    const res = await api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    const fileUrl = res.data?.data?.file_url || res.data?.file_url;
    if (fileUrl) return fileUrl;
  } catch (err) {
    console.warn('Backend image upload failed, falling back to data URL:', err);
  }

  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

export type CalloutType = 'tip' | 'warning' | 'note' | 'target';

export function getCalloutSnippet(type: CalloutType): string {
  const configs = {
    tip: {
      icon: '💡',
      title: 'Совет',
      border: 'border-l-4 border-indigo-500',
      bg: 'bg-indigo-50/90 dark:bg-indigo-950/40',
      textColor: 'text-indigo-950 dark:text-indigo-100',
      desc: 'Полезная рекомендация или подсказка к материалу...',
    },
    warning: {
      icon: '⚠️',
      title: 'Важно',
      border: 'border-l-4 border-amber-500',
      bg: 'bg-amber-50/90 dark:bg-amber-950/40',
      textColor: 'text-amber-950 dark:text-amber-100',
      desc: 'Обратите особое внимание на этот момент...',
    },
    note: {
      icon: '📝',
      title: 'Заметка',
      border: 'border-l-4 border-emerald-500',
      bg: 'bg-emerald-50/90 dark:bg-emerald-950/40',
      textColor: 'text-emerald-950 dark:text-emerald-100',
      desc: 'Дополнительные сведения или справочная информация...',
    },
    target: {
      icon: '🎯',
      title: 'Цель / Результат',
      border: 'border-l-4 border-purple-500',
      bg: 'bg-purple-50/90 dark:bg-purple-950/40',
      textColor: 'text-purple-950 dark:text-purple-100',
      desc: 'Ожидаемый результат или ключевой вывод темы...',
    },
  };
  const c = configs[type];
  return `
    <div class="my-4 p-4 rounded-2xl ${c.bg} ${c.border} flex items-start gap-3 shadow-xs">
      <span class="text-xl select-none leading-none pt-0.5">${c.icon}</span>
      <div class="text-sm leading-relaxed ${c.textColor} flex-1">
        <strong>${c.title}:</strong> ${c.desc}
      </div>
    </div>
    <p><br/></p>
  `;
}

/**
 * Unified RichTextCanvasEditor:
 * In-place Word/Google Docs style text editing on the canvas with:
 * - Paragraphs (P), Headings (H1, H2, H3), Blockquote
 * - Alignment (Left, Center, Right, Justify)
 * - Font styling (Bold, Italic, Underline, Strikethrough, Code, Highlighter)
 * - Lists (Bullet, Numbered) and Indentation (Indent, Outdent)
 * - Image Upload & Direct Clipboard Paste (Ctrl+V)
 * - Document Import (.docx / .md)
 * - Colored Callout Boxes (💡 Совет, ⚠️ Важно, 📝 Заметка, 🎯 Цель)
 * - Horizontal Divider, Clear formatting
 */
interface RichTextCanvasEditorProps {
  htmlContent: string;
  onChange: (html: string) => void;
  isEditing?: boolean;
  onFocusBlock?: () => void;
  defaultAlign?: string;
  placeholder?: string;
  className?: string;
  dragRef?: any;
}

export function RichTextCanvasEditor({
  htmlContent,
  onChange,
  isEditing = false,
  onFocusBlock,
  defaultAlign = 'left',
  placeholder = 'Начните вводить текст лекции...',
  className = '',
  dragRef,
}: RichTextCanvasEditorProps) {
  const editorRef = useRef<HTMLDivElement>(null);
  const isInitializedRef = useRef(false);
  const initialHtmlRef = useRef<string>(htmlContent || '');
  const lastHtmlRef = useRef<string>(htmlContent || '');
  const savedSelectionRef = useRef<Range | null>(null);

  // Image insertion state (Upload vs URL)
  const [showImageModal, setShowImageModal] = useState(false);
  const [imageTab, setImageTab] = useState<'upload' | 'url'>('upload');
  const [imageUrl, setImageUrl] = useState('');
  const [imageCaption, setImageCaption] = useState('');
  const [imageWidth, setImageWidth] = useState<'100%' | '75%' | '50%'>('100%');
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [isUploadingImage, setIsUploadingImage] = useState(false);
  const imageFileInputRef = useRef<HTMLInputElement>(null);

  // Link & Color & Callouts state
  const [showLinkModal, setShowLinkModal] = useState(false);
  const [linkUrl, setLinkUrl] = useState('');
  const [showColorPicker, setShowColorPicker] = useState(false);
  const [showCalloutMenu, setShowCalloutMenu] = useState(false);

  // Document import state (.docx / .md)
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isImporting, setIsImporting] = useState(false);
  const [importStatus, setImportStatus] = useState<string>('');
  const [isDragOver, setIsDragOver] = useState(false);
  const [importToast, setImportToast] = useState<string | null>(null);

  const getEditorDoc = useCallback(() => {
    return editorRef.current?.ownerDocument || (typeof document !== 'undefined' ? document : null);
  }, []);

  const getEditorWin = useCallback(() => {
    const doc = getEditorDoc();
    return doc?.defaultView || (typeof window !== 'undefined' ? window : null);
  }, [getEditorDoc]);

  // Save current selection range inside the editor
  const saveSelection = useCallback(() => {
    const win = getEditorWin();
    if (!win) return;
    const sel = win.getSelection();
    if (sel && sel.rangeCount > 0) {
      const range = sel.getRangeAt(0);
      if (
        editorRef.current &&
        (editorRef.current === range.commonAncestorContainer || editorRef.current.contains(range.commonAncestorContainer))
      ) {
        savedSelectionRef.current = range.cloneRange();
      }
    }
  }, [getEditorWin]);

  // Restore saved selection before applying commands
  const restoreSelection = useCallback(() => {
    const win = getEditorWin();
    if (!win || !editorRef.current) return;
    const sel = win.getSelection();
    if (!sel || !savedSelectionRef.current) return;
    if (sel.rangeCount > 0) {
      const curRange = sel.getRangeAt(0);
      if (
        curRange.startContainer === savedSelectionRef.current.startContainer &&
        curRange.startOffset === savedSelectionRef.current.startOffset &&
        curRange.endContainer === savedSelectionRef.current.endContainer &&
        curRange.endOffset === savedSelectionRef.current.endOffset
      ) {
        return;
      }
    }
    try {
      sel.removeAllRanges();
      sel.addRange(savedSelectionRef.current);
    } catch (e) {}
  }, [getEditorWin]);

  // Track selection changes across both editor ownerDocument and top document
  useEffect(() => {
    const doc = getEditorDoc();
    const handleSelectionChange = () => {
      saveSelection();
    };
    doc?.addEventListener('selectionchange', handleSelectionChange);
    if (typeof document !== 'undefined' && document !== doc) {
      document.addEventListener('selectionchange', handleSelectionChange);
    }
    return () => {
      doc?.removeEventListener('selectionchange', handleSelectionChange);
      if (typeof document !== 'undefined' && document !== doc) {
        document.removeEventListener('selectionchange', handleSelectionChange);
      }
    };
  }, [getEditorDoc, saveSelection]);

  // Initialize once and set clean default paragraph separator
  useEffect(() => {
    if (editorRef.current && !isInitializedRef.current) {
      editorRef.current.innerHTML = htmlContent || '';
      lastHtmlRef.current = htmlContent || '';
      isInitializedRef.current = true;
    }
    const doc = getEditorDoc();
    try {
      doc?.execCommand('defaultParagraphSeparator', false, 'p');
      doc?.execCommand('styleWithCSS', false, 'true');
    } catch (e) {}
  }, [getEditorDoc, htmlContent]);

  // Sync external changes ONLY when prop actually changed from outside and user is not currently typing inside
  useEffect(() => {
    if (editorRef.current && isInitializedRef.current) {
      const doc = getEditorDoc();
      const isFocused =
        Boolean(doc?.activeElement && (doc.activeElement === editorRef.current || editorRef.current.contains(doc.activeElement)));
      if (!isFocused && htmlContent !== undefined && htmlContent !== lastHtmlRef.current) {
        lastHtmlRef.current = htmlContent;
        editorRef.current.innerHTML = htmlContent || '';
      }
    }
  }, [getEditorDoc, htmlContent]);

  const triggerChange = useCallback(() => {
    if (!editorRef.current) return;
    const newHtml = editorRef.current.innerHTML;
    lastHtmlRef.current = newHtml;
    onChange(newHtml);
  }, [onChange]);

  const exec = (command: string, value: string | undefined = undefined) => {
    const doc = getEditorDoc();
    if (!doc || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    try {
      doc.execCommand('styleWithCSS', false, 'false');
    } catch (e) {}
    doc.execCommand(command, false, value);
    saveSelection();
    triggerChange();
  };

  const insertCustomHtml = (htmlSnippet: string) => {
    const doc = getEditorDoc();
    if (!doc || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    doc.execCommand('insertHTML', false, htmlSnippet);
    saveSelection();
    triggerChange();
  };

  // Document Import (.docx / .md)
  const handleFileImport = async (file: File) => {
    const ext = file.name.toLowerCase();
    if (!ext.endsWith('.docx') && !ext.endsWith('.md') && !ext.endsWith('.markdown')) {
      alert('Пожалуйста, выберите файл документа Microsoft Word (.docx) или Markdown (.md)');
      return;
    }

    setIsImporting(true);
    setImportStatus('Чтение файла...');
    try {
      const res = await convertDocumentToHtml(file, (msg) => setImportStatus(msg));

      if (editorRef.current) {
        const currentHtml = editorRef.current.innerHTML.trim();
        const isDefault =
          !currentHtml ||
          currentHtml === '<p>Пустой текст</p>' ||
          currentHtml === '<p><br></p>' ||
          currentHtml === '<p>Начните вводить текст лекции...</p>';
        const newHtml = isDefault ? res.html : `${currentHtml}<hr class="my-6 border-slate-200 dark:border-slate-800" />${res.html}`;
        editorRef.current.innerHTML = newHtml;
        lastHtmlRef.current = newHtml;
        onChange(newHtml);
      }

      setImportToast(`Документ успешно импортирован (${res.wordCount} слов, ${res.tablesCount} таблиц)`);
      setTimeout(() => setImportToast(null), 4000);
    } catch (err: any) {
      console.error('Import failed', err);
      alert(err.message || 'Ошибка импорта документа');
    } finally {
      setIsImporting(false);
      setImportStatus('');
    }
  };

  // Keyboard navigation, Puck isolation, and Clipboard Paste (Ctrl+V) handler
  useEffect(() => {
    const el = editorRef.current;
    if (!el) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      // Isolate keydown from Puck hotkeys
      e.stopPropagation();

      // Tab key: Indentation in lists or paragraph indent in text
      if (e.key === 'Tab') {
        e.preventDefault();
        const doc = getEditorDoc();
        const win = getEditorWin();
        const sel = win?.getSelection();

        let isInList = false;
        if (sel && sel.rangeCount > 0) {
          let node: Node | null = sel.getRangeAt(0).startContainer;
          while (node && node !== editorRef.current) {
            if (node.nodeName === 'LI' || node.nodeName === 'UL' || node.nodeName === 'OL') {
              isInList = true;
              break;
            }
            node = node.parentNode;
          }
        }

        if (isInList) {
          if (e.shiftKey) {
            exec('outdent');
          } else {
            exec('indent');
          }
        } else {
          doc?.execCommand('insertHTML', false, '&emsp;&emsp;');
          triggerChange();
        }
        return;
      }

      // Hotkeys: Ctrl+B, Ctrl+I, Ctrl+U, Ctrl+Z, Ctrl+Y
      if (e.ctrlKey || e.metaKey) {
        const k = e.key.toLowerCase();
        if (k === 'b' || k === 'и') {
          e.preventDefault();
          exec('bold');
          return;
        }
        if (k === 'i' || k === 'ш') {
          e.preventDefault();
          exec('italic');
          return;
        }
        if (k === 'u' || k === 'г') {
          e.preventDefault();
          exec('underline');
          return;
        }
        if (k === 'z' || k === 'я') {
          if (e.shiftKey) {
            e.preventDefault();
            exec('redo');
          } else {
            e.preventDefault();
            exec('undo');
          }
          return;
        }
        if (k === 'y' || k === 'н') {
          e.preventDefault();
          exec('redo');
          return;
        }
      }
    };

    // Intercept image pastes from clipboard (Ctrl+V screenshot / copied image)
    const handlePaste = async (e: ClipboardEvent) => {
      const items = e.clipboardData?.items;
      if (!items) return;

      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        if (item.type.indexOf('image') !== -1) {
          e.preventDefault();
          e.stopPropagation();
          const file = item.getAsFile();
          if (!file) continue;

          try {
            const imgUrl = await uploadOrConvertImage(file);
            const snippet = `
              <figure class="my-4 mx-auto w-full text-center">
                <img src="${imgUrl}" alt="Изображение из буфера" class="rounded-2xl shadow-md border border-slate-200 dark:border-slate-800 mx-auto object-cover max-h-[460px]" />
              </figure>
              <p><br/></p>
            `;
            insertCustomHtml(snippet);
          } catch (err) {
            console.error('Failed to paste image:', err);
          }
          return;
        }
      }
    };

    const stopPropagation = (e: KeyboardEvent) => {
      e.stopPropagation();
    };

    el.addEventListener('keydown', handleKeyDown, false);
    el.addEventListener('keyup', stopPropagation, false);
    el.addEventListener('keypress', stopPropagation, false);
    el.addEventListener('paste', handlePaste, false);

    return () => {
      el.removeEventListener('keydown', handleKeyDown, false);
      el.removeEventListener('keyup', stopPropagation, false);
      el.removeEventListener('keypress', stopPropagation, false);
      el.removeEventListener('paste', handlePaste, false);
    };
  }, [getEditorDoc, getEditorWin, triggerChange]);

  const formatHeading = (tag: 'p' | 'h1' | 'h2' | 'h3') => {
    const doc = getEditorDoc();
    if (!doc || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    try {
      doc.execCommand('styleWithCSS', false, 'false');
    } catch (e) {}
    let ok = doc.execCommand('formatBlock', false, `<${tag}>`);
    if (!ok) {
      doc.execCommand('formatBlock', false, tag);
    }
    saveSelection();
    triggerChange();
  };

  const applyInlineStyle = (styleObj: Record<string, string>) => {
    const doc = getEditorDoc();
    const win = getEditorWin();
    if (!doc || !win || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    const sel = win.getSelection();
    if (!sel || sel.rangeCount === 0) return;
    const range = sel.getRangeAt(0);
    if (!editorRef.current.contains(range.commonAncestorContainer) && editorRef.current !== range.commonAncestorContainer) return;

    if (range.collapsed) {
      let el: HTMLElement | null =
        (range.startContainer.nodeType === Node.ELEMENT_NODE
          ? range.startContainer
          : range.startContainer.parentElement) as HTMLElement;
      while (el && el !== editorRef.current && el.parentElement !== editorRef.current) {
        el = el.parentElement;
      }
      if (el && el !== editorRef.current) {
        Object.assign(el.style, styleObj);
      } else {
        Object.assign(editorRef.current.style, styleObj);
      }
      triggerChange();
      return;
    }

    const contents = range.extractContents();
    const span = doc.createElement('span');
    Object.assign(span.style, styleObj);
    span.appendChild(contents);
    range.insertNode(span);

    const newRange = doc.createRange();
    newRange.selectNodeContents(span);
    sel.removeAllRanges();
    sel.addRange(newRange);
    savedSelectionRef.current = newRange.cloneRange();
    triggerChange();
  };

  const formatFont = (fontFamily: string) => {
    applyInlineStyle({ fontFamily });
  };

  const formatSize = (sizePx: string) => {
    applyInlineStyle({ fontSize: sizePx });
  };

  const formatColor = (color: string) => {
    applyInlineStyle({ color });
  };

  const formatAlign = (alignCmd: 'justifyLeft' | 'justifyCenter' | 'justifyRight' | 'justifyFull') => {
    const doc = getEditorDoc();
    const win = getEditorWin();
    if (!doc || !win || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    try {
      doc.execCommand('styleWithCSS', false, 'false');
    } catch (e) {}
    doc.execCommand(alignCmd, false);
    const sel = win.getSelection();
    if (!sel || sel.rangeCount === 0 || sel.isCollapsed) {
      const map: Record<string, string> = {
        justifyLeft: 'left',
        justifyCenter: 'center',
        justifyRight: 'right',
        justifyFull: 'justify',
      };
      if (map[alignCmd]) {
        editorRef.current.style.textAlign = map[alignCmd];
      }
    }
    saveSelection();
    triggerChange();
  };

  const formatHighlight = (color: string = '#fef08a') => {
    applyInlineStyle({ backgroundColor: color });
  };

  const formatCode = () => {
    const doc = getEditorDoc();
    const win = getEditorWin();
    if (!doc || !win || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    const sel = win.getSelection();
    if (sel && !sel.isCollapsed && sel.rangeCount > 0) {
      const range = sel.getRangeAt(0);
      if (editorRef.current.contains(range.commonAncestorContainer)) {
        const contents = range.extractContents();
        const code = doc.createElement('code');
        code.className = 'px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-indigo-600 dark:text-indigo-400 font-mono text-xs';
        code.appendChild(contents);
        range.insertNode(code);

        const newRange = doc.createRange();
        newRange.selectNodeContents(code);
        sel.removeAllRanges();
        sel.addRange(newRange);
        savedSelectionRef.current = newRange.cloneRange();
        triggerChange();
        return;
      }
    }
    formatHeading('p');
    exec('formatBlock', '<pre>');
  };

  const clearFormatting = () => {
    const doc = getEditorDoc();
    const win = getEditorWin();
    if (!doc || !win || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    try {
      doc.execCommand('styleWithCSS', false, 'false');
    } catch (e) {}
    doc.execCommand('removeFormat', false);
    doc.execCommand('unlink', false);
    const sel = win.getSelection();
    if (sel && !sel.isCollapsed && sel.rangeCount > 0) {
      const range = sel.getRangeAt(0);
      if (editorRef.current.contains(range.commonAncestorContainer)) {
        let el = range.commonAncestorContainer as HTMLElement;
        if (el.nodeType !== Node.ELEMENT_NODE) el = el.parentElement as HTMLElement;
        while (el && el !== editorRef.current) {
          if (el.tagName === 'SPAN' || el.tagName === 'FONT' || el.tagName === 'MARK') {
            el.removeAttribute('style');
            el.removeAttribute('face');
            el.removeAttribute('size');
            el.removeAttribute('color');
          }
          el = el.parentElement as HTMLElement;
        }
      }
    }
    saveSelection();
    triggerChange();
  };

  const handleInsertImage = async () => {
    let finalUrl = imageUrl.trim();

    if (imageTab === 'upload' && imageFile) {
      setIsUploadingImage(true);
      try {
        finalUrl = await uploadOrConvertImage(imageFile);
      } catch (err) {
        alert('Не удалось загрузить изображение');
        setIsUploadingImage(false);
        return;
      }
      setIsUploadingImage(false);
    }

    if (!finalUrl) return;

    const widthClass =
      imageWidth === '50%' ? 'max-w-[50%]' : imageWidth === '75%' ? 'max-w-[75%]' : 'w-full';

    const snippet = `
      <figure class="my-4 mx-auto ${widthClass} text-center">
        <img src="${finalUrl}" alt="${imageCaption || 'Иллюстрация'}" class="rounded-2xl shadow-md border border-slate-200 dark:border-slate-800 mx-auto object-cover max-h-[460px]" />
        ${imageCaption ? `<figcaption class="mt-1.5 text-xs text-slate-500 dark:text-slate-400 font-medium">${imageCaption}</figcaption>` : ''}
      </figure>
      <p><br/></p>
    `;
    insertCustomHtml(snippet);
    setImageUrl('');
    setImageCaption('');
    setImageFile(null);
    setImagePreview(null);
    setShowImageModal(false);
  };

  const handleInsertCallout = (type: CalloutType) => {
    insertCustomHtml(getCalloutSnippet(type));
    setShowCalloutMenu(false);
  };

  const handleInsertLink = () => {
    if (!linkUrl.trim()) return;
    let url = linkUrl.trim();
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      url = 'https://' + url;
    }
    exec('createLink', url);
    setLinkUrl('');
    setShowLinkModal(false);
  };

  if (!isEditing) {
    return (
      <div
        className={`prose dark:prose-invert max-w-none text-slate-800 dark:text-slate-200 leading-relaxed ${className}`}
        style={{ textAlign: defaultAlign as any }}
        dangerouslySetInnerHTML={{ __html: htmlContent || '<p>Пустой текст</p>' }}
      />
    );
  }

  const preventBtnFocus = {
    onPointerDown: (e: React.PointerEvent) => {
      e.preventDefault();
      e.stopPropagation();
    },
    onMouseDown: (e: React.MouseEvent) => {
      e.preventDefault();
      e.stopPropagation();
    },
  };

  return (
    <div
      className="my-3 rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm overflow-hidden transition-all relative"
      onClick={() => onFocusBlock?.()}
    >
      {/* Embedded Scoped Headings & Lists Styles */}
      <style dangerouslySetInnerHTML={{ __html: EDITOR_SCOPED_STYLES }} />

      {/* Sleek Canvas Formatting Ribbon (Unified with Word editor) */}
      <div
        data-puck-overlay-portal="true"
        onPointerDown={(e) => {
          const target = e.target as HTMLElement;
          if (target && (target.tagName === 'INPUT' || target.tagName === 'SELECT' || target.tagName === 'TEXTAREA')) {
            return;
          }
          e.preventDefault();
          e.stopPropagation();
        }}
        onMouseDown={(e) => {
          const target = e.target as HTMLElement;
          if (target && (target.tagName === 'INPUT' || target.tagName === 'SELECT' || target.tagName === 'TEXTAREA')) {
            return;
          }
          e.preventDefault();
          e.stopPropagation();
        }}
        className="bg-slate-100/90 dark:bg-slate-900/90 px-3 py-1.5 border-b border-slate-200 dark:border-slate-800 flex flex-wrap items-center gap-1.5 text-slate-700 dark:text-slate-300 select-none text-xs"
      >
        {/* Block Badge */}
        <div
          ref={dragRef}
          data-puck-drag-handle={dragRef ? 'true' : undefined}
          className="flex items-center gap-1.5 px-2 py-1 rounded-lg bg-indigo-600 text-white text-[11px] font-bold select-none shadow-2xs"
          title="Текстовый блок (Canvas)"
        >
          <Type size={13} />
          <span>Текст</span>
        </div>

        {/* Undo / Redo */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Отменить действие (Ctrl+Z)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('undo'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300"
          >
            <Undo size={14} />
          </button>
          <button
            type="button"
            title="Повторить действие (Ctrl+Y)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('redo'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300"
          >
            <Redo size={14} />
          </button>
        </div>

        {/* Headings Hierarchy: P (16px), H1 (32px), H2 (24px), H3 (20px) */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Обычный текст (16px)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatHeading('p'); }}
            className="edtech-inline-btn px-2 py-1 text-xs font-bold rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            P
          </button>
          <button
            type="button"
            title="Заголовок 1 (Крупный, 32px)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatHeading('h1'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400 font-extrabold"
          >
            <Heading1 size={14} />
          </button>
          <button
            type="button"
            title="Заголовок 2 (Средний, 24px)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatHeading('h2'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400 font-bold"
          >
            <Heading2 size={14} />
          </button>
          <button
            type="button"
            title="Подзаголовок 3 (Малый, 20px)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatHeading('h3'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400 font-semibold"
          >
            <Heading3 size={14} />
          </button>
        </div>

        {/* Font Family (Шрифт) */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Шрифт: Стандартный без засечек"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatFont('sans-serif'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[11px] font-medium rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            Sans
          </button>
          <button
            type="button"
            title="Шрифт: С засечками (Serif / Книга)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatFont('Georgia, serif'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[11px] font-serif rounded hover:bg-slate-100 dark:hover:bg-slate-800 italic"
          >
            Serif
          </button>
          <button
            type="button"
            title="Шрифт: Моноширинный (Код)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatFont('monospace'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[11px] font-mono rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            Mono
          </button>
        </div>

        {/* Font Size (Размер текста) */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Размер шрифта: Мелкий"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatSize('13px'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[10px] font-bold rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            A-
          </button>
          <button
            type="button"
            title="Размер шрифта: Стандартный"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatSize('16px'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-xs font-bold rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            A
          </button>
          <button
            type="button"
            title="Размер шрифта: Крупный"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatSize('22px'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-xs font-black rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400"
          >
            A+
          </button>
        </div>

        {/* Text Color Palette */}
        <div className="relative">
          <button
            type="button"
            title="Цвет текста"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); setShowColorPicker(!showColorPicker); }}
            className="edtech-inline-btn p-1.5 bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400"
          >
            <Palette size={14} />
          </button>
          {showColorPicker && (
            <div
              className="absolute top-full left-0 mt-1 p-2 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl shadow-xl flex items-center gap-1.5 z-50"
              onPointerDownCapture={(e) => e.stopPropagation()}
              onMouseDownCapture={(e) => e.stopPropagation()}
            >
              {[
                { label: 'Стандартный', color: 'inherit', bg: 'bg-slate-800 dark:bg-slate-200' },
                { label: 'Индиго', color: '#4f46e5', bg: 'bg-indigo-600' },
                { label: 'Изумруд', color: '#059669', bg: 'bg-emerald-600' },
                { label: 'Янтарь', color: '#d97706', bg: 'bg-amber-600' },
                { label: 'Красный', color: '#e11d48', bg: 'bg-rose-600' },
                { label: 'Голубой', color: '#0284c7', bg: 'bg-sky-600' },
                { label: 'Фиолетовый', color: '#9333ea', bg: 'bg-purple-600' },
              ].map((c) => (
                <button
                  key={c.color}
                  type="button"
                  title={c.label}
                  {...preventBtnFocus}
                  onClick={(e) => {
                    e.stopPropagation();
                    formatColor(c.color);
                    setShowColorPicker(false);
                  }}
                  className={`w-5 h-5 rounded-full ${c.bg} border-2 border-white dark:border-slate-900 shadow-xs hover:scale-110 transition-transform`}
                />
              ))}
            </div>
          )}
        </div>

        {/* Alignment: Left, Center, Right, Justify */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="По левому краю"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyLeft'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignLeft size={14} />
          </button>
          <button
            type="button"
            title="По центру"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyCenter'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignCenter size={14} />
          </button>
          <button
            type="button"
            title="По правому краю"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyRight'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignRight size={14} />
          </button>
          <button
            type="button"
            title="По ширине"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyFull'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignJustify size={14} />
          </button>
        </div>

        {/* Font styling: Bold, Italic, Underline, Strikethrough, Code, Highlighter */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Жирный (Ctrl+B)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('bold'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 font-extrabold"
          >
            <Bold size={14} />
          </button>
          <button
            type="button"
            title="Курсив (Ctrl+I)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('italic'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 italic"
          >
            <Italic size={14} />
          </button>
          <button
            type="button"
            title="Подчеркнутый (Ctrl+U)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('underline'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Underline size={14} />
          </button>
          <button
            type="button"
            title="Зачеркнутый"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('strikeThrough'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Strikethrough size={14} />
          </button>
          <button
            type="button"
            title="Код (inline / block)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatCode(); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 font-mono"
          >
            <Code size={14} />
          </button>
          <button
            type="button"
            title="Выделитель желтым"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatHighlight('#fef08a'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-yellow-100 dark:hover:bg-yellow-950/60 text-yellow-600"
          >
            <Highlighter size={14} />
          </button>
        </div>

        {/* Lists & Indentation */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Маркированный список"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('insertUnorderedList'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400"
          >
            <List size={14} />
          </button>
          <button
            type="button"
            title="Нумерованный список"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('insertOrderedList'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400"
          >
            <ListOrdered size={14} />
          </button>
          <button
            type="button"
            title="Увеличить отступ (Tab)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('indent'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Indent size={14} />
          </button>
          <button
            type="button"
            title="Уменьшить отступ (Shift+Tab)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('outdent'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Outdent size={14} />
          </button>
        </div>

        {/* Rich inserts: Quote, Callout, Image, Link, Import Word/MD, Divider, Clear */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Цитата (Блок с акцентом)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('formatBlock', '<blockquote>'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-amber-600"
          >
            <Quote size={14} />
          </button>

          {/* Callout Dropdown (💡 Совет, ⚠️ Важно, 📝 Заметка, 🎯 Цель) */}
          <div className="relative">
            <button
              type="button"
              title="Вставить цветную врезку (Callout)"
              {...preventBtnFocus}
              onClick={(e) => { e.stopPropagation(); setShowCalloutMenu(!showCalloutMenu); }}
              className="edtech-inline-btn p-1 rounded hover:bg-indigo-50 dark:hover:bg-indigo-950 text-indigo-600"
            >
              <Info size={14} />
            </button>
            {showCalloutMenu && (
              <div
                className="absolute top-full left-0 mt-1 p-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl shadow-2xl flex flex-col gap-1 z-50 min-w-[200px]"
                onPointerDownCapture={(e) => e.stopPropagation()}
                onMouseDownCapture={(e) => e.stopPropagation()}
              >
                <button
                  type="button"
                  {...preventBtnFocus}
                  onClick={(e) => { e.stopPropagation(); handleInsertCallout('tip'); }}
                  className="flex items-center gap-2 px-2.5 py-1.5 text-xs rounded-lg hover:bg-indigo-50 dark:hover:bg-indigo-950/60 text-slate-700 dark:text-slate-200 text-left"
                >
                  <span className="text-base">💡</span>
                  <div>
                    <div className="font-bold text-indigo-600 dark:text-indigo-400">Совет</div>
                    <div className="text-[10px] text-slate-400">Полезная рекомендация</div>
                  </div>
                </button>
                <button
                  type="button"
                  {...preventBtnFocus}
                  onClick={(e) => { e.stopPropagation(); handleInsertCallout('warning'); }}
                  className="flex items-center gap-2 px-2.5 py-1.5 text-xs rounded-lg hover:bg-amber-50 dark:hover:bg-amber-950/60 text-slate-700 dark:text-slate-200 text-left"
                >
                  <span className="text-base">⚠️</span>
                  <div>
                    <div className="font-bold text-amber-600 dark:text-amber-400">Важно</div>
                    <div className="text-[10px] text-slate-400">Внимание к деталям</div>
                  </div>
                </button>
                <button
                  type="button"
                  {...preventBtnFocus}
                  onClick={(e) => { e.stopPropagation(); handleInsertCallout('note'); }}
                  className="flex items-center gap-2 px-2.5 py-1.5 text-xs rounded-lg hover:bg-emerald-50 dark:hover:bg-emerald-950/60 text-slate-700 dark:text-slate-200 text-left"
                >
                  <span className="text-base">📝</span>
                  <div>
                    <div className="font-bold text-emerald-600 dark:text-emerald-400">Заметка</div>
                    <div className="text-[10px] text-slate-400">Информационный блок</div>
                  </div>
                </button>
                <button
                  type="button"
                  {...preventBtnFocus}
                  onClick={(e) => { e.stopPropagation(); handleInsertCallout('target'); }}
                  className="flex items-center gap-2 px-2.5 py-1.5 text-xs rounded-lg hover:bg-purple-50 dark:hover:bg-purple-950/60 text-slate-700 dark:text-slate-200 text-left"
                >
                  <span className="text-base">🎯</span>
                  <div>
                    <div className="font-bold text-purple-600 dark:text-purple-400">Цель</div>
                    <div className="text-[10px] text-slate-400">Ожидаемый результат</div>
                  </div>
                </button>
              </div>
            )}
          </div>

          <button
            type="button"
            title="Вставить картинку (загрузка с диска или URL)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); setShowImageModal(!showImageModal); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-sky-600"
          >
            <ImageIcon size={14} />
          </button>
          <button
            type="button"
            title="Вставить ссылку"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); setShowLinkModal(!showLinkModal); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-blue-600"
          >
            <LinkIcon size={14} />
          </button>

          {/* Import Word / Markdown Button */}
          <button
            type="button"
            title="Импортировать документ Word (.docx) или Markdown (.md)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); fileInputRef.current?.click(); }}
            className="edtech-inline-btn px-2 py-1 rounded text-purple-700 dark:text-purple-300 hover:bg-purple-50 dark:hover:bg-purple-950/60 font-semibold flex items-center gap-1"
          >
            <FileUp size={13} />
            <span className="text-[11px] hidden sm:inline">Импорт .docx / .md</span>
          </button>
          <input
            ref={fileInputRef}
            type="file"
            accept=".docx,.md,.markdown"
            className="hidden"
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (file) handleFileImport(file);
              e.target.value = '';
            }}
          />

          <button
            type="button"
            title="Разделительная линия"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('insertHorizontalRule'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-500"
          >
            <Minus size={14} />
          </button>
          <button
            type="button"
            title="Очистить форматирование"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); clearFormatting(); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-rose-500"
          >
            <RemoveFormatting size={14} />
          </button>
        </div>
      </div>

      {/* Popover Modal for Insert Image (Upload & URL Tabs) */}
      {showImageModal && (
        <div
          data-puck-overlay-portal="true"
          onPointerDownCapture={(e) => e.stopPropagation()}
          onMouseDownCapture={(e) => e.stopPropagation()}
          className="p-3 bg-indigo-50/90 dark:bg-slate-950 border-b border-indigo-200 dark:border-indigo-900 space-y-3"
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <span className="text-xs font-bold text-indigo-900 dark:text-indigo-200">Вставка изображения</span>
              {/* Tab Selector */}
              <div className="flex items-center bg-white dark:bg-slate-900 rounded-lg p-0.5 border border-indigo-200 dark:border-indigo-800 text-[11px]">
                <button
                  type="button"
                  onClick={() => setImageTab('upload')}
                  className={`px-2 py-0.5 rounded font-bold ${
                    imageTab === 'upload'
                      ? 'bg-indigo-600 text-white shadow-2xs'
                      : 'text-slate-600 dark:text-slate-300 hover:text-indigo-600'
                  }`}
                >
                  📁 С устройства
                </button>
                <button
                  type="button"
                  onClick={() => setImageTab('url')}
                  className={`px-2 py-0.5 rounded font-bold ${
                    imageTab === 'url'
                      ? 'bg-indigo-600 text-white shadow-2xs'
                      : 'text-slate-600 dark:text-slate-300 hover:text-indigo-600'
                  }`}
                >
                  🔗 По ссылке
                </button>
              </div>
            </div>
            <button
              onClick={() => {
                setShowImageModal(false);
                setImageFile(null);
                setImagePreview(null);
              }}
              className="p-1 rounded-lg hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-500"
            >
              <X size={14} />
            </button>
          </div>

          {imageTab === 'upload' ? (
            <div className="space-y-2">
              <input
                ref={imageFileInputRef}
                type="file"
                accept="image/*"
                className="hidden"
                onChange={(e) => {
                  const file = e.target.files?.[0];
                  if (file) {
                    setImageFile(file);
                    const reader = new FileReader();
                    reader.onload = () => setImagePreview(reader.result as string);
                    reader.readAsDataURL(file);
                  }
                }}
              />
              {imagePreview ? (
                <div className="flex items-center gap-3 p-2 bg-white dark:bg-slate-900 rounded-xl border border-indigo-200 dark:border-indigo-800">
                  <img src={imagePreview} alt="Preview" className="w-16 h-16 object-cover rounded-lg border" />
                  <div className="flex-1 text-xs">
                    <p className="font-bold text-slate-800 dark:text-slate-200 truncate">{imageFile?.name}</p>
                    <p className="text-[10px] text-slate-400">
                      {imageFile?.size ? `${(imageFile.size / 1024).toFixed(1)} KB` : ''}
                    </p>
                    <button
                      type="button"
                      onClick={() => imageFileInputRef.current?.click()}
                      className="mt-1 text-[11px] text-indigo-600 hover:underline font-semibold"
                    >
                      Выбрать другой файл
                    </button>
                  </div>
                </div>
              ) : (
                <div
                  onClick={() => imageFileInputRef.current?.click()}
                  className="p-4 rounded-xl border-2 border-dashed border-indigo-300 dark:border-indigo-800 hover:border-indigo-500 bg-white/60 dark:bg-slate-900/60 text-center cursor-pointer transition-colors"
                >
                  <Upload size={20} className="mx-auto text-indigo-600 dark:text-indigo-400 mb-1" />
                  <p className="text-xs font-bold text-slate-700 dark:text-slate-300">
                    Нажмите, чтобы выбрать файл изображения
                  </p>
                  <p className="text-[10px] text-slate-400">PNG, JPG, WebP, GIF или вставьте через Ctrl+V прямо в текст</p>
                </div>
              )}
            </div>
          ) : (
            <input
              type="text"
              value={imageUrl}
              onChange={(e) => setImageUrl(e.target.value)}
              placeholder="Вставьте URL изображения (https://...)..."
              className="px-3 py-1.5 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-800 dark:text-slate-100 w-full"
            />
          )}

          <div className="grid grid-cols-1 md:grid-cols-3 gap-2 items-center">
            <input
              type="text"
              value={imageCaption}
              onChange={(e) => setImageCaption(e.target.value)}
              placeholder="Подпись к картинке (опционально)..."
              className="px-3 py-1.5 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-800 dark:text-slate-100 md:col-span-2"
            />
            <div className="flex items-center gap-2">
              <select
                value={imageWidth}
                onChange={(e) => setImageWidth(e.target.value as any)}
                className="px-2 py-1.5 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-800 dark:text-slate-100 flex-1"
              >
                <option value="100%">100% ширины</option>
                <option value="75%">75% ширины</option>
                <option value="50%">50% ширины</option>
              </select>
              <button
                type="button"
                onClick={handleInsertImage}
                disabled={(!imageFile && !imageUrl.trim()) || isUploadingImage}
                className="px-4 py-1.5 rounded-xl bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold flex items-center gap-1.5 disabled:opacity-40"
              >
                {isUploadingImage ? <Loader2 size={13} className="animate-spin" /> : <Check size={14} />}
                Вставить
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Popover Modal for Insert Link */}
      {showLinkModal && (
        <div
          data-puck-overlay-portal="true"
          onPointerDownCapture={(e) => e.stopPropagation()}
          onMouseDownCapture={(e) => e.stopPropagation()}
          className="p-3 bg-indigo-50/90 dark:bg-slate-950 border-b border-indigo-200 dark:border-indigo-900 flex items-center gap-2"
        >
          <input
            type="text"
            value={linkUrl}
            onChange={(e) => setLinkUrl(e.target.value)}
            placeholder="Введите URL ссылки (напр. https://example.com)..."
            className="px-3 py-1.5 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-800 dark:text-slate-100 flex-1"
          />
          <button
            type="button"
            onClick={handleInsertLink}
            disabled={!linkUrl.trim()}
            className="px-4 py-1.5 rounded-xl bg-indigo-600 text-white text-xs font-bold disabled:opacity-40"
          >
            Применить
          </button>
          <button
            onClick={() => setShowLinkModal(false)}
            className="p-1.5 rounded-xl hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-500"
          >
            <X size={14} />
          </button>
        </div>
      )}

      {/* Document Sheet Container with Drag-and-Drop and Overlays */}
      <div
        className="relative"
        onDragOver={(e) => {
          e.preventDefault();
          e.stopPropagation();
          setIsDragOver(true);
        }}
        onDragLeave={(e) => {
          e.preventDefault();
          e.stopPropagation();
          setIsDragOver(false);
        }}
        onDrop={(e) => {
          e.preventDefault();
          e.stopPropagation();
          setIsDragOver(false);
          const file = e.dataTransfer.files?.[0];
          if (file) handleFileImport(file);
        }}
      >
        {/* Drag & Drop Overlay */}
        {isDragOver && (
          <div className="absolute inset-0 z-30 bg-indigo-500/10 backdrop-blur-xs border-4 border-dashed border-indigo-500 rounded-2xl flex flex-col items-center justify-center p-6 text-center animate-in fade-in">
            <div className="w-12 h-12 rounded-2xl bg-indigo-600 text-white flex items-center justify-center mb-2 shadow-lg">
              <FileUp size={24} />
            </div>
            <h4 className="text-sm font-extrabold text-indigo-900 dark:text-indigo-200">
              Отпустите файл Word (.docx) или Markdown (.md)
            </h4>
            <p className="text-xs text-indigo-700 dark:text-indigo-300 mt-0.5">
              Текст и форматирование будут добавлены в этот текстовый блок
            </p>
          </div>
        )}

        {/* Importing Progress Overlay */}
        {isImporting && (
          <div className="absolute inset-0 z-40 bg-white/80 dark:bg-slate-900/80 backdrop-blur-xs flex flex-col items-center justify-center p-6 text-center animate-in fade-in">
            <div className="w-10 h-10 rounded-2xl bg-indigo-100 dark:bg-indigo-950/80 text-indigo-600 flex items-center justify-center mb-2 shadow-md">
              <Loader2 size={20} className="animate-spin" />
            </div>
            <h4 className="text-xs font-extrabold text-slate-900 dark:text-white">Конвертация документа...</h4>
            <p className="text-[11px] text-indigo-600 dark:text-indigo-400 font-semibold mt-0.5">
              {importStatus || 'Пожалуйста, подождите...'}
            </p>
          </div>
        )}

        {/* Editable Body */}
        <div
          ref={editorRef}
          contentEditable
          suppressContentEditableWarning
          dangerouslySetInnerHTML={{ __html: initialHtmlRef.current }}
          data-puck-overlay-portal="true"
          onPointerDownCapture={(e) => e.stopPropagation()}
          onMouseDownCapture={(e) => e.stopPropagation()}
          onKeyDownCapture={(e) => e.stopPropagation()}
          onKeyUpCapture={(e) => e.stopPropagation()}
          onPointerDown={(e) => e.stopPropagation()}
          onMouseDown={(e) => e.stopPropagation()}
          onMouseUp={saveSelection}
          onKeyUp={saveSelection}
          onSelect={saveSelection}
          onFocus={() => {
            onFocusBlock?.();
            saveSelection();
          }}
          onInput={triggerChange}
          onBlur={triggerChange}
          style={{ userSelect: 'text', WebkitUserSelect: 'text', pointerEvents: 'auto' }}
          className="edtech-inline-editable p-5 min-h-[140px] focus:outline-none prose dark:prose-invert max-w-none text-slate-800 dark:text-slate-200 leading-relaxed text-sm md:text-base cursor-text select-text"
          data-placeholder={placeholder}
        />

        {/* Import Success Toast */}
        {importToast && (
          <div className="absolute bottom-3 right-3 z-50 px-3.5 py-2 rounded-xl bg-indigo-600 text-white text-xs font-bold shadow-xl flex items-center gap-2 animate-in slide-in-from-bottom-2">
            <Check size={14} />
            <span>{importToast}</span>
          </div>
        )}
      </div>
    </div>
  );
}

/**
 * Word-style Rich Text Editor Component for full lecture articles with images, formatting, lists, callouts.
 */
interface RichTextWordEditorProps {
  htmlContent: string;
  onChange: (html: string) => void;
  title?: string;
  onTitleChange?: (title: string) => void;
  isEditing?: boolean;
  dragRef?: any;
  onFocusBlock?: () => void;
}

export function RichTextWordEditor({
  htmlContent,
  onChange,
  title,
  onTitleChange,
  isEditing = false,
  dragRef,
  onFocusBlock,
}: RichTextWordEditorProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const editorRef = useRef<HTMLDivElement>(null);
  const isInitializedRef = useRef(false);
  const initialHtmlRef = useRef<string>(htmlContent || '');
  const lastHtmlRef = useRef<string>(htmlContent || '');
  const savedSelectionRef = useRef<Range | null>(null);

  const [showImageModal, setShowImageModal] = useState(false);
  const [imageTab, setImageTab] = useState<'upload' | 'url'>('upload');
  const [imageUrl, setImageUrl] = useState('');
  const [imageCaption, setImageCaption] = useState('');
  const [imageWidth, setImageWidth] = useState<'100%' | '75%' | '50%'>('100%');
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [isUploadingImage, setIsUploadingImage] = useState(false);
  const imageFileInputRef = useRef<HTMLInputElement>(null);

  const [showLinkModal, setShowLinkModal] = useState(false);
  const [linkUrl, setLinkUrl] = useState('');
  const [showColorPicker, setShowColorPicker] = useState(false);
  const [showCalloutMenu, setShowCalloutMenu] = useState(false);

  // File import state
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isImporting, setIsImporting] = useState(false);
  const [importStatus, setImportStatus] = useState<string>('');
  const [isDragOver, setIsDragOver] = useState(false);
  const [importToast, setImportToast] = useState<string | null>(null);

  const handleFileImport = async (file: File) => {
    const ext = file.name.toLowerCase();
    if (!ext.endsWith('.docx') && !ext.endsWith('.md') && !ext.endsWith('.markdown')) {
      alert('Пожалуйста, выберите файл документа Microsoft Word (.docx) или Markdown (.md)');
      return;
    }

    setIsImporting(true);
    setImportStatus('Чтение файла...');
    try {
      const res = await convertDocumentToHtml(file, (msg) => setImportStatus(msg));

      if (editorRef.current) {
        const currentHtml = editorRef.current.innerHTML.trim();
        const isDefault = !currentHtml || currentHtml === '<p>Пустой текст лекции</p>' || currentHtml === '<p><br></p>';
        const newHtml = isDefault ? res.html : `${currentHtml}<hr class="my-6 border-slate-200 dark:border-slate-800" />${res.html}`;
        editorRef.current.innerHTML = newHtml;
        lastHtmlRef.current = newHtml;
        onChange(newHtml);
      }

      if (res.extractedTitle && onTitleChange && !title) {
        onTitleChange(res.extractedTitle);
      }

      setImportToast(`Документ успешно импортирован (${res.wordCount} слов, ${res.tablesCount} таблиц)`);
      setTimeout(() => setImportToast(null), 4000);
    } catch (err: any) {
      console.error('Import failed', err);
      alert(err.message || 'Ошибка импорта документа');
    } finally {
      setIsImporting(false);
      setImportStatus('');
    }
  };

  const getEditorDoc = useCallback(() => {
    return editorRef.current?.ownerDocument || (typeof document !== 'undefined' ? document : null);
  }, []);

  const getEditorWin = useCallback(() => {
    const doc = getEditorDoc();
    return doc?.defaultView || (typeof window !== 'undefined' ? window : null);
  }, [getEditorDoc]);

  // Save current selection range inside the editor
  const saveSelection = useCallback(() => {
    const win = getEditorWin();
    if (!win) return;
    const sel = win.getSelection();
    if (sel && sel.rangeCount > 0) {
      const range = sel.getRangeAt(0);
      if (
        editorRef.current &&
        (editorRef.current === range.commonAncestorContainer || editorRef.current.contains(range.commonAncestorContainer))
      ) {
        savedSelectionRef.current = range.cloneRange();
      }
    }
  }, [getEditorWin]);

  // Restore saved selection before applying commands
  const restoreSelection = useCallback(() => {
    const win = getEditorWin();
    if (!win || !editorRef.current) return;
    const sel = win.getSelection();
    if (!sel || !savedSelectionRef.current) return;
    if (sel.rangeCount > 0) {
      const curRange = sel.getRangeAt(0);
      if (
        curRange.startContainer === savedSelectionRef.current.startContainer &&
        curRange.startOffset === savedSelectionRef.current.startOffset &&
        curRange.endContainer === savedSelectionRef.current.endContainer &&
        curRange.endOffset === savedSelectionRef.current.endOffset
      ) {
        return;
      }
    }
    try {
      sel.removeAllRanges();
      sel.addRange(savedSelectionRef.current);
    } catch (e) {}
  }, [getEditorWin]);

  // Track selection changes across both editor ownerDocument and top document
  useEffect(() => {
    const doc = getEditorDoc();
    const handleSelectionChange = () => {
      saveSelection();
    };
    doc?.addEventListener('selectionchange', handleSelectionChange);
    if (typeof document !== 'undefined' && document !== doc) {
      document.addEventListener('selectionchange', handleSelectionChange);
    }
    return () => {
      doc?.removeEventListener('selectionchange', handleSelectionChange);
      if (typeof document !== 'undefined' && document !== doc) {
        document.removeEventListener('selectionchange', handleSelectionChange);
      }
    };
  }, [getEditorDoc, saveSelection]);

  // Set initial content on mount once
  useEffect(() => {
    if (editorRef.current && !isInitializedRef.current) {
      editorRef.current.innerHTML = htmlContent || '';
      lastHtmlRef.current = htmlContent || '';
      isInitializedRef.current = true;
    }
    const doc = getEditorDoc();
    try {
      doc?.execCommand('defaultParagraphSeparator', false, 'p');
      doc?.execCommand('styleWithCSS', false, 'true');
    } catch (e) {}
  }, [getEditorDoc, htmlContent]);

  // Sync external changes ONLY when prop actually changed from outside and user is not typing inside
  useEffect(() => {
    if (editorRef.current && isInitializedRef.current) {
      const doc = getEditorDoc();
      const isFocused =
        Boolean(doc?.activeElement && (doc.activeElement === editorRef.current || editorRef.current.contains(doc.activeElement)));
      if (!isFocused && htmlContent !== undefined && htmlContent !== lastHtmlRef.current) {
        lastHtmlRef.current = htmlContent;
        editorRef.current.innerHTML = htmlContent || '';
      }
    }
  }, [getEditorDoc, htmlContent]);

  const triggerChange = useCallback(() => {
    if (!editorRef.current) return;
    const newHtml = editorRef.current.innerHTML;
    lastHtmlRef.current = newHtml;
    onChange(newHtml);
  }, [onChange]);

  const exec = (command: string, value: string | undefined = undefined) => {
    const doc = getEditorDoc();
    if (!doc || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    try {
      doc.execCommand('styleWithCSS', false, 'false');
    } catch (e) {}
    doc.execCommand(command, false, value);
    saveSelection();
    triggerChange();
  };

  // Native keyboard event isolation, Tab indent in lists/text, and Clipboard Paste (Ctrl+V) handler
  useEffect(() => {
    const el = editorRef.current;
    if (!el) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      // Isolate keydown from Puck hotkeys
      e.stopPropagation();

      // Tab key: Indentation in lists or paragraph indent in text
      if (e.key === 'Tab') {
        e.preventDefault();
        const doc = getEditorDoc();
        const win = getEditorWin();
        const sel = win?.getSelection();

        let isInList = false;
        if (sel && sel.rangeCount > 0) {
          let node: Node | null = sel.getRangeAt(0).startContainer;
          while (node && node !== editorRef.current) {
            if (node.nodeName === 'LI' || node.nodeName === 'UL' || node.nodeName === 'OL') {
              isInList = true;
              break;
            }
            node = node.parentNode;
          }
        }

        if (isInList) {
          if (e.shiftKey) {
            exec('outdent');
          } else {
            exec('indent');
          }
        } else {
          doc?.execCommand('insertHTML', false, '&emsp;&emsp;');
          triggerChange();
        }
        return;
      }

      // Hotkeys: Ctrl+B, Ctrl+I, Ctrl+U, Ctrl+Z, Ctrl+Y
      if (e.ctrlKey || e.metaKey) {
        const k = e.key.toLowerCase();
        if (k === 'b' || k === 'и') {
          e.preventDefault();
          exec('bold');
          return;
        }
        if (k === 'i' || k === 'ш') {
          e.preventDefault();
          exec('italic');
          return;
        }
        if (k === 'u' || k === 'г') {
          e.preventDefault();
          exec('underline');
          return;
        }
        if (k === 'z' || k === 'я') {
          if (e.shiftKey) {
            e.preventDefault();
            exec('redo');
          } else {
            e.preventDefault();
            exec('undo');
          }
          return;
        }
        if (k === 'y' || k === 'н') {
          e.preventDefault();
          exec('redo');
          return;
        }
      }
    };

    // Intercept image pastes from clipboard (Ctrl+V screenshot / copied image)
    const handlePaste = async (e: ClipboardEvent) => {
      const items = e.clipboardData?.items;
      if (!items) return;

      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        if (item.type.indexOf('image') !== -1) {
          e.preventDefault();
          e.stopPropagation();
          const file = item.getAsFile();
          if (!file) continue;

          try {
            const imgUrl = await uploadOrConvertImage(file);
            const snippet = `
              <figure class="my-6 mx-auto w-full text-center">
                <img src="${imgUrl}" alt="Изображение из буфера" class="rounded-2xl shadow-lg border border-slate-200 dark:border-slate-800 mx-auto object-cover max-h-[480px]" />
              </figure>
              <p><br/></p>
            `;
            insertCustomHtml(snippet);
          } catch (err) {
            console.error('Failed to paste image:', err);
          }
          return;
        }
      }
    };

    const stopPropagation = (e: KeyboardEvent) => {
      e.stopPropagation();
    };

    el.addEventListener('keydown', handleKeyDown, false);
    el.addEventListener('keyup', stopPropagation, false);
    el.addEventListener('keypress', stopPropagation, false);
    el.addEventListener('paste', handlePaste, false);

    return () => {
      el.removeEventListener('keydown', handleKeyDown, false);
      el.removeEventListener('keyup', stopPropagation, false);
      el.removeEventListener('keypress', stopPropagation, false);
      el.removeEventListener('paste', handlePaste, false);
    };
  }, [getEditorDoc, getEditorWin, triggerChange]);

  const formatHeading = (tag: 'p' | 'h1' | 'h2' | 'h3') => {
    const doc = getEditorDoc();
    if (!doc || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    try {
      doc.execCommand('styleWithCSS', false, 'false');
    } catch (e) {}
    let ok = doc.execCommand('formatBlock', false, `<${tag}>`);
    if (!ok) {
      doc.execCommand('formatBlock', false, tag);
    }
    saveSelection();
    triggerChange();
  };

  const applyInlineStyle = (styleObj: Record<string, string>) => {
    const doc = getEditorDoc();
    const win = getEditorWin();
    if (!doc || !win || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    const sel = win.getSelection();
    if (!sel || sel.rangeCount === 0) return;
    const range = sel.getRangeAt(0);
    if (!editorRef.current.contains(range.commonAncestorContainer) && editorRef.current !== range.commonAncestorContainer) return;

    if (range.collapsed) {
      let el: HTMLElement | null =
        (range.startContainer.nodeType === Node.ELEMENT_NODE
          ? range.startContainer
          : range.startContainer.parentElement) as HTMLElement;
      while (el && el !== editorRef.current && el.parentElement !== editorRef.current) {
        el = el.parentElement;
      }
      if (el && el !== editorRef.current) {
        Object.assign(el.style, styleObj);
      } else {
        Object.assign(editorRef.current.style, styleObj);
      }
      triggerChange();
      return;
    }

    const contents = range.extractContents();
    const span = doc.createElement('span');
    Object.assign(span.style, styleObj);
    span.appendChild(contents);
    range.insertNode(span);

    const newRange = doc.createRange();
    newRange.selectNodeContents(span);
    sel.removeAllRanges();
    sel.addRange(newRange);
    savedSelectionRef.current = newRange.cloneRange();
    triggerChange();
  };

  const formatFont = (fontFamily: string) => {
    applyInlineStyle({ fontFamily });
  };

  const formatSize = (sizePx: string) => {
    applyInlineStyle({ fontSize: sizePx });
  };

  const formatColor = (color: string) => {
    applyInlineStyle({ color });
  };

  const formatAlign = (alignCmd: 'justifyLeft' | 'justifyCenter' | 'justifyRight' | 'justifyFull') => {
    const doc = getEditorDoc();
    const win = getEditorWin();
    if (!doc || !win || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    try {
      doc.execCommand('styleWithCSS', false, 'false');
    } catch (e) {}
    doc.execCommand(alignCmd, false);
    const sel = win.getSelection();
    if (!sel || sel.rangeCount === 0 || sel.isCollapsed) {
      const map: Record<string, string> = {
        justifyLeft: 'left',
        justifyCenter: 'center',
        justifyRight: 'right',
        justifyFull: 'justify',
      };
      if (map[alignCmd]) {
        editorRef.current.style.textAlign = map[alignCmd];
      }
    }
    saveSelection();
    triggerChange();
  };

  const formatHighlight = (color: string = '#fef08a') => {
    applyInlineStyle({ backgroundColor: color });
  };

  const formatCode = () => {
    const doc = getEditorDoc();
    const win = getEditorWin();
    if (!doc || !win || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    const sel = win.getSelection();
    if (sel && !sel.isCollapsed && sel.rangeCount > 0) {
      const range = sel.getRangeAt(0);
      if (editorRef.current.contains(range.commonAncestorContainer)) {
        const contents = range.extractContents();
        const code = doc.createElement('code');
        code.className = 'px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-indigo-600 dark:text-indigo-400 font-mono text-xs';
        code.appendChild(contents);
        range.insertNode(code);

        const newRange = doc.createRange();
        newRange.selectNodeContents(code);
        sel.removeAllRanges();
        sel.addRange(newRange);
        savedSelectionRef.current = newRange.cloneRange();
        triggerChange();
        return;
      }
    }
    formatHeading('p');
    exec('formatBlock', '<pre>');
  };

  const clearFormatting = () => {
    const doc = getEditorDoc();
    const win = getEditorWin();
    if (!doc || !win || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    try {
      doc.execCommand('styleWithCSS', false, 'false');
    } catch (e) {}
    doc.execCommand('removeFormat', false);
    doc.execCommand('unlink', false);
    const sel = win.getSelection();
    if (sel && !sel.isCollapsed && sel.rangeCount > 0) {
      const range = sel.getRangeAt(0);
      if (editorRef.current.contains(range.commonAncestorContainer)) {
        let el = range.commonAncestorContainer as HTMLElement;
        if (el.nodeType !== Node.ELEMENT_NODE) el = el.parentElement as HTMLElement;
        while (el && el !== editorRef.current) {
          if (el.tagName === 'SPAN' || el.tagName === 'FONT' || el.tagName === 'MARK') {
            el.removeAttribute('style');
            el.removeAttribute('face');
            el.removeAttribute('size');
            el.removeAttribute('color');
          }
          el = el.parentElement as HTMLElement;
        }
      }
    }
    saveSelection();
    triggerChange();
  };

  const insertCustomHtml = (htmlSnippet: string) => {
    const doc = getEditorDoc();
    if (!doc || !editorRef.current) return;
    editorRef.current.focus();
    restoreSelection();
    doc.execCommand('insertHTML', false, htmlSnippet);
    saveSelection();
    triggerChange();
  };

  const handleInsertImage = async () => {
    let finalUrl = imageUrl.trim();

    if (imageTab === 'upload' && imageFile) {
      setIsUploadingImage(true);
      try {
        finalUrl = await uploadOrConvertImage(imageFile);
      } catch (err) {
        alert('Не удалось загрузить изображение');
        setIsUploadingImage(false);
        return;
      }
      setIsUploadingImage(false);
    }

    if (!finalUrl) return;

    const widthClass =
      imageWidth === '50%' ? 'max-w-[50%]' : imageWidth === '75%' ? 'max-w-[75%]' : 'w-full';

    const snippet = `
      <figure class="my-6 mx-auto ${widthClass} text-center">
        <img src="${finalUrl}" alt="${imageCaption || 'Иллюстрация к лекции'}" class="rounded-2xl shadow-lg border border-slate-200 dark:border-slate-800 mx-auto object-cover max-h-[480px]" />
        ${imageCaption ? `<figcaption class="mt-2 text-xs text-slate-500 dark:text-slate-400 font-medium">${imageCaption}</figcaption>` : ''}
      </figure>
      <p><br/></p>
    `;
    insertCustomHtml(snippet);
    setImageUrl('');
    setImageCaption('');
    setImageFile(null);
    setImagePreview(null);
    setShowImageModal(false);
  };

  const handleInsertCallout = (type: CalloutType) => {
    insertCustomHtml(getCalloutSnippet(type));
    setShowCalloutMenu(false);
  };

  const handleInsertLink = () => {
    if (!linkUrl.trim()) return;
    let url = linkUrl.trim();
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      url = 'https://' + url;
    }
    exec('createLink', url);
    setLinkUrl('');
    setShowLinkModal(false);
  };

  if (!isEditing) {
    return (
      <div className="my-6 space-y-4">
        {title && (
          <h2 className="text-2xl font-extrabold text-slate-900 dark:text-white tracking-tight">
            {title}
          </h2>
        )}
        <div
          className="prose dark:prose-invert max-w-none text-slate-700 dark:text-slate-300 leading-relaxed text-sm md:text-base"
          dangerouslySetInnerHTML={{ __html: htmlContent || '<p>Пустой текст лекции</p>' }}
        />
      </div>
    );
  }

  const preventBtnFocus = {
    onPointerDown: (e: React.PointerEvent) => {
      e.preventDefault();
      e.stopPropagation();
    },
    onMouseDown: (e: React.MouseEvent) => {
      e.preventDefault();
      e.stopPropagation();
    },
  };

  return (
    <div
      ref={containerRef}
      data-puck-overlay-portal="true"
      onClick={() => onFocusBlock?.()}
      className="edtech-inline-editable my-6 rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-lg overflow-hidden transition-all relative"
    >
      {/* Embedded Scoped Headings & Lists Styles */}
      <style dangerouslySetInnerHTML={{ __html: EDITOR_SCOPED_STYLES }} />

      {/* Top Banner & Header */}
      <div className="bg-slate-100/90 dark:bg-slate-800/90 px-5 py-3 border-b border-slate-200 dark:border-slate-700 flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          {/* Badge */}
          <div
            ref={dragRef}
            data-puck-drag-handle={dragRef ? 'true' : undefined}
            className="flex items-center gap-1.5 px-2.5 py-1 rounded-xl bg-indigo-600 text-white text-xs font-bold select-none shadow-2xs"
            title="Интерактивная статья / лекция"
          >
            <FileText size={14} />
            <span>Лекция / Word</span>
          </div>
        </div>

        {/* Section title inline editable */}
        {onTitleChange && (
          <div className="flex items-center gap-2 flex-1 max-w-md ml-auto" data-puck-overlay-portal="true">
            <span className="text-[11px] font-bold text-slate-500 whitespace-nowrap">Заголовок раздела:</span>
            <input
              type="text"
              value={title || ''}
              onChange={(e) => onTitleChange(e.target.value)}
              onPointerDownCapture={(e) => e.stopPropagation()}
              onMouseDownCapture={(e) => e.stopPropagation()}
              onClickCapture={(e) => e.stopPropagation()}
              data-puck-overlay-portal="true"
              placeholder="Введите название темы..."
              className="edtech-inline-editable px-2.5 py-1 text-xs font-bold rounded-lg bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-700 w-full focus:outline-indigo-500 text-slate-800 dark:text-slate-100"
            />
          </div>
        )}
      </div>

      {/* Word-style Ribbon Toolbar */}
      <div
        data-puck-overlay-portal="true"
        onPointerDown={(e) => {
          const target = e.target as HTMLElement;
          if (target && (target.tagName === 'INPUT' || target.tagName === 'SELECT' || target.tagName === 'TEXTAREA')) {
            return;
          }
          e.preventDefault();
          e.stopPropagation();
        }}
        onMouseDown={(e) => {
          const target = e.target as HTMLElement;
          if (target && (target.tagName === 'INPUT' || target.tagName === 'SELECT' || target.tagName === 'TEXTAREA')) {
            return;
          }
          e.preventDefault();
          e.stopPropagation();
        }}
        className="bg-slate-50 dark:bg-slate-950 px-4 py-2 border-b border-slate-200 dark:border-slate-800 flex flex-wrap items-center gap-1.5 text-slate-700 dark:text-slate-300 select-none text-xs"
      >
        {/* Undo / Redo */}
        <div className="flex items-center bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Отменить действие (Ctrl+Z)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('undo'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300"
          >
            <Undo size={14} />
          </button>
          <button
            type="button"
            title="Повторить действие (Ctrl+Y)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('redo'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300"
          >
            <Redo size={14} />
          </button>
        </div>

        {/* Style / Headings Dropdown */}
        <div className="flex items-center bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-0.5" data-puck-overlay-portal="true">
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Обычный текст (Абзац)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatHeading('p'); }}
            className="edtech-inline-btn px-2 py-1 text-xs font-bold rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            P
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Заголовок H1"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatHeading('h1'); }}
            className="edtech-inline-btn p-1 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400 font-extrabold"
          >
            <Heading1 size={15} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Заголовок H2"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatHeading('h2'); }}
            className="edtech-inline-btn p-1 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400 font-bold"
          >
            <Heading2 size={15} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Подзаголовок H3"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatHeading('h3'); }}
            className="edtech-inline-btn p-1 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400"
          >
            <Heading3 size={15} />
          </button>
        </div>

        {/* Font Family (Шрифт) */}
        <div className="flex items-center bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Шрифт: Без засечек (Sans)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatFont('sans-serif'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[11px] font-medium rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            Sans
          </button>
          <button
            type="button"
            title="Шрифт: С засечками (Serif / Книга)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatFont('Georgia, serif'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[11px] font-serif rounded hover:bg-slate-100 dark:hover:bg-slate-800 italic"
          >
            Serif
          </button>
          <button
            type="button"
            title="Шрифт: Моноширинный (Mono)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatFont('monospace'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[11px] font-mono rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            Mono
          </button>
        </div>

        {/* Font Size (Размер текста) */}
        <div className="flex items-center bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Размер: Мелкий"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatSize('13px'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[10px] font-bold rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            A-
          </button>
          <button
            type="button"
            title="Размер: Стандартный"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatSize('16px'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-xs font-bold rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            A
          </button>
          <button
            type="button"
            title="Размер: Крупный"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatSize('22px'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-xs font-black rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400"
          >
            A+
          </button>
        </div>

        {/* Text Color Palette */}
        <div className="relative">
          <button
            type="button"
            title="Цвет текста"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); setShowColorPicker(!showColorPicker); }}
            className="edtech-inline-btn p-1.5 bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400"
          >
            <Palette size={14} />
          </button>
          {showColorPicker && (
            <div
              className="absolute top-full left-0 mt-1 p-2 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl shadow-xl flex items-center gap-1.5 z-50"
              onPointerDownCapture={(e) => e.stopPropagation()}
              onMouseDownCapture={(e) => e.stopPropagation()}
            >
              {[
                { label: 'Стандартный', color: 'inherit', bg: 'bg-slate-800 dark:bg-slate-200' },
                { label: 'Индиго', color: '#4f46e5', bg: 'bg-indigo-600' },
                { label: 'Изумруд', color: '#059669', bg: 'bg-emerald-600' },
                { label: 'Янтарь', color: '#d97706', bg: 'bg-amber-600' },
                { label: 'Красный', color: '#e11d48', bg: 'bg-rose-600' },
                { label: 'Голубой', color: '#0284c7', bg: 'bg-sky-600' },
                { label: 'Фиолетовый', color: '#9333ea', bg: 'bg-purple-600' },
              ].map((c) => (
                <button
                  key={c.color}
                  type="button"
                  title={c.label}
                  {...preventBtnFocus}
                  onClick={(e) => {
                    e.stopPropagation();
                    formatColor(c.color);
                    setShowColorPicker(false);
                  }}
                  className={`w-5 h-5 rounded-full ${c.bg} border-2 border-white dark:border-slate-900 shadow-xs hover:scale-110 transition-transform`}
                />
              ))}
            </div>
          )}
        </div>

        {/* Basic formatting: Bold, Italic, Underline, Strikethrough, Code, Highlighter */}
        <div className="flex items-center bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-0.5" data-puck-overlay-portal="true">
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Жирный (Ctrl+B)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('bold'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 font-extrabold"
          >
            <Bold size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Курсив (Ctrl+I)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('italic'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 italic"
          >
            <Italic size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Подчеркнутый (Ctrl+U)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('underline'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Underline size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Зачеркнутый"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('strikeThrough'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Strikethrough size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Моноширинный код"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatCode(); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 font-mono"
          >
            <Code size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Выделитель маркером"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatHighlight('#fef08a'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-yellow-100 dark:hover:bg-yellow-900/40 text-yellow-600"
          >
            <Highlighter size={14} />
          </button>
        </div>

        {/* Alignment */}
        <div className="flex items-center bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-0.5" data-puck-overlay-portal="true">
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="По левому краю"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyLeft'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignLeft size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="По центру"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyCenter'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignCenter size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="По правому краю"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyRight'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignRight size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="По ширине"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyFull'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignJustify size={14} />
          </button>
        </div>

        {/* Lists & Indents */}
        <div className="flex items-center bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-0.5" data-puck-overlay-portal="true">
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Маркированный список"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('insertUnorderedList'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <List size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Нумерованный список"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('insertOrderedList'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <ListOrdered size={14} />
          </button>
          <button
            type="button"
            title="Увеличить отступ"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('indent'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Indent size={14} />
          </button>
          <button
            type="button"
            title="Уменьшить отступ"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('outdent'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Outdent size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Цитата"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('formatBlock', '<blockquote>'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Quote size={14} />
          </button>
        </div>

        {/* Insert Media / Elements */}
        <div className="flex items-center gap-1" data-puck-overlay-portal="true">
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Вставить картинку"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); setShowImageModal(true); }}
            className="edtech-inline-btn px-2.5 py-1.5 rounded-xl bg-indigo-50 dark:bg-indigo-950/60 text-indigo-700 dark:text-indigo-300 hover:bg-indigo-100 dark:hover:bg-indigo-900 text-xs font-bold flex items-center gap-1.5 transition-colors border border-indigo-200 dark:border-indigo-800"
          >
            <ImageIcon size={14} />
            <span>Картинка</span>
          </button>

          {/* Callout Dropdown (💡 Совет, ⚠️ Важно, 📝 Заметка, 🎯 Цель) */}
          <div className="relative">
            <button
              type="button"
              data-puck-overlay-portal="true"
              title="Вставить цветную врезку (Callout)"
              {...preventBtnFocus}
              onClick={(e) => { e.stopPropagation(); setShowCalloutMenu(!showCalloutMenu); }}
              className="edtech-inline-btn px-2.5 py-1.5 rounded-xl bg-amber-50 dark:bg-amber-950/60 text-amber-700 dark:text-amber-300 hover:bg-amber-100 dark:hover:bg-amber-900 text-xs font-bold flex items-center gap-1.5 transition-colors border border-amber-200 dark:border-amber-800"
            >
              <Info size={14} />
              <span>Врезка</span>
            </button>
            {showCalloutMenu && (
              <div
                className="absolute top-full left-0 mt-1 p-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl shadow-2xl flex flex-col gap-1 z-50 min-w-[200px]"
                onPointerDownCapture={(e) => e.stopPropagation()}
                onMouseDownCapture={(e) => e.stopPropagation()}
              >
                <button
                  type="button"
                  {...preventBtnFocus}
                  onClick={(e) => { e.stopPropagation(); handleInsertCallout('tip'); }}
                  className="flex items-center gap-2 px-2.5 py-1.5 text-xs rounded-lg hover:bg-indigo-50 dark:hover:bg-indigo-950/60 text-slate-700 dark:text-slate-200 text-left"
                >
                  <span className="text-base">💡</span>
                  <div>
                    <div className="font-bold text-indigo-600 dark:text-indigo-400">Совет</div>
                    <div className="text-[10px] text-slate-400">Полезная рекомендация</div>
                  </div>
                </button>
                <button
                  type="button"
                  {...preventBtnFocus}
                  onClick={(e) => { e.stopPropagation(); handleInsertCallout('warning'); }}
                  className="flex items-center gap-2 px-2.5 py-1.5 text-xs rounded-lg hover:bg-amber-50 dark:hover:bg-amber-950/60 text-slate-700 dark:text-slate-200 text-left"
                >
                  <span className="text-base">⚠️</span>
                  <div>
                    <div className="font-bold text-amber-600 dark:text-amber-400">Важно</div>
                    <div className="text-[10px] text-slate-400">Внимание к деталям</div>
                  </div>
                </button>
                <button
                  type="button"
                  {...preventBtnFocus}
                  onClick={(e) => { e.stopPropagation(); handleInsertCallout('note'); }}
                  className="flex items-center gap-2 px-2.5 py-1.5 text-xs rounded-lg hover:bg-emerald-50 dark:hover:bg-emerald-950/60 text-slate-700 dark:text-slate-200 text-left"
                >
                  <span className="text-base">📝</span>
                  <div>
                    <div className="font-bold text-emerald-600 dark:text-emerald-400">Заметка</div>
                    <div className="text-[10px] text-slate-400">Информационный блок</div>
                  </div>
                </button>
                <button
                  type="button"
                  {...preventBtnFocus}
                  onClick={(e) => { e.stopPropagation(); handleInsertCallout('target'); }}
                  className="flex items-center gap-2 px-2.5 py-1.5 text-xs rounded-lg hover:bg-purple-50 dark:hover:bg-purple-950/60 text-slate-700 dark:text-slate-200 text-left"
                >
                  <span className="text-base">🎯</span>
                  <div>
                    <div className="font-bold text-purple-600 dark:text-purple-400">Цель</div>
                    <div className="text-[10px] text-slate-400">Ожидаемый результат</div>
                  </div>
                </button>
              </div>
            )}
          </div>

          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Вставить ссылку"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); setShowLinkModal(true); }}
            className="edtech-inline-btn p-1.5 rounded-xl hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300"
          >
            <LinkIcon size={14} />
          </button>

          {/* Import Word / Markdown Button */}
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Импортировать документ Word (.docx) или Markdown (.md)"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); fileInputRef.current?.click(); }}
            className="edtech-inline-btn px-2.5 py-1.5 rounded-xl bg-purple-50 dark:bg-purple-950/60 text-purple-700 dark:text-purple-300 hover:bg-purple-100 dark:hover:bg-purple-900 text-xs font-bold flex items-center gap-1.5 transition-colors border border-purple-200 dark:border-purple-800 cursor-pointer"
          >
            <FileUp size={14} />
            <span>Импорт .docx / .md</span>
          </button>
          <input
            ref={fileInputRef}
            type="file"
            accept=".docx,.md,.markdown"
            className="hidden"
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (file) handleFileImport(file);
              e.target.value = '';
            }}
          />

          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Горизонтальная черта"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); exec('insertHorizontalRule'); }}
            className="edtech-inline-btn p-1.5 rounded-xl hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300"
          >
            <Minus size={14} />
          </button>

          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Очистить форматирование"
            {...preventBtnFocus}
            onClick={(e) => { e.stopPropagation(); clearFormatting(); }}
            className="edtech-inline-btn p-1.5 rounded-xl hover:bg-rose-100 dark:hover:bg-rose-950/50 text-slate-400 hover:text-rose-500"
          >
            <RotateCcw size={14} />
          </button>
        </div>
      </div>

      {/* Popover Modal for Insert Image (Upload & URL Tabs) */}
      {showImageModal && (
        <div className="p-4 bg-indigo-50/90 dark:bg-slate-950 border-b border-indigo-200 dark:border-indigo-900 flex flex-col gap-3 animate-in fade-in">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <span className="text-xs font-black uppercase text-indigo-900 dark:text-indigo-300 flex items-center gap-1.5">
                <ImageIcon size={14} /> Вставка изображения
              </span>
              {/* Tab Selector */}
              <div className="flex items-center bg-white dark:bg-slate-900 rounded-lg p-0.5 border border-indigo-200 dark:border-indigo-800 text-[11px]">
                <button
                  type="button"
                  onClick={() => setImageTab('upload')}
                  className={`px-2 py-0.5 rounded font-bold ${
                    imageTab === 'upload'
                      ? 'bg-indigo-600 text-white shadow-2xs'
                      : 'text-slate-600 dark:text-slate-300 hover:text-indigo-600'
                  }`}
                >
                  📁 С устройства
                </button>
                <button
                  type="button"
                  onClick={() => setImageTab('url')}
                  className={`px-2 py-0.5 rounded font-bold ${
                    imageTab === 'url'
                      ? 'bg-indigo-600 text-white shadow-2xs'
                      : 'text-slate-600 dark:text-slate-300 hover:text-indigo-600'
                  }`}
                >
                  🔗 По ссылке
                </button>
              </div>
            </div>
            <button
              onClick={() => {
                setShowImageModal(false);
                setImageFile(null);
                setImagePreview(null);
              }}
              className="p-1 rounded-lg hover:bg-indigo-200 dark:hover:bg-slate-800 text-slate-500"
            >
              <X size={14} />
            </button>
          </div>

          {imageTab === 'upload' ? (
            <div className="space-y-2">
              <input
                ref={imageFileInputRef}
                type="file"
                accept="image/*"
                className="hidden"
                onChange={(e) => {
                  const file = e.target.files?.[0];
                  if (file) {
                    setImageFile(file);
                    const reader = new FileReader();
                    reader.onload = () => setImagePreview(reader.result as string);
                    reader.readAsDataURL(file);
                  }
                }}
              />
              {imagePreview ? (
                <div className="flex items-center gap-3 p-2 bg-white dark:bg-slate-900 rounded-xl border border-indigo-200 dark:border-indigo-800">
                  <img src={imagePreview} alt="Preview" className="w-16 h-16 object-cover rounded-lg border" />
                  <div className="flex-1 text-xs">
                    <p className="font-bold text-slate-800 dark:text-slate-200 truncate">{imageFile?.name}</p>
                    <p className="text-[10px] text-slate-400">
                      {imageFile?.size ? `${(imageFile.size / 1024).toFixed(1)} KB` : ''}
                    </p>
                    <button
                      type="button"
                      onClick={() => imageFileInputRef.current?.click()}
                      className="mt-1 text-[11px] text-indigo-600 hover:underline font-semibold"
                    >
                      Выбрать другой файл
                    </button>
                  </div>
                </div>
              ) : (
                <div
                  onClick={() => imageFileInputRef.current?.click()}
                  className="p-4 rounded-xl border-2 border-dashed border-indigo-300 dark:border-indigo-800 hover:border-indigo-500 bg-white/60 dark:bg-slate-900/60 text-center cursor-pointer transition-colors"
                >
                  <Upload size={20} className="mx-auto text-indigo-600 dark:text-indigo-400 mb-1" />
                  <p className="text-xs font-bold text-slate-700 dark:text-slate-300">
                    Нажмите, чтобы выбрать файл изображения
                  </p>
                  <p className="text-[10px] text-slate-400">PNG, JPG, WebP, GIF или вставьте через Ctrl+V прямо в текст</p>
                </div>
              )}
            </div>
          ) : (
            <input
              type="text"
              value={imageUrl}
              onChange={(e) => setImageUrl(e.target.value)}
              placeholder="Вставьте URL изображения (https://...)..."
              className="px-3 py-1.5 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-800 dark:text-slate-100 w-full"
            />
          )}

          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <input
              type="text"
              value={imageCaption}
              onChange={(e) => setImageCaption(e.target.value)}
              placeholder="Подпись к картинке (опционально)"
              className="px-3 py-1.5 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-800 dark:text-slate-100 w-full md:col-span-2"
            />
            <div className="flex items-center gap-2">
              <select
                value={imageWidth}
                onChange={(e) => setImageWidth(e.target.value as any)}
                className="px-3 py-1.5 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-800 dark:text-slate-100"
              >
                <option value="100%">100% ширины</option>
                <option value="75%">75% ширины</option>
                <option value="50%">50% ширины</option>
              </select>
              <button
                type="button"
                onClick={handleInsertImage}
                disabled={(!imageFile && !imageUrl.trim()) || isUploadingImage}
                className="px-4 py-1.5 rounded-xl bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold flex items-center gap-1.5 disabled:opacity-40"
              >
                {isUploadingImage ? <Loader2 size={13} className="animate-spin" /> : <Check size={14} />}
                Вставить
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Popover Modal for Insert Link */}
      {showLinkModal && (
        <div className="p-3 bg-indigo-50/90 dark:bg-slate-950 border-b border-indigo-200 dark:border-indigo-900 flex items-center gap-2">
          <input
            type="text"
            value={linkUrl}
            onChange={(e) => setLinkUrl(e.target.value)}
            placeholder="Введите URL ссылки (напр. https://example.com)..."
            className="px-3 py-1.5 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-800 dark:text-slate-100 flex-1"
          />
          <button
            type="button"
            onClick={handleInsertLink}
            disabled={!linkUrl.trim()}
            className="px-4 py-1.5 rounded-xl bg-indigo-600 text-white text-xs font-bold disabled:opacity-40"
          >
            Применить
          </button>
          <button
            onClick={() => setShowLinkModal(false)}
            className="p-1.5 rounded-xl hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-500"
          >
            <X size={14} />
          </button>
        </div>
      )}

      {/* Relative container for Document Sheet with Drag-and-Drop and Overlays */}
      <div
        className="relative"
        onDragOver={(e) => {
          e.preventDefault();
          e.stopPropagation();
          setIsDragOver(true);
        }}
        onDragLeave={(e) => {
          e.preventDefault();
          e.stopPropagation();
          setIsDragOver(false);
        }}
        onDrop={(e) => {
          e.preventDefault();
          e.stopPropagation();
          setIsDragOver(false);
          const file = e.dataTransfer.files?.[0];
          if (file) handleFileImport(file);
        }}
      >
        {/* Drag & Drop Overlay */}
        {isDragOver && (
          <div className="absolute inset-0 z-30 bg-purple-500/10 backdrop-blur-xs border-4 border-dashed border-purple-500 rounded-2xl flex flex-col items-center justify-center p-6 text-center animate-in fade-in">
            <div className="w-14 h-14 rounded-2xl bg-purple-600 text-white flex items-center justify-center mb-3 shadow-lg">
              <FileUp size={28} />
            </div>
            <h4 className="text-base font-extrabold text-purple-900 dark:text-purple-200">
              Отпустите файл Word (.docx) или Markdown (.md)
            </h4>
            <p className="text-xs text-purple-700 dark:text-purple-300 mt-1 max-w-sm">
              Текст, таблицы и иллюстрации будут автоматически сконвертированы и добавлены в этот блок лекции.
            </p>
          </div>
        )}

        {/* Importing Progress Overlay */}
        {isImporting && (
          <div className="absolute inset-0 z-40 bg-white/80 dark:bg-slate-900/80 backdrop-blur-xs flex flex-col items-center justify-center p-6 text-center animate-in fade-in">
            <div className="w-12 h-12 rounded-2xl bg-purple-100 dark:bg-purple-950/80 text-purple-600 flex items-center justify-center mb-3 shadow-md">
              <Loader2 size={24} className="animate-spin" />
            </div>
            <h4 className="text-sm font-extrabold text-slate-900 dark:text-white">
              Конвертация документа...
            </h4>
            <p className="text-xs text-purple-600 dark:text-purple-400 font-semibold mt-1">
              {importStatus || 'Пожалуйста, подождите...'}
            </p>
          </div>
        )}

        {/* Word Editable Document Sheet */}
        <div
          ref={editorRef}
          contentEditable
          suppressContentEditableWarning
          dangerouslySetInnerHTML={{ __html: initialHtmlRef.current }}
          data-puck-overlay-portal="true"
          onPointerDownCapture={(e) => e.stopPropagation()}
          onMouseDownCapture={(e) => e.stopPropagation()}
          onKeyDownCapture={(e) => e.stopPropagation()}
          onKeyUpCapture={(e) => e.stopPropagation()}
          onPointerDown={(e) => e.stopPropagation()}
          onMouseDown={(e) => e.stopPropagation()}
          onMouseUp={saveSelection}
          onKeyUp={saveSelection}
          onSelect={saveSelection}
          onFocus={() => {
            onFocusBlock?.();
            saveSelection();
          }}
          onInput={triggerChange}
          onBlur={triggerChange}
          style={{ userSelect: 'text', WebkitUserSelect: 'text', pointerEvents: 'auto' }}
          className="edtech-inline-editable min-h-[300px] p-8 text-slate-800 dark:text-slate-100 leading-relaxed text-sm md:text-base outline-none focus:ring-2 focus:ring-indigo-500/40 prose dark:prose-invert max-w-none cursor-text select-text"
          data-placeholder="Начните писать лекцию прямо здесь или перетащите файл Word (.docx) / Markdown (.md)..."
        />

        {/* Import Success Toast */}
        {importToast && (
          <div className="absolute bottom-4 right-4 z-50 px-4 py-2.5 rounded-xl bg-purple-600 text-white text-xs font-bold shadow-xl flex items-center gap-2 animate-in slide-in-from-bottom-2">
            <Check size={16} />
            <span>{importToast}</span>
          </div>
        )}
      </div>
    </div>
  );
}
