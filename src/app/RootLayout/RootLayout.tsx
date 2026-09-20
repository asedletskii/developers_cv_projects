import { PropsWithChildren } from 'react'
import { Box } from '@mui/material'
import { useLocation } from 'react-router'
import { useTranslate } from '@shared/hooks'
import { Footer } from './Footer'
import { Header } from './Header'

export const RootLayout = ({ children }: PropsWithChildren) => {
  const translate = useTranslate('Footer.Developers')
  const { pathname } = useLocation()
  const isHeaderHidden = pathname === '/cv' || pathname.startsWith('/cv/')
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
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        backgroundColor: 'surface.main',
      }}
    >
      {!isHeaderHidden && <Header />}

      <Box
        sx={{
          display: 'flex',
          flex: 1,
          backgroundColor: 'surface.main',
          // pb: { xs: 8, sm: 12 },
        }}
      >
        {children}
      </Box>

      <Footer appName={'Developers CV'} developers={developers} />
    </Box>
  )
}
