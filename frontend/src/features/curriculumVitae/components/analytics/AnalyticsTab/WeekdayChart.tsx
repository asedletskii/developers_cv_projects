import { glassCard } from '@features/curriculumVitae/styles'
import { DailyPoint } from '@features/curriculumVitae/types'
import { groupByWeekday } from '@features/curriculumVitae/utils'
import { Stack, Typography } from '@mui/material'
import { useLanguage, useTranslate } from '@shared/hooks'
import { BarList } from './BarList'

const MIN_VIEWS = 3
const MONDAY = new Date(2024, 0, 1)

type Props = {
  days: DailyPoint[]
}

export const WeekdayChart = ({ days }: Props) => {
  const translate = useTranslate('Lk.Analytics.Weekday')
  const { language } = useLanguage()

  const weekdays = groupByWeekday(days)
  const total = weekdays.reduce((sum, { views }) => sum + views, 0)

  if (total < MIN_VIEWS) return null

  const names = weekdays.map((_, index) =>
    new Intl.DateTimeFormat(language, { weekday: 'long' }).format(
      new Date(2024, 0, MONDAY.getDate() + index),
    ),
  )
  const best = weekdays.reduce(
    (top, { views }, index) => (views > weekdays[top].views ? index : top),
    0,
  )

  return (
    <Stack spacing={2} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Stack spacing={0.25}>
        <Typography variant="h6">{translate('title')}</Typography>
        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
          {translate('best', { day: names[best] })}
        </Typography>
      </Stack>

      <BarList
        rows={weekdays.map(({ views, reads }, index) => ({
          key: String(index),
          label: names[index].charAt(0).toUpperCase() + names[index].slice(1),
          value: views,
          caption: translate('reads', { count: reads }),
        }))}
      />
    </Stack>
  )
}
