'use client';

import React, { useState, useEffect, useTransition } from 'react';
import ProtectedRoute from '@/components/ProtectedRoute';
import Sidebar from '@/components/layout/Sidebar';
import TopNavbar from '@/components/layout/TopNavbar';
import {
  AdminUser,
  fetchAdminUsers,
  updateUserRole,
  updateUserBan,
} from '@/lib/admin';
import {
  Shield,
  Users,
  UserCheck,
  GraduationCap,
  Ban,
  Search,
  Filter,
  AlertTriangle,
  CheckCircle,
  Loader2,
  ChevronLeft,
  ChevronRight,
  MoreVertical,
  Check
} from 'lucide-react';

export default function AdminUsersPage() {
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [isLoading, setIsLoading] = useState(true);

  // Filters
  const [search, setSearch] = useState('');
  const [roleFilter, setRoleFilter] = useState('all');
  const [onlyBanned, setOnlyBanned] = useState(false);

  // Ban confirmation modal
  const [banModalUser, setBanModalUser] = useState<AdminUser | null>(null);
  const [isUpdatingBan, setIsUpdatingBan] = useState(false);

  // Action status message
  const [statusMessage, setStatusMessage] = useState<{ text: string; type: 'success' | 'error' } | null>(null);

  const loadUsers = async () => {
    setIsLoading(true);
    try {
      const data = await fetchAdminUsers({
        search: search.trim() || undefined,
        role: roleFilter !== 'all' ? roleFilter : undefined,
        is_banned: onlyBanned ? true : undefined,
        page,
        page_size: pageSize,
      });
      setUsers(data.users);
      setTotal(data.total);
    } catch (err) {
      console.error('Failed to load admin users', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadUsers();
  }, [page, roleFilter, onlyBanned]);

  // Handle Search submit / debounce
  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    loadUsers();
  };

  const handleRoleSelect = async (user: AdminUser, newRole: string) => {
    if (user.role === newRole) return;
    try {
      await updateUserRole(user.id, newRole);
      setUsers((prev) =>
        prev.map((u) => (u.id === user.id ? { ...u, role: newRole } : u))
      );
      setStatusMessage({
        text: `Роль пользователя ${user.username} успешно изменена на ${newRole}`,
        type: 'success',
      });
      setTimeout(() => setStatusMessage(null), 3500);
    } catch {
      setStatusMessage({
        text: 'Не удалось обновить роль пользователя',
        type: 'error',
      });
      setTimeout(() => setStatusMessage(null), 3500);
    }
  };

  const handleConfirmBanToggle = async () => {
    if (!banModalUser) return;
    setIsUpdatingBan(true);
    try {
      const targetBannedState = !banModalUser.is_banned;
      await updateUserBan(banModalUser.id, targetBannedState);
      setUsers((prev) =>
        prev.map((u) => (u.id === banModalUser.id ? { ...u, is_banned: targetBannedState } : u))
      );
      setStatusMessage({
        text: `Пользователь ${banModalUser.username} успешно ${
          targetBannedState ? 'заблокирован' : 'разблокирован'
        }`,
        type: 'success',
      });
      setBanModalUser(null);
      setTimeout(() => setStatusMessage(null), 3500);
    } catch {
      setStatusMessage({
        text: 'Не удалось обновить статус блокировки',
        type: 'error',
      });
      setTimeout(() => setStatusMessage(null), 3500);
    } finally {
      setIsUpdatingBan(false);
    }
  };

  // Stats calculation
  const totalCount = total;
  const activeCount = users.filter((u) => !u.is_banned).length;
  const teachersCount = users.filter((u) => u.role === 'teacher' || u.role === 'author').length;
  const bannedCount = users.filter((u) => u.is_banned).length;

  const getRoleBadge = (role: string) => {
    switch (role) {
      case 'admin':
        return 'bg-amber-100 dark:bg-amber-950/80 text-amber-800 dark:text-amber-300 border-amber-300 dark:border-amber-800';
      case 'teacher':
        return 'bg-indigo-100 dark:bg-indigo-950/80 text-indigo-700 dark:text-indigo-300 border-indigo-200 dark:border-indigo-800';
      case 'author':
        return 'bg-purple-100 dark:bg-purple-950/80 text-purple-700 dark:text-purple-300 border-purple-200 dark:border-purple-800';
      default:
        return 'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-700';
    }
  };

  return (
    <ProtectedRoute allowedRoles={['admin']}>
      <div className="flex min-h-screen bg-slate-50 dark:bg-slate-950">
        <Sidebar />

        <div className="flex-1 flex flex-col min-w-0 ml-20 lg:ml-64">
          <TopNavbar
            title="Панель администратора"
            subtitle="Управление пользователями, распределение ролей и безопасность"
          />

          <main className="p-6 sm:p-8 max-w-7xl w-full mx-auto space-y-6 flex-1">
            {/* Status Toast Alert */}
            {statusMessage && (
              <div
                className={`p-4 rounded-2xl flex items-center justify-between text-xs font-bold border transition-all animate-in slide-in-from-top-2 ${
                  statusMessage.type === 'success'
                    ? 'bg-emerald-50 dark:bg-emerald-950/60 border-emerald-200 text-emerald-800 dark:text-emerald-300'
                    : 'bg-rose-50 dark:bg-rose-950/60 border-rose-200 text-rose-800 dark:text-rose-300'
                }`}
              >
                <div className="flex items-center gap-2">
                  {statusMessage.type === 'success' ? <CheckCircle size={16} /> : <AlertTriangle size={16} />}
                  <span>{statusMessage.text}</span>
                </div>
                <button
                  onClick={() => setStatusMessage(null)}
                  className="opacity-70 hover:opacity-100"
                >
                  ✕
                </button>
              </div>
            )}

            {/* Metrics Overview */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
              <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-5 shadow-xs flex items-center gap-4">
                <div className="w-12 h-12 rounded-2xl bg-indigo-50 dark:bg-indigo-950/60 text-indigo-600 dark:text-indigo-400 flex items-center justify-center flex-shrink-0">
                  <Users size={22} />
                </div>
                <div>
                  <span className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">
                    Всего пользователей
                  </span>
                  <span className="text-2xl font-black text-slate-900 dark:text-white">
                    {totalCount}
                  </span>
                </div>
              </div>

              <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-5 shadow-xs flex items-center gap-4">
                <div className="w-12 h-12 rounded-2xl bg-emerald-50 dark:bg-emerald-950/60 text-emerald-600 dark:text-emerald-400 flex items-center justify-center flex-shrink-0">
                  <UserCheck size={22} />
                </div>
                <div>
                  <span className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">
                    Активные аккаунты
                  </span>
                  <span className="text-2xl font-black text-slate-900 dark:text-white">
                    {activeCount}
                  </span>
                </div>
              </div>

              <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-5 shadow-xs flex items-center gap-4">
                <div className="w-12 h-12 rounded-2xl bg-purple-50 dark:bg-purple-950/60 text-purple-600 dark:text-purple-400 flex items-center justify-center flex-shrink-0">
                  <GraduationCap size={22} />
                </div>
                <div>
                  <span className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">
                    Преподаватели / Авторы
                  </span>
                  <span className="text-2xl font-black text-slate-900 dark:text-white">
                    {teachersCount}
                  </span>
                </div>
              </div>

              <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-5 shadow-xs flex items-center gap-4">
                <div className="w-12 h-12 rounded-2xl bg-rose-50 dark:bg-rose-950/60 text-rose-600 dark:text-rose-400 flex items-center justify-center flex-shrink-0">
                  <Ban size={22} />
                </div>
                <div>
                  <span className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">
                    Заблокированные
                  </span>
                  <span className="text-2xl font-black text-rose-600 dark:text-rose-400">
                    {bannedCount}
                  </span>
                </div>
              </div>
            </div>

            {/* Filters Bar */}
            <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-5 shadow-xs flex flex-col md:flex-row items-stretch md:items-center justify-between gap-4">
              <form onSubmit={handleSearchSubmit} className="flex-1 flex gap-2">
                <div className="relative flex-1">
                  <Search size={16} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input
                    type="text"
                    placeholder="Поиск по имени, username или email..."
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    className="w-full pl-10 pr-4 py-2.5 bg-slate-50 dark:bg-slate-800/70 border border-slate-200 dark:border-slate-700 rounded-2xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                  />
                </div>
                <button
                  type="submit"
                  className="px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-2xl text-xs font-bold transition-colors cursor-pointer"
                >
                  Найти
                </button>
              </form>

              <div className="flex flex-wrap items-center gap-3">
                <div className="flex items-center gap-2">
                  <Filter size={15} className="text-slate-400" />
                  <span className="text-xs font-bold text-slate-600 dark:text-slate-400">Роль:</span>
                  <select
                    value={roleFilter}
                    onChange={(e) => {
                      setRoleFilter(e.target.value);
                      setPage(1);
                    }}
                    className="px-3 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:outline-none"
                  >
                    <option value="all">Все роли</option>
                    <option value="student">Студент</option>
                    <option value="teacher">Преподаватель</option>
                    <option value="author">Автор курса</option>
                    <option value="admin">Администратор</option>
                  </select>
                </div>

                <label className="flex items-center gap-2 text-xs font-bold text-slate-700 dark:text-slate-300 cursor-pointer select-none">
                  <input
                    type="checkbox"
                    checked={onlyBanned}
                    onChange={(e) => {
                      setOnlyBanned(e.target.checked);
                      setPage(1);
                    }}
                    className="w-4 h-4 rounded text-rose-600 focus:ring-rose-500"
                  />
                  <span>Только заблокированные</span>
                </label>
              </div>
            </div>

            {/* Users Table */}
            <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl shadow-xs overflow-hidden">
              <div className="overflow-x-auto">
                <table className="w-full text-left text-xs">
                  <thead className="bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-800 text-[11px] uppercase tracking-wider text-slate-400 font-bold">
                    <tr>
                      <th className="py-4 px-6">Пользователь</th>
                      <th className="py-4 px-6">Назначенная роль</th>
                      <th className="py-4 px-6">Статус</th>
                      <th className="py-4 px-6">Регистрация</th>
                      <th className="py-4 px-6 text-right">Действия</th>
                    </tr>
                  </thead>

                  <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
                    {isLoading ? (
                      <tr>
                        <td colSpan={5} className="py-12 text-center text-slate-400">
                          <Loader2 size={24} className="animate-spin mx-auto mb-2 text-indigo-600" />
                          <span>Загрузка списка пользователей...</span>
                        </td>
                      </tr>
                    ) : users.length === 0 ? (
                      <tr>
                        <td colSpan={5} className="py-12 text-center text-slate-400">
                          <Users size={32} className="mx-auto mb-2 opacity-50" />
                          <span>Пользователи по заданным критериям не найдены</span>
                        </td>
                      </tr>
                    ) : (
                      users.map((user) => (
                        <tr
                          key={user.id}
                          className="hover:bg-slate-50/50 dark:hover:bg-slate-800/40 transition-colors"
                        >
                          {/* User info */}
                          <td className="py-4 px-6">
                            <div className="flex items-center gap-3">
                              <div className="w-9 h-9 rounded-xl bg-gradient-to-tr from-indigo-500 to-purple-600 text-white font-bold text-xs flex items-center justify-center flex-shrink-0 shadow-2xs">
                                {user.first_name?.[0] || user.username[0]?.toUpperCase() || 'U'}
                              </div>
                              <div className="min-w-0">
                                <div className="font-extrabold text-slate-900 dark:text-white truncate">
                                  {user.first_name ? `${user.first_name} ${user.last_name || ''}` : user.username}
                                </div>
                                <div className="text-[11px] text-slate-400 truncate">
                                  {user.email} • @{user.username}
                                </div>
                              </div>
                            </div>
                          </td>

                          {/* Role Select */}
                          <td className="py-4 px-6">
                            <div className="inline-flex items-center">
                              <select
                                value={user.role}
                                onChange={(e) => handleRoleSelect(user, e.target.value)}
                                className={`px-2.5 py-1 rounded-xl text-xs font-bold border cursor-pointer focus:outline-none transition-all ${getRoleBadge(
                                  user.role
                                )}`}
                              >
                                <option value="student">Студент</option>
                                <option value="teacher">Преподаватель</option>
                                <option value="author">Автор курса</option>
                                <option value="admin">Администратор</option>
                              </select>
                            </div>
                          </td>

                          {/* Status */}
                          <td className="py-4 px-6">
                            {user.is_banned ? (
                              <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[10px] font-extrabold bg-rose-100 dark:bg-rose-950/80 text-rose-700 dark:text-rose-300 border border-rose-200 dark:border-rose-900">
                                <Ban size={12} />
                                <span>Заблокирован</span>
                              </span>
                            ) : (
                              <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[10px] font-extrabold bg-emerald-100 dark:bg-emerald-950/80 text-emerald-700 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-900">
                                <CheckCircle size={12} />
                                <span>Активен</span>
                              </span>
                            )}
                          </td>

                          {/* Registered date */}
                          <td className="py-4 px-6 text-slate-500 text-[11px]">
                            {new Date(user.created_at).toLocaleDateString('ru-RU')}
                          </td>

                          {/* Actions */}
                          <td className="py-4 px-6 text-right">
                            <button
                              type="button"
                              onClick={() => setBanModalUser(user)}
                              className={`px-3 py-1.5 rounded-xl font-bold text-xs transition-colors cursor-pointer border ${
                                user.is_banned
                                  ? 'bg-emerald-50 text-emerald-600 border-emerald-200 hover:bg-emerald-100 dark:bg-emerald-950/60 dark:border-emerald-900 dark:text-emerald-300'
                                  : 'bg-rose-50 text-rose-600 border-rose-200 hover:bg-rose-100 dark:bg-rose-950/60 dark:border-rose-900 dark:text-rose-300'
                              }`}
                            >
                              {user.is_banned ? 'Разблокировать' : 'Заблокировать'}
                            </button>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>

              {/* Pagination */}
              {total > pageSize && (
                <div className="p-4 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between text-xs text-slate-500">
                  <span>
                    Показано {(page - 1) * pageSize + 1} - {Math.min(page * pageSize, total)} из {total}
                  </span>
                  <div className="flex gap-2">
                    <button
                      type="button"
                      disabled={page <= 1}
                      onClick={() => setPage(page - 1)}
                      className="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 disabled:opacity-40 hover:bg-slate-50 dark:hover:bg-slate-800"
                    >
                      <ChevronLeft size={16} />
                    </button>
                    <button
                      type="button"
                      disabled={page * pageSize >= total}
                      onClick={() => setPage(page + 1)}
                      className="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 disabled:opacity-40 hover:bg-slate-50 dark:hover:bg-slate-800"
                    >
                      <ChevronRight size={16} />
                    </button>
                  </div>
                </div>
              )}
            </div>
          </main>
        </div>

        {/* Ban / Unban Confirmation Modal */}
        {banModalUser && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4 animate-in fade-in duration-150">
            <div className="bg-white dark:bg-slate-900 rounded-3xl p-6 sm:p-8 max-w-md w-full border border-slate-200 dark:border-slate-800 shadow-2xl space-y-5">
              <div
                className={`w-12 h-12 rounded-2xl flex items-center justify-center mx-auto ${
                  banModalUser.is_banned
                    ? 'bg-emerald-50 dark:bg-emerald-950/60 text-emerald-600'
                    : 'bg-rose-50 dark:bg-rose-950/60 text-rose-600'
                }`}
              >
                {banModalUser.is_banned ? <CheckCircle size={24} /> : <AlertTriangle size={24} />}
              </div>

              <div className="text-center space-y-2">
                <h3 className="text-lg font-extrabold text-slate-900 dark:text-white">
                  {banModalUser.is_banned ? 'Разблокировать аккаунт?' : 'Заблокировать пользователя?'}
                </h3>
                <p className="text-xs text-slate-500 leading-relaxed">
                  {banModalUser.is_banned
                    ? `Пользователь ${banModalUser.username} (${banModalUser.email}) вновь получит доступ к платформе и учебным материалам.`
                    : `Пользователь ${banModalUser.username} (${banModalUser.email}) не сможет авторизоваться в системе и выполнять задания до момента разблокировки.`}
                </p>
              </div>

              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  disabled={isUpdatingBan}
                  onClick={() => setBanModalUser(null)}
                  className="flex-1 py-3 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors"
                >
                  Отмена
                </button>
                <button
                  type="button"
                  disabled={isUpdatingBan}
                  onClick={handleConfirmBanToggle}
                  className={`flex-1 py-3 rounded-xl text-white text-xs font-bold shadow-md transition-all flex items-center justify-center gap-2 cursor-pointer ${
                    banModalUser.is_banned
                      ? 'bg-emerald-600 hover:bg-emerald-700 shadow-emerald-600/20'
                      : 'bg-rose-600 hover:bg-rose-700 shadow-rose-600/20'
                  }`}
                >
                  {isUpdatingBan ? (
                    <Loader2 size={15} className="animate-spin" />
                  ) : banModalUser.is_banned ? (
                    'Разблокировать'
                  ) : (
                    'Да, заблокировать'
                  )}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </ProtectedRoute>
  );
}
