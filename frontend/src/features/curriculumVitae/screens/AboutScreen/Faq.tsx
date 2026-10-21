import { glassCard } from '@features/curriculumVitae/styles'
import { ExpandMoreRounded } from '@mui/icons-material'
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Container,
  Stack,
  Typography,
} from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { FAQ } from './constants'

export const Faq = () => {
  const translate = useTranslate('AboutPage.Faq')

  return (
    <Container maxWidth="md" sx={{ py: { xs: 4, md: 10 } }}>
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

      <Stack className="reveal-group" spacing={1.5}>
        {FAQ.map((key) => (
          <Accordion
            key={key}
            className="reveal-item"
            disableGutters
            elevation={0}
            sx={{
              ...glassCard,
              overflow: 'hidden',
              '&::before': { display: 'none' },
              transition: 'border-color .25s ease',
              '&.Mui-expanded': { borderColor: 'primary.main' },
            }}
          >
            <AccordionSummary expandIcon={<ExpandMoreRounded />} sx={{ px: 3, py: 0.5 }}>
              <Typography sx={{ fontWeight: 600 }}>{translate(`${key}.q`)}</Typography>
            </AccordionSummary>

            <AccordionDetails sx={{ px: 3, pb: 3 }}>
              <Typography sx={{ color: 'text.secondary' }}>{translate(`${key}.a`)}</Typography>
            </AccordionDetails>
          </Accordion>
        ))}
      </Stack>
    </Container>
  )
}
