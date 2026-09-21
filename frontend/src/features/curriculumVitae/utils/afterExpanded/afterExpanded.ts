import { SectionId } from '@features/curriculumVitae/types'
import { prefersReducedMotion } from '@shared/animations'

/** Вызывает callback, когда закончилась анимация раскрытия секции (или сразу, если анимации нет). */
export const afterExpanded = (id: SectionId, callback: () => void) => {
  const content = document.querySelector<HTMLElement>(`#${id} [data-collapsible]`)

  if (!content || prefersReducedMotion()) {
    requestAnimationFrame(callback)
    return
  }

  let isDone = false

  const finish = () => {
    if (isDone) return
    isDone = true
    content.removeEventListener('transitionend', onTransitionEnd)
    callback()
  }

  const onTransitionEnd = (event: TransitionEvent) => {
    if (event.propertyName === 'grid-template-rows') finish()
  }

  content.addEventListener('transitionend', onTransitionEnd)
  window.setTimeout(finish, 700)
}
