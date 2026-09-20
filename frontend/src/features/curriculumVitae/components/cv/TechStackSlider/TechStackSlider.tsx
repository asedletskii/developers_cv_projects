import { useMemo } from 'react'
import { useTechMarquee } from '@features/curriculumVitae/components/cv/TechStackSlider/hooks/useTechMarquee'
import { Technology } from '@features/curriculumVitae/types'
import { Box } from '@mui/material'
import { MIN_ITEMS_PER_HALF } from './constants'
import { Item } from './Item'

type Props = {
  technologies: Technology[]
}

export const TechStackSlider = ({ technologies }: Props) => {
  const trackRef = useTechMarquee()

  const items = useMemo(() => {
    if (!technologies.length) return []

    const repeats = Math.ceil(MIN_ITEMS_PER_HALF / technologies.length)
    const half = Array.from({ length: repeats }, () => technologies).flat()

    return [
      ...half.map(({ name }) => ({ name, isClone: false })),
      ...half.map(({ name }) => ({ name, isClone: true })),
    ]
  }, [technologies])

  if (!items.length) return null

  return (
    <Box data-print-hidden sx={{ overflow: 'hidden', py: { xs: 6, md: 10 } }}>
      <Box
        sx={{
          bgcolor: 'primary.main',
          transform: 'rotate(-2deg) scale(1.05)',
          py: { xs: 2, md: 3 },
          boxShadow: 2,
        }}
      >
        <Box
          ref={trackRef}
          sx={{
            display: 'flex',
            width: 'max-content',
            alignItems: 'center',
            willChange: 'transform',
          }}
        >
          {items.map(({ name, isClone }, index) => (
            <Item key={`${name}-${index}`} name={name} isClone={isClone} />
          ))}
        </Box>
      </Box>
    </Box>
  )
}
