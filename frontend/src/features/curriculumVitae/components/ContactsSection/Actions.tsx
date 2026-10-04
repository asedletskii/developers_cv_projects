import { useTurnstile } from '@features/curriculumVitae/hooks'
import { ContactLink } from '@features/curriculumVitae/utils'
import { Box, Stack } from '@mui/material'
import { ContactButton } from './ContactButton'

type Props = {
  resumeId: string
  links: ContactLink[]
}

export const Actions = ({ resumeId, links }: Props) => {
  const { containerRef, getToken } = useTurnstile()

  return (
    <>
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

      <Box
        ref={containerRef}
        sx={{ position: 'relative', display: 'flex', justifyContent: 'center' }}
      />
    </>
  )
}
