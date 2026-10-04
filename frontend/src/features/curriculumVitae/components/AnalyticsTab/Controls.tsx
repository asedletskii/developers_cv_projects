import { AnalyticsPeriod, ResumeLink } from '@features/curriculumVitae/types'
import { RefreshRounded } from '@mui/icons-material'
import {
  IconButton,
  MenuItem,
  Stack,
  TextField,
  ToggleButton,
  ToggleButtonGroup,
  Tooltip,
} from '@mui/material'
import { useTranslate } from '@shared/hooks'

const PERIODS: AnalyticsPeriod[] = [7, 30, 90]

type Props = {
  links: ResumeLink[]
  scope: string
  allValue: string
  period: AnalyticsPeriod
  onScopeChange: (scope: string) => void
  onPeriodChange: (period: AnalyticsPeriod) => void
  onRefresh: () => void
}

export const Controls = ({
  links,
  scope,
  allValue,
  period,
  onScopeChange,
  onPeriodChange,
  onRefresh,
}: Props) => {
  const translate = useTranslate('Lk.Analytics')

  return (
    <Stack direction="row" spacing={1.5} alignItems="center" flexWrap="wrap" useFlexGap>
      <TextField
        select
        size="small"
        value={scope}
        onChange={(event) => onScopeChange(event.target.value)}
        sx={{ minWidth: 200 }}
      >
        <MenuItem value={allValue}>{translate('allLinks')}</MenuItem>
        {links.map((link) => (
          <MenuItem key={link.id} value={link.id}>
            {link.name}
          </MenuItem>
        ))}
      </TextField>

      <ToggleButtonGroup
        exclusive
        size="small"
        value={period}
        onChange={(_, value: AnalyticsPeriod | null) => value && onPeriodChange(value)}
        aria-label={translate('period')}
      >
        {PERIODS.map((days) => (
          <ToggleButton key={days} value={days} sx={{ textTransform: 'none', px: 1.5 }}>
            {translate('days', { count: days })}
          </ToggleButton>
        ))}
      </ToggleButtonGroup>

      <Tooltip title={translate('refresh')}>
        <IconButton onClick={onRefresh} aria-label={translate('refresh')}>
          <RefreshRounded />
        </IconButton>
      </Tooltip>
    </Stack>
  )
}
