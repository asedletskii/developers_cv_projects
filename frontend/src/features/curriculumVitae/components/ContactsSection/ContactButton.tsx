import { ContactKey } from '@features/curriculumVitae/types'
import {
  ContactLink,
  getContactDisplayValue,
  getContactHref,
} from '@features/curriculumVitae/utils'
import { ContentCopyRounded } from '@mui/icons-material'
import { Button, CircularProgress, IconButton, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { CONTACT_BUTTON_SX } from './constants'
import { useContactButton } from './hooks'

type Props = {
  resumeId: string
  link: ContactLink
  getToken: () => Promise<string | undefined>
}

export const ContactButton = ({ resumeId, link, getToken }: Props) => {
  const translate = useTranslate('CurriculumVitae.Contacts')
  const { key, icon, isHidden } = link
  const { value, isRevealing, reveal, copy, trackOpen } = useContactButton(resumeId, link, getToken)

  if (!value) {
    return (
      <Button
        onClick={reveal}
        data-analytics-id={`contact-${key}`}
        disabled={isRevealing}
        startIcon={isRevealing ? <CircularProgress size={18} color="inherit" /> : icon}
        size="large"
        sx={CONTACT_BUTTON_SX}
      >
        {translate(`Reveal.${key}`)}
      </Button>
    )
  }

  return (
    <Stack
      direction="row"
      alignItems="center"
      spacing={0.75}
      sx={{ width: { xs: '100%', sm: 'auto' } }}
    >
      <Button
        component="a"
        href={getContactHref(key, value)}
        data-analytics-id={`contact-${key}`}
        onClick={trackOpen}
        target={key === ContactKey.Telegram ? '_blank' : undefined}
        rel="noopener noreferrer"
        startIcon={icon}
        size="large"
        sx={CONTACT_BUTTON_SX}
      >
        {isHidden ? getContactDisplayValue(key, value) : translate(key)}
      </Button>

      {isHidden && (
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
