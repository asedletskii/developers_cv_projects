import { PropsWithChildren } from 'react'
import { Box, SxProps, Theme } from '@mui/material'
import { useTiltCard } from './hooks'

type Props = PropsWithChildren<{
  sx?: SxProps<Theme>
}>

export const TiltCard = ({ children, sx }: Props) => {
  const cardRef = useTiltCard()

  return (
    <Box
      ref={cardRef}
      sx={[
        {
          position: 'relative',
          willChange: 'transform',
          '&::after': {
            content: '""',
            position: 'absolute',
            inset: 0,
            borderRadius: 'inherit',
            pointerEvents: 'none',
            opacity: 0,
            transition: 'opacity .3s',
            background:
              'radial-gradient(circle at var(--mx, 50%) var(--my, 50%), rgba(255,255,255,.16), transparent 55%)',
          },
          '&:hover::after': { opacity: 1 },
        },
        ...(Array.isArray(sx) ? sx : [sx]),
      ]}
    >
      {children}
    </Box>
  )
}
