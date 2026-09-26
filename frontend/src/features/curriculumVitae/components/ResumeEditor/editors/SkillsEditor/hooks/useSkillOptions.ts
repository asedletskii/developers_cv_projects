import { HardSkill, SoftSkill } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'

/** Справочники навыков с подписями на текущем языке: профессиональные и личные. */
export const useSkillOptions = () => {
  const translate = useTranslate('CurriculumVitae.Skills')

  return {
    hardSkillOptions: Object.values(HardSkill).map((skill) => ({
      value: skill,
      label: translate(`HardSkills.${skill}`),
    })),
    softSkillOptions: Object.values(SoftSkill).map((skill) => ({
      value: skill,
      label: translate(`SoftSkills.${skill}`),
    })),
  }
}
