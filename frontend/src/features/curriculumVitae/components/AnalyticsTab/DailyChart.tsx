import { glassCard } from '@features/curriculumVitae/styles'
import { DailyPoint } from '@features/curriculumVitae/types'
import { alpha, Box, Stack, Typography } from '@mui/material'
import { useLanguage, useTranslate } from '@shared/hooks'

const MIN_VIEWS_FOR_CHART = 3
const WIDTH = 600
const HEIGHT = 140

type Props = {
  days: DailyPoint[]
}

export const DailyChart = ({ days }: Props) => {
  const translate = useTranslate('Lk.Analytics.Chart')
  const { language } = useLanguage()

  const total = days.reduce((sum, day) => sum + day.views, 0)
  const formatDay = (date: string) =>
    new Intl.DateTimeFormat(language, { day: 'numeric', month: 'short' }).format(new Date(date))

  // Пустая сетка из нулей ничего не говорит — показываем график, только когда есть что показать
  if (total < MIN_VIEWS_FOR_CHART) {
    return (
      <Typography variant="body2" sx={{ color: 'text.secondary' }}>
        {translate('tooShort')}
      </Typography>
    )
  }

  const max = Math.max(...days.map((day) => day.views), 1)
  const slot = WIDTH / days.length
  const barWidth = Math.max(slot * 0.7, 1)
  const barHeight = (value: number) => (value / max) * HEIGHT

  return (
    <Stack spacing={1} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Stack direction="row" justifyContent="space-between" alignItems="center">
        <Typography variant="h6">{translate('title')}</Typography>

        <Stack direction="row" spacing={2}>
          {(['views', 'reads'] as const).map((key) => (
            <Stack key={key} direction="row" spacing={0.75} alignItems="center">
              <Box
                sx={{
                  width: 10,
                  height: 10,
                  borderRadius: 0.5,
                  bgcolor: ({ palette }) =>
                    key === 'views' ? alpha(palette.primary.main, 0.35) : palette.primary.main,
                }}
              />
              <Typography variant="caption">{translate(key)}</Typography>
            </Stack>
          ))}
        </Stack>
      </Stack>

      <Box
        component="svg"
        viewBox={`0 0 ${WIDTH} ${HEIGHT}`}
        role="img"
        aria-label={translate('title')}
        sx={{ width: '100%', height: 'auto', overflow: 'visible' }}
      >
        {days.map((day, index) => {
          const x = index * slot + (slot - barWidth) / 2

          return (
            <g key={day.date}>
              <title>{`${formatDay(day.date)}: ${day.views} / ${day.reads}`}</title>
              <Box
                component="rect"
                x={x}
                y={HEIGHT - barHeight(day.views)}
                width={barWidth}
                height={barHeight(day.views)}
                rx={2}
                sx={{ fill: ({ palette }) => alpha(palette.primary.main, 0.35) }}
              />
              <Box
                component="rect"
                x={x}
                y={HEIGHT - barHeight(day.reads)}
                width={barWidth}
                height={barHeight(day.reads)}
                rx={2}
                sx={{ fill: ({ palette }) => palette.primary.main }}
              />
            </g>
          )
        })}
      </Box>

      <Stack direction="row" justifyContent="space-between">
        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
          {formatDay(days[0].date)}
        </Typography>
        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
          {formatDay(days[days.length - 1].date)}
        </Typography>
      </Stack>
    </Stack>
  )
}
