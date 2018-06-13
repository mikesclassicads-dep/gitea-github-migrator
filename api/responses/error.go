package responses

type Error struct {
	Description string
}

func (e Error) Error() string {
	return e.Description
}
