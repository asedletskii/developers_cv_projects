import { MouseEvent } from 'react'
import { ContentCopyRounded, MoreVertRounded, OpenInNewRounded } from '@mui/icons-material'
import { Button, Chip, IconButton, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  name: string
  url: string
  isActive: boolean
  isCopied: boolean
  stats?: { views: number; reads: number }
  onCopy: () => void
  onOpenMenu: (e: MouseEvent<HTMLElement>) => void
}

export const LinkRow = ({ name, url, isActive, isCopied, stats, onCopy, onOpenMenu }: Props) => {
  const translate = useTranslate('Lk.Links')

  return (
    <Stack
      direction="row"
      alignItems="center"
      spacing={1}
      sx={{
        p: 1.5,
        border: 1,
        borderColor: 'divider',
        borderRadius: 2,
        opacity: isActive ? 1 : 0.6,
      }}
    >
      <Stack direction="row" alignItems="center" spacing={1} sx={{ flex: 1, minWidth: 0 }}>
        <Stack sx={{ minWidth: 0 }}>
          <Typography sx={{ fontWeight: 600 }}>{name}</Typography>

          <Typography variant="body2" noWrap sx={{ color: 'text.secondary' }}>
            {url}
          </Typography>

          {stats && (
            <Typography variant="caption" sx={{ color: 'text.secondary' }}>
              {translate('stats', stats)}
            </Typography>
          )}
        </Stack>

        {!isActive && <Chip size="small" label={translate('inactive')} />}
      </Stack>

      <Button
        size="small"
        startIcon={<ContentCopyRounded />}
        onClick={onCopy}
        sx={{ textTransform: 'none', whiteSpace: 'nowrap' }}
      >
        {translate(isCopied ? 'copied' : 'copy')}
      </Button>

      {isActive && (
        <IconButton
          size="small"
          component="a"
          href={url}
          target="_blank"
          rel="noopener noreferrer"
          aria-label={translate('open')}
        >
          <OpenInNewRounded fontSize="small" />
        </IconButton>
      )}

      <IconButton size="small" aria-label={translate('actions')} onClick={onOpenMenu}>
        <MoreVertRounded fontSize="small" />
      </IconButton>
    </Stack>
  )
}
