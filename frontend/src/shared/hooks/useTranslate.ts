import { useCallback, useContext } from 'react'
import { LocalizationContext } from '@app/providers'
import { _localeManager } from '@shared/locale/localeManager'
import { AppLocaleDictionary, NestedKeyOf, PathValue } from '@shared/types'

const getTranslateHook = <Dictionary extends Record<string, any>>() => {
  const useTranslate = <Prefix extends NestedKeyOf<Dictionary>>(prefix: Prefix) => {
    useLocalization()

    return useCallback(
      (key: NestedKeyOf<PathValue<Dictionary, Prefix>>, params?: Record<string, unknown>) => {
        const fullKey = `${prefix}.${key}`
        return _localeManager.translate(fullKey, params)
      },
      [prefix],
    )
  }

  return useTranslate
}

const useLocalization = () => {
  const ctx = useContext(LocalizationContext)
  if (!ctx) {
    throw new Error('useLocalization must be used within <LocalizationProvider />')
  }
  return ctx
}

export const useTranslate = getTranslateHook<AppLocaleDictionary>()
