import { useNavigate } from 'react-router'
import './Signin.css'

export function Signin() {

  const navigate = useNavigate();

  const handleClick = () => {
    navigate('/sign-in');
  }
 

  return(
    <>
    <div className='signin'>
      <button className="button-signin" onClick={handleClick}>
        Sign in
      </button>
    </div>
    </>
  )
}
