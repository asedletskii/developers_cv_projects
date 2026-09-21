import { useContext } from 'react'
import { ColorModeContext } from '@app/providers'

export const useColorMode = () => {
  const ctx = useContext(ColorModeContext)
  if (!ctx) {
    throw new Error('useColorMode must be used within <ColorModeProvider />')
  }
  return ctx
}
