import React from 'react';
import './ProtectedNavbar.css';
import { SignedOut, SignedIn, UserButton } from '@clerk/clerk-react';
import { Signin } from '../auth/Signin';

const ProtectedNavbar: React.FC = () => {
  return (
    <nav className="navbar">
      <div className="navbar-brand">
        <a href="/">CareerCraft</a>
      </div>
      <ul className="navbar-links">
        <li><a href="/app">Find your career</a></li>
        <li><a href="/results">Test Results</a></li>
      </ul>
      <div className="navbar-buttons">
        <SignedOut>
          <Signin />
        </SignedOut>
        <SignedIn>
          <UserButton />
        </SignedIn>
      </div>
    </nav>
  );
};

export default ProtectedNavbar;
