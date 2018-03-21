package main

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/api"

	_ "github.com/ccsdsmo/malgo/mal/transport/tcp"
)

const (
	providerURLPort string = "maltcp://127.0.0.1:12400/provider"
)

type RequestProvider struct {
	ctx  *Context
	cctx *ClientContext
}

func newRequestProvider() (*RequestProvider, error) {
	ctx, err := NewContext(providerURLPort)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "provider")
	if err != nil {
		return nil, err
	}

	provider := &RequestProvider{ctx, cctx}

	// Register handler
	/*requestHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(RequestTransaction)
			fmt.Println("\t>>> requestHandler receives: ", string(msg.Body))
			transaction.Reply([]byte("reply message"), false)
		} else {
			fmt.Println("receive: nil")
		}

		return nil
	}*/

	//cctx.RegisterRequestHandler(2, 1, 2, 0, requestHandler)

	return provider, nil
}

func (provider *RequestProvider) close() {
	provider.ctx.Close()
}

func main() {
	fmt.Println("#### PROVIDER ####")

	rand.Seed(time.Now().UnixNano())

	fmt.Println("provider url and port =", providerURLPort)

	provider, err := newRequestProvider()
	if err != nil {
		fmt.Println("Error creating provider, ", err)
		return
	}
	defer provider.close()

	time.Sleep(120 * time.Second)
}
