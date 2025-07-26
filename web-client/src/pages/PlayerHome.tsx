import React from 'react'
import Header from '../components/Header'
import GameList from '../components/GameList'
import Footer from '../components/Footer'

const PlayerHome = () => {
  return (
    <div className="min-h-screen flex flex-col">
        
        <main className="flex-grow bg-gray-100">
            <GameList />
        </main>
        
    </div>
  )
}

export default PlayerHome