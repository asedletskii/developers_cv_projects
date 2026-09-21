import { ContactKey } from '@features/curriculumVitae/types'
import { getContactLinks } from './getContactLinks'

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

  it('собирает правильные href', () => {
    const links = getContactLinks({
      telegram: 'a_sedletskii',
      email: 'anton@example.com',
      phone: '+79999999999',
    })

    expect(links.map(({ href }) => href)).toEqual([
      'https://t.me/a_sedletskii',
      'mailto:anton@example.com',
      'tel:+79999999999',
    ])
  })

  it('отдаёт значение для показа: ник с @, почту и номер как есть', () => {
    const links = getContactLinks({
      telegram: 'a_sedletskii',
      email: 'anton@example.com',
      phone: '+79999999999',
    })

    expect(links.map(({ value }) => value)).toEqual([
      '@a_sedletskii',
      'anton@example.com',
      '+79999999999',
    ])
  })

  it('у каждой ссылки есть иконка', () => {
    const links = getContactLinks({ telegram: 'a_sedletskii', email: 'anton@example.com' })

    expect(links.every(({ icon }) => !!icon)).toBe(true)
  })

  it('пропускает незаполненные контакты', () => {
    const links = getContactLinks({ email: 'anton@example.com' })

    expect(links.map(({ key }) => key)).toEqual([ContactKey.Email])
  })

  it('пустая строка считается незаполненным контактом', () => {
    expect(getContactLinks({ telegram: '', email: '', phone: '' })).toEqual([])
  })

  it('без контактов возвращает пустой список', () => {
    expect(getContactLinks({})).toEqual([])
  })
})
