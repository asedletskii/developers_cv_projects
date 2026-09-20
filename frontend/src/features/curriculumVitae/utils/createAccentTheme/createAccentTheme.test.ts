import { createTheme } from '@mui/material/styles'
import { paletteDark } from '@shared/styles/paletteDark'
import { createAccentTheme } from './createAccentTheme'

const base = createTheme()

describe('createAccentTheme', () => {
  it('накладывает основной цвет и считает оттенки', () => {
    const theme = createAccentTheme(base, { primary: '#2E9E6B' })

    expect(theme).not.toBeNull()
    expect(theme?.palette.primary.main).toBe('#2E9E6B')
    expect(theme?.palette.primary.light).toBeTruthy()
    expect(theme?.palette.primary.dark).toBeTruthy()
    expect(theme?.palette.primary.contrastText).toBeTruthy()
  })

  it('не трогает второй акцент, если он не задан', () => {
    const theme = createAccentTheme(base, { primary: '#2E9E6B' })

    expect(theme?.palette.secondary.main).toBe(base.palette.secondary.main)
  })

  it('красит оба акцента', () => {
    const theme = createAccentTheme(base, { primary: '#2E9E6B', secondary: '#7C5CFA' })

    expect(theme?.palette.primary.main).toBe('#2E9E6B')
    expect(theme?.palette.secondary.main).toBe('#7C5CFA')
  })

  it('принимает короткую запись цвета', () => {
    expect(createAccentTheme(base, { primary: '#0af' })?.palette.primary.main).toBe('#0af')
  })

  it('невалидный цвет игнорирует, валидный применяет', () => {
    const theme = createAccentTheme(base, { primary: 'red', secondary: '#7C5CFA' })

    expect(theme?.palette.primary.main).toBe(base.palette.primary.main)
    expect(theme?.palette.secondary.main).toBe('#7C5CFA')
  })

  it.each([undefined, {}, { primary: 'red' }, { primary: 'javascript:alert(1)' }])(
    'если применять нечего (%p), возвращает null',
    (colors) => {
      expect(createAccentTheme(base, colors)).toBeNull()
    },
  )

  it('сохраняет светлую или тёмную схему базовой темы', () => {
    const dark = createTheme({ palette: paletteDark })

    expect(createAccentTheme(dark, { primary: '#2E9E6B' })?.palette.mode).toBe('dark')
  })

  it('подбирает контрастный цвет текста: тёмный на светлом фоне и светлый на тёмном', () => {
    const onLight = createAccentTheme(base, { primary: '#FFEB3B' })
    const onDark = createAccentTheme(base, { primary: '#1A237E' })

    expect(onLight?.palette.primary.contrastText).toBe('rgba(0, 0, 0, 0.87)')
    expect(onDark?.palette.primary.contrastText).toBe('#fff')
  })

  it('не меняет базовую тему', () => {
    const before = base.palette.primary.main

    createAccentTheme(base, { primary: '#2E9E6B' })

    expect(base.palette.primary.main).toBe(before)
  })
})
