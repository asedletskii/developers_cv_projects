import { Box, useTheme } from '@mui/material'
import { useSparkles } from './hooks'

type Props = {
  count?: number
  tone?: 'accent' | 'white'

  /** Зерно раскладки: одно и то же зерно всегда даёт одинаковое расположение звёзд */
  seed?: number
}

export const SparkleField = ({ count = 36, tone = 'accent', seed = 1 }: Props) => {
  const theme = useTheme()
  const color = tone === 'white' ? '#FFFFFF' : theme.palette.primary.light
  const { rootRef, sparkles } = useSparkles({ count, seed })

  return (
    <Box
      ref={rootRef}
      aria-hidden
      sx={{ position: 'absolute', inset: 0, overflow: 'hidden', pointerEvents: 'none' }}
    >
      {sparkles.map((sparkle, index) => (
        <Box
          key={index}
          className="sparkle"
          sx={{
            position: 'absolute',
            left: `${sparkle.x}%`,
            top: `${sparkle.y}%`,
            width: sparkle.size,
            height: sparkle.size,
            borderRadius: '50%',
            opacity: 0.6,
            bgcolor: color,
            boxShadow: `0 0 ${sparkle.size * 3}px ${color}`,
          }}
        />
      ))}
    </Box>
  )
}
