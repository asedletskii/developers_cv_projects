import { CreateLinkPrompt } from '@features/curriculumVitae/components/lk/CreateLinkPrompt'
import { useResumeLinks } from '@features/curriculumVitae/hooks'
import { glassCard } from '@features/curriculumVitae/styles'
import { Resume, ResumeEditorTab, ResumeStatistics } from '@features/curriculumVitae/types'
import { formatVisitTime } from '@features/curriculumVitae/utils'
import { Button, LinearProgress, Stack, Typography } from '@mui/material'
import { useLanguage, useTranslate } from '@shared/hooks'
import dayjs from 'dayjs'
import { ChecklistItem } from './ChecklistItem'
import { useResumeChecklist } from './hooks'

type Props = {
  view: Resume
  linksState: ReturnType<typeof useResumeLinks>
  statistics?: ResumeStatistics
  onOpenTab: (tab: ResumeEditorTab) => void
  onOpenAnalytics: () => void
}

export const MainTab = ({ view, linksState, statistics, onOpenTab, onOpenAnalytics }: Props) => {
  const translate = useTranslate('Lk.Main')
  const { language } = useLanguage()
  const { items, percent } = useResumeChecklist(view)

  return (
    <Stack spacing={2}>
      {!linksState.isLoading && !linksState.hasLoadError && linksState.links.length === 0 && (
        <CreateLinkPrompt onCreate={linksState.create} />
      )}

      {statistics && linksState.links.length > 0 && (
        <Stack spacing={1} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
          <Stack direction="row" alignItems="center" justifyContent="space-between">
            <Typography variant="h6">{translate('Stats.title')}</Typography>

            <Button onClick={onOpenAnalytics} sx={{ textTransform: 'none' }}>
              {translate('Stats.more')}
            </Button>
          </Stack>

          <Typography>
            {statistics.views === 0
              ? translate('Stats.empty')
              : translate('Stats.summary', {
                  views: statistics.views,
                  reads: statistics.reads,
                  contacts: statistics.contactViews,
                })}
          </Typography>

          {statistics.lastViewAt && (
            <Typography variant="body2" sx={{ color: 'text.secondary' }}>
              {translate('Stats.last', {
                time: formatVisitTime(statistics.lastViewAt, language, { isCapitalized: false }),
              })}
            </Typography>
          )}
        </Stack>
      )}

      <Stack spacing={3} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
        <Stack spacing={1}>
          <Typography variant="h6">{translate('completeness', { percent })}</Typography>

          <LinearProgress
            variant="determinate"
            value={percent}
            sx={{ height: 8, borderRadius: 4 }}
          />

          <Typography variant="body2" sx={{ color: 'text.secondary' }}>
            {translate('updated', { time: dayjs(view.updatedAt).fromNow() })}
          </Typography>
        </Stack>

        <Stack spacing={0.5}>
          {items.map(([tab, isDone]) => (
            <ChecklistItem key={tab} tab={tab} isDone={isDone} onFill={() => onOpenTab(tab)} />
          ))}
        </Stack>
      </Stack>
    </Stack>
  )
}
