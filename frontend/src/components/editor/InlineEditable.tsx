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
  Palette
} from 'lucide-react';

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

/**
 * Unified RichTextCanvasEditor:
 * In-place Word/Google Docs style text editing on the canvas with:
 * - Paragraphs (P), Headings (H1, H2, H3), Blockquote
 * - Alignment (Left, Center, Right, Justify)
 * - Font styling (Bold, Italic, Underline, Strikethrough, Code, Highlighter)
 * - Lists (Bullet, Numbered) and Indentation (Indent, Outdent)
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
  const lastHtmlRef = useRef<string>(htmlContent || '');
  const savedSelectionRef = useRef<Range | null>(null);

  const [showImageModal, setShowImageModal] = useState(false);
  const [imageUrl, setImageUrl] = useState('');
  const [imageCaption, setImageCaption] = useState('');
  const [showLinkModal, setShowLinkModal] = useState(false);
  const [linkUrl, setLinkUrl] = useState('');
  const [showColorPicker, setShowColorPicker] = useState(false);

  // Save current selection range inside the editor
  const saveSelection = useCallback(() => {
    if (typeof window === 'undefined') return;
    const sel = window.getSelection();
    if (sel && sel.rangeCount > 0) {
      const range = sel.getRangeAt(0);
      if (editorRef.current && editorRef.current.contains(range.commonAncestorContainer)) {
        savedSelectionRef.current = range.cloneRange();
      }
    }
  }, []);

  // Restore saved selection before applying commands
  const restoreSelection = useCallback(() => {
    if (typeof window === 'undefined' || !editorRef.current) return;
    const sel = window.getSelection();
    if (editorRef.current && document.activeElement !== editorRef.current) {
      editorRef.current.focus();
    }
    if (savedSelectionRef.current && sel) {
      try {
        sel.removeAllRanges();
        sel.addRange(savedSelectionRef.current);
      } catch (e) {}
    }
  }, []);

  // Track selection changes across the document
  useEffect(() => {
    const handleSelectionChange = () => {
      saveSelection();
    };
    document.addEventListener('selectionchange', handleSelectionChange);
    return () => {
      document.removeEventListener('selectionchange', handleSelectionChange);
    };
  }, [saveSelection]);

  // Initialize once and set clean default paragraph separator
  useEffect(() => {
    if (editorRef.current && !isInitializedRef.current) {
      editorRef.current.innerHTML = htmlContent || '';
      lastHtmlRef.current = htmlContent || '';
      isInitializedRef.current = true;
    }
    try {
      document.execCommand('defaultParagraphSeparator', false, 'p');
      document.execCommand('styleWithCSS', false, 'true');
    } catch (e) {}
  }, []);

  // Sync external changes ONLY when prop actually changed from outside
  useEffect(() => {
    if (editorRef.current && isInitializedRef.current) {
      if (htmlContent !== undefined && htmlContent !== lastHtmlRef.current) {
        lastHtmlRef.current = htmlContent;
        editorRef.current.innerHTML = htmlContent || '';
      }
    }
  }, [htmlContent]);

  const triggerChange = useCallback(() => {
    if (!editorRef.current) return;
    const newHtml = editorRef.current.innerHTML;
    lastHtmlRef.current = newHtml;
    onChange(newHtml);
  }, [onChange]);

  const exec = (command: string, value: string | undefined = undefined) => {
    if (!editorRef.current) return;
    restoreSelection();
    try {
      document.execCommand('styleWithCSS', false, 'true');
    } catch (e) {}
    document.execCommand(command, false, value);
    saveSelection();
    triggerChange();
  };

  // Native keyboard event isolation to stop Puck from blocking Backspace, Delete, Ctrl+Z
  useEffect(() => {
    const el = editorRef.current;
    if (!el) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      // Isolate keydown from Puck hotkeys
      e.stopPropagation();

      if (e.key === 'Tab') {
        e.preventDefault();
        document.execCommand('insertHTML', false, '&nbsp;&nbsp;&nbsp;&nbsp;');
        triggerChange();
        return;
      }

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

    const stopPropagation = (e: KeyboardEvent) => {
      e.stopPropagation();
    };

    el.addEventListener('keydown', handleKeyDown, false);
    el.addEventListener('keyup', stopPropagation, false);
    el.addEventListener('keypress', stopPropagation, false);

    return () => {
      el.removeEventListener('keydown', handleKeyDown, false);
      el.removeEventListener('keyup', stopPropagation, false);
      el.removeEventListener('keypress', stopPropagation, false);
    };
  }, [triggerChange]);

  const formatHeading = (tag: 'p' | 'h1' | 'h2' | 'h3') => {
    if (!editorRef.current) return;
    restoreSelection();
    try {
      document.execCommand('styleWithCSS', false, 'true');
    } catch (e) {}
    let ok = document.execCommand('formatBlock', false, `<${tag}>`);
    if (!ok) {
      document.execCommand('formatBlock', false, tag);
    }
    saveSelection();
    triggerChange();
  };

  const formatFont = (fontFamily: string) => {
    if (!editorRef.current) return;
    restoreSelection();
    const sel = window.getSelection();
    if (sel && !sel.isCollapsed) {
      try {
        document.execCommand('styleWithCSS', false, 'true');
      } catch (e) {}
      document.execCommand('fontName', false, fontFamily);
    } else if (sel && sel.anchorNode && editorRef.current) {
      let el: HTMLElement | null =
        (sel.anchorNode.nodeType === Node.ELEMENT_NODE
          ? sel.anchorNode
          : sel.anchorNode.parentElement) as HTMLElement;
      while (el && el !== editorRef.current && el.parentElement !== editorRef.current) {
        el = el.parentElement;
      }
      if (el && el !== editorRef.current) {
        el.style.fontFamily = fontFamily;
      } else {
        editorRef.current.style.fontFamily = fontFamily;
      }
    }
    saveSelection();
    triggerChange();
  };

  const formatSize = (sizeCmd: string, sizePx: string) => {
    if (!editorRef.current) return;
    restoreSelection();
    const sel = window.getSelection();
    if (sel && !sel.isCollapsed) {
      try {
        document.execCommand('styleWithCSS', false, 'true');
      } catch (e) {}
      document.execCommand('fontSize', false, sizeCmd);
    } else if (sel && sel.anchorNode && editorRef.current) {
      let el: HTMLElement | null =
        (sel.anchorNode.nodeType === Node.ELEMENT_NODE
          ? sel.anchorNode
          : sel.anchorNode.parentElement) as HTMLElement;
      while (el && el !== editorRef.current && el.parentElement !== editorRef.current) {
        el = el.parentElement;
      }
      if (el && el !== editorRef.current) {
        el.style.fontSize = sizePx;
      }
    }
    saveSelection();
    triggerChange();
  };

  const formatColor = (color: string) => {
    if (!editorRef.current) return;
    restoreSelection();
    const sel = window.getSelection();
    if (sel && !sel.isCollapsed) {
      try {
        document.execCommand('styleWithCSS', false, 'true');
      } catch (e) {}
      document.execCommand('foreColor', false, color);
    } else if (sel && sel.anchorNode && editorRef.current) {
      let el: HTMLElement | null =
        (sel.anchorNode.nodeType === Node.ELEMENT_NODE
          ? sel.anchorNode
          : sel.anchorNode.parentElement) as HTMLElement;
      while (el && el !== editorRef.current && el.parentElement !== editorRef.current) {
        el = el.parentElement;
      }
      if (el && el !== editorRef.current) {
        el.style.color = color;
      }
    }
    saveSelection();
    triggerChange();
  };

  const formatAlign = (alignCmd: 'justifyLeft' | 'justifyCenter' | 'justifyRight' | 'justifyFull') => {
    if (!editorRef.current) return;
    restoreSelection();
    try {
      document.execCommand('styleWithCSS', false, 'true');
    } catch (e) {}
    document.execCommand(alignCmd, false);
    const sel = window.getSelection();
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
    if (!editorRef.current) return;
    restoreSelection();
    try {
      document.execCommand('styleWithCSS', false, 'true');
    } catch (e) {}
    let ok = document.execCommand('hiliteColor', false, color);
    if (!ok) {
      document.execCommand('backColor', false, color);
    }
    saveSelection();
    triggerChange();
  };

  const insertCustomHtml = (htmlSnippet: string) => {
    if (!editorRef.current) return;
    restoreSelection();
    document.execCommand('insertHTML', false, htmlSnippet);
    saveSelection();
    triggerChange();
  };

  const handleInsertImage = () => {
    if (!imageUrl.trim()) return;
    const snippet = `
      <figure class="my-4 mx-auto text-center max-w-full">
        <img src="${imageUrl}" alt="${imageCaption || 'Иллюстрация'}" class="rounded-2xl shadow-md border border-slate-200 dark:border-slate-800 mx-auto object-cover max-h-[420px]" />
        ${imageCaption ? `<figcaption class="mt-1.5 text-xs text-slate-500 dark:text-slate-400 font-medium">${imageCaption}</figcaption>` : ''}
      </figure>
      <p><br/></p>
    `;
    insertCustomHtml(snippet);
    setImageUrl('');
    setImageCaption('');
    setShowImageModal(false);
  };

  const handleInsertCallout = () => {
    const snippet = `
      <div class="my-4 p-4 rounded-2xl bg-indigo-50/80 dark:bg-indigo-950/40 border-l-4 border-indigo-500 flex items-start gap-3">
        <span class="text-xl">💡</span>
        <div class="text-sm leading-relaxed text-slate-800 dark:text-slate-200 flex-1">
          <strong>Важное замечание:</strong> введите сюда ключевой тезис или вывод к теме...
        </div>
      </div>
      <p><br/></p>
    `;
    insertCustomHtml(snippet);
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

  return (
    <div
      className="my-3 rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm overflow-hidden transition-all"
      onClick={() => onFocusBlock?.()}
    >
      {/* Sleek Canvas Formatting Ribbon */}
      <div
        data-puck-overlay-portal="true"
        onPointerDown={(e) => e.stopPropagation()}
        onMouseDown={(e) => {
          e.stopPropagation();
          const target = e.target as HTMLElement;
          if (target && (target.tagName === 'INPUT' || target.tagName === 'SELECT' || target.tagName === 'TEXTAREA')) {
            return;
          }
          e.preventDefault();
        }}
        className="bg-slate-100/90 dark:bg-slate-900/90 px-3 py-1.5 border-b border-slate-200 dark:border-slate-800 flex flex-wrap items-center gap-1.5 text-slate-700 dark:text-slate-300 select-none text-xs"
      >
        {/* Drag Handle */}
        {dragRef && (
          <div
            ref={dragRef}
            data-puck-drag-handle="true"
            className="flex items-center gap-1.5 px-2 py-1 rounded-lg bg-indigo-600 hover:bg-indigo-700 text-white text-[11px] font-bold transition-colors cursor-grab active:cursor-grabbing select-none"
            title="Зажмите и потяните, чтобы переместить блок"
          >
            <GripVertical size={13} />
            <span>Текст</span>
          </div>
        )}

        {/* Undo / Redo */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Отменить действие (Ctrl+Z)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('undo'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300"
          >
            <Undo size={14} />
          </button>
          <button
            type="button"
            title="Повторить действие (Ctrl+Y)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('redo'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300"
          >
            <Redo size={14} />
          </button>
        </div>

        {/* Paragraph & Headings */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Обычный абзац (P)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatHeading('p'); }}
            className="edtech-inline-btn px-2 py-1 text-xs font-bold rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            P
          </button>
          <button
            type="button"
            title="Заголовок H1"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatHeading('h1'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400 font-extrabold"
          >
            <Heading1 size={14} />
          </button>
          <button
            type="button"
            title="Заголовок H2"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatHeading('h2'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400 font-bold"
          >
            <Heading2 size={14} />
          </button>
          <button
            type="button"
            title="Подзаголовок H3"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatHeading('h3'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400"
          >
            <Heading3 size={14} />
          </button>
        </div>

        {/* Font Family (Шрифт) */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Шрифт: Стандартный без засечек"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatFont('sans-serif'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[11px] font-medium rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            Sans
          </button>
          <button
            type="button"
            title="Шрифт: С засечками (Serif / Книга)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatFont('Georgia, serif'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[11px] font-serif rounded hover:bg-slate-100 dark:hover:bg-slate-800 italic"
          >
            Serif
          </button>
          <button
            type="button"
            title="Шрифт: Моноширинный (Код)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
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
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatSize('2', '13px'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[10px] font-bold rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            A-
          </button>
          <button
            type="button"
            title="Размер шрифта: Стандартный"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatSize('3', '16px'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-xs font-bold rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            A
          </button>
          <button
            type="button"
            title="Размер шрифта: Крупный"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatSize('5', '22px'); }}
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
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
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
                  onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
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
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyLeft'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignLeft size={14} />
          </button>
          <button
            type="button"
            title="По центру"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyCenter'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignCenter size={14} />
          </button>
          <button
            type="button"
            title="По правому краю"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyRight'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignRight size={14} />
          </button>
          <button
            type="button"
            title="По ширине"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyFull'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignJustify size={14} />
          </button>
        </div>

        {/* Font styling: Bold, Italic, Underline, Strikethrough, Code */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Жирный (Ctrl+B)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('bold'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 font-extrabold"
          >
            <Bold size={14} />
          </button>
          <button
            type="button"
            title="Курсив (Ctrl+I)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('italic'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 italic"
          >
            <Italic size={14} />
          </button>
          <button
            type="button"
            title="Подчеркнутый (Ctrl+U)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('underline'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Underline size={14} />
          </button>
          <button
            type="button"
            title="Зачеркнутый"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('strikeThrough'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Strikethrough size={14} />
          </button>
          <button
            type="button"
            title="Код"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatHeading('p'); exec('formatBlock', '<pre>'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 font-mono"
          >
            <Code size={14} />
          </button>
        </div>

        {/* Lists & Indentation */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Маркированный список"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('insertUnorderedList'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <List size={14} />
          </button>
          <button
            type="button"
            title="Нумерованный список"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('insertOrderedList'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <ListOrdered size={14} />
          </button>
          <button
            type="button"
            title="Увеличить отступ"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('indent'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Indent size={14} />
          </button>
          <button
            type="button"
            title="Уменьшить отступ"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('outdent'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Outdent size={14} />
          </button>
        </div>

        {/* Rich inserts: Quote, Marker, Callout, Image, Link, Divider, Clear */}
        <div className="flex items-center bg-white dark:bg-slate-950 rounded-lg border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Цитата"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('formatBlock', '<blockquote>'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-amber-600"
          >
            <Quote size={14} />
          </button>
          <button
            type="button"
            title="Выделитель желтым"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatHighlight('#fef08a'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-yellow-100 dark:hover:bg-yellow-950/60 text-yellow-600"
          >
            <Highlighter size={14} />
          </button>
          <button
            type="button"
            title="Вставить важную заметку (Callout)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); handleInsertCallout(); }}
            className="edtech-inline-btn p-1 rounded hover:bg-indigo-50 dark:hover:bg-indigo-950 text-indigo-600"
          >
            <Info size={14} />
          </button>
          <button
            type="button"
            title="Вставить картинку по ссылке"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); setShowImageModal(!showImageModal); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-sky-600"
          >
            <ImageIcon size={14} />
          </button>
          <button
            type="button"
            title="Вставить ссылку"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); setShowLinkModal(!showLinkModal); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-blue-600"
          >
            <LinkIcon size={14} />
          </button>
          <button
            type="button"
            title="Разделительная линия"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('insertHorizontalRule'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-500"
          >
            <Minus size={14} />
          </button>
          <button
            type="button"
            title="Очистить форматирование"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('removeFormat'); }}
            className="edtech-inline-btn p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-rose-500"
          >
            <RemoveFormatting size={14} />
          </button>
        </div>
      </div>

      {/* Popover Modal for Insert Image */}
      {showImageModal && (
        <div
          data-puck-overlay-portal="true"
          onPointerDownCapture={(e) => e.stopPropagation()}
          onMouseDownCapture={(e) => e.stopPropagation()}
          className="p-3 bg-indigo-50/90 dark:bg-slate-950 border-b border-indigo-200 dark:border-indigo-900 space-y-2"
        >
          <div className="flex items-center justify-between">
            <span className="text-xs font-bold text-indigo-900 dark:text-indigo-200">Вставка изображения</span>
            <button
              onClick={() => setShowImageModal(false)}
              className="p-1 rounded-lg hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-500"
            >
              <X size={14} />
            </button>
          </div>
          <div className="flex flex-col gap-2">
            <input
              type="text"
              value={imageUrl}
              onChange={(e) => setImageUrl(e.target.value)}
              placeholder="Вставьте URL изображения (https://...)..."
              className="px-3 py-1.5 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-800 dark:text-slate-100 w-full"
            />
            <input
              type="text"
              value={imageCaption}
              onChange={(e) => setImageCaption(e.target.value)}
              placeholder="Подпись к картинке (опционально)..."
              className="px-3 py-1.5 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-800 dark:text-slate-100 w-full"
            />
            <div className="flex justify-end">
              <button
                type="button"
                onClick={handleInsertImage}
                disabled={!imageUrl.trim()}
                className="px-4 py-1.5 rounded-xl bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold flex items-center gap-1.5 disabled:opacity-40"
              >
                <Check size={14} /> Вставить
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

      {/* Editable Body */}
      <div
        ref={editorRef}
        contentEditable
        suppressContentEditableWarning
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
  const lastHtmlRef = useRef<string>(htmlContent || '');
  const savedSelectionRef = useRef<Range | null>(null);

  const [showImageModal, setShowImageModal] = useState(false);
  const [imageUrl, setImageUrl] = useState('');
  const [imageCaption, setImageCaption] = useState('');
  const [imageWidth, setImageWidth] = useState<'100%' | '75%' | '50%'>('100%');
  const [showLinkModal, setShowLinkModal] = useState(false);
  const [linkUrl, setLinkUrl] = useState('');
  const [showColorPicker, setShowColorPicker] = useState(false);

  // Save current selection range inside the editor
  const saveSelection = useCallback(() => {
    if (typeof window === 'undefined') return;
    const sel = window.getSelection();
    if (sel && sel.rangeCount > 0) {
      const range = sel.getRangeAt(0);
      if (editorRef.current && editorRef.current.contains(range.commonAncestorContainer)) {
        savedSelectionRef.current = range.cloneRange();
      }
    }
  }, []);

  // Restore saved selection before applying commands
  const restoreSelection = useCallback(() => {
    if (typeof window === 'undefined' || !editorRef.current) return;
    const sel = window.getSelection();
    if (editorRef.current && document.activeElement !== editorRef.current) {
      editorRef.current.focus();
    }
    if (savedSelectionRef.current && sel) {
      try {
        sel.removeAllRanges();
        sel.addRange(savedSelectionRef.current);
      } catch (e) {}
    }
  }, []);

  // Track selection changes across the document
  useEffect(() => {
    const handleSelectionChange = () => {
      saveSelection();
    };
    document.addEventListener('selectionchange', handleSelectionChange);
    return () => {
      document.removeEventListener('selectionchange', handleSelectionChange);
    };
  }, [saveSelection]);

  // Set initial content on mount once
  useEffect(() => {
    if (editorRef.current && !isInitializedRef.current) {
      editorRef.current.innerHTML = htmlContent || '';
      lastHtmlRef.current = htmlContent || '';
      isInitializedRef.current = true;
    }
    try {
      document.execCommand('defaultParagraphSeparator', false, 'p');
      document.execCommand('styleWithCSS', false, 'true');
    } catch (e) {}
  }, []);

  // Sync external changes ONLY when prop actually changed from outside
  useEffect(() => {
    if (editorRef.current && isInitializedRef.current) {
      if (htmlContent !== undefined && htmlContent !== lastHtmlRef.current) {
        lastHtmlRef.current = htmlContent;
        editorRef.current.innerHTML = htmlContent || '';
      }
    }
  }, [htmlContent]);

  const triggerChange = useCallback(() => {
    if (!editorRef.current) return;
    const newHtml = editorRef.current.innerHTML;
    lastHtmlRef.current = newHtml;
    onChange(newHtml);
  }, [onChange]);

  const exec = (command: string, value: string | undefined = undefined) => {
    if (!editorRef.current) return;
    restoreSelection();
    try {
      document.execCommand('styleWithCSS', false, 'true');
    } catch (e) {}
    document.execCommand(command, false, value);
    saveSelection();
    triggerChange();
  };

  // Native keyboard event isolation to stop Puck from blocking Backspace, Delete, Ctrl+Z
  useEffect(() => {
    const el = editorRef.current;
    if (!el) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      // Isolate keydown from Puck hotkeys
      e.stopPropagation();

      if (e.key === 'Tab') {
        e.preventDefault();
        document.execCommand('insertHTML', false, '&nbsp;&nbsp;&nbsp;&nbsp;');
        triggerChange();
        return;
      }

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

    const stopPropagation = (e: KeyboardEvent) => {
      e.stopPropagation();
    };

    el.addEventListener('keydown', handleKeyDown, false);
    el.addEventListener('keyup', stopPropagation, false);
    el.addEventListener('keypress', stopPropagation, false);

    return () => {
      el.removeEventListener('keydown', handleKeyDown, false);
      el.removeEventListener('keyup', stopPropagation, false);
      el.removeEventListener('keypress', stopPropagation, false);
    };
  }, [triggerChange]);

  const formatHeading = (tag: 'p' | 'h1' | 'h2' | 'h3') => {
    if (!editorRef.current) return;
    restoreSelection();
    try {
      document.execCommand('styleWithCSS', false, 'true');
    } catch (e) {}
    let ok = document.execCommand('formatBlock', false, `<${tag}>`);
    if (!ok) {
      document.execCommand('formatBlock', false, tag);
    }
    saveSelection();
    triggerChange();
  };

  const formatFont = (fontFamily: string) => {
    if (!editorRef.current) return;
    restoreSelection();
    const sel = window.getSelection();
    if (sel && !sel.isCollapsed) {
      try {
        document.execCommand('styleWithCSS', false, 'true');
      } catch (e) {}
      document.execCommand('fontName', false, fontFamily);
    } else if (sel && sel.anchorNode && editorRef.current) {
      let el: HTMLElement | null =
        (sel.anchorNode.nodeType === Node.ELEMENT_NODE
          ? sel.anchorNode
          : sel.anchorNode.parentElement) as HTMLElement;
      while (el && el !== editorRef.current && el.parentElement !== editorRef.current) {
        el = el.parentElement;
      }
      if (el && el !== editorRef.current) {
        el.style.fontFamily = fontFamily;
      } else {
        editorRef.current.style.fontFamily = fontFamily;
      }
    }
    saveSelection();
    triggerChange();
  };

  const formatSize = (sizeCmd: string, sizePx: string) => {
    if (!editorRef.current) return;
    restoreSelection();
    const sel = window.getSelection();
    if (sel && !sel.isCollapsed) {
      try {
        document.execCommand('styleWithCSS', false, 'true');
      } catch (e) {}
      document.execCommand('fontSize', false, sizeCmd);
    } else if (sel && sel.anchorNode && editorRef.current) {
      let el: HTMLElement | null =
        (sel.anchorNode.nodeType === Node.ELEMENT_NODE
          ? sel.anchorNode
          : sel.anchorNode.parentElement) as HTMLElement;
      while (el && el !== editorRef.current && el.parentElement !== editorRef.current) {
        el = el.parentElement;
      }
      if (el && el !== editorRef.current) {
        el.style.fontSize = sizePx;
      }
    }
    saveSelection();
    triggerChange();
  };

  const formatColor = (color: string) => {
    if (!editorRef.current) return;
    restoreSelection();
    const sel = window.getSelection();
    if (sel && !sel.isCollapsed) {
      try {
        document.execCommand('styleWithCSS', false, 'true');
      } catch (e) {}
      document.execCommand('foreColor', false, color);
    } else if (sel && sel.anchorNode && editorRef.current) {
      let el: HTMLElement | null =
        (sel.anchorNode.nodeType === Node.ELEMENT_NODE
          ? sel.anchorNode
          : sel.anchorNode.parentElement) as HTMLElement;
      while (el && el !== editorRef.current && el.parentElement !== editorRef.current) {
        el = el.parentElement;
      }
      if (el && el !== editorRef.current) {
        el.style.color = color;
      }
    }
    saveSelection();
    triggerChange();
  };

  const formatAlign = (alignCmd: 'justifyLeft' | 'justifyCenter' | 'justifyRight' | 'justifyFull') => {
    if (!editorRef.current) return;
    restoreSelection();
    try {
      document.execCommand('styleWithCSS', false, 'true');
    } catch (e) {}
    document.execCommand(alignCmd, false);
    const sel = window.getSelection();
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
    if (!editorRef.current) return;
    restoreSelection();
    try {
      document.execCommand('styleWithCSS', false, 'true');
    } catch (e) {}
    let ok = document.execCommand('hiliteColor', false, color);
    if (!ok) {
      document.execCommand('backColor', false, color);
    }
    saveSelection();
    triggerChange();
  };

  const insertCustomHtml = (htmlSnippet: string) => {
    if (!editorRef.current) return;
    restoreSelection();
    document.execCommand('insertHTML', false, htmlSnippet);
    saveSelection();
    triggerChange();
  };

  const handleInsertImage = () => {
    if (!imageUrl.trim()) return;
    const widthClass =
      imageWidth === '50%' ? 'max-w-[50%]' : imageWidth === '75%' ? 'max-w-[75%]' : 'w-full';

    const snippet = `
      <figure class="my-6 mx-auto ${widthClass} text-center">
        <img src="${imageUrl}" alt="${imageCaption || 'Иллюстрация к лекции'}" class="rounded-2xl shadow-lg border border-slate-200 dark:border-slate-800 mx-auto object-cover max-h-[480px]" />
        ${imageCaption ? `<figcaption class="mt-2 text-xs text-slate-500 dark:text-slate-400 font-medium">${imageCaption}</figcaption>` : ''}
      </figure>
      <p><br/></p>
    `;
    insertCustomHtml(snippet);
    setImageUrl('');
    setImageCaption('');
    setShowImageModal(false);
  };

  const handleInsertCallout = () => {
    const snippet = `
      <div class="my-5 p-4 rounded-2xl bg-indigo-50/80 dark:bg-indigo-950/40 border-l-4 border-indigo-500 shadow-sm flex items-start gap-3">
        <span class="text-xl">💡</span>
        <div class="text-sm leading-relaxed text-slate-800 dark:text-slate-200 flex-1">
          <strong>Важное замечание:</strong> введите сюда ключевой вывод, правило или формулу к лекции...
        </div>
      </div>
      <p><br/></p>
    `;
    insertCustomHtml(snippet);
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

  return (
    <div
      ref={containerRef}
      data-puck-overlay-portal="true"
      onClick={() => onFocusBlock?.()}
      className="edtech-inline-editable my-6 rounded-3xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-lg overflow-hidden transition-all"
    >
      {/* Top Banner & Header */}
      <div className="bg-slate-100/90 dark:bg-slate-800/90 px-5 py-3 border-b border-slate-200 dark:border-slate-700 flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          {/* Drag Handle */}
          <div
            ref={dragRef}
            data-puck-drag-handle="true"
            className="flex items-center gap-1.5 px-2.5 py-1 rounded-xl bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold transition-colors cursor-grab active:cursor-grabbing select-none"
            title="Зажмите и потяните, чтобы переместить блок по уроку"
          >
            <GripVertical size={14} />
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
        onPointerDown={(e) => e.stopPropagation()}
        onMouseDown={(e) => {
          e.stopPropagation();
          const target = e.target as HTMLElement;
          if (target && (target.tagName === 'INPUT' || target.tagName === 'SELECT' || target.tagName === 'TEXTAREA')) {
            return;
          }
          e.preventDefault();
        }}
        className="bg-slate-50 dark:bg-slate-950 px-4 py-2 border-b border-slate-200 dark:border-slate-800 flex flex-wrap items-center gap-1.5 text-slate-700 dark:text-slate-300 select-none text-xs"
      >
        {/* Undo / Redo */}
        <div className="flex items-center bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-0.5">
          <button
            type="button"
            title="Отменить действие (Ctrl+Z)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('undo'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300"
          >
            <Undo size={14} />
          </button>
          <button
            type="button"
            title="Повторить действие (Ctrl+Y)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
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
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatHeading('p'); }}
            className="edtech-inline-btn px-2 py-1 text-xs font-bold rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            P
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Заголовок H1"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatHeading('h1'); }}
            className="edtech-inline-btn p-1 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400 font-extrabold"
          >
            <Heading1 size={15} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Заголовок H2"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatHeading('h2'); }}
            className="edtech-inline-btn p-1 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-indigo-600 dark:text-indigo-400 font-bold"
          >
            <Heading2 size={15} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Подзаголовок H3"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
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
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatFont('sans-serif'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[11px] font-medium rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            Sans
          </button>
          <button
            type="button"
            title="Шрифт: С засечками (Serif / Книга)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatFont('Georgia, serif'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[11px] font-serif rounded hover:bg-slate-100 dark:hover:bg-slate-800 italic"
          >
            Serif
          </button>
          <button
            type="button"
            title="Шрифт: Моноширинный (Mono)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
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
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatSize('2', '13px'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-[10px] font-bold rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            A-
          </button>
          <button
            type="button"
            title="Размер: Стандартный"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatSize('3', '16px'); }}
            className="edtech-inline-btn px-1.5 py-0.5 text-xs font-bold rounded hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            A
          </button>
          <button
            type="button"
            title="Размер: Крупный"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatSize('5', '22px'); }}
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
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
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
                  onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
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
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('bold'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 font-extrabold"
          >
            <Bold size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Курсив (Ctrl+I)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('italic'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 italic"
          >
            <Italic size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Подчеркнутый (Ctrl+U)"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('underline'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Underline size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Зачеркнутый"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('strikeThrough'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Strikethrough size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Моноширинный код"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatHeading('p'); exec('formatBlock', '<pre>'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 font-mono"
          >
            <Code size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Выделитель маркером"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
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
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyLeft'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignLeft size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="По центру"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyCenter'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignCenter size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="По правому краю"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); formatAlign('justifyRight'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <AlignRight size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="По ширине"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
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
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('insertUnorderedList'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <List size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Нумерованный список"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('insertOrderedList'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <ListOrdered size={14} />
          </button>
          <button
            type="button"
            title="Увеличить отступ"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('indent'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Indent size={14} />
          </button>
          <button
            type="button"
            title="Уменьшить отступ"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('outdent'); }}
            className="edtech-inline-btn p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <Outdent size={14} />
          </button>
          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Цитата"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
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
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); setShowImageModal(true); }}
            className="edtech-inline-btn px-2.5 py-1.5 rounded-xl bg-indigo-50 dark:bg-indigo-950/60 text-indigo-700 dark:text-indigo-300 hover:bg-indigo-100 dark:hover:bg-indigo-900 text-xs font-bold flex items-center gap-1.5 transition-colors border border-indigo-200 dark:border-indigo-800"
          >
            <ImageIcon size={14} />
            <span>Картинка</span>
          </button>

          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Вставить врезку / заметку"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); handleInsertCallout(); }}
            className="edtech-inline-btn px-2.5 py-1.5 rounded-xl bg-amber-50 dark:bg-amber-950/60 text-amber-700 dark:text-amber-300 hover:bg-amber-100 dark:hover:bg-amber-900 text-xs font-bold flex items-center gap-1.5 transition-colors border border-amber-200 dark:border-amber-800"
          >
            <Info size={14} />
            <span>Заметка</span>
          </button>

          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Вставить ссылку"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); setShowLinkModal(true); }}
            className="edtech-inline-btn p-1.5 rounded-xl hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300"
          >
            <LinkIcon size={14} />
          </button>

          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Горизонтальная черта"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('insertHorizontalRule'); }}
            className="edtech-inline-btn p-1.5 rounded-xl hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300"
          >
            <Minus size={14} />
          </button>

          <button
            type="button"
            data-puck-overlay-portal="true"
            title="Очистить форматирование"
            onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); }}
            onClick={(e) => { e.stopPropagation(); exec('removeFormat'); }}
            className="edtech-inline-btn p-1.5 rounded-xl hover:bg-rose-100 dark:hover:bg-rose-950/50 text-slate-400 hover:text-rose-500"
          >
            <RotateCcw size={14} />
          </button>
        </div>
      </div>

      {/* Popover Modal for Insert Image */}
      {showImageModal && (
        <div className="p-4 bg-indigo-50/90 dark:bg-slate-950 border-b border-indigo-200 dark:border-indigo-900 flex flex-col gap-3 animate-in fade-in">
          <div className="flex items-center justify-between">
            <span className="text-xs font-black uppercase text-indigo-900 dark:text-indigo-300 flex items-center gap-1.5">
              <ImageIcon size={14} /> Вставка изображения
            </span>
            <button
              onClick={() => setShowImageModal(false)}
              className="p-1 rounded-lg hover:bg-indigo-200 dark:hover:bg-slate-800 text-slate-500"
            >
              <X size={14} />
            </button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <input
              type="text"
              value={imageUrl}
              onChange={(e) => setImageUrl(e.target.value)}
              placeholder="URL изображения (https://...)"
              className="px-3 py-1.5 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-800 dark:text-slate-100 w-full"
            />
            <input
              type="text"
              value={imageCaption}
              onChange={(e) => setImageCaption(e.target.value)}
              placeholder="Подпись к картинке (опционально)"
              className="px-3 py-1.5 text-xs rounded-xl bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-700 text-slate-800 dark:text-slate-100 w-full"
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
                disabled={!imageUrl.trim()}
                className="px-4 py-1.5 rounded-xl bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold flex items-center gap-1.5 disabled:opacity-40"
              >
                <Check size={14} /> Вставить
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

      {/* Word Editable Document Sheet */}
      <div
        ref={editorRef}
        contentEditable
        suppressContentEditableWarning
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
        data-placeholder="Начните писать лекцию прямо здесь..."
      />
    </div>
  );
}
