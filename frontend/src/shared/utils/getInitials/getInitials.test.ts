import { getInitials } from './getInitials'

describe('getInitials', () => {
  it('берёт первые буквы двух слов', () => {
    expect(getInitials('Седлецкий Антон')).toBe('СА')
  })

  it('не падает на пустой строке и лишних пробелах', () => {
    expect(getInitials('')).toBe('')
    expect(getInitials('  anton   sedletskii  ')).toBe('AS')
  })
})
