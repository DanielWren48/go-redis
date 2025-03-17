package main

const (
	CommandSET    = "set"
	CommandGET    = "get"
	CommandHELLO  = "hello"
	CommandClient = "client"
)

type Command interface {
}

type SetCommand struct {
	key, val []byte
}

type ClientCommand struct {
	value string
}

type HelloCommand struct {
	value string
}

type GetCommand struct {
	key []byte
}
