'use client';

import React, { useState, useEffect, useRef } from 'react';
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
  CheckCircle2,
  Camera,
  Upload,
  Loader2,
  Trash2,
  Lock,
  KeyRound,
  Eye,
  EyeOff,
  ShieldAlert,
  AlertCircle
} from 'lucide-react';

export default function ProfileAndSettingsPage() {
  const { user, setUser } = useAuth();
  const { theme, setTheme } = useTheme();

  const [activeTab, setActiveTab] = useState<'profile' | 'security' | 'preferences'>('profile');
  const [isEditing, setIsEditing] = useState(false);

  // Form states
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [headline, setHeadline] = useState('');
  const [bio, setBio] = useState('');
  const [avatarUrl, setAvatarUrl] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [saveSuccess, setSaveSuccess] = useState(false);
  const [isUploadingAvatar, setIsUploadingAvatar] = useState(false);
  const [avatarError, setAvatarError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Change password states
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showOldPassword, setShowOldPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isChangingPassword, setIsChangingPassword] = useState(false);
  const [passwordError, setPasswordError] = useState<string | null>(null);
  const [passwordSuccess, setPasswordSuccess] = useState(false);

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setPasswordError(null);
    setPasswordSuccess(false);

    if (!oldPassword) {
      setPasswordError('Пожалуйста, введите текущий пароль');
      return;
    }

    if (newPassword.length < 8) {
      setPasswordError('Новый пароль должен содержать не менее 8 символов');
      return;
    }

    if (newPassword === oldPassword) {
      setPasswordError('Новый пароль не должен совпадать с текущим');
      return;
    }

    if (newPassword !== confirmPassword) {
      setPasswordError('Новый пароль и подтверждение не совпадают');
      return;
    }

    setIsChangingPassword(true);
    try {
      await api.post('/auth/change-password', {
        old_password: oldPassword,
        new_password: newPassword,
      });

      setOldPassword('');
      setNewPassword('');
      setConfirmPassword('');
      setPasswordSuccess(true);
      setTimeout(() => setPasswordSuccess(false), 5000);
    } catch (err: any) {
      console.error('Change password failed', err);
      const errMsg =
        err.response?.data?.message ||
        err.response?.data?.error ||
        'Не удалось изменить пароль. Проверьте правильность текущего пароля.';
      setPasswordError(errMsg);
    } finally {
      setIsChangingPassword(false);
    }
  };

  const resolveUrl = (url: string) => {
    if (!url) return '';
    if (url.startsWith('http://') || url.startsWith('https://')) return url;
    const baseUrl = (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8082/api/v1').replace(/\/api\/v1\/?$/, '');
    return `${baseUrl}${url.startsWith('/') ? '' : '/'}${url}`;
  };

  const handleAvatarFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (file.size > 5 * 1024 * 1024) {
      setAvatarError('Размер фото не должен превышать 5 МБ');
      return;
    }

    setIsUploadingAvatar(true);
    setAvatarError(null);

    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('category', 'avatar');

      const res = await api.post('/upload', formData, {
        params: { category: 'avatar' },
        headers: { 'Content-Type': 'multipart/form-data' },
      });

      const data = res.data.data || res.data;
      const uploadedUrl = data.file_url || data.url || data.FileUrl;
      if (uploadedUrl) {
        const fullUrl = resolveUrl(uploadedUrl);
        setAvatarUrl(fullUrl);
      }
    } catch (err: any) {
      console.error('Avatar upload failed', err);
      setAvatarError(err.response?.data?.message || err.response?.data?.error || 'Не удалось загрузить фото');
    } finally {
      setIsUploadingAvatar(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  useEffect(() => {
    if (user) {
      setFirstName(user.first_name || '');
      setLastName(user.last_name || '');
      setHeadline(user.headline || '');
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
        headline: headline,
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
        <div className="flex flex-wrap gap-2 p-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl w-fit shadow-xs">
          <button
            onClick={() => setActiveTab('profile')}
            className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer ${
              activeTab === 'profile'
                ? 'bg-indigo-600 text-white shadow-sm'
                : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'
            }`}
          >
            <UserIcon size={16} />
            <span>Профиль пользователя</span>
          </button>
          <button
            onClick={() => setActiveTab('security')}
            className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer ${
              activeTab === 'security'
                ? 'bg-indigo-600 text-white shadow-sm'
                : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'
            }`}
          >
            <Lock size={16} />
            <span>Безопасность и пароль</span>
          </button>
          <button
            onClick={() => setActiveTab('preferences')}
            className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer ${
              activeTab === 'preferences'
                ? 'bg-indigo-600 text-white shadow-sm'
                : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'
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
                <div
                  className="relative group w-20 h-20 rounded-3xl bg-gradient-to-tr from-indigo-500 to-purple-600 text-white flex items-center justify-center font-extrabold text-2xl shadow-md flex-shrink-0 overflow-hidden cursor-pointer"
                  onClick={() => fileInputRef.current?.click()}
                  title="Нажмите, чтобы загрузить новое фото"
                >
                  {avatarUrl ? (
                    <img src={avatarUrl} alt="Avatar" className="w-full h-full object-cover rounded-3xl" />
                  ) : (
                    <span>
                      {user?.first_name?.[0] || user?.username?.[0] || 'U'}
                      {user?.last_name?.[0] || ''}
                    </span>
                  )}

                  {/* Hover / Upload overlay */}
                  <div className="absolute inset-0 bg-black/55 text-white flex flex-col items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity rounded-3xl">
                    {isUploadingAvatar ? (
                      <Loader2 size={20} className="animate-spin text-white" />
                    ) : (
                      <>
                        <Camera size={18} />
                        <span className="text-[9px] font-bold mt-1">Фото</span>
                      </>
                    )}
                  </div>
                </div>

                <input
                  type="file"
                  ref={fileInputRef}
                  onChange={handleAvatarFileChange}
                  accept="image/png,image/jpeg,image/webp,image/gif"
                  className="hidden"
                />

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
                  <div className="flex items-center justify-between">
                    <label className="text-xs font-bold text-slate-700 dark:text-slate-300">Аватар профиля</label>
                    <button
                      type="button"
                      onClick={() => fileInputRef.current?.click()}
                      disabled={isUploadingAvatar}
                      className="text-xs font-bold text-indigo-600 hover:text-indigo-700 dark:text-indigo-400 flex items-center gap-1.5"
                    >
                      {isUploadingAvatar ? (
                        <>
                          <Loader2 size={13} className="animate-spin" />
                          <span>Загрузка...</span>
                        </>
                      ) : (
                        <>
                          <Upload size={13} />
                          <span>Загрузить фото с диска</span>
                        </>
                      )}
                    </button>
                  </div>
                  <div className="flex gap-2">
                    <input
                      type="url"
                      value={avatarUrl}
                      onChange={(e) => setAvatarUrl(e.target.value)}
                      placeholder="https://example.com/avatar.jpg или выберите файл"
                      className="flex-1 p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                    />
                    {avatarUrl && (
                      <button
                        type="button"
                        onClick={() => setAvatarUrl('')}
                        className="px-3 py-2 rounded-xl border border-rose-200 dark:border-rose-900/50 text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors flex items-center gap-1 text-xs font-bold"
                        title="Удалить фото"
                      >
                        <Trash2 size={14} />
                      </button>
                    )}
                  </div>
                  {avatarError && (
                    <p className="text-[11px] text-rose-600 dark:text-rose-400 font-semibold">{avatarError}</p>
                  )}
                </div>

                <div className="space-y-1.5">
                  <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                    Специализация / Профессиональный статус (Headline)
                  </label>
                  <input
                    type="text"
                    value={headline}
                    onChange={(e) => setHeadline(e.target.value)}
                    placeholder="Например: Эксперт ЕГЭ по русскому языку, стаж 10 лет"
                    maxLength={150}
                    className="w-full p-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                  />
                  <p className="text-[10px] text-slate-400">
                    Отображается в карточке автора на лендингах ваших курсов
                  </p>
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

                {user?.headline && (
                  <div className="p-4 bg-slate-50 dark:bg-slate-800/60 rounded-2xl border border-slate-200 dark:border-slate-800 space-y-1 md:col-span-2">
                    <span className="text-[10px] uppercase font-bold text-slate-400">Специализация / Статус автора</span>
                    <p className="text-xs font-bold text-indigo-600 dark:text-indigo-400">
                      {user.headline}
                    </p>
                  </div>
                )}

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

        {/* Tab 2: Security & Password */}
        {activeTab === 'security' && (
          <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-8 shadow-sm space-y-8">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-slate-100 dark:border-slate-800">
              <div className="space-y-1">
                <div className="flex items-center gap-2">
                  <h3 className="text-lg font-extrabold text-slate-900 dark:text-white">
                    Безопасность и смена пароля
                  </h3>
                  <span className="px-2.5 py-0.5 rounded-lg text-[10px] font-extrabold uppercase bg-emerald-50 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300 flex items-center gap-1">
                    <ShieldCheck size={12} /> Защищено
                  </span>
                </div>
                <p className="text-xs text-slate-500">
                  Обновите пароль для защиты вашего аккаунта и доступа к курсам
                </p>
              </div>

              {passwordSuccess && (
                <div className="px-4 py-2 rounded-xl bg-emerald-50 dark:bg-emerald-950/80 border border-emerald-200 dark:border-emerald-800 text-emerald-700 dark:text-emerald-300 text-xs font-bold flex items-center gap-1.5 animate-in fade-in">
                  <CheckCircle2 size={16} />
                  <span>Пароль успешно обновлен!</span>
                </div>
              )}
            </div>

            <form onSubmit={handleChangePassword} className="max-w-xl space-y-5">
              {passwordError && (
                <div className="p-4 rounded-2xl bg-rose-50 dark:bg-rose-950/60 border border-rose-200 dark:border-rose-900 text-rose-600 dark:text-rose-400 text-xs font-medium flex items-start gap-2.5">
                  <AlertCircle size={16} className="mt-0.5 flex-shrink-0" />
                  <span>{passwordError}</span>
                </div>
              )}

              {/* Current Password */}
              <div className="space-y-1.5">
                <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                  Текущий пароль *
                </label>
                <div className="relative">
                  <input
                    type={showOldPassword ? 'text' : 'password'}
                    required
                    value={oldPassword}
                    onChange={(e) => setOldPassword(e.target.value)}
                    placeholder="Введите ваш текущий пароль"
                    className="w-full p-3 pr-10 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                  />
                  <button
                    type="button"
                    onClick={() => setShowOldPassword(!showOldPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 cursor-pointer"
                    title={showOldPassword ? 'Скрыть пароль' : 'Показать пароль'}
                  >
                    {showOldPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </button>
                </div>
              </div>

              {/* New Password */}
              <div className="space-y-1.5">
                <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                  Новый пароль *
                </label>
                <div className="relative">
                  <input
                    type={showNewPassword ? 'text' : 'password'}
                    required
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    placeholder="Минимум 8 символов"
                    className="w-full p-3 pr-10 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                  />
                  <button
                    type="button"
                    onClick={() => setShowNewPassword(!showNewPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 cursor-pointer"
                    title={showNewPassword ? 'Скрыть пароль' : 'Показать пароль'}
                  >
                    {showNewPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </button>
                </div>

                {/* Password strength indicator */}
                {newPassword && (
                  <div className="space-y-1 pt-1">
                    <div className="flex items-center justify-between text-[10px]">
                      <span className="font-semibold text-slate-400">Надежность пароля:</span>
                      <span className={`font-bold ${
                        newPassword.length >= 10 && /[A-Z]/.test(newPassword) && /[0-9]/.test(newPassword)
                          ? 'text-emerald-600'
                          : newPassword.length >= 8
                          ? 'text-amber-500'
                          : 'text-rose-500'
                      }`}>
                        {newPassword.length >= 10 && /[A-Z]/.test(newPassword) && /[0-9]/.test(newPassword)
                          ? 'Надежный'
                          : newPassword.length >= 8
                          ? 'Средний'
                          : 'Слабый (менее 8 символов)'}
                      </span>
                    </div>
                    <div className="h-1.5 w-full bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden flex gap-1">
                      <div className={`h-full rounded-full transition-all duration-300 ${
                        newPassword.length >= 8 ? 'bg-amber-500 flex-1' : 'bg-rose-500 w-1/3'
                      }`} />
                      <div className={`h-full rounded-full transition-all duration-300 flex-1 ${
                        newPassword.length >= 10 && (/[0-9]/.test(newPassword) || /[A-Z]/.test(newPassword))
                          ? 'bg-amber-500'
                          : 'bg-transparent'
                      }`} />
                      <div className={`h-full rounded-full transition-all duration-300 flex-1 ${
                        newPassword.length >= 10 && /[A-Z]/.test(newPassword) && /[0-9]/.test(newPassword)
                          ? 'bg-emerald-500'
                          : 'bg-transparent'
                      }`} />
                    </div>
                  </div>
                )}
              </div>

              {/* Confirm New Password */}
              <div className="space-y-1.5">
                <label className="text-xs font-bold text-slate-700 dark:text-slate-300">
                  Повторите новый пароль *
                </label>
                <div className="relative">
                  <input
                    type={showConfirmPassword ? 'text' : 'password'}
                    required
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder="Повторите новый пароль"
                    className="w-full p-3 pr-10 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-xs font-medium focus:ring-2 focus:ring-indigo-500 focus:outline-none"
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 cursor-pointer"
                    title={showConfirmPassword ? 'Скрыть пароль' : 'Показать пароль'}
                  >
                    {showConfirmPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </button>
                </div>
                {confirmPassword && newPassword !== confirmPassword && (
                  <p className="text-[11px] text-rose-500 font-semibold">
                    Пароли не совпадают
                  </p>
                )}
              </div>

              {/* Security notice */}
              <div className="p-4 bg-slate-50 dark:bg-slate-800/60 rounded-2xl border border-slate-200 dark:border-slate-800 flex items-start gap-3 text-xs text-slate-500 dark:text-slate-400">
                <ShieldAlert size={18} className="text-amber-500 flex-shrink-0 mt-0.5" />
                <p>
                  После успешной смены пароля вам потребуется заново войти в аккаунт на всех других устройствах.
                </p>
              </div>

              <div className="pt-2">
                <button
                  type="submit"
                  disabled={isChangingPassword || (confirmPassword ? newPassword !== confirmPassword : false)}
                  className="px-6 py-3 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed text-white text-xs font-bold rounded-xl shadow-md transition-all flex items-center gap-2 cursor-pointer"
                >
                  {isChangingPassword ? (
                    <>
                      <Loader2 size={16} className="animate-spin" />
                      <span>Обновление пароля...</span>
                    </>
                  ) : (
                    <>
                      <KeyRound size={16} />
                      <span>Обновить пароль</span>
                    </>
                  )}
                </button>
              </div>
            </form>
          </div>
        )}

        {/* Tab 3: Preferences */}
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
