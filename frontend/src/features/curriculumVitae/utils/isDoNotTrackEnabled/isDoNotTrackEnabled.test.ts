import { isDoNotTrackEnabled } from './isDoNotTrackEnabled'

describe('isDoNotTrackEnabled', () => {
  it('Do Not Track = 1', () => {
    expect(isDoNotTrackEnabled({ doNotTrack: '1' })).toBe(true)
  })

  it('Global Privacy Control', () => {
    expect(isDoNotTrackEnabled({ globalPrivacyControl: true })).toBe(true)
  })

  it('без сигналов или с DNT = 0 — отслеживание разрешено', () => {
    expect(isDoNotTrackEnabled({})).toBe(false)
    expect(isDoNotTrackEnabled({ doNotTrack: '0' })).toBe(false)
    expect(isDoNotTrackEnabled({ doNotTrack: null, globalPrivacyControl: false })).toBe(false)
    expect(isDoNotTrackEnabled(undefined)).toBe(false)
  })
})
