import { HardSkill, SoftSkill } from '@features/curriculumVitae/types'
import { Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { SelectField } from './SelectField'

/** Секция «Навыки»: профессиональные и личные, оба — фиксированные наборы значений. */
export const SkillsEditor = () => {
  const translate = useTranslate('Lk.Edit.Skills')
  const translateSkills = useTranslate('CurriculumVitae.Skills')

  return (
    <Stack spacing={2}>
      <SelectField
        name="hardSkills"
        multiple
        label={translate('hard')}
        options={Object.values(HardSkill).map((skill) => ({
          value: skill,
          label: translateSkills(`HardSkills.${skill}`),
        }))}
      />

      <SelectField
        name="softSkills"
        multiple
        label={translate('soft')}
        options={Object.values(SoftSkill).map((skill) => ({
          value: skill,
          label: translateSkills(`SoftSkills.${skill}`),
        }))}
      />
    </Stack>
  )
}
