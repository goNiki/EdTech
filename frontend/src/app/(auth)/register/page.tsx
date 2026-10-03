'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/store/useAuth';
import { api } from '@/lib/api';
import {
  Mail,
  Lock,
  User as UserIcon,
  Eye,
  EyeOff,
  Check,
  ArrowRight
} from 'lucide-react';

export default function RegisterPage() {
  const router = useRouter();
  const login = useAuth((state) => state.login);

  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [termsAgreed, setTermsAgreed] = useState(true);
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!termsAgreed) {
      setError('Необходимо согласиться с правилами платформы.');
      return;
    }

    setError('');
    setIsLoading(true);

    try {
      // 1. Register
      await api.post('/auth/register', {
        email,
        password,
        username: username.trim() || email.split('@')[0],
      });

      // 2. Login
      const { data } = await api.post('/auth/login', { email, password });
      login(data.access_token, data.refresh_token);
      await useAuth.getState().fetchUser();

      const currentUser = useAuth.getState().user;
      if (currentUser?.role === 'teacher') {
        router.push('/teacher/courses');
      } else {
        router.push('/dashboard');
      }
    } catch (err: any) {
      if (err.message === 'Network Error' || err.code === 'ERR_NETWORK') {
        setError('Нет связи с сервером. Убедитесь, что бэкенд запущен.');
      } else {
        setError(
          err.response?.data?.message ||
            err.response?.data?.error ||
            'Ошибка при регистрации. Проверьте введенные данные.'
        );
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full max-w-5xl grid lg:grid-cols-12 gap-8 items-stretch bg-white dark:bg-slate-900 border border-slate-200/90 dark:border-slate-800 rounded-3xl p-4 sm:p-8 lg:p-10 shadow-xl shadow-indigo-500/5">
      {/* Left Promotional Info Panel (Desktop Only) */}
      <div className="hidden lg:flex lg:col-span-5 flex-col justify-between p-8 rounded-2xl bg-gradient-to-br from-indigo-600 via-indigo-700 to-indigo-900 text-white relative overflow-hidden shadow-lg shadow-indigo-600/20">
        {/* Decorative blur elements */}
        <div className="absolute -top-16 -right-16 w-52 h-52 bg-cyan-400/20 rounded-full blur-2xl pointer-events-none" />
        <div className="absolute -bottom-16 -left-16 w-52 h-52 bg-indigo-400/25 rounded-full blur-2xl pointer-events-none" />

        <div className="relative z-10 space-y-6">
          <div className="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-white/15 backdrop-blur-md border border-white/20 text-white text-xs font-semibold">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
            Быстрый старт за 30 секунд
          </div>

          <div className="space-y-3">
            <h2 className="text-2xl sm:text-3xl font-extrabold text-white leading-tight">
              Создайте профиль и начните <span className="text-cyan-300">обучение</span>
            </h2>
            <p className="text-indigo-100 text-sm leading-relaxed font-normal">
              Персональный кабинет ученика, неограниченный доступ к срезам знаний и конструктору интерактивных тестов.
            </p>
          </div>

          {/* Mini Live Test Scorecard Preview */}
          <div className="p-4 rounded-2xl bg-white/10 backdrop-blur-md border border-white/15 space-y-3">
            <div className="flex items-center justify-between text-xs text-indigo-100">
              <span className="font-semibold text-white">Уровень доступа</span>
              <span className="text-emerald-300 font-medium">Бесплатный аккаунт</span>
            </div>
            <div className="space-y-2">
              <div className="flex justify-between text-xs text-indigo-100 font-medium">
                <span>Доступно курсов</span>
                <span className="text-white font-bold">600+ программ</span>
              </div>
              <div className="w-full bg-black/20 h-2 rounded-full overflow-hidden">
                <div className="bg-gradient-to-r from-cyan-300 to-emerald-400 h-full w-[100%] rounded-full" />
              </div>
            </div>
          </div>
        </div>

        {/* Features inside Left Panel */}
        <div className="relative z-10 pt-6 border-t border-white/15 grid grid-cols-2 gap-4 text-xs text-indigo-100 font-medium">
          <div className="flex items-center gap-2">
            <Check size={16} className="text-emerald-300 flex-shrink-0" />
            <span>Сохранение прогресса</span>
          </div>
          <div className="flex items-center gap-2">
            <Check size={16} className="text-emerald-300 flex-shrink-0" />
            <span>Личная аналитика</span>
          </div>
          <div className="flex items-center gap-2">
            <Check size={16} className="text-emerald-300 flex-shrink-0" />
            <span>Именные сертификаты</span>
          </div>
          <div className="flex items-center gap-2">
            <Check size={16} className="text-emerald-300 flex-shrink-0" />
            <span>Поддержка 24/7</span>
          </div>
        </div>
      </div>

      {/* Right Form Section */}
      <div className="lg:col-span-7 flex flex-col justify-center px-2 sm:px-6">
        {/* Mode Toggle Tabs */}
        <div className="bg-slate-100 dark:bg-slate-800 p-1.5 rounded-2xl border border-slate-200/80 dark:border-slate-700 flex items-center mb-8">
          <Link
            href="/login"
            className="flex-1 py-3 text-center text-sm font-bold rounded-xl transition-all text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white"
          >
            Вход
          </Link>
          <Link
            href="/register"
            className="flex-1 py-3 text-center text-sm font-bold rounded-xl transition-all bg-white dark:bg-indigo-600 text-indigo-600 dark:text-white shadow-xs"
          >
            Регистрация
          </Link>
        </div>

        {/* Heading */}
        <div className="mb-6 space-y-1">
          <h1 className="text-2xl sm:text-3xl font-extrabold text-slate-900 dark:text-white tracking-tight">
            Создать профиль
          </h1>
          <p className="text-sm text-slate-600 dark:text-slate-400">
            Заполните данные для создания аккаунта на платформе
          </p>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <div className="p-3.5 bg-rose-50 dark:bg-rose-950/60 border border-rose-200 dark:border-rose-800 rounded-xl text-xs font-semibold text-rose-600 dark:text-rose-400">
              {error}
            </div>
          )}

          {/* Username */}
          <div className="space-y-1.5">
            <label className="block text-xs font-semibold text-slate-700 dark:text-slate-300">
              Имя пользователя (Никнейм)
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                <UserIcon size={16} />
              </div>
              <input
                type="text"
                required
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="alex_developer"
                className="w-full pl-10 pr-4 py-3 bg-slate-50 dark:bg-slate-800 hover:bg-slate-50/80 focus:bg-white dark:focus:bg-slate-900 border border-slate-200 dark:border-slate-700 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 rounded-xl text-sm text-slate-900 dark:text-white placeholder-slate-400 transition-all outline-none"
              />
            </div>
          </div>

          {/* Email */}
          <div className="space-y-1.5">
            <label className="block text-xs font-semibold text-slate-700 dark:text-slate-300">
              Электронная почта (Email)
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                <Mail size={16} />
              </div>
              <input
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="name@example.com"
                className="w-full pl-10 pr-4 py-3 bg-slate-50 dark:bg-slate-800 hover:bg-slate-50/80 focus:bg-white dark:focus:bg-slate-900 border border-slate-200 dark:border-slate-700 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 rounded-xl text-sm text-slate-900 dark:text-white placeholder-slate-400 transition-all outline-none"
              />
            </div>
          </div>

          {/* Password */}
          <div className="space-y-1.5">
            <label className="block text-xs font-semibold text-slate-700 dark:text-slate-300">
              Пароль
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                <Lock size={16} />
              </div>
              <input
                type={showPassword ? 'text' : 'password'}
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                className="w-full pl-10 pr-11 py-3 bg-slate-50 dark:bg-slate-800 hover:bg-slate-50/80 focus:bg-white dark:focus:bg-slate-900 border border-slate-200 dark:border-slate-700 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 rounded-xl text-sm text-slate-900 dark:text-white placeholder-slate-400 transition-all outline-none"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 focus:outline-none cursor-pointer"
              >
                {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
              </button>
            </div>
          </div>

          {/* Terms Agreement */}
          <div className="flex items-start gap-2 pt-1">
            <input
              type="checkbox"
              id="termsAgree"
              checked={termsAgreed}
              onChange={(e) => setTermsAgreed(e.target.checked)}
              className="mt-1 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500/20 cursor-pointer"
            />
            <label htmlFor="termsAgree" className="text-xs text-slate-600 dark:text-slate-400 cursor-pointer">
              Я соглашаюсь с правилами платформы и обработкой персональных данных
            </label>
          </div>

          {/* Submit Button */}
          <button
            type="submit"
            disabled={isLoading}
            className="w-full py-3.5 bg-indigo-600 hover:bg-indigo-700 active:scale-95 text-white font-bold rounded-xl shadow-lg shadow-indigo-600/30 transition-all text-sm mt-2 flex items-center justify-center gap-2 cursor-pointer"
          >
            <span>{isLoading ? 'Создание профиля...' : 'Зарегистрироваться бесплатно'}</span>
            {!isLoading && <ArrowRight size={16} />}
          </button>
        </form>

        {/* Footer info */}
        <div className="mt-8 text-center text-xs text-slate-500 dark:text-slate-400">
          Уже есть аккаунт?{' '}
          <Link href="/login" className="text-indigo-600 dark:text-indigo-400 font-bold hover:underline">
            Войти в систему
          </Link>
        </div>
      </div>
    </div>
  );
}
