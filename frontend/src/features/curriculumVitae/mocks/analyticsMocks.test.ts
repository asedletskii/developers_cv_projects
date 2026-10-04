import { buildClicks, buildDaily, buildStatistics, buildVisitsPage } from './analyticsMocks'

describe('analyticsMocks', () => {
  it('daily заполняет каждый день периода', () => {
    expect(buildDaily('2026-09-28', '2026-10-02').map((day) => day.date)).toEqual([
      '2026-09-28',
      '2026-09-29',
      '2026-09-30',
      '2026-10-01',
      '2026-10-02',
    ])
  })

  it('статистика согласована с графиком', () => {
    const days = buildDaily('2026-09-01', '2026-09-30')

    expect(buildStatistics('2026-09-01', '2026-09-30').views).toBe(
      days.reduce((sum, day) => sum + day.views, 0),
    )
  })

  it('лента отдаётся страницами до конца', () => {
    const first = buildVisitsPage()
    const second = buildVisitsPage(first.nextCursor!)
    const third = buildVisitsPage(second.nextCursor!)

    expect([first.items.length, second.items.length, third.items.length]).toEqual([20, 20, 5])
    expect(third.nextCursor).toBeNull()
  })

  it('клики по прежним версиям секции не входят в сетку, пока их не попросили', () => {
    expect(buildClicks('about', 'lg').staleTotal).toBeGreaterThan(0)
    expect(buildClicks('about', 'lg', true).staleTotal).toBe(0)
  })
})
