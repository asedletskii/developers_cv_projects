import { KeyboardArrowUpRounded } from '@mui/icons-material'
import { Button } from '@mui/material'
import { useScrollToTop, useTranslate } from '@shared/hooks'

/** Кнопка «Наверх»: плавно прокручивает страницу к началу. */
export const BackToTop = () => {
  const translate = useTranslate('Footer')
  const scrollToTop = useScrollToTop()

  return (
    <Button
      onClick={scrollToTop}
      endIcon={<KeyboardArrowUpRounded />}
      sx={{ textTransform: 'none', color: 'gray.light', '&:hover': { color: 'primary.main' } }}
    >
      {translate('backToTop')}
    </Button>
  )
}
