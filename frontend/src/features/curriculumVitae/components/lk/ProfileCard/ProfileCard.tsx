import { useState } from 'react'
import { Portrait } from '@features/curriculumVitae/components/cv/AboutSection/Portrait'
import { EditProfileDialog } from '@features/curriculumVitae/components/lk/EditProfileDialog'
import { glassCard, textAccent } from '@features/curriculumVitae/styles'
import { LkProfile } from '@features/curriculumVitae/types'
import { EditRounded } from '@mui/icons-material'
import { Chip, IconButton, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useProfileCard } from './hooks'

type Props = {
  userId: string
  profile: LkProfile
}

export const ProfileCard = ({ userId, profile }: Props) => {
  const translate = useTranslate('Lk.Profile')
  const { localized, fio, initials, isCompact, avatarSize } = useProfileCard(profile)
  const [isEditing, setIsEditing] = useState(false)

  return (
    <Stack
      direction={isCompact ? 'row' : 'column'}
      spacing={isCompact ? 2.5 : 2}
      alignItems="center"
      textAlign={isCompact ? 'left' : 'center'}
      sx={{ p: isCompact ? 2.5 : 4, position: 'relative', ...glassCard }}
    >
      <IconButton
        aria-label={translate('edit')}
        onClick={() => setIsEditing(true)}
        sx={{ position: 'absolute', top: 8, right: 8 }}
      >
        <EditRounded fontSize="small" />
      </IconButton>

      <Portrait src={profile.photoPath ?? ''} alt={fio} initials={initials} size={avatarSize} />

      <Stack spacing={0.5} alignItems={isCompact ? 'flex-start' : 'center'} sx={{ minWidth: 0 }}>
        <Typography variant="h5" sx={{ fontSize: isCompact ? '1.15rem' : undefined }}>
          {fio}
        </Typography>

        {localized.position && <Typography sx={textAccent}>{localized.position}</Typography>}

        {localized.location && <Chip size="small" label={localized.location} />}
      </Stack>

      {isEditing && (
        <EditProfileDialog userId={userId} profile={profile} onClose={() => setIsEditing(false)} />
      )}
    </Stack>
  )
}
