import { DailyPoint } from '@features/curriculumVitae/types'

export const toPercent = (part: number, total: number) =>
  total > 0 ? Math.round((part / total) * 100) : null

export const getDelta = (current: number, previous: number) =>
  previous > 0 ? Math.round(((current - previous) / previous) * 100) : null

export const groupByWeekday = (days: DailyPoint[]) => {
  const totals = Array.from({ length: 7 }, () => ({ views: 0, reads: 0 }))

  for (const { date, views, reads } of days) {
    const weekday = (new Date(`${date}T12:00:00`).getDay() + 6) % 7
    totals[weekday].views += views
    totals[weekday].reads += reads
  }

  return totals
}

export const getPeakDay = (days: DailyPoint[]) =>
  days.reduce<DailyPoint | undefined>(
    (peak, day) => (day.views > (peak?.views ?? 0) ? day : peak),
    undefined,
  )
