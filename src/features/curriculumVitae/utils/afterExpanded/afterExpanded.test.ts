import { SectionId } from '@features/curriculumVitae/types'
import { prefersReducedMotion } from '@shared/animations'
import { afterExpanded } from './afterExpanded'

jest.mock('@shared/animations', () => ({ prefersReducedMotion: jest.fn() }))

const mockedReducedMotion = prefersReducedMotion as jest.Mock

const transitionEnd = (element: Element, propertyName: string) =>
  element.dispatchEvent(Object.assign(new Event('transitionend'), { propertyName }))

describe('afterExpanded', () => {
  let content: HTMLElement

  beforeEach(() => {
    jest.useFakeTimers()
    mockedReducedMotion.mockReturnValue(false)

    document.body.innerHTML = '<section id="courses"><div data-collapsible></div></section>'
    content = document.querySelector('[data-collapsible]') as HTMLElement
  })

  afterEach(() => {
    jest.useRealTimers()
  })

  it('вызывает callback, когда закончилась анимация высоты', () => {
    const callback = jest.fn()

    afterExpanded(SectionId.Courses, callback)
    expect(callback).not.toHaveBeenCalled()

    transitionEnd(content, 'grid-template-rows')

    expect(callback).toHaveBeenCalledTimes(1)
  })

  it('не реагирует на окончание других переходов', () => {
    const callback = jest.fn()

    afterExpanded(SectionId.Courses, callback)
    transitionEnd(content, 'opacity')

    expect(callback).not.toHaveBeenCalled()
  })

  it('вызывает callback только один раз', () => {
    const callback = jest.fn()

    afterExpanded(SectionId.Courses, callback)
    transitionEnd(content, 'grid-template-rows')
    transitionEnd(content, 'grid-template-rows')

    expect(callback).toHaveBeenCalledTimes(1)
  })

  it('если переход не сработал, вызывает callback по таймауту 700 мс', () => {
    const callback = jest.fn()

    afterExpanded(SectionId.Courses, callback)

    jest.advanceTimersByTime(699)
    expect(callback).not.toHaveBeenCalled()

    jest.advanceTimersByTime(1)
    expect(callback).toHaveBeenCalledTimes(1)
  })

  it('после срабатывания таймаута окончание перехода не вызывает callback повторно', () => {
    const callback = jest.fn()

    afterExpanded(SectionId.Courses, callback)
    jest.advanceTimersByTime(700)
    transitionEnd(content, 'grid-template-rows')

    expect(callback).toHaveBeenCalledTimes(1)
  })

  it('при reduced motion вызывает callback на следующем кадре, не дожидаясь перехода', () => {
    mockedReducedMotion.mockReturnValue(true)
    const callback = jest.fn()

    afterExpanded(SectionId.Courses, callback)
    expect(callback).not.toHaveBeenCalled()

    jest.advanceTimersByTime(20)

    expect(callback).toHaveBeenCalledTimes(1)
  })

  it('если секции нет на странице, вызывает callback на следующем кадре', () => {
    document.body.innerHTML = ''
    const callback = jest.fn()

    afterExpanded(SectionId.Courses, callback)
    jest.advanceTimersByTime(20)

    expect(callback).toHaveBeenCalledTimes(1)
  })
})
