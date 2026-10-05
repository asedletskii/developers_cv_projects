import { CurriculumVitaeSliceState } from '@features/curriculumVitae/types'

export const getCurriculumVitaeInitialState = (): CurriculumVitaeSliceState => {
  return {
    resume: {
      id: '',
      title: '',
      language: 'ru',
      overview: {
        greeting: { first: '', second: '', secondHighlight: '', third: '' },
        sections: [],
      },
      languages: [],
      technologies: [],
      experiences: [],
      educations: [],
      courses: [],
      activities: [],
      projects: [],
    },
    viewId: null,
    isLoading: true,
    error: null,
    currentRequestId: null,
  }
}
