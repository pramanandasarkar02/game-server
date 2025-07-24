import { Axios } from 'axios';
import React, { useEffect } from 'react'
import type { PlayerInformation } from '../types/player';

type Props = {}

const Profile = (props: Props) => {
    const [player, setPlayer] = React.useState<PlayerInformation | null>(null);

    const extractProfileInformation = () => {
        // get profile information from backend
        const  playerId = localStorage.getItem('playerId');
        const token = localStorage.getItem('token');

        const response = Axios.get(`${baseURL}/player/info/${playerId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        })

        const playerInfo: PlayerInformation = response.data;
        console.log(playerInfo);
        setPlayer(playerInfo);

    }


    useEffect(() => {
        // document.title = "Profile"

    }, [])

    return (
    <div>
        <div>
            {/* left */}
            player profile
            <h2>{player?.userName}</h2>
            <h2>{player?.id}</h2>
            <h2>{player?.level}</h2>
            <h2>{player?.score}</h2>
            <h2>{player?.matchHistory}</h2>

        </div>
        <div>
            {/* right */}
        </div>


    </div>
  )
}

export default Profile