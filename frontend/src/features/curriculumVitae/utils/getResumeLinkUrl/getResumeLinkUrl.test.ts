import { getResumeLinkUrl } from './getResumeLinkUrl'

describe('getResumeLinkUrl', () => {
  it('собирает адрес из origin и id', () => {
    expect(getResumeLinkUrl('abc', 'https://cv.example.ru')).toBe('https://cv.example.ru/cv/abc')
  })

  it('экранирует id', () => {
    expect(getResumeLinkUrl('a b/c', 'https://x.ru')).toBe('https://x.ru/cv/a%20b%2Fc')
  })
})
