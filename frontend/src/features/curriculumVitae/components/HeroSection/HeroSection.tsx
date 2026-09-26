import { useGreetingAnimation } from '@features/curriculumVitae/components/HeroSection/hooks/useGreetingAnimation'
import {
  Availability as AvailabilityData,
  Greeting as GreetingData,
} from '@features/curriculumVitae/types'
import { Box, Container } from '@mui/material'
import { Availability } from './Availability'
import { Background } from './Background'
import { Cta } from './Cta'
import { Greeting } from './Greeting'
import { ScrollHint } from './ScrollHint'

type Props = {
  greeting: GreetingData
  resumeLink?: string
  availability?: AvailabilityData
  updatedAt?: string
}

/** Первый экран: приветствие на весь экран, кнопка резюме и живой фон. */
export const HeroSection = ({ greeting, resumeLink, availability, updatedAt }: Props) => {
  const rootRef = useGreetingAnimation()

  const scrollToNext = () => {
    rootRef.current?.nextElementSibling?.scrollIntoView({ behavior: 'smooth' })
  }

  return (
    <Box
      ref={rootRef}
      data-hero
      sx={{
        position: 'relative',
        display: 'flex',
        alignItems: 'center',
        minHeight: '100svh',
        width: '100%',
        overflow: 'hidden',
      }}
    >
      <Background />

      <Container maxWidth="lg" sx={{ position: 'relative', zIndex: 1 }}>
        <Greeting greeting={greeting} />

        <Availability availability={availability} updatedAt={updatedAt} />

        <Cta resumeLink={resumeLink} />
      </Container>

      <ScrollHint onClick={scrollToNext} />
    </Box>
  )
}
