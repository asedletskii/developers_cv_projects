import { Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useSkillOptions } from './hooks'
import { SkillsField } from './SkillsField'

export const SkillsEditor = () => {
  const translate = useTranslate('Lk.Edit.Skills')
  const { hardSkillOptions, softSkillOptions, getHardSkillLabel, getSoftSkillLabel } =
    useSkillOptions()

  return (
    <Stack spacing={2}>
      <SkillsField
        name="skills.hardSkills"
        label={translate('hard')}
        helperText={translate('hint')}
        options={hardSkillOptions}
        getOptionLabel={getHardSkillLabel}
      />

      <SkillsField
        name="skills.softSkills"
        label={translate('soft')}
        helperText={translate('hint')}
        options={softSkillOptions}
        getOptionLabel={getSoftSkillLabel}
      />
    </Stack>
  )
}
