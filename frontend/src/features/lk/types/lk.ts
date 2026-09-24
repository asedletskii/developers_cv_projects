/**
 * Черновой контракт: бэкенд для личного кабинета ещё не существует, фронт пока работает на
 * моках (см. `api/lk.mock.ts`). Реальные `GET /api/v1/lk/profile` и `GET /api/v1/lk/resumes`
 * (userId — заголовком `X-Mock-User-ID`, не в пути) см. в docs/lk-api.md.
 */

/** Одна строка в списке резюме пользователя */
export type LkResumeSummary = {
  id: string
  title: string
}

/** Общие данные владельца ЛК — то, что показывается в левой колонке */
export type LkProfile = {
  userId: string
  name: string
  surname: string
  patronymic?: string
  position?: string
  location?: string
  photoPath?: string
}

