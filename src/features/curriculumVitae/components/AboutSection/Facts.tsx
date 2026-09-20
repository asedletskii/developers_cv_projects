import { Person } from '@features/curriculumVitae/types'
import {
  CakeOutlined,
  FlightTakeoffOutlined,
  LocationOnOutlined,
  TranslateOutlined,
  WorkOutline,
} from '@mui/icons-material'
import { Chip, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { formatAge } from '@shared/utils'

type Props = Pick<Person, 'dateOfBirth' | 'location' | 'relocation' | 'workFormat' | 'languages'>

/** Короткие факты о человеке в виде чипов; появляются каскадом (см. useScrollReveal). */
export const Facts = ({ dateOfBirth, location, relocation, workFormat, languages }: Props) => {
  const translate = useTranslate('CurriculumVitae.About')

  return (
    <Stack
      className="reveal-group"
      direction="row"
      flexWrap="wrap"
      gap={1}
      justifyContent={{ xs: 'center', sm: 'flex-start' }}
    >
      <Chip className="reveal-item" icon={<LocationOnOutlined />} label={location} />

      <Chip className="reveal-item" icon={<CakeOutlined />} label={formatAge(dateOfBirth)} />

      <Chip
        className="reveal-item"
        icon={<FlightTakeoffOutlined />}
        label={translate(relocation ? 'relocation' : 'noRelocation')}
      />

      {workFormat.map((format) => (
        <Chip
          key={format}
          className="reveal-item"
          icon={<WorkOutline />}
          label={translate(`WorkFormat.${format}`)}
        />
      ))}

      {languages.map(({ language, grade }) => (
        <Chip
          key={language}
          className="reveal-item"
          icon={<TranslateOutlined />}
          label={[translate(`Languages.${language}`), grade].filter(Boolean).join(' · ')}
        />
      ))}
    </Stack>
  )
}
