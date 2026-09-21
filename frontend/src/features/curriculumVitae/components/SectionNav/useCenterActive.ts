import { RefObject, useLayoutEffect } from 'react'

/** На узких экранах прокручивает саму навигацию так, чтобы активный пункт был по центру. */
export const useCenterActive = (
  pillRef: RefObject<HTMLElement | null>,
  activeId: string | null,
) => {
  useLayoutEffect(() => {
    const pill = pillRef.current
    const active = pill?.querySelector<HTMLElement>('[data-active="true"]')
    if (!pill || !active) return

    pill.scrollTo({
      left: active.offsetLeft - pill.clientWidth / 2 + active.clientWidth / 2,
      behavior: 'smooth',
    })
  }, [pillRef, activeId])
}
