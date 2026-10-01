import { gsap } from './gsap.config'

type SpringOptions = {
  damping?: number
  response?: number
  precision?: number
}

export type Spring = {
  to(target: number): void
  set(value: number): void
  kill(): void
}

export const createSpring = (
  initial: number,
  onUpdate: (value: number) => void,
  { damping = 1, response = 0.4, precision = 0.01 }: SpringOptions = {},
): Spring => {
  const stiffness = ((2 * Math.PI) / response) ** 2
  const drag = (4 * Math.PI * damping) / response

  let value = initial
  let target = initial
  let velocity = 0
  let running = false

  const tick = (_time: number, deltaMs: number) => {
    let remaining = Math.min(deltaMs / 1000, 0.05)

    while (remaining > 0) {
      const dt = Math.min(remaining, 1 / 120)
      velocity += (stiffness * (target - value) - drag * velocity) * dt
      value += velocity * dt
      remaining -= dt
    }

    const isSettled = Math.abs(target - value) < precision && Math.abs(velocity) < precision

    if (isSettled) {
      value = target
      velocity = 0
      stop()
    }

    onUpdate(value)
  }

  const start = () => {
    if (running) return
    running = true
    gsap.ticker.add(tick)
  }

  const stop = () => {
    running = false
    gsap.ticker.remove(tick)
  }

  return {
    to(next) {
      target = next
      start()
    },
    set(next) {
      value = target = next
      velocity = 0
      onUpdate(next)
    },
    kill: stop,
  }
}
