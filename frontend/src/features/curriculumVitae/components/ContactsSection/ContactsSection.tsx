import { SparkleField } from '@features/curriculumVitae/components/SparkleField'
import { ContactKey, Contacts, ContactsText, SectionId } from '@features/curriculumVitae/types'
import { getContactLinks } from '@features/curriculumVitae/utils'
import { Stack } from '@mui/material'
import { keyframes } from '@mui/material/styles'
import { useTranslate } from '@shared/hooks'
import { Actions } from './Actions'
import { Headline } from './Headline'

const drift = keyframes`
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
`

type Props = {
  contacts: Contacts
  text?: ContactsText
  hidden?: ContactKey[]
}

export const ContactsSection = ({ contacts, text, hidden = [] }: Props) => {
  const translate = useTranslate('CurriculumVitae.Contacts')

  const links = getContactLinks(contacts)

  if (!links.length) return null

  return (
    <Stack
      id={SectionId.Contacts}
      component="section"
      className="reveal"
      spacing={3}
      sx={{
        my: { xs: 4, md: 6 },
        scrollMarginTop: 72,
        p: { xs: 3, md: 6 },
        borderRadius: 4,
        alignItems: 'center',
        textAlign: 'center',
        position: 'relative',
        overflow: 'hidden',
        color: 'primary.contrastText',
        animation: `${drift} 12s ease-in-out infinite`,
        '@media (prefers-reduced-motion: reduce)': { animation: 'none' },
        background: ({ palette }) =>
          `linear-gradient(135deg, ${palette.primary.light}, ${palette.primary.main} 55%, ${palette.primary.dark})`,
        backgroundSize: '200% 200%',
      }}
    >
      <SparkleField count={26} tone="white" seed={7} />

      <Headline
        title={text?.title ?? translate('title')}
        subtitle={text?.subtitle ?? translate('subtitle')}
      />

      <Actions links={links} hidden={hidden} />
    </Stack>
  )
}
