import { ResumeStatus } from '@features/curriculumVitae/types'
import { ChipProps } from '@mui/material'

/** Цвет чипа по статусу: опубликовано — заметно зелёным, остальное — приглушённо */
export const RESUME_STATUS_COLOR: Record<ResumeStatus, ChipProps['color']> = {
  draft: 'default',
  opened: 'success',
  closed: 'warning',
}
