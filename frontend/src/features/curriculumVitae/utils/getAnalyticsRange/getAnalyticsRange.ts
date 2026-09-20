import { AnalyticsPeriod } from '@features/curriculumVitae/types'

const toIsoDate = (date: Date) => date.toLocaleDateString('sv-SE')

export const getAnalyticsRange = (days: AnalyticsPeriod, now = new Date()) => {
  const from = new Date(now)
  from.setDate(from.getDate() - (days - 1))

  return {
    from: toIsoDate(from),
    to: toIsoDate(now),
    tz: Intl.DateTimeFormat().resolvedOptions().timeZone,
  }
}

export const getPreviousRange = (days: AnalyticsPeriod, now = new Date()) => {
  const { from } = getAnalyticsRange(days, now)
  const end = new Date(`${from}T12:00:00`)
  end.setDate(end.getDate() - 1)

  return getAnalyticsRange(days, end)
}
