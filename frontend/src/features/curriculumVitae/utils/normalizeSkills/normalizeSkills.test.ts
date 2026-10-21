import { normalizeSkills, SKILL_MAX_LENGTH } from './normalizeSkills'

describe('normalizeSkills', () => {
  it('обрезает пробелы, убирает пустые и повторы без учёта регистра, порядок сохраняет', () => {
    expect(normalizeSkills([' testing ', '', '  ', 'Docker', 'docker', 'Testing'])).toEqual([
      'testing',
      'Docker',
    ])
  })

  it('режет слишком длинный текст', () => {
    expect(normalizeSkills(['a'.repeat(SKILL_MAX_LENGTH + 10)])[0]).toHaveLength(SKILL_MAX_LENGTH)
  })
})
