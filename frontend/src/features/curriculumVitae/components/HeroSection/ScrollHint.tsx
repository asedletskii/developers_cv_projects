import { Box } from '@mui/material'
import { WHEEL } from './constants'

type Props = {
  onClick(): void
}

export const ScrollHint = ({ onClick }: Props) => (
  <Box
    className="hero-hint"
    data-print-hidden
    role="button"
    tabIndex={0}
    aria-label="Scroll down"
    onClick={onClick}
    onKeyDown={(event) => event.key === 'Enter' && onClick()}
    sx={{
      cursor: 'pointer',
      position: 'absolute',
      bottom: { xs: 24, md: 40 },
      left: '50%',
      ml: '-13px',
      width: 26,
      height: 42,
      border: 2,
      borderColor: 'gray.light',
      borderRadius: 13,
      display: 'flex',
      justifyContent: 'center',
      pt: '8px',
    }}
  >
    <Box
      sx={{
        width: 4,
        height: 8,
        borderRadius: 2,
        bgcolor: 'primary.main',
        animation: `${WHEEL} 1.8s ease-in-out infinite`,
        '@media (prefers-reduced-motion: reduce)': { animation: 'none' },
      }}
    />
  </Box>
)
