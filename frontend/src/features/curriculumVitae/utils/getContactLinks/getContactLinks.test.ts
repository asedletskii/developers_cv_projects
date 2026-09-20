import { ContactKey } from '@features/curriculumVitae/types'
import { getContactDisplayValue, getContactHref, getContactLinks } from './getContactLinks'

// Иконки не нужны: подменяем, чтобы не загружать всю библиотеку
jest.mock('@mui/icons-material', () => ({
  Telegram: () => null,
  EmailOutlined: () => null,
  PhoneOutlined: () => null,
}))

describe('getContactLinks', () => {
  it('строит ссылки для всех заполненных контактов в фиксированном порядке', () => {
    const links = getContactLinks({
      phone: '+79999999999',
      email: 'anton@example.com',
      telegram: 'a_sedletskii',
    })

    expect(links.map(({ key }) => key)).toEqual([
      ContactKey.Telegram,
      ContactKey.Email,
      ContactKey.Phone,
    ])
  })

  it('отдаёт значение как есть, без форматирования', () => {
    const links = getContactLinks({
      telegram: 'a_sedletskii',
      email: 'anton@example.com',
      phone: '+79999999999',
    })

    expect(links.map(({ value }) => value)).toEqual([
      'a_sedletskii',
      'anton@example.com',
      '+79999999999',
    ])
  })

  it('у каждой ссылки есть иконка', () => {
    const links = getContactLinks({ telegram: 'a_sedletskii', email: 'anton@example.com' })

    expect(links.every(({ icon }) => !!icon)).toBe(true)
  })

  it('пропускает незаполненные и нескрытые контакты', () => {
    const links = getContactLinks({ email: 'anton@example.com' })

    expect(links.map(({ key }) => key)).toEqual([ContactKey.Email])
  })

  it('пустая строка считается незаполненным контактом', () => {
    expect(getContactLinks({ telegram: '', email: '', phone: '' })).toEqual([])
  })

  it('без контактов возвращает пустой список', () => {
    expect(getContactLinks({})).toEqual([])
  })

  it('скрытый контакт без значения тоже попадает в список — value придёт по клику', () => {
    const [link] = getContactLinks({}, [ContactKey.Phone])

    expect(link).toMatchObject({ key: ContactKey.Phone, isHidden: true, value: undefined })
    expect(link.icon).toBeTruthy()
  })

  it('помечает открытые контакты как нескрытые', () => {
    const [link] = getContactLinks({ email: 'anton@example.com' }, [ContactKey.Phone])

    expect(link.isHidden).toBe(false)
  })
})

describe('getContactHref: значение от автора не меняет ссылку', () => {
  it('email: параметры mailto кодируются, @ остаётся', () => {
    expect(getContactHref(ContactKey.Email, 'me@example.ru')).toBe('mailto:me@example.ru')
    expect(getContactHref(ContactKey.Email, 'a@b.ru?cc=x&body=y')).toBe(
      'mailto:a@b.ru%3Fcc%3Dx%26body%3Dy',
    )
  })

  it('телефон: только цифры и + * #', () => {
    expect(getContactHref(ContactKey.Phone, '+7 (999) 888-22-33')).toBe('tel:+79998882233')
    expect(getContactHref(ContactKey.Phone, '+7999;phone-context=evil')).toBe('tel:+7999')
  })

  it('telegram: @ срезается, путь кодируется', () => {
    expect(getContactHref(ContactKey.Telegram, '@a_b')).toBe('https://t.me/a_b')
    expect(getContactHref(ContactKey.Telegram, 'a/../b?x')).toBe('https://t.me/a%2F..%2Fb%3Fx')
  })
})

describe('getContactHref', () => {
  it('собирает адрес по ключу контакта', () => {
    expect(getContactHref(ContactKey.Telegram, 'a_sedletskii')).toBe('https://t.me/a_sedletskii')
    expect(getContactHref(ContactKey.Email, 'anton@example.com')).toBe('mailto:anton@example.com')
    expect(getContactHref(ContactKey.Phone, '+79999999999')).toBe('tel:+79999999999')
  })
})

describe('getContactDisplayValue', () => {
  it('у телеграма добавляет @, у остальных отдаёт значение как есть', () => {
    expect(getContactDisplayValue(ContactKey.Telegram, 'a_sedletskii')).toBe('@a_sedletskii')
    expect(getContactDisplayValue(ContactKey.Email, 'anton@example.com')).toBe('anton@example.com')
  })
})
