import { Person, SectionId } from '@features/curriculumVitae/types'
import { useFieldArray, useFormContext } from 'react-hook-form'

/**
 * Настройки секций: порядок массива `sections` = порядок на странице. «Обо мне» закреплена: первая,
 * без переключателей, остальные выше неё не поднимаются.
 *
 * `keyName: '_key'` — без этого RHF использует поле `id` самих данных (`SectionSettings.id`,
 * это `SectionId`, а не случайный ключ) как свой служебный React-key и затирает его случайным.
 */
export const useSectionsEditor = () => {
  const { control } = useFormContext<Person>()
  const { fields, move, update } = useFieldArray({ control, name: 'sections', keyName: '_key' })

  return {
    fields,
    isPinned: (id: SectionId) => id === SectionId.About,
    canMoveUp: (index: number) => index > 1,
    canMoveDown: (index: number) => index > 0 && index < fields.length - 1,
    moveUp: (index: number) => move(index, index - 1),
    moveDown: (index: number) => move(index, index + 1),
    toggleEnabled: (index: number) => {
      const { _key, ...section } = fields[index]
      update(index, { ...section, enabled: !section.enabled })
    },
    toggleCollapsed: (index: number) => {
      const { _key, ...section } = fields[index]
      update(index, { ...section, collapsed: !section.collapsed })
    },
  }
}
