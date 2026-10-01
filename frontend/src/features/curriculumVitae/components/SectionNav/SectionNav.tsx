import { useMemo, useRef } from 'react'
import { useActiveSection } from '@features/curriculumVitae/components/SectionNav/hooks/useActiveSection'
import { useCenterActive } from '@features/curriculumVitae/components/SectionNav/hooks/useCenterActive'
import { useReveal } from '@features/curriculumVitae/components/SectionNav/hooks/useReveal'
import { useSectionCollapse } from '@features/curriculumVitae/hooks'
import { glassCard } from '@features/curriculumVitae/styles'
import { SectionId } from '@features/curriculumVitae/types'
import { afterExpanded, scrollToSection } from '@features/curriculumVitae/utils'
import { Box } from '@mui/material'
import { Item } from './Item'

type Props = {
  items: { id: SectionId; label: string }[]
}

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
