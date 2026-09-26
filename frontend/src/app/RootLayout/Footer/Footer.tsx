import { Box, Container, Divider, Stack, Typography } from '@mui/material'
import { alpha } from '@mui/material/styles'
import { LanguageSwitch, ThemeToggle } from '@shared/components'
import { useTranslate } from '@shared/hooks'
import { getBrandGradient } from '@shared/styles/gradients'
import { BackToTop } from './BackToTop'
import { Brand } from './Brand'
import { DeveloperCard } from './DeveloperCard'
import { useCurrentYear } from './hooks'
import { SiteLink } from './SiteLink'

type Props = {
  appName: string
  developers: {
    name: string
    role: string
    telegram?: string
  }[]
  externalSiteLink?: {
    url: string
    label: string
  }
}

/** Футер: название проекта с описанием, карточки разработчиков и нижняя строка с годом и «Наверх». */
export const Footer = ({ appName, developers, externalSiteLink }: Props) => {
  const translate = useTranslate('Footer')
  const year = useCurrentYear()

  return (
    <Box
      component="footer"
      sx={{ position: 'relative', overflow: 'hidden', bgcolor: 'surface.dark' }}
    >
      {/* Градиентная линия над футером и два размытых световых пятна за содержимым */}
      <Box
        aria-hidden
        sx={(theme) => ({
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          height: 2,
          background: getBrandGradient(theme),
        })}
      />

      <Box
        aria-hidden
        sx={({ palette }) => ({
          position: 'absolute',
          inset: 0,
          pointerEvents: 'none',
          background: [
            `radial-gradient(420px circle at 8% 0%, ${alpha(palette.primary.main, 0.16)}, transparent 70%)`,
            `radial-gradient(380px circle at 95% 100%, ${alpha(palette.secondary.main, 0.14)}, transparent 70%)`,
          ].join(', '),
        })}
      />

      <Container sx={{ position: 'relative', py: { xs: 5, md: 7 } }}>
        <Stack direction={{ xs: 'column', md: 'row' }} spacing={{ xs: 4, md: 8 }}>
          <Brand appName={appName} />

          <Stack spacing={2} flex={1}>
            <Typography
              variant="caption"
              sx={{ color: 'gray.light', letterSpacing: '0.12em', textTransform: 'uppercase' }}
            >
              {translate('developers')}
            </Typography>

            <Stack direction="row" flexWrap="wrap" useFlexGap gap={2}>
              {developers.map((developer) => (
                <DeveloperCard key={developer.name} {...developer} />
              ))}
            </Stack>
          </Stack>
        </Stack>

        <Divider sx={{ my: { xs: 3, md: 4 } }} />

        <Stack
          direction={{ xs: 'column', sm: 'row' }}
          alignItems={{ xs: 'flex-start', sm: 'center' }}
          justifyContent="space-between"
          spacing={1}
        >
          <Typography variant="caption" sx={{ color: 'gray.light' }}>
            {`© ${year} ${appName} · ${translate('madeWith')}`}
          </Typography>

          <Stack direction="row" alignItems="center" spacing={2}>
            {externalSiteLink && <SiteLink {...externalSiteLink} />}

            <LanguageSwitch />

            <ThemeToggle sx={{ color: 'gray.light', '&:hover': { color: 'primary.main' } }} />

            <BackToTop />
          </Stack>
        </Stack>
      </Container>
    </Box>
  )
}
