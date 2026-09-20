type Props = {
  name: string
  surname: string
  patronymic?: string
}

export const getPersonFIO = ({ ...props }: Props) => {
  let fio = props.surname + ` ${props.name}`

  if (props.patronymic) {
    fio += ` ${props.patronymic}`
  }

  return fio
}
