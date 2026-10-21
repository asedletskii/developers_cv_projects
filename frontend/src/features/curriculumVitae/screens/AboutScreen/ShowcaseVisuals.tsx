import {
  CheckCircleRounded,
  DescriptionRounded,
  LockRounded,
  ShieldRounded,
} from '@mui/icons-material'
import { Box, Chip, Stack } from '@mui/material'
import { DROP_FILE, GROW_BAR, HEAT_PULSE, RING, STACK_IN, STILL } from './constants'

const stage = {
  position: 'relative',
  height: 220,
  width: '100%',
  borderRadius: 3,
  overflow: 'hidden',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  bgcolor: 'action.hover',
} as const

const BAR_HEIGHTS = [40, 65, 50, 85, 60, 95, 70]
const HEAT = [0.25, 0.9, 0.5, 0.15, 0.7, 1, 0.35, 0.2, 0.8, 0.45, 0.1, 0.6]
const HEAT_COLORS = ['success.main', 'warning.main', 'error.main']

export const BuilderVisual = () => (
  <Stack spacing={1.25} sx={{ ...stage, flexDirection: 'column', px: 4 }}>
    {[70, 90, 55].map((width, index) => (
      <Stack
        key={width}
        direction="row"
        spacing={1.5}
        alignItems="center"
        sx={{
          width: '100%',
          p: 1.5,
          borderRadius: 2,
          bgcolor: 'background.paper',
          boxShadow: 1,
          animation: `${STACK_IN} 4.5s ease-in-out ${index * 0.5}s infinite`,
          ...STILL,
        }}
      >
        <Box sx={{ width: 10, height: 10, borderRadius: '50%', bgcolor: 'primary.main' }} />

        <Box sx={{ height: 8, width: `${width}%`, borderRadius: 4, bgcolor: 'action.selected' }} />
      </Stack>
    ))}
  </Stack>
)

export const AnalyticsVisual = () => (
  <Stack direction="row" spacing={3} sx={{ ...stage, px: 4 }}>
    <Stack direction="row" spacing={1} alignItems="flex-end" sx={{ height: 130, flex: 1 }}>
      {BAR_HEIGHTS.map((height, index) => (
        <Box
          key={index}
          sx={{
            flex: 1,
            height: `${height}%`,
            borderRadius: '6px 6px 0 0',
            bgcolor: 'primary.main',
            transformOrigin: 'bottom',
            animation: `${GROW_BAR} 0.9s ease-out ${index * 0.08}s both`,
            ...STILL,
          }}
        />
      ))}
    </Stack>

    <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(4, 22px)', gap: 0.75 }}>
      {HEAT.map((level, index) => (
        <Box
          key={index}
          sx={{
            width: 22,
            height: 22,
            borderRadius: 1,
            bgcolor: HEAT_COLORS[level > 0.66 ? 2 : level > 0.33 ? 1 : 0],
            opacity: 0.3 + level * 0.7,
            animation: `${HEAT_PULSE} 2.4s ease-in-out ${(index % 5) * 0.3}s infinite`,
            ...STILL,
          }}
        />
      ))}
    </Box>
  </Stack>
)

export const ProtectionVisual = () => (
  <Stack spacing={2.5} alignItems="center" sx={{ ...stage, flexDirection: 'column' }}>
    <Box
      sx={{ position: 'relative', display: 'grid', placeItems: 'center', width: 80, height: 80 }}
    >
      {[0, 1.2].map((delay) => (
        <Box
          key={delay}
          sx={{
            position: 'absolute',
            inset: 0,
            borderRadius: '50%',
            border: 2,
            borderColor: 'primary.main',
            animation: `${RING} 2.4s ease-out ${delay}s infinite`,
            ...STILL,
          }}
        />
      ))}
      <ShieldRounded sx={{ fontSize: 56, color: 'primary.main' }} />
    </Box>

    <Stack direction="row" spacing={1.5} alignItems="center">
      <Chip icon={<LockRounded />} label="••••••••••" variant="outlined" />

      <Chip
        icon={<CheckCircleRounded />}
        label="Cloudflare"
        color="success"
        sx={{ animation: `${HEAT_PULSE} 2.4s ease-in-out infinite`, ...STILL }}
      />
    </Stack>
  </Stack>
)

export const DownloadVisual = () => (
  <Stack alignItems="center" sx={{ ...stage, flexDirection: 'column', justifyContent: 'center' }}>
    <DescriptionRounded
      sx={{
        fontSize: 72,
        color: 'primary.main',
        animation: `${DROP_FILE} 2.2s ease-in infinite`,
        ...STILL,
      }}
    />
    <Box
      sx={{
        width: 120,
        height: 18,
        borderRadius: '0 0 12px 12px',
        border: 2,
        borderTop: 0,
        borderColor: 'text.secondary',
      }}
    />

    <Chip label=".html" size="small" sx={{ mt: 1.5 }} />
  </Stack>
)

export const VISUALS = {
  builder: BuilderVisual,
  analytics: AnalyticsVisual,
  protection: ProtectionVisual,
  download: DownloadVisual,
} as const
