import { formatDuration } from './formatDuration'

describe('formatDuration', () => {
  it('секунды отдельно, если меньше минуты', () => {
    expect(formatDuration(45, 'en')).toBe('45s')
  })

  it('минуты и секунды', () => {
    expect(formatDuration(131, 'en')).toBe('2m 11s')
  })

  it('ровные минуты без секунд', () => {
    expect(formatDuration(120, 'en')).toBe('2m')
  })
})
