import { SectionId, SectionSettings } from '@features/curriculumVitae/types'

export const DEFAULT_SECTIONS: SectionSettings[] = Object.values(SectionId).map((id) => ({
  id,
  enabled: true,
}))
