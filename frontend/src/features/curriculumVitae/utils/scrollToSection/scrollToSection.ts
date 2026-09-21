import { SectionId } from '@features/curriculumVitae/types'
import { gsap } from '@shared/animations'

/** Отступ сверху, чтобы плавающая навигация не закрывала заголовок секции */
const SCROLL_OFFSET = 72

/** Плавно прокручивает страницу к секции с учётом отступа под навигацию. */
export const scrollToSection = (id: SectionId) => {
  const section = document.getElementById(id)
  if (!section) return

  // Ещё не проявившаяся секция сдвинута анимацией появления — цель считаем без этого сдвига
  const revealOffset = Number(gsap.getProperty(section, 'y')) || 0
  const top = section.getBoundingClientRect().top + window.scrollY - revealOffset - SCROLL_OFFSET

  window.scrollTo({ top, behavior: 'smooth' })
}
