import { useMemo } from 'react'
import { glassCard } from '@features/curriculumVitae/styles'
import { Box } from '@mui/material'
import { alpha } from '@mui/material/styles'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { Link as RouterLink, useLocation } from 'react-router'

export const Navigation = () => {
  const { pathname } = useLocation()
  const translate = useTranslate('Header')

  const routes = useMemo(
    () => [
      { label: translate('Tabs.about'), to: ROUTES.home },
      { label: translate('Tabs.lk'), to: ROUTES.lk },
    ],
    [translate],
  )

  return (
    <Box
      component="nav"
      id="back-to-top-anchor"
      sx={{ display: 'inline-flex', gap: 0.5, p: 0.5, ...glassCard, borderRadius: 999 }}
    >
      {routes.map((r) => {
        const isActive = r.to === ROUTES.home ? pathname === r.to : pathname.startsWith(r.to)

        return (
          <Box
            key={r.to}
            component={RouterLink}
            to={r.to}
            sx={{
              px: 1.75,
              py: 0.75,
              borderRadius: 999,
              typography: 'body2',
              fontWeight: isActive ? 600 : 500,
              textDecoration: 'none',
              color: isActive ? 'primary.contrastText' : 'text.primary',
              bgcolor: isActive ? 'primary.main' : 'transparent',
              transition: 'background-color 200ms, color 200ms',
              '&:hover': {
                bgcolor: isActive
                  ? 'primary.main'
                  : ({ palette }) => alpha(palette.text.primary, 0.08),
              },
            }}
          >
            {r.label}
          </Box>
        )
      })}
    </Box>
  )
}
