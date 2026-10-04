import { glassCard } from '@features/curriculumVitae/styles'
import { LinkStatistics } from '@features/curriculumVitae/types'
import { formatVisitTime } from '@features/curriculumVitae/utils'
import { Box, Stack, Tooltip, Typography } from '@mui/material'
import { useLanguage, useTranslate } from '@shared/hooks'

type Props = {
  statistics: LinkStatistics
}

export const Summary = ({ statistics }: Props) => {
  const translate = useTranslate('Lk.Analytics.Summary')
  const { language } = useLanguage()

  const tiles = [
    { key: 'views', value: statistics.views },
    { key: 'reads', value: statistics.reads, hint: translate('readsHint') },
    { key: 'visitors', value: statistics.uniqueVisitors },
    { key: 'contacts', value: statistics.contactViews },
    { key: 'downloads', value: statistics.downloads },
  ] as const

  const { firstViewAt, lastViewAt, filtered } = statistics

  return (
    <Stack spacing={1.5}>
      <Box
        sx={{
          display: 'grid',
          gap: 1.5,
          gridTemplateColumns: { xs: 'repeat(2, 1fr)', md: 'repeat(5, 1fr)' },
        }}
      >
        {tiles.map(({ key, value, ...rest }) => (
          <Tooltip key={key} title={'hint' in rest ? rest.hint : ''} placement="top">
            <Stack spacing={0.25} sx={{ p: 2, ...glassCard }}>
              <Typography variant="h4" sx={{ fontWeight: 700, lineHeight: 1.1 }}>
                {value}
              </Typography>
              <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                {translate(key)}
              </Typography>
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
