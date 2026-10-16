import { Box } from '@mui/material'
import { alpha } from '@mui/material/styles'
import { useCursorGlow } from './hooks'

export const CursorGlow = () => {
  const glowRef = useCursorGlow()

  return (
    <Box
      ref={glowRef}
      aria-hidden
      data-print-hidden
      sx={{
        position: 'fixed',
        top: 0,
        left: 0,
        width: 620,
        height: 620,
        borderRadius: '50%',
        opacity: 0,
        pointerEvents: 'none',
        zIndex: -1,
        background: ({ palette }) =>
          `radial-gradient(circle, ${alpha(palette.primary.main, 0.2)}, transparent 65%)`,
      }}
    />
  )
}
