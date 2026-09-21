import { Box } from '@mui/material'
import { keyframes } from '@mui/material/styles'

type Props = {
  onClick(): void
}

const wheel = keyframes`
  0% { transform: translateY(0); opacity: 1; }
  70% { transform: translateY(14px); opacity: 0; }
  100% { transform: translateY(0); opacity: 0; }
`

/** Индикатор «листайте вниз» внизу hero; по клику прокручивает к следующей секции. */
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
        animation: `${wheel} 1.8s ease-in-out infinite`,
        '@media (prefers-reduced-motion: reduce)': { animation: 'none' },
      }}
    />
  </Box>
)
