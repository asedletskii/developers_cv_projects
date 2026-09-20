import { SectionId } from '@features/curriculumVitae/types'
import { withPinnedAbout } from './withPinnedAbout'

describe('withPinnedAbout', () => {
  it('ставит «Обо мне» первой и включённой, что бы ни пришло', () => {
    const result = withPinnedAbout([
      { id: SectionId.Skills, enabled: true },
      { id: SectionId.About, enabled: false },
    ])

    expect(result).toEqual([
      { id: SectionId.About, enabled: true },
      { id: SectionId.Skills, enabled: true },
    ])
  })

  it('добавляет «Обо мне», если её не было', () => {
    expect(withPinnedAbout([])).toEqual([{ id: SectionId.About, enabled: true }])
  })

  it('вид, выбранный для «Обо мне», сохраняется', () => {
    expect(withPinnedAbout([{ id: SectionId.About, enabled: false, style: 'third' }])).toEqual([
      { id: SectionId.About, enabled: true, style: 'third' },
    ])
  })
})
