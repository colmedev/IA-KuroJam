import { Navigate } from 'react-router-dom';
import { useAuth } from '@clerk/clerk-react';

export function ProtectedRoute({ element }) {
  const {isLoaded, userId} = useAuth();
 
    console.log(isLoaded, userId);
  if (!isLoaded) {
        return
  }

  if (!userId) {
    return <Navigate to="/sign-in" />;
  }

  return element;
}
