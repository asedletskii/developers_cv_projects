import { getRuntimeConfig } from './getRuntimeConfig'

describe('getRuntimeConfig', () => {
  it('без конфига — дефолты (уважаем «Не отслеживать»)', () => {
    expect(getRuntimeConfig(undefined)).toEqual({ respectDoNotTrack: true })
    expect(getRuntimeConfig({})).toEqual({ respectDoNotTrack: true })
  })

  it('берёт значение из конфига', () => {
    expect(getRuntimeConfig({ respectDoNotTrack: false })).toEqual({ respectDoNotTrack: false })
    expect(getRuntimeConfig({ respectDoNotTrack: true })).toEqual({ respectDoNotTrack: true })
  })

  it('значение не того типа игнорируется — остаётся дефолт', () => {
    expect(getRuntimeConfig({ respectDoNotTrack: 'false' })).toEqual({ respectDoNotTrack: true })
    expect(getRuntimeConfig({ respectDoNotTrack: 0 })).toEqual({ respectDoNotTrack: true })
  })
})
