import dayjs from 'dayjs'

/** «март 2023» / «Mar 2023» в текущей локали dayjs */
export const formatMonth = (date: string) => dayjs(date).format('MMM YYYY')
