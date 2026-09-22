import { useState } from 'react'
import { ContactKey } from '@features/curriculumVitae/types'
import { ContactLink } from '@features/curriculumVitae/utils'
import { ContentCopyRounded } from '@mui/icons-material'
import { Button, IconButton, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { ToastType } from '@shared/types'
import { showToast } from '@shared/utils'

type Props = {
  link: ContactLink

  /** Скрытый контакт: ни значения, ни ссылки нет в разметке, пока пользователь не нажмёт «Показать» */
  isHidden: boolean
}

/** Кнопка связи. Скрытый контакт раскрывается по клику: становится ссылкой со значением и получает «копировать». */
export const ContactButton = ({ link, isHidden }: Props) => {
  const translate = useTranslate('CurriculumVitae.Contacts')
  const [isRevealed, setIsRevealed] = useState(!isHidden)

  const { key, href, icon, value } = link

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(value)
      showToast({ messageKey: 'copied', type: ToastType.Success })
    } catch {
      showToast({ messageKey: 'copyFailed', type: ToastType.Error })
    }
  }

  const buttonSx = {
    flex: { xs: 1, sm: 'none' },
    textTransform: 'none',
    borderRadius: 6,
    px: 3,
    bgcolor: 'common.white',
    color: 'primary.dark',
    '&:hover': { bgcolor: 'common.white', transform: 'translateY(-2px)' },
    transition: 'transform .2s',
  } as const

  return (
    <Stack
      direction="row"
      alignItems="center"
      spacing={0.75}
      sx={{ width: { xs: '100%', sm: 'auto' } }}
    >
      {isRevealed ? (
        <Button
          component="a"
          href={href}
          target={key === ContactKey.Telegram ? '_blank' : undefined}
          rel="noopener noreferrer"
          startIcon={icon}
          size="large"
          sx={buttonSx}
        >
          {isHidden ? value : translate(key)}
        </Button>
      ) : (
        <Button onClick={() => setIsRevealed(true)} startIcon={icon} size="large" sx={buttonSx}>
          {translate(`Reveal.${key}`)}
        </Button>
      )}

      {isHidden && isRevealed && (
        <IconButton
          onClick={copy}
          aria-label={translate('copy')}
          title={translate('copy')}
          sx={{ color: 'inherit', bgcolor: 'rgba(255, 255, 255, 0.2)' }}
        >
          <ContentCopyRounded />
        </IconButton>
      )}
    </Stack>
  )
}
