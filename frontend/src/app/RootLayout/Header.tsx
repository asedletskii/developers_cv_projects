import { useLogout, useSession } from '@features/auth/hooks'
import { LogoutRounded } from '@mui/icons-material'
import { AppBar, Box, Button, Container, IconButton, useMediaQuery } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { useNavigate } from 'react-router'
import { Navigation } from './Navigation'

type Props = {
  isOverlay?: boolean
}

export const Header = ({ isOverlay = false }: Props) => {
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
    <AppBar
      position={isOverlay ? 'absolute' : 'static'}
      color="transparent"
      elevation={0}
      sx={{ zIndex: 2 }}
    >
      <Container
        maxWidth="xl"
        sx={{ display: 'grid', gridTemplateColumns: '1fr auto 1fr', alignItems: 'center', py: 1.5 }}
      >
        <Box />

        <Navigation />

        <Box sx={{ justifySelf: 'end' }}>
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
        </Box>
      </Container>
    </AppBar>
  )
}
