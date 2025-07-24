import React from 'react'
import Header from '../components/Header'
import GameList from '../components/GameList'
import Footer from '../components/Footer'

const PlayerHome = () => {
  return (
    <div className="min-h-screen flex flex-col">
        <Header />
        <main className="flex-grow bg-gray-100">
            <GameList />
        </main>
        <Footer />
    </div>
  )
}

export default PlayerHome