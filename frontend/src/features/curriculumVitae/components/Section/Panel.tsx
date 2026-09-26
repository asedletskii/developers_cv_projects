import { PropsWithChildren } from 'react'
import { SectionId } from '@features/curriculumVitae/types'
import { Box } from '@mui/material'
import { contentSx, DURATION, EASE } from './constants'
import { getContentId, getTitleId } from '@features/curriculumVitae/utils/getSectionElementId'
import { usePanelAnimation } from './hooks'

type Props = PropsWithChildren<{
  id: SectionId
  isCollapsed: boolean
}>

/**
 * Сворачиваемое тело секции. Высота анимируется через grid-template-rows: 0fr ↔ 1fr,
 * без измерений в JS.
 */
export const Panel = ({ id, isCollapsed, children }: Props) => {
  const { isAnimating, handleTransitionEnd } = usePanelAnimation(isCollapsed)

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
