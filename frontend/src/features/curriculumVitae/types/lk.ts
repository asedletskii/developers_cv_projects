/** Личный кабинет (`/lk/*`, контракт — docs/api). */

/** Языки самого резюме (поле `language`), не путать с LanguageType — языками, которыми владеет человек */
export const RESUME_LANGUAGES = ['ru', 'en', 'fr', 'de', 'kz', 'zh', 'es'] as const
export type ResumeLanguage = (typeof RESUME_LANGUAGES)[number]

/** Одна строка в списке резюме пользователя */
export type LkResumeSummary = {
  id: string
  title: string
  language: ResumeLanguage
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

/** Секции, у которых на резюме ровно одна запись (`overview` есть всегда, остальные создаются POST-ом) */
export type SingletonSection = 'overview' | 'profile' | 'contacts' | 'skills'

/** Секции-списки: `/lk/resume/:id/<kind>` и `/:itemId` */
export type ListSection =
  | 'languages'
  | 'technologies'
  | 'experiences'
  | 'educations'
  | 'courses'
  | 'activities'
  | 'projects'
