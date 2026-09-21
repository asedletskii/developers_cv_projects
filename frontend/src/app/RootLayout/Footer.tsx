import { Box, Container, Divider, Stack, Typography } from '@mui/material'
import { useScreenSizeValue, useTranslate } from '@shared/hooks'
import { Brand } from './Brand'
import { Developer } from './Developer'
import { SiteLink } from './SiteLink'

type Props = {
  appName: string
  developers: {
    name: string
    role: string
    telegram?: string
  }[]
  support?: {
    info: string
    email?: string
    telephone?: string
  }
  externalSiteLink?: {
    url: string
    label: string
  }
}

export const Footer = ({ appName, developers, externalSiteLink }: Props) => {
  const translate = useTranslate('Footer')

  const paddingY = useScreenSizeValue('md', 3, 4)
  const layoutDirection = useScreenSizeValue('md', 'column', 'row')
  const layoutSpacing = useScreenSizeValue('md', 3, 9)

  return (
    <Box
      sx={{
        backgroundSize: 'auto 100%',
        backgroundColor: 'surface.dark',
      }}
    >
      <Container sx={{ paddingY }}>
        <Stack direction={layoutDirection} spacing={layoutSpacing} alignItems="flex-start">
          <Brand appName={appName} />

          <Stack spacing={1} flex={1}>
            <Typography variant="h5" color="gray.dark">
              {translate('developers')}
            </Typography>

            <Stack
              spacing={1.5}
              alignItems="flex-start"
              direction={layoutDirection}
              divider={<Divider flexItem orientation="vertical" sx={{ bgcolor: 'surface.dark' }} />}
            >
              {developers.map((developer) => (
                <Developer key={developer.name} {...developer} />
              ))}

              {/* Контакты поддержки пока не используются — заготовка оставлена на будущее:
              {support.telephone && (
                <Link
                  href={`tel:${support.telephone}`}
                  sx={({ typography }) => ({
                    ...typography.caption,
                    ':hover': { color: 'gray.dark' },
                  })}
                >
                  {support.telephone}
                </Link>
              )}

              {support.email && (
                <Link
                  href={`mailto:${support.email}`}
                  sx={({ typography }) => ({
                    ...typography.caption,
                    ':hover': { color: 'gray.dark' },
                  })}
                >
                  {support.email}
                </Link>
              )} */}
            </Stack>
          </Stack>

          {externalSiteLink && <SiteLink {...externalSiteLink} />}
        </Stack>
      </Container>
    </Box>
  )
}
