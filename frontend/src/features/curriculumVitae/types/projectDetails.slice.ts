import { SliceCaseReducers } from '@shared/types'
import { Project } from './curriculumVitae'

export type ProjectDetailsSliceState = {
  project: Project
  isLoading: boolean
}

export type ProjectDetailsSliceReducers = SliceCaseReducers<ProjectDetailsSliceState> & {
  clearProjectDetailsState(state: ProjectDetailsSliceState): void
}
