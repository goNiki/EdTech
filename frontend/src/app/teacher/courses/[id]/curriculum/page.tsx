'use client';

import React, { useEffect, useState, use } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { api } from '@/lib/api';
import TopNavbar from '@/components/layout/TopNavbar';
import CourseAnalyticsCards from '@/components/teacher/CourseAnalyticsCards';
import StudentsTable, { StudentItem } from '@/components/teacher/StudentsTable';
import PendingHomeworksQueue, { PendingHWItem } from '@/components/teacher/PendingHomeworksQueue';
import ModalGradeHW from '@/components/teacher/ModalGradeHW';
import ModalStudentDrilldown from '@/components/teacher/ModalStudentDrilldown';
import ModalAddStudent from '@/components/teacher/ModalAddStudent';
import ModalEditModule from '@/components/teacher/ModalEditModule';
import ModalEditLesson from '@/components/teacher/ModalEditLesson';
import ModalCreateModule from '@/components/teacher/ModalCreateModule';
import ModalCreateLesson from '@/components/teacher/ModalCreateLesson';
import { useAuth } from '@/store/useAuth';
import {
  Layers,
  Users,
  Settings,
  PlusCircle,
  Edit3,
  Sparkles,
  ArrowUp,
  ArrowDown,
  BookOpen,
  CheckCircle2,
  FileText,
  HelpCircle,
  Award,
  Globe,
  Lock,
  EyeOff,
  ShieldAlert
} from 'lucide-react';

