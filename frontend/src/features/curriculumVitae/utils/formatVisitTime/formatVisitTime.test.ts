import { formatVisitTime } from './formatVisitTime'

const now = new Date(2026, 8, 30, 18, 0)

describe('formatVisitTime', () => {
  it('сегодня и вчера — словами', () => {
    expect(formatVisitTime(new Date(2026, 8, 30, 9, 5).toISOString(), 'en', { now })).toBe(
      'Today, 09:05 AM',
    )
    expect(formatVisitTime(new Date(2026, 8, 29, 16, 40).toISOString(), 'en', { now })).toBe(
      'Yesterday, 04:40 PM',
    )
  })

  it('без заглавной буквы для середины фразы', () => {
    expect(
      formatVisitTime(new Date(2026, 8, 29, 16, 40).toISOString(), 'en', {
        now,
        isCapitalized: false,
      }),
    ).toBe('yesterday, 04:40 PM')
  })

  it('старше — датой', () => {
    expect(formatVisitTime(new Date(2026, 8, 12, 8, 14).toISOString(), 'en', { now })).toBe(
      'Sep 12, 08:14 AM',
    )
  })
})
