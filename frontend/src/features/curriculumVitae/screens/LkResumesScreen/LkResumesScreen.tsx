import { ResumeList } from '@features/curriculumVitae/components'
import { useLkOutletContext, useLkResumes } from '@features/curriculumVitae/hooks'
import { Alert, LinearProgress } from '@mui/material'
import { useTranslate } from '@shared/hooks'

const LkResumesScreen = () => {
  const translate = useTranslate('Lk.Error')
  const { userId } = useLkOutletContext()
  const { resumes, error } = useLkResumes(userId)

  if (error) return <Alert severity="error">{translate('failedTitle')}</Alert>
  if (!resumes || !userId) return <LinearProgress />

  return <ResumeList userId={userId} resumes={resumes} />
}

export default LkResumesScreen
