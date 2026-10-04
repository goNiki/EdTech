import mammoth from 'mammoth';

export type QuizBlockType =
  | 'single'
  | 'multi'
  | 'match'
  | 'dropdown'
  | 'input'
  | 'sequence'
  | 'essay'
  | 'file';

export interface ParsedQuizOption {
  text: string;
  isCorrect: boolean;
  explain?: string;
}

export interface ParsedQuizPair {
  left: string;
  right: string;
}

export interface ParsedQuizSequenceItem {
  text: string;
}

export interface ParsedQuiz {
  id: string;
  rawIndex: number;
  question: string;
  type: QuizBlockType;
  points: number;
  options?: ParsedQuizOption[];
  pairs?: ParsedQuizPair[];
  templateText?: string;
  blanksCount?: number;
  items?: ParsedQuizSequenceItem[];
  rubric?: string;
  sampleAnswer?: string;
  fileTitle?: string;
  instructions?: string;
  allowedTypes?: string;
  maxSizeMB?: number;
  warnings: string[];
}

export const SAMPLE_QUIZ_TEMPLATES = [
  {
    name: '📋 Все типы (Мега-пример)',
    text: `[SINGLE]
Какая архитектура изолирует доменную логику от внешних фреймворков и БД?
A. Active Record
B. Clean Architecture *
C. Monolithic Core
D. Transaction Script
POINTS: 10

[MULTI]
Выберите протоколы прикладного уровня:
[x] HTTP/2
[x] gRPC
[ ] TCP
[x] WebSockets
POINTS: 10

[MATCH]
Сопоставьте шаблоны и их определения:
AST Tree :: Иерархическое синтаксическое дерево
Block Props :: Параметры компонентов и данных
Puck Canvas :: Интерактивная область визуального рендеринга
POINTS: 15

[DROPDOWN]
Заполните пропуски в предложении:
Спецификация для REST API называется {OpenAPI; Buf, WSDL}, а для RPC систем чаще всего используется {gRPC; GraphQL, SOAP}.
POINTS: 10

[INPUT]
Вставьте точные термины:
Команда для переключения веток в Git: git [checkout]. Для отката изменений рабочей директории используется git [reset].
POINTS: 10

[SEQUENCE]
Расставьте этапы обработки запроса в правильном порядке:
-> Прием HTTP запроса роутером
-> Валидация JWT токена в Middleware
-> Вызов бизнес-логики в сервисе
-> Запись в БД и возврат JSON ответа
POINTS: 15

[ESSAY]
Опишите преимущества применения Чистой Архитектуры в микросервисах:
РУБРИКА: Полнота ответа, понимание инверсии зависимостей, тестопригодность.
ОТВЕТ: Независимость от внешних библиотек, легкость написания модульных тестов и быстрая замена адаптеров БД.
POINTS: 25

[FILE]
Лабораторная работа: Реализация роутинга на Chi
ИНСТРУКЦИЯ: Прикрепите архив с исходным кодом (.zip) или PDF отчет с листингом обработчиков.
ФОРМАТЫ: .zip, .tar.gz, .pdf, .go
ЛИМИТ: 25
POINTS: 50`,
  },
  {
    name: '🔘 1. Один выбор (Single)',
    text: `[SINGLE]
1. Какой HTTP-статус означает успешное создание ресурса на сервере?
A. 200 OK
B. 201 Created *
C. 204 No Content
D. 301 Moved Permanently
POINTS: 5`,
  },
  {
    name: '☑️ 2. Множественный выбор (Multi)',
    text: `[MULTI]
1. Какие принципы входят в аббревиатуру SOLID?
[x] Single Responsibility
[x] Open-Closed Principle
[ ] Over-Engineering Pattern
[x] Dependency Inversion
POINTS: 10`,
  },
  {
    name: '🔄 3. Сопоставление пар (Matching)',
    text: `[MATCH]
Сопоставьте HTTP методы с их свойствами:
GET :: Безопасный и идемпотентный метод получения данных
POST :: Неидемпотентный метод создания ресурса
PUT :: Идемпотентный метод полной замены ресурса
DELETE :: Метод удаления ресурса
POINTS: 15`,
  },
  {
    name: '🔽 4. Выпадающие списки (Dropdown)',
    text: `[DROPDOWN]
Заполните термины в предложении:
В протоколе OAuth 2.0 за выдачу ключей отвечает {Authorization Server; Resource Owner, Client App}, а защищенные данные предоставляет {Resource Server; User Agent, Identity Proxy}.
POINTS: 10`,
  },
  {
    name: '✏️ 5. Пропуски в тексте (Input)',
    text: `[INPUT]
Заполните пропущенные команды CLI:
Для создания коммита в Git используется команда git [commit], а для отправки в удаленный репозиторий — git [push].
POINTS: 10`,
  },
  {
    name: '🔢 6. Последовательность (Sequence)',
    text: `[SEQUENCE]
Расставьте жизненный цикл запроса в веб-сервере:
-> Разрешение DNS имени в IP адрес
-> Установка TCP соединения и TLS рукопожатие
-> Передача HTTP заголовков и тела запроса
-> Обработка запроса контроллером бэкенда
-> Отправка HTTP ответа клиенту
POINTS: 15`,
  },
  {
    name: '📝 7. Развернутый ответ (Essay)',
    text: `[ESSAY]
Объясните концепцию горутин (Goroutines) и каналов в Go:
РУБРИКА: Понимание конкурентности vs параллелизма, CSP модель, управление утечками памяти.
ОТВЕТ: Горутины — легковесные потоки под управлением Go runtime. Каналы служат для безопасной передачи данных между ними по модели CSP.
POINTS: 20`,
  },
  {
    name: '📁 8. Загрузка файла (File Upload)',
    text: `[FILE]
Практическое задание: Реализация микросервиса авторизации
ИНСТРУКЦИЯ: Загрузите архив репозитория с исходным кодом, Dockerfile и файлом README с инструкцией по запуску.
ФОРМАТЫ: .zip, .tar.gz, .rar
ЛИМИТ: 50
POINTS: 30`,
  },
];

