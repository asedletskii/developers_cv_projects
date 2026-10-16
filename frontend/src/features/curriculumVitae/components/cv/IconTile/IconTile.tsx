import { PropsWithChildren } from 'react'
import { Stack } from '@mui/material'

type Props = PropsWithChildren<{
  color?: string
  size?: number
}>

export const IconTile = ({ color = 'primary.main', size = 48, children }: Props) => (
  <Stack
    sx={{
      flexShrink: 0,
      width: size,
      height: size,
      borderRadius: 2,
      alignItems: 'center',
      justifyContent: 'center',
      color,
      bgcolor: 'surface.main',
    }}
  >
    {children}
  </Stack>
)
