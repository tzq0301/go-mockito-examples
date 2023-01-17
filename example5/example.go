package example

import "log"

type Service interface {
	PingPong(ping string) (pong string)
}

func NewService() Service {
	return &serviceImpl{}
}

type serviceImpl struct{}

func (e serviceImpl) PingPong(ping string) (pong string) {
	log.Println("Begin Do some RPC Call !!!")
	pong = ping
	log.Println("End   Do some RPC Call !!!")
	return
}

// --------------------------------------------------------

type Handler struct {
	service Service
}

func (h Handler) PingPong(ping string) (pong string) {
	pong = h.service.PingPong(ping)
	return
}
