export const generateId = (() => {
  let id = 0

  return () => {
    id++
    return id
  }
})()
