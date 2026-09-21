import { useLayoutEffect } from 'react'
import { ScrollSmoother } from '@shared/animations'

export const useSmoothScroll = () => {
  useLayoutEffect(() => {
    const smoother = ScrollSmoother.create({
      wrapper: '#smooth-wrapper',
      content: '#smooth-content',
      smooth: 1.2,
      effects: true,
    })

    return () => smoother.kill()
  }, [])
}