/**
 * Парсер сырого текста во все 8 типов интерактивных заданий
 */
export function parseQuizzesFromText(rawText: string): ParsedQuiz[] {
  if (!rawText.trim()) return [];

  const rawLines = rawText.split(/\r?\n/).map((l) => l.trim());
  const blocks: string[][] = [];
  let currentBlock: string[] = [];

  const isTypeTag = (line: string) =>
    /^\[(SINGLE|MULTI|MATCH|MATCHING|DROPDOWN|SELECT|INPUT|BLANK|SEQUENCE|ORDER|ESSAY|TEXT|OPEN|FILE|FILE_UPLOAD|UPLOAD)\]/i.test(
      line
    );

  const isDivider = (line: string) => /^={3,}|^-{3,}/.test(line);

  for (const line of rawLines) {
    if (isDivider(line)) {
      if (currentBlock.length > 0) {
        blocks.push(currentBlock);
        currentBlock = [];
      }
      continue;
    }

    if (!line) {
      if (currentBlock.length > 0) {
        blocks.push(currentBlock);
        currentBlock = [];
      }
      continue;
    }

    // Если начинается явный тег типа [TYPE]
    if (isTypeTag(line) && currentBlock.length > 0) {
      blocks.push(currentBlock);
      currentBlock = [];
    }

    // Если начинается нумерованный вопрос 1. или Вопрос 1: и в блоке уже есть данные
    const isNumberedStart = /^(?:#+\s*|\d+[\.\)]\s*|Вопрос\s*\d+[\.:]?\s*)/i.test(line);
    const hasBodyInCurrent = currentBlock.some(
      (l) =>
        isTypeTag(l) ||
        /^(?:[a-zA-Zа-яА-Я][\.\)]|\[[xX\s]\]|\*|\-|::|->|ANSWER|ОТВЕТ|РУБРИКА|ИНСТРУКЦИЯ)/i.test(l) ||
        /\{[^{}]+\}/.test(l)
    );

    if (isNumberedStart && currentBlock.length > 0 && hasBodyInCurrent) {
      blocks.push(currentBlock);
      currentBlock = [];
    }

    currentBlock.push(line);
  }

  if (currentBlock.length > 0) {
    blocks.push(currentBlock);
  }

  const parsedQuizzes: ParsedQuiz[] = [];

  blocks.forEach((block, idx) => {
    if (block.length === 0) return;

    const warnings: string[] = [];
    let detectedType: QuizBlockType | null = null;
    let explicitPoints: number | null = null;

    // 1. Поиск явного тега [TYPE]
    for (let i = 0; i < block.length; i++) {
      const line = block[i];
      const tagMatch = line.match(
        /^\[(SINGLE|MULTI|MATCH|MATCHING|DROPDOWN|SELECT|INPUT|BLANK|SEQUENCE|ORDER|ESSAY|TEXT|OPEN|FILE|FILE_UPLOAD|UPLOAD)\]/i
      );
      if (tagMatch) {
        const rawTag = tagMatch[1].toUpperCase();
        if (rawTag === 'SINGLE') detectedType = 'single';
        else if (rawTag === 'MULTI') detectedType = 'multi';
        else if (rawTag === 'MATCH' || rawTag === 'MATCHING') detectedType = 'match';
        else if (rawTag === 'DROPDOWN' || rawTag === 'SELECT') detectedType = 'dropdown';
        else if (rawTag === 'INPUT' || rawTag === 'BLANK') detectedType = 'input';
        else if (rawTag === 'SEQUENCE' || rawTag === 'ORDER') detectedType = 'sequence';
        else if (rawTag === 'ESSAY' || rawTag === 'TEXT' || rawTag === 'OPEN') detectedType = 'essay';
        else if (rawTag === 'FILE' || rawTag === 'FILE_UPLOAD' || rawTag === 'UPLOAD') detectedType = 'file';
        break;
      }
    }

    // 2. Поиск явных баллов POINTS: X
    for (const line of block) {
      const pMatch = line.match(/^(?:POINTS|БАЛЛЫ|БАЛЛ|MAX_POINTS)[\s:]+(\d+)/i);
      if (pMatch) {
        explicitPoints = parseInt(pMatch[1], 10);
      }
    }

    // Фильтруем служебные строки тегов и POINTS
    const cleanLines = block.filter((line) => {
      if (isTypeTag(line)) return false;
      if (/^(?:POINTS|БАЛЛЫ|БАЛЛ|MAX_POINTS)[\s:]+\d+/i.test(line)) return false;
      return true;
    });

    if (cleanLines.length === 0) return;

    // 3. Эвристический анализ, если тег не указан
    if (!detectedType) {
      if (cleanLines.some((l) => l.includes('::'))) {
        detectedType = 'match';
      } else if (cleanLines.some((l) => /\{[^{};]+;[^{}]+\}/.test(l) || /\{[^{}]+\}/.test(l))) {
        detectedType = 'dropdown';
      } else if (cleanLines.some((l) => /\[([^\]xX\s]+)\]/.test(l))) {
        detectedType = 'input';
      } else if (cleanLines.some((l) => /^->\s+|=>\s+/.test(l))) {
        detectedType = 'sequence';
      } else if (cleanLines.some((l) => /^(?:РУБРИКА|RUBRIC|КРИТЕРИИ):/i.test(l))) {
        detectedType = 'essay';
      } else if (cleanLines.some((l) => /^(?:ИНСТРУКЦИЯ|ФОРМАТЫ|ЛИМИТ):/i.test(l))) {
        detectedType = 'file';
      } else {
        // Проверяем наличие чекбоксов [x] или букв A. B. C.
        const checkboxCount = cleanLines.filter((l) => /^\[[xX]\]/i.test(l)).length;
        if (checkboxCount > 1) {
          detectedType = 'multi';
        } else {
          detectedType = 'single';
        }
      }
    }

    // 4. Парсинг каждого типа
    if (detectedType === 'match') {
      let question = '';
      const pairs: ParsedQuizPair[] = [];

      for (const line of cleanLines) {
        if (line.includes('::')) {
          const parts = line.split('::');
          pairs.push({
            left: parts[0].trim(),
            right: parts.slice(1).join('::').trim(),
          });
        } else if (!question) {
          question = line.replace(/^(?:#+\s*|\d+[\.\)]\s*|Вопрос\s*\d+[\.:]?\s*)/i, '').trim();
        }
      }

      if (pairs.length < 2) {
        warnings.push('Для сопоставления требуется минимум 2 пары (Левая часть :: Правая часть).');
      }

      parsedQuizzes.push({
        id: `quiz-match-${Date.now()}-${idx}`,
        rawIndex: idx + 1,
        type: 'match',
        question: question || `Сопоставление терминов ${idx + 1}`,
        points: explicitPoints ?? 15,
        pairs,
        warnings,
      });
      return;
    }

    if (detectedType === 'dropdown') {
      let question = '';
      let templateText = '';

      // Если в первой строке нет фигурных скобок, считаем её вопросом
      if (!/\{[^{}]+\}/.test(cleanLines[0]) && cleanLines.length > 1) {
        question = cleanLines[0].replace(/^(?:#+\s*|\d+[\.\)]\s*|Вопрос\s*\d+[\.:]?\s*)/i, '').trim();
        templateText = cleanLines.slice(1).join('\n');
      } else {
        templateText = cleanLines.join('\n');
        question = 'Заполните пропуски в предложении';
      }

      const blanks = templateText.match(/\{[^{}]+\}/g) || [];
      if (blanks.length === 0) {
        warnings.push('Не найдено ни одного пропуска с вариантами в формате {Правильный; Ошибка1, Ошибка2}.');
      }

      parsedQuizzes.push({
        id: `quiz-dropdown-${Date.now()}-${idx}`,
        rawIndex: idx + 1,
        type: 'dropdown',
        question,
        points: explicitPoints ?? 10,
        templateText,
        blanksCount: blanks.length,
        warnings,
      });
      return;
    }

    if (detectedType === 'input') {
      let question = '';
      let templateText = '';

      if (!/\[([^\]]+)\]/.test(cleanLines[0]) && cleanLines.length > 1) {
        question = cleanLines[0].replace(/^(?:#+\s*|\d+[\.\)]\s*|Вопрос\s*\d+[\.:]?\s*)/i, '').trim();
        templateText = cleanLines.slice(1).join('\n');
      } else {
        templateText = cleanLines.join('\n');
        question = 'Вставьте пропущенные термины';
      }

      const blanks = (templateText.match(/\[([^\]]+)\]/g) || []).filter(
        (b) => !/^\[[xX\s]\]$/.test(b)
      );
      if (blanks.length === 0) {
        warnings.push('Не найдено ни одного поля для ручного ввода в формате [термин].');
      }

      parsedQuizzes.push({
        id: `quiz-input-${Date.now()}-${idx}`,
        rawIndex: idx + 1,
        type: 'input',
        question,
        points: explicitPoints ?? 10,
        templateText,
        blanksCount: blanks.length,
        warnings,
      });
      return;
    }

    if (detectedType === 'sequence') {
      let question = '';
      const items: ParsedQuizSequenceItem[] = [];

      for (const line of cleanLines) {
        const arrowMatch = line.match(/^(?:->|=>)\s+(.+)$/);
        const numMatch = line.match(/^\d+[\.\)]\s+(.+)$/);

        if (arrowMatch) {
          items.push({ text: arrowMatch[1].trim() });
        } else if (items.length > 0 && numMatch) {
          items.push({ text: numMatch[1].trim() });
        } else if (!question) {
          question = line.replace(/^(?:#+\s*|\d+[\.\)]\s*|Вопрос\s*\d+[\.:]?\s*)/i, '').trim();
        } else {
          // Если вопрос уже есть, а стрелочки нет
          items.push({ text: line.replace(/^(?:->|=>|\d+[\.\)])\s*/, '').trim() });
        }
      }

      if (items.length < 2) {
        warnings.push('Для последовательности нужно минимум 2 шага (-> Шаг 1).');
      }

      parsedQuizzes.push({
        id: `quiz-sequence-${Date.now()}-${idx}`,
        rawIndex: idx + 1,
        type: 'sequence',
        question: question || `Расставьте шаги в правильном порядке`,
        points: explicitPoints ?? 15,
        items,
        warnings,
      });
      return;
    }

    if (detectedType === 'essay') {
      let question = '';
      let rubric = '';
      let sampleAnswer = '';

      for (const line of cleanLines) {
        const rMatch = line.match(/^(?:РУБРИКА|RUBRIC|КРИТЕРИИ)[\s:]+(.+)$/i);
        const sMatch = line.match(/^(?:ОТВЕТ|SAMPLE|ЭТАЛОН)[\s:]+(.+)$/i);

        if (rMatch) {
          rubric = rMatch[1].trim();
        } else if (sMatch) {
          sampleAnswer = sMatch[1].trim();
        } else if (!question) {
          question = line.replace(/^(?:#+\s*|\d+[\.\)]\s*|Вопрос\s*\d+[\.:]?\s*)/i, '').trim();
        } else {
          question += ' ' + line;
        }
      }

      parsedQuizzes.push({
        id: `quiz-essay-${Date.now()}-${idx}`,
        rawIndex: idx + 1,
        type: 'essay',
        question: question || `Развернутый ответ на вопрос ${idx + 1}`,
        points: explicitPoints ?? 25,
        rubric: rubric || 'Полнота раскрытия темы, аргументация и практические примеры.',
        sampleAnswer,
        warnings,
      });
      return;
    }

    if (detectedType === 'file') {
      let title = '';
      let instructions = '';
      let allowedTypes = '.zip, .pdf, .docx';
      let maxSizeMB = 25;

      for (const line of cleanLines) {
        const iMatch = line.match(/^(?:ИНСТРУКЦИЯ|INSTRUCTION|ОПИСАНИЕ)[\s:]+(.+)$/i);
        const fMatch = line.match(/^(?:ФОРМАТЫ|FORMATS)[\s:]+(.+)$/i);
        const lMatch = line.match(/^(?:ЛИМИТ|LIMIT|MAX_MB)[\s:]+(\d+)/i);

        if (iMatch) {
          instructions = iMatch[1].trim();
        } else if (fMatch) {
          allowedTypes = fMatch[1].trim();
        } else if (lMatch) {
          maxSizeMB = parseInt(lMatch[1], 10);
        } else if (!title) {
          title = line.replace(/^(?:#+\s*|\d+[\.\)]\s*|Задание\s*\d+[\.:]?\s*)/i, '').trim();
        } else {
          if (instructions) instructions += ' ' + line;
          else instructions = line;
        }
      }

      parsedQuizzes.push({
        id: `quiz-file-${Date.now()}-${idx}`,
        rawIndex: idx + 1,
        type: 'file',
        question: title || `Практическое задание с загрузкой файла ${idx + 1}`,
        fileTitle: title || `Практическое задание с загрузкой файла`,
        instructions: instructions || 'Прикрепите файл выполненной работы для проверки преподавателем.',
        allowedTypes,
        maxSizeMB,
        points: explicitPoints ?? 50,
        warnings,
      });
      return;
    }

    // 5. Обработка single и multi
    let questionText = '';
    const rawOptions: { text: string; isCorrect: boolean; key?: string }[] = [];
    let answerKey = '';
    const questionLines: string[] = [];
    let foundFirstOption = false;

    for (let i = 0; i < cleanLines.length; i++) {
      const line = cleanLines[i];

      const answerMatch = line.match(/^(?:ANSWER|ОТВЕТ|ПРАВИЛЬНЫЙ\s+ОТВЕТ)[\s:]+(.+)$/i);
      if (answerMatch) {
        answerKey = answerMatch[1].trim().toUpperCase();
        continue;
      }

      if (questionLines.length === 0) {
        questionLines.push(line);
        continue;
      }

      const checkboxMatch = line.match(/^\[([xX\s])\]\s+(.+)$/);
      if (checkboxMatch) {
        foundFirstOption = true;
        const isCorrect = checkboxMatch[1].toLowerCase() === 'x';
        rawOptions.push({
          text: checkboxMatch[2].trim(),
          isCorrect,
        });
        continue;
      }

      const letterOptionMatch = line.match(/^([a-zA-Zа-яА-Я])[\.\)]\s+(.+)$/);
      if (letterOptionMatch) {
        foundFirstOption = true;
        const key = letterOptionMatch[1].toUpperCase();
        let optText = letterOptionMatch[2].trim();
        let isCorrect = false;

        if (optText.endsWith('*') || optText.startsWith('*')) {
          isCorrect = true;
          optText = optText.replace(/^\*|\*$/g, '').trim();
        } else if (/\((?:верно|правильно|правильный|correct)\)/i.test(optText)) {
          isCorrect = true;
          optText = optText.replace(/\((?:верно|правильно|правильный|correct)\)/gi, '').trim();
        } else if (optText.endsWith('+')) {
          isCorrect = true;
          optText = optText.replace(/\+$/, '').trim();
        }

        rawOptions.push({
          text: optText,
          isCorrect,
          key,
        });
        continue;
      }

      const numberOptionMatch = line.match(/^(\d+)[\.\)]\s+(.+)$/);
      if (numberOptionMatch && (foundFirstOption || numberOptionMatch[1] === '1')) {
        foundFirstOption = true;
        const key = numberOptionMatch[1];
        let optText = numberOptionMatch[2].trim();
        let isCorrect = false;

        if (optText.endsWith('*') || optText.startsWith('*')) {
          isCorrect = true;
          optText = optText.replace(/^\*|\*$/g, '').trim();
        } else if (/\((?:верно|правильно|правильный|correct)\)/i.test(optText)) {
          isCorrect = true;
          optText = optText.replace(/\((?:верно|правильно|правильный|correct)\)/gi, '').trim();
        } else if (optText.endsWith('+')) {
          isCorrect = true;
          optText = optText.replace(/\+$/, '').trim();
        }

        rawOptions.push({
          text: optText,
          isCorrect,
          key,
        });
        continue;
      }

      const dashOptionMatch = line.match(/^\-\s+(.+)$/);
      if (dashOptionMatch && (foundFirstOption || line.startsWith('- '))) {
        foundFirstOption = true;
        let optText = dashOptionMatch[1].trim();
        let isCorrect = false;
        if (optText.endsWith('*') || optText.startsWith('*')) {
          isCorrect = true;
          optText = optText.replace(/^\*|\*$/g, '').trim();
        }
        rawOptions.push({
          text: optText,
          isCorrect,
        });
        continue;
      }

      if (!foundFirstOption) {
        questionLines.push(line);
      } else {
        if (rawOptions.length > 0) {
          rawOptions[rawOptions.length - 1].text += ' ' + line;
        }
      }
    }

    questionText = questionLines
      .join(' ')
      .replace(/^(?:#+\s*|Вопрос\s*\d+[\.:]?\s*|\d+[\.\)]\s*)/i, '')
      .trim();

    if (answerKey) {
      const keys = answerKey.split(/[,;\s]+/).filter(Boolean);
      rawOptions.forEach((opt) => {
        if (opt.key && keys.includes(opt.key)) {
          opt.isCorrect = true;
        }
      });
    }

    if (!questionText) {
      questionText = `Вопрос ${idx + 1}`;
      warnings.push('Формулировка вопроса была пустой или не распознана.');
    }

    if (rawOptions.length < 2) {
      warnings.push('Менее 2 вариантов ответов. Вопрос может быть неполным.');
    }

    const correctCount = rawOptions.filter((o) => o.isCorrect).length;
    if (correctCount === 0) {
      warnings.push('Не указан ни один правильный ответ (отметьте звездочкой * или [x]).');
    }

    const finalType: 'single' | 'multi' =
      detectedType === 'multi' || correctCount > 1 ? 'multi' : 'single';

    parsedQuizzes.push({
      id: `quiz-${finalType}-${Date.now()}-${idx}`,
      rawIndex: idx + 1,
      question: questionText,
      type: finalType,
      points: explicitPoints ?? (finalType === 'multi' ? 10 : 5),
      options: rawOptions.map((o) => ({
        text: o.text,
        isCorrect: o.isCorrect,
      })),
      warnings,
    });
  });

  return parsedQuizzes;
}

/**
 * Чтение текста из загруженного файла (.txt, .md, .docx)
 */
export async function readTextFromFile(file: File): Promise<string> {
  const fileName = file.name.toLowerCase();

  if (fileName.endsWith('.docx')) {
    const arrayBuffer = await file.arrayBuffer();
    const result = await mammoth.extractRawText({ arrayBuffer });
    return result.value;
  }

  return await file.text();
}

/**
 * Преобразование распознанных квизов во все 8 типов блоков Puck Editor
 */
export function convertQuizzesToPuckBlocks(quizzes: ParsedQuiz[]): any[] {
  return quizzes.map((q, idx) => {
    const rnd = Math.random().toString(36).substring(2, 7);

    switch (q.type) {
      case 'multi':
        return {
          type: 'QuizMultiBlock',
          props: {
            id: `QuizMultiBlock-${Date.now()}-${idx}-${rnd}`,
            question: q.question,
            points: q.points || 10,
            options: (q.options || []).map((opt) => ({
              text: opt.text,
              isCorrect: opt.isCorrect ? 'true' : 'false',
            })),
          },
        };

      case 'match':
        return {
          type: 'QuizMatchBlock',
          props: {
            id: `QuizMatchBlock-${Date.now()}-${idx}-${rnd}`,
            question: q.question,
            points: q.points || 15,
            pairs: (q.pairs || []).map((p) => ({
              left: p.left,
              right: p.right,
            })),
          },
        };

      case 'dropdown':
        return {
          type: 'QuizDropdownBlankBlock',
          props: {
            id: `QuizDropdownBlankBlock-${Date.now()}-${idx}-${rnd}`,
            question: q.question,
            points: q.points || 10,
            templateText: q.templateText || '',
          },
        };

      case 'input':
        return {
          type: 'QuizInputBlankBlock',
          props: {
            id: `QuizInputBlankBlock-${Date.now()}-${idx}-${rnd}`,
            question: q.question,
            points: q.points || 10,
            templateText: q.templateText || '',
          },
        };

      case 'sequence':
        return {
          type: 'QuizSequenceBlock',
          props: {
            id: `QuizSequenceBlock-${Date.now()}-${idx}-${rnd}`,
            question: q.question,
            points: q.points || 15,
            items: (q.items || []).map((item) => ({
              text: item.text,
            })),
          },
        };

      case 'essay':
        return {
          type: 'QuizEssayBlock',
          props: {
            id: `QuizEssayBlock-${Date.now()}-${idx}-${rnd}`,
            question: q.question,
            points: q.points || 25,
            rubric: q.rubric || 'Полнота и глубина ответа, практические примеры.',
            sampleAnswer: q.sampleAnswer || '',
          },
        };

      case 'file':
        return {
          type: 'FileUploadBlock',
          props: {
            id: `FileUploadBlock-${Date.now()}-${idx}-${rnd}`,
            title: q.fileTitle || q.question || 'Практическое задание',
            instructions: q.instructions || 'Прикрепите файл решения для проверки.',
            points: q.points || 50,
            allowedTypes: q.allowedTypes || '.zip, .pdf, .docx',
            maxSizeMB: q.maxSizeMB || 25,
          },
        };

      case 'single':
      default:
        return {
          type: 'QuizSingleBlock',
          props: {
            id: `QuizSingleBlock-${Date.now()}-${idx}-${rnd}`,
            question: q.question,
            points: q.points || 5,
            options: (q.options || []).map((opt) => ({
              text: opt.text,
              isCorrect: opt.isCorrect ? 'true' : 'false',
              explain: opt.isCorrect ? 'Правильный ответ!' : undefined,
            })),
          },
        };
    }
  });
}
