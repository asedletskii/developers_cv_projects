import { createContext, PropsWithChildren, useCallback, useEffect, useMemo, useState } from 'react'
import '@shared/styles/index.css'
import { CssBaseline, ThemeProvider as MuiThemeProvider } from '@mui/material'
import { darkTheme } from '@shared/styles/darkTheme'
import { lightTheme } from '@shared/styles/lightTheme'
import { ColorMode } from '@shared/types'

export const ColorModeContext = createContext<{
  colorMode: ColorMode
  changeColorMode(newColorMode: ColorMode): void
} | null>(null)

export const ColorModeProvider = ({ children }: PropsWithChildren) => {
  const [colorMode, setColorMode] = useState<ColorMode>(() => {
    const stored = window.localStorage.getItem('colorMode')
    if (stored === ColorMode.Light || stored === ColorMode.Dark) {
      return stored
    }

    if (window.matchMedia?.('(prefers-color-scheme: dark)').matches) {
      return ColorMode.Dark
    }

    return ColorMode.Light
  })

  const handleColorModeChange = useCallback((newColorMode: ColorMode) => {
    setColorMode(newColorMode)
  }, [])

  const theme = colorMode === ColorMode.Light ? lightTheme : darkTheme

  const contextValue = useMemo(
    () => ({
      colorMode,
      changeColorMode: handleColorModeChange,
    }),
    [colorMode, handleColorModeChange],
  )

  useEffect(() => {
    window.localStorage.setItem('colorMode', colorMode)
  }, [colorMode])

  useEffect(() => {
    const media = window.matchMedia('(prefers-color-scheme: dark)')

    if (media.addEventListener) {
      const handleChange = (e: MediaQueryListEvent | MediaQueryList) => {
        setColorMode(e.matches ? ColorMode.Dark : ColorMode.Light)
      }

      media.addEventListener('change', handleChange)
      return () => media.removeEventListener('change', handleChange)
    }
  }, [])

  return (
    <ColorModeContext.Provider value={contextValue}>
      <MuiThemeProvider theme={theme}>
        <CssBaseline />

        {children}
      </MuiThemeProvider>
    </ColorModeContext.Provider>
  )
}
