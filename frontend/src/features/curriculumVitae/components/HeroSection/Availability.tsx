import { Box, Stack, Typography } from '@mui/material'
import { keyframes } from '@mui/material/styles'
import {
  Availability as AvailabilityData,
  AvailabilityStatus,
} from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'
import dayjs from 'dayjs'

type Props = {
  availability?: AvailabilityData
  updatedAt?: string
}

const pulse = keyframes`
  0% { box-shadow: 0 0 0 0 currentColor; opacity: 1; }
  70% { box-shadow: 0 0 0 8px transparent; opacity: 0.7; }
  100% { box-shadow: 0 0 0 0 transparent; opacity: 1; }
`

/** Цвет индикатора: зелёный — ищу, жёлтый — рассматриваю, серый — не ищу */
const STATUS_COLORS: Record<AvailabilityStatus, string> = {
  [AvailabilityStatus.Open]: 'success.main',
  [AvailabilityStatus.Considering]: 'alert.main',
  [AvailabilityStatus.Busy]: 'gray.light',
}

/** Плашка «Открыт к предложениям» с датой выхода, заметкой и временем последнего обновления. */
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
            animation: status === AvailabilityStatus.Open ? `${pulse} 2s ease-out infinite` : 'none',
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
