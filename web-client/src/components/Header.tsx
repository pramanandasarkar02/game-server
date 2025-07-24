import React from 'react'
import { useNavigate } from 'react-router-dom'
import { FaUserCircle } from 'react-icons/fa'

type Props = {}

const Header = (props: Props) => {
    const navigate = useNavigate()
    
    return (
        <header className="bg-gray-800 text-white p-4">
            <div className="container mx-auto flex justify-between items-center">
                <div className="text-2xl font-bold cursor-pointer" onClick={() => navigate('/')}>
                    Game Server
                </div>
                <div 
                    className="flex items-center space-x-2 cursor-pointer hover:text-gray-300 transition"
                    onClick={() => navigate('/profile')}
                >
                    <FaUserCircle size={24} />
                    <span>Profile</span>
                </div>
            </div>
        </header>
    )
}

export default Header