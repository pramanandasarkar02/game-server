import React, { useEffect, useState } from 'react'
import axios from 'axios'
import { 
  FaTrophy, 
  FaHistory, 
  FaUserAlt, 
  FaIdCard,
  FaCrown,
  FaChartLine,
  FaSpinner
} from 'react-icons/fa'
import { useNavigate } from 'react-router-dom'
import type { PlayerInformation } from '../types/player'

const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:4000/api'

const Profile = () => {
    const [player, setPlayer] = useState<PlayerInformation | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')
    const navigate = useNavigate()

    useEffect(() => {
        const fetchProfileInformation = async () => {
            try {
                setLoading(true)
                setError('')
                
                const playerId = localStorage.getItem('playerId')
                const token = localStorage.getItem('token')

                if (!playerId || !token) {
                    throw new Error('Authentication required')
                }

                const response = await axios.get(`${baseURL}/player/info/${playerId}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                })

                if (response.data) {
                    setPlayer(response.data)
                } else {
                    throw new Error('No player data received')
                }
            } catch (err) {
                const errorMessage = err instanceof Error ? err.message : 'Failed to fetch profile information'
                setError(errorMessage)
                
                // Redirect to login if unauthorized
                if (axios.isAxiosError(err) && err.response?.status === 401) {
                    localStorage.removeItem('token')
                    localStorage.removeItem('playerId')
                    navigate('/login')
                }
            } finally {
                setLoading(false)
            }
        }

        fetchProfileInformation()
    }, [navigate])

    if (loading) {
        return (
            <div className="flex flex-col items-center justify-center min-h-[60vh]">
                <FaSpinner className="animate-spin text-blue-500 text-4xl mb-4" />
                <p className="text-gray-600">Loading your profile...</p>
            </div>
        )
    }

    if (error) {
        return (
            <div className="container mx-auto py-8">
                <div className="max-w-4xl mx-auto bg-red-50 border border-red-200 rounded-lg p-6">
                    <h2 className="text-xl font-bold text-red-600 mb-2">Error Loading Profile</h2>
                    <p className="text-red-500 mb-4">{error}</p>
                    <button 
                        onClick={() => window.location.reload()}
                        className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700 transition"
                    >
                        Try Again
                    </button>
                </div>
            </div>
        )
    }

    return (
        <div className="container mx-auto py-8 px-4">
            <div className="max-w-4xl mx-auto bg-white rounded-xl shadow-md overflow-hidden">
                <div className="md:flex">
                    {/* Left Side - Profile Info */}
                    <div className="md:w-1/2 p-6 bg-gradient-to-br from-blue-50 to-gray-50">
                        <div className="flex items-center mb-6">
                            <div className="bg-gradient-to-r from-blue-500 to-blue-600 text-white p-3 rounded-full mr-4 shadow-md">
                                <FaUserAlt size={24} />
                            </div>
                            <div>
                                <h2 className="text-2xl font-bold text-gray-800">{player?.userName}</h2>
                                <p className="text-gray-500 text-sm">Active Player</p>
                            </div>
                        </div>

                        <div className="space-y-4">
                            <div className="flex items-center bg-white p-3 rounded-lg shadow-sm">
                                <div className="bg-gray-100 p-2 rounded-full mr-3">
                                    <FaIdCard className="text-gray-600" />
                                </div>
                                <div>
                                    <p className="text-sm text-gray-500">Player ID</p>
                                    <p className="font-medium text-gray-800">{player?.id}</p>
                                </div>
                            </div>

                            <div className="flex items-center bg-white p-3 rounded-lg shadow-sm">
                                <div className="bg-yellow-100 p-2 rounded-full mr-3">
                                    <FaCrown className="text-yellow-600" />
                                </div>
                                <div>
                                    <p className="text-sm text-gray-500">Level</p>
                                    <p className="font-medium text-gray-800">
                                        {player?.level}
                                        <span className="ml-2 text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full">
                                            {player?.level && player.level > 10 ? 'Expert' : player?.level && player.level > 5 ? 'Intermediate' : 'Beginner'}
                                        </span>
                                    </p>
                                </div>
                            </div>

                            <div className="flex items-center bg-white p-3 rounded-lg shadow-sm">
                                <div className="bg-green-100 p-2 rounded-full mr-3">
                                    <FaChartLine className="text-green-600" />
                                </div>
                                <div>
                                    <p className="text-sm text-gray-500">Score</p>
                                    <p className="font-medium text-gray-800">{player?.score}</p>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Right Side - Match History */}
                    <div className="md:w-1/2 p-6">
                        <div className="flex items-center justify-between mb-4">
                            <h3 className="text-xl font-bold text-gray-800 flex items-center">
                                <FaHistory className="mr-2 text-blue-500" />
                                Match History
                            </h3>
                            {player?.matchHistory && player.matchHistory.length > 0 && (
                                <span className="text-sm bg-gray-100 text-gray-600 px-2 py-1 rounded-full">
                                    {player.matchHistory.length} matches
                                </span>
                            )}
                        </div>
                        
                        {player?.matchHistory && player.matchHistory.length > 0 ? (
                            <ul className="divide-y divide-gray-100">
                                {player.matchHistory.map((match, index) => (
                                    <li key={index} className="py-3 hover:bg-gray-50 px-2 rounded transition">
                                        <div className="flex items-center">
                                            <div className={`flex-shrink-0 h-2 w-2 rounded-full mr-3 ${
                                                match.includes('Win') ? 'bg-green-500' : 
                                                match.includes('Loss') ? 'bg-red-500' : 'bg-gray-500'
                                            }`}></div>
                                            <div>
                                                <p className="text-gray-800">{match}</p>
                                                <p className="text-xs text-gray-500 mt-1">2 days ago</p>
                                            </div>
                                        </div>
                                    </li>
                                ))}
                            </ul>
                        ) : (
                            <div className="text-center py-8">
                                <div className="bg-gray-100 p-4 rounded-lg inline-block">
                                    <FaHistory className="text-gray-400 text-3xl mx-auto mb-2" />
                                    <p className="text-gray-500">No match history available</p>
                                    <p className="text-sm text-gray-400 mt-1">Play some games to see your history</p>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Profile