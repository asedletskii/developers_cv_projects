import { Avatar, Box } from '@mui/material'
import { FLOAT, SPIN } from './constants'

type Props = {
  src: string
  alt: string
  initials: string
}

export const Portrait = ({ src, alt, initials }: Props) => (
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
        inset: -14,
        borderRadius: '50%',
        border: '4px dashed',
        borderColor: 'primary.light',
        animation: `${SPIN} 24s linear infinite`,
        '@media (prefers-reduced-motion: reduce)': { animation: 'none' },
      }}
    />

    <Avatar
      src={src}
      alt={alt}
      sx={{
        width: 144,
        height: 144,
        fontSize: '3rem',
        fontWeight: 700,
        color: 'primary.contrastText',
        bgcolor: 'primary.main',
      }}
    >
      {initials}
    </Avatar>
  </Box>
)
