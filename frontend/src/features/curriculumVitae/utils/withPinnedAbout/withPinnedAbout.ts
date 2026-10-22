import { SectionId, SectionSettings } from '@features/curriculumVitae/types'

export const withPinnedAbout = (sections: SectionSettings[]): SectionSettings[] => {
  const about = sections.find(({ id }) => id === SectionId.About)

  return [
    { id: SectionId.About, enabled: true, ...(about?.style ? { style: about.style } : {}) },
    ...sections.filter(({ id }) => id !== SectionId.About),
  ]
}
