import './Button.css'
type ButtonProps = {
  text: string;
  link?: string;
}

export const Button = ({ text, link }: ButtonProps) => {
  return (
    <button className='button'>
      <a href={link}>{text}</a>
    </button>
  )
}