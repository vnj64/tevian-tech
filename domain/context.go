package domain

type Context interface {
	Connection() Connection
	Make() Context
	Services() Services
}
