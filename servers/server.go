package servers

const (
	APIVersion = "v1"
)

type Server interface {
	Run(addr string) error
}
