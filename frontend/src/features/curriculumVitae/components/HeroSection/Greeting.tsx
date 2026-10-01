import { Greeting as GreetingData } from '@features/curriculumVitae/types'
import { Stack, Typography } from '@mui/material'
import { splitText } from '@shared/animations'
import { ColoredTypography } from '@shared/components'
import { SHIMMER } from './constants'

type Props = {
  greeting: GreetingData
}

export const Greeting = ({ greeting }: Props) => {
  return (
    <Stack
      spacing={{ xs: 1.5, md: 2.5 }}
      sx={{ alignItems: 'center', justifyContent: 'center', textAlign: 'center' }}
    >
      {greeting.first && (
        <Typography variant="h1" sx={{ fontSize: 'clamp(3rem, 10vw, 6.5rem)', lineHeight: 1.05 }}>
          {splitText(greeting.first, 'letter')}
        </Typography>
      )}

      {(greeting.second || greeting.secondHighlight) && (
        <ColoredTypography
          className="name"
          firstPart={greeting.second}
          coloredPart={greeting.secondHighlight}
          variant="h2"
          color="primary.main"
          sx={{ fontSize: 'clamp(1.75rem, 5.5vw, 3.75rem)' }}
          coloredSx={({ palette }) => ({
            background: `linear-gradient(90deg, ${palette.primary.dark}, ${palette.primary.light}, ${palette.primary.dark})`,
            backgroundSize: '200% 100%',
            WebkitBackgroundClip: 'text',
            backgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            animation: `${SHIMMER} 6s linear infinite`,
            '@media (prefers-reduced-motion: reduce)': { animation: 'none' },
          })}
        />
      )}

      {greeting.third && (
        <Typography
          className="subtitle"
          variant="h6"
          sx={{ color: 'gray.light', fontSize: 'clamp(1.05rem, 2.4vw, 1.6rem)' }}
        >
          {greeting.third}
        </Typography>
      )}
    </Stack>
  )
}
