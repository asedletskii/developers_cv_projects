import { toOptionalNumber, toOptionalString } from './toOptional'

describe('toOptional', () => {
  it('пустая строка и undefined — не задано', () => {
    expect(toOptionalString('')).toBeUndefined()
    expect(toOptionalString(undefined)).toBeUndefined()
    expect(toOptionalString('a')).toBe('a')
  })

  it('NaN и undefined — не задано, число сохраняется', () => {
    expect(toOptionalNumber(NaN)).toBeUndefined()
    expect(toOptionalNumber(undefined)).toBeUndefined()
    expect(toOptionalNumber(0)).toBe(0)
  })
})
