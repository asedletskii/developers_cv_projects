export const safeExternalUrl = (value?: string | null, allowRelative = false) => {
  if (!value) return undefined

  try {
    const url = allowRelative ? new URL(value, window.location.origin) : new URL(value)
    return url.protocol === 'http:' || url.protocol === 'https:' ? url.href : undefined
  } catch {
    return undefined
  }
}
