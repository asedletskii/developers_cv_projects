import { alpha, Theme } from '@mui/material/styles'

const isLight = ({ palette }: Theme) => palette.mode === 'light'

export const glassCard = {
  borderRadius: 3,
  border: 1,
  borderColor: 'divider',
  bgcolor: (theme: Theme) => alpha(theme.palette.surface.light, isLight(theme) ? 0.72 : 0.6),
  backdropFilter: 'blur(20px) saturate(180%)',
  WebkitBackdropFilter: 'blur(20px) saturate(180%)',
  boxShadow: (theme: Theme) =>
    `inset 0 1px 0 ${alpha(theme.palette.common.white, isLight(theme) ? 0.9 : 0.1)}, 0 8px 32px ${alpha(theme.palette.common.black, isLight(theme) ? 0.05 : 0.3)}`,

  '@media (prefers-reduced-transparency: reduce)': {
    bgcolor: 'surface.light',
    backdropFilter: 'none',
    WebkitBackdropFilter: 'none',
  },

  '@media (prefers-contrast: more)': {
    borderColor: 'gray.main',
  },
}
