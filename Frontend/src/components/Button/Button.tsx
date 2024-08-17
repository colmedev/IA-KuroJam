import './Button.css'
type ButtonProps = {
  text: string;
  link?: string;
  isInverted?: boolean;
}

export const Button = ({ text, link, isInverted}: ButtonProps) => {
  return (
    <button className={`button ${isInverted ? 'button--inverted' : ''}`}>
      <a href={link}>{text}</a>
    </button>
  )
}