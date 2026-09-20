import { Box, CircularProgress } from '@mui/material'
import { keyframes } from '@mui/material/styles'

const FADE_IN = keyframes`
  from { opacity: 0; }
  to { opacity: 1; }
`

export const PageLoader = () => (
  <Box
    role="status"
    aria-live="polite"
    sx={{
      flex: 1,
      display: 'grid',
      placeItems: 'center',
      minHeight: '60vh',
      animation: `${FADE_IN} 0.4s ease-out 0.2s both`,
      '@media (prefers-reduced-motion: reduce)': { animation: 'none' },
    }}
  >
    <CircularProgress size={48} thickness={4} />
  </Box>
)
