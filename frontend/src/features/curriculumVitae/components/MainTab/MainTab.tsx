import { glassCard } from '@features/curriculumVitae/styles'
import { ResumeEditorTab, ResumeView } from '@features/curriculumVitae/types'
import { LinearProgress, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import dayjs from 'dayjs'
import { ChecklistItem } from './ChecklistItem'
import { useResumeChecklist } from './hooks'

type Props = {
  view: ResumeView

  /** Открыть вкладку раздела «Редактирование» */
  onOpenTab: (tab: ResumeEditorTab) => void
}

/** «Главная» резюме: насколько оно заполнено и что осталось сделать. */
export const MainTab = ({ view, onOpenTab }: Props) => {
  const translate = useTranslate('Lk.Main')
  const { items, percent } = useResumeChecklist(view)

  return (
    <Stack spacing={3} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Stack spacing={1}>
        <Typography variant="h6">{translate('completeness', { percent })}</Typography>

        <LinearProgress variant="determinate" value={percent} sx={{ height: 8, borderRadius: 4 }} />

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
  )
}
