import Layout from './Layout/Layout'
import { Card, Cta, Footer } from './components'

function App() {

  return (
    <>
    <div className='container'>
    <Layout />

    <div className="career__path">
      <h2>Explore Different Career Paths</h2>
      <p>Find the perfect career that aligns with your passions and strengths.</p>
      <div className='career__cards'>
        <Card 
        icon='briefCase'
        title='Business and Finance'
        description='Explore careers in accounting, marketing, management, and more.'
        />
        <Card 
        icon='Code'
        title='Business and Finance'
        description='Explore careers in accounting, marketing, management, and more.'
        />
        <Card 
        icon='Heart'
        title='Business and Finance'
        description='Explore careers in accounting, marketing, management, and more.'
        />
        <Card 
        icon='Palette'
        title='Business and Finance'
        description='Explore careers in accounting, marketing, management, and more.'
        />
        <Card 
        icon='School'
        title='Business and Finance'
        description='Explore careers in accounting, marketing, management, and more.'
        />
        <Card 
        icon='rrhh'
        title='Business and Finance'
        description='Explore careers in accounting, marketing, management, and more.'
        />
        
      </div>
    </div>

    <Cta />

    <Footer />
    </div>

    </>
  )
}

export default App
