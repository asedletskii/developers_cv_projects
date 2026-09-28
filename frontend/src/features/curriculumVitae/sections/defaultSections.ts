import { SectionId, SectionSettings } from '@features/curriculumVitae/types'

/**
 * Порядок секций по умолчанию, если API не прислал настройки — совпадает с порядком объявления
 * `SectionId`. Вынесено из `sectionDefinitions.tsx` в свой файл: там же лежат JSX-рендеры всех
 * секций страницы, а этот порядок нужен и вне контекста рендера (форма редактора, её тесты) —
 * незачем тащить туда компоненты ради одного массива id.
 */
export const DEFAULT_SECTIONS: SectionSettings[] = Object.values(SectionId).map((id) => ({
  id,
  enabled: true,
}))
