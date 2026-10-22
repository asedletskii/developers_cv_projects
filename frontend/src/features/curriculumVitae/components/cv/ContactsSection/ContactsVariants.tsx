import { richCard } from '@features/curriculumVitae/styles'
import { ContactKey, Contacts, ContactsText, SectionId } from '@features/curriculumVitae/types'
import { getContactLinks } from '@features/curriculumVitae/utils'
import { Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useParams } from 'react-router'
import { Actions } from './Actions'

type Props = {
  contacts: Contacts
  text?: ContactsText
  hidden?: ContactKey[]
}

const useContactsData = ({ contacts, text, hidden = [] }: Props) => {
  const translate = useTranslate('CurriculumVitae.Contacts')
  const { cvId } = useParams()

  return {
    cvId,
    links: getContactLinks(contacts, hidden),
    title: text?.title ?? translate('title'),
    subtitle: text?.subtitle ?? translate('subtitle'),
  }
}

export const ContactsRich = (props: Props) => {
  const { cvId, links, title, subtitle } = useContactsData(props)

  if (!cvId || !links.length) return null

  return (
    <Stack
      id={SectionId.Contacts}
      component="section"
      className="reveal"
      direction={{ xs: 'column', md: 'row' }}
      spacing={3}
      alignItems={{ xs: 'flex-start', md: 'center' }}
      justifyContent="space-between"
      sx={{ my: { xs: 4, md: 6 }, scrollMarginTop: 72, p: { xs: 3, md: 5 }, ...richCard }}
    >
      <Stack spacing={1} sx={{ position: 'relative', zIndex: 1 }}>
        {title && (
          <Typography variant="h3" sx={{ fontSize: { xs: '1.75rem', md: '2.25rem' } }}>
            {title}
          </Typography>
        )}
        {subtitle && <Typography sx={{ color: 'text.secondary' }}>{subtitle}</Typography>}
      </Stack>

      <Actions resumeId={cvId} links={links} />
    </Stack>
  )
}

export const ContactsCompact = (props: Props) => {
  const { cvId, links, title } = useContactsData(props)

  if (!cvId || !links.length) return null

  return (
    <Stack
      id={SectionId.Contacts}
      component="section"
      className="reveal"
      direction={{ xs: 'column', sm: 'row' }}
      spacing={2}
      alignItems={{ xs: 'flex-start', sm: 'center' }}
      justifyContent="space-between"
      sx={{ my: 3, scrollMarginTop: 72, pt: 2, borderTop: 1, borderColor: 'divider' }}
    >
      {title && <Typography variant="h6">{title}</Typography>}

      <Actions resumeId={cvId} links={links} />
    </Stack>
  )
}
