import { useState } from 'react'
import './App.css'

import MainPage from './components/MainPage'
import LandingPage from './components/LandingPage'
import { Link } from 'react-router-dom'
import {Routes , Route } from "react-router-dom" 

function App() {

  return (
    <>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/app" element={<MainPage />} />
      </Routes>
    </>
  )
}

export default App
