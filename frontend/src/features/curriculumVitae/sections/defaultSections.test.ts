import { SectionId } from '@features/curriculumVitae/types'
import { DEFAULT_SECTIONS } from './defaultSections'

describe('DEFAULT_SECTIONS', () => {
  it('содержит каждую секцию ровно один раз, включённой', () => {
    expect(DEFAULT_SECTIONS.map(({ id }) => id)).toEqual(Object.values(SectionId))
    expect(DEFAULT_SECTIONS.every(({ enabled }) => enabled)).toBe(true)
  })

  it('«Обо мне» идёт первой', () => {
    expect(DEFAULT_SECTIONS[0].id).toBe(SectionId.About)
  })
})
