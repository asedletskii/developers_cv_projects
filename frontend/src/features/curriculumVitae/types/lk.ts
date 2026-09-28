/** Личный кабинет (`/lk/*`, контракт — docs/api). */

/** Языки самого резюме (поле `language`), не путать с LanguageType — языками, которыми владеет человек */
export const RESUME_LANGUAGES = ['ru', 'en', 'fr', 'de', 'kz', 'zh', 'es'] as const
export type ResumeLanguage = (typeof RESUME_LANGUAGES)[number]

/** Статус резюме: `draft` — только что созданное, `opened` — публичная ссылка работает,
 *  `closed` — сняли с публикации. Новое резюме всегда создаётся `draft`; в `opened` перевести
 *  можно и из `draft`, и из `closed`, обратно в `draft` — нельзя (см. `Resume.Publish`/`Close` на бэке). */
export const RESUME_STATUSES = ['draft', 'opened', 'closed'] as const
export type ResumeStatus = (typeof RESUME_STATUSES)[number]

/** Одна строка в списке резюме пользователя. `version` — оптимистическая блокировка всего
 *  резюме, нужна для любого `PUT /lk/resume/:id...`. */
export type LkResumeSummary = {
  id: string
  title: string
  language: ResumeLanguage
  status: ResumeStatus
  version: number
}

/** Имя, фамилия и т.д. на одном языке интерфейса — у владельца ЛК их два набора, `ru` и `en` */
export type LkLocalizedProfile = {
  name?: string
  surname?: string
  patronymic?: string
  position?: string
  location?: string
}

/** Общие данные владельца ЛК — то, что показывается в левой колонке. При миграции на локализацию
 *  всё, что было, ушло в `ru`; `en` у существующих пользователей пока пусто. */
export type LkProfile = {
  userId: string
  ru: LkLocalizedProfile
  en: LkLocalizedProfile
  photoPath?: string
}

/** Статус ссылки: `active` — открывается, `inactive` — отключена (без удаления, можно включить
 *  обратно). `version` — оптимистическая блокировка для `PUT /lk/link/:id/status`. */
export const LINK_STATUSES = ['active', 'inactive'] as const
export type LinkStatus = (typeof LINK_STATUSES)[number]

/** Именованная ссылка на резюме — свой UUID, а не резюме с меткой, чтобы адресат не мог её
 *  подменить, а бэк по нему считал просмотры отдельно от публичного адреса. */
export type ResumeLink = {
  id: string
  name: string
  status: LinkStatus
  version: number
  createdAt: string
  updatedAt: string
}

/** Секции, у которых на резюме ровно одна запись — сохраняются целиком через `PUT` */
export type SingletonSection = 'overview' | 'profile' | 'contacts' | 'skills'

/** Секции-списки — тоже сохраняются целиком через `PUT` со всем массивом сразу */
export type ListSection =
  | 'languages'
  | 'technologies'
  | 'experiences'
  | 'educations'
  | 'courses'
  | 'activities'
  | 'projects'
