import { Portrait } from '@features/curriculumVitae/components/AboutSection/Portrait'
import { glassCard } from '@features/curriculumVitae/styles'
import { LkProfile } from '@features/lk/types'
import { Chip, Stack, Typography } from '@mui/material'
import { getPersonFIO } from '@shared/utils'

type Props = {
  profile: LkProfile
}

/** Левая колонка ЛК: фото, ФИО и общие данные владельца. */
export const ProfileCard = ({ profile }: Props) => {
  const fio = getPersonFIO(profile)

  return (
    <Stack spacing={2} alignItems="center" textAlign="center" sx={{ p: 4, ...glassCard }}>
      <Portrait
        src={profile.photoPath ?? ''}
        alt={fio}
        initials={`${profile.name[0]}${profile.surname[0]}`}
      />

      <Stack spacing={0.5}>
        <Typography variant="h5">{fio}</Typography>

        {profile.position && (
          <Typography sx={{ color: 'primary.main' }}>{profile.position}</Typography>
        )}
      </Stack>

      {profile.location && <Chip label={profile.location} />}
    </Stack>
  )
}
