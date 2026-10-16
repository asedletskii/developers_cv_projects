import { useGreetingAnimation } from '@features/curriculumVitae/components/cv/HeroSection/hooks/useGreetingAnimation'
import { useDownloadResume } from '@features/curriculumVitae/hooks'
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
  availability?: AvailabilityData
  updatedAt?: string
}

export const HeroSection = ({ greeting, availability, updatedAt }: Props) => {
  const rootRef = useGreetingAnimation()
  const { download, isDownloading } = useDownloadResume()

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

        <Cta onClick={download} isLoading={isDownloading} />
      </Container>

      <ScrollHint onClick={scrollToNext} />
    </Box>
  )
}
