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
import { LkResumeSummary, ResumeStatus } from './lk'

type Item<T> = T & { id?: string }

export type ResumeOverview = {
  theme?: ThemeSettings
  greeting: Greeting
  availability?: Availability
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
  hiddenContacts: ContactKey[]
  contactsText?: ContactsText
}

export type ResumeSkills = { hardSkills: HardSkill[]; softSkills: SoftSkill[] }

export type ResumeLanguageItem = Item<{ language: LanguageType; grade?: LanguageGrade }>

export type ResumeTechnologyItem = Item<Technology>

export type Resume = Omit<LkResumeSummary, 'version' | 'status'> & {
  createdAt?: string
  updatedAt?: string
  version?: number
  status?: ResumeStatus
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

/** Значение скрытого контакта, которое вернул бэкенд */
export type RevealedContact = { contact: ContactKey; value: string }
