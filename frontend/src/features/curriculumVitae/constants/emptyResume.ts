import { ResumeContacts, ResumeProfile, ResumeSkills } from '@features/curriculumVitae/types'

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
