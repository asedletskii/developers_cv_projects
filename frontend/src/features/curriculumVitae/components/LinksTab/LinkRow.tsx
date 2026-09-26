import { ContentCopyRounded, MoreVertRounded, OpenInNewRounded } from '@mui/icons-material'
import { Button, IconButton, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { MouseEvent } from 'react'

type Props = {
  name: string
  url: string
  isCopied: boolean
  onCopy: () => void

  /** Публичная ссылка открывается по клику; именованные бэк пока не разрешает в резюме */
  isPublic?: boolean

  /** Меню действий (переименовать, удалить) — только у именованных ссылок */
  onOpenMenu?: (e: MouseEvent<HTMLElement>) => void
}

/** Одна строка списка ссылок: название, адрес и действия. */
export const LinkRow = ({ name, url, isCopied, onCopy, isPublic, onOpenMenu }: Props) => {
  const translate = useTranslate('Lk.Links')

  return (
    <Stack
      direction="row"
      alignItems="center"
      spacing={1}
      sx={{ p: 1.5, border: 1, borderColor: 'divider', borderRadius: 2 }}
    >
      <Stack sx={{ flex: 1, minWidth: 0 }}>
        <Typography sx={{ fontWeight: 600 }}>{name}</Typography>

        <Typography variant="body2" noWrap sx={{ color: 'text.secondary' }}>
          {url}
        </Typography>
      </Stack>

      <Button
        size="small"
        startIcon={<ContentCopyRounded />}
        onClick={onCopy}
        sx={{ textTransform: 'none', whiteSpace: 'nowrap' }}
      >
        {translate(isCopied ? 'copied' : 'copy')}
      </Button>

      {isPublic && (
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

      {onOpenMenu && (
        <IconButton size="small" aria-label={translate('actions')} onClick={onOpenMenu}>
          <MoreVertRounded fontSize="small" />
        </IconButton>
      )}
    </Stack>
  )
}
