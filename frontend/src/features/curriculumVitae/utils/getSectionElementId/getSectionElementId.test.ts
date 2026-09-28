import { getContentId, getTitleId } from './getSectionElementId'

describe('getSectionElementId', () => {
  it('строит id заголовка и тела по id секции', () => {
    expect(getTitleId('experience')).toBe('experience-title')
    expect(getContentId('experience')).toBe('experience-content')
  })
})
