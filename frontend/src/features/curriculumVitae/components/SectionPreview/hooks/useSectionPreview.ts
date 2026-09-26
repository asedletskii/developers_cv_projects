import { SECTION_DEFINITIONS } from '@features/curriculumVitae/sections'
import { Resume, SectionId } from '@features/curriculumVitae/types'
import { useFormContext, useWatch } from 'react-hook-form'

/** Данные предпросмотра секции: значения формы «как сейчас набрано» и есть ли что показывать. */
export const useSectionPreview = (section: SectionId) => {
  const { control } = useFormContext<Resume>()
  const resume = useWatch({ control }) as Resume

  const definition = SECTION_DEFINITIONS[section]

  return { resume, definition, hasContent: definition.hasContent(resume) }
}
