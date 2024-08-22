import './Card.css'
interface CardProps {
  icon: string;
  title: string;
  description: string;
}

export const Card = ({icon, title, description} : CardProps) => {

  

  return (
    <div className="card__container">
      <img src={`/svg/${icon}.svg`} alt="icon" className="card__icon"/>
      <h3 className="card__title">{title}</h3>
      <p className="card__content">{description}</p>
    </div>
  )
}