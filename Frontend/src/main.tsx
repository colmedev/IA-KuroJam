import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { createRoot } from 'react-dom/client'
import { ClerkProvider } from '@clerk/clerk-react'
import App from './App.tsx'
import './index.css'
import { ProtectedRoute } from './components/ProtectedRoute/ProtectedRoute.tsx'
import SignInPage from './pages/SignInPage.tsx'
import Chat from './pages/Chat.tsx'

const PUBLISHABLE_KEY = import.meta.env.VITE_PUBLISHED_KEY;

if (!PUBLISHABLE_KEY) {
  throw new Error('Missing Publishable Key')
}

createRoot(document.getElementById('root')!).render(
  <BrowserRouter>
    <ClerkProvider 
        publishableKey={PUBLISHABLE_KEY} 
        afterSignOutUrl="/" 
    >
      <Routes>
        <Route path="/" element={<App />} />
        <Route path='/app' element={<ProtectedRoute element={<Chat />}/>} />
        <Route path="/sign-in" element={<SignInPage />} />
      </Routes>
    </ClerkProvider>
  </BrowserRouter>,
)
