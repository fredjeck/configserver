package config

type Repository struct {
	Name            string
	Url             string
	Token           string
	RefreshInterval int
	Grants          Grants
}
type Repositories []Repository

type Grant struct {
	ClientId     string
	ClientSecret string
}
type Grants []Grant
