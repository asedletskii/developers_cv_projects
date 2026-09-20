import { useState } from 'react'
import { glassCard } from '@features/curriculumVitae/styles'
import { CheckCircleRounded } from '@mui/icons-material'
import { Box, Container, Stack, Tab, Tabs, Typography, useMediaQuery } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { SHOWCASE, SLIDE_UP, STILL } from './constants'
import { VISUALS } from './ShowcaseVisuals'

type Key = (typeof SHOWCASE)[number]['key']

const Panel = ({ item }: { item: Key }) => {
  const translate = useTranslate('AboutPage.Showcase')
  const Visual = VISUALS[item]

  return (
    <Box
      sx={{
        display: 'grid',
        gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' },
        gap: { xs: 3, md: 5 },
        alignItems: 'center',
      }}
    >
      <Stack spacing={2}>
        <Typography variant="h4" sx={{ fontSize: 'clamp(1.4rem, 3vw, 2rem)' }}>
          {translate(`${item}.title`)}
        </Typography>

        <Typography sx={{ color: 'text.secondary' }}>{translate(`${item}.text`)}</Typography>

        <Stack component="ul" spacing={1} sx={{ m: 0, p: 0, listStyle: 'none' }}>
          {(['p1', 'p2', 'p3'] as const).map((point) => (
            <Stack key={point} component="li" direction="row" spacing={1.25} alignItems="center">
              <CheckCircleRounded fontSize="small" color="primary" />

              <Typography>{translate(`${item}.${point}`)}</Typography>
            </Stack>
          ))}
        </Stack>
      </Stack>

      <Visual />
    </Box>
  )
}

export const Showcase = () => {
  const translate = useTranslate('AboutPage.Showcase')
  const isMobile = useMediaQuery((theme) => theme.breakpoints.down('md'))
  const [active, setActive] = useState<Key>(SHOWCASE[0].key)

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

      {isMobile ? (
        <Stack className="reveal-group" spacing={2.5}>
          {SHOWCASE.map(({ key }) => (
            <Box key={key} className="reveal-item" sx={{ p: 2.5, ...glassCard }}>
              <Panel item={key} />
            </Box>
          ))}
        </Stack>
      ) : (
        <Stack className="reveal" spacing={3} sx={{ p: 4, ...glassCard }}>
          <Tabs value={active} onChange={(_, next) => setActive(next)}>
            {SHOWCASE.map(({ key, Icon }) => (
              <Tab
                key={key}
                value={key}
                icon={<Icon />}
                iconPosition="start"
                label={translate(`${key}.tab`)}
                sx={{ textTransform: 'none', fontSize: '1rem' }}
              />
            ))}
          </Tabs>

          <Box
            key={active}
            role="tabpanel"
            sx={{ animation: `${SLIDE_UP} 0.5s ease-out both`, ...STILL }}
          >
            <Panel item={active} />
          </Box>
        </Stack>
      )}
    </Container>
  )
}
