/** Детерминированный ГПСЧ (mulberry32): рендер остаётся чистым, а раскладка выглядит случайной. */
export const createRandom = (seed: number) => {
  let state = seed

  return () => {
    state += 0x6d2b79f5
    let value = Math.imul(state ^ (state >>> 15), 1 | state)
    value ^= value + Math.imul(value ^ (value >>> 7), 61 | value)

    return ((value ^ (value >>> 14)) >>> 0) / 4294967296
  }
}
