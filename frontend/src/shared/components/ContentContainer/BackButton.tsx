import { ArrowBackOutlined } from '@mui/icons-material'
import { IconButton, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useNavigate } from 'react-router'

/** Возвращает на предыдущую страницу истории. */
export const BackButton = () => {
  const navigate = useNavigate()
  const translate = useTranslate('ContentContainer')

  return (
    <IconButton
      onClick={() => navigate(-1)}
      sx={{
        display: 'flex',
        justifyContent: 'flex-start',
        width: 'max-content',
        borderRadius: '6px',
        gap: 1,
        left: -10,
        color: 'text.primary',
      }}
    >
      <ArrowBackOutlined />

      <Typography>{translate('backButton')}</Typography>
    </IconButton>
  )
}
