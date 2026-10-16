import { Avatar, Box } from '@mui/material'
import { FLOAT, SPIN } from './constants'

type Props = {
  src: string
  alt: string
  initials: string
  size?: number
}

export const Portrait = ({ src, alt, initials, size = 144 }: Props) => (
  <Box
    data-preview="profile"
    sx={{
      position: 'relative',
      flexShrink: 0,
      animation: `${FLOAT} 5s ease-in-out infinite`,
      '@media (prefers-reduced-motion: reduce)': { animation: 'none' },
    }}
  >
    <Box
      aria-hidden
      sx={{
        position: 'absolute',
        inset: -Math.round(size / 10),
        borderRadius: '50%',
        border: `${Math.max(2, Math.round(size / 36))}px dashed`,
        borderColor: 'primary.light',
        animation: `${SPIN} 24s linear infinite`,
        '@media (prefers-reduced-motion: reduce)': { animation: 'none' },
      }}
    />

    <Avatar
      src={src}
      alt={alt}
      sx={{
        width: size,
        height: size,
        fontSize: `${size / 48}rem`,
        fontWeight: 700,
        color: 'primary.contrastText',
        bgcolor: 'primary.main',
      }}
    >
      {initials}
    </Avatar>
  </Box>
)
