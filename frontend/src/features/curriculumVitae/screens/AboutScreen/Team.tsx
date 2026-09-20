import { glassCard } from '@features/curriculumVitae/styles'
import { Telegram } from '@mui/icons-material'
import { Box, Container, Link, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { getInitials } from '@shared/utils'

const DEVELOPERS = [
  { key: 'sedletskii', roleKey: 'frontend', telegram: 'a_sedletskii' },
  { key: 'kosterin', roleKey: 'backend', telegram: 'KostDmitry' },
] as const

export const Team = () => {
  const translate = useTranslate('AboutPage.Team')
  const translateDevelopers = useTranslate('Footer.Developers')

  return (
    <Container maxWidth="md" sx={{ py: { xs: 4, md: 10 } }}>
      <Stack
        spacing={1}
        alignItems="center"
        textAlign="center"
        className="reveal"
        sx={{ mb: { xs: 3, md: 6 } }}
      >
        <Typography variant="h3" sx={{ fontSize: 'clamp(1.75rem, 4vw, 2.5rem)' }}>
          {translate('title')}
        </Typography>
        <Typography sx={{ color: 'text.secondary' }}>{translate('subtitle')}</Typography>
      </Stack>

      <Stack
        className="reveal-group"
        direction={{ xs: 'column', sm: 'row' }}
        spacing={2.5}
        justifyContent="center"
      >
        {DEVELOPERS.map(({ key, roleKey, telegram }) => {
          const name = translateDevelopers(key)

          return (
            <Stack
              key={key}
              className="reveal-item"
              spacing={2}
              alignItems="center"
              textAlign="center"
              sx={{ flex: 1, p: 4, ...glassCard }}
            >
              <Box
                sx={{
                  width: 80,
                  height: 80,
                  borderRadius: '50%',
                  display: 'grid',
                  placeItems: 'center',
                  fontWeight: 700,
                  fontSize: '1.5rem',
                  color: 'primary.contrastText',
                  bgcolor: 'primary.main',
                }}
              >
                {getInitials(name)}
              </Box>

              <Stack spacing={0.25}>
                <Typography variant="h6">{name}</Typography>
                <Typography
                  sx={({ palette }) => ({
                    color: palette.mode === 'light' ? 'text.secondary' : 'primary.main',
                  })}
                >
                  {translateDevelopers(`Roles.${roleKey}`)}
                </Typography>
              </Stack>

              <Link
                href={`https://t.me/${telegram}`}
                target="_blank"
                rel="noopener noreferrer"
                underline="none"
                sx={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  gap: 1,
                  px: 1.5,
                  py: 0.5,
                  borderRadius: 999,
                  color: 'link.main',
                  bgcolor: 'action.hover',
                  transition: 'background-color .2s, color .2s',
                  '&:hover': { bgcolor: 'link.main', color: '#FFF' },
                }}
              >
                <Telegram fontSize="small" />
                {`@${telegram}`}
              </Link>
            </Stack>
          )
        })}
      </Stack>
    </Container>
  )
}
