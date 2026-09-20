export const MIN_VISIBLE_RATIO = 0.4
export const MIN_VIEWPORT_SHARE = 0.5

export const isSectionVisibleEnough = (
  intersectionRatio: number,
  visibleHeight: number,
  viewportHeight: number,
) => intersectionRatio >= MIN_VISIBLE_RATIO || visibleHeight >= viewportHeight * MIN_VIEWPORT_SHARE

type SectionRect = { id: string; top: number; bottom: number }

export const pickMostVisibleSection = (rects: SectionRect[], viewportHeight: number) => {
  let best: { id: string; visible: number } | undefined

  for (const { id, top, bottom } of rects) {
    const visible = Math.min(bottom, viewportHeight) - Math.max(top, 0)

    if (visible > 0 && (!best || visible > best.visible)) best = { id, visible }
  }

  return best?.id
}
