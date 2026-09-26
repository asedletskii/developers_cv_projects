import { ContactKey } from '@features/curriculumVitae/types'
import { ContactLink } from '@features/curriculumVitae/utils'
import { Stack } from '@mui/material'
import { ContactButton } from './ContactButton'

type Props = {
  links: ContactLink[]
  hidden: ContactKey[]
}

export const Actions = ({ links, hidden }: Props) => (
  <Stack
    direction={{ xs: 'column', sm: 'row' }}
    spacing={1.5}
    width={{ xs: '100%', sm: 'auto' }}
    sx={{ position: 'relative' }}
  >
    {links.map((link) => (
      <ContactButton key={link.key} link={link} isHidden={hidden.includes(link.key)} />
    ))}
  </Stack>
)
