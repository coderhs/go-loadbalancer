package balancer

type Backend struct {
	URL     string
	Healthy bool
}
