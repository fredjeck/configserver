package config

type Repository struct {
	Name            string
	Url             string
	Token           string
	RefreshInterval int
	Clients         []string
}
type Repositories []Repository
