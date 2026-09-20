/** «2 мин 11 с» / «2m 11s» — через Intl, без своих словарей */
export const formatDuration = (totalSeconds: number, locale: string) => {
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = Math.round(totalSeconds % 60)
  const unit = (value: number, name: 'minute' | 'second') =>
    new Intl.NumberFormat(locale, { style: 'unit', unit: name, unitDisplay: 'narrow' }).format(
      value,
    )

  if (minutes === 0) return unit(seconds, 'second')
  if (seconds === 0) return unit(minutes, 'minute')
  return `${unit(minutes, 'minute')} ${unit(seconds, 'second')}`
}
