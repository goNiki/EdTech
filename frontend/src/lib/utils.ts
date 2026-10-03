import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

const RU_TO_EN_MAP: Record<string, string> = {
  а: 'a', б: 'b', в: 'v', г: 'g', д: 'd', е: 'e', ё: 'yo', ж: 'zh',
  з: 'z', и: 'i', й: 'y', к: 'k', л: 'l', м: 'm', н: 'n', о: 'o',
  п: 'p', р: 'r', с: 's', т: 't', у: 'u', ф: 'f', х: 'h', ц: 'ts',
  ч: 'ch', ш: 'sh', щ: 'sch', ъ: '', ы: 'y', ь: '', э: 'e', ю: 'yu', я: 'ya',
};

/**
 * Транслитерирует кириллические символы в латиницу по стандарту ГОСТ/ISO
 */
export function transliterate(text: string): string {
  if (!text) return '';
  return text
    .split('')
    .map((char) => {
      const lower = char.toLowerCase();
      if (Object.prototype.hasOwnProperty.call(RU_TO_EN_MAP, lower)) {
        const mapped = RU_TO_EN_MAP[lower];
        return char === char.toUpperCase() && mapped.length > 0
          ? mapped.charAt(0).toUpperCase() + mapped.slice(1)
          : mapped;
      }
      return char;
    })
    .join('');
}

/**
 * Генерирует SEO-friendly URL-слаг из произвольного текста (включая кириллицу)
 */
export function slugify(text: string): string {
  if (!text) return '';
  return transliterate(text)
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9 -]/g, '')
    .replace(/[\s_]+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-+|-+$/g, '');
}

/**
 * Проверяет соответствие слага безопасному формату URL [a-z0-9-]
 */
export function isValidSlug(slug: string): boolean {
  if (!slug) return false;
  return /^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(slug);
}
