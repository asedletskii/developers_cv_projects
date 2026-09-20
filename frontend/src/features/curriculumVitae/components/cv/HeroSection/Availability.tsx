import {
  Availability as AvailabilityData,
  AvailabilityStatus,
} from '@features/curriculumVitae/types'
import { Box, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import dayjs from 'dayjs'
import { PULSE, STATUS_COLORS } from './constants'

type Props = {
  availability?: AvailabilityData
  updatedAt?: string
}

export const Availability = ({ availability, updatedAt }: Props) => {
  const translate = useTranslate('CurriculumVitae.Availability')

  if (!availability) return null

  const { status, from, note } = availability

  const statusText = [
    translate(status),
    from && translate('from', { date: dayjs(from).format('D MMMM YYYY') }),
  ]
    .filter(Boolean)
    .join(' · ')

  const details = [note, updatedAt && translate('updated', { time: dayjs(updatedAt).fromNow() })]
    .filter(Boolean)
    .join(' · ')

  return (
    <Stack className="hero-availability" spacing={0.75} alignItems="center" sx={{ mt: 3 }}>
      <Stack
        direction="row"
        alignItems="center"
        spacing={1.25}
        sx={{
          px: 2,
          py: 0.75,
          borderRadius: 999,
          bgcolor: 'action.hover',
          border: 1,
          borderColor: 'divider',
        }}
      >
        <Box
          aria-hidden
          sx={{
            width: 10,
            height: 10,
            borderRadius: '50%',
            color: STATUS_COLORS[status],
            bgcolor: 'currentColor',
            animation:
              status === AvailabilityStatus.Open ? `${PULSE} 2s ease-out infinite` : 'none',
            '@media (prefers-reduced-motion: reduce)': { animation: 'none' },
          }}
        />

        <Typography sx={{ fontWeight: 600, fontSize: '0.9375rem' }}>{statusText}</Typography>
      </Stack>

      {details && (
        <Typography variant="caption" sx={{ color: 'gray.light' }}>
          {details}
        </Typography>
      )}
    </Stack>
  )
}
