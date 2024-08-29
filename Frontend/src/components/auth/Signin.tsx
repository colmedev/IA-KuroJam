import './Signin.css'
import { SignedOut } from "@clerk/clerk-react"

export function Signin() {


  const handleClick = () => {
    const dialog = document.querySelector('dialog')
    dialog?.showModal()
  }


  const handleClose = (e: React.MouseEvent) => {
    e.stopPropagation()
    const dialog = document.querySelector('dialog')
    dialog?.close()
  }

  return(
    <>
    <div className='signin'>
      <button className="button-signin" onClick={handleClick}>
        Sign in
      </button>
      <dialog onClick={handleClose} className="signin__dialog">
      </dialog>
    <SignedOut />
    </div>
    </>
  )
}