import { ContactKey } from '@features/curriculumVitae/types'
import { ContactLink } from '@features/curriculumVitae/utils'
import { ContentCopyRounded } from '@mui/icons-material'
import { Button, IconButton, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { CONTACT_BUTTON_SX } from './constants'
import { useContactButton } from './hooks'

type Props = {
  link: ContactLink
  isHidden: boolean
}

export const ContactButton = ({ link, isHidden }: Props) => {
  const translate = useTranslate('CurriculumVitae.Contacts')
  const { key, href, icon, value } = link
  const { isRevealed, reveal, copy } = useContactButton(value, isHidden)

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
          sx={CONTACT_BUTTON_SX}
        >
          {isHidden ? value : translate(key)}
        </Button>
      ) : (
        <Button onClick={reveal} startIcon={icon} size="large" sx={CONTACT_BUTTON_SX}>
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
