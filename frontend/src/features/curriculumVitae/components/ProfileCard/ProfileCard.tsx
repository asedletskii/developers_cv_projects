import { Portrait } from '@features/curriculumVitae/components/AboutSection/Portrait'
import { glassCard } from '@features/curriculumVitae/styles'
import { LkProfile } from '@features/curriculumVitae/types'
import { Chip, Stack, Typography } from '@mui/material'
import { useProfileCard } from './hooks'

type Props = {
  profile: LkProfile
}

/** Левая колонка ЛК: фото, ФИО и общие данные владельца. На узком экране — компактная строка. */
export const ProfileCard = ({ profile }: Props) => {
  const { fio, initials, isCompact, avatarSize } = useProfileCard(profile)

  return (
    <Stack
      direction={isCompact ? 'row' : 'column'}
      spacing={isCompact ? 2.5 : 2}
      alignItems="center"
      textAlign={isCompact ? 'left' : 'center'}
      sx={{ p: isCompact ? 2.5 : 4, ...glassCard }}
    >
      <Portrait src={profile.photoPath ?? ''} alt={fio} initials={initials} size={avatarSize} />

      <Stack spacing={0.5} alignItems={isCompact ? 'flex-start' : 'center'} sx={{ minWidth: 0 }}>
        <Typography variant="h5" sx={{ fontSize: isCompact ? '1.15rem' : undefined }}>
          {fio}
        </Typography>

        {profile.position && (
          <Typography sx={{ color: 'primary.main' }}>{profile.position}</Typography>
        )}

        {profile.location && <Chip size="small" label={profile.location} />}
      </Stack>
    </Stack>
  )
}
