import { CurriculumVitaeSliceState } from '@features/curriculumVitae/types'

export function getCurriculumVitaeInitialState(): CurriculumVitaeSliceState {
  return {
    person: {
      id: '',
      name: '',
      surname: '',
      position: '',
      sections: [],
      greeting: {
        first: '',
        second: '',
        secondHighlight: '',
        third: '',
      },
      photoPath: '',
      dateOfBirth: '',
      location: '',
      relocation: false,
      workFormat: [],
      contacts: {
        phone: undefined,
        telegram: undefined,
        email: undefined,
      },
      languages: [],
      experience: [],
      education: [],
      courses: [],
      activities: [],
      hardSkills: [],
      softSkills: [],
      technologies: [],
      projects: [],
    },
    isLoading: true,
    error: null,
    currentRequestId: null,
  }
}
