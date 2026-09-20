import { HardSkill, SoftSkill } from '@features/curriculumVitae/types'
import { isHardSkill, isSoftSkill } from '@features/curriculumVitae/utils'
import { useTranslate } from '@shared/hooks'

export const useSkillOptions = () => {
  const translate = useTranslate('CurriculumVitae.Skills')

  return {
    hardSkillOptions: Object.values(HardSkill) as string[],
    softSkillOptions: Object.values(SoftSkill) as string[],
    getHardSkillLabel: (skill: string) =>
      isHardSkill(skill) ? translate(`HardSkills.${skill}`) : skill,
    getSoftSkillLabel: (skill: string) =>
      isSoftSkill(skill) ? translate(`SoftSkills.${skill}`) : skill,
  }
}
