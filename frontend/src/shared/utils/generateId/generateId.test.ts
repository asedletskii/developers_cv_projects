import { generateId } from './generateId'

describe('generateId', () => {
  it('каждый следующий id на единицу больше', () => {
    const first = generateId()

    expect(generateId()).toBe(first + 1)
    expect(generateId()).toBe(first + 2)
  })

  it('id не повторяются', () => {
    const ids = Array.from({ length: 100 }, generateId)

    expect(new Set(ids).size).toBe(100)
  })
})
