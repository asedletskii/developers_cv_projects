import { CurriculumVitae, LoadFailure, LoadPlaceholder } from '@features/curriculumVitae/components'
import { usePublicResume } from '@features/curriculumVitae/hooks'
import { useParams } from 'react-router'

const CurriculumVitaeScreen = () => {
  const { cvId } = useParams()
  const { isLoading, error } = usePublicResume(cvId)

  if (isLoading) return <LoadPlaceholder />
  if (error) return <LoadFailure error={error} />

  return <CurriculumVitae />
}

export default CurriculumVitaeScreen
