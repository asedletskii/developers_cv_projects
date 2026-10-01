import { ResumeEditorTab } from '@features/curriculumVitae/types'
import { CheckCircleRounded, RadioButtonUncheckedRounded } from '@mui/icons-material'
import { Button, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  tab: ResumeEditorTab
  isDone: boolean
  onFill: () => void
}

export const ChecklistItem = ({ tab, isDone, onFill }: Props) => {
  const translate = useTranslate('Lk.Main')
  const translateTabs = useTranslate('Lk.Edit.Tabs')

  return (
    <Stack direction="row" alignItems="center" spacing={1.5}>
      {isDone ? (
        <CheckCircleRounded fontSize="small" color="success" />
      ) : (
        <RadioButtonUncheckedRounded fontSize="small" sx={{ color: 'text.disabled' }} />
      )}

      <Typography sx={{ flex: 1, color: isDone ? 'text.secondary' : 'text.primary' }}>
        {translateTabs(tab)}
      </Typography>

      {!isDone && (
        <Button size="small" onClick={onFill} sx={{ textTransform: 'none' }}>
          {translate('fill')}
        </Button>
      )}
    </Stack>
  )
}
