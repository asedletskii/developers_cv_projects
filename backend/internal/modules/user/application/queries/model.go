package queries

type Profile struct {
	ID         string  `json:"userId"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic,omitempty"`
	Position   *string `json:"position,omitempty"`
	Location   *string `json:"location,omitempty"`
	PhotoPath  *string `json:"photoPath,omitempty"`
}
