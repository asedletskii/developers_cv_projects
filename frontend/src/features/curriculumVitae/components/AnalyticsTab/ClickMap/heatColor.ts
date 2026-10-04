/** Цвет клетки по доле от максимума (0..1): зелёный → жёлтый → оранжевый → красный.
 *  Идём по оттенку HSL от 120° до 0° — промежуточные цвета получаются сами. */
export const getHeatColor = (ratio: number, alpha = 0.6) => {
  const clamped = Math.min(Math.max(ratio, 0), 1)
  return `hsla(${Math.round(120 * (1 - clamped))}, 90%, 48%, ${alpha})`
}

/** Градиент для легенды — та же шкала, что и у клеток */
export const HEAT_LEGEND_GRADIENT = `linear-gradient(90deg, ${[0, 0.33, 0.66, 1]
  .map((ratio) => getHeatColor(ratio, 1))
  .join(', ')})`
