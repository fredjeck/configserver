package config

type Repository struct {
	Url             string
	RefreshInterval int
	Grants          Grants
}
type Repositories []Repository

type Grant struct {
	ClientId     string
	ClientSecret string
}
type Grants []Grant
