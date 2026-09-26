import { SectionId, SectionSettings } from '@features/curriculumVitae/types'

export const withPinnedAbout = (sections: SectionSettings[]): SectionSettings[] => [
  { id: SectionId.About, enabled: true },
  ...sections.filter(({ id }) => id !== SectionId.About),
]
