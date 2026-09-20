import { ReactElement } from 'react'
import { ResumeLanguageItem, ResumeProfile } from '@features/curriculumVitae/types'
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

type Props = Pick<ResumeProfile, 'dateOfBirth' | 'location' | 'relocation' | 'workFormat'> & {
  languages: ResumeLanguageItem[]
}

type Fact = { key: string; icon: ReactElement; label: string; preview: 'profile' | 'languages' }

export const useFacts = ({
  dateOfBirth,
  location,
  relocation,
  workFormat,
  languages,
}: Props): Fact[] => {
  const translate = useTranslate('CurriculumVitae.About')

  return [
    { key: 'location', icon: <LocationOnOutlined />, label: location, preview: 'profile' },
    { key: 'age', icon: <CakeOutlined />, label: formatAge(dateOfBirth), preview: 'profile' },
    {
      key: 'relocation',
      icon: <FlightTakeoffOutlined />,
      label: translate(relocation ? 'relocation' : 'noRelocation'),
      preview: 'profile',
    },
    ...workFormat.map((format) => ({
      key: format,
      icon: <WorkOutline />,
      label: translate(`WorkFormat.${format}`),
      preview: 'profile' as const,
    })),
    ...languages.map(({ language, grade }) => ({
      key: language,
      icon: <TranslateOutlined />,
      label: [translate(`Languages.${language}`), grade].filter(Boolean).join(' · '),
      preview: 'languages' as const,
    })),
  ]
}

export const Facts = (props: Props) => {
  const facts = useFacts(props)

  return (
    <Stack
      className="reveal-group"
      direction="row"
      flexWrap="wrap"
      gap={1}
      justifyContent={{ xs: 'center', sm: 'flex-start' }}
    >
      {facts.map(({ key, icon, label, preview }) => (
        <Chip key={key} className="reveal-item" data-preview={preview} icon={icon} label={label} />
      ))}
    </Stack>
  )
}
