import { Technology } from './curriculumVitae'

export type TechnologyFilterValue = {
  selected: string | null
  isSelected(name: string): boolean
  toggle(name: string): void
  clear(): void
  isDimmed(stack: Technology[]): boolean
}
