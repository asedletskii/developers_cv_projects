import { SECTION_DEFINITIONS } from '@features/curriculumVitae/sections'
import { Person, SectionId } from '@features/curriculumVitae/types'
import { useFormContext, useWatch } from 'react-hook-form'

/** Данные предпросмотра секции: значения формы «как сейчас набрано» и есть ли что показывать. */
export const useSectionPreview = (section: SectionId) => {
  const { control } = useFormContext<Person>()
  const person = useWatch({ control }) as Person

  const definition = SECTION_DEFINITIONS[section]

  return { person, definition, hasContent: definition.hasContent(person) }
}
