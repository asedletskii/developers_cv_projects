import { OpenInNewOutlined } from '@mui/icons-material'
import { Link } from '@mui/material'

type Props = {
  url: string
  label: string
}

export const SiteLink = ({ url, label }: Props) => (
  <Link
    href={url}
    target="_blank"
    sx={({ typography }) => ({
      display: 'flex',
      alignItems: 'center',
      ...typography.h4,
      textDecoration: 'none',
      padding: '12px 14px',
      borderRadius: 1,
      color: 'gray.dark',
      backgroundColor: 'surface.dark',
      ':hover': {
        backgroundColor: 'surface.light',
      },
    })}
  >
    {label}

    <OpenInNewOutlined sx={{ color: 'gray.main', marginLeft: 1 }} />
  </Link>
)
