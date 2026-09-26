import { SectionId } from './enums'

/** Состояние аккордеона секций страницы. */
export type SectionCollapseValue = {
  isCollapsed(id: SectionId): boolean
  toggle(id: SectionId): void
  expand(id: SectionId): void
}
