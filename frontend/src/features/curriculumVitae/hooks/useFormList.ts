import { Resume } from '@features/curriculumVitae/types'
import { useFieldArray, useFormContext } from 'react-hook-form'

type ListName = 'educations' | 'experiences' | 'courses' | 'activities' | 'projects'

/**
 * Список записей в форме редактора (опыт, образование, курсы, …): добавление пустой записи,
 * удаление и перестановка. `itemProps(index)` — готовые свойства карточки записи.
 */
export const useFormList = <Name extends ListName>(
  name: Name,
  createItem: () => Resume[Name][number],
) => {
  const { control, register } = useFormContext<Resume>()
  const { fields, append, remove, move } = useFieldArray({ control, name: name as 'educations' })

  const addItem = () => append(createItem() as Resume['educations'][number])

  const itemProps = (index: number) => ({
    canMoveUp: index > 0,
    canMoveDown: index < fields.length - 1,
    onMoveUp: () => move(index, index - 1),
    onMoveDown: () => move(index, index + 1),
    onRemove: () => remove(index),
  })

  return { control, register, fields: fields as { id: string }[], addItem, itemProps }
}