export default function TeacherCourseManagementPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const searchParams = useSearchParams();
  const { user } = useAuth();

  // Tab State: 'editor' | 'students'
  const initialTab = searchParams.get('tab') === 'students' ? 'students' : 'editor';
  const [activeTab, setActiveTab] = useState<'editor' | 'students'>(initialTab);
  const [activeSubTab, setActiveSubTab] = useState<'list' | 'pending'>('list');

  // Course & Structure state
  const [courseData, setCourseData] = useState<any>(null);
  const [sections, setSections] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isTogglingStatus, setIsTogglingStatus] = useState(false);

  // RBAC permission check
  const canEdit = React.useMemo(() => {
    if (courseData?.permissions?.can_edit !== undefined) {
      return Boolean(courseData.permissions.can_edit);
    }
    if (courseData?.can_edit !== undefined) {
      return Boolean(courseData.can_edit);
    }
    if (user?.role === 'admin') return true;
    const creatorId = courseData?.CreatedBy ?? courseData?.created_by;
    if (creatorId && user?.id && user.id !== creatorId && user?.role !== 'teacher') {
      return false;
    }
    return true;
  }, [courseData, user]);

  // Analytics & Students state
  const [analytics, setAnalytics] = useState({
    total_students: 0,
    avg_progress: 0,
    avg_score: 0,
    pending_homeworks_count: 0,
  });
  const [students, setStudents] = useState<StudentItem[]>([]);
  const [pendingHWs, setPendingHWs] = useState<PendingHWItem[]>([]);

  // Modals state
  const [isCreateModuleOpen, setIsCreateModuleOpen] = useState(false);
  const [createLessonTarget, setCreateLessonTarget] = useState<{
    sectionId: number;
    sectionTitle: string;
    nextIndex: number;
  } | null>(null);
  const [selectedModule, setSelectedModule] = useState<any | null>(null);
  const [selectedLesson, setSelectedLesson] = useState<any | null>(null);
  const [drilldownStudentId, setDrilldownStudentId] = useState<number | null>(null);
  const [gradingHW, setGradingHW] = useState<PendingHWItem | null>(null);
  const [isAddStudentOpen, setIsAddStudentOpen] = useState(false);

  // Notifications
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);
  const showToast = (message: string, type: 'success' | 'error' = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 4000);
  };

  const fetchCourseData = async () => {
    try {
      // 1. Structure
      const structRes = await api.get(`/courses/${id}/structure`);
      const sData = structRes.data.data || structRes.data;
      setCourseData(sData.Course || sData.course);
      setSections(sData.Sections || sData.sections || []);

      // 2. Analytics
      const anRes = await api.get(`/courses/${id}/analytics`).catch(() => null);
      if (anRes) {
        const aData = anRes.data.data || anRes.data;
        setAnalytics({
          total_students: aData.total_students || 0,
          avg_progress: Math.round(aData.avg_progress_percent || 0),
          avg_score: Math.round(aData.avg_score || 0),
          pending_homeworks_count: aData.pending_homeworks_count || 0,
        });
      }

      // 3. Students list
      const stRes = await api.get(`/courses/${id}/students`).catch(() => null);
      if (stRes) {
        const stData = stRes.data.data || stRes.data;
        const list = (stData.students || stData || []).map((s: any) => ({
          id: s.id || s.user_id,
          name: s.name || `${s.first_name || ''} ${s.last_name || ''}`.trim() || s.username,
          email: s.email,
          username: s.username,
          progress: Math.round(s.progress_percentage ?? s.progress_percent ?? s.progress ?? 0),
          score: Math.round(s.average_score ?? s.score ?? 0),
          completed_lessons: s.completed_lessons ?? s.completed_lessons_count ?? 0,
          has_pending: s.has_pending_homeworks ?? s.has_pending ?? false,
        }));
        setStudents(list);
      }

      // 4. Pending Homeworks
      const hwRes = await api.get(`/courses/${id}/grading/pending`).catch(() => null);
      if (hwRes) {
        const hwData = hwRes.data.data || hwRes.data;
        setPendingHWs(hwData.items || hwData || []);
      }
    } catch (err) {
      console.error('Failed to load course management data', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchCourseData();
  }, [id]);

  const handleTogglePublishCourse = async () => {
    if (!canEdit) {
      showToast('У вас нет прав на редактирование этого курса. Изменения не сохранены.', 'error');
      return;
    }
    setIsTogglingStatus(true);
    try {
      const isCurrentlyPublished = (courseData?.Status || courseData?.status) === 'published';

      if (!isCurrentlyPublished) {
        await api.post(`/courses/${id}/publish`);
        showToast('Курс и все связанные модули и уроки успешно опубликованы!');
      } else {
        await api.patch(`/courses/${id}/status`, {
          status: 'draft',
        });
        showToast('Курс и все связанные модули переведены в статус «Черновик»');
      }
      await fetchCourseData();
    } catch (err: any) {
      console.error('Failed to change course status', err);
      if (err.response?.status === 403) {
        showToast('У вас нет прав на редактирование этого курса. Изменения не сохранены.', 'error');
      } else {
        showToast(err.response?.data?.message || err.response?.data?.error || 'Ошибка при изменении статуса курса', 'error');
      }
    } finally {
      setIsTogglingStatus(false);
    }
  };

  const handleMoveSection = async (fromIdx: number, toIdx: number) => {
    if (!canEdit) {
      showToast('У вас нет прав на редактирование этого курса. Изменения не сохранены.', 'error');
      return;
    }
    if (toIdx < 0 || toIdx >= sections.length) return;
    const copy = [...sections];
    const [moved] = copy.splice(fromIdx, 1);
    copy.splice(toIdx, 0, moved);
    setSections(copy);

    try {
      const secIds = copy.map((s: any) => Number((s.Section || s.section)?.ID || (s.Section || s.section)?.id));
      await api.put(`/courses/${id}/reorder-sections`, {
        item_ids: secIds,
      });
      showToast('Порядок модулей обновлен');
    } catch (err: any) {
      console.error('Failed to reorder sections', err);
      if (err.response?.status === 403) {
        showToast('У вас нет прав на редактирование этого курса. Изменения не сохранены.', 'error');
      } else {
        showToast('Ошибка при изменении порядка модулей', 'error');
      }
      fetchCourseData();
    }
  };

  const handleMoveLessonInList = async (sectionId: number, fromIdx: number, toIdx: number, lessonsList: any[]) => {
    if (!canEdit) {
      showToast('У вас нет прав на редактирование этого курса. Изменения не сохранены.', 'error');
      return;
    }
    if (toIdx < 0 || toIdx >= lessonsList.length) return;
    const copy = [...lessonsList];
    const [moved] = copy.splice(fromIdx, 1);
    copy.splice(toIdx, 0, moved);

    try {
      const lessonIds = copy.map((l: any) => Number(l.ID || l.id));
      await api.put(`/sections/${sectionId}/reorder-lessons`, {
        item_ids: lessonIds,
      });
      showToast('Порядок уроков сохранен');
      fetchCourseData();
    } catch (err: any) {
      console.error('Failed to reorder lessons', err);
      if (err.response?.status === 403) {
        showToast('У вас нет прав на редактирование этого курса. Изменения не сохранены.', 'error');
      } else {
        showToast('Ошибка при изменении порядка уроков', 'error');
      }
      fetchCourseData();
    }
  };

  const handleRemoveStudent = async (student: StudentItem) => {
    if (!canEdit) {
      showToast('У вас нет прав на управление студентами этого курса.', 'error');
      return;
    }
    try {
      await api.delete(`/courses/${id}/students/${student.id}`);
      showToast(`Студент ${student.name} успешно исключен из курса`);
      setStudents((prev) => prev.filter((s) => s.id !== student.id));
      setAnalytics((prev: any) =>
        prev ? { ...prev, total_students: Math.max(0, (prev.total_students || 1) - 1) } : prev
      );
    } catch (err: any) {
      console.error('Failed to remove student', err);
      if (err.response?.status === 403) {
        showToast('Недостаточно прав для исключения студента из курса', 'error');
      } else {
        showToast(
          err.response?.data?.message || err.response?.data?.error || 'Ошибка при исключении студента',
          'error'
        );
      }
    }
  };

  if (isLoading) {
    return (
      <div className="flex flex-col min-h-screen">
        <TopNavbar title="Загрузка управления курсом..." />
        <main className="p-8 max-w-7xl w-full mx-auto animate-pulse space-y-6">
          <div className="h-44 bg-slate-200 dark:bg-slate-800 rounded-3xl" />
          <div className="h-64 bg-slate-200 dark:bg-slate-800 rounded-3xl" />
        </main>
      </div>
    );
  }

  const isPublished = (courseData?.Status || courseData?.status) === 'published';

  return (
    <div className="flex flex-col min-h-screen">
      <TopNavbar
        title={courseData?.Title || courseData?.title || 'Управление курсом'}
        subtitle="Редактирование структуры, публикация и работа со студентами"
      />

      <main className="p-8 max-w-7xl w-full mx-auto space-y-8 flex-1">
        {/* Course Header Banner */}
        <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-sm flex flex-col md:flex-row justify-between md:items-center gap-6">
          <div className="space-y-2">
            <div className="flex flex-wrap items-center gap-2">
              <span
                className={`px-3 py-1 rounded-xl text-xs font-black uppercase tracking-wider ${
                  isPublished
                    ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300 ring-1 ring-emerald-500/20'
                    : 'bg-amber-100 dark:bg-amber-950 text-amber-700 dark:text-amber-300 ring-1 ring-amber-500/20'
                }`}
              >
                {isPublished ? '● Опубликован' : '○ Черновик'}
              </span>
              <span className="px-2.5 py-0.5 rounded-lg text-[10px] font-bold bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300">
                Видимость: {courseData?.Visibility || courseData?.visibility || 'public'}
              </span>
              <span className="px-2.5 py-0.5 rounded-lg text-[10px] font-bold bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400">
                Категория: {courseData?.CategoryID || courseData?.category_id || 'IT'}
              </span>
            </div>
            <h2 className="text-2xl font-extrabold text-slate-900 dark:text-white">
              {courseData?.Title || courseData?.title}
            </h2>
            <p className="text-xs text-slate-500 max-w-2xl">
              {courseData?.Description || courseData?.description || 'Нет описания курса.'}
            </p>
          </div>

          <div className="flex flex-wrap items-center gap-3 flex-shrink-0">
            {/* Single Course Publication Action Button */}
            {isPublished ? (
              <button
                onClick={handleTogglePublishCourse}
                disabled={isTogglingStatus || !canEdit}
                className="px-5 py-2.5 rounded-2xl bg-amber-500 hover:bg-amber-600 disabled:opacity-40 disabled:cursor-not-allowed text-white text-xs font-bold transition-all shadow-md flex items-center gap-2"
                title={!canEdit ? 'У вас нет прав на редактирование этого курса' : undefined}
              >
                <EyeOff size={16} />
                <span>{isTogglingStatus ? 'Обновление...' : 'Снять с публикации (В черновик)'}</span>
              </button>
            ) : (
              <button
                onClick={handleTogglePublishCourse}
                disabled={isTogglingStatus || !canEdit}
                className="px-6 py-2.5 rounded-2xl bg-emerald-600 hover:bg-emerald-700 disabled:opacity-40 disabled:cursor-not-allowed text-white text-xs font-extrabold transition-all shadow-lg shadow-emerald-600/20 flex items-center gap-2"
                title={!canEdit ? 'У вас нет прав на редактирование этого курса' : undefined}
              >
                <Sparkles size={16} />
                <span>{isTogglingStatus ? 'Публикация...' : 'Опубликовать курс'}</span>
              </button>
            )}

            <button
              onClick={() => router.push(`/teacher/courses/${id}/settings`)}
              className="px-4 py-2.5 rounded-2xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 hover:bg-slate-50 dark:hover:bg-slate-700 text-xs font-bold text-slate-700 dark:text-slate-200 transition-all flex items-center gap-2 shadow-xs"
            >
              <Settings size={16} />
              <span>Настройки курса</span>
            </button>
          </div>
        </div>

        {/* Read-Only Warning Banner */}
        {!canEdit && (
          <div className="bg-amber-500/10 border border-amber-500/30 rounded-3xl p-5 flex items-center justify-between gap-4 text-amber-900 dark:text-amber-200 shadow-xs">
            <div className="flex items-center gap-3">
              <div className="p-2.5 rounded-2xl bg-amber-500/20 text-amber-600 dark:text-amber-400 flex-shrink-0">
                <Lock size={20} />
              </div>
              <div>
                <h4 className="text-sm font-extrabold">Режим «Только чтение» (View-Only)</h4>
                <p className="text-xs text-amber-700/90 dark:text-amber-300/90 mt-0.5">
                  У вас нет прав на изменение структуры, уроков или публикацию этого курса.
                </p>
              </div>
            </div>
            <span className="px-3 py-1 bg-amber-500/20 text-amber-800 dark:text-amber-300 text-[10px] font-black uppercase tracking-wider rounded-xl ring-1 ring-amber-500/30">
              Только просмотр
            </span>
          </div>
        )}

        {/* 2 Main Tabs Switcher */}
        <div className="flex gap-2 p-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-fit shadow-xs">
          <button
            onClick={() => setActiveTab('editor')}
            className={`flex items-center gap-2 px-5 py-2.5 rounded-xl text-xs font-bold transition-all ${
              activeTab === 'editor'
                ? 'bg-indigo-600 text-white shadow-sm'
                : 'text-slate-600 dark:text-slate-400 hover:text-slate-900'
            }`}
          >
            <Layers size={16} />
            <span>Программа курса и уроки</span>
          </button>
          <button
            onClick={() => setActiveTab('students')}
            className={`flex items-center gap-2 px-5 py-2.5 rounded-xl text-xs font-bold transition-all relative ${
              activeTab === 'students'
                ? 'bg-indigo-600 text-white shadow-sm'
                : 'text-slate-600 dark:text-slate-400 hover:text-slate-900'
            }`}
          >
            <Users size={16} />
            <span>Студенты и проверка ДЗ</span>
            {analytics.pending_homeworks_count > 0 && (
              <span className="w-2 h-2 rounded-full bg-rose-500 animate-pulse ml-1" />
            )}
          </button>
        </div>

        {/* TAB 1: CURRICULUM EDITOR */}
        {activeTab === 'editor' && (
          <div className="space-y-6">
            <div className="flex justify-between items-center">
              <div>
                <h3 className="text-lg font-extrabold text-slate-900 dark:text-white">
                  Структура учебных модулей
                </h3>
                <p className="text-xs text-slate-500">
                  Перемещайте модули и уроки стрелками, настраивайте параметры и наполняйте контентом в Puck
                </p>
              </div>
              {canEdit && (
                <button
                  onClick={() => setIsCreateModuleOpen(true)}
                  className="px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-2xl transition-all shadow-md flex items-center gap-1.5"
                >
                  <PlusCircle size={16} />
                  <span>+ Добавить модуль</span>
                </button>
              )}
            </div>

            {sections.length === 0 ? (
              <div className="p-16 text-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 shadow-sm space-y-4">
                <div className="w-12 h-12 rounded-2xl bg-indigo-50 dark:bg-slate-800 text-indigo-600 flex items-center justify-center mx-auto">
                  <Layers size={24} />
                </div>
                <h4 className="text-base font-bold text-slate-800 dark:text-slate-200">
                  В курсе пока нет модулей
                </h4>
                <p className="text-xs text-slate-500 max-w-sm mx-auto">
                  {canEdit
                    ? 'Нажмите кнопку создания модуля, чтобы задать название, цели и наполнить курс уроками.'
                    : 'В данном курсе ещё нет добавленных модулей.'}
                </p>
                {canEdit && (
                  <button
                    onClick={() => setIsCreateModuleOpen(true)}
                    className="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-xl shadow-md"
                  >
                    Создать первый модуль
                  </button>
                )}
              </div>
            ) : (
              <div className="space-y-6">
                {sections.map((secWrap: any, sIdx: number) => {
                  const sec = secWrap.Section || secWrap.section;
                  const lessons = secWrap.Lessons || secWrap.lessons || [];
                  const secId = sec?.ID || sec?.id;
                  const isModPub = (sec?.Status || sec?.status) === 'published';

                  return (
                    <div
                      key={secId || sIdx}
                      className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-sm space-y-5"
                    >
                      {/* Module Header */}
                      <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-3 pb-4 border-b border-slate-100 dark:border-slate-800">
                        <div className="flex items-center gap-3">
                          {/* Module Reorder buttons */}
                          <div className="flex flex-col gap-0.5 text-slate-400 bg-slate-50 dark:bg-slate-800 p-1 rounded-lg border border-slate-200 dark:border-slate-700">
                            <button
                              disabled={sIdx === 0 || !canEdit}
                              onClick={() => handleMoveSection(sIdx, sIdx - 1)}
                              className="hover:text-indigo-600 disabled:opacity-20 text-[10px]"
                              title="Переместить модуль выше"
                            >
                              ▲
                            </button>
                            <button
                              disabled={sIdx === sections.length - 1 || !canEdit}
                              onClick={() => handleMoveSection(sIdx, sIdx + 1)}
                              className="hover:text-indigo-600 disabled:opacity-20 text-[10px]"
                              title="Переместить модуль ниже"
                            >
                              ▼
                            </button>
                          </div>

                          <div>
                            <div className="flex items-center gap-2 flex-wrap">
                              <span className="bg-indigo-600 text-white text-[10px] font-black px-2.5 py-0.5 rounded-lg shadow-2xs">
                                Модуль {sIdx + 1}
                              </span>
                              <h4 className="text-base font-extrabold text-slate-900 dark:text-white">
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
                            <p className="text-xs text-slate-500 mt-1">
                              {sec?.Description || sec?.description || 'Нет описания модуля'}
                            </p>
                          </div>
                        </div>

                        <button
                          onClick={() => setSelectedModule(secWrap)}
                          className="px-3.5 py-1.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 hover:bg-slate-100 text-slate-700 dark:text-slate-300 text-xs font-bold rounded-xl transition-colors shadow-2xs flex items-center gap-1.5 w-fit"
                        >
                          <Edit3 size={14} />
                          <span>Настройки модуля</span>
                        </button>
                      </div>

                      {/* Lessons List in Module */}
                      <div className="space-y-2.5 pl-2">
                        {lessons.map((l: any, lesIdx: number) => {
                          const lId = l.ID || l.id;
                          const isLesPub = (l.Status || l.status) === 'published';

                          return (
                            <div
                              key={lId || lesIdx}
                              className="flex flex-col sm:flex-row sm:items-center justify-between p-3.5 bg-slate-50 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700 rounded-2xl hover:border-indigo-300 dark:hover:border-indigo-600 transition-all gap-3 shadow-xs"
                            >
                              <div className="flex items-center gap-3 flex-1 min-w-0">
                                {/* Lesson Reorder buttons */}
                                <div className="flex flex-col gap-0.5 text-slate-400 bg-white dark:bg-slate-900 p-1 rounded-lg border border-slate-200 dark:border-slate-700 flex-shrink-0">
                                  <button
                                    disabled={lesIdx === 0 || !canEdit}
                                    onClick={() => handleMoveLessonInList(secId, lesIdx, lesIdx - 1, lessons)}
                                    className="hover:text-indigo-600 disabled:opacity-20 text-[10px]"
                                    title="Переместить урок выше"
                                  >
                                    ▲
                                  </button>
                                  <button
                                    disabled={lesIdx === lessons.length - 1 || !canEdit}
                                    onClick={() => handleMoveLessonInList(secId, lesIdx, lesIdx + 1, lessons)}
                                    className="hover:text-indigo-600 disabled:opacity-20 text-[10px]"
                                    title="Переместить урок ниже"
                                  >
                                    ▼
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

                              <div className="flex items-center gap-2 flex-shrink-0">
                                <button
                                  onClick={() => setSelectedLesson(l)}
                                  className="px-3 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300 text-[11px] font-bold rounded-xl hover:bg-slate-100 transition-colors shadow-2xs flex items-center gap-1"
                                >
                                  <Edit3 size={13} />
                                  <span>Изменить урок</span>
                                </button>

                                <button
                                  onClick={() => router.push(`/teacher/lessons/${lId}/edit`)}
                                  className="px-3 py-1.5 bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 text-[11px] font-bold rounded-xl hover:bg-indigo-100 transition-colors shadow-2xs ring-1 ring-indigo-200 dark:ring-indigo-800 flex items-center gap-1"
                                >
                                  <Sparkles size={13} />
                                  <span>Редактор PUCK</span>
                                </button>
                              </div>
                            </div>
                          );
                        })}

                        {canEdit && (
                          <button
                            onClick={() => setCreateLessonTarget({
                              sectionId: secId,
                              sectionTitle: sec?.Title || sec?.title,
                              nextIndex: lessons.length + 1,
                            })}
                            className="w-full py-2.5 border-2 border-dashed border-slate-200 dark:border-slate-700 text-slate-500 hover:text-indigo-600 hover:border-indigo-400 rounded-2xl text-xs font-bold transition-all bg-white/40 dark:bg-slate-900/40 flex items-center justify-center gap-1.5 mt-2"
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
            )}
          </div>
        )}

        {/* TAB 2: STUDENTS & HOMEWORKS */}
        {activeTab === 'students' && (
          <div className="space-y-6">
            {/* 4 Analytics Metrics */}
            <CourseAnalyticsCards
              totalStudents={analytics.total_students}
              avgProgress={analytics.avg_progress}
              avgScore={analytics.avg_score}
              pendingCount={analytics.pending_homeworks_count}
            />

            {/* Sub-tab Navigation */}
            <div className="flex gap-2 p-1 bg-slate-100 dark:bg-slate-800 rounded-xl w-fit">
              <button
                onClick={() => setActiveSubTab('list')}
                className={`px-4 py-1.5 text-xs font-bold rounded-lg transition-all ${
                  activeSubTab === 'list'
                    ? 'bg-white dark:bg-indigo-600 text-indigo-600 dark:text-white shadow-xs'
                    : 'text-slate-600 dark:text-slate-400 hover:text-slate-900'
                }`}
              >
                Список группы ({students.length})
              </button>
              <button
                onClick={() => setActiveSubTab('pending')}
                className={`px-4 py-1.5 text-xs font-bold rounded-lg transition-all flex items-center gap-1.5 ${
                  activeSubTab === 'pending'
                    ? 'bg-white dark:bg-indigo-600 text-indigo-600 dark:text-white shadow-xs'
                    : 'text-slate-600 dark:text-slate-400 hover:text-slate-900'
                }`}
              >
                <span>Очередь проверки ДЗ</span>
                {pendingHWs.length > 0 && (
                  <span className="px-1.5 py-0.2 bg-rose-500 text-white rounded-full text-[9px]">
                    {pendingHWs.length}
                  </span>
                )}
              </button>
            </div>

            {/* Sub-tab Content */}
            {activeSubTab === 'list' ? (
              <StudentsTable
                students={students}
                onOpenDrilldown={(stId) => setDrilldownStudentId(stId)}
                onAddStudent={() => setIsAddStudentOpen(true)}
                onRemoveStudent={handleRemoveStudent}
              />
            ) : (
              <PendingHomeworksQueue
                items={pendingHWs}
                onOpenGradeModal={(hw) => setGradingHW(hw)}
              />
            )}
          </div>
        )}
      </main>

      {/* Modals */}
      <ModalCreateModule
        courseId={Number(id)}
        isOpen={isCreateModuleOpen}
        orderIndex={sections.length + 1}
        onClose={() => setIsCreateModuleOpen(false)}
        onCreated={() => {
          showToast('Новый модуль успешно создан!');
          fetchCourseData();
        }}
      />

      <ModalCreateLesson
        courseId={Number(id)}
        sectionId={createLessonTarget?.sectionId || null}
        sectionTitle={createLessonTarget?.sectionTitle}
        orderIndex={createLessonTarget?.nextIndex || 1}
        isOpen={Boolean(createLessonTarget)}
        onClose={() => setCreateLessonTarget(null)}
        onCreated={() => {
          showToast('Новый урок успешно создан!');
          fetchCourseData();
        }}
      />

      <ModalEditModule
        moduleData={selectedModule}
        isOpen={Boolean(selectedModule)}
        onClose={() => setSelectedModule(null)}
        onSaved={fetchCourseData}
      />

      <ModalEditLesson
        lessonData={selectedLesson}
        isOpen={Boolean(selectedLesson)}
        onClose={() => setSelectedLesson(null)}
        onSaved={fetchCourseData}
      />

      <ModalStudentDrilldown
        courseId={Number(id)}
        studentId={drilldownStudentId}
        onClose={() => setDrilldownStudentId(null)}
      />

      <ModalGradeHW
        hw={gradingHW}
        onClose={() => setGradingHW(null)}
        onGraded={fetchCourseData}
      />

      <ModalAddStudent
        courseId={Number(id)}
        isOpen={isAddStudentOpen}
        onClose={() => setIsAddStudentOpen(false)}
        onAdded={fetchCourseData}
      />

      {/* Toast Notification */}
      {toast && (
        <div
          className={`fixed bottom-6 right-6 z-50 text-white text-xs font-bold px-4 py-3 rounded-2xl shadow-xl flex items-center gap-2 animate-bounce ${
            toast.type === 'error' ? 'bg-rose-600' : 'bg-indigo-600'
          }`}
        >
          {toast.type === 'error' ? <ShieldAlert size={16} /> : <CheckCircle2 size={16} />}
          <span>{toast.message}</span>
        </div>
      )}
    </div>
  );
}
