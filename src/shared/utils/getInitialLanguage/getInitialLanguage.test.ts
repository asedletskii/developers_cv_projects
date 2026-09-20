import { Languages } from '@shared/types'
import { getInitialLanguage } from './getInitialLanguage'

// Словарь нужен только как список поддерживаемых языков
jest.mock('@app/locale', () => ({ appLocaleDict: { ru: {}, en: {} } }))

const setNavigatorLanguage = (value: string) =>
  Object.defineProperty(window.navigator, 'language', { value, configurable: true })

describe('getInitialLanguage', () => {
  beforeEach(() => {
    localStorage.clear()
    setNavigatorLanguage('en-US')
  })

  it('берёт сохранённый язык', () => {
    localStorage.setItem('language', 'ru')

    expect(getInitialLanguage()).toBe(Languages.Ru)
  })

  it('сохранённый язык важнее языка браузера', () => {
    localStorage.setItem('language', 'ru')
    setNavigatorLanguage('en-US')

    expect(getInitialLanguage()).toBe(Languages.Ru)
  })

  it('неподдерживаемый сохранённый язык игнорируется', () => {
    localStorage.setItem('language', 'de')
    setNavigatorLanguage('ru-RU')

    expect(getInitialLanguage()).toBe(Languages.Ru)
  })

  it('без сохранённого языка берёт язык браузера', () => {
    setNavigatorLanguage('en-GB')

    expect(getInitialLanguage(Languages.Ru)).toBe(Languages.En)
  })

  it('учитывает регистр в языке браузера', () => {
    setNavigatorLanguage('RU-ru')

    expect(getInitialLanguage(Languages.En)).toBe(Languages.Ru)
  })

  it('неподдерживаемый язык браузера даёт язык по умолчанию (русский)', () => {
    setNavigatorLanguage('de-DE')

    expect(getInitialLanguage()).toBe(Languages.Ru)
  })

  it('неподдерживаемый язык браузера даёт переданный запасной', () => {
    setNavigatorLanguage('de-DE')

    expect(getInitialLanguage(Languages.En)).toBe(Languages.En)
  })
})
