import { ResumeStatus } from '@features/curriculumVitae/types'
import { ChipProps } from '@mui/material'

export const RESUME_STATUS_COLOR: Record<ResumeStatus, ChipProps['color']> = {
  draft: 'default',
  opened: 'success',
  closed: 'warning',
}
