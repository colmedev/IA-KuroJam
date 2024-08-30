import {SignedIn, SignedOut, UserButton } from '@clerk/clerk-react'
import './Navbar.css'
import { Signin } from '../auth/Signin';

export function Navbar() { 

  return (
    <>
        <div className='Navbar'>
            <SignedOut>
                <Signin />
            </SignedOut>
            <SignedIn>
                <UserButton />
            </SignedIn>
        </div>
    </>
  )
}
