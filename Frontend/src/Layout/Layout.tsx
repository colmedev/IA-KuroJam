import { Button, Navbar } from '../components';
import './Layout.css'
//TODO Improve the layout component adding pictures or content 
const Layout = () => {
  return (
    <>
    <div className='layout__container'>
      <Navbar />
      <div className="layout__content">
        <h1 className="layout__title">Discover your ideal career with AI.</h1>

        <p className="layout__description">Unlock your full potential and find the perfect career path with our AI-powered guidance.</p>

        <Button 
        text="Get Started" 
        link="/app" />
      </div>
    </div>
    </>
  );
};

export default Layout;