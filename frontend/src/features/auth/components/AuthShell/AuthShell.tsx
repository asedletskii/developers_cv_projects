import { PropsWithChildren } from 'react'
import { SparkleField } from '@features/curriculumVitae/components/cv/SparkleField'
import { Box } from '@mui/material'

export const AuthShell = ({ children }: PropsWithChildren) => (
  <Box
    sx={{
      position: 'relative',
      flex: 1,
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      minHeight: '70vh',
      px: 2,
      py: { xs: 6, md: 8 },
      overflow: 'hidden',
    }}
  >
    <Box
      aria-hidden
      sx={{
        position: 'absolute',
        inset: 0,
        maskImage: 'linear-gradient(to bottom, transparent, #000 30%, #000 70%, transparent)',
        WebkitMaskImage: 'linear-gradient(to bottom, transparent, #000 30%, #000 70%, transparent)',
        background: ({ palette }) =>
          `radial-gradient(720px circle at 50% 45%, ${palette.primary.main}22, transparent 70%)`,
      }}
    />
    <SparkleField count={24} seed={5} />

    <Box sx={{ position: 'relative', zIndex: 1, width: '100%' }}>{children}</Box>
  </Box>
)
