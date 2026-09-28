import { LkProfile } from '@features/curriculumVitae/types'
import { Languages } from '@shared/types'
import { getLocalizedProfile } from './getLocalizedProfile'

const profile: LkProfile = {
  userId: '1',
  ru: { name: 'Антон', surname: 'Седлецкий', position: 'Frontend разработчик' },
  en: { name: 'Anton' },
}

describe('getLocalizedProfile', () => {
  it('отдаёт русские поля для ru', () => {
    expect(getLocalizedProfile(profile, Languages.Ru)).toEqual({
      name: 'Антон',
      surname: 'Седлецкий',
      patronymic: undefined,
      position: 'Frontend разработчик',
      location: undefined,
    })
  })

  it('для en берёт то, что заполнено на английском', () => {
    expect(getLocalizedProfile(profile, Languages.En).name).toBe('Anton')
  })

  it('незаполненные на en поля подставляет из ru', () => {
    expect(getLocalizedProfile(profile, Languages.En).surname).toBe('Седлецкий')
    expect(getLocalizedProfile(profile, Languages.En).position).toBe('Frontend разработчик')
  })
})
