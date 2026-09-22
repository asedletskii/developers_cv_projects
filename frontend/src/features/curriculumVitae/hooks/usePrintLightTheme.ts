import { useEffect, useRef } from 'react'
import { useColorMode } from '@shared/hooks'
import { ColorMode } from '@shared/types'
import { flushSync } from 'react-dom'

/**
 * Печать всегда в светлой теме (иначе белый текст на белой бумаге): переключаем перед печатью
 * и возвращаем после. Работает и для кнопки, и для Ctrl+P.
 */
export const usePrintLightTheme = () => {
  const { colorMode, changeColorMode } = useColorMode()

  const currentModeRef = useRef(colorMode)
  // eslint-disable-next-line react-hooks/refs
  currentModeRef.current = colorMode

  const modeToRestoreRef = useRef<typeof colorMode | null>(null)

  useEffect(() => {
    const onBeforePrint = () => {
      modeToRestoreRef.current = currentModeRef.current
      flushSync(() => changeColorMode(ColorMode.Light))
    }

    const onAfterPrint = () => {
      if (modeToRestoreRef.current) {
        changeColorMode(modeToRestoreRef.current)
        modeToRestoreRef.current = null
      }
    }

    window.addEventListener('beforeprint', onBeforePrint)
    window.addEventListener('afterprint', onAfterPrint)

    return () => {
      window.removeEventListener('beforeprint', onBeforePrint)
      window.removeEventListener('afterprint', onAfterPrint)
    }
  }, [changeColorMode])
}
