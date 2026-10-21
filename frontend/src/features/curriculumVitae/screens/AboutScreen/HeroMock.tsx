import { TiltCard } from '@features/curriculumVitae/components/cv/TiltCard'
import { glassCard } from '@features/curriculumVitae/styles'
import { EmailOutlined, PhoneOutlined, Telegram } from '@mui/icons-material'
import { Box, Chip, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { DRAW_LINE, FLOAT_Y, RING, STILL } from './constants'

const LINES = [92, 74, 84, 58]

const FLOATING = [
  { label: 'React', top: -16, right: { xs: -4, sm: -28 }, delay: 0 },
  { label: 'TypeScript', bottom: 64, left: { xs: -4, sm: -40 }, delay: 1.1 },
  { label: 'Go', bottom: -18, right: 20, delay: 2 },
]

export const HeroMock = () => {
  const translate = useTranslate('AboutPage.Hero.Mock')
  const name = translate('name')

  return (
    <Box sx={{ position: 'relative', width: '100%', maxWidth: 340, mx: 'auto' }}>
      <TiltCard sx={{ ...glassCard, p: 3 }}>
        <Stack spacing={2.5}>
          <Stack direction="row" spacing={2} alignItems="center">
            <Stack
              alignItems="center"
              justifyContent="center"
              sx={{
                width: 56,
                height: 56,
                borderRadius: '50%',
                bgcolor: 'primary.main',
                color: 'common.black',
                fontWeight: 700,
              }}
            >
              {name
                .split(' ')
                .map((part) => part[0])
                .join('')}
            </Stack>

            <Stack>
              <Typography sx={{ fontWeight: 700 }}>{name}</Typography>
              <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                {translate('role')}
              </Typography>
            </Stack>
          </Stack>

          <Stack direction="row" spacing={1} alignItems="center">
            <Box sx={{ position: 'relative', width: 10, height: 10 }}>
              <Box
                sx={{
                  position: 'absolute',
                  inset: 0,
                  borderRadius: '50%',
                  bgcolor: 'success.main',
                  animation: `${RING} 2s ease-out infinite`,
                  ...STILL,
                }}
              />
              <Box
                sx={{
                  position: 'absolute',
                  inset: 0,
                  borderRadius: '50%',
                  bgcolor: 'success.main',
                }}
              />
            </Box>
            <Typography variant="body2">{translate('available')}</Typography>
          </Stack>

          <Stack spacing={1.25}>
            {LINES.map((width, index) => (
              <Box
                key={width}
                sx={{
                  height: 8,
                  width: `${width}%`,
                  borderRadius: 4,
                  bgcolor: 'action.selected',
                  transformOrigin: 'left center',
                  animation: `${DRAW_LINE} 0.9s ease-out ${0.4 + index * 0.25}s both`,
                  ...STILL,
                }}
              />
            ))}
          </Stack>

          <Stack direction="row" spacing={1.25}>
            {[Telegram, EmailOutlined, PhoneOutlined].map((Icon, index) => (
              <Stack
                key={index}
                alignItems="center"
                justifyContent="center"
                sx={{
                  width: 40,
                  height: 40,
                  borderRadius: '50%',
                  border: 1,
                  borderColor: 'divider',
                  color: 'primary.main',
                }}
              >
                <Icon fontSize="small" />
              </Stack>
            ))}
          </Stack>
        </Stack>
      </TiltCard>

      {FLOATING.map(({ label, delay, ...position }) => (
        <Chip
          key={label}
          label={label}
          color="primary"
          sx={{
            position: 'absolute',
            ...position,
            boxShadow: 3,
            animation: `${FLOAT_Y} 5s ease-in-out ${delay}s infinite`,
            ...STILL,
          }}
        />
      ))}
    </Box>
  )
}
