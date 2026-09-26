import { useState } from 'react'
import { Person } from '@features/curriculumVitae/types'
import { FieldPath, useController, useFormContext } from 'react-hook-form'

/** Поле цвета: значение формы (или цвет по умолчанию) и состояние всплывающей палитры. */
export const useColorField = (name: FieldPath<Person>, defaultColor: string) => {
  const { control } = useFormContext<Person>()
  const { field } = useController({ control, name })
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null)

  const color = (field.value as string) || defaultColor

  return {
    color,
    anchorEl,
    openPicker: (target: HTMLElement) => setAnchorEl(target),
    closePicker: () => setAnchorEl(null),
    onChange: (next: string) => field.onChange(next),
    onBlur: field.onBlur,
  }
}
