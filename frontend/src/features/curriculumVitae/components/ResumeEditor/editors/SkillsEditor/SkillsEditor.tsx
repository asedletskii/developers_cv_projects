import { Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useSkillOptions } from './hooks'
import { SkillsField } from './SkillsField'

/** Вкладка «Навыки»: профессиональные и личные, оба — фиксированные наборы значений. */
export const SkillsEditor = () => {
  const translate = useTranslate('Lk.Edit.Skills')
  const { hardSkillOptions, softSkillOptions } = useSkillOptions()

  return (
    <Stack spacing={2}>
      <SkillsField name="skills.hardSkills" label={translate('hard')} options={hardSkillOptions} />

      <SkillsField name="skills.softSkills" label={translate('soft')} options={softSkillOptions} />
    </Stack>
  )
}
