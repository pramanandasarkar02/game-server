import { createContext, useState } from "react"

interface Player {
    name: string
    token: string
}


interface AuthContextType {
    player: Player | null 
    setPlayer: React.Dispatch<React.SetStateAction<Player | null>>
}


export const AuthContext = createContext<AuthContextType>({
    player: null,
    setPlayer: () => {}
})


interface AuthContextProviderProps {
    children: React.ReactNode
}

export const AuthContextProvider: React.FC<AuthContextProviderProps> = ({ children }) => {
    const [player, setPlayer] = useState<Player | null>(null)
    return (
        <AuthContext.Provider value={{ player, setPlayer }}>
            {children}
        </AuthContext.Provider>
    )
}