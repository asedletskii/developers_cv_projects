import { useEffect } from 'react'
import { CurriculumVitae, LoadFailure, LoadPlaceholder } from '@features/curriculumVitae/components'
import {
  clearCurriculumVitaeState,
  getPersonData,
  selectIsPersonDataLoading,
  selectPersonDataError,
} from '@features/curriculumVitae/store/curriculumVitae'
import { useLanguage } from '@shared/hooks'
import { useAppDispatch, useAppSelector } from '@shared/store'
import { useParams } from 'react-router'

export default function CurriculumVitaeScreen() {
  const dispatch = useAppDispatch()
  const { cvId } = useParams()
  const { language } = useLanguage()

  const isLoading = useAppSelector(selectIsPersonDataLoading)
  const error = useAppSelector(selectPersonDataError)

  // Резюме выбирается по id из ссылки, а тексты приходят на выбранном языке —
  // поэтому и смена id, и смена языка запрашивают данные заново
  useEffect(() => {
    if (cvId) dispatch(getPersonData({ cvId, lang: language }))
  }, [dispatch, cvId, language])

  useEffect(() => {
    return () => {
      dispatch(clearCurriculumVitaeState())
    }
  }, [dispatch])

  if (isLoading) return <LoadPlaceholder />

  if (error) return <LoadFailure error={error} />

  return <CurriculumVitae />
}
