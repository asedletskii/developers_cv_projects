import dayjs from 'dayjs'

/** «2 года 3 месяца» / «2 years 3 months» между датами (без конца — до сегодня) */
export function formatDuration(start: string, end?: string) {
  const totalMonths = Math.max(1, dayjs(end).diff(dayjs(start), 'month') + 1)
  const years = Math.floor(totalMonths / 12)
  const months = totalMonths % 12

  const format = (value: number, unit: 'year' | 'month') =>
    new Intl.NumberFormat(dayjs.locale(), { style: 'unit', unit, unitDisplay: 'long' }).format(
      value,
    )

  return [years && format(years, 'year'), months && format(months, 'month')]
    .filter(Boolean)
    .join(' ')
}
