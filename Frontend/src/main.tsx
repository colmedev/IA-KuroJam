import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import { Chat } from './pages/Chat.tsx'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import './index.css'

createRoot(document.getElementById('root')!).render(

  <BrowserRouter>
    <Routes>
        <Route path="/" element={<App />} />
        <Route path='/app' element={<Chat />} />
    </Routes>
  </BrowserRouter>,
)
