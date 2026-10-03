'use client';

import React, { useState, useEffect } from 'react';
import { useAuth } from '@/store/useAuth';
import { api } from '@/lib/api';
import { useTheme } from 'next-themes';
import TopNavbar from '@/components/layout/TopNavbar';
import {
  User as UserIcon,
  Settings,
  ShieldCheck,
  Calendar,
  Clock,
  Edit2,
  Save,
  Moon,
  Sun,
  Laptop,
  Sparkles,
  CheckCircle2
} from 'lucide-react';

export default function ProfileAndSettingsPage() {
  const { user, setUser } = useAuth();
  const { theme, setTheme } = useTheme();

  const [activeTab, setActiveTab] = useState<'profile' | 'preferences'>('profile');
  const [isEditing, setIsEditing] = useState(false);

  // Form states
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [bio, setBio] = useState('');
  const [avatarUrl, setAvatarUrl] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [saveSuccess, setSaveSuccess] = useState(false);

  useEffect(() => {
    if (user) {
      setFirstName(user.first_name || '');
      setLastName(user.last_name || '');
      setBio(user.bio || '');
      setAvatarUrl(user.avatar_url || '');
    }
  }, [user]);

  const handleSaveProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    try {
      const res = await api.patch('/auth/profile', {
        first_name: firstName,
        last_name: lastName,
        bio: bio,
        avatar_url: avatarUrl,
      });
      const updated = res.data.data || res.data;
      if (updated) {
        setUser({ ...user, ...updated });
      }
      setIsEditing(false);
      setSaveSuccess(true);
      setTimeout(() => setSaveSuccess(false), 3000);
    } catch (err: any) {
      alert(err.response?.data?.message || err.response?.data?.error || 'Ошибка при сохранении профиля');
    } finally {
      setIsSaving(false);
    }
  };

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return '—';
    try {
      return new Date(dateStr).toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
      });
    } catch {
      return dateStr;
    }
  };

  return (
    <div className="flex flex-col min-h-screen">
      <TopNavbar title="Личный кабинет и настройки" subtitle="Управление профилем и персонализация аккаунта" />

      <main className="p-8 max-w-4xl w-full mx-auto space-y-8 flex-1">
        {/* Navigation Tabs */}
        <div className="flex gap-2 p-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-fit shadow-xs">
          <button
            onClick={() => setActiveTab('profile')}
            className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all ${
              activeTab === 'profile'
                ? 'bg-indigo-600 text-white shadow-sm'
                : 'text-slate-600 dark:text-slate-400 hover:text-slate-900'
            }`}
          >
            <UserIcon size={16} />
            <span>Профиль пользователя</span>
          </button>
          <button
            onClick={() => setActiveTab('preferences')}
            className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all ${
              activeTab === 'preferences'
                ? 'bg-indigo-600 text-white shadow-sm'
                : 'text-slate-600 dark:text-slate-400 hover:text-slate-900'
            }`}
          >
            <Settings size={16} />
            <span>Внешний вид и настройки</span>
          </button>
        </div>

        {/* Tab 1: Profile */}
        {activeTab === 'profile' && (
          <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-8 shadow-sm space-y-8">
            {/* Header with Avatar */}
            <div className="flex flex-col sm:flex-row items-center justify-between gap-6 pb-6 border-b border-slate-100 dark:border-slate-800">
              <div className="flex items-center gap-5">
                <div className="w-20 h-20 rounded-3xl bg-gradient-to-tr from-indigo-500 to-purple-600 text-white flex items-center justify-center font-extrabold text-2xl shadow-md flex-shrink-0">
                  {avatarUrl ? (
                    <img src={avatarUrl} alt="Avatar" className="w-full h-full object-cover rounded-3xl" />
                  ) : (
                    <span>
                      {user?.first_name?.[0] || user?.username?.[0] || 'U'}
                      {user?.last_name?.[0] || ''}
                    </span>
                  )}
                </div>
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <h2 className="text-xl font-extrabold text-slate-900 dark:text-white">
                      {user?.first_name ? `${user.first_name} ${user.last_name || ''}` : user?.username}
                    </h2>
                    <span className="px-2.5 py-0.5 rounded-lg text-[10px] font-extrabold uppercase tracking-wider bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400">
                      {user?.role || 'Студент'}
                    </span>
                  </div>
                  <p className="text-xs text-slate-400">@{user?.username}</p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                {saveSuccess && (
                  <span className="text-xs text-emerald-600 font-bold flex items-center gap-1">
                    <CheckCircle2 size={16} /> Сохранено
                  </span>
                )}
                <button
                  onClick={() => setIsEditing(!isEditing)}
                  className="px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 hover:bg-slate-50 dark:hover:bg-slate-700 text-xs font-bold text-slate-700 dark:text-slate-200 transition-all flex items-center gap-1.5 shadow-xs"
                >
                  <Edit2 size={14} />
                  <span>{isEditing ? 'Отмена' : 'Редактировать'}</span>
                </button>
              </div>
            </div>

            {/* Profile Content / Form */}
            {isEditing ? (
              <form onSubmit={handleSaveProfile} className="space-y-6">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="space-y-1.5">
                    <label className="text-xs font-bold text-slate-700 dark:text-slate-300">Имя</label>
                    <input
                      type="text"
                      value={firstName}
                      onChange={(e) => setFirstName(e.target.value)}
                      className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                    />
                  </div>
                  <div className="space-y-1.5">
                    <label className="text-xs font-bold text-slate-700 dark:text-slate-300">Фамилия</label>
                    <input
                      type="text"
                      value={lastName}
                      onChange={(e) => setLastName(e.target.value)}
                      className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                    />
                  </div>
                </div>

                <div className="space-y-1.5">
                  <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                    Email <span className="text-slate-400 font-normal">(заблокирован для изменения)</span>
                  </label>
                  <input
                    type="email"
                    value={user?.email || ''}
                    disabled
                    className="w-full p-3 bg-slate-100 dark:bg-slate-800/50 border border-slate-200 dark:border-slate-800 rounded-xl text-xs text-slate-500 cursor-not-allowed"
                  />
                </div>

                <div className="space-y-1.5">
                  <label className="text-xs font-bold text-slate-700 dark:text-slate-300">Ссылка на аватарку (URL)</label>
                  <input
                    type="url"
                    value={avatarUrl}
                    onChange={(e) => setAvatarUrl(e.target.value)}
                    placeholder="https://example.com/avatar.jpg"
                    className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                  />
                </div>

                <div className="space-y-1.5">
                  <label className="text-xs font-bold text-slate-700 dark:text-slate-300">О себе (Bio)</label>
                  <textarea
                    rows={4}
                    value={bio}
                    onChange={(e) => setBio(e.target.value)}
                    placeholder="Расскажите о своем опыте и целях обучения..."
                    className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                  />
                </div>

                <div className="flex justify-end gap-3 pt-4">
                  <button
                    type="button"
                    onClick={() => setIsEditing(false)}
                    className="px-5 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-600 dark:text-slate-400 hover:bg-slate-50"
                  >
                    Отмена
                  </button>
                  <button
                    type="submit"
                    disabled={isSaving}
                    className="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold rounded-xl shadow-md transition-all flex items-center gap-2"
                  >
                    <Save size={16} />
                    <span>{isSaving ? 'Сохранение...' : 'Сохранить изменения'}</span>
                  </button>
                </div>
              </form>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="p-4 bg-slate-50 dark:bg-slate-800/60 rounded-2xl border border-slate-200 dark:border-slate-800 space-y-1">
                  <span className="text-[10px] uppercase font-bold text-slate-400">Email адрес</span>
                  <div className="flex items-center gap-2">
                    <p className="text-xs font-bold text-slate-900 dark:text-white">{user?.email}</p>
                    <span className="px-2 py-0.5 rounded text-[9px] font-extrabold bg-emerald-100 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300 flex items-center gap-0.5">
                      <ShieldCheck size={12} /> Подтвержден
                    </span>
                  </div>
                </div>

                <div className="p-4 bg-slate-50 dark:bg-slate-800/60 rounded-2xl border border-slate-200 dark:border-slate-800 space-y-1">
                  <span className="text-[10px] uppercase font-bold text-slate-400">Дата регистрации</span>
                  <div className="flex items-center gap-2">
                    <Calendar size={14} className="text-indigo-500" />
                    <p className="text-xs font-bold text-slate-900 dark:text-white">
                      {formatDate(user?.created_at)}
                    </p>
                  </div>
                </div>

                <div className="p-4 bg-slate-50 dark:bg-slate-800/60 rounded-2xl border border-slate-200 dark:border-slate-800 space-y-1 md:col-span-2">
                  <span className="text-[10px] uppercase font-bold text-slate-400">О себе</span>
                  <p className="text-xs text-slate-600 dark:text-slate-300 leading-relaxed">
                    {user?.bio || 'Информация о себе пока не заполнена.'}
                  </p>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Tab 2: Preferences */}
        {activeTab === 'preferences' && (
          <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-8 shadow-sm space-y-8">
            <div className="space-y-1">
              <h3 className="text-lg font-bold text-slate-900 dark:text-white">Тема оформления</h3>
              <p className="text-xs text-slate-500">Выберите комфортный режим отображения интерфейса</p>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <button
                onClick={() => setTheme('light')}
                className={`p-4 rounded-2xl border text-left space-y-3 transition-all ${
                  theme === 'light'
                    ? 'border-indigo-600 bg-indigo-50/50 dark:bg-indigo-950/30 ring-2 ring-indigo-500/20'
                    : 'border-slate-200 dark:border-slate-800 hover:border-indigo-300'
                }`}
              >
                <div className="w-10 h-10 rounded-xl bg-white border border-slate-200 flex items-center justify-center text-amber-500 shadow-xs">
                  <Sun size={20} />
                </div>
                <div>
                  <h4 className="text-xs font-bold text-slate-900 dark:text-white">Светлая тема</h4>
                  <p className="text-[10px] text-slate-400">Классический светлый фон</p>
                </div>
              </button>

              <button
                onClick={() => setTheme('dark')}
                className={`p-4 rounded-2xl border text-left space-y-3 transition-all ${
                  theme === 'dark'
                    ? 'border-indigo-600 bg-indigo-50/50 dark:bg-indigo-950/30 ring-2 ring-indigo-500/20'
                    : 'border-slate-200 dark:border-slate-800 hover:border-indigo-300'
                }`}
              >
                <div className="w-10 h-10 rounded-xl bg-slate-800 border border-slate-700 flex items-center justify-center text-indigo-400 shadow-xs">
                  <Moon size={20} />
                </div>
                <div>
                  <h4 className="text-xs font-bold text-slate-900 dark:text-white">Темная тема</h4>
                  <p className="text-[10px] text-slate-400">Снижает нагрузку на глаза</p>
                </div>
              </button>

              <button
                onClick={() => setTheme('system')}
                className={`p-4 rounded-2xl border text-left space-y-3 transition-all ${
                  theme === 'system'
                    ? 'border-indigo-600 bg-indigo-50/50 dark:bg-indigo-950/30 ring-2 ring-indigo-500/20'
                    : 'border-slate-200 dark:border-slate-800 hover:border-indigo-300'
                }`}
              >
                <div className="w-10 h-10 rounded-xl bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex items-center justify-center text-slate-600 dark:text-slate-300 shadow-xs">
                  <Laptop size={20} />
                </div>
                <div>
                  <h4 className="text-xs font-bold text-slate-900 dark:text-white">Системная</h4>
                  <p className="text-[10px] text-slate-400">Синхронизация с ОС</p>
                </div>
              </button>
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
