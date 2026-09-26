import { LoadFailure, ProfileCard, ProfilePlaceholder } from '@features/curriculumVitae/components'
import { useLkProfile } from '@features/curriculumVitae/hooks'
import { Container, Stack } from '@mui/material'
import { Outlet, useParams } from 'react-router'

/**
 * Каркас личного кабинета: ФИО и общие данные слева (1/4) всегда на месте, а правая часть (3/4)
 * меняется по адресу — список резюме (`/cv/lk/:userId`) или редактор резюме (`.../:cvId`).
 */
const LkProfileScreen = () => {
  const { userId } = useParams()
  const { profile, isLoading, error } = useLkProfile(userId)

  if (isLoading) return <ProfilePlaceholder />
  if (error) return <LoadFailure error={error} subject="profile" />

  // Сюда не попадаем: без ошибки и загрузки profile уже гарантированно загружен
  if (!profile) return null

  return (
    <Container maxWidth="lg" sx={{ py: { xs: 3, md: 6 } }}>
      <Stack
        direction={{ xs: 'column', md: 'row' }}
        spacing={{ xs: 3, md: 4 }}
        alignItems="flex-start"
      >
        <Stack sx={{ width: { xs: '100%', md: '25%' }, flexShrink: 0 }}>
          <ProfileCard profile={profile} />
        </Stack>

        <Stack sx={{ width: '100%', flex: 1, minWidth: 0 }}>
          <Outlet />
        </Stack>
      </Stack>
    </Container>
  )
}

export default LkProfileScreen
