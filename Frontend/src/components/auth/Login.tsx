import { useSignUp } from "@clerk/clerk-react"

//TODO: Add functional login
export const Login = () => {
  const { setActive, signUp, isLoaded } = useSignUp()

  const HandleLogin = () => {
    
  } 

  return (
    <button className="button-login" onClick={HandleLogin}>
      
    </button>
  )
}