import { useState } from 'react'
import { useLogout } from '@features/auth/hooks'
import { Portrait } from '@features/curriculumVitae/components/AboutSection/Portrait'
import { EditProfileDialog } from '@features/curriculumVitae/components/EditProfileDialog'
import { glassCard } from '@features/curriculumVitae/styles'
import { LkProfile } from '@features/curriculumVitae/types'
import { EditRounded, LogoutRounded } from '@mui/icons-material'
import { Chip, IconButton, Stack, Typography } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { useNavigate } from 'react-router'
import { useProfileCard } from './hooks'

type Props = {
  userId: string
  profile: LkProfile
}

export const ProfileCard = ({ userId, profile }: Props) => {
  const translate = useTranslate('Lk.Profile')
  const translateHeader = useTranslate('Header')
  const { localized, fio, initials, isCompact, avatarSize } = useProfileCard(profile)
  const [isEditing, setIsEditing] = useState(false)
  const { signOut, isLoggingOut } = useLogout()
  const navigate = useNavigate()

  // Глобальный Header скрыт на всех /cv/*-страницах, поэтому выход из аккаунта живёт здесь —
  // это единственное место в ЛК, которое остаётся на экране при любом разделе.
  const handleLogout = async () => {
    await signOut()
    navigate(ROUTES.home)
  }

  return (
    <Stack
      direction={isCompact ? 'row' : 'column'}
      spacing={isCompact ? 2.5 : 2}
      alignItems="center"
      textAlign={isCompact ? 'left' : 'center'}
      sx={{ p: isCompact ? 2.5 : 4, position: 'relative', ...glassCard }}
    >
      <Stack direction="row" spacing={0.5} sx={{ position: 'absolute', top: 8, right: 8 }}>
        <IconButton aria-label={translate('edit')} onClick={() => setIsEditing(true)}>
          <EditRounded fontSize="small" />
        </IconButton>

        <IconButton
          aria-label={translateHeader('logout')}
          onClick={handleLogout}
          disabled={isLoggingOut}
        >
          <LogoutRounded fontSize="small" />
        </IconButton>
      </Stack>

      <Portrait src={profile.photoPath ?? ''} alt={fio} initials={initials} size={avatarSize} />

      <Stack spacing={0.5} alignItems={isCompact ? 'flex-start' : 'center'} sx={{ minWidth: 0 }}>
        <Typography variant="h5" sx={{ fontSize: isCompact ? '1.15rem' : undefined }}>
          {fio}
        </Typography>

        {localized.position && (
          <Typography sx={{ color: 'primary.main' }}>{localized.position}</Typography>
        )}

        {localized.location && <Chip size="small" label={localized.location} />}
      </Stack>

      {isEditing && (
        <EditProfileDialog userId={userId} profile={profile} onClose={() => setIsEditing(false)} />
      )}
    </Stack>
  )
}
