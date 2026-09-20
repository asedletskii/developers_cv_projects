import { alpha, Theme } from '@mui/material/styles'
import { glassCard } from './glassCard'

export const richCard = {
  ...glassCard,
  position: 'relative',
  overflow: 'hidden',
  transition: 'transform .25s ease, border-color .25s ease',

  '&::before': {
    content: '""',
    position: 'absolute',
    inset: '0 0 auto 0',
    height: 3,
    background: ({ palette }: Theme) =>
      `linear-gradient(90deg, ${palette.primary.light}, ${palette.primary.dark})`,
  },

  '&::after': {
    content: '""',
    position: 'absolute',
    top: -80,
    right: -80,
    width: 200,
    height: 200,
    borderRadius: '50%',
    pointerEvents: 'none',
    background: ({ palette }: Theme) =>
      `radial-gradient(circle, ${alpha(palette.primary.main, 0.16)}, transparent 70%)`,
  },

  '&:hover': { transform: 'translateY(-3px)', borderColor: 'primary.main' },

  '@media (prefers-reduced-motion: reduce)': {
    transition: 'none',
    '&:hover': { transform: 'none' },
  },
} as const
