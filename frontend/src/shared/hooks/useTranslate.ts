import { useCallback, useContext } from 'react'
import { LocalizationContext } from '@app/providers'
import { _localeManager } from '@shared/locale/localeManager'
import { AppLocaleDictionary, NestedKeyOf, PathValue } from '@shared/types'

export const useTranslate = getTranslateHook<AppLocaleDictionary>()

function getTranslateHook<Dictionary extends Record<string, any>>() {
  return function useTranslate<Prefix extends NestedKeyOf<Dictionary>>(prefix: Prefix) {
    useLocalization()

    return useCallback(
      (key: NestedKeyOf<PathValue<Dictionary, Prefix>>, params?: Record<string, unknown>) => {
        const fullKey = `${prefix}.${key}`
        return _localeManager.translate(fullKey, params)
      },
      [prefix],
    )
  }
}

function useLocalization() {
  const ctx = useContext(LocalizationContext)
  if (!ctx) {
    throw new Error('useLocalization must be used within <LocalizationProvider />')
  }
  return ctx
}
