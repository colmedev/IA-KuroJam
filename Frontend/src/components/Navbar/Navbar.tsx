import './Navbar.css'
import { Button }  from '../Button/Button'
export function Navbar() {
  
  const links = [
    {
      link: "#",
      text: "Home"
    }
  ]

  return (
    <>
    <ul className="Navbar">
      {
        links.map((link, index) => (
          <li key={index} className="Navbar__item">
            <a href={link.link}>{link.text}</a>
          </li>
        ))
      }
      <Button text="Get Started" link="/app"/>
    </ul>
    </>
  )
}

