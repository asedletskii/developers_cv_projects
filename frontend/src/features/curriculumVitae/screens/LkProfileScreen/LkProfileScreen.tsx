import { useEffect } from 'react'
import { LoadFailure, ProfileCard, ProfilePlaceholder, ResumeList } from '@features/curriculumVitae/components'
import {
  clearLkState,
  getLkProfile,
  getLkResumes,
  selectIsLkProfileLoading,
  selectIsLkResumesLoading,
  selectLkProfile,
  selectLkProfileError,
  selectLkResumes,
  selectLkResumesError,
} from '@features/curriculumVitae/store/lk'
import { Container, Stack } from '@mui/material'
import { useAppDispatch, useAppSelector } from '@shared/store'
import { useParams } from 'react-router'

/** Личный кабинет: ФИО и общие данные слева (1/4), список резюме пользователя справа (3/4). */
export default function LkProfileScreen() {
  const dispatch = useAppDispatch()
  const { userId } = useParams()

  const isProfileLoading = useAppSelector(selectIsLkProfileLoading)
  const isResumesLoading = useAppSelector(selectIsLkResumesLoading)
  const profileError = useAppSelector(selectLkProfileError)
  const resumesError = useAppSelector(selectLkResumesError)
  const profile = useAppSelector(selectLkProfile)
  const resumes = useAppSelector(selectLkResumes)

  useEffect(() => {
    if (!userId) return

    const profileRequest = dispatch(getLkProfile(userId))
    const resumesRequest = dispatch(getLkResumes(userId))

    return () => {
      profileRequest.abort()
      resumesRequest.abort()
    }
  }, [dispatch, userId])

  useEffect(() => {
    return () => {
      dispatch(clearLkState())
    }
  }, [dispatch])

  if (isProfileLoading || isResumesLoading) return <ProfilePlaceholder />

  const error = profileError ?? resumesError
  if (error) return <LoadFailure error={error} subject="profile" />

  // Сюда не попадаем: без ошибки и загрузки profile/resumes уже гарантированно загружены
  if (!profile || !resumes) return null

  return (
    <Container maxWidth="lg" sx={{ py: 6 }}>
      <Stack direction={{ xs: 'column', md: 'row' }} spacing={4} alignItems="flex-start">
        <Stack sx={{ width: { xs: '100%', md: '25%' }, flexShrink: 0 }}>
          <ProfileCard profile={profile} />
        </Stack>

        <Stack sx={{ width: '100%', flex: 1 }}>
          <ResumeList userId={profile.userId} resumes={resumes} />
        </Stack>
      </Stack>
    </Container>
  )
}
