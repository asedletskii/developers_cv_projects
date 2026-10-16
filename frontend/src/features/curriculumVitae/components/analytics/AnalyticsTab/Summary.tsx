import { glassCard } from '@features/curriculumVitae/styles'
import { AnalyticsPeriod, LinkStatistics } from '@features/curriculumVitae/types'
import {
  formatDuration,
  formatVisitTime,
  getDelta,
  toPercent,
} from '@features/curriculumVitae/utils'
import { Box, Stack, Tooltip, Typography } from '@mui/material'
import { useLanguage, useTranslate } from '@shared/hooks'

type Props = {
  statistics: LinkStatistics
  previous?: LinkStatistics
  period: AnalyticsPeriod
}

type DeltaProps = { delta: number | null; hint: string }

const Delta = ({ delta, hint }: DeltaProps) => {
  if (delta === null) return null

  const color = delta > 0 ? 'success.main' : delta < 0 ? 'error.main' : 'text.secondary'
  const arrow = delta > 0 ? '▲' : delta < 0 ? '▼' : '•'

  return (
    <Tooltip title={hint} placement="bottom">
      <Typography variant="caption" sx={{ color, fontWeight: 600, cursor: 'default' }}>
        {`${arrow} ${Math.abs(delta)}%`}
      </Typography>
    </Tooltip>
  )
}

export const Summary = ({ statistics, previous, period }: Props) => {
  const translate = useTranslate('Lk.Analytics.Summary')
  const { language } = useLanguage()

  const tiles = [
    { key: 'views', value: statistics.views, before: previous?.views },
    {
      key: 'reads',
      value: statistics.reads,
      before: previous?.reads,
      hint: translate('readsHint'),
    },
    { key: 'visitors', value: statistics.uniqueVisitors, before: previous?.uniqueVisitors },
    { key: 'contacts', value: statistics.contactViews, before: previous?.contactViews },
    { key: 'downloads', value: statistics.downloads, before: previous?.downloads },
  ] as const

  const { firstViewAt, lastViewAt, filtered } = statistics
  const readRate = toPercent(statistics.reads, statistics.views)
  const contactRate = toPercent(statistics.contactViews, statistics.reads)

  const insights = [
    {
      key: 'readRate',
      value: readRate === null ? '—' : `${readRate}%`,
      hint: translate('readRateHint'),
    },
    {
      key: 'medianTime',
      value: statistics.reads ? formatDuration(statistics.medianActiveSeconds, language) : '—',
      hint: translate('medianTimeHint'),
    },
    {
      key: 'returning',
      value: translate('returningValue', {
        count: statistics.returningVisitors,
        total: statistics.uniqueVisitors,
      }),
      hint: translate('returningHint'),
    },
    {
      key: 'contactRate',
      value: contactRate === null ? '—' : translate('contactRateValue', { percent: contactRate }),
      hint: translate('contactRateHint'),
    },
  ] as const

  return (
    <Stack spacing={1.5}>
      <Box
        sx={{
          display: 'grid',
          gap: 1.5,
          gridTemplateColumns: { xs: 'repeat(2, 1fr)', md: 'repeat(5, 1fr)' },
        }}
      >
        {tiles.map(({ key, value, before, ...rest }, index) => (
          <Tooltip key={key} title={'hint' in rest ? rest.hint : ''} placement="top">
            <Stack
              spacing={0.25}
              sx={{
                p: 2,
                ...(index === tiles.length - 1 && { gridColumn: { xs: '1 / -1', md: 'auto' } }),
                ...glassCard,
              }}
            >
              <Stack direction="row" alignItems="baseline" justifyContent="space-between">
                <Typography variant="h4" sx={{ fontWeight: 700, lineHeight: 1.1 }}>
                  {value}
                </Typography>

                {before !== undefined && (
                  <Delta
                    delta={getDelta(value, before)}
                    hint={translate('vsPrevious', { days: period, before })}
                  />
                )}
              </Stack>

              <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                {translate(key)}
              </Typography>
            </Stack>
          </Tooltip>
        ))}
      </Box>

      <Box
        sx={{
          display: 'grid',
          gap: 1.5,
          p: 2,
          gridTemplateColumns: { xs: 'repeat(2, 1fr)', md: 'repeat(4, 1fr)' },
          ...glassCard,
        }}
      >
        {insights.map(({ key, value, hint }) => (
          <Tooltip key={key} title={hint} placement="top">
            <Stack>
              <Typography variant="caption" sx={{ color: 'text.secondary' }}>
                {translate(key)}
              </Typography>
              <Typography sx={{ fontWeight: 600 }}>{value}</Typography>
            </Stack>
          </Tooltip>
        ))}
      </Box>

      <Stack spacing={0.25}>
        {firstViewAt && lastViewAt && (
          <Typography variant="caption" sx={{ color: 'text.secondary' }}>
            {translate('firstLast', {
              first: formatVisitTime(firstViewAt, language, { isCapitalized: false }),
              last: formatVisitTime(lastViewAt, language, { isCapitalized: false }),
            })}
          </Typography>
        )}

        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
          {translate('filtered', { bot: filtered.bot, owner: filtered.owner })}
        </Typography>
      </Stack>
    </Stack>
  )
}
