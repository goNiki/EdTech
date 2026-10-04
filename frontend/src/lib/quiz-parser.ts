import mammoth from 'mammoth';

export interface ParsedQuizOption {
  text: string;
  isCorrect: boolean;
  explain?: string;
}

export interface ParsedQuiz {
  id: string;
  rawIndex: number;
  question: string;
  type: 'single' | 'multi';
  points: number;
  options: ParsedQuizOption[];
  warnings: string[];
}

export const SAMPLE_QUIZ_TEMPLATES = [
  {
    name: 'Aiken / Moodle (стандартный)',
    text: `1. Какой паттерн обеспечивает чистую изоляцию бизнес-логики от БД и фреймворков?
A. Active Record
B. Clean Architecture
C. Model-View-Controller
D. Monolithic Core
ANSWER: B

2. Какие протоколы передачи данных поддерживаются в современных распределенных системах?
A. HTTP/2 и HTTP/3
B. gRPC (Protocol Buffers)
C. FTP для стриминга
D. WebSockets
ANSWER: A, B, D`,
  },
  {
    name: 'Звездочки (*)',
    text: `1. Что такое идемпотентность HTTP-метода?
a) Возможность вызывать метод параллельно без задержек
b) Повторный вызов метода приводит к тому же состоянию системы *
c) Обязательное шифрование тела запроса
d) Ответ сервера всегда со статусом 200

2. Выберите принципы SOLID:
a) Single Responsibility *
b) Open-Closed Principle *
c) Over-Engineering
d) Dependency Inversion *`,
  },
  {
    name: 'Markdown чекбоксы ([x])',
    text: `1. Какая команда Git создает новую ветку и сразу переключается на нее?
[ ] git branch -d new-feature
[x] git checkout -b new-feature
[ ] git push origin new-feature
[ ] git reset --hard

2. Какие типы данных являются ссылочными в языке Go?
[x] Слайсы (Slices)
[x] Мапы (Maps)
[ ] Числовые типы (int64)
[x] Каналы (Channels)`,
  },
];

/**
 * Парсер сырого текста в структурированные вопросы квиза
 */
export function parseQuizzesFromText(rawText: string): ParsedQuiz[] {
  if (!rawText.trim()) return [];

  // Разбиваем текст на блоки вопросов (разделенные пустой строкой или нумерацией)
  const lines = rawText.split(/\r?\n/).map((l) => l.trim());
  const blocks: string[][] = [];
  let currentBlock: string[] = [];

  for (const line of lines) {
    if (!line) {
      if (currentBlock.length > 0) {
        blocks.push(currentBlock);
        currentBlock = [];
      }
      continue;
    }

    // Если начинается новый вопрос (например "1. " или "Вопрос 2:" или "# Вопрос"), и текущий блок уже содержит варианты
    const isNewQuestionStart = /^(?:#+\s*|\d+[\.\)]\s*|Вопрос\s*\d+[\.:]?\s*)/i.test(line);
    const hasOptionsInCurrent = currentBlock.some((l) =>
      /^(?:[a-zA-Zа-яА-Я][\.\)]|\[[xX\s]\]|\*|\-)/.test(l) || /^ANSWER:/i.test(l)
    );

    if (isNewQuestionStart && currentBlock.length > 0 && hasOptionsInCurrent) {
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

    let questionText = '';
    const rawOptions: { text: string; isCorrect: boolean; key?: string }[] = [];
    let answerKey = '';
    const warnings: string[] = [];

    // Ищем вопрос и строки ответов
    const questionLines: string[] = [];
    let foundFirstOption = false;

    for (let i = 0; i < block.length; i++) {
      const line = block[i];

      // Проверка на строку ответа вида "ANSWER: B" или "Ответ: А, Б"
      const answerMatch = line.match(/^(?:ANSWER|ОТВЕТ|ПРАВИЛЬНЫЙ\s+ОТВЕТ)[\s:]+(.+)$/i);
      if (answerMatch) {
        answerKey = answerMatch[1].trim().toUpperCase();
        continue;
      }

      // Проверка на Markdown чекбокс: [x] или [ ]
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

      // Проверка на стандартный вариант: "A. Текст", "1) Текст", "- Текст"
      const optionMatch = line.match(/^([a-zA-Zа-яА-Я\d])[\.\)]\s+(.+)$/);
      if (optionMatch) {
        foundFirstOption = true;
        const key = optionMatch[1].toUpperCase();
        let optText = optionMatch[2].trim();
        let isCorrect = false;

        // Проверяем звездочку или маркер (верно)
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

      // Если варианты еще не начались, это часть формулировки вопроса
      if (!foundFirstOption) {
        questionLines.push(line);
      } else {
        // Дополнительный текст к последнему варианту или не распознанная строка
        if (rawOptions.length > 0) {
          rawOptions[rawOptions.length - 1].text += ' ' + line;
        }
      }
    }

    questionText = questionLines.join(' ').replace(/^\d+[\.\)]\s*/, '').trim();

    // Если был указан ANSWER: B или ANSWER: A, C, расставляем флаги
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

    const quizType: 'single' | 'multi' = correctCount > 1 ? 'multi' : 'single';

    parsedQuizzes.push({
      id: `quiz-imported-${Date.now()}-${idx}`,
      rawIndex: idx + 1,
      question: questionText,
      type: quizType,
      points: quizType === 'multi' ? 10 : 5,
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
 * Преобразование распознанных квизов в валидные блоки Puck Editor
 */
export function convertQuizzesToPuckBlocks(quizzes: ParsedQuiz[]): any[] {
  return quizzes.map((q, idx) => {
    const uniqueId = `${q.type === 'multi' ? 'QuizMultiBlock' : 'QuizSingleBlock'}-${Date.now()}-${idx}-${Math.random().toString(36).substr(2, 5)}`;

    if (q.type === 'multi') {
      return {
        type: 'QuizMultiBlock',
        props: {
          id: uniqueId,
          question: q.question,
          points: q.points || 10,
          options: q.options.map((opt) => ({
            text: opt.text,
            isCorrect: opt.isCorrect ? 'true' : 'false',
          })),
        },
      };
    }

    return {
      type: 'QuizSingleBlock',
      props: {
        id: uniqueId,
        question: q.question,
        points: q.points || 5,
        options: q.options.map((opt) => ({
          text: opt.text,
          isCorrect: opt.isCorrect ? 'true' : 'false',
          explain: opt.isCorrect ? 'Правильный ответ!' : undefined,
        })),
      },
    };
  });
}
