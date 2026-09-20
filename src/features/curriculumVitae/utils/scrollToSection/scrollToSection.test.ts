import { SectionId } from '@features/curriculumVitae/types'
import { gsap } from '@shared/animations'
import { scrollToSection } from './scrollToSection'

jest.mock('@shared/animations', () => ({ gsap: { getProperty: jest.fn() } }))

const mockedGetProperty = gsap.getProperty as jest.Mock

const SCROLL_OFFSET = 72

const mountSection = (id: string, top: number) => {
  document.body.innerHTML = `<section id="${id}"></section>`
  const section = document.getElementById(id) as HTMLElement
  jest.spyOn(section, 'getBoundingClientRect').mockReturnValue({ top } as DOMRect)
  return section
}

describe('scrollToSection', () => {
  beforeEach(() => {
    window.scrollTo = jest.fn()
    Object.defineProperty(window, 'scrollY', { value: 1000, configurable: true })
    mockedGetProperty.mockReturnValue(0)
  })

  it('плавно прокручивает так, чтобы секция встала под навигацией', () => {
    mountSection(SectionId.About, 500)

    scrollToSection(SectionId.About)

    expect(window.scrollTo).toHaveBeenCalledWith({
      top: 500 + 1000 - SCROLL_OFFSET,
      behavior: 'smooth',
    })
  })

  it('не учитывает сдвиг секции от анимации появления', () => {
    // Секция ещё не проявилась и смещена на 56px вниз
    mockedGetProperty.mockReturnValue(56)
    mountSection(SectionId.Projects, 500)

    scrollToSection(SectionId.Projects)

    expect(window.scrollTo).toHaveBeenCalledWith({
      top: 500 + 1000 - 56 - SCROLL_OFFSET,
      behavior: 'smooth',
    })
  })

  it('если GSAP не вернул число, считает сдвиг нулевым', () => {
    mockedGetProperty.mockReturnValue(undefined)
    mountSection(SectionId.About, 200)

    scrollToSection(SectionId.About)

    expect(window.scrollTo).toHaveBeenCalledWith({
      top: 200 + 1000 - SCROLL_OFFSET,
      behavior: 'smooth',
    })
  })

  it('если секции нет на странице, ничего не делает', () => {
    document.body.innerHTML = ''

    scrollToSection(SectionId.About)

    expect(window.scrollTo).not.toHaveBeenCalled()
  })
})
