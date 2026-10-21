import { Resume } from '@features/curriculumVitae/types'
import { normalizeSkills } from '@features/curriculumVitae/utils'
import { FieldPath, useController, useFormContext } from 'react-hook-form'

export const useSkillsField = (
  name: Extract<FieldPath<Resume>, 'skills.hardSkills' | 'skills.softSkills'>,
) => {
  const { control } = useFormContext<Resume>()
  const { field } = useController({ control, name })

  return {
    selected: (field.value ?? []) as string[],
    onChange: (next: string[]) => field.onChange(normalizeSkills(next)),
    onBlur: field.onBlur,
  }
}
