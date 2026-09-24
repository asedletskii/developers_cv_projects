import { parseLkError } from './parseLkError'

const axiosError = (status: number, data?: unknown) => ({
  isAxiosError: true,
  response: { status, data },
})

describe('parseLkError', () => {
  it('404 становится notFound', () => {
    expect(parseLkError(axiosError(404)).kind).toBe('notFound')
  })

  it('любой другой статус ответа становится cannotGetLk', () => {
    expect(parseLkError(axiosError(500)).kind).toBe('cannotGetLk')
    expect(parseLkError(axiosError(400)).kind).toBe('cannotGetLk')
  })

  it('сохраняет HTTP-статус ответа', () => {
    expect(parseLkError(axiosError(500)).status).toBe(500)
  })

  it('достаёт requestId из тела ошибки бэкенда', () => {
    const error = axiosError(500, { error: { code: 'INTERNAL_ERROR', requestId: 'req-123' } })

    expect(parseLkError(error).requestId).toBe('req-123')
  })

  it('без requestId в теле — requestId не заполнен', () => {
    expect(parseLkError(axiosError(500, {})).requestId).toBeUndefined()
  })

  it('сетевая ошибка без ответа (нет соединения) — cannotGetLk без статуса', () => {
    const result = parseLkError({ isAxiosError: true, response: undefined })

    expect(result).toEqual({ kind: 'cannotGetLk' })
  })

  it('не-axios ошибка тоже сводится к cannotGetLk без статуса', () => {
    expect(parseLkError(new Error('boom'))).toEqual({ kind: 'cannotGetLk' })
  })

  it('совсем не Error (например, строка) не падает', () => {
    expect(parseLkError('oops')).toEqual({ kind: 'cannotGetLk' })
  })
})
