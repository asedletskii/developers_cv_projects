import { Stack } from '@mui/material'

/** Линия справа от заголовка; прорисовывается при появлении секции (см. useScrollReveal). */
export const Line = () => (
  <Stack
    className="section-line"
    sx={{ flex: 1, height: '2px', bgcolor: 'divider', borderRadius: 1 }}
  />
)
