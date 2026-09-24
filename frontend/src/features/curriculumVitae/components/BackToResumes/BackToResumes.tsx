import { glassCard } from '@features/curriculumVitae/styles'
import { ArrowBackOutlined, EditRounded } from '@mui/icons-material'
import { Button, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Link as RouterLink } from 'react-router'

type Props = {
  userId: string

  /** id резюме — есть только на странице предпросмотра, там же нужна кнопка «Редактировать» */
  cvId?: string
}

/**
 * Плавающая кнопка возврата к списку резюме ЛК (и, на предпросмотре, кнопка редактирования рядом).
 * Слева, тем же уровнем z-index и тем же стеклянным стилем, что и ControlsBar/SectionNav резюме
 * справа и по центру — не перекрываются.
 */
export const BackToResumes = ({ userId, cvId }: Props) => {
  const translate = useTranslate('Lk.Resume')

  const pillSx = {
    ...glassCard,
    borderRadius: 999,
    px: 2,
    textTransform: 'none',
  }

  return (
    <Stack
      data-print-hidden
      direction="row"
      spacing={1}
      sx={{ position: 'fixed', top: 12, left: 12, zIndex: 1201 }}
    >
      <Button
        component={RouterLink}
        to={`/cv/lk/${encodeURIComponent(userId)}`}
        startIcon={<ArrowBackOutlined />}
        sx={pillSx}
      >
        {translate('back')}
      </Button>

      {cvId && (
        <Button
          component={RouterLink}
          to={`/cv/lk/${encodeURIComponent(userId)}/${encodeURIComponent(cvId)}`}
          startIcon={<EditRounded />}
          sx={pillSx}
        >
          {translate('edit')}
        </Button>
      )}
    </Stack>
  )
}
