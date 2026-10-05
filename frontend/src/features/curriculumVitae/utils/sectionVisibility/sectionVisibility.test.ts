import { isSectionVisibleEnough, pickMostVisibleSection } from './sectionVisibility'

describe('isSectionVisibleEnough', () => {
  it('достаточно 40 % высоты секции', () => {
    expect(isSectionVisibleEnough(0.4, 100, 800)).toBe(true)
    expect(isSectionVisibleEnough(0.39, 100, 800)).toBe(false)
  })

  it('у высокой секции хватает половины экрана', () => {
    expect(isSectionVisibleEnough(0.1, 400, 800)).toBe(true)
    expect(isSectionVisibleEnough(0.1, 399, 800)).toBe(false)
  })
})

describe('pickMostVisibleSection', () => {
  it('выбирает секцию с наибольшей видимой частью', () => {
    const rects = [
      { id: 'about', top: -500, bottom: 100 },
      { id: 'experience', top: 100, bottom: 900 },
    ]

    expect(pickMostVisibleSection(rects, 800)).toBe('experience')
  })

  it('учитывает только часть, попавшую в экран', () => {
    const rects = [
      { id: 'tall', top: -2000, bottom: 300 },
      { id: 'next', top: 300, bottom: 1000 },
    ]

    expect(pickMostVisibleSection(rects, 800)).toBe('next')
  })

  it('без секций на экране — undefined', () => {
    expect(pickMostVisibleSection([{ id: 'x', top: 900, bottom: 1200 }], 800)).toBeUndefined()
    expect(pickMostVisibleSection([], 800)).toBeUndefined()
  })
})
