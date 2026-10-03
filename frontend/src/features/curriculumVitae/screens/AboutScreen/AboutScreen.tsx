import { CursorGlow } from '@features/curriculumVitae/components/CursorGlow'
import { Box } from '@mui/material'
import { useDocumentTitle, useScrollReveal, useTranslate } from '@shared/hooks'
import { Cta } from './Cta'
import { Features } from './Features'
import { Hero } from './Hero'
import { Steps } from './Steps'
import { Team } from './Team'
import { TechStack } from './TechStack'

const AboutScreen = () => {
  const translate = useTranslate('AboutPage.Hero')
  const containerRef = useScrollReveal()
  useDocumentTitle(translate('title'))

  return (
    <Box ref={containerRef} sx={{ position: 'relative', width: '100%', overflowX: 'clip' }}>
      <CursorGlow />

      <Hero />
      <Steps />
      <Features />
      <TechStack />
      <Team />
      <Cta />
    </Box>
  )
}

export default AboutScreen
