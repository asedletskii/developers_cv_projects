import { Chip, Container, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { TECH_STACK } from './constants'

export const TechStack = () => {
  const translate = useTranslate('AboutPage.Stack')

  return (
    <Container maxWidth="md" sx={{ py: { xs: 4, md: 6 }, textAlign: 'center' }}>
      <Typography variant="h5" className="reveal" sx={{ mb: 3 }}>
        {translate('title')}
      </Typography>

      <Stack
        className="reveal-group"
        direction="row"
        flexWrap="wrap"
        justifyContent="center"
        gap={1.25}
      >
        {TECH_STACK.map((name) => (
          <Chip
            key={name}
            className="reveal-item"
            label={name}
            sx={{ fontSize: '0.9rem', px: 1, py: 2.5 }}
          />
        ))}
      </Stack>
    </Container>
  )
}
