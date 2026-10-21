import { HardSkill, SoftSkill } from '@features/curriculumVitae/types'

const HARD_SKILLS = new Set<string>(Object.values(HardSkill))
const SOFT_SKILLS = new Set<string>(Object.values(SoftSkill))

export const isHardSkill = (value: string): value is HardSkill => HARD_SKILLS.has(value)

export const isSoftSkill = (value: string): value is SoftSkill => SOFT_SKILLS.has(value)
