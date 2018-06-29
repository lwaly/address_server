package models

const (
	ErrDatabase = -1
	ErrSystem   = -2
	ErrDupRows  = -3
	ErrNotFound = -4
	ErrInput    = -5
)

const (
	db             = "address_server"
	address_server = "address_server"
	address_record = "address_record"
	address_gps = "address_gps"
)

const (
	VALID   = 0
	INVALID = 1
)

const (
	CITY     = 1
	PROVINCE = 2
	COUNTRY  = 3
	ALL      = 4
)
