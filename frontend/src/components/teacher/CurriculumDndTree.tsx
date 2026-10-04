'use client';

import React, { useState, useRef } from 'react';
import { useRouter } from 'next/navigation';
import {
  GripVertical,
  Edit3,
  Sparkles,
  PlusCircle,
  ChevronUp,
  ChevronDown,
} from 'lucide-react';
import { api } from '@/lib/api';

export interface CurriculumDndTreeProps {
  courseId: string | number;
  sections: any[];
  canEdit: boolean;
  onSectionsChange: (newSections: any[]) => void;
  onEditModule: (secWrap: any) => void;
  onEditLesson: (lesson: any) => void;
  onCreateLesson: (target: { sectionId: number; sectionTitle: string; nextIndex: number }) => void;
  showToast: (msg: string, type?: 'success' | 'error') => void;
  fetchCourseData: () => Promise<void>;
}

export default function CurriculumDndTree({
  courseId,
  sections,
  canEdit,
  onSectionsChange,
  onEditModule,
  onEditLesson,
  onCreateLesson,
  showToast,
  fetchCourseData,
}: CurriculumDndTreeProps) {
  const router = useRouter();

  // Drag state for modules (sections)
  const [draggedSectionIndex, setDraggedSectionIndex] = useState<number | null>(null);
  const [dropSectionIndex, setDropSectionIndex] = useState<number | null>(null);
  const draggedSectionRef = useRef<number | null>(null);

  // Drag state for lessons
  const [draggedLesson, setDraggedLesson] = useState<{
    sectionIndex: number;
    lessonIndex: number;
    lessonId: number;
  } | null>(null);
  const [dropLessonTarget, setDropLessonTarget] = useState<{
    sectionIndex: number;
    lessonIndex: number;
  } | null>(null);
  const draggedLessonRef = useRef<{
    sectionIndex: number;
    lessonIndex: number;
    lessonId: number;
  } | null>(null);

  const getSection = (secWrap: any) => secWrap?.Section || secWrap?.section || {};
  const getLessons = (secWrap: any) => secWrap?.Lessons || secWrap?.lessons || [];
  const getSectionId = (secWrap: any) => Number(getSection(secWrap)?.ID || getSection(secWrap)?.id);
  const getLessonId = (l: any) => Number(l?.ID || l?.id);

  // --- MANUAL MOVE (FALLBACK FOR ARROW BUTTONS) ---
  const handleMoveSection = async (fromIdx: number, toIdx: number) => {
    if (!canEdit) {
      showToast('У вас нет прав на редактирование этого курса.', 'error');
      return;
    }
    if (toIdx < 0 || toIdx >= sections.length || fromIdx === toIdx) return;

    const previousSections = [...sections];
    const copy = [...sections];
    const [moved] = copy.splice(fromIdx, 1);
    copy.splice(toIdx, 0, moved);
    onSectionsChange(copy);

    try {
      const secIds = copy.map((s) => getSectionId(s));
      await api.put(`/courses/${courseId}/reorder-sections`, { item_ids: secIds });
      showToast('Порядок модулей обновлен');
    } catch (err: any) {
      console.error('Failed to reorder sections', err);
      onSectionsChange(previousSections);
      showToast('Ошибка при изменении порядка модулей', 'error');
    }
  };

  const handleMoveLessonInList = async (
    sectionId: number,
    fromIdx: number,
    toIdx: number,
    lessonsList: any[]
  ) => {
    if (!canEdit) {
      showToast('У вас нет прав на редактирование этого курса.', 'error');
      return;
    }
    if (toIdx < 0 || toIdx >= lessonsList.length || fromIdx === toIdx) return;

    const previousSections = JSON.parse(JSON.stringify(sections));
    const newSections = sections.map((secWrap) => {
      if (getSectionId(secWrap) === sectionId) {
        const copyLessons = [...getLessons(secWrap)];
        const [moved] = copyLessons.splice(fromIdx, 1);
        copyLessons.splice(toIdx, 0, moved);
        return {
          ...secWrap,
          Lessons: copyLessons,
          lessons: copyLessons,
        };
      }
      return secWrap;
    });

    onSectionsChange(newSections);

    try {
      const targetSec = newSections.find((s) => getSectionId(s) === sectionId);
      const lessonIds = getLessons(targetSec).map((l: any) => getLessonId(l));
      await api.put(`/sections/${sectionId}/reorder-lessons`, { item_ids: lessonIds });
      showToast('Порядок уроков сохранен');
    } catch (err: any) {
      console.error('Failed to reorder lessons', err);
      onSectionsChange(previousSections);
      showToast('Ошибка при изменении порядка уроков', 'error');
    }
  };

  // --- MODULE DRAG AND DROP ---
  const handleSectionDragStart = (e: React.DragEvent, sIdx: number) => {
    if (!canEdit) return;
    draggedSectionRef.current = sIdx;
    setDraggedSectionIndex(sIdx);
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', `section:${sIdx}`);
  };

  const handleSectionDragOver = (e: React.DragEvent, sIdx: number) => {
    const currentDragged = draggedSectionRef.current ?? draggedSectionIndex;
    if (currentDragged === null || currentDragged === sIdx) return;
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
    if (dropSectionIndex !== sIdx) {
      setDropSectionIndex(sIdx);
    }
  };

  const handleSectionDrop = async (e: React.DragEvent, targetIdx: number) => {
    e.preventDefault();
    e.stopPropagation();

    const fromIdx = draggedSectionRef.current ?? draggedSectionIndex;
    draggedSectionRef.current = null;
    setDraggedSectionIndex(null);
    setDropSectionIndex(null);

    if (fromIdx === null || fromIdx === targetIdx) {
      return;
    }

    await handleMoveSection(fromIdx, targetIdx);
  };

  // --- LESSON DRAG AND DROP ---
  const handleLessonDragStart = (
    e: React.DragEvent,
    sectionIndex: number,
    lessonIndex: number,
    lessonId: number
  ) => {
    if (!canEdit) return;
    e.stopPropagation();
    const item = { sectionIndex, lessonIndex, lessonId };
    draggedLessonRef.current = item;
    setDraggedLesson(item);
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', `lesson:${lessonId}`);
  };

  const handleLessonDragOver = (
    e: React.DragEvent,
    sectionIndex: number,
    lessonIndex: number
  ) => {
    const currentDragged = draggedLessonRef.current ?? draggedLesson;
    if (!currentDragged) return;
    e.preventDefault();
    e.stopPropagation();
    e.dataTransfer.dropEffect = 'move';

    if (
      !dropLessonTarget ||
      dropLessonTarget.sectionIndex !== sectionIndex ||
      dropLessonTarget.lessonIndex !== lessonIndex
    ) {
      setDropLessonTarget({ sectionIndex, lessonIndex });
    }
  };

  const handleLessonDrop = async (
    e: React.DragEvent,
    targetSectionIndex: number,
    targetLessonIndex: number
  ) => {
    e.preventDefault();
    e.stopPropagation();

    const dragged = draggedLessonRef.current ?? draggedLesson;
    draggedLessonRef.current = null;
    setDraggedLesson(null);
    setDropLessonTarget(null);

    if (!dragged) {
      return;
    }

    const { sectionIndex: sourceSecIdx, lessonIndex: sourceLesIdx, lessonId } = dragged;

    // If dropped in identical position, no-op
    if (sourceSecIdx === targetSectionIndex && sourceLesIdx === targetLessonIndex) {
      return;
    }

    const previousSections = JSON.parse(JSON.stringify(sections));

    // Case 1: Reordering within the SAME module
    if (sourceSecIdx === targetSectionIndex) {
      const targetSecWrap = sections[sourceSecIdx];
      const secId = getSectionId(targetSecWrap);
      const currLessons = getLessons(targetSecWrap);
      await handleMoveLessonInList(secId, sourceLesIdx, targetLessonIndex, currLessons);
      return;
    }

    // Case 2: Moving ACROSS modules (source -> target)
    const newSections = JSON.parse(JSON.stringify(sections));
    const sourceLessons = getLessons(newSections[sourceSecIdx]);
    const targetLessons = getLessons(newSections[targetSectionIndex]);

    const [movedLesson] = sourceLessons.splice(sourceLesIdx, 1);
    const sourceSecId = getSectionId(newSections[sourceSecIdx]);
    const targetSecId = getSectionId(newSections[targetSectionIndex]);

    // Update moved lesson's section reference
    movedLesson.section_id = targetSecId;
    movedLesson.SectionID = targetSecId;

    // Insert at target position
    targetLessons.splice(targetLessonIndex, 0, movedLesson);

    // Sync back
    newSections[sourceSecIdx].Lessons = sourceLessons;
    newSections[sourceSecIdx].lessons = sourceLessons;
    newSections[targetSectionIndex].Lessons = targetLessons;
    newSections[targetSectionIndex].lessons = targetLessons;

    onSectionsChange(newSections);
    showToast('Урок перенесен в другой модуль');

    try {
      // 1. Update lesson's parent section
      await api.patch(`/lessons/${lessonId}`, {
        section_id: targetSecId,
      });

      // 2. Reorder target section lessons
      const targetIds = targetLessons.map((l: any) => getLessonId(l));
      if (targetIds.length > 0) {
        await api.put(`/sections/${targetSecId}/reorder-lessons`, {
          item_ids: targetIds,
        });
      }

      // 3. Reorder source section lessons
      const sourceIds = sourceLessons.map((l: any) => getLessonId(l));
      if (sourceIds.length > 0) {
        await api.put(`/sections/${sourceSecId}/reorder-lessons`, {
          item_ids: sourceIds,
        });
      }

      showToast('Порядок уроков сохранен');
    } catch (err: any) {
      console.error('Failed to move lesson between sections', err);
      onSectionsChange(previousSections);
      showToast('Ошибка при перемещении урока между модулями', 'error');
      fetchCourseData();
    }
  };

  const handleDragEnd = () => {
    draggedSectionRef.current = null;
    draggedLessonRef.current = null;
    setDraggedSectionIndex(null);
    setDropSectionIndex(null);
    setDraggedLesson(null);
    setDropLessonTarget(null);
  };

  return (
    <div className="space-y-6" onDragEnd={handleDragEnd}>
      {sections.map((secWrap: any, sIdx: number) => {
        const sec = getSection(secWrap);
        const lessons = getLessons(secWrap);
        const secId = getSectionId(secWrap);
        const isModPub = (sec?.Status || sec?.status) === 'published';

        const isSectionBeingDragged = draggedSectionIndex === sIdx;
        const isSectionDropTarget = dropSectionIndex === sIdx && draggedSectionIndex !== sIdx;

        return (
          <div
            key={secId || sIdx}
            onDragOver={(e) => handleSectionDragOver(e, sIdx)}
            onDrop={(e) => handleSectionDrop(e, sIdx)}
            className={`bg-white dark:bg-slate-900 border rounded-3xl p-6 shadow-sm space-y-5 transition-all duration-200 ${
              isSectionBeingDragged
                ? 'opacity-40 scale-[0.99] border-dashed border-indigo-400 dark:border-indigo-600'
                : isSectionDropTarget
                ? 'border-indigo-500 ring-2 ring-indigo-500/20 shadow-lg'
                : 'border-slate-200 dark:border-slate-800'
            }`}
          >
            {/* Module Header (Draggable for Module Reordering) */}
            <div
              draggable={canEdit && draggedLesson === null}
              onDragStart={(e) => handleSectionDragStart(e, sIdx)}
              className="flex flex-col sm:flex-row justify-between sm:items-center gap-3 pb-4 border-b border-slate-100 dark:border-slate-800 cursor-grab active:cursor-grabbing"
            >
              <div className="flex items-center gap-3 min-w-0">
                {/* Grip Handle for Section DnD */}
                {canEdit && (
                  <div
                    className="p-1.5 rounded-lg text-slate-400 hover:text-indigo-600 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors flex-shrink-0"
                    title="Зажмите и перетащите модуль, чтобы изменить его позицию"
                  >
                    <GripVertical size={18} />
                  </div>
                )}

                {/* Module Reorder buttons (Mobile & Click Fallback) */}
                <div
                  className="flex flex-col gap-0.5 text-slate-400 bg-slate-50 dark:bg-slate-800 p-1 rounded-lg border border-slate-200 dark:border-slate-700 flex-shrink-0"
                  onClick={(e) => e.stopPropagation()}
                >
                  <button
                    disabled={sIdx === 0 || !canEdit}
                    onClick={() => handleMoveSection(sIdx, sIdx - 1)}
                    className="hover:text-indigo-600 disabled:opacity-20 text-[10px] p-0.5"
                    title="Переместить модуль выше"
                    type="button"
                  >
                    <ChevronUp size={12} />
                  </button>
                  <button
                    disabled={sIdx === sections.length - 1 || !canEdit}
                    onClick={() => handleMoveSection(sIdx, sIdx + 1)}
                    className="hover:text-indigo-600 disabled:opacity-20 text-[10px] p-0.5"
                    title="Переместить модуль ниже"
                    type="button"
                  >
                    <ChevronDown size={12} />
                  </button>
                </div>

                <div className="min-w-0">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className="bg-indigo-600 text-white text-[10px] font-black px-2.5 py-0.5 rounded-lg shadow-2xs">
                      Модуль {sIdx + 1}
                    </span>
                    <h4 className="text-base font-extrabold text-slate-900 dark:text-white truncate">
                      {sec?.Title || sec?.title}
                    </h4>
                    <span
                      className={`text-[9px] font-extrabold px-2 py-0.5 rounded-md uppercase tracking-wider ${
                        isModPub
                          ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300'
                          : 'bg-amber-100 dark:bg-amber-950 text-amber-700 dark:text-amber-300'
                      }`}
                    >
                      {isModPub ? 'Опубликован' : 'Черновик'}
                    </span>
                  </div>
                  <p className="text-xs text-slate-500 mt-1 line-clamp-1">
                    {sec?.Description || sec?.description || 'Нет описания модуля'}
                  </p>
                </div>
              </div>

              <div onClick={(e) => e.stopPropagation()}>
                <button
                  type="button"
                  onClick={() => onEditModule(secWrap)}
                  className="px-3.5 py-1.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 hover:bg-slate-100 text-slate-700 dark:text-slate-300 text-xs font-bold rounded-xl transition-colors shadow-2xs flex items-center gap-1.5 w-fit flex-shrink-0 cursor-pointer"
                >
                  <Edit3 size={14} />
                  <span>Настройки модуля</span>
                </button>
              </div>
            </div>

            {/* Lessons List in Module */}
            <div
              className="space-y-2.5 pl-2 sm:pl-3"
              onDragOver={(e) => {
                if ((draggedLessonRef.current || draggedLesson) && lessons.length === 0) {
                  handleLessonDragOver(e, sIdx, 0);
                }
              }}
              onDrop={(e) => {
                if ((draggedLessonRef.current || draggedLesson) && lessons.length === 0) {
                  handleLessonDrop(e, sIdx, 0);
                }
              }}
            >
              {lessons.length === 0 && (
                <div
                  onDragOver={(e) => {
                    if (draggedLessonRef.current || draggedLesson) {
                      handleLessonDragOver(e, sIdx, 0);
                    }
                  }}
                  onDrop={(e) => {
                    if (draggedLessonRef.current || draggedLesson) {
                      handleLessonDrop(e, sIdx, 0);
                    }
                  }}
                  className={`p-6 border-2 border-dashed rounded-2xl text-center text-xs text-slate-400 transition-all ${
                    dropLessonTarget?.sectionIndex === sIdx
                      ? 'border-indigo-500 bg-indigo-50/50 dark:bg-indigo-950/30 text-indigo-600 font-bold'
                      : 'border-slate-200 dark:border-slate-800'
                  }`}
                >
                  В этом модуле пока нет уроков. Перетащите урок сюда или добавьте новый.
                </div>
              )}

              {lessons.map((l: any, lesIdx: number) => {
                const lId = getLessonId(l);
                const isLesPub = (l.Status || l.status) === 'published';
                const isLessonBeingDragged =
                  draggedLesson?.sectionIndex === sIdx && draggedLesson?.lessonIndex === lesIdx;
                const isLessonDropTarget =
                  dropLessonTarget?.sectionIndex === sIdx &&
                  dropLessonTarget?.lessonIndex === lesIdx &&
                  !isLessonBeingDragged;

                return (
                  <React.Fragment key={lId || lesIdx}>
                    {/* Visual Insertion Placeholder before target item */}
                    {isLessonDropTarget && (
                      <div className="h-1 bg-indigo-500 rounded-full my-1 shadow-sm animate-pulse" />
                    )}

                    <div
                      draggable={canEdit && draggedSectionIndex === null}
                      onDragStart={(e) => handleLessonDragStart(e, sIdx, lesIdx, lId)}
                      onDragOver={(e) => handleLessonDragOver(e, sIdx, lesIdx)}
                      onDrop={(e) => handleLessonDrop(e, sIdx, lesIdx)}
                      className={`flex flex-col sm:flex-row sm:items-center justify-between p-3.5 bg-slate-50 dark:bg-slate-800/60 border rounded-2xl transition-all gap-3 shadow-xs cursor-grab active:cursor-grabbing ${
                        isLessonBeingDragged
                          ? 'opacity-40 scale-[0.98] border-indigo-400 bg-indigo-50/30 dark:bg-indigo-950/30'
                          : isLessonDropTarget
                          ? 'border-indigo-500 shadow-md ring-2 ring-indigo-500/20'
                          : 'border-slate-200 dark:border-slate-700 hover:border-indigo-300 dark:hover:border-indigo-600'
                      }`}
                    >
                      <div className="flex items-center gap-3 flex-1 min-w-0">
                        {/* Grip Handle for Lesson DnD */}
                        {canEdit && (
                          <div
                            className="p-1 rounded text-slate-400 hover:text-indigo-600 hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors flex-shrink-0"
                            title="Зажмите и перетащите урок в любую позицию или модуль"
                          >
                            <GripVertical size={16} />
                          </div>
                        )}

                        {/* Lesson Reorder buttons (Fallback) */}
                        <div
                          className="flex flex-col gap-0.5 text-slate-400 bg-white dark:bg-slate-900 p-1 rounded-lg border border-slate-200 dark:border-slate-700 flex-shrink-0"
                          onClick={(e) => e.stopPropagation()}
                        >
                          <button
                            disabled={lesIdx === 0 || !canEdit}
                            onClick={() => handleMoveLessonInList(secId, lesIdx, lesIdx - 1, lessons)}
                            className="hover:text-indigo-600 disabled:opacity-20 text-[10px] p-0.5"
                            title="Переместить урок выше"
                            type="button"
                          >
                            <ChevronUp size={11} />
                          </button>
                          <button
                            disabled={lesIdx === lessons.length - 1 || !canEdit}
                            onClick={() => handleMoveLessonInList(secId, lesIdx, lesIdx + 1, lessons)}
                            className="hover:text-indigo-600 disabled:opacity-20 text-[10px] p-0.5"
                            title="Переместить урок ниже"
                            type="button"
                          >
                            <ChevronDown size={11} />
                          </button>
                        </div>

                        <div className="min-w-0 flex-1">
                          <div className="flex items-center gap-2 flex-wrap">
                            <h5 className="text-xs font-bold text-slate-900 dark:text-white truncate">
                              {lesIdx + 1}. {l.Title || l.title}
                            </h5>
                            <span
                              className={`text-[9px] font-extrabold px-1.5 py-0.5 rounded uppercase tracking-wider ${
                                isLesPub
                                  ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300'
                                  : 'bg-amber-100 dark:bg-amber-950 text-amber-700 dark:text-amber-300'
                              }`}
                            >
                              {isLesPub ? 'Опубликован' : 'Черновик'}
                            </span>
                            <span className="text-[9px] bg-indigo-100 dark:bg-indigo-950 text-indigo-700 dark:text-indigo-300 px-1.5 py-0.5 rounded font-bold uppercase">
                              {l.Type || l.type || 'lecture'}
                            </span>
                            {(l.IsFree || l.is_free) && (
                              <span className="text-[9px] bg-blue-100 dark:bg-blue-950 text-blue-700 dark:text-blue-300 px-1.5 py-0.5 rounded font-bold">
                                Demo
                              </span>
                            )}
                          </div>
                          <p className="text-[11px] text-slate-500 line-clamp-1 mt-0.5">
                            {l.Description || l.description || 'Нет описания'}
                          </p>
                        </div>
                      </div>

                      <div className="flex items-center gap-2 flex-shrink-0" onClick={(e) => e.stopPropagation()}>
                        <button
                          type="button"
                          onClick={() => onEditLesson(l)}
                          className="px-3 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300 text-[11px] font-bold rounded-xl hover:bg-slate-100 transition-colors shadow-2xs flex items-center gap-1 cursor-pointer"
                        >
                          <Edit3 size={13} />
                          <span>Изменить</span>
                        </button>

                        <button
                          type="button"
                          onClick={() => router.push(`/teacher/lessons/${lId}/edit`)}
                          className="px-3 py-1.5 bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 text-[11px] font-bold rounded-xl hover:bg-indigo-100 transition-colors shadow-2xs ring-1 ring-indigo-200 dark:ring-indigo-800 flex items-center gap-1 cursor-pointer"
                        >
                          <Sparkles size={13} />
                          <span>Редактор PUCK</span>
                        </button>
                      </div>
                    </div>
                  </React.Fragment>
                );
              })}

              {/* Bottom drop target for appending to end of module */}
              {draggedLesson && (
                <div
                  onDragOver={(e) => handleLessonDragOver(e, sIdx, lessons.length)}
                  onDrop={(e) => handleLessonDrop(e, sIdx, lessons.length)}
                  className={`py-2 text-center text-[11px] rounded-xl border border-dashed transition-all ${
                    dropLessonTarget?.sectionIndex === sIdx &&
                    dropLessonTarget?.lessonIndex === lessons.length
                      ? 'border-indigo-500 bg-indigo-50/50 dark:bg-indigo-950/30 text-indigo-600 font-bold'
                      : 'border-transparent text-slate-400'
                  }`}
                >
                  {dropLessonTarget?.sectionIndex === sIdx &&
                  dropLessonTarget?.lessonIndex === lessons.length
                    ? 'Вставить в конец модуля'
                    : ''}
                </div>
              )}

              {canEdit && (
                <button
                  type="button"
                  onClick={() =>
                    onCreateLesson({
                      sectionId: secId,
                      sectionTitle: sec?.Title || sec?.title,
                      nextIndex: lessons.length + 1,
                    })
                  }
                  className="w-full py-2.5 border-2 border-dashed border-slate-200 dark:border-slate-700 text-slate-500 hover:text-indigo-600 hover:border-indigo-400 rounded-2xl text-xs font-bold transition-all bg-white/40 dark:bg-slate-900/40 flex items-center justify-center gap-1.5 mt-2 cursor-pointer"
                >
                  <PlusCircle size={14} />
                  <span>+ Добавить урок в модуль</span>
                </button>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}
