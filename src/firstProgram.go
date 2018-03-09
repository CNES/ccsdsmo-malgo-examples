package main

// Imports
import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	. "github.com/ccsdsmo/malgo/src/mal"
	. "github.com/ccsdsmo/malgo/src/mal/api"
	_ "github.com/ccsdsmo/malgo/src/mal/transport/tcp"
)

// Constantes
const (
	consumerURL string = "maltcp://127.0.0.1:"
	providerURL string = "maltcp://127.0.0.1:"
)

// Variables globales
var consumerURLPort = consumerURL
var providerURLPort = providerURL

// SendProvider :
type SendProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

// SubmitProvider :
type SubmitProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

// RequestProvider :
type RequestProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

// Création du provider en fonction du moyen de communication
func newSendProvider() (*SendProvider, error) {
	ctx, err := NewContext(providerURLPort)
	if err != nil {
		fmt.Println("context")
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "provider")
	if err != nil {
		fmt.Println("client context")
		return nil, err
	}
	provider := &SendProvider{ctx, cctx, 0}

	// Register handler
	sendHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			fmt.Println("\t>>> sendHandler receives: ", string(msg.Body))
			provider.nbmsg++
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	cctx.RegisterSendHandler(2, 1, 2, 0, sendHandler)

	return provider, nil
}

func newSubmitProvider() (*SubmitProvider, error) {
	ctx, err := NewContext(providerURLPort)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "provider")
	if err != nil {
		return nil, err
	}
	provider := &SubmitProvider{ctx, cctx, 0}

	// Register handler
	submitHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(SubmitTransaction)
			fmt.Println("\t>>> submitHandler receives: ", string(msg.Body))
			provider.nbmsg++
			transaction.Ack(nil, false)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}

	cctx.RegisterSubmitHandler(2, 1, 2, 0, submitHandler)

	return provider, nil
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
	provider := &RequestProvider{ctx, cctx, 0}

	// Register handler
	requestHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(RequestTransaction)
			fmt.Println("\t>>> requestHandler receives: ", string(msg.Body))
			provider.nbmsg++
			transaction.Reply([]byte("reply message"), false)
		} else {
			fmt.Println("receive: nil")
		}

		return nil
	}

	cctx.RegisterRequestHandler(2, 1, 2, 0, requestHandler)

	return provider, nil
}

// Méthodes permettant de clôturer un provider
func (provider *SendProvider) close() {
	provider.ctx.Close()
}

func (provider *SubmitProvider) close() {
	provider.ctx.Close()
}

func (provider *RequestProvider) close() {
	provider.ctx.Close()
}

//
func send(msg ...[]byte) error {
	// Waiting for the previous socket to close
	time.Sleep(250 * time.Millisecond)

	provider, err := newSendProvider()
	if err != nil {
		return err
	}
	defer provider.close()

	consumerCtx, err := NewContext(consumerURLPort)
	if err != nil {
		return err
	}
	defer consumerCtx.Close()

	consumer, err := NewClientContext(consumerCtx, "consumer")
	if err != nil {
		return err
	}

	// Create submit operation
	firstOp := consumer.NewSendOperation(provider.cctx.Uri, 2, 1, 2, 0)
	// Call send method
	firstOp.Send(msg[0])

	// Create submit operation
	secondOp := consumer.NewSendOperation(provider.cctx.Uri, 2, 1, 2, 0)
	// Call send method
	secondOp.Send(msg[1])

	// Waits for message reception
	time.Sleep(250 * time.Millisecond)

	if provider.nbmsg != 2 {
		fmt.Printf("Error: Received %d messages, expected %d.\n", provider.nbmsg, 2)
	}

	return nil
}

