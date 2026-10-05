import { useAnalytics, useAnalyticsNotice } from '@features/curriculumVitae/hooks'
import { glassCard } from '@features/curriculumVitae/styles'
import { InsightsRounded } from '@mui/icons-material'
import { Button, Slide, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

/** Плашка справа снизу: страницу смотрят со статистикой. Оформлена как диалог приложения, но без
 *  затемнения и без закрытия по клику вне — пропадает только по кнопке, дальше не показывается. */
export const AnalyticsNotice = () => {
  const translate = useTranslate('CurriculumVitae.AnalyticsNotice')
  const { viewId } = useAnalytics()
  const { isDismissed, dismiss } = useAnalyticsNotice()

  // Нет `viewId` — статистика не собирается (предпросмотр, ошибка на бэке, «Не отслеживать»)
  if (!viewId) return null

  return (
    <Slide in={!isDismissed} direction="up" appear mountOnEnter unmountOnExit>
      <Stack
        data-print-hidden
        role="dialog"
        aria-label={translate('title')}
        spacing={1.5}
        sx={{
          position: 'fixed',
          right: { xs: 12, sm: 24 },
          bottom: { xs: 12, sm: 24 },
          left: { xs: 12, sm: 'auto' },
          maxWidth: { sm: 380 },
          p: 2.5,
          zIndex: 1200,
          ...glassCard,
          backgroundImage: 'none',
        }}
      >
        <Stack direction="row" spacing={1} alignItems="center">
          <InsightsRounded color="primary" fontSize="small" />
          <Typography variant="h6" sx={{ fontSize: '1.05rem' }}>
            {translate('title')}
          </Typography>
        </Stack>

        <Typography variant="body2">{translate('text')}</Typography>

        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
          {translate('dnt')}
        </Typography>

        <Button
          onClick={dismiss}
          variant="contained"
          sx={{ alignSelf: 'flex-end', borderRadius: 999, px: 3, textTransform: 'none' }}
        >
          {translate('ok')}
        </Button>
      </Stack>
    </Slide>
  )
}
