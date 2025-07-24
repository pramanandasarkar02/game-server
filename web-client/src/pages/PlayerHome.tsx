import React from 'react'
import Header from '../components/Header'
import GameList from '../components/GameList'
import Footer from '../components/Footer'

type Props = {}

const PlayerHome = (props: Props) => {
  return (
    <div>
        <div>
            {/* header */}
            <Header />
        </div>
        <div>
            {/* Games */}
            <GameList />
        </div>
        <div>
            {/* footer  */}
            <Footer />
        </div>
    </div>
  )
}

export default PlayerHome