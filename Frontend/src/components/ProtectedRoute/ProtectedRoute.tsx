import { Navigate } from 'react-router-dom';
import { useAuth } from '@clerk/clerk-react';
import { Footer } from '../Footer/Footer';
import ProtectedNavbar from '../ProtectedNavbar/ProtectedNavbar';

interface ProtectedRouteProps {
    element: React.ReactElement
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ element }) => {
  const {isLoaded, userId} = useAuth();
 
  if (!isLoaded) {
        return
  }

  if (!userId) {
    return <Navigate to="/sign-in" />;
  }

  return (
    <div className='container'>
        <ProtectedNavbar/> 
        {element}
        <Footer/>
    </div>
  );
}
