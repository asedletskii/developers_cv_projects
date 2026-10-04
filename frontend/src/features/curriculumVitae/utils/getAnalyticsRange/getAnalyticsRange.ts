import { AnalyticsPeriod } from '@features/curriculumVitae/types'

const toIsoDate = (date: Date) => date.toLocaleDateString('sv-SE')

/** Параметры периода для ручек аналитики: `from`/`to` — даты в часовом поясе браузера (YYYY-MM-DD) */
export const getAnalyticsRange = (days: AnalyticsPeriod, now = new Date()) => {
  const from = new Date(now)
  from.setDate(from.getDate() - (days - 1))

  return {
    from: toIsoDate(from),
    to: toIsoDate(now),
    tz: Intl.DateTimeFormat().resolvedOptions().timeZone,
  }
}
