import React from 'react'
import { FaCopyright } from 'react-icons/fa'

const Footer = () => {
  return (
    <footer className="bg-gray-800 text-white p-4 mt-8">
        <div className="container mx-auto flex items-center justify-center space-x-2">
            <FaCopyright />
            <span>{new Date().getFullYear()} Game Server. All rights reserved.</span>
        </div>
    </footer>
  )
}

export default Footer