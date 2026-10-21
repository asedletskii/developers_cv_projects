export const SKILL_MAX_LENGTH = 50

export const normalizeSkills = (skills: string[]): string[] => {
  const seen = new Set<string>()

  return skills
    .map((skill) => skill.trim().slice(0, SKILL_MAX_LENGTH))
    .filter((skill) => {
      const key = skill.toLowerCase()
      if (!skill || seen.has(key)) return false
      seen.add(key)
      return true
    })
}
