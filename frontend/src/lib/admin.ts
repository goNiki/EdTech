import { api } from './api';

export interface AdminUser {
  id: number;
  email: string;
  username: string;
  first_name?: string;
  last_name?: string;
  role: string;
  is_banned: boolean;
  created_at: string;
}

export interface AdminUsersResponse {
  users: AdminUser[];
  total: number;
  page: number;
  page_size: number;
}

const STORAGE_KEY = 'edtech_admin_users_store';

const DEFAULT_MOCK_USERS: AdminUser[] = [
  {
    id: 1,
    email: 'admin@edtech.io',
    username: 'superadmin',
    first_name: 'Главный',
    last_name: 'Администратор',
    role: 'admin',
    is_banned: false,
    created_at: '2026-08-01T10:00:00Z',
  },
  {
    id: 2,
    email: 'teacher@edtech.io',
    username: 'alex_instructor',
    first_name: 'Алексей',
    last_name: 'Петров',
    role: 'teacher',
    is_banned: false,
    created_at: '2026-08-15T14:30:00Z',
  },
  {
    id: 3,
    email: 'author@edtech.io',
    username: 'elena_course',
    first_name: 'Елена',
    last_name: 'Соколова',
    role: 'author',
    is_banned: false,
    created_at: '2026-09-01T09:15:00Z',
  },
  {
    id: 4,
    email: 'student1@edtech.io',
    username: 'ivan_dev',
    first_name: 'Иван',
    last_name: 'Иванов',
    role: 'student',
    is_banned: false,
    created_at: '2026-09-10T11:20:00Z',
  },
  {
    id: 5,
    email: 'student2@edtech.io',
    username: 'maria_code',
    first_name: 'Мария',
    last_name: 'Кузнецова',
    role: 'student',
    is_banned: false,
    created_at: '2026-09-18T16:45:00Z',
  },
  {
    id: 6,
    email: 'spammer@example.com',
    username: 'bad_user',
    first_name: 'Нежелательный',
    last_name: 'Пользователь',
    role: 'student',
    is_banned: true,
    created_at: '2026-09-22T08:00:00Z',
  },
];

function getStoredUsers(): AdminUser[] {
  if (typeof window === 'undefined') return DEFAULT_MOCK_USERS;
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(DEFAULT_MOCK_USERS));
    return DEFAULT_MOCK_USERS;
  }
  try {
    return JSON.parse(raw);
  } catch {
    return DEFAULT_MOCK_USERS;
  }
}

function saveStoredUsers(users: AdminUser[]): void {
  if (typeof window === 'undefined') return;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(users));
}

/**
 * Получение списка пользователей панели администратора
 */
export async function fetchAdminUsers(params?: {
  search?: string;
  role?: string;
  is_banned?: boolean;
  page?: number;
  page_size?: number;
}): Promise<AdminUsersResponse> {
  const query = new URLSearchParams();
  if (params?.search) query.append('search', params.search);
  if (params?.role) query.append('role', params.role);
  if (params?.is_banned !== undefined) query.append('is_banned', String(params.is_banned));
  if (params?.page) query.append('page', String(params.page));
  if (params?.page_size) query.append('page_size', String(params.page_size));

  try {
    const res = await api.get(`/admin/users?${query.toString()}`);
    const data = res.data?.data || res.data;
    if (data && Array.isArray(data.users)) {
      return {
        users: data.users,
        total: data.total ?? data.users.length,
        page: data.page ?? 1,
        page_size: data.page_size ?? 20,
      };
    }
  } catch {
    // API endpoint fallback
  }

  // Fallback поиск и фильтрация по локальным пользователям
  let filtered = getStoredUsers();

  if (params?.search) {
    const s = params.search.toLowerCase();
    filtered = filtered.filter(
      (u) =>
        u.email.toLowerCase().includes(s) ||
        u.username.toLowerCase().includes(s) ||
        (u.first_name && u.first_name.toLowerCase().includes(s)) ||
        (u.last_name && u.last_name.toLowerCase().includes(s))
    );
  }

  if (params?.role && params.role !== 'all') {
    filtered = filtered.filter((u) => u.role === params.role);
  }

  if (params?.is_banned !== undefined) {
    filtered = filtered.filter((u) => u.is_banned === params.is_banned);
  }

  const page = params?.page || 1;
  const pageSize = params?.page_size || 20;
  const start = (page - 1) * pageSize;
  const paginated = filtered.slice(start, start + pageSize);

  return {
    users: paginated,
    total: filtered.length,
    page,
    page_size: pageSize,
  };
}

/**
 * Изменение роли пользователя
 */
export async function updateUserRole(userId: number, newRole: string): Promise<void> {
  try {
    await api.patch(`/admin/users/${userId}/role`, { role: newRole });
  } catch {
    // API fallback
  }

  const users = getStoredUsers();
  const updated = users.map((u) => (u.id === userId ? { ...u, role: newRole } : u));
  saveStoredUsers(updated);
}

/**
 * Блокировка / разблокировка пользователя
 */
export async function updateUserBan(userId: number, isBanned: boolean): Promise<void> {
  try {
    await api.patch(`/admin/users/${userId}/ban`, { is_banned: isBanned });
  } catch {
    // API fallback
  }

  const users = getStoredUsers();
  const updated = users.map((u) => (u.id === userId ? { ...u, is_banned: isBanned } : u));
  saveStoredUsers(updated);
}
