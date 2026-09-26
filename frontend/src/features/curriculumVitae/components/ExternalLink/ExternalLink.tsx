import { PropsWithChildren } from 'react'
import { OpenInNewOutlined } from '@mui/icons-material'
import { Button } from '@mui/material'

type Props = PropsWithChildren<{
  href: string
}>

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
