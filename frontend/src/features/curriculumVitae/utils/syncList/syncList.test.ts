import { syncById, syncByIndex } from './syncList'

const makeOps = () => {
  const calls: string[] = []

  return {
    calls,
    ops: {
      create: async (item: { id: string }) => void calls.push(`create ${item.id}`),
      update: async (id: string) => void calls.push(`update ${id}`),
      remove: async (id: string) => void calls.push(`remove ${id}`),
    },
  }
}

describe('syncById', () => {
  it('удаляет пропавшие, обновляет известные, создаёт новые', async () => {
    const { calls, ops } = makeOps()

    await syncById([{ id: 'a' }, { id: 'b' }], [{ id: 'b' }, { id: 'new' }], ops)

    expect(calls).toEqual(['remove a', 'update b', 'create new'])
  })
})

describe('syncByIndex', () => {
  it('i-я запись формы обновляет i-ю серверную, лишние создаёт и удаляет с хвоста', async () => {
    const created = makeOps()
    await syncByIndex([{ id: 's0' }], [{ id: 'x' }, { id: 'y' }], created.ops)
    expect(created.calls).toEqual(['update s0', 'create y'])

    const removed = makeOps()
    await syncByIndex([{ id: 's0' }, { id: 's1' }], [{ id: 'x' }], removed.ops)
    expect(removed.calls).toEqual(['update s0', 'remove s1'])
  })
})
