import { Person } from '@features/curriculumVitae/types'
import { FieldPath, useController, useFormContext } from 'react-hook-form'

type Option = { value: string; label: string }

/** Поле-список навыков: выбранные значения формы ↔ варианты справочника. */
export const useSkillsField = (
  name: Extract<FieldPath<Person>, 'hardSkills' | 'softSkills'>,
  options: Option[],
) => {
  const { control } = useFormContext<Person>()
  const { field } = useController({ control, name })

  const values = (field.value ?? []) as string[]

  return {
    selected: options.filter((option) => values.includes(option.value)),
    onChange: (selected: Option[]) => field.onChange(selected.map((option) => option.value)),
    onBlur: field.onBlur,
  }
}
