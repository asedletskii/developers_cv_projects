import { Telegram } from '@mui/icons-material'
import { Link, Stack, Typography } from '@mui/material'

type Props = {
  name: string
  role: string
  telegram?: string
}

export const Developer = ({ name, role, telegram }: Props) => (
  <Stack spacing={0.5}>
    <Stack direction={'row'} sx={{ alignItems: 'end' }} spacing={1}>
      <Typography>{name}</Typography>

      <Typography variant="caption">{role}</Typography>
    </Stack>

    {telegram && (
      <Link
        href={`https://t.me/${telegram}`}
        target="_blank"
        sx={{
          display: 'flex',
          alignItems: 'center',
          textDecoration: 'none',
          width: 'max-content',
          bgcolor: 'surface.light',
          p: '4px 12px',
          borderRadius: '24px',
          m: 0,
          gap: 1,
          ':hover': {
            bgcolor: 'link.main',
            '&>svg': {
              color: '#FFF',
            },
          },
        }}
      >
        <Telegram
          sx={{
            color: 'link.main',
            ':hover': {
              color: '#FFFFFF',
            },
          }}
        />

        {`@${telegram}`}
      </Link>
    )}
  </Stack>
)
