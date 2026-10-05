import { getHeatColor } from './heatColor'

describe('getHeatColor', () => {
  it('редкие клики — зелёные, частые — красные', () => {
    expect(getHeatColor(0)).toContain('hsla(120,')
    expect(getHeatColor(1)).toContain('hsla(0,')
  })

  it('середина шкалы — жёлтая', () => {
    expect(getHeatColor(0.5)).toContain('hsla(60,')
  })

  it('значения вне 0..1 обрезаются', () => {
    expect(getHeatColor(5)).toBe(getHeatColor(1))
    expect(getHeatColor(-1)).toBe(getHeatColor(0))
  })
})
