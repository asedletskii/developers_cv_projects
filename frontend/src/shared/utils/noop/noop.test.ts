import { noop } from './noop'

describe('noop', () => {
  it('ничего не возвращает', () => {
    expect(noop()).toBeUndefined()
  })

  it('это одна и та же функция при каждом импорте', () => {
    expect(typeof noop).toBe('function')
    expect(noop).toBe(noop)
  })
})
