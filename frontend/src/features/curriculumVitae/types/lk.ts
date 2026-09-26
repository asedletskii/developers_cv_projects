import { Availability, Greeting, SectionSettings, ThemeSettings } from './curriculumVitae'

/** Личный кабинет (`/lk/*`, контракт — docs/lk-api.md). */

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

/** Вложенные элементы ответа `GET /lk/resume/:id`: у каждого есть свой `id` для PUT/DELETE */
type WithId = { id: string }
type ViewTechnology = { name: string; grade?: string }

/** Резюме владельца целиком — каждая секция редактируется своей ручкой (см. docs/lk-api.md) */
export type ResumeView = LkResumeSummary & {
  createdAt: string
  updatedAt: string
  overview: {
    theme?: ThemeSettings
    greeting: Greeting
    availability?: Availability
    resumeLink?: string
    sections: SectionSettings[]
  }
  profile?: {
    name: string
    surname: string
    patronymic?: string
    position: string
    photoPath: string
    dateOfBirth: string
    location: string
    relocation: boolean
    workFormat: string[]
  }
  contacts?: {
    phone?: string
    telegram?: string
    email?: string
    hiddenContacts: string[]
    contactsText?: { title?: string; subtitle?: string }
  }
  skills?: { hardSkills: string[]; softSkills: string[] }
  languages: (WithId & { language: string; grade?: string })[]
  technologies: (WithId & ViewTechnology)[]
  experiences: (WithId & {
    company: string
    position: string
    location?: string
    startDate: string
    endDate?: string
    description?: string
    achievements: string[]
    technologies: ViewTechnology[]
  })[]
  educations: (WithId & {
    location: string
    university: string
    grade: number
    programName: string
    profileName?: string
    faculty: string
    course?: number
    endDate?: string
  })[]
  courses: (WithId & {
    name: string
    organization: string
    endDate: string
    hours?: number
    certificateLink?: string
  })[]
  activities: (WithId & {
    type: string
    title: string
    description?: string
    date?: string
    link?: string
  })[]
  projects: (WithId & {
    name: string
    description?: string
    role?: string
    screenshots: string[]
    link?: string
    gitLink?: string
    stack: ViewTechnology[]
    contributors: {
      name: string
      surname: string
      role: string
      contacts: { phone?: string; telegram?: string; email?: string }
    }[]
  })[]
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
