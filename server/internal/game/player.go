package game


type PlayerStatus int

const (
	PLAYING PlayerStatus = iota
	ONLINE
	OFFLINE
)


type Player struct{
	Username string		`json:"username"` 
	UserId 	string `json:"userId"`
	PlayerStatus PlayerStatus	
}

type SignupRequest struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogOutRequest struct{
	Username string `json:"username"`
}

type LoginRequest struct{
	Username string `json:"username"`
	Password string `json:"password"`
}


func(p *Player) Login(request LoginRequest){
	//  write login logic here
	p.PlayerStatus = ONLINE
}

func(p *Player)Logout(request LogOutRequest){

}

func(p *Player)SignUp(request SignupRequest){

}



