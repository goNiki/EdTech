import mammoth from 'mammoth';
import { marked } from 'marked';
import { api } from './api';

export interface ImportResult {
  html: string;
  wordCount: number;
  tablesCount: number;
  extractedTitle?: string;
  imagesCount: number;
}

/**
 * Очистка HTML от артефактов Microsoft Word (mso-*, пустые теги, служебные стили)
 */
export function cleanWordHtml(html: string): string {
  let cleaned = html;

  // Удаляем xml / o:p теги
  cleaned = cleaned.replace(/<o:p[\s\S]*?<\/o:p>/gi, '');
  cleaned = cleaned.replace(/<\/?\w+:[^>]*>/gi, '');

  // Удаляем стили mso-*
  cleaned = cleaned.replace(/style="[^"]*mso-[^"]*"/gi, '');
  cleaned = cleaned.replace(/style='[^']*mso-[^']*'/gi, '');

  // Стилизуем таблицы под современный Tailwind
  cleaned = cleaned.replace(
    /<table/gi,
    '<table class="w-full my-4 border-collapse border border-slate-200 dark:border-slate-700 text-xs text-left"'
  );
  cleaned = cleaned.replace(
    /<th/gi,
    '<th class="border border-slate-200 dark:border-slate-700 bg-slate-100 dark:bg-slate-800 p-2.5 font-bold text-slate-800 dark:text-slate-200"'
  );
  cleaned = cleaned.replace(
    /<td/gi,
    '<td class="border border-slate-200 dark:border-slate-700 p-2.5 text-slate-700 dark:text-slate-300"'
  );

  // Стилизуем цитаты и списки
  cleaned = cleaned.replace(
    /<blockquote/gi,
    '<blockquote class="border-l-4 border-indigo-500 pl-4 py-1 my-3 text-slate-600 dark:text-slate-300 italic bg-indigo-50/40 dark:bg-indigo-950/20 rounded-r-xl"'
  );

  return cleaned.trim();
}

/**
 * Преобразование base64 DataURL в Blob
 */
function dataURItoBlob(dataURI: string): Blob {
  const arr = dataURI.split(',');
  const mime = arr[0].match(/:(.*?);/)?.[1] || 'image/png';
  const bstr = atob(arr[1]);
  let n = bstr.length;
  const u8arr = new Uint8Array(n);
  while (n--) {
    u8arr[n] = bstr.charCodeAt(n);
  }
  return new Blob([u8arr], { type: mime });
}

/**
 * Загрузка пачки изображений на сервер
 */
export async function uploadBatchImages(images: { id: string; dataUrl: string }[]): Promise<Record<string, string>> {
  if (images.length === 0) return {};

  const urlMap: Record<string, string> = {};

  try {
    const formData = new FormData();
    images.forEach((img, idx) => {
      const blob = dataURItoBlob(img.dataUrl);
      const ext = blob.type.split('/')[1] || 'png';
      formData.append('files[]', blob, `imported_image_${idx + 1}.${ext}`);
    });
    formData.append('category', 'lesson_media');

    const res = await api.post('/upload/batch', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });

    const data = res.data?.data || res.data;
    if (data?.uploaded && Array.isArray(data.uploaded)) {
      data.uploaded.forEach((item: any, idx: number) => {
        if (images[idx] && item.file_url) {
          urlMap[images[idx].id] = item.file_url;
        }
      });
    }
  } catch (err) {
    console.warn('Batch image upload failed, falling back to individual or local URLs:', err);
    // Если бэкенд batch еще не задеплоен, сохраняем оптимизированный base64
    images.forEach((img) => {
      urlMap[img.id] = img.dataUrl;
    });
  }

  return urlMap;
}

/**
 * Основная функция конвертации файла (.docx или .md) в чистый семантический HTML
 */
export async function convertDocumentToHtml(
  file: File,
  onProgress?: (status: string) => void
): Promise<ImportResult> {
  const fileName = file.name.toLowerCase();
  let rawHtml = '';
  const extractedImages: { id: string; dataUrl: string }[] = [];

  if (fileName.endsWith('.docx')) {
    onProgress?.('Чтение файла Word (.docx)...');
    const arrayBuffer = await file.arrayBuffer();

    onProgress?.('Парсинг стилей, таблиц и структуры...');
    const result = await mammoth.convertToHtml(
      { arrayBuffer },
      {
        convertImage: mammoth.images.inline((element) => {
          return element.read('base64').then((imageBuffer: string) => {
            const dataUrl = `data:${element.contentType};base64,${imageBuffer}`;
            const id = `img_${Math.random().toString(36).substr(2, 9)}`;
            extractedImages.push({ id, dataUrl });
            return { src: `__PLACEHOLDER_${id}__` };
          });
        }),
      }
    );

    rawHtml = result.value;

    if (extractedImages.length > 0) {
      onProgress?.(`Загрузка иллюстраций (${extractedImages.length} шт.)...`);
      const urlMap = await uploadBatchImages(extractedImages);
      extractedImages.forEach((img) => {
        const targetUrl = urlMap[img.id] || img.dataUrl;
        rawHtml = rawHtml.replace(`__PLACEHOLDER_${img.id}__`, targetUrl);
      });
    }
  } else if (fileName.endsWith('.md') || fileName.endsWith('.markdown')) {
    onProgress?.('Парсинг разметки Markdown...');
    const text = await file.text();
    rawHtml = await marked.parse(text);
  } else {
    throw new Error('Неподдерживаемый формат файла. Поддерживаются только .docx и .md');
  }

  onProgress?.('Очистка разметки...');
  const cleanedHtml = cleanWordHtml(rawHtml);

  // Вычисляем статистику
  const plainText = cleanedHtml.replace(/<[^>]+>/g, ' ').trim();
  const wordCount = plainText.length > 0 ? plainText.split(/\s+/).length : 0;
  const tablesCount = (cleanedHtml.match(/<table/gi) || []).length;

  // Извлекаем возможный первый заголовок H1
  let extractedTitle: string | undefined;
  const h1Match = cleanedHtml.match(/<h1[^>]*>([\s\S]*?)<\/h1>/i);
  if (h1Match && h1Match[1]) {
    extractedTitle = h1Match[1].replace(/<[^>]+>/g, '').trim();
  } else {
    // имя файла без расширения
    extractedTitle = file.name.replace(/\.[^/.]+$/, '');
  }

  return {
    html: cleanedHtml,
    wordCount,
    tablesCount,
    extractedTitle,
    imagesCount: extractedImages.length,
  };
}
