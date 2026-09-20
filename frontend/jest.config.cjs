/** @type {import('jest').Config} */
module.exports = {
  testEnvironment: 'jsdom',
  roots: ['<rootDir>/src'],
  setupFiles: ['<rootDir>/jest.setup.cjs'],
  testMatch: ['**/*.test.ts', '**/*.test.tsx'],

  // типы проверяет `yarn typecheck`, тесты только транспилируются: так быстрее и без дублирования ошибок
  transform: {
    '^.+\\.tsx?$': ['ts-jest', { tsconfig: '<rootDir>/tsconfig.jest.json', diagnostics: false }],
  },

  // те же алиасы, что в tsconfig.paths.json
  moduleNameMapper: {
    '^@app$': '<rootDir>/src/app',
    '^@app/(.*)$': '<rootDir>/src/app/$1',
    '^@shared/(.*)$': '<rootDir>/src/shared/$1',
    '^@features/auth/(.*)$': '<rootDir>/src/features/auth/$1',
    '^@features/curriculumVitae/(.*)$': '<rootDir>/src/features/curriculumVitae/$1',
  },
}
