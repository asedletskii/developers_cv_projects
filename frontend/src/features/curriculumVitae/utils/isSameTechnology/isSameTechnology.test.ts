import { isSameTechnology } from './isSameTechnology'

describe('isSameTechnology', () => {
  it('совпадающие названия — одна технология', () => {
    expect(isSameTechnology('React', 'React')).toBe(true)
  })

  it('регистр не важен', () => {
    expect(isSameTechnology('TypeScript', 'typescript')).toBe(true)
  })

  it('пробелы по краям не важны', () => {
    expect(isSameTechnology('  Node.js ', 'node.js')).toBe(true)
  })

  it('разные названия — разные технологии', () => {
    expect(isSameTechnology('React', 'React Native')).toBe(false)
    expect(isSameTechnology('React', 'React.js')).toBe(false)
  })

  it('пробелы внутри названия важны', () => {
    expect(isSameTechnology('Redux Toolkit', 'ReduxToolkit')).toBe(false)
  })
})
