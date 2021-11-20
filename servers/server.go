package servers

type Server interface {
	Run(addr string) error
}

