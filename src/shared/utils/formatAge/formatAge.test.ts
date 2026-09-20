import { formatAge } from './formatAge'

describe('formatAge', () => {
  beforeEach(() => {
    jest.useFakeTimers().setSystemTime(new Date('2026-09-21T12:00:00Z'))
  })

  afterEach(() => {
    jest.useRealTimers()
  })

  it('считает полные годы', () => {
    expect(formatAge('1998-04-12')).toBe('28 years')
  })

  it('до дня рождения в этом году возраст ещё на год меньше', () => {
    expect(formatAge('1998-12-01')).toBe('27 years')
  })

  it('в день рождения возраст уже вырос', () => {
    expect(formatAge('1998-09-21')).toBe('28 years')
  })

  it('использует единственное число для одного года', () => {
    expect(formatAge('2025-09-01')).toBe('1 year')
  })
})
