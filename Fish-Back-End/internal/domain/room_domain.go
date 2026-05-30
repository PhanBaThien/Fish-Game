package domain

type CreateRoomRequest struct {
	Name        string  `json:"name"        binding:"required"`
	MaxPlayers  int32   `json:"max_players" binding:"required,min=1"`
	Description *string `json:"description"`
	RTP         float64 `json:"rtp"         binding:"omitempty,min=0,max=1"`
}

type UpdateRoomRequest struct {
	Name        *string  `json:"name"`
	MaxPlayers  *int32   `json:"max_players" binding:"omitempty,min=1"`
	Description *string  `json:"description"`
	RTP         *float64 `json:"rtp"         binding:"omitempty,min=0,max=1"`
}
