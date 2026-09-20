import { DRAW, PULSE_RING, STILL, TRAVEL } from '@features/auth/components/AuthShell/constants'
import { glassCard } from '@features/curriculumVitae/styles'
import { Telegram as TelegramIcon } from '@mui/icons-material'
import { Box, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

const LINES = [88, 64, 76]

export const LoginFlow = () => {
  const translate = useTranslate('Auth')

  return (
    <Stack spacing={2} sx={{ p: 3, ...glassCard }}>
      <Stack direction="row" alignItems="center" spacing={2}>
        <Stack alignItems="center" spacing={1} sx={{ width: 76, flexShrink: 0 }}>
          <Box sx={{ position: 'relative', width: 56, height: 56 }}>
            <Box
              sx={{
                position: 'absolute',
                inset: 0,
                borderRadius: '50%',
                border: 2,
                borderColor: 'primary.main',
                animation: `${PULSE_RING} 2.4s ease-out infinite`,
                ...STILL,
              }}
            />
            <Stack
              alignItems="center"
              justifyContent="center"
              sx={{
                position: 'absolute',
                inset: 0,
                borderRadius: '50%',
                bgcolor: 'primary.main',
                color: 'common.black',
              }}
            >
              <TelegramIcon />
            </Stack>
          </Box>
          <Typography variant="caption" sx={{ color: 'text.secondary', textAlign: 'center' }}>
            {translate('flowTelegram')}
          </Typography>
        </Stack>

        <Box
          sx={{
            flex: 1,
            height: 56,
            display: 'flex',
            alignItems: 'center',
            alignSelf: 'flex-start',
          }}
        >
          <Box
            aria-hidden
            sx={{
              position: 'relative',
              width: '100%',
              height: 2,
              backgroundImage: ({ palette }) =>
                `repeating-linear-gradient(90deg, ${palette.divider} 0 6px, transparent 6px 12px)`,
            }}
          >
            <Box
              sx={{
                position: 'absolute',
                top: -4,
                width: 10,
                height: 10,
                borderRadius: '50%',
                bgcolor: 'primary.main',
                boxShadow: ({ palette }) => `0 0 10px ${palette.primary.main}`,
                animation: `${TRAVEL} 2.4s ease-in-out infinite`,
                ...STILL,
              }}
            />
          </Box>
        </Box>

        <Stack alignItems="center" spacing={1} sx={{ width: 96, flexShrink: 0 }}>
          <Stack
            spacing={0.75}
            sx={{
              width: 80,
              height: 56,
              p: 1,
              borderRadius: 2,
              border: 1,
              borderColor: 'divider',
              justifyContent: 'center',
            }}
          >
            {LINES.map((width, index) => (
              <Box
                key={width}
                sx={{
                  height: 5,
                  width: `${width}%`,
                  borderRadius: 3,
                  bgcolor: 'primary.main',
                  opacity: 0.7,
                  transformOrigin: 'left center',
                  animation: `${DRAW} 0.8s ease-out ${0.3 + index * 0.25}s both`,
                  ...STILL,
                }}
              />
            ))}
          </Stack>
          <Typography variant="caption" sx={{ color: 'text.secondary', textAlign: 'center' }}>
            {translate('flowCabinet')}
          </Typography>
        </Stack>
      </Stack>
    </Stack>
  )
}
