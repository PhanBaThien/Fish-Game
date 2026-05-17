package domain

type CreateFishRequest struct {
	Name             string  `json:"name"              binding:"required"`
	Health           int32   `json:"health"            binding:"required,min=1"`
	RewardMultiplier int32   `json:"reward_multiplier" binding:"required,min=1"`
	Speed            float64 `json:"speed"             binding:"required,min=0"`
	AssetPath        string  `json:"asset_path"        binding:"required"`
}

type UpdateFishRequest struct {
	Name             *string  `json:"name"`
	Health           *int32   `json:"health"            binding:"omitempty,min=1"`
	RewardMultiplier *int32   `json:"reward_multiplier" binding:"omitempty,min=1"`
	Speed            *float64 `json:"speed"             binding:"omitempty,min=0"`
	AssetPath        *string  `json:"asset_path"`
}
