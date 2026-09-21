import { Avatar, Box } from '@mui/material'
import { keyframes } from '@mui/material/styles'

type Props = {
  src: string
  alt: string
  initials: string
}

const float = keyframes`
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
`

const spin = keyframes`
  to { transform: rotate(360deg); }
`

/** Круглое фото: медленно парит, вокруг вращается пунктирное кольцо. Без фото — инициалы. */
export const Portrait = ({ src, alt, initials }: Props) => (
  <Box
    sx={{
      position: 'relative',
      flexShrink: 0,
      animation: `${float} 5s ease-in-out infinite`,
      '@media (prefers-reduced-motion: reduce)': { animation: 'none' },
    }}
  >
    <Box
      aria-hidden
      sx={{
        position: 'absolute',
        inset: -14,
        borderRadius: '50%',
        border: '2px dashed',
        borderColor: 'primary.light',
        animation: `${spin} 24s linear infinite`,
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
