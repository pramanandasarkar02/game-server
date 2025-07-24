import React, { useContext, useEffect, useState } from 'react'
import axios from 'axios'
import { FaTrophy, FaHistory, FaUserAlt, FaIdCard } from 'react-icons/fa'
import type { PlayerInformation } from '../types/player'
import { useNavigate } from 'react-router-dom'
import { AuthContext } from '../contexts/AuthContext'

const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:4000/api'

const Profile = () => {
    const { player} = useContext(AuthContext)
    const [playerInfo, setPlayerInfo] = useState<PlayerInformation | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')
    const navigate = useNavigate()

    useEffect(() => {
        const extractProfileInformation = async () => {
            try {
                console.log("printing player ....");
                // console.log(player);

                //  get player id and token from local storage
                const playerId = localStorage.getItem('playerId')
                const token = localStorage.getItem('token')

                console.log(`playerId: ${playerId}, token: ${token}`);

                const response = await axios.get(`${baseURL}/player/info/${playerId}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                })

                // setPlayer(response.data)
                setPlayerInfo(response.data.player)

                // print player info
                console.log(response.data)
            } catch (err: unknown) {
                if (err instanceof Error) {
                    setError(err.message)
                } else {
                    setError('Failed to fetch profile information')
                }
                
                // Handle 401 unauthorized
                if (axios.isAxiosError(err) && err.response?.status === 401) {
                    localStorage.removeItem('token')
                    localStorage.removeItem('playerId')
                    navigate('/register')
                }
            } finally {
                setLoading(false)
            }
        }

        extractProfileInformation()
    }, [navigate])

    if (loading) {
        return <div className="text-center py-8">Loading...</div>
    }

    if (error) {
        return <div className="text-center py-8 text-red-500">{error}</div>
    }

    if (!playerInfo) {
        return <div className="text-center py-8 text-red-500">Player data not available</div>
    }

    return (
        <div className="container mx-auto py-8">
            <div className="max-w-4xl mx-auto bg-white rounded-lg shadow-md overflow-hidden">
                <div className="md:flex">
                    {/* Left Side - Profile Info */}
                    <div className="md:w-1/2 p-6 bg-gray-100">
                        <div className="flex items-center mb-6">
                            <div className="bg-blue-500 text-white p-3 rounded-full mr-4">
                                <FaUserAlt size={24} />
                            </div>
                            <h2 className="text-2xl font-bold">{playerInfo?.userName}</h2>
                        </div>

                        <div className="space-y-4">
                            <div className="flex items-center">
                                <FaIdCard className="text-gray-500 mr-2" />
                                <span className="font-medium">ID:</span>
                                <span className="ml-2">{playerInfo?.id}</span>
                            </div>

                            <div className="flex items-center">
                                <FaTrophy className="text-yellow-500 mr-2" />
                                <span className="font-medium">Level:</span>
                                <span className="ml-2">{playerInfo?.level}</span>
                            </div>

                            <div className="flex items-center">
                                <FaTrophy className="text-blue-500 mr-2" />
                                <span className="font-medium">Score:</span>
                                <span className="ml-2">{playerInfo?.score}</span>
                            </div>
                        </div>
                    </div>

                    {/* Right Side - Match History */}
                    <div className="md:w-1/2 p-6">
                        <h3 className="text-xl font-bold mb-4 flex items-center">
                            <FaHistory className="mr-2" />
                            Match History
                        </h3>
                        {playerInfo?.matchHistory && playerInfo.matchHistory.length > 0 ? (
                            <ul className="divide-y divide-gray-200">
                                {playerInfo.matchHistory.map((match, index) => (
                                    <li key={index} className="py-2">
                                        {match}
                                    </li>
                                ))}
                            </ul>
                        ) : (
                            <p className="text-gray-500">No match history available</p>
                        )}
                    </div>
                </div>

                <div className='text-right'>
                    <button className="position-absolute m-2  bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition" onClick={() => navigate('/home')}>
                        back to home
                    </button>
                </div>
            </div>
        </div>
    )
}

export default Profile