import { transliterate, slugify, isValidSlug } from './utils';

function assertEqual(actual: string | boolean, expected: string | boolean, message: string) {
  if (actual !== expected) {
    console.error(`❌ FAILED: ${message}. Expected: "${expected}", Got: "${actual}"`);
    process.exit(1);
  } else {
    console.log(`✅ PASSED: ${message}`);
  }
}

console.log('--- Testing transliterate & slugify ---');

// 1. Basic Cyrillic transliteration
assertEqual(
  transliterate('Основы Golang для начинающих'),
  'Osnovy Golang dlya nachinayuschih',
  'Transliterate Cyrillic string'
);

// 2. Slugify with Russian text and spaces
assertEqual(
  slugify('Архитектура микросервисов и Go'),
  'arhitektura-mikroservisov-i-go',
  'Slugify Russian title'
);

// 3. Punctuation removal and single hyphens
assertEqual(
  slugify('Курс: "Go с нуля!" (Полный гид, 2026?!)'),
  'kurs-go-s-nulya-polnyy-gid-2026',
  'Slugify with complex punctuation'
);

// 4. Repeated hyphens and trailing hyphens
assertEqual(
  slugify('---Тест---и---проверка---'),
  'test-i-proverka',
  'Trim leading and trailing hyphens'
);

// 5. isValidSlug verification
assertEqual(isValidSlug('arhitektura-mikroservisov-i-go'), true, 'Valid slug check');
assertEqual(isValidSlug('valid-slug-123'), true, 'Valid slug with numbers');
assertEqual(isValidSlug('Invalid Slug'), false, 'Invalid slug with spaces');
assertEqual(isValidSlug('invalid_slug'), false, 'Invalid slug with underscores');
assertEqual(isValidSlug('-leading-hyphen'), false, 'Invalid slug with leading hyphen');
assertEqual(isValidSlug('trailing-hyphen-'), false, 'Invalid slug with trailing hyphen');
assertEqual(isValidSlug('double--hyphen'), false, 'Invalid slug with double hyphen');
assertEqual(isValidSlug('русский-слаг'), false, 'Invalid slug with Cyrillic');

console.log('\n🎉 ALL UTILS TESTS PASSED SUCCESSFULLY!');
