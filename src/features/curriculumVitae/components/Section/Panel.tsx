import { PropsWithChildren, TransitionEvent, useEffect, useRef, useState } from 'react'
import { Box } from '@mui/material'
import { SectionId } from '@features/curriculumVitae/types'
import { prefersReducedMotion, ScrollTrigger } from '@shared/animations'
import { contentSx, DURATION, EASE, getContentId, getTitleId } from './constants'

type Props = PropsWithChildren<{
  id: SectionId
  isCollapsed: boolean
}>

/**
 * Сворачиваемое тело секции. Высота анимируется через grid-template-rows: 0fr ↔ 1fr,
 * без измерений в JS.
 */
export const Panel = ({ id, isCollapsed, children }: Props) => {
  // Пока идёт анимация, содержимое обрезается по высоте; в покое — нет, чтобы не резать тени карточек
  const [isAnimating, setIsAnimating] = useState(false)
  const previousCollapsedRef = useRef(isCollapsed)

  useEffect(() => {
    if (previousCollapsedRef.current === isCollapsed) return
    previousCollapsedRef.current = isCollapsed

    if (prefersReducedMotion()) {
      ScrollTrigger.refresh(true)
    } else {
      setIsAnimating(true)
    }
  }, [isCollapsed])

  const handleTransitionEnd = (event: TransitionEvent) => {
    if (event.propertyName !== 'grid-template-rows') return

    setIsAnimating(false)
    // Высота страницы изменилась — пересчитываем триггеры, но не посреди плавного скролла (safe)
    ScrollTrigger.refresh(true)
  }

  return (
    <Box
      data-collapsible
      onTransitionEnd={handleTransitionEnd}
      sx={{
        display: 'grid',
        gridTemplateRows: isCollapsed ? '0fr' : '1fr',
        transition: `grid-template-rows ${DURATION} ${EASE}`,
        '@media (prefers-reduced-motion: reduce)': { transition: 'none' },
      }}
    >
      <Box
        id={getContentId(id)}
        role="region"
        aria-labelledby={getTitleId(id)}
        sx={{
          minHeight: 0,
          overflow: isCollapsed || isAnimating ? 'hidden' : 'visible',
          // Свёрнутое содержимое недоступно ни для фокуса, ни для скринридеров
          visibility: isCollapsed ? 'hidden' : 'visible',
          opacity: isCollapsed ? 0 : 1,
          transition: isCollapsed
            ? `opacity 0.25s, visibility 0s linear ${DURATION}`
            : 'opacity 0.3s 0.05s',
          '@media (prefers-reduced-motion: reduce)': { transition: 'none' },
        }}
      >
        <Box sx={contentSx}>{children}</Box>
      </Box>
    </Box>
  )
}
