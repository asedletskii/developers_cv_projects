import { Section } from '@features/curriculumVitae/components/Section'
import { glassCard } from '@features/curriculumVitae/styles'
import { Person, SectionId } from '@features/curriculumVitae/types'
import { Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { getPersonFIO } from '@shared/utils'
import { Facts } from './Facts'
import { Portrait } from './Portrait'
import { Summary } from './Summary'

type Props = Pick<
  Person,
  | 'name'
  | 'surname'
  | 'patronymic'
  | 'position'
  | 'photoPath'
  | 'dateOfBirth'
  | 'location'
  | 'relocation'
  | 'workFormat'
  | 'languages'
>

export const AboutSection = ({
  name,
  surname,
  patronymic,
  position,
  photoPath,
  ...facts
}: Props) => {
  const translate = useTranslate('CurriculumVitae.About')

  const fio = getPersonFIO({ name, surname, patronymic })

  return (
    <Section id={SectionId.About} title={translate('title')} collapsible={false}>
      <Stack
        direction={{ xs: 'column', sm: 'row' }}
        spacing={{ xs: 3, sm: 4 }}
        alignItems={{ xs: 'center', sm: 'flex-start' }}
        sx={{ p: { xs: 3, sm: 4 }, ...glassCard }}
      >
        <Portrait src={photoPath} alt={fio} initials={`${name[0]}${surname[0]}`} />

        <Stack spacing={2} sx={{ minWidth: 0 }}>
          <Summary fio={fio} position={position} />

          <Facts {...facts} />
        </Stack>
      </Stack>
    </Section>
  )
}
