import { isHexColor } from './isHexColor'

describe('isHexColor', () => {
  it.each(['#fff', '#FFF', '#0af', '#2E9E6B', '#2e9e6b'])('принимает %s', (value) => {
    expect(isHexColor(value)).toBe(true)
  })

  it.each([
    'fff',
    '#ff',
    '#ffff',
    '#fffff',
    '#12345g',
    'red',
    'rgb(1, 2, 3)',
    '#2E9E6B ',
    '',
    'javascript:alert(1)',
  ])('отклоняет %p', (value) => {
    expect(isHexColor(value)).toBe(false)
  })

  it.each([null, undefined, 123, {}, []])('отклоняет не строку %p', (value) => {
    expect(isHexColor(value)).toBe(false)
  })
})
