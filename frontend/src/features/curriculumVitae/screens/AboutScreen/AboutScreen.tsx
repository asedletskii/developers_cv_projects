import { CursorGlow } from '@features/curriculumVitae/components/cv/CursorGlow'
import { ScrollProgress } from '@features/curriculumVitae/components/cv/ScrollProgress'
import { Box } from '@mui/material'
import { useDocumentTitle, useScrollReveal, useTranslate } from '@shared/hooks'
import { Cta } from './Cta'
import { Faq } from './Faq'
import { Features } from './Features'
import { Hero } from './Hero'
import { Showcase } from './Showcase'
import { Stats } from './Stats'
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

      <ScrollProgress />

      <Hero />

      <Stats />

      <Steps />

      <Showcase />

      <Features />

      <TechStack />

      <Faq />

      <Team />

      <Cta />
    </Box>
  )
}

export default AboutScreen
