import { ResumeList } from '@features/curriculumVitae/components'
import { useLkResumes } from '@features/curriculumVitae/hooks'
import { Alert, LinearProgress } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useParams } from 'react-router'

const LkResumesScreen = () => {
  const translate = useTranslate('Lk.Error')
  const { userId } = useParams()
  const { resumes, error } = useLkResumes(userId)

  if (error) return <Alert severity="error">{translate('failedTitle')}</Alert>
  if (!resumes || !userId) return <LinearProgress />

  return <ResumeList userId={userId} resumes={resumes} />
}

export default LkResumesScreen
