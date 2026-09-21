import dayjs from 'dayjs'

/** «27 лет» / «27 years» в текущей локали dayjs */
export function formatAge(dateOfBirth: string) {
  const age = dayjs().diff(dayjs(dateOfBirth), 'year')

  return new Intl.NumberFormat(dayjs.locale(), {
    style: 'unit',
    unit: 'year',
    unitDisplay: 'long',
  }).format(age)
}
