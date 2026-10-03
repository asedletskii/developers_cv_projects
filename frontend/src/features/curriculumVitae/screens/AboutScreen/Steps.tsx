import { Container, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useIsDesktopStory } from './hooks/useIsDesktopStory'
import { RoadStory } from './RoadStory'
import { StoryCarousel } from './StoryCarousel'

export const Steps = () => {
  const translate = useTranslate('AboutPage.Steps')
  const isDesktopStory = useIsDesktopStory()

  return (
    <Container maxWidth="lg" sx={{ py: { xs: 4, md: 10 } }}>
      <Stack
        spacing={1}
        alignItems="center"
        textAlign="center"
        className="reveal"
        sx={{ mb: { xs: 3, md: 5 } }}
      >
        <Typography variant="h3" sx={{ fontSize: 'clamp(1.75rem, 4vw, 2.5rem)' }}>
          {translate('title')}
        </Typography>
        <Typography sx={{ color: 'text.secondary' }}>{translate('subtitle')}</Typography>
      </Stack>

      {isDesktopStory ? <RoadStory /> : <StoryCarousel />}
    </Container>
  )
}
