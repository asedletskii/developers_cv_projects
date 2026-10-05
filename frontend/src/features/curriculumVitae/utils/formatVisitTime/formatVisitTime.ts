const startOfDay = (date: Date) => new Date(date.getFullYear(), date.getMonth(), date.getDate())

const capitalize = (text: string) => text.charAt(0).toUpperCase() + text.slice(1)

/** «Вчера, 16:40» для последних дней, «29 сент., 16:40» для остальных.
 *  `isCapitalized: false` — для встраивания в середину фразы («последний вчера, 16:40») */
export const formatVisitTime = (
  iso: string,
  locale: string,
  { now = new Date(), isCapitalized = true } = {},
) => {
  const date = new Date(iso)
  const time = new Intl.DateTimeFormat(locale, { hour: '2-digit', minute: '2-digit' }).format(date)
  const dayDiff = Math.round((startOfDay(date).getTime() - startOfDay(now).getTime()) / 864e5)

  if (dayDiff === 0 || dayDiff === -1) {
    const day = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' }).format(dayDiff, 'day')
    return `${isCapitalized ? capitalize(day) : day}, ${time}`
  }

  const day = new Intl.DateTimeFormat(locale, { day: 'numeric', month: 'short' }).format(date)
  return `${day}, ${time}`
}
