/** Личный кабинет (`/lk/*`, контракт — docs/lk-api.md). */

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

