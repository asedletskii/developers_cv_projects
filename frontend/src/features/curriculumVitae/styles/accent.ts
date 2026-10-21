import { alpha, Theme } from '@mui/material/styles'
import { glassCard } from './glassCard'

const isLight = ({ palette }: Theme) => palette.mode === 'light'

export const iconTile = {
  bgcolor: (theme: Theme) =>
    isLight(theme) ? theme.palette.primary.main : theme.palette.action.hover,
  color: (theme: Theme) =>
    isLight(theme) ? theme.palette.primary.contrastText : theme.palette.primary.main,
}

export const badge = {
  bgcolor: (theme: Theme) => (isLight(theme) ? theme.palette.primary.main : 'transparent'),
  color: (theme: Theme) =>
    isLight(theme) ? theme.palette.primary.contrastText : theme.palette.primary.main,
  borderColor: (theme: Theme) =>
    isLight(theme) ? 'transparent' : alpha(theme.palette.primary.main, 0.5),
}

export const secondaryButton = {
  color: 'text.primary',
  borderColor: 'divider',
  '&:hover': { borderColor: 'text.secondary', bgcolor: 'action.hover' },
}

export const textAccent = {
  color: (theme: Theme) =>
    isLight(theme) ? theme.palette.text.primary : theme.palette.primary.main,
}

export const dialogPaper = {
  ...glassCard,
  bgcolor: 'surface.light',
  backdropFilter: 'none',
  WebkitBackdropFilter: 'none',
  backgroundImage: 'none',
}
