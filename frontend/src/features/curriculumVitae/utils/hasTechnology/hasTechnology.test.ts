import { hasTechnology } from './hasTechnology'

const stack = [{ name: 'React' }, { name: 'TypeScript' }, { name: 'Node.js' }]

describe('hasTechnology', () => {
  it('находит технологию в стеке', () => {
    expect(hasTechnology(stack, 'React')).toBe(true)
  })

  it('ищет без учёта регистра и пробелов по краям', () => {
    expect(hasTechnology(stack, ' typescript ')).toBe(true)
  })

  it('не находит отсутствующую технологию', () => {
    expect(hasTechnology(stack, 'Vue')).toBe(false)
  })

  it('не считает совпадением часть названия', () => {
    expect(hasTechnology(stack, 'Type')).toBe(false)
  })

  it('пустой стек ничего не содержит', () => {
    expect(hasTechnology([], 'React')).toBe(false)
  })
})
