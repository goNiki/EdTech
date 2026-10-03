import { api } from '@/lib/api';

export interface Category {
  id: number;
  name: string;
  slug: string;
  description?: string;
  icon_url?: string;
  parent_id?: number | null;
  courses_count?: number;
}

export const DEFAULT_CATEGORIES: Category[] = [
  {
    id: 1,
    name: 'Программирование',
    slug: 'programming',
    description: 'Курсы по разработке ПО и системному программированию',
  },
  {
    id: 2,
    name: 'Веб-разработка',
    slug: 'web-development',
    description: 'Frontend, Backend и Fullstack разработка',
  },
  {
    id: 3,
    name: 'Data Science',
    slug: 'data-science',
    description: 'Наука о данных, машинное обучение и ИИ',
  },
  {
    id: 4,
    name: 'Дизайн',
    slug: 'design',
    description: 'UI/UX, веб-дизайн и компьютерная графика',
  },
  {
    id: 5,
    name: 'Бизнес',
    slug: 'business',
    description: 'Бизнес, управление проектами и маркетинг',
  },
];

/**
 * Загружает список активных категорий с бэкенда (GET /api/v1/categories)
 * с безопасным fallback на эталонный справочник при временной недоступности бэкенда.
 */
export async function fetchCategories(): Promise<Category[]> {
  try {
    const res = await api.get('/categories');
    const data = res.data?.data || res.data;
    const list = data?.categories || (Array.isArray(data) ? data : null);

    if (Array.isArray(list) && list.length > 0) {
      return list.map((item: any) => ({
        id: Number(item.id || item.ID),
        name: String(item.name || item.Name || ''),
        slug: String(item.slug || item.Slug || ''),
        description: item.description || item.Description,
        icon_url: item.icon_url || item.IconURL,
        parent_id: item.parent_id ?? item.ParentID ?? null,
        courses_count: item.courses_count ?? item.CoursesCount,
      }));
    }
  } catch (err) {
    // В случае если эндпоинт категорий еще не поднят или недоступен,
    // используем эталонный справочник без прерывания пользовательского сценария
    console.warn('GET /categories fallback to DEFAULT_CATEGORIES:', err);
  }

  return DEFAULT_CATEGORIES;
}

/**
 * Получает наименование категории по её идентификатору
 */
export function getCategoryName(categoryId?: number | null, categories: Category[] = DEFAULT_CATEGORIES): string | null {
  if (!categoryId) return null;
  const found = categories.find((c) => c.id === categoryId);
  return found ? found.name : null;
}
