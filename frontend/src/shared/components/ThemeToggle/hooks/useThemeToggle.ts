import { useColorMode } from '@shared/hooks'
import { ColorMode } from '@shared/types'

/** Смена светлой/тёмной темы: на полсекунды включает плавный переход цветов по всей странице. */
export const useThemeToggle = () => {
  const { colorMode, changeColorMode } = useColorMode()

  const isDark = colorMode === ColorMode.Dark

  const toggle = () => {
    const root = document.documentElement

    root.classList.add('theme-transition')
    window.setTimeout(() => root.classList.remove('theme-transition'), 500)

    changeColorMode(isDark ? ColorMode.Light : ColorMode.Dark)
  }

  return { isDark, toggle }
}
