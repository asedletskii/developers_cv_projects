/** Допускаем только #RGB и #RRGGBB: цвет приходит с бэка, и в тему должно попадать только валидное значение */
export const isHexColor = (value: unknown): value is string =>
  typeof value === 'string' && /^#([0-9a-f]{3}|[0-9a-f]{6})$/i.test(value)
