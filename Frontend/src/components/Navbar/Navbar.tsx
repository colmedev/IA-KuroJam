import './Navbar.css'
import { Button }  from '../Button/Button'
export function Navbar() {
  

  return (
    <>
    <ul className="Navbar">
      <li className='Navbar__item'>
      <Button text="Login" link="#"/>
      </li>

      <li className='Navbar__item'>
      <Button text="Sign Up" link="#"/>
      </li>
    </ul>
    </>
  )
}

