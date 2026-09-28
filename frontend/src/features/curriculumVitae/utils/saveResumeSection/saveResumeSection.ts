import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { Resume, ResumeEditorTab } from '@features/curriculumVitae/types'
import { buildSectionRequest } from './buildSectionRequest'

export type { SectionRequest } from './buildSectionRequest'
export { buildSectionRequest } from './buildSectionRequest'

/** Отправляет на бэк ровно одну секцию из значений формы — тело собирает `buildSectionRequest`. */
export const saveResumeSection = (
  userId: string,
  view: Resume,
  section: ResumeEditorTab,
  values: Resume,
): Promise<void> => {
  const { section: apiSection, body } = buildSectionRequest(view, section, values)
  return curriculumVitaeApi.replaceSection(userId, view.id, apiSection, body)
}
