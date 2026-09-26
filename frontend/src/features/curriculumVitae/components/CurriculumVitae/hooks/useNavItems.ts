import { useMemo } from 'react'
import { SECTION_DEFINITIONS } from '@features/curriculumVitae/sections'
import { SectionId, SectionSettings } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'

type NavKey = Parameters<ReturnType<typeof useTranslate<'CurriculumVitae.Nav'>>>[0]

export const useNavItems = (sections: SectionSettings[]) => {
  const translate = useTranslate('CurriculumVitae.Nav')

  return useMemo(
    () =>
      sections
        .filter(({ id }) => SECTION_DEFINITIONS[id].inNav)
        .map(({ id }): { id: SectionId; label: string } => ({
          id,
          label: translate(id as NavKey),
        })),
    [sections, translate],
  )
}
