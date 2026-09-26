import { Telegram } from '@mui/icons-material'
import { Box, Link, Stack, Typography } from '@mui/material'
import { alpha } from '@mui/material/styles'
import { getBrandGradient } from '@shared/styles/gradients'
import { getInitials } from '@shared/utils'
import { AVATAR_SIZE, CARD_MIN_WIDTH } from './constants'

type Props = {
  name: string
  role: string
  telegram?: string
}

/** Карточка разработчика: аватар с инициалами, имя, роль и ссылка на Telegram. Приподнимается при наведении. */
export const DeveloperCard = ({ name, role, telegram }: Props) => (
  <Stack
    spacing={2}
    sx={({ palette }) => ({
      flex: '1 1 ' + CARD_MIN_WIDTH + 'px',
      minWidth: 0,
      p: 2.5,
      borderRadius: 4,
      border: 1,
      borderColor: 'divider',
      bgcolor: alpha(palette.surface.light, 0.5),
      transition: 'transform .25s ease, border-color .25s ease, box-shadow .25s ease',
      '&:hover': {
        transform: 'translateY(-4px)',
        borderColor: alpha(palette.primary.main, 0.6),
        boxShadow: `0 12px 32px ${alpha(palette.primary.main, 0.18)}`,
      },
      '@media (prefers-reduced-motion: reduce)': {
        transition: 'none',
        '&:hover': { transform: 'none' },
      },
    })}
  >
    <Stack direction="row" alignItems="center" spacing={1.5}>
      <Box
        aria-hidden
        sx={(theme) => ({
          width: AVATAR_SIZE,
          height: AVATAR_SIZE,
          flexShrink: 0,
          borderRadius: '50%',
          display: 'grid',
          placeItems: 'center',
          fontWeight: 700,
          color: theme.palette.common.white,
          background: getBrandGradient(theme, 135),
        })}
      >
        {getInitials(name)}
      </Box>

      <Stack sx={{ minWidth: 0 }}>
        <Typography sx={{ fontWeight: 600 }}>{name}</Typography>

        <Typography variant="caption" sx={{ color: 'gray.light' }}>
          {role}
        </Typography>
      </Stack>
    </Stack>

    {telegram && (
      <Link
        href={`https://t.me/${telegram}`}
        target="_blank"
        rel="noopener noreferrer"
        underline="none"
        sx={{
          display: 'inline-flex',
          alignItems: 'center',
          gap: 1,
          width: 'max-content',
          px: 1.5,
          py: 0.5,
          borderRadius: 999,
          color: 'link.main',
          bgcolor: 'action.hover',
          transition: 'background-color .2s, color .2s',
          '&:hover': { bgcolor: 'link.main', color: '#FFF' },
        }}
      >
        <Telegram fontSize="small" />

        {`@${telegram}`}
      </Link>
    )}
  </Stack>
)
