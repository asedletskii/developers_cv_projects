import { useCallback, useEffect, useRef, useState } from 'react'
import { SectionId } from '@features/curriculumVitae/types'

const SECTION_IDS = new Set<string>(Object.values(SectionId))
const HOST_ATTRIBUTE = 'data-click-map-host'

export type EmbeddedSection = { id: string; host: HTMLElement }

/** Находит секции страницы резюме внутри iframe (тот же origin) и вешает в каждую пустой слой,
 *  куда рисуется тепловая карта. Слой — потомок секции, поэтому прокручивается вместе с ней, а
 *  раскладка на нужной ширине считается самим приложением — предпросмотр настоящий, не макет. */
export const useEmbeddedSections = () => {
  const iframeRef = useRef<HTMLIFrameElement>(null)
  const observerRef = useRef<MutationObserver | null>(null)
  const [sections, setSections] = useState<EmbeddedSection[]>([])

  const sync = useCallback(() => {
    const doc = iframeRef.current?.contentDocument
    if (!doc?.body) return

    const found = [...doc.querySelectorAll<HTMLElement>('section[id]')]
      .filter((element) => SECTION_IDS.has(element.id))
      .map((element) => {
        let host = element.querySelector<HTMLElement>(`:scope > [${HOST_ATTRIBUTE}]`)

        if (!host) {
          host = doc.createElement('div')
          host.setAttribute(HOST_ATTRIBUTE, '')
          Object.assign(host.style, {
            position: 'absolute',
            inset: '0',
            pointerEvents: 'none',
            zIndex: '5',
          })
          if (doc.defaultView?.getComputedStyle(element).position === 'static') {
            element.style.position = 'relative'
          }
          element.appendChild(host)
        }

        return { id: element.id, host }
      })

    setSections((current) =>
      current.length === found.length && current.every((item, i) => item.host === found[i].host)
        ? current
        : found,
    )
  }, [])

  const onLoad = useCallback(() => {
    sync()

    observerRef.current?.disconnect()
    const body = iframeRef.current?.contentDocument?.body
    if (!body) return

    // Секции появляются не сразу: страница сначала грузит резюме, потом рисует его
    observerRef.current = new MutationObserver(() => requestAnimationFrame(sync))
    observerRef.current.observe(body, { childList: true, subtree: true })
  }, [sync])

  useEffect(() => () => observerRef.current?.disconnect(), [])

  return { iframeRef, onLoad, sections }
}
