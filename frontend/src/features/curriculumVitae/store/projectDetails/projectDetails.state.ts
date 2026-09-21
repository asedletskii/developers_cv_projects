import { ProjectDetailsSliceState } from '@features/curriculumVitae/types'

export function getProjectDetailsInitialState(): ProjectDetailsSliceState {
  return {
    project: {
      id: '',
      name: '',
      stack: [],
      contributors: [],
    },
    isLoading: true,
  }
}
