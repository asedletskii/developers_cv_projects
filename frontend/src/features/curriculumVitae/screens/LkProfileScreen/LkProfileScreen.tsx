import { LoginPrompt } from '@features/auth/components/LoginPrompt'
import { useSession } from '@features/auth/hooks'
import { LoadFailure, ProfileCard, ProfilePlaceholder } from '@features/curriculumVitae/components'
import { EditProfileDialog } from '@features/curriculumVitae/components/lk/EditProfileDialog'
import { useLkProfile } from '@features/curriculumVitae/hooks'
import { Container, Stack } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useNoIndex, useTranslate } from '@shared/hooks'
import { Outlet, useNavigate } from 'react-router'

const LkProfileScreen = () => {
  const translate = useTranslate('Lk.Profile')
  const navigate = useNavigate()
  useNoIndex()
  const { userId, isLoading: isSessionLoading, isAuthenticated } = useSession()
  const {
    profile,
    isLoading: isProfileLoading,
    error,
  } = useLkProfile(isAuthenticated ? (userId ?? undefined) : undefined)

  if (isSessionLoading) return <ProfilePlaceholder />
  if (!isAuthenticated || !userId) return <LoginPrompt />
  if (isProfileLoading) return <ProfilePlaceholder />

  if (error?.kind === 'notFound' && !profile) {
    return (
      <EditProfileDialog
        userId={userId}
        profile={{ userId, ru: {}, en: {} }}
        intro={translate('createText')}
        onClose={() => navigate(ROUTES.home)}
      />
    )
  }

  if (error) return <LoadFailure error={error} subject="profile" />

  if (!profile) return null

  return (
    <Container maxWidth="lg" sx={{ py: { xs: 3, md: 6 } }}>
      <Stack
        direction={{ xs: 'column', md: 'row' }}
        spacing={{ xs: 3, md: 4 }}
        alignItems="flex-start"
      >
        <Stack sx={{ width: { xs: '100%', md: '25%' }, flexShrink: 0 }}>
          <ProfileCard userId={profile.userId} profile={profile} />
        </Stack>

        <Stack sx={{ width: '100%', flex: 1, minWidth: 0 }}>
          <Outlet context={{ userId: profile.userId }} />
        </Stack>
      </Stack>
    </Container>
  )
}

export default LkProfileScreen
