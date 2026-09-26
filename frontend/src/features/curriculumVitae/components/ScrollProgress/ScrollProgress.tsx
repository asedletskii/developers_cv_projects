import { Box } from '@mui/material'
import { useScrollProgress } from './hooks'

/** Тонкая полоска прогресса прокрутки страницы вверху экрана. */
export const ScrollProgress = () => {
  const barRef = useScrollProgress()

  return (
    <Box
      aria-hidden
      data-print-hidden
      sx={{ position: 'fixed', top: 0, left: 0, right: 0, height: 3, zIndex: 1300 }}
    >
      <Box
        ref={barRef}
        sx={{
          height: '100%',
          transformOrigin: 'left center',
          background: ({ palette }) =>
            `linear-gradient(90deg, ${palette.primary.main}, ${palette.secondary.main})`,
        }}
      />
    </Box>
  )
}
