import { useLogout, useSession } from '@features/auth/hooks'
import { LogoutRounded } from '@mui/icons-material'
import { AppBar, Box, Button, Container, IconButton, useMediaQuery } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { useNavigate } from 'react-router'
import { Navigation } from './Navigation'

export const Header = () => {
  const translate = useTranslate('Header')
  const isMobile = useMediaQuery((theme) => theme.breakpoints.down('sm'))
  const { isAuthenticated } = useSession()
  const { signOut, isLoggingOut } = useLogout()
  const navigate = useNavigate()

  const handleLogout = async () => {
    if (isLoggingOut) return
    await signOut()
    navigate(ROUTES.home)
  }

  return (
    <AppBar position="static" color="transparent">
      <Container maxWidth="xl" sx={{ display: 'flex', alignItems: 'center' }}>
        <Navigation />

        <Box sx={{ flexGrow: 1 }} />

        {isAuthenticated &&
          (isMobile ? (
            <IconButton
              aria-label={translate('logout')}
              onClick={handleLogout}
              disabled={isLoggingOut}
              sx={{
                bgcolor: 'transparent',
                '&:hover': { bgcolor: 'transparent' },
                typography: 'caption',
                color: 'secondary.main',
              }}
            >
              <LogoutRounded />
            </IconButton>
          ) : (
            <Button
              onClick={handleLogout}
              disabled={isLoggingOut}
              sx={{
                bgcolor: 'transparent',
                '&:hover': { bgcolor: 'transparent' },
                typography: 'caption',
              }}
              endIcon={<LogoutRounded />}
            >
              {translate('logout')}
            </Button>
          ))}
      </Container>
    </AppBar>
  )
}
