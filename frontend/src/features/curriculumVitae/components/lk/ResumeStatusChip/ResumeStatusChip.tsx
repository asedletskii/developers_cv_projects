import { ResumeStatus } from '@features/curriculumVitae/types'
import { Chip } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { RESUME_STATUS_COLOR } from './constants'

type Props = {
  status: ResumeStatus
}

export const ResumeStatusChip = ({ status }: Props) => {
  const translate = useTranslate('Lk.Resumes.Status')

  return <Chip label={translate(status)} color={RESUME_STATUS_COLOR[status]} size="small" />
}
