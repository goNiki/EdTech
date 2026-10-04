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
 * Добавление встроенных стилей text-align для гарантированного сохранения выравнивания
 */
function addInlineAlignmentStyles(html: string): string {
  let res = html;

  const alignments: [string, string][] = [
    ['text-center', 'text-align: center;'],
    ['text-right', 'text-align: right;'],
    ['text-justify', 'text-align: justify;'],
    ['text-left', 'text-align: left;'],
  ];

  alignments.forEach(([cls, styleRule]) => {
    const tagRegex = new RegExp(`(<(p|h1|h2|h3|h4|h5|h6|div|td|th)[^>]*class="[^"]*\\b${cls}\\b[^"]*"[^>]*>)`, 'gi');
    res = res.replace(tagRegex, (fullTag) => {
      if (/style\s*=\s*"/i.test(fullTag)) {
        if (!/text-align\s*:/i.test(fullTag)) {
          return fullTag.replace(/style="([^"]*)"/i, `style="${styleRule} $1"`);
        }
        return fullTag;
      } else if (/style\s*=\s*'/i.test(fullTag)) {
        if (!/text-align\s*:/i.test(fullTag)) {
          return fullTag.replace(/style='([^']*)'/i, `style='${styleRule} $1'`);
        }
        return fullTag;
      } else {
        return fullTag.replace(/>$/, ` style="${styleRule}">`);
      }
    });
  });

  return res;
}

/**
 * Очистка HTML от артефактов Microsoft Word (mso-*, пустые теги, служебные стили)
 * с бережным сохранением выравнивания, шрифтов, цветов и таблиц.
 */
export function cleanWordHtml(html: string): string {
  let cleaned = html;

  // Удаляем xml / o:p теги
  cleaned = cleaned.replace(/<o:p[\s\S]*?<\/o:p>/gi, '');
  cleaned = cleaned.replace(/<\/?\w+:[^>]*>/gi, '');

  // Удаляем исключительно служебные декларации mso-*, сохраняя валидные CSS правила
  cleaned = cleaned.replace(/mso-[^;:"']+(;)?/gi, '');
  cleaned = cleaned.replace(/style="\s*"/gi, '');
  cleaned = cleaned.replace(/style='\s*'/gi, '');

  // Навешиваем inline-стили выравнивания на блоки с классами text-center/right/justify
  cleaned = addInlineAlignmentStyles(cleaned);

  // Стилизуем таблицы под современный Tailwind без принудительного перебивания выравнивания
  cleaned = cleaned.replace(
    /<table/gi,
    '<table class="w-full my-4 border-collapse border border-slate-300 dark:border-slate-700 text-xs shadow-2xs rounded-xl overflow-hidden"'
  );
  cleaned = cleaned.replace(
    /<th/gi,
    '<th class="border border-slate-300 dark:border-slate-700 bg-slate-100 dark:bg-slate-800 p-2.5 font-bold text-slate-800 dark:text-slate-200"'
  );
  cleaned = cleaned.replace(
    /<td/gi,
    '<td class="border border-slate-300 dark:border-slate-700 p-2.5 text-slate-700 dark:text-slate-300"'
  );

  // Стилизуем цитаты и списки
  cleaned = cleaned.replace(
    /<blockquote/gi,
    '<blockquote class="border-l-4 border-indigo-500 pl-4 py-1.5 my-3 text-slate-600 dark:text-slate-300 italic bg-indigo-50/40 dark:bg-indigo-950/20 rounded-r-xl"'
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

    const docxStyleMap = [
      'u => u',
      'strike => s',
      'p.Heading1-center => h1.text-center:fresh',
      'p.Heading1-right => h1.text-right:fresh',
      'p.Heading1-justify => h1.text-justify:fresh',
      'p.Heading2-center => h2.text-center:fresh',
      'p.Heading2-right => h2.text-right:fresh',
      'p.Heading2-justify => h2.text-justify:fresh',
      'p.Heading3-center => h3.text-center:fresh',
      'p.Heading3-right => h3.text-right:fresh',
      'p.Heading3-justify => h3.text-justify:fresh',
      'p.Title-center => h1.text-center:fresh',
      'p.Title-right => h1.text-right:fresh',
      'p.align-center => p.text-center:fresh',
      'p.align-right => p.text-right:fresh',
      'p.align-justify => p.text-justify:fresh',
      'p.align-both => p.text-justify:fresh',
      "p[style-name='Heading 1'] => h1:fresh",
      "p[style-name='Heading 2'] => h2:fresh",
      "p[style-name='Heading 3'] => h3:fresh",
      "p[style-name='Title'] => h1:fresh",
    ];

    const docxTransform = mammoth.transforms.paragraph((p: any) => {
      const align = p.alignment;
      if (!align) return p;
      let normAlign = String(align).toLowerCase();
      if (normAlign === 'both') normAlign = 'justify';
      if (normAlign === 'end') normAlign = 'right';
      if (normAlign !== 'center' && normAlign !== 'right' && normAlign !== 'justify') return p;
      const styleId = String(p.styleId || '');
      const styleName = String(p.styleName || '');
      if (/heading\s*1/i.test(styleId) || /heading\s*1/i.test(styleName)) return { ...p, styleId: `Heading1-${normAlign}` };
      if (/heading\s*2/i.test(styleId) || /heading\s*2/i.test(styleName)) return { ...p, styleId: `Heading2-${normAlign}` };
      if (/heading\s*3/i.test(styleId) || /heading\s*3/i.test(styleName)) return { ...p, styleId: `Heading3-${normAlign}` };
      if (/title/i.test(styleId) || /title/i.test(styleName)) return { ...p, styleId: `Title-${normAlign}` };
      return { ...p, styleId: `align-${normAlign}` };
    });

    const result = await mammoth.convertToHtml(
      { arrayBuffer },
      {
        styleMap: docxStyleMap,
        transformDocument: docxTransform,
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
