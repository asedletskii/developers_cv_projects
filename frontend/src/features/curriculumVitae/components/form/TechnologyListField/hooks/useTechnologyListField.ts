import { Person, Technology } from '@features/curriculumVitae/types'
import { isSameTechnology } from '@features/curriculumVitae/utils'
import { useController, useFormContext } from 'react-hook-form'

type Name = 'technologies' | `experience.${number}.stack` | `projects.${number}.stack`

/** Поле-список технологий: значения формы ↔ названия в выпадающем списке. */
export const useTechnologyListField = (name: Name) => {
  const { control } = useFormContext<Person>()
  const { field } = useController({ control, name })

  const value = (field.value ?? []) as Technology[]

  const onChange = (names: string[]) => {
    const unique = names
      .map((item) => item.trim())
      .filter(
        (item, index, all) =>
          item && all.findIndex((other) => isSameTechnology(other, item)) === index,
      )

    // Существующие записи оставляем как есть (в них может быть уровень), новые — по имени
    field.onChange(unique.map((item) => value.find((tech) => tech.name === item) ?? { name: item }))
  }

  return { names: value.map(({ name: technology }) => technology), onChange, onBlur: field.onBlur }
}
