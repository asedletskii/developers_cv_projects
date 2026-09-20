import { safeExternalUrl } from './safeExternalUrl'

describe('safeExternalUrl', () => {
  it('пускает только http(s)', () => {
    expect(safeExternalUrl('https://github.com/a')).toBe('https://github.com/a')
    expect(safeExternalUrl('http://example.ru')).toBe('http://example.ru/')
  })

  it('отбрасывает опасные и неразборчивые адреса', () => {
    expect(safeExternalUrl('javascript:alert(1)')).toBeUndefined()
    expect(safeExternalUrl('JaVaScRiPt:alert(1)')).toBeUndefined()
    expect(safeExternalUrl('data:text/html,<script>1</script>')).toBeUndefined()
    expect(safeExternalUrl('//evil.example')).toBeUndefined()
    expect(safeExternalUrl('not a url')).toBeUndefined()
    expect(safeExternalUrl('')).toBeUndefined()
    expect(safeExternalUrl(undefined)).toBeUndefined()
  })

  it('путь без хоста — только когда разрешено, и он ведёт на этот же сайт', () => {
    expect(safeExternalUrl('/uploads/a.png')).toBeUndefined()
    expect(safeExternalUrl('/uploads/a.png', true)).toBe(`${window.location.origin}/uploads/a.png`)
    expect(safeExternalUrl('javascript:alert(1)', true)).toBeUndefined()
  })
})
