import { ResumeContacts, ResumeProfile, ResumeSkills } from '@features/curriculumVitae/types'

/**
 * Значения отсутствующих секций резюме: `profile`, `contacts` и `skills` у нового резюме не
 * приходят. Компоненты читают `resume.profile ?? EMPTY_PROFILE` и своих значений по умолчанию
 * не заводят.
 */
export const EMPTY_PROFILE: ResumeProfile = {
  name: '',
  surname: '',
  position: '',
  photoPath: '',
  dateOfBirth: '',
  location: '',
  relocation: false,
  workFormat: [],
}

export const EMPTY_CONTACTS: ResumeContacts = { hiddenContacts: [] }

export const EMPTY_SKILLS: ResumeSkills = { hardSkills: [], softSkills: [] }
