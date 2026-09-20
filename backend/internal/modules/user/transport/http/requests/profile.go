package requests

type LocalizedUserProfileRequest struct {
	Name       *string `json:"name"`
	Surname    *string `json:"surname"`
	Patronymic *string `json:"patronymic"`
	Position   *string `json:"position"`
	Location   *string `json:"location"`
}

type ReplaceUserProfileRequest struct {
	RU       LocalizedUserProfileRequest `json:"ru"`
	EN       LocalizedUserProfileRequest `json:"en"`
	PhotoURL *string                     `json:"photoPath"`
}
