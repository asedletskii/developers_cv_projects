import { AccentTheme } from '@features/curriculumVitae/components/cv/AccentTheme'
import { Availability } from '@features/curriculumVitae/components/cv/HeroSection/Availability'
import { Cta } from '@features/curriculumVitae/components/cv/HeroSection/Cta'
import { Greeting } from '@features/curriculumVitae/components/cv/HeroSection/Greeting'
import { glassCard } from '@features/curriculumVitae/styles'
import { Box, Stack } from '@mui/material'
import { useHeroPreview } from './hooks'

export const HeroPreview = () => {
  const { greeting, availability, theme } = useHeroPreview()

  return (
    <AccentTheme colors={theme}>
      <Stack
        aria-hidden
        justifyContent="center"
        sx={{ p: 3, minHeight: 200, overflow: 'hidden', ...glassCard }}
      >
        <Box sx={{ zoom: 0.55, pointerEvents: 'none' }}>
          <Greeting greeting={greeting} />

          {availability?.status && <Availability availability={availability} />}

          <Cta />
        </Box>
      </Stack>
    </AccentTheme>
  )
}
