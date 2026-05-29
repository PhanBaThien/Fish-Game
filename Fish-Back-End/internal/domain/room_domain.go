package domain

type CreateRoomRequest struct {
	Name        string  `json:"name"        binding:"required"`
	MinBet      int64   `json:"min_bet"     binding:"required,min=0"`
	MaxPlayers  int32   `json:"max_players" binding:"required,min=1"`
	Description *string `json:"description"`
	RTP         float64 `json:"rtp"         binding:"omitempty,min=0,max=1"`
}

type UpdateRoomRequest struct {
	Name        *string  `json:"name"`
	MinBet      *int64   `json:"min_bet"     binding:"omitempty,min=0"`
	MaxPlayers  *int32   `json:"max_players" binding:"omitempty,min=1"`
	Description *string  `json:"description"`
	RTP         *float64 `json:"rtp"         binding:"omitempty,min=0,max=1"`
}
