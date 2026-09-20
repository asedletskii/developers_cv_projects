import { useMemo } from 'react'
import { Tab, Tabs } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Link as RouterLink, useLocation } from 'react-router'

export const Navigation = () => {
  const { pathname } = useLocation()
  const translate = useTranslate('Header')

  const routes = useMemo(
    () => [
      { label: translate('Tabs.main'), to: '/main' },
      { label: translate('Tabs.about'), to: '/about' },
    ],
    [translate],
  )

  const value = useMemo(() => {
    const idx = routes.findIndex((r) =>
      r.to === '/' ? pathname === '/' : pathname.startsWith(r.to),
    )
    return idx === -1 ? false : idx
  }, [pathname, routes])

  return (
    <Tabs
      id="back-to-top-anchor"
      value={value}
      slotProps={{
        indicator: {
          sx: { bgcolor: 'primary.main' },
        },
      }}
    >
      {routes.map((r) => (
        <Tab
          key={r.to}
          label={r.label}
          component={RouterLink}
          to={r.to}
          disableRipple
          sx={{
            minWidth: 'auto',
            px: { xs: 1, sm: 2 },
          }}
        />
      ))}
    </Tabs>
  )
}
