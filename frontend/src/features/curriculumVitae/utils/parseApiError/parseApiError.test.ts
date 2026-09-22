import { parseApiError } from './parseApiError'

const axiosError = (status: number, data?: unknown) => ({
  isAxiosError: true,
  response: { status, data },
})

describe('parseApiError', () => {
  it('404 становится notFound', () => {
    expect(parseApiError(axiosError(404)).kind).toBe('notFound')
  })

  it('любой другой статус ответа становится cannotGetPersonData', () => {
    expect(parseApiError(axiosError(500)).kind).toBe('cannotGetPersonData')
    expect(parseApiError(axiosError(400)).kind).toBe('cannotGetPersonData')
  })

  it('сохраняет HTTP-статус ответа', () => {
    expect(parseApiError(axiosError(500)).status).toBe(500)
  })

  it('достаёт requestId из тела ошибки бэкенда', () => {
    const error = axiosError(500, { error: { code: 'INTERNAL_ERROR', requestId: 'req-123' } })

    expect(parseApiError(error).requestId).toBe('req-123')
  })

  it('без requestId в теле — requestId не заполнен', () => {
    expect(parseApiError(axiosError(500, {})).requestId).toBeUndefined()
  })

  it('сетевая ошибка без ответа (нет соединения) — cannotGetPersonData без статуса', () => {
    const result = parseApiError({ isAxiosError: true, response: undefined })

    expect(result).toEqual({ kind: 'cannotGetPersonData' })
  })

  it('не-axios ошибка тоже сводится к cannotGetPersonData без статуса', () => {
    expect(parseApiError(new Error('boom'))).toEqual({ kind: 'cannotGetPersonData' })
  })

  it('совсем не Error (например, строка) не падает', () => {
    expect(parseApiError('oops')).toEqual({ kind: 'cannotGetPersonData' })
  })
})
