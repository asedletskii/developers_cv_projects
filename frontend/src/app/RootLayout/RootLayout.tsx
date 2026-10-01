import { PropsWithChildren } from 'react'
import { Box } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { matchPath, useLocation } from 'react-router'
import { Footer } from './Footer'
import { Header } from './Header'

export const RootLayout = ({ children }: PropsWithChildren) => {
  const translate = useTranslate('Footer.Developers')
  const { pathname } = useLocation()
  const isPreview = !!matchPath(ROUTES.lkResumePreview, pathname)
  const isLk = !!matchPath({ path: ROUTES.lk, end: false }, pathname)
  const isPublicResume = !isLk && !!matchPath(ROUTES.cvResume, pathname)
  const isChromeHidden = isPreview || isPublicResume
  const developers = [
    {
      name: translate('sedletskii'),
      role: translate('Roles.frontend'),
      telegram: 'a_sedletskii',
    },
    {
      name: translate('kosterin'),
      role: translate('Roles.backend'),
      telegram: 'KostDmitry',
    },
  ]

  return (
    <Box
      sx={{
        flex: 1,
        display: 'flex',
        flexDirection: 'column',
        backgroundColor: 'surface.main',
      }}
    >
      {!isChromeHidden && <Header />}

      <Box
        sx={{
          display: 'flex',
          flex: 1,
          minHeight: '100vh',
          '@supports (min-height: 100dvh)': { minHeight: '100dvh' },
          backgroundColor: 'surface.main',
        }}
      >
        {children}
      </Box>

      {!isChromeHidden && <Footer appName={'Developers CV'} developers={developers} />}
    </Box>
  )
}
