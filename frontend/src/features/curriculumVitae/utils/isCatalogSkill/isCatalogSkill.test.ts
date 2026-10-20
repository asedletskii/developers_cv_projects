import { isHardSkill, isSoftSkill } from './isCatalogSkill'

describe('isCatalogSkill', () => {
  it('узнаёт значения справочника и не путает группы', () => {
    expect(isHardSkill('testing')).toBe(true)
    expect(isSoftSkill('teamwork')).toBe(true)
    expect(isHardSkill('teamwork')).toBe(false)
    expect(isSoftSkill('testing')).toBe(false)
  })

  it('свободный текст не из справочника', () => {
    expect(isHardSkill('Сборка в Docker')).toBe(false)
    expect(isSoftSkill('')).toBe(false)
  })
})
