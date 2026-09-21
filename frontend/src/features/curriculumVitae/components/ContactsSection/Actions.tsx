import { ContactKey } from '@features/curriculumVitae/types'
import { Stack } from '@mui/material'
import { ContactButton } from './ContactButton'
import { ContactLink } from '@features/curriculumVitae/utils'

type Props = {
  links: ContactLink[]

  /** Какие контакты показываются только по клику */
  hidden: ContactKey[]
}

/** Кнопки связи: обычные ссылки или скрытые контакты, которые раскрываются по клику. */
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
