package aurora

type Project struct {
	Cluster      string `flag:"cluster, The Aurora cluster"`
	Contact      string `flag:"contact, User"`
	Role         string `flag:"role, The role"`
	Environment  string `flag:"environment, The environment"`
	IsProduction bool   `flag:"prod, True if this is for production"`
}
