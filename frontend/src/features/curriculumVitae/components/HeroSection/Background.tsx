import { SparkleField } from '@features/curriculumVitae/components/SparkleField'
import { Box } from '@mui/material'
import { ORBS } from './constants'
import { useHeroBackground } from './hooks/useHeroBackground'

/** Плавающие световые пятна: живут своей жизнью, реагируют на мышь и уезжают медленнее контента при скролле. */
export const Background = () => {
  const rootRef = useHeroBackground()

  return (
    <Box
      ref={rootRef}
      aria-hidden
      data-print-hidden
      sx={{ position: 'absolute', inset: 0, overflow: 'hidden', pointerEvents: 'none' }}
    >
      <Box className="hero-bg" sx={{ position: 'absolute', inset: 0 }}>
        {ORBS.map((orb, index) => (
          <Box
            key={index}
            className="orb-parallax"
            data-depth={orb.depth}
            sx={{ position: 'absolute', width: orb.size, height: orb.size, ...orb.position }}
          >
            <Box
              className="orb"
              sx={{
                width: '100%',
                height: '100%',
                borderRadius: '50%',
                bgcolor: orb.color,
                opacity: orb.opacity,
                filter: 'blur(80px)',
              }}
            />
          </Box>
        ))}

        <SparkleField count={44} />
      </Box>
    </Box>
  )
}
