package routes

type flash struct {
	notice string
	alert  string

	custom map[string]string
}
