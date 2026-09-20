import { useEffect, useState } from 'react'

const DESKTOP_QUERY = '(min-width: 900px)'
const REDUCE_MOTION_QUERY = '(prefers-reduced-motion: reduce)'

export const useIsDesktopStory = () => {
  const [isDesktopStory, setIsDesktopStory] = useState(
    () =>
      window.matchMedia(DESKTOP_QUERY).matches && !window.matchMedia(REDUCE_MOTION_QUERY).matches,
  )

  useEffect(() => {
    const desktopQuery = window.matchMedia(DESKTOP_QUERY)
    const reduceMotionQuery = window.matchMedia(REDUCE_MOTION_QUERY)

    const update = () => setIsDesktopStory(desktopQuery.matches && !reduceMotionQuery.matches)

    desktopQuery.addEventListener('change', update)
    reduceMotionQuery.addEventListener('change', update)
    return () => {
      desktopQuery.removeEventListener('change', update)
      reduceMotionQuery.removeEventListener('change', update)
    }
  }, [])

  return isDesktopStory
}
