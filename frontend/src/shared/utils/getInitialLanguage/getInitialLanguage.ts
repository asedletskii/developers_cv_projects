import { appLocaleDict } from '@app/locale'
import { Languages } from '@shared/types'

export const getInitialLanguage = (fallback: Languages = Languages.Ru): Languages => {
  const stored = localStorage.getItem('language') as Languages | null
  if (stored && appLocaleDict[stored]) return stored

  const nav = navigator.language?.slice(0, 2).toLowerCase() as Languages
  if (nav && appLocaleDict[nav]) return nav

  return fallback
}
