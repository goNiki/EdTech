'use client';

import React, { useState, Suspense } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/store/useAuth';
import { api } from '@/lib/api';
import {
  Mail,
  Lock,
  Eye,
  EyeOff,
  Check,
  ArrowRight,
  Loader2
} from 'lucide-react';

function LoginFormContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const login = useAuth((state) => state.login);

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [rememberMe, setRememberMe] = useState(true);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      const { data } = await api.post('/auth/login', { email, password });
      await login(data.access_token, data.refresh_token);
      
      const currentUser = useAuth.getState().user;
      const viewMode = useAuth.getState().viewMode;

      const redirectParam = searchParams.get('redirect');
      if (redirectParam && redirectParam.startsWith('/')) {
        router.replace(redirectParam);
        return;
      }

      const isStaff = ['teacher', 'author', 'admin'].includes(currentUser?.role || '');
      if (isStaff && viewMode !== 'student') {
        router.replace('/teacher/courses');
      } else {
        router.replace('/dashboard');
      }
    } catch (err: any) {
      if (err.message === 'Network Error' || err.code === 'ERR_NETWORK') {
        setError('Нет связи с сервером. Убедитесь, что бэкенд запущен.');
      } else {
        setError(err.response?.data?.message || err.response?.data?.error || 'Неверный email или пароль');
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
            Единый аккаунт экосистемы ED
          </div>

          <div className="space-y-3">
            <h2 className="text-2xl sm:text-3xl font-extrabold text-white leading-tight">
              Учитесь, создавайте курсы и проходите <span className="text-cyan-300">тестирования</span>
            </h2>
            <p className="text-indigo-100 text-sm leading-relaxed font-normal">
              Автоматическая проверка ответов, срезы знаний, персональная статистика и сертификаты сразу после прохождения.
            </p>
          </div>

          {/* Mini Live Test Scorecard Preview */}
          <div className="p-4 rounded-2xl bg-white/10 backdrop-blur-md border border-white/15 space-y-3">
            <div className="flex items-center justify-between text-xs text-indigo-100">
              <span className="font-semibold text-white">Срез знаний</span>
              <span className="text-emerald-300 font-medium">100% автопроверка</span>
            </div>
            <div className="space-y-2">
              <div className="flex justify-between text-xs text-indigo-100 font-medium">
                <span>Точность ответов</span>
                <span className="text-white font-bold">92%</span>
              </div>
              <div className="w-full bg-black/20 h-2 rounded-full overflow-hidden">
                <div className="bg-gradient-to-r from-cyan-300 to-emerald-400 h-full w-[92%] rounded-full" />
              </div>
            </div>
          </div>
        </div>

        {/* Features inside Left Panel */}
        <div className="relative z-10 pt-6 border-t border-white/15 grid grid-cols-2 gap-4 text-xs text-indigo-100 font-medium">
          <div className="flex items-center gap-2">
            <Check size={16} className="text-emerald-300 flex-shrink-0" />
            <span>Мгновенный доступ</span>
          </div>
          <div className="flex items-center gap-2">
            <Check size={16} className="text-emerald-300 flex-shrink-0" />
            <span>Защита данных</span>
          </div>
          <div className="flex items-center gap-2">
            <Check size={16} className="text-emerald-300 flex-shrink-0" />
            <span>600+ тестов</span>
          </div>
          <div className="flex items-center gap-2">
            <Check size={16} className="text-emerald-300 flex-shrink-0" />
            <span>Умные отчеты</span>
          </div>
        </div>
      </div>

      {/* Right Form Section */}
      <div className="lg:col-span-7 flex flex-col justify-center px-2 sm:px-6">
        {/* Mode Toggle Tabs */}
        <div className="bg-slate-100 dark:bg-slate-800 p-1.5 rounded-2xl border border-slate-200/80 dark:border-slate-700 flex items-center mb-8">
          <Link
            href="/login"
            className="flex-1 py-3 text-center text-sm font-bold rounded-xl transition-all bg-white dark:bg-indigo-600 text-indigo-600 dark:text-white shadow-xs"
          >
            Вход
          </Link>
          <Link
            href="/register"
            className="flex-1 py-3 text-center text-sm font-bold rounded-xl transition-all text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white"
          >
            Регистрация
          </Link>
        </div>

        {/* Heading */}
        <div className="mb-6 space-y-1">
          <h1 className="text-2xl sm:text-3xl font-extrabold text-slate-900 dark:text-white tracking-tight">
            Войти в аккаунт
          </h1>
          <p className="text-sm text-slate-600 dark:text-slate-400">
            Введите ваш email и пароль для доступа к тестам и курсам
          </p>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <div className="p-3.5 bg-rose-50 dark:bg-rose-950/60 border border-rose-200 dark:border-rose-800 rounded-xl text-xs font-semibold text-rose-600 dark:text-rose-400">
              {error}
            </div>
          )}

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
            <div className="flex items-center justify-between">
              <label className="block text-xs font-semibold text-slate-700 dark:text-slate-300">
                Пароль
              </label>
              <span className="text-xs text-indigo-600 dark:text-indigo-400 hover:text-indigo-700 font-semibold cursor-pointer">
                Забыли пароль?
              </span>
            </div>
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

          {/* Remember Me */}
          <div className="flex items-center justify-between pt-1">
            <label className="flex items-center gap-2 cursor-pointer text-xs text-slate-600 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200">
              <input
                type="checkbox"
                checked={rememberMe}
                onChange={(e) => setRememberMe(e.target.checked)}
                className="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500/20 cursor-pointer"
              />
              <span>Запомнить меня на этом устройстве</span>
            </label>
          </div>

          {/* Submit Button */}
          <button
            type="submit"
            disabled={isLoading}
            className="w-full py-3.5 bg-indigo-600 hover:bg-indigo-700 active:scale-95 text-white font-bold rounded-xl shadow-lg shadow-indigo-600/30 transition-all text-sm mt-2 flex items-center justify-center gap-2 cursor-pointer"
          >
            <span>{isLoading ? 'Проверка данных...' : 'Войти в личный кабинет'}</span>
            {!isLoading && <ArrowRight size={16} />}
          </button>
        </form>

        {/* Footer info */}
        <div className="mt-8 text-center text-xs text-slate-500 dark:text-slate-400">
          Еще нет аккаунта?{' '}
          <Link href="/register" className="text-indigo-600 dark:text-indigo-400 font-bold hover:underline">
            Зарегистрироваться бесплатно
          </Link>
        </div>
      </div>
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense
      fallback={
        <div className="w-full max-w-5xl h-96 flex items-center justify-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800">
          <Loader2 className="w-8 h-8 animate-spin text-indigo-600" />
        </div>
      }
    >
      <LoginFormContent />
    </Suspense>
  );
}
