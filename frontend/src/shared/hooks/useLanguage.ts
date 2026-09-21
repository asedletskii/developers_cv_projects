import { useContext } from 'react'
import { LocalizationContext } from '@app/providers'

export const useLanguage = () => {
  const ctx = useContext(LocalizationContext)
  if (!ctx) {
    throw new Error('useLanguage must be used within <LocalizationProvider />')
  }
  return ctx
}
