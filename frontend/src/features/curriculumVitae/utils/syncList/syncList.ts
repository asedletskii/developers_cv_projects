export type ListOps<F> = {
  create(item: F, index: number): Promise<unknown>
  update(id: string, item: F, index: number): Promise<void>
  remove(id: string): Promise<void>
}

/** Записи с собственным id (у новых он клиентский — бэку неизвестен): лишние на сервере
 *  удаляем, известные обновляем, остальные создаём в порядке формы. */
export const syncById = async <F extends { id: string }>(
  server: { id: string }[],
  form: F[],
  ops: ListOps<F>,
) => {
  const serverIds = new Set(server.map(({ id }) => id))
  const formIds = new Set(form.map(({ id }) => id))

  for (const { id } of server) if (!formIds.has(id)) await ops.remove(id)

  for (const [index, item] of form.entries()) {
    if (serverIds.has(item.id)) await ops.update(item.id, item, index)
    else await ops.create(item, index)
  }
}

/** Записи без id в форме (языки, технологии; порядок хранит сервер, новые дописываются в конец): i-я запись формы
 *  становится i-й записью сервера, лишние создаются/удаляются с хвоста. */
export const syncByIndex = async <F>(server: { id: string }[], form: F[], ops: ListOps<F>) => {
  for (const [index, item] of form.entries()) {
    const existing = server[index]
    if (existing) await ops.update(existing.id, item, index)
    else await ops.create(item, index)
  }

  for (const { id } of server.slice(form.length)) await ops.remove(id)
}
