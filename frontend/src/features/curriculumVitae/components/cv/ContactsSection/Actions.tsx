import { useTurnstileToken } from '@features/curriculumVitae/hooks'
import { ContactLink } from '@features/curriculumVitae/utils'
import { Stack } from '@mui/material'
import { ContactButton } from './ContactButton'

type Props = {
  resumeId: string
  links: ContactLink[]
}

export const Actions = ({ resumeId, links }: Props) => {
  const getToken = useTurnstileToken()

  return (
    <Stack
      direction={{ xs: 'column', sm: 'row' }}
      spacing={1.5}
      width={{ xs: '100%', sm: 'auto' }}
      sx={{ position: 'relative' }}
    >
      {links.map((link) => (
        <ContactButton key={link.key} resumeId={resumeId} link={link} getToken={getToken} />
      ))}
    </Stack>
  )
}
