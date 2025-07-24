import React from 'react'
import { Navigate, useNavigate } from 'react-router-dom'

type Props = {}

const Header = (props: Props) => {
    const navigate = useNavigate()
  return (
    <div>
        <div>
            Game Server
        </div>
        <div onClick={()=> {navigate('/profile')}}>
            User Profile Icon
        </div>
    </div>
  )
}

export default Header