import { mergeLocaleDicts } from './mergeLocaleDicts'

describe('mergeLocaleDicts', () => {
  it('объединяет словари по каждому языку', () => {
    const merged = mergeLocaleDicts(
      { ru: { Header: { title: 'Шапка' } }, en: { Header: { title: 'Header' } } },
      { ru: { Footer: { title: 'Подвал' } }, en: { Footer: { title: 'Footer' } } },
    )

    expect(merged).toEqual({
      ru: { Header: { title: 'Шапка' }, Footer: { title: 'Подвал' } },
      en: { Header: { title: 'Header' }, Footer: { title: 'Footer' } },
    })
  })

  it('при одинаковом ключе верхнего уровня побеждает последний словарь', () => {
    const merged = mergeLocaleDicts(
      { ru: { Title: 'первый' }, en: { Title: 'first' } },
      { ru: { Title: 'второй' }, en: { Title: 'second' } },
    )

    expect(merged.ru.Title).toBe('второй')
    expect(merged.en.Title).toBe('second')
  })

  it('объединение неглубокое: вложенный объект заменяется целиком', () => {
    const merged = mergeLocaleDicts(
      { ru: { Block: { a: 'a', b: 'b' } }, en: { Block: { a: 'a', b: 'b' } } },
      { ru: { Block: { a: 'новое' } }, en: { Block: { a: 'new' } } },
    )

    expect(merged.ru.Block).toEqual({ a: 'новое' })
  })

  it('не меняет исходные словари', () => {
    const first = { ru: { A: 'а' }, en: { A: 'a' } }
    const second = { ru: { B: 'б' }, en: { B: 'b' } }

    mergeLocaleDicts(first, second)

    expect(first).toEqual({ ru: { A: 'а' }, en: { A: 'a' } })
  })

  it('со одним словарём возвращает его содержимое', () => {
    const dict = { ru: { A: 'а' }, en: { A: 'a' } }

    expect(mergeLocaleDicts(dict)).toEqual(dict)
  })
})
