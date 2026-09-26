import { OpenInNewOutlined } from '@mui/icons-material'
import { Link } from '@mui/material'

type Props = {
  url: string
  label: string
}

/** Ссылка на внешний сайт в нижней строке футера. */
export const SiteLink = ({ url, label }: Props) => (
  <Link
    href={url}
    target="_blank"
    rel="noopener noreferrer"
    underline="none"
    sx={{
      display: 'inline-flex',
      alignItems: 'center',
      gap: 0.5,
      color: 'gray.light',
      '&:hover': { color: 'primary.main' },
    }}
  >
    {label}

    <OpenInNewOutlined fontSize="small" />
  </Link>
)
