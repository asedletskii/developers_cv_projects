import { orderSectionStats } from './orderSectionStats'

const stat = (section: string) => ({ section })

describe('orderSectionStats', () => {
  it('следует порядку секций на странице, а не порядку ответа', () => {
    const result = orderSectionStats(
      [stat('about'), stat('experience'), stat('projects')],
      ['about', 'projects', 'experience'],
    )

    expect(result.map((item) => item.section)).toEqual(['about', 'projects', 'experience'])
  })

  it('секции, которых больше нет на странице, идут в конец', () => {
    const result = orderSectionStats(
      [stat('courses'), stat('about'), stat('projects')],
      ['about', 'projects'],
    )

    expect(result.map((item) => item.section)).toEqual(['about', 'projects', 'courses'])
  })

  it('не мутирует исходный массив', () => {
    const items = [stat('b'), stat('a')]
    orderSectionStats(items, ['a', 'b'])

    expect(items.map((item) => item.section)).toEqual(['b', 'a'])
  })
})
