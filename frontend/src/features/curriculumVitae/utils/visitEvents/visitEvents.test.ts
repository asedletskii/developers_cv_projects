import { isVisibleAction } from './visitEvents'

describe('isVisibleAction', () => {
  it('действия посетителя показываем', () => {
    ;[
      'contact_view',
      'contact_copy',
      'contact_click',
      'tech_filter',
      'download',
      'print',
      'external_click',
    ].forEach((type) => expect(isVisibleAction({ type })).toBe(true))
  })

  it('служебные и неизвестные события — нет', () => {
    expect(isVisibleAction({ type: 'section_view' })).toBe(false)
    expect(isVisibleAction({ type: 'something_new' })).toBe(false)
  })
})
