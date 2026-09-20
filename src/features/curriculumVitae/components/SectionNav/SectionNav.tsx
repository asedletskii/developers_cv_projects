import { useMemo, useRef } from 'react'
import { Box } from '@mui/material'
import { useSectionCollapse } from '@features/curriculumVitae/components/SectionCollapseProvider'
import { useActiveSection } from '@features/curriculumVitae/hooks'
import { glassCard } from '@features/curriculumVitae/styles'
import { SectionId } from '@features/curriculumVitae/types'
import { Item } from './Item'
import { afterExpanded, scrollToSection } from '@features/curriculumVitae/utils'
import { useCenterActive } from './useCenterActive'
import { useReveal } from './useReveal'

type Props = {
  items: { id: SectionId; label: string }[]
}

/**
 * Плавающая навигация по секциям. Появляется, когда hero уходит вверх,
 * подсвечивает текущую секцию и по клику плавно доскролливает до неё.
 */
export const SectionNav = ({ items }: Props) => {
  const pillRef = useRef<HTMLDivElement | null>(null)
  const ids = useMemo(() => items.map(({ id }) => id), [items])
  const { activeId, select } = useActiveSection(ids)
  const { expand, isCollapsed } = useSectionCollapse()

  useReveal(pillRef, items[0]?.id)
  useCenterActive(pillRef, activeId)

  const navigate = (id: SectionId) => {
    select(id)

    if (!isCollapsed(id)) {
      scrollToSection(id)
      return
    }

    // Пока секция свёрнута, страница короче, и браузер обрежет скролл по текущему максимуму,
    // поэтому сначала разворачиваем, а прокручиваем, когда высота уже настоящая
    expand(id)
    afterExpanded(id, () => scrollToSection(id))
  }

  return (
    <Box
      component="nav"
      sx={{
        position: 'fixed',
        top: 12,
        left: { xs: 0, md: 200 },
        right: 200,
        zIndex: 1200,
        display: 'flex',
        justifyContent: 'center',
        pl: 1.5,
        pointerEvents: 'none',
      }}
    >
      <Box
        ref={pillRef}
        sx={{
          ...glassCard,
          display: 'flex',
          gap: 0.5,
          p: 0.5,
          maxWidth: '100%',
          borderRadius: 999,
          overflowX: 'auto',
          pointerEvents: 'auto',
          scrollbarWidth: 'none',
          '&::-webkit-scrollbar': { display: 'none' },
        }}
      >
        {items.map(({ id, label }) => (
          <Item key={id} label={label} isActive={id === activeId} onClick={() => navigate(id)} />
        ))}
      </Box>
    </Box>
  )
}
