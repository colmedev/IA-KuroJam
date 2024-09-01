import { SignIn } from '@clerk/clerk-react';


export default function SignInPage() {
  return(
        <div className='center-container'>
            <SignIn routing="path" path="/sign-in" fallbackRedirectUrl={"/app"}/>
        </div>
  )
}
