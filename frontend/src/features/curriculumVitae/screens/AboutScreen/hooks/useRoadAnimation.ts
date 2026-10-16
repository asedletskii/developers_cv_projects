import { useLayoutEffect, useRef } from 'react'
import { gsap } from '@shared/animations'

export const useRoadAnimation = (stepsCount: number) => {
  const stageRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    const stage = stageRef.current
    if (!stage || stepsCount < 2) return

    const ctx = gsap.context(() => {
      const path = stage.querySelector<SVGPathElement>('.road-path')
      const nodes = gsap.utils.toArray<SVGCircleElement>('.road-node', stage)
      const frames = gsap.utils.toArray<HTMLElement>('.road-frame', stage)
      if (!path || !frames.length) return

      const length = path.getTotalLength()
      gsap.set(path, { strokeDasharray: length, strokeDashoffset: length })
      gsap.set(nodes, { scale: 0.6, opacity: 0.3, transformOrigin: '50% 50%' })
      gsap.set(frames, { position: 'absolute', inset: 0, opacity: 0, y: 24 })
      gsap.set(frames[0], { opacity: 1, y: 0 })
      gsap.set(nodes[0], { scale: 1.3, opacity: 1 })

      const duration = frames.length - 1

      const tl = gsap.timeline({
        scrollTrigger: { trigger: stage, start: 'top top', end: 'bottom bottom', scrub: 0.6 },
      })

      tl.to(path, { strokeDashoffset: 0, ease: 'none', duration }, 0)

      frames.forEach((frame, index) => {
        if (index === 0) return
        const at = index - 1 + 0.4

        tl.to(frames[index - 1], { opacity: 0, y: -24, duration: 0.3 }, at)
        tl.to(frame, { opacity: 1, y: 0, duration: 0.3 }, at)
        tl.to(nodes[index - 1], { scale: 1, opacity: 0.4, duration: 0.3 }, at)
        tl.to(nodes[index], { scale: 1.3, opacity: 1, duration: 0.3 }, at)
      })
    }, stageRef)

    return () => ctx.revert()
  }, [stepsCount])

  return stageRef
}
