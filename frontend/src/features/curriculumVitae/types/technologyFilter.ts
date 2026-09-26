import { Technology } from './curriculumVitae'

/** Фильтр страницы по технологии. */
export type TechnologyFilterValue = {
  /** Выбранная технология или null, если фильтра нет */
  selected: string | null
  isSelected(name: string): boolean
  toggle(name: string): void
  clear(): void
  /** Нужно ли приглушить элемент со стеком: фильтр включён, а технологии в стеке нет */
  isDimmed(stack: Technology[]): boolean
}
