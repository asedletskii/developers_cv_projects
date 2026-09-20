import { byStyle, resolveSectionStyle } from './sectionStyles'

describe('resolveSectionStyle', () => {
  it('известные виды возвращает как есть', () => {
    expect(resolveSectionStyle('second')).toBe('second')
    expect(resolveSectionStyle('third')).toBe('third')
  })

  it('без значения и с незнакомым значением — основной вид', () => {
    expect(resolveSectionStyle()).toBe('first')
    expect(resolveSectionStyle('')).toBe('first')
    expect(resolveSectionStyle('fourth')).toBe('first')
  })
})

describe('byStyle', () => {
  it('выбирает вариант по виду', () => {
    expect(byStyle('third', { first: 1, second: 2, third: 3 })).toBe(3)
  })
})