func submit(msg ...[]byte) error {
	// Waiting for the previous socket to close
	time.Sleep(250 * time.Millisecond)

	provider, err := newSubmitProvider()
	if err != nil {
		return err
	}
	defer provider.close()

	consumerCtx, err := NewContext(consumerURLPort)
	if err != nil {
		return err
	}
	defer consumerCtx.Close()

	consumer, err := NewClientContext(consumerCtx, "consumer")
	if err != nil {
		return err
	}

	// Create submit operation
	firstOp := consumer.NewSubmitOperation(provider.cctx.Uri, 2, 1, 2, 0)
	// Call submit method
	_, errFirstOp := firstOp.Submit(msg[0])
	if errFirstOp != nil {
		return errFirstOp
	}
	fmt.Println("\t>>> First Submit: OK")

	// Create submit operation
	secondOp := consumer.NewSubmitOperation(provider.cctx.Uri, 2, 1, 2, 0)
	// Call submit method
	_, errSecondOp := secondOp.Submit(msg[1])
	if errSecondOp != nil {
		return errSecondOp
	}
	fmt.Println("\t>>> Second Submit: OK")

	// Waits for message reception
	time.Sleep(250 * time.Millisecond)

	if provider.nbmsg != 2 {
		fmt.Printf("Error: Received %d messages, expected %d.\n", provider.nbmsg, len(msg))
	}

	return nil
}

func request(msg ...[]byte) error {
	// Waiting for the previous socket to close
	time.Sleep(250 * time.Millisecond)

	provider, err := newRequestProvider()
	if err != nil {
		return err
	}
	defer provider.close()

	consumerCtx, err := NewContext(consumerURLPort)
	if err != nil {
		return err
	}
	defer consumerCtx.Close()

	consumer, err := NewClientContext(consumerCtx, "consumer")
	if err != nil {
		return err
	}

	// Create first request operation
	firstOp := consumer.NewRequestOperation(provider.cctx.Uri, 2, 1, 2, 0)
	// Call request method
	ret1, err := firstOp.Request(msg[0])
	if err != nil {
		return err
	}
	fmt.Println("\t>>> Request1: OK, ", string(ret1.Body))

	// Create second request operation
	secondOp := consumer.NewRequestOperation(provider.cctx.Uri, 2, 1, 2, 0)
	// Call request method
	ret2, err := secondOp.Request(msg[1])
	if err != nil {
		return err
	}
	fmt.Println("\t>>> Request2: OK, ", string(ret2.Body))

	if provider.nbmsg != 2 {
		fmt.Printf("Error: Received %d messages, expected %d.\n", provider.nbmsg, len(msg))
	}

	return nil
}

// Fonction principale: lancement des échanges entre différents acteurs à l'aide de différents moyens
func main() {
	rand.Seed(time.Now().UnixNano())

	var consumerPort = 1023 + rand.Intn(15000)
	var providerPort = 1023 + rand.Intn(15000)

	consumerURLPort += strconv.Itoa(consumerPort)
	providerURLPort += strconv.Itoa(providerPort)

	fmt.Println("consumer port = ", consumerPort)
	fmt.Println("provider port = ", providerPort)

	// -- SEND --
	// Send variables
	var msgSend1 = []byte("send_test_1")
	var msgSend2 = []byte("send_test_2")

	// Call send method
	errSend := send(msgSend1, msgSend2)
	if errSend != nil {
		fmt.Println("Error: problem with send function -> ", errSend)
	}

	// -- SUBMIT --
	// Submit variables
	var msgSubmit1 = []byte("submit_test_1")
	var msgSubmit2 = []byte("submit_test_2")

	// Call sumbit method
	errSubmit := submit(msgSubmit1, msgSubmit2)
	if errSubmit != nil {
		fmt.Println("Error: problem with submit function -> ", errSubmit)
	}

	// -- REQUEST --
	// Request variables
	var msgRequest1 = []byte("request_test_1")
	var msgRequest2 = []byte("request_test_2")

	// Call request mathod
	errRequest := request(msgRequest1, msgRequest2)
	if errRequest != nil {
		fmt.Println("Error: problem with request function -> ", errRequest)
	}

}
