import { PropsWithChildren } from 'react'
import { Stack } from '@mui/material'

type Props = PropsWithChildren<{
  /** Цвет иконки (ключ палитры) */
  color?: string
  size?: number
}>

/** Квадратная плашка с иконкой слева от заголовка карточки. */
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
