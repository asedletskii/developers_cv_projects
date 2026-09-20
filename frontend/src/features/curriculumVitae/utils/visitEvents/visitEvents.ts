export const VISIBLE_ACTION_TYPES: ReadonlySet<string> = new Set([
  'contact_view',
  'contact_copy',
  'contact_click',
  'tech_filter',
  'download',
  'print',
  'external_click',
])

export const isVisibleAction = ({ type }: { type: string }) => VISIBLE_ACTION_TYPES.has(type)
