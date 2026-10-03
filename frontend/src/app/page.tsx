'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/store/useAuth';
import { api } from '@/lib/api';
import {
  Sparkles,
  CheckCircle2,
  XCircle,
  ArrowRight,
  BookOpen,
  Layers,
  Award,
  BarChart3,
  ShieldCheck,
  Zap,
  Clock,
  Check,
  X,
  Menu,
  User as UserIcon,
  Mail,
  Lock,
  Eye,
  EyeOff,
  ChevronRight,
  FileSpreadsheet
} from 'lucide-react';

export default function HomePage() {
  const router = useRouter();
  const { user, isAuthenticated, login } = useAuth();

  // Navigation & Drawer
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  // How it works role toggle
  const [activeRoleTab, setActiveRoleTab] = useState<'student' | 'teacher'>('student');

  // Hero interactive quiz state
  const [heroSelectedOption, setHeroSelectedOption] = useState<number | null>(null);

  // Auth Modal state
  const [isAuthModalOpen, setIsAuthModalOpen] = useState(false);
  const [authMode, setAuthMode] = useState<'login' | 'register'>('login');
  const [authEmail, setAuthEmail] = useState('');
  const [authPassword, setAuthPassword] = useState('');
  const [authUsername, setAuthUsername] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [authError, setAuthError] = useState('');
  const [authLoading, setAuthLoading] = useState(false);
  const [authNotice, setAuthNotice] = useState<string | null>(null);

  const openAuth = (mode: 'login' | 'register') => {
    setAuthMode(mode);
    setAuthError('');
    setAuthNotice(null);
    setIsAuthModalOpen(true);
  };

  const handleAuthSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setAuthError('');
    setAuthLoading(true);

    try {
      if (authMode === 'login') {
        const { data } = await api.post('/auth/login', {
          email: authEmail,
          password: authPassword,
        });
        login(data.access_token, data.refresh_token);
        await useAuth.getState().fetchUser();
        setAuthNotice(`Вход выполнен успешно! Добро пожаловать, ${authEmail}`);
        setTimeout(() => {
          setIsAuthModalOpen(false);
          router.push('/dashboard');
        }, 1000);
      } else {
        await api.post('/auth/register', {
          email: authEmail,
          password: authPassword,
          username: authUsername || authEmail.split('@')[0],
        });
        // Auto-login after registration
        const { data } = await api.post('/auth/login', {
          email: authEmail,
          password: authPassword,
        });
        login(data.access_token, data.refresh_token);
        await useAuth.getState().fetchUser();
        setAuthNotice(`Регистрация успешна! Добро пожаловать, ${authEmail}`);
        setTimeout(() => {
          setIsAuthModalOpen(false);
          router.push('/dashboard');
        }, 1000);
      }
    } catch (err: any) {
      if (err.message === 'Network Error' || err.code === 'ERR_NETWORK') {
        setAuthError('Нет связи с сервером. Убедитесь, что бэкенд запущен.');
      } else {
        setAuthError(err.response?.data?.message || err.response?.data?.error || 'Ошибка авторизации. Проверьте данные.');
      }
    } finally {
      setAuthLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-slate-950 text-slate-800 dark:text-slate-100 antialiased selection:bg-indigo-500 selection:text-white flex flex-col font-sans transition-colors duration-200">
      {/* Sticky Top Navbar */}
      <header className="sticky top-0 z-50 bg-white/95 dark:bg-slate-900/95 backdrop-blur-md border-b border-slate-200/80 dark:border-slate-800 transition-all shadow-xs">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-20">
            {/* Logo */}
            <Link href="/" className="flex items-center gap-3 group">
              <div className="w-11 h-11 rounded-2xl bg-gradient-to-tr from-indigo-600 via-indigo-500 to-cyan-400 flex items-center justify-center text-white font-extrabold text-xl shadow-lg shadow-indigo-500/25 group-hover:scale-105 transition-transform">
                ED
              </div>
              <div className="flex flex-col">
                <span className="font-black text-xl tracking-tight text-slate-900 dark:text-white group-hover:text-indigo-600 transition-colors">
                  ED.Learn
                </span>
                <span className="text-[10px] uppercase font-bold tracking-widest text-indigo-500">
                  Курсы & Тестирование
                </span>
              </div>
            </Link>

            {/* Desktop Navigation Links */}
            <nav className="hidden md:flex items-center gap-1 lg:gap-2">
              <a
                href="#about"
                className="px-3.5 py-2 text-sm font-medium text-slate-600 dark:text-slate-300 hover:text-indigo-600 hover:bg-indigo-50/60 dark:hover:bg-slate-800 rounded-xl transition-colors"
              >
                О платформе
              </a>
              <a
                href="#advantages"
                className="px-3.5 py-2 text-sm font-medium text-slate-600 dark:text-slate-300 hover:text-indigo-600 hover:bg-indigo-50/60 dark:hover:bg-slate-800 rounded-xl transition-colors"
              >
                Преимущества
              </a>
              <a
                href="#how-it-works"
                className="px-3.5 py-2 text-sm font-medium text-slate-600 dark:text-slate-300 hover:text-indigo-600 hover:bg-indigo-50/60 dark:hover:bg-slate-800 rounded-xl transition-colors"
              >
                Как пользоваться
              </a>
              <a
                href="#courses"
                className="px-3.5 py-2 text-sm font-medium text-slate-600 dark:text-slate-300 hover:text-indigo-600 hover:bg-indigo-50/60 dark:hover:bg-slate-800 rounded-xl transition-colors"
              >
                Курсы и тесты
              </a>
              <a
                href="#contacts"
                className="px-3.5 py-2 text-sm font-medium text-slate-600 dark:text-slate-300 hover:text-indigo-600 hover:bg-indigo-50/60 dark:hover:bg-slate-800 rounded-xl transition-colors"
              >
                Контакты
              </a>
            </nav>

            {/* Auth Actions */}
            <div className="hidden sm:flex items-center gap-3">
              {isAuthenticated ? (
                <div className="flex items-center gap-3">
                  <Link
                    href="/dashboard"
                    className="px-5 py-2.5 text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 rounded-xl shadow-md shadow-indigo-600/20 transition-all flex items-center gap-2"
                  >
                    <UserIcon size={16} />
                    <span>Личный кабинет</span>
                  </Link>
                  {user?.role === 'teacher' && (
                    <Link
                      href="/teacher/courses"
                      className="px-4 py-2.5 text-sm font-bold text-slate-700 dark:text-slate-200 bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 rounded-xl transition-all"
                    >
                      Преподавание
                    </Link>
                  )}
                </div>
              ) : (
                <>
                  <button
                    onClick={() => openAuth('login')}
                    className="px-4 py-2.5 text-sm font-semibold text-slate-700 dark:text-slate-200 hover:text-indigo-600 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-xl transition-all cursor-pointer"
                  >
                    Войти
                  </button>
                  <button
                    onClick={() => openAuth('register')}
                    className="px-5 py-2.5 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-700 active:scale-95 shadow-md shadow-indigo-600/30 rounded-xl transition-all flex items-center gap-2 cursor-pointer"
                  >
                    <span>Регистрация</span>
                    <ArrowRight size={16} />
                  </button>
                </>
              )}
            </div>

            {/* Mobile menu toggle */}
            <div className="md:hidden flex items-center">
              <button
                onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                className="p-2.5 rounded-xl text-slate-600 dark:text-slate-300 hover:text-slate-900 hover:bg-slate-100 dark:hover:bg-slate-800"
              >
                <Menu size={24} />
              </button>
            </div>
          </div>
        </div>

        {/* Mobile Drawer */}
        {isMobileMenuOpen && (
          <div className="md:hidden bg-white/95 dark:bg-slate-900/95 border-b border-slate-200 dark:border-slate-800 px-4 pt-2 pb-6 space-y-2 backdrop-blur-md">
            <a
              href="#about"
              onClick={() => setIsMobileMenuOpen(false)}
              className="block px-4 py-3 rounded-xl text-base font-medium text-slate-700 dark:text-slate-300 hover:bg-indigo-50 dark:hover:bg-slate-800 hover:text-indigo-600"
            >
              О платформе
            </a>
            <a
              href="#advantages"
              onClick={() => setIsMobileMenuOpen(false)}
              className="block px-4 py-3 rounded-xl text-base font-medium text-slate-700 dark:text-slate-300 hover:bg-indigo-50 dark:hover:bg-slate-800 hover:text-indigo-600"
            >
              Преимущества
            </a>
            <a
              href="#how-it-works"
              onClick={() => setIsMobileMenuOpen(false)}
              className="block px-4 py-3 rounded-xl text-base font-medium text-slate-700 dark:text-slate-300 hover:bg-indigo-50 dark:hover:bg-slate-800 hover:text-indigo-600"
            >
              Как пользоваться
            </a>
            <a
              href="#courses"
              onClick={() => setIsMobileMenuOpen(false)}
              className="block px-4 py-3 rounded-xl text-base font-medium text-slate-700 dark:text-slate-300 hover:bg-indigo-50 dark:hover:bg-slate-800 hover:text-indigo-600"
            >
              Курсы и тесты
            </a>
            <a
              href="#contacts"
              onClick={() => setIsMobileMenuOpen(false)}
              className="block px-4 py-3 rounded-xl text-base font-medium text-slate-700 dark:text-slate-300 hover:bg-indigo-50 dark:hover:bg-slate-800 hover:text-indigo-600"
            >
              Контакты
            </a>
            <div className="pt-4 flex flex-col gap-2.5 border-t border-slate-100 dark:border-slate-800">
              {isAuthenticated ? (
                <Link
                  href="/dashboard"
                  className="w-full py-3 text-center font-bold text-white bg-indigo-600 rounded-xl shadow-md"
                >
                  Личный кабинет
                </Link>
              ) : (
                <>
                  <button
                    onClick={() => {
                      setIsMobileMenuOpen(false);
                      openAuth('login');
                    }}
                    className="w-full py-3 text-center font-medium text-slate-700 dark:text-slate-200 bg-slate-100 dark:bg-slate-800 rounded-xl"
                  >
                    Войти
                  </button>
                  <button
                    onClick={() => {
                      setIsMobileMenuOpen(false);
                      openAuth('register');
                    }}
                    className="w-full py-3 text-center font-semibold text-white bg-indigo-600 rounded-xl shadow-md"
                  >
                    Регистрация
                  </button>
                </>
              )}
            </div>
          </div>
        )}
      </header>

      {/* Hero Section */}
      <section
        id="hero"
        className="relative pt-12 pb-20 lg:pt-20 lg:pb-28 overflow-hidden bg-gradient-to-b from-indigo-50/60 via-slate-50 to-slate-50 dark:from-slate-900 dark:via-slate-950 dark:to-slate-950 border-b border-slate-200/60 dark:border-slate-800/80"
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid lg:grid-cols-12 gap-12 items-center">
            {/* Left Col */}
            <div className="lg:col-span-7 space-y-6 text-center lg:text-left">
              <div className="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-indigo-50 dark:bg-indigo-950/80 border border-indigo-200/80 dark:border-indigo-800 text-indigo-700 dark:text-indigo-300 text-xs sm:text-sm font-semibold tracking-wide">
                <span className="flex h-2 w-2 rounded-full bg-emerald-500 animate-pulse"></span>
                Умная система онлайн-тестирования и курсов
              </div>

              <h1 className="text-4xl sm:text-5xl lg:text-6xl font-extrabold text-slate-900 dark:text-white tracking-tight leading-[1.15]">
                Платформа для{' '}
                <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 via-indigo-500 to-cyan-500 dark:from-indigo-400 dark:via-indigo-300 dark:to-cyan-400">
                  создания курсов
                </span>{' '}
                и автоматической проверки знаний
              </h1>

              <p className="text-lg sm:text-xl text-slate-600 dark:text-slate-300 max-w-2xl mx-auto lg:mx-0 leading-relaxed font-normal">
                Создавайте интерактивные обучающие материалы, проводите комплексные тестирования со срезом знаний и моментальной автопроверкой результатов.
              </p>

              {/* Action Buttons */}
              <div className="pt-2 flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-4">
                <button
                  onClick={() => openAuth('register')}
                  className="w-full sm:w-auto px-8 py-4 bg-indigo-600 hover:bg-indigo-700 active:scale-95 text-white font-bold rounded-2xl shadow-xl shadow-indigo-600/30 transition-all flex items-center justify-center gap-3 group cursor-pointer"
                >
                  <span>Начать бесплатно</span>
                  <ArrowRight className="group-hover:translate-x-1 transition-transform" size={18} />
                </button>
                <a
                  href="#how-it-works"
                  className="w-full sm:w-auto px-7 py-4 bg-white dark:bg-slate-900 hover:bg-slate-100 dark:hover:bg-slate-800 border border-slate-300/80 dark:border-slate-700 text-slate-700 dark:text-slate-200 font-semibold rounded-2xl shadow-xs transition-all text-center"
                >
                  Как работает платформа?
                </a>
              </div>

              {/* Trust Metrics */}
              <div className="pt-6 grid grid-cols-3 gap-4 border-t border-slate-200/80 dark:border-slate-800 max-w-md mx-auto lg:mx-0">
                <div>
                  <p className="text-2xl lg:text-3xl font-black text-slate-900 dark:text-white">45,000+</p>
                  <p className="text-xs text-slate-500 dark:text-slate-400 font-medium">Пройденных тестов</p>
                </div>
                <div>
                  <p className="text-2xl lg:text-3xl font-black text-indigo-600 dark:text-indigo-400">600+</p>
                  <p className="text-xs text-slate-500 dark:text-slate-400 font-medium">Курсов и программ</p>
                </div>
                <div>
                  <p className="text-2xl lg:text-3xl font-black text-emerald-600 dark:text-emerald-400">100%</p>
                  <p className="text-xs text-slate-500 dark:text-slate-400 font-medium">Автоматический срез</p>
                </div>
              </div>
            </div>

            {/* Right Col: Interactive Mockup */}
            <div className="lg:col-span-5 relative">
              <div className="relative mx-auto max-w-md lg:max-w-none">
                <div className="bg-gradient-to-br from-indigo-500 to-cyan-500 rounded-3xl p-1 shadow-2xl shadow-indigo-500/20">
                  <div className="bg-slate-900 rounded-[22px] overflow-hidden p-6 text-white space-y-5">
                    <div className="flex items-center justify-between border-b border-slate-800 pb-4">
                      <div className="flex items-center gap-2">
                        <span className="w-2.5 h-2.5 rounded-full bg-emerald-400"></span>
                        <span className="text-xs text-slate-400 font-medium">Тестирование: Модуль 3</span>
                      </div>
                      <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-indigo-500/20 text-indigo-300 text-xs font-semibold">
                        <Clock size={14} />
                        <span>14:20 мин</span>
                      </span>
                    </div>

                    <div className="space-y-3">
                      <div className="flex justify-between items-center text-xs text-slate-400">
                        <span>Вопрос 4 из 10</span>
                        <span className="text-cyan-400 font-semibold">Сложность: Средняя</span>
                      </div>
                      <h4 className="text-sm font-semibold text-slate-100">
                        Какой метод валидации наиболее эффективен для проверки остаточных знаний?
                      </h4>
                    </div>

                    {/* Interactive Choices */}
                    <div className="space-y-2.5">
                      <button
                        onClick={() => setHeroSelectedOption(0)}
                        className={`w-full text-left p-3 rounded-xl border text-xs transition-all flex items-center justify-between cursor-pointer ${
                          heroSelectedOption === 0
                            ? 'bg-rose-950/60 border-rose-500/60 text-rose-200'
                            : 'bg-slate-800/80 hover:bg-slate-700/80 border-slate-700 text-slate-300'
                        }`}
                      >
                        <span>A) Однократное устное собеседование</span>
                        <span
                          className={`w-4 h-4 rounded-full flex items-center justify-center text-[10px] ${
                            heroSelectedOption === 0
                              ? 'bg-rose-500 text-white'
                              : 'border border-slate-600'
                          }`}
                        >
                          {heroSelectedOption === 0 ? '✕' : ''}
                        </span>
                      </button>

                      <button
                        onClick={() => setHeroSelectedOption(1)}
                        className={`w-full text-left p-3 rounded-xl border text-xs transition-all flex items-center justify-between cursor-pointer ${
                          heroSelectedOption === 1 || heroSelectedOption === null
                            ? 'bg-indigo-950/60 border-indigo-500/50 text-indigo-200'
                            : 'bg-slate-800/80 border-slate-700 text-slate-300'
                        }`}
                      >
                        <span>B) Интервальное онлайн-тестирование со срезом</span>
                        <span className="w-4 h-4 rounded-full bg-emerald-500 text-white flex items-center justify-center text-[10px]">
                          ✓
                        </span>
                      </button>

                      <button
                        onClick={() => setHeroSelectedOption(2)}
                        className={`w-full text-left p-3 rounded-xl border text-xs transition-all flex items-center justify-between cursor-pointer ${
                          heroSelectedOption === 2
                            ? 'bg-rose-950/60 border-rose-500/60 text-rose-200'
                            : 'bg-slate-800/80 hover:bg-slate-700/80 border-slate-700 text-slate-300'
                        }`}
                      >
                        <span>C) Пассивное чтение конспектов</span>
                        <span
                          className={`w-4 h-4 rounded-full flex items-center justify-center text-[10px] ${
                            heroSelectedOption === 2
                              ? 'bg-rose-500 text-white'
                              : 'border border-slate-600'
                          }`}
                        >
                          {heroSelectedOption === 2 ? '✕' : ''}
                        </span>
                      </button>
                    </div>

                    <div className="space-y-2 pt-2">
                      <div className="flex justify-between text-xs font-medium text-slate-400">
                        <span>Прогресс прохождения</span>
                        <span className="text-indigo-400 font-bold">85% баллов</span>
                      </div>
                      <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                        <div className="bg-gradient-to-r from-indigo-500 to-emerald-400 h-full w-[85%] rounded-full"></div>
                      </div>
                    </div>

                    <button
                      onClick={() => openAuth('register')}
                      className="w-full py-3 bg-gradient-to-r from-indigo-500 to-cyan-500 text-white rounded-xl font-bold text-sm hover:opacity-95 transition-opacity cursor-pointer"
                    >
                      Создать свой первый тест
                    </button>
                  </div>
                </div>

                {/* Floating Badge */}
                <div className="absolute -bottom-5 -left-5 bg-white dark:bg-slate-800 p-4 rounded-2xl shadow-xl border border-slate-100 dark:border-slate-700 flex items-center gap-3">
                  <div className="p-2.5 bg-emerald-100 dark:bg-emerald-950 text-emerald-600 rounded-xl">
                    <CheckCircle2 size={24} />
                  </div>
                  <div>
                    <p className="text-xs text-slate-500 dark:text-slate-400 font-medium">Мгновенная проверка</p>
                    <p className="text-sm font-bold text-slate-900 dark:text-white">0.2 сек на расчет баллов</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Section: About Us */}
      <section id="about" className="py-20 bg-white dark:bg-slate-900 border-b border-slate-200/80 dark:border-slate-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="max-w-3xl mx-auto text-center space-y-4">
            <span className="text-xs font-bold text-indigo-600 dark:text-indigo-400 tracking-wider uppercase bg-indigo-50 dark:bg-indigo-950/80 px-3.5 py-1.5 rounded-full">
              О платформе ED
            </span>
            <h2 className="text-3xl sm:text-4xl font-extrabold text-slate-900 dark:text-white">
              Единое пространство для структурированных курсов и тестов
            </h2>
            <p className="text-slate-600 dark:text-slate-300 text-base sm:text-lg">
              ED помогает авторам быстро собирать обучающие программы и тесты любой сложности, а ученикам — последовательно осваивать теорию и проверять знания на практике.
            </p>
          </div>

          <div className="mt-14 grid md:grid-cols-3 gap-8">
            <div className="p-8 rounded-3xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200/70 dark:border-slate-700 hover:shadow-lg transition-all space-y-3">
              <div className="w-12 h-12 rounded-2xl bg-indigo-100 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 flex items-center justify-center font-bold text-xl mb-4">
                01
              </div>
              <h3 className="text-xl font-bold text-slate-900 dark:text-white">Конструктор курсов и уроков</h3>
              <p className="text-slate-600 dark:text-slate-300 text-sm leading-relaxed">
                Структурируйте материал по главам и модулям, прикрепляйте практические файлы, ссылки и инструкции для пошагового прохождения.
              </p>
            </div>

            <div className="p-8 rounded-3xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200/70 dark:border-slate-700 hover:shadow-lg transition-all space-y-3">
              <div className="w-12 h-12 rounded-2xl bg-cyan-100 dark:bg-cyan-950 text-cyan-600 dark:text-cyan-400 flex items-center justify-center font-bold text-xl mb-4">
                02
              </div>
              <h3 className="text-xl font-bold text-slate-900 dark:text-white">Генератор тестов и срезов</h3>
              <p className="text-slate-600 dark:text-slate-300 text-sm leading-relaxed">
                Тесты с одиночным и множественным выбором, вводом ответов, рандомизацией вопросов, таймерами и защитой от списывания.
              </p>
            </div>

            <div className="p-8 rounded-3xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200/70 dark:border-slate-700 hover:shadow-lg transition-all space-y-3">
              <div className="w-12 h-12 rounded-2xl bg-emerald-100 dark:bg-emerald-950 text-emerald-600 dark:text-emerald-400 flex items-center justify-center font-bold text-xl mb-4">
                03
              </div>
              <h3 className="text-xl font-bold text-slate-900 dark:text-white">Автопроверка и статистика</h3>
              <p className="text-slate-600 dark:text-slate-300 text-sm leading-relaxed">
                Мгновенный расчет баллов, подробный отчет об ошибках для студента и глубокая аналитика успеваемости для преподавателя.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Section: Advantages */}
      <section id="advantages" className="py-20 bg-slate-50 dark:bg-slate-950 border-b border-slate-200/80 dark:border-slate-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-3xl mx-auto space-y-4">
            <span className="text-xs font-bold text-indigo-600 dark:text-indigo-400 tracking-wider uppercase bg-indigo-100/70 dark:bg-indigo-950/80 px-3.5 py-1.5 rounded-full">
              Наши преимущества
            </span>
            <h2 className="text-3xl sm:text-4xl font-extrabold text-slate-900 dark:text-white">
              Почему для курсов и тестов выбирают ED
            </h2>
            <p className="text-slate-600 dark:text-slate-300 text-base sm:text-lg">
              Инструменты, которые экономят часы преподавателей на проверку заданий и делают учебу понятной для студентов.
            </p>
          </div>

          <div className="mt-14 grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {/* Card 1 */}
            <div className="bg-white dark:bg-slate-900 p-7 rounded-3xl border border-slate-200 dark:border-slate-800 shadow-xs hover:shadow-xl hover:-translate-y-1 transition-all group">
              <div className="w-14 h-14 rounded-2xl bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 flex items-center justify-center mb-6 group-hover:bg-indigo-600 group-hover:text-white transition-colors">
                <Zap size={28} />
              </div>
              <h3 className="text-lg font-bold text-slate-900 dark:text-white mb-2">Мгновенная автопроверка тестов</h3>
              <p className="text-slate-600 dark:text-slate-400 text-sm leading-relaxed">
                Никакого ручного пересчета баллов: студент видит свой результат, правильные пояснения и рекомендации сразу после завершения теста.
              </p>
            </div>

            {/* Card 2 */}
            <div className="bg-white dark:bg-slate-900 p-7 rounded-3xl border border-slate-200 dark:border-slate-800 shadow-xs hover:shadow-xl hover:-translate-y-1 transition-all group">
              <div className="w-14 h-14 rounded-2xl bg-cyan-50 dark:bg-cyan-950 text-cyan-600 dark:text-cyan-400 flex items-center justify-center mb-6 group-hover:bg-cyan-600 group-hover:text-white transition-colors">
                <Layers size={28} />
              </div>
              <h3 className="text-lg font-bold text-slate-900 dark:text-white mb-2">Банк вопросов и вариантов</h3>
              <p className="text-slate-600 dark:text-slate-400 text-sm leading-relaxed">
                Создавайте базы из сотен вопросов с автоперемешиванием — каждый ученик получает уникальный вариант теста.
              </p>
            </div>

            {/* Card 3 */}
            <div className="bg-white dark:bg-slate-900 p-7 rounded-3xl border border-slate-200 dark:border-slate-800 shadow-xs hover:shadow-xl hover:-translate-y-1 transition-all group">
              <div className="w-14 h-14 rounded-2xl bg-emerald-50 dark:bg-emerald-950 text-emerald-600 dark:text-emerald-400 flex items-center justify-center mb-6 group-hover:bg-emerald-600 group-hover:text-white transition-colors">
                <BarChart3 size={28} />
              </div>
              <h3 className="text-lg font-bold text-slate-900 dark:text-white mb-2">Глубокая аналитика срезов</h3>
              <p className="text-slate-600 dark:text-slate-400 text-sm leading-relaxed">
                Детальные графики по каждому вопросу: вы сразу видите, какие темы вызывают наибольшие затруднения у группы.
              </p>
            </div>

            {/* Card 4 */}
            <div className="bg-white dark:bg-slate-900 p-7 rounded-3xl border border-slate-200 dark:border-slate-800 shadow-xs hover:shadow-xl hover:-translate-y-1 transition-all group">
              <div className="w-14 h-14 rounded-2xl bg-amber-50 dark:bg-amber-950 text-amber-600 dark:text-amber-400 flex items-center justify-center mb-6 group-hover:bg-amber-600 group-hover:text-white transition-colors">
                <ShieldCheck size={28} />
              </div>
              <h3 className="text-lg font-bold text-slate-900 dark:text-white mb-2">Контроль честности прохождения</h3>
              <p className="text-slate-600 dark:text-slate-400 text-sm leading-relaxed">
                Ограничение по времени на тест или отдельный вопрос, блокировка переключения вкладок и фиксация попыток.
              </p>
            </div>

            {/* Card 5 */}
            <div className="bg-white dark:bg-slate-900 p-7 rounded-3xl border border-slate-200 dark:border-slate-800 shadow-xs hover:shadow-xl hover:-translate-y-1 transition-all group">
              <div className="w-14 h-14 rounded-2xl bg-purple-50 dark:bg-purple-950 text-purple-600 dark:text-purple-400 flex items-center justify-center mb-6 group-hover:bg-purple-600 group-hover:text-white transition-colors">
                <Award size={28} />
              </div>
              <h3 className="text-lg font-bold text-slate-900 dark:text-white mb-2">Автогенерация сертификатов</h3>
              <p className="text-slate-600 dark:text-slate-400 text-sm leading-relaxed">
                При успешной сдаче итогового тестирования система автоматически формирует персонализированный именной сертификат.
              </p>
            </div>

            {/* Card 6 */}
            <div className="bg-white dark:bg-slate-900 p-7 rounded-3xl border border-slate-200 dark:border-slate-800 shadow-xs hover:shadow-xl hover:-translate-y-1 transition-all group">
              <div className="w-14 h-14 rounded-2xl bg-rose-50 dark:bg-rose-950 text-rose-600 dark:text-rose-400 flex items-center justify-center mb-6 group-hover:bg-rose-600 group-hover:text-white transition-colors">
                <FileSpreadsheet size={28} />
              </div>
              <h3 className="text-lg font-bold text-slate-900 dark:text-white mb-2">Быстрый запуск без кода</h3>
              <p className="text-slate-600 dark:text-slate-400 text-sm leading-relaxed">
                Импорт вопросов, интуитивный визуальный редактор блоков и публикация курса для группы за пару минут.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Section: How It Works */}
      <section id="how-it-works" className="py-20 bg-white dark:bg-slate-900 border-b border-slate-200/80 dark:border-slate-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-3xl mx-auto space-y-4">
            <span className="text-xs font-bold text-indigo-600 dark:text-indigo-400 tracking-wider uppercase bg-indigo-50 dark:bg-indigo-950/80 px-3.5 py-1.5 rounded-full">
              Простой процесс
            </span>
            <h2 className="text-3xl sm:text-4xl font-extrabold text-slate-900 dark:text-white">
              Как начать пользоваться платформой ED
            </h2>
            <p className="text-slate-600 dark:text-slate-300 text-base sm:text-lg">
              Прозрачный процесс для эффективного обучения и точной оценки знаний.
            </p>
          </div>

          {/* Role Selector Tabs */}
          <div className="mt-10 flex justify-center">
            <div className="bg-slate-100 dark:bg-slate-800 p-1.5 rounded-2xl inline-flex border border-slate-200 dark:border-slate-700">
              <button
                onClick={() => setActiveRoleTab('student')}
                className={`px-6 py-2.5 rounded-xl font-semibold text-sm transition-all cursor-pointer ${
                  activeRoleTab === 'student'
                    ? 'bg-white dark:bg-indigo-600 text-indigo-600 dark:text-white shadow-xs'
                    : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'
                }`}
              >
                Для студентов и учеников
              </button>
              <button
                onClick={() => setActiveRoleTab('teacher')}
                className={`px-6 py-2.5 rounded-xl font-semibold text-sm transition-all cursor-pointer ${
                  activeRoleTab === 'teacher'
                    ? 'bg-white dark:bg-cyan-600 text-cyan-600 dark:text-white shadow-xs'
                    : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'
                }`}
              >
                Для авторов и преподавателей
              </button>
            </div>
          </div>

          {/* Student Steps */}
          {activeRoleTab === 'student' && (
            <div className="mt-12 grid md:grid-cols-4 gap-6">
              <div className="relative bg-slate-50 dark:bg-slate-800/70 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-700 space-y-2">
                <span className="text-4xl font-black text-indigo-300 dark:text-indigo-900">01</span>
                <h4 className="text-lg font-bold text-slate-900 dark:text-white">Регистрация профиля</h4>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Создайте аккаунт, чтобы сохранять прогресс и историю всех тестирований.
                </p>
              </div>
              <div className="relative bg-slate-50 dark:bg-slate-800/70 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-700 space-y-2">
                <span className="text-4xl font-black text-indigo-300 dark:text-indigo-900">02</span>
                <h4 className="text-lg font-bold text-slate-900 dark:text-white">Изучение теории</h4>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Проходите интерактивные модули курса в удобном для вас темпе.
                </p>
              </div>
              <div className="relative bg-slate-50 dark:bg-slate-800/70 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-700 space-y-2">
                <span className="text-4xl font-black text-indigo-300 dark:text-indigo-900">03</span>
                <h4 className="text-lg font-bold text-slate-900 dark:text-white">Сдача тестов и срезов</h4>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Отвечайте на вопросы, проверяйте знания и моментально получайте анализ ошибок.
                </p>
              </div>
              <div className="relative bg-slate-50 dark:bg-slate-800/70 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-700 space-y-2">
                <span className="text-4xl font-black text-indigo-300 dark:text-indigo-900">04</span>
                <h4 className="text-lg font-bold text-slate-900 dark:text-white">Итоги и сертификат</h4>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Получите цифровой сертификат с баллом и рекомендациями по развитию.
                </p>
              </div>
            </div>
          )}

          {/* Teacher Steps */}
          {activeRoleTab === 'teacher' && (
            <div className="mt-12 grid md:grid-cols-4 gap-6">
              <div className="relative bg-slate-50 dark:bg-slate-800/70 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-700 space-y-2">
                <span className="text-4xl font-black text-cyan-300 dark:text-cyan-900">01</span>
                <h4 className="text-lg font-bold text-slate-900 dark:text-white">Создание курса</h4>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Оформите разделы, уроки и соберите учебный контент в визуальном редакторе Puck.
                </p>
              </div>
              <div className="relative bg-slate-50 dark:bg-slate-800/70 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-700 space-y-2">
                <span className="text-4xl font-black text-cyan-300 dark:text-cyan-900">02</span>
                <h4 className="text-lg font-bold text-slate-900 dark:text-white">Настройка тестирования</h4>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Добавьте вопросы, шкалу оценивания, таймер и условия прохождения.
                </p>
              </div>
              <div className="relative bg-slate-50 dark:bg-slate-800/70 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-700 space-y-2">
                <span className="text-4xl font-black text-cyan-300 dark:text-cyan-900">03</span>
                <h4 className="text-lg font-bold text-slate-900 dark:text-white">Назначение группам</h4>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Отправьте ссылку ученикам или добавьте их в группу курса в пару кликов.
                </p>
              </div>
              <div className="relative bg-slate-50 dark:bg-slate-800/70 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-700 space-y-2">
                <span className="text-4xl font-black text-cyan-300 dark:text-cyan-900">04</span>
                <h4 className="text-lg font-bold text-slate-900 dark:text-white">Сводная ведомость</h4>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Просматривайте журнал с оценками, прогрессом и проверяйте домашние задания.
                </p>
              </div>
            </div>
          )}
        </div>
      </section>

      {/* Section: Popular Courses */}
      <section id="courses" className="py-20 bg-slate-50 dark:bg-slate-950 border-b border-slate-200/80 dark:border-slate-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col md:flex-row md:items-end justify-between mb-12 gap-4">
            <div>
              <span className="text-xs font-bold text-indigo-600 dark:text-indigo-400 tracking-wider uppercase bg-indigo-100/70 dark:bg-indigo-950 px-3.5 py-1.5 rounded-full">
                Каталог программ
              </span>
              <h2 className="text-3xl sm:text-4xl font-extrabold text-slate-900 dark:text-white mt-3">
                Популярные курсы и тесты
              </h2>
            </div>
            <Link
              href="/courses"
              className="text-indigo-600 dark:text-indigo-400 font-bold text-sm hover:underline flex items-center gap-1"
            >
              Смотреть каталог программ &rarr;
            </Link>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            {/* Course 1 */}
            <div className="bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 overflow-hidden shadow-xs hover:shadow-xl transition-all flex flex-col">
              <div className="h-48 bg-gradient-to-tr from-indigo-900 to-indigo-600 p-6 flex flex-col justify-between text-white relative">
                <div className="flex justify-between items-center">
                  <span className="px-3 py-1 bg-white/20 backdrop-blur-md rounded-full text-xs font-semibold">
                    IT & Разработка
                  </span>
                  <span className="text-xs bg-indigo-800/80 px-2.5 py-1 rounded-full">12 тестов</span>
                </div>
                <div>
                  <span className="text-xs text-indigo-200">18 учебных модулей</span>
                  <h3 className="text-xl font-bold">Frontend и JavaScript с автотестами</h3>
                </div>
              </div>
              <div className="p-6 flex-1 flex flex-col justify-between space-y-4">
                <p className="text-sm text-slate-600 dark:text-slate-300">
                  Пошаговый курс с интерактивными проверками кода, промежуточными квизами и итоговым квалификационным тестом.
                </p>
                <div className="pt-4 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between">
                  <div>
                    <span className="text-xs text-slate-400">Срезы знаний</span>
                    <p className="text-sm font-bold text-slate-900 dark:text-white">Сертификат по итогам</p>
                  </div>
                  <Link
                    href="/courses"
                    className="px-4 py-2 bg-indigo-50 dark:bg-indigo-950 hover:bg-indigo-600 hover:text-white text-indigo-600 dark:text-indigo-400 font-semibold rounded-xl text-sm transition-colors"
                  >
                    Пройти тест
                  </Link>
                </div>
              </div>
            </div>

            {/* Course 2 */}
            <div className="bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 overflow-hidden shadow-xs hover:shadow-xl transition-all flex flex-col">
              <div className="h-48 bg-gradient-to-tr from-cyan-900 to-cyan-600 p-6 flex flex-col justify-between text-white relative">
                <div className="flex justify-between items-center">
                  <span className="px-3 py-1 bg-white/20 backdrop-blur-md rounded-full text-xs font-semibold">
                    Маркетинг & Анализ
                  </span>
                  <span className="text-xs bg-cyan-800/80 px-2.5 py-1 rounded-full">8 тестов</span>
                </div>
                <div>
                  <span className="text-xs text-cyan-200">10 учебных модулей</span>
                  <h3 className="text-xl font-bold">Digital-маркетинг: Аттестационный курс</h3>
                </div>
              </div>
              <div className="p-6 flex-1 flex flex-col justify-between space-y-4">
                <p className="text-sm text-slate-600 dark:text-slate-300">
                  Практические тесты по аналитике трафика, таргетированной рекламе и расчету unit-экономики.
                </p>
                <div className="pt-4 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between">
                  <div>
                    <span className="text-xs text-slate-400">Срезы знаний</span>
                    <p className="text-sm font-bold text-slate-900 dark:text-white">Балльный рейтинг</p>
                  </div>
                  <Link
                    href="/courses"
                    className="px-4 py-2 bg-cyan-50 dark:bg-cyan-950 hover:bg-cyan-600 hover:text-white text-cyan-600 dark:text-cyan-400 font-semibold rounded-xl text-sm transition-colors"
                  >
                    Пройти тест
                  </Link>
                </div>
              </div>
            </div>

            {/* Course 3 */}
            <div className="bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 overflow-hidden shadow-xs hover:shadow-xl transition-all flex flex-col">
              <div className="h-48 bg-gradient-to-tr from-emerald-900 to-emerald-600 p-6 flex flex-col justify-between text-white relative">
                <div className="flex justify-between items-center">
                  <span className="px-3 py-1 bg-white/20 backdrop-blur-md rounded-full text-xs font-semibold">
                    Data & Базы данных
                  </span>
                  <span className="text-xs bg-emerald-800/80 px-2.5 py-1 rounded-full">15 тестов</span>
                </div>
                <div>
                  <span className="text-xs text-emerald-200">14 учебных модулей</span>
                  <h3 className="text-xl font-bold">SQL и Анализ данных: Практикум</h3>
                </div>
              </div>
              <div className="p-6 flex-1 flex flex-col justify-between space-y-4">
                <p className="text-sm text-slate-600 dark:text-slate-300">
                  Банк из 200+ тестовых вопросов по SQL-запросам, оптимизации и аналитике с мгновенной проверкой.
                </p>
                <div className="pt-4 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between">
                  <div>
                    <span className="text-xs text-slate-400">Срезы знаний</span>
                    <p className="text-sm font-bold text-slate-900 dark:text-white">Автопроверка</p>
                  </div>
                  <Link
                    href="/courses"
                    className="px-4 py-2 bg-emerald-50 dark:bg-emerald-950 hover:bg-emerald-600 hover:text-white text-emerald-600 dark:text-emerald-400 font-semibold rounded-xl text-sm transition-colors"
                  >
                    Пройти тест
                  </Link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-16 bg-gradient-to-tr from-slate-900 via-indigo-950 to-slate-900 text-white relative overflow-hidden">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
          <div className="max-w-4xl mx-auto bg-gradient-to-r from-indigo-600 to-cyan-600 rounded-3xl p-8 sm:p-12 shadow-2xl flex flex-col md:flex-row items-center justify-between gap-8">
            <div className="space-y-3 text-center md:text-left">
              <h3 className="text-2xl sm:text-3xl font-extrabold text-white">
                Готовы создать курс или пройти тестирование?
              </h3>
              <p className="text-indigo-100 text-sm sm:text-base max-w-md">
                Зарегистрируйтесь прямо сейчас, чтобы получить доступ к личному кабинету, каталогу тестов и конструктору курсов!
              </p>
            </div>
            <div className="w-full md:w-auto flex-shrink-0">
              <button
                onClick={() => openAuth('register')}
                className="w-full md:w-auto px-8 py-4 bg-white hover:bg-slate-100 text-indigo-900 font-black rounded-2xl shadow-lg active:scale-95 transition-all text-base whitespace-nowrap cursor-pointer"
              >
                Зарегистрироваться бесплатно
              </button>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer id="contacts" className="bg-slate-900 text-slate-300 pt-16 pb-12 border-t border-slate-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-10 pb-12 border-b border-slate-800">
            {/* Brand Info */}
            <div className="lg:col-span-2 space-y-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-2xl bg-indigo-500 flex items-center justify-center text-white font-extrabold text-lg">
                  ED
                </div>
                <span className="font-black text-xl text-white tracking-tight">ED.Learn</span>
              </div>
              <p className="text-sm text-slate-400 max-w-sm leading-relaxed">
                Цифровая платформа для онлайн-курсов, организации тестирований и срезов знаний.
              </p>
              <div className="flex gap-3 pt-2">
                <span className="w-9 h-9 rounded-xl bg-slate-800 hover:bg-indigo-600 flex items-center justify-center text-slate-300 hover:text-white transition-colors cursor-pointer font-bold text-xs">
                  TG
                </span>
                <span className="w-9 h-9 rounded-xl bg-slate-800 hover:bg-indigo-600 flex items-center justify-center text-slate-300 hover:text-white transition-colors cursor-pointer font-bold text-xs">
                  VK
                </span>
                <span className="w-9 h-9 rounded-xl bg-slate-800 hover:bg-indigo-600 flex items-center justify-center text-slate-300 hover:text-white transition-colors cursor-pointer font-bold text-xs">
                  YT
                </span>
              </div>
            </div>

            {/* Navigation Quick Links */}
            <div className="space-y-3">
              <h4 className="text-white font-bold text-sm tracking-wider uppercase">Навигация</h4>
              <ul className="space-y-2 text-sm text-slate-400">
                <li><a href="#about" className="hover:text-white transition-colors">О платформе</a></li>
                <li><a href="#advantages" className="hover:text-white transition-colors">Преимущества</a></li>
                <li><a href="#how-it-works" className="hover:text-white transition-colors">Как пользоваться</a></li>
                <li><a href="#courses" className="hover:text-white transition-colors">Курсы и тесты</a></li>
              </ul>
            </div>

            {/* For Users */}
            <div className="space-y-3">
              <h4 className="text-white font-bold text-sm tracking-wider uppercase">Личный кабинет</h4>
              <ul className="space-y-2 text-sm text-slate-400">
                <li><button onClick={() => openAuth('login')} className="hover:text-white transition-colors cursor-pointer">Вход в систему</button></li>
                <li><button onClick={() => openAuth('register')} className="hover:text-white transition-colors cursor-pointer">Регистрация студента</button></li>
                <li><button onClick={() => openAuth('register')} className="hover:text-white transition-colors cursor-pointer">Кабинет преподавателя</button></li>
                <li><Link href="/courses" className="hover:text-white transition-colors">Каталог курсов</Link></li>
              </ul>
            </div>

            {/* Contacts Details */}
            <div className="space-y-3">
              <h4 className="text-white font-bold text-sm tracking-wider uppercase">Контакты</h4>
              <ul className="space-y-2 text-sm text-slate-400">
                <li className="flex items-center gap-2">
                  <span className="text-indigo-400 font-semibold">Email:</span>
                  <span className="text-white">support@ed-platform.io</span>
                </li>
                <li className="flex items-center gap-2">
                  <span className="text-indigo-400 font-semibold">Телефон:</span>
                  <span className="text-white">+7 (800) 555-35-35</span>
                </li>
                <li className="flex items-center gap-2">
                  <span className="text-indigo-400 font-semibold">Адрес:</span>
                  <span>Москва, Пресненская наб., 12</span>
                </li>
                <li>
                  <span className="inline-block mt-2 px-2.5 py-1 bg-emerald-500/10 text-emerald-400 text-xs rounded-lg font-medium">
                    ● Поддержка 24/7 онлайн
                  </span>
                </li>
              </ul>
            </div>
          </div>

          {/* Bottom Copyright */}
          <div className="mt-8 flex flex-col sm:flex-row items-center justify-between text-xs text-slate-500 gap-4">
            <p>&copy; 2026 ED.Learn Platform. Все права защищены.</p>
            <div className="flex gap-6">
              <span className="hover:underline cursor-pointer">Политика конфиденциальности</span>
              <span className="hover:underline cursor-pointer">Пользовательское соглашение</span>
            </div>
          </div>
        </div>
      </footer>

      {/* Built-in Auth Modal */}
      {isAuthModalOpen && (
        <div className="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-white dark:bg-slate-900 w-full max-w-md rounded-3xl p-6 sm:p-8 shadow-2xl border border-slate-100 dark:border-slate-800 relative transition-all animate-in fade-in zoom-in-95">
            {/* Close Button */}
            <button
              onClick={() => setIsAuthModalOpen(false)}
              className="absolute top-5 right-5 text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 p-2 rounded-xl bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 transition-colors cursor-pointer"
            >
              <X size={18} />
            </button>

            {/* Modal Header */}
            <div className="flex items-center gap-3 border-b border-slate-100 dark:border-slate-800 pb-4 mb-6">
              <div className="w-10 h-10 rounded-2xl bg-indigo-600 text-white font-extrabold flex items-center justify-center text-lg shadow-md shadow-indigo-600/30">
                ED
              </div>
              <div>
                <h3 className="text-xl font-bold text-slate-900 dark:text-white">
                  {authMode === 'login' ? 'С возвращением!' : 'Создать профиль'}
                </h3>
                <p className="text-xs text-slate-500 dark:text-slate-400">
                  {authMode === 'login'
                    ? 'Введите свои данные для доступа к платформе'
                    : 'Присоединяйтесь к платформе ED за 30 секунд'}
                </p>
              </div>
            </div>

            {/* Mode Switcher */}
            <div className="flex bg-slate-100 dark:bg-slate-800 p-1 rounded-2xl mb-6">
              <button
                onClick={() => {
                  setAuthMode('login');
                  setAuthError('');
                }}
                className={`flex-1 py-2 text-center text-sm font-semibold rounded-xl transition-all cursor-pointer ${
                  authMode === 'login'
                    ? 'bg-white dark:bg-indigo-600 text-indigo-600 dark:text-white shadow-xs'
                    : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'
                }`}
              >
                Вход
              </button>
              <button
                onClick={() => {
                  setAuthMode('register');
                  setAuthError('');
                }}
                className={`flex-1 py-2 text-center text-sm font-semibold rounded-xl transition-all cursor-pointer ${
                  authMode === 'register'
                    ? 'bg-white dark:bg-indigo-600 text-indigo-600 dark:text-white shadow-xs'
                    : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'
                }`}
              >
                Регистрация
              </button>
            </div>

            {/* Form */}
            <form onSubmit={handleAuthSubmit} className="space-y-4">
              {authError && (
                <div className="p-3 rounded-xl bg-rose-50 dark:bg-rose-950/60 border border-rose-200 dark:border-rose-800 text-rose-600 dark:text-rose-400 text-xs font-semibold">
                  {authError}
                </div>
              )}

              {authNotice && (
                <div className="p-3 rounded-xl bg-emerald-50 dark:bg-emerald-950/60 border border-emerald-200 dark:border-emerald-800 text-emerald-600 dark:text-emerald-400 text-xs font-semibold">
                  {authNotice}
                </div>
              )}

              {authMode === 'register' && (
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700 dark:text-slate-300">
                    Имя пользователя (Никнейм)
                  </label>
                  <input
                    type="text"
                    required
                    value={authUsername}
                    onChange={(e) => setAuthUsername(e.target.value)}
                    placeholder="alex_developer"
                    className="w-full px-4 py-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 text-slate-900 dark:text-white transition-all"
                  />
                </div>
              )}

              <div className="space-y-1">
                <label className="text-xs font-semibold text-slate-700 dark:text-slate-300">
                  Электронная почта
                </label>
                <input
                  type="email"
                  required
                  value={authEmail}
                  onChange={(e) => setAuthEmail(e.target.value)}
                  placeholder="example@ed.ru"
                  className="w-full px-4 py-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 text-slate-900 dark:text-white transition-all"
                />
              </div>

              <div className="space-y-1">
                <div className="flex justify-between items-center">
                  <label className="text-xs font-semibold text-slate-700 dark:text-slate-300">
                    Пароль
                  </label>
                </div>
                <div className="relative">
                  <input
                    type={showPassword ? 'text' : 'password'}
                    required
                    value={authPassword}
                    onChange={(e) => setAuthPassword(e.target.value)}
                    placeholder="••••••••"
                    className="w-full px-4 py-3 pr-10 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 text-slate-900 dark:text-white transition-all"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 cursor-pointer"
                  >
                    {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </button>
                </div>
              </div>

              <button
                type="submit"
                disabled={authLoading}
                className="w-full py-3.5 bg-indigo-600 hover:bg-indigo-700 active:scale-95 text-white font-bold rounded-xl shadow-lg shadow-indigo-600/30 transition-all text-sm mt-2 flex items-center justify-center gap-2 cursor-pointer"
              >
                <span>
                  {authLoading
                    ? 'Подождите...'
                    : authMode === 'login'
                    ? 'Войти в личный кабинет'
                    : 'Зарегистрироваться бесплатно'}
                </span>
              </button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
