import { glassCard } from '@features/curriculumVitae/styles'
import { DeviceType, LinkStatistics } from '@features/curriculumVitae/types'
import { toPercent } from '@features/curriculumVitae/utils'
import { Box, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { BarList } from './BarList'

type Props = {
  statistics: LinkStatistics
}

export const Breakdown = ({ statistics }: Props) => {
  const translate = useTranslate('Lk.Analytics.Breakdown')
  const translateDevice = useTranslate('Lk.Analytics.Device')

  const devices = (Object.values(DeviceType) as DeviceType[])
    .map((device) => ({
      key: device,
      label: translateDevice(device),
      value: statistics.devices[device] ?? 0,
    }))
    .filter((row) => row.value > 0)

  const withShare = (rows: { key: string; label: string; value: number }[]) => {
    const total = rows.reduce((sum, { value }) => sum + value, 0)
    return rows.map((row) => ({ ...row, caption: `${toPercent(row.value, total) ?? 0}%` }))
  }

  const referrers = statistics.referrers.map(({ host, views }) => ({
    key: host ?? 'direct',
    label: host ?? translate('direct'),
    value: views,
  }))

  return (
    <Box
      sx={{ display: 'grid', gap: 2.5, gridTemplateColumns: { xs: '1fr', md: 'repeat(2, 1fr)' } }}
    >
      {[
        { title: translate('devices'), rows: withShare(devices) },
        { title: translate('referrers'), rows: withShare(referrers) },
      ].map(({ title, rows }) => (
        <Stack key={title} spacing={2} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
          <Typography variant="h6">{title}</Typography>
          <BarList rows={rows} />
        </Stack>
      ))}
    </Box>
  )
}
