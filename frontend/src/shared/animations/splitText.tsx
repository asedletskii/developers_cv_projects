export const splitText = (text: string, className: string) => {
  return text.split('').map((char, index) => (
    <span
      key={index}
      style={{
        display: 'inline-block',
      }}
      className={className}
    >
      {char}
    </span>
  ))
}
