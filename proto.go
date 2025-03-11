package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const (
	CommandSET = "set"
	CommandGET = "get"
)

type Command interface {
}

type SetCommand struct {
	key, val string
}

func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Read %s\n", v.Type())

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandSET:
					fmt.Println(len(v.Array()))
					if (len(v.Array())) != 3 {
						return nil, fmt.Errorf("invalid set command")
					}
					setCmd := &SetCommand{
						key: value.Array()[1].String(),
						val: value.Array()[2].String(),
					}
					return setCmd, nil
				}
			}
		}
	}

	return "foo", nil
}
