import { getPersonFIO } from './getPersonFIO'

describe('getPersonFIO', () => {
  it('собирает фамилию, имя и отчество', () => {
    expect(getPersonFIO({ name: 'Антон', surname: 'Седлецкий', patronymic: 'Александрович' })).toBe(
      'Седлецкий Антон Александрович',
    )
  })

  it('без отчества возвращает фамилию и имя', () => {
    expect(getPersonFIO({ name: 'Дмитрий', surname: 'Костерин' })).toBe('Костерин Дмитрий')
  })

  it('пустое отчество не добавляет лишний пробел', () => {
    expect(getPersonFIO({ name: 'Дмитрий', surname: 'Костерин', patronymic: '' })).toBe(
      'Костерин Дмитрий',
    )
  })
})
