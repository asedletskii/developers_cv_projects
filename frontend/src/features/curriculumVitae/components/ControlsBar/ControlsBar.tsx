import { glassCard } from '@features/curriculumVitae/styles'
import { Stack } from '@mui/material'
import { ThemeToggle } from '@shared/components'
import { ResumeButton } from './ResumeButton'

type Props = {
  resumeLink?: string
}

export const ControlsBar = ({ resumeLink }: Props) => (
  <Stack
    data-print-hidden
    direction="row"
    alignItems="center"
    spacing={0.5}
    sx={{
      ...glassCard,
      position: 'fixed',
      top: 12,
      right: 12,
      zIndex: 1201,
      p: 0.5,
      borderRadius: 999,
    }}
  >
    <ResumeButton resumeLink={resumeLink} />

    <ThemeToggle />
  </Stack>
)
