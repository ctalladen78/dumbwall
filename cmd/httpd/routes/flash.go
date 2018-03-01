package routes

type flash struct {
	Notice string
	Alert  string

	Errors []string

	Custom map[string]string
}
