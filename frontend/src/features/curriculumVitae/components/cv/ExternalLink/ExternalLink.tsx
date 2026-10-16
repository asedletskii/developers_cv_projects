import { PropsWithChildren } from 'react'
import { useAnalytics } from '@features/curriculumVitae/hooks'
import { AnalyticsEventType } from '@features/curriculumVitae/types'
import { OpenInNewOutlined } from '@mui/icons-material'
import { Button } from '@mui/material'

type Props = PropsWithChildren<{
  href: string
  kind?: string
}>

const getHostname = (href: string) => {
  try {
    return new URL(href).hostname
  } catch {
    return ''
  }
}

export const ExternalLink = ({ href, kind = 'other', children }: Props) => {
  const { track } = useAnalytics()

  return (
    <Button
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      data-analytics-id={`external-${kind}`}
      onClick={() =>
        track({ type: AnalyticsEventType.ExternalClick, kind, target: getHostname(href) })
      }
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
}
