import { useEffect, useState } from 'react'

export const useElementWidth = (ref: React.RefObject<HTMLElement>) => {
  const [width, setWidth] = useState(0)

  useEffect(() => {
    if (!ref.current) return

    const observer = new ResizeObserver(([entry]) => {
      setWidth(entry.contentRect.width)
    })

    observer.observe(ref.current)

    return () => observer.disconnect()
  }, [ref])

  return width
}
