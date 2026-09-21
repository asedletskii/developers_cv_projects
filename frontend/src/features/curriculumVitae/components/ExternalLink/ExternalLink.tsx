import { PropsWithChildren } from 'react'
import { OpenInNewOutlined } from '@mui/icons-material'
import { Button } from '@mui/material'

type Props = PropsWithChildren<{
  href: string
}>

/**
 * Текстовая ссылка на внешний ресурс. Внутренний отступ даёт заметную подсветку при наведении,
 * а сдвиг влево на ту же величину возвращает текст на линию выравнивания с остальным содержимым
 * карточки. Именно position, а не margin: Stack обнуляет margin у своих детей.
 */
export const ExternalLink = ({ href, children }: Props) => (
  <Button
    href={href}
    target="_blank"
    rel="noopener noreferrer"
    endIcon={<OpenInNewOutlined />}
    sx={{
      alignSelf: 'flex-start',
      position: 'relative',
      left: ({ spacing }) => spacing(-1.5),
      px: 1.5,
      py: 0.75,
      textTransform: 'none',
    }}
  >
    {children}
  </Button>
)
