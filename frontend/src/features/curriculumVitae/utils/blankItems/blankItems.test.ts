import { blankExperience, blankProject } from './blankItems'

describe('blankItems', () => {
  it('каждая новая запись получает свой id', () => {
    expect(blankExperience().id).not.toBe(blankExperience().id)
  })

  it('у пустого проекта пустой стэк', () => {
    expect(blankProject().stack).toEqual([])
  })
})
