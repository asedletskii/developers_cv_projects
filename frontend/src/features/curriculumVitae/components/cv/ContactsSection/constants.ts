import { keyframes } from '@mui/material/styles'

/** фон блока контактов плавно переливается */
export const DRIFT = keyframes`
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
`

export const CONTACT_BUTTON_SX = {
  flex: { xs: 1, sm: 'none' },
  textTransform: 'none',
  borderRadius: 6,
  px: 3,
  bgcolor: 'common.white',
  color: 'primary.dark',
  '&:hover': { bgcolor: 'common.white', transform: 'translateY(-2px)' },
  transition: 'transform .2s',
  '&.Mui-disabled': { bgcolor: 'common.white', color: 'primary.dark', opacity: 0.85 },
} as const
