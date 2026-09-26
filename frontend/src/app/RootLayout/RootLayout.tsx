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
  const isHeaderHidden = !!matchPath({ path: ROUTES.cv, end: false }, pathname)
  // Публичная страница резюме и предпросмотр — самостоятельные страницы, футер проекта им не нужен;
  // в личном кабинете (список и редактирование) он остаётся
  const isFooterHidden = [ROUTES.cvResume, ROUTES.lkResumePreview].some((path) =>
    matchPath(path, pathname),
  )
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
      {!isHeaderHidden && <Header />}

      {/* Место под страницу резервируем целиком: короткая страница (список из одного резюме, экран
          загрузки) не подтягивает футер к себе и не «прыгает» при появлении данных */}
      <Box
        sx={{
          display: 'flex',
          flex: 1,
          minHeight: '100vh',
          // dvh учитывает адресную строку мобильного браузера; vh — запасное значение
          '@supports (min-height: 100dvh)': { minHeight: '100dvh' },
          backgroundColor: 'surface.main',
        }}
      >
        {children}
      </Box>

      {!isFooterHidden && <Footer appName={'Developers CV'} developers={developers} />}
    </Box>
  )
}
