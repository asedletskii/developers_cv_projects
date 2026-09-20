import { Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { formatDuration, formatMonth } from '@shared/utils'

type Props = {
  startDate: string
  endDate?: string
}

/** «март 2023 — по настоящее время · 3 года 7 месяцев» */
export const Dates = ({ startDate, endDate }: Props) => {
  const translate = useTranslate('CurriculumVitae.Experience')

  return (
    <Typography variant="caption" sx={{ color: 'gray.light' }}>
      {formatMonth(startDate)} — {endDate ? formatMonth(endDate) : translate('present')} ·{' '}
      {formatDuration(startDate, endDate)}
    </Typography>
  )
}
