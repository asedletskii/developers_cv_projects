import { AppBar, Box, Container } from '@mui/material'
import { Navigation } from './Navigation'

export const Header = () => {
  // Выход из аккаунта пока не нужен — заготовка оставлена на будущее
  //   const translate = useTranslate('Header')
  //   const isMobile = useMediaQuery((theme) => theme.breakpoints.down('sm'))
  //   const [isLoggingOut, setIsLoggingOut] = useState(false)
  //
  //   const handleLogout = async () => {
  //     if (isLoggingOut) return
  //     setIsLoggingOut(true)
  //     try {
  //       await logoutHelper(dispatch)
  //     } finally {
  //       setIsLoggingOut(false)
  //     }
  //   }

  return (
    <AppBar position="static" color="transparent">
      <Container maxWidth="xl" sx={{ display: 'flex', alignItems: 'center' }}>
        <Navigation />

        <Box sx={{ flexGrow: 1 }} />

        {/* {isMobile ? (
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
        )} */}
      </Container>
    </AppBar>
  )
}
