import {
  Activity,
  Availability,
  Contacts,
  ContactsText,
  Course,
  Education,
  Experience,
  Greeting,
  Project,
  SectionSettings,
  Technology,
  ThemeSettings,
} from './curriculumVitae'
import { ContactKey, HardSkill, LanguageGrade, LanguageType, SoftSkill, WorkFormat } from './enums'
import { LkResumeSummary } from './lk'

/** Запись списка — у записей с бэкенда есть `id`, у только что добавленных в форме его ещё нет */
type Item<T> = T & { id?: string }

export type ResumeOverview = {
  theme?: ThemeSettings
  greeting: Greeting
  availability?: Availability
  resumeLink?: string
  sections: SectionSettings[]
}

export type ResumeProfile = {
  name: string
  surname: string
  patronymic?: string
  position: string
  photoPath: string
  dateOfBirth: string
  location: string
  relocation: boolean
  workFormat: WorkFormat[]
}

export type ResumeContacts = Contacts & {
  /** Контакты, которые показываются только по клику (защита от сборщиков почты и телефонов) */
  hiddenContacts: ContactKey[]
  contactsText?: ContactsText
}

export type ResumeSkills = { hardSkills: HardSkill[]; softSkills: SoftSkill[] }

export type ResumeLanguageItem = Item<{ language: LanguageType; grade?: LanguageGrade }>

export type ResumeTechnologyItem = Item<Technology>

/**
 * Резюме целиком: единая форма и для публичной страницы (`GET /cv/:id`), и для кабинета
 * (`GET /lk/resume/:id`), и для форм редактора. `profile`, `contacts`, `skills` у нового
 * резюме отсутствуют — читаем их через `EMPTY_*` из `constants/emptyResume.ts`.
 */
export type Resume = LkResumeSummary & {
  createdAt?: string
  updatedAt?: string
  overview: ResumeOverview
  profile?: ResumeProfile
  contacts?: ResumeContacts
  skills?: ResumeSkills
  languages: ResumeLanguageItem[]
  technologies: ResumeTechnologyItem[]
  experiences: Item<Experience>[]
  educations: Item<Education>[]
  courses: Item<Course>[]
  activities: Item<Activity>[]
  projects: Item<Project>[]
}
