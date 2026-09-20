import { useTheme } from '@mui/material'

export const useThemedValue = (light: any, dark: any) => {
  const {
    palette: { mode },
  } = useTheme()

  return mode === 'light' ? light : dark
}
