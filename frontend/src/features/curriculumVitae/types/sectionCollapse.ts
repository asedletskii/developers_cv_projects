import { SectionId } from './enums'

export type SectionCollapseValue = {
  isCollapsed(id: SectionId): boolean
  toggle(id: SectionId): void
  expand(id: SectionId): void
}
