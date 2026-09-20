import { ButtonBase } from '@mui/material'
import { alpha } from '@mui/material/styles'

type Props = {
  label: string
  isActive: boolean
  onClick(): void
}

export const Item = ({ label, isActive, onClick }: Props) => (
  <ButtonBase
    data-active={isActive}
    onClick={onClick}
    sx={{
      flexShrink: 0,
      px: 1.75,
      minHeight: 40,
      borderRadius: 999,
      typography: 'body1',
      fontSize: '0.875rem',
      fontWeight: isActive ? 600 : 500,
      color: isActive ? 'primary.contrastText' : 'text.primary',
      bgcolor: isActive ? 'primary.main' : 'transparent',
      transition: 'background-color 200ms, color 200ms, scale 120ms ease-out',
      '&:hover': {
        bgcolor: isActive ? 'primary.main' : ({ palette }) => alpha(palette.text.primary, 0.08),
      },
    }}
  >
    {label}
  </ButtonBase>
)
