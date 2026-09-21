import { createRandom } from './createRandom'

const take = (random: () => number, count: number) => Array.from({ length: count }, random)

describe('createRandom', () => {
  it('отдаёт числа в диапазоне [0, 1)', () => {
    const values = take(createRandom(1), 1000)

    expect(values.every((value) => value >= 0 && value < 1)).toBe(true)
  })

  it('с одним и тем же зерном даёт одну и ту же последовательность', () => {
    expect(take(createRandom(42), 20)).toEqual(take(createRandom(42), 20))
  })

  it('с разными зёрнами даёт разные последовательности', () => {
    expect(take(createRandom(1), 20)).not.toEqual(take(createRandom(2), 20))
  })

  it('не повторяет значение подряд', () => {
    const [first, second, third] = take(createRandom(7), 3)

    expect(new Set([first, second, third]).size).toBe(3)
  })

  it('независимые генераторы не влияют друг на друга', () => {
    const a = createRandom(5)
    const b = createRandom(5)

    a()
    a()

    expect(take(b, 3)).toEqual(take(createRandom(5), 3))
  })
})
