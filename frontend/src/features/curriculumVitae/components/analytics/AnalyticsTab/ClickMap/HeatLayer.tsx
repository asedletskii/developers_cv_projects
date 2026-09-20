import { ClickGrid } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'
import { getHeatColor } from './heatColor'

type Props = {
  title: string
  grid: ClickGrid
  max: number
  scale: number
}

/** Тепловая сетка поверх одной секции. Рисуется в документе iframe, где нет стилей MUI/emotion,
 *  поэтому только инлайн-стили. Берём клетки с кликами (до 400 на секцию) */
export const HeatLayer = ({ title, grid, max, scale }: Props) => {
  const translate = useTranslate('Lk.Analytics.ClickMap')
  const size = `${100 / grid.grid}%`

  return (
    <>
      {grid.cells.map(({ x, y, count }) => (
        <div
          key={`${x}-${y}`}
          title={String(count)}
          style={{
            position: 'absolute',
            left: `${(x / grid.grid) * 100}%`,
            top: `${(y / grid.grid) * 100}%`,
            width: size,
            height: size,
            background: getHeatColor(count / max),
          }}
        />
      ))}

      <span
        style={{
          position: 'absolute',
          top: 6 / scale,
          right: 6 / scale,
          padding: '2px 8px',
          borderRadius: 4,
          font: `${Math.round(13 / scale)}px/1.4 sans-serif`,
          color: '#fff',
          background: 'rgba(0, 0, 0, 0.65)',
        }}
      >
        {title} · {translate('clicks', { count: grid.total })}
      </span>
    </>
  )
}
