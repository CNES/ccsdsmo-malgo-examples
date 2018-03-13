package main

// Imports
import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/api"
	_ "github.com/ccsdsmo/malgo/mal/transport/tcp"
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

// InvokeProvider :
type InvokeProvider struct {
	ctx *Context
	cctx *ClientContext
	nbmsg int
}

// ProgressProvider :
type ProgressProvider struct {
	ctx *Context
	cctx *ClientContext
	uri *URI
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

func newInvokeProvider () (*InvokeProvider, error) {
	ctx, err := NewContext(providerURLPort)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "provider")
	if err != nil {
		return nil, err
	}

	provider := &InvokeProvider{ctx, cctx, 0}

	// Register handler
	invokeHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(InvokeTransaction)
			fmt.Println("\t>>> invokeHandler receives: ", string(msg.Body))

			provider.nbmsg++
			transaction.Ack(nil, false)
			time.Sleep(250 * time.Millisecond)
			transaction.Reply(msg.Body, false)
		} else {
			fmt.Println("receive: nil")
		}

		return nil
	}
	cctx.RegisterInvokeHandler(2, 1, 2, 0, invokeHandler)

	return provider, nil
}

func newProgressProvider () (*ProgressProvider, error) {
	ctx, err := NewContext(providerURLPort)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "provider")
	if err != nil {
		return nil, err
	}

	provider := &ProgressProvider{ctx, cctx, cctx.Uri, 0}

	// Register handler
	// Handler 1
	progressHandler1 := func (msg *Message, t Transaction) error {
		provider.nbmsg++
		if msg != nil {
			fmt.Println("\t>>> progressHandler1 receives: ", string(msg.Body))

			transaction := t.(ProgressTransaction)

			transaction.Ack(nil, false)
			for i := 0; i < 10; i++ {
				transaction.Update([]byte(fmt.Sprintf("message1.#%d", i)), false)
			}
			transaction.Reply([]byte("last message1"), false)
		} else {
			fmt.Println("receive: nil")
		}

		return nil
	}

	// Register Progress handler 1
	cctx.RegisterProgressHandler(2, 1, 2, 0, progressHandler1)

	// Handler 2
	progressHandler2 := func (msg *Message, t Transaction) error {
		provider.nbmsg++
		if msg != nil {
			fmt.Println("\t>>> progressHandler2 receives: ", string(msg.Body))

			transaction := t.(ProgressTransaction)

			transaction.Ack(nil, false)
			for i := 0; i < 5; i++ {
				transaction.Update([]byte(fmt.Sprintf("message2.#%d", i)), false)
			}
			transaction.Reply([]byte("last message2"), false)
		} else {
			fmt.Println("receive: nil")
		}

		return nil
	}

	// Register Progress handler 2
	cctx.RegisterProgressHandler(2, 1, 2, 1, progressHandler2)

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

func (provider *InvokeProvider) close() {
	provider.ctx.Close()
}

func (provider *ProgressProvider) close() {
	provider.cctx.Close()
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

	if provider.nbmsg != len(msg) {
		fmt.Printf("Error: Received %d messages, expected %d.\n", provider.nbmsg, len(msg))
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

	if provider.nbmsg != len(msg) {
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

	// Waits for message reception
	time.Sleep(250 * time.Millisecond)

	if provider.nbmsg != len(msg) {
		fmt.Printf("Error: Received %d messages, expected %d.\n", provider.nbmsg, len(msg))
	}

	return nil
}

func invoke(msg... []byte) error {
	// Waiting for the previous socket to close
	time.Sleep(250 * time.Millisecond)

	provider, err := newInvokeProvider()
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
		return nil
	}

	firstOp := consumer.NewInvokeOperation(provider.cctx.Uri, 2, 1, 2, 0)
	_, errFirstOp := firstOp.Invoke(msg[0])
	if errFirstOp != nil {
		return errFirstOp
	}

	respFirstOp, errRespFirstOp := firstOp.GetResponse()
	if errRespFirstOp != nil {
		return errRespFirstOp
	}

	fmt.Println("\t>>> Invoke1: OK, ", string(respFirstOp.Body))

	secondOp := consumer.NewInvokeOperation(provider.cctx.Uri, 2, 1, 2, 0)
	_, errSecondOp := secondOp.Invoke(msg[1])
	if errSecondOp != nil {
		return errSecondOp
	}

	respSecondOp, errRespSecondOp := secondOp.GetResponse()
	if errRespSecondOp != nil {
		return errRespSecondOp
	}

	fmt.Println("\t>>> Invoke2: OK, ", string(respSecondOp.Body))

	// Waits for message reception
	time.Sleep(250 * time.Millisecond)

	if provider.nbmsg != len(msg) {
		fmt.Printf("Error: Received %d messages, expected %d.\n", provider.nbmsg, len(msg))
	}

	return nil
}

func progress(msg... []byte) error {
	// Waiting for the previous socket to close
	time.Sleep(250 * time.Millisecond)

	provider, err := newProgressProvider()
	if err != nil {
		return err
	}
	defer provider.close()

	consumer_ctx, err := NewContext(consumerURLPort)
	if err != nil {
		return err
	}
	defer consumer_ctx.Close()

	consumer, err := NewClientContext(consumer_ctx, "consumer")
	if err != nil {
		return err
	}

	nbmsg := 0

	// Call first progress operation
	firstOp := consumer.NewProgressOperation(provider.cctx.Uri, 2, 1, 2, 0)
	firstOp.Progress(msg[0])

	fmt.Println("\t>>> Progress1: OK")

	updt, err := firstOp.GetUpdate()
	if err != nil {
		return err
	}

	for updt != nil {
		nbmsg++
		fmt.Println("\t>>> Progress1: Update -> ", string(updt.Body))
		updt, err = firstOp.GetUpdate()
		if err != nil {
			return err
		}
	}

	resp, err := firstOp.GetResponse()
	if err != nil {
		return err
	}
	nbmsg++
	fmt.Println("\t>>> Progress1: Response -> ", string(resp.Body))

	if nbmsg != 11 {
		fmt.Printf("Error: Received %d messages, expected %d.\n", nbmsg, 11)
	}

	// Call second progress operation
	secondOp := consumer.NewProgressOperation(provider.cctx.Uri, 2, 1, 2, 1)
	secondOp.Progress(msg[0])

	fmt.Println("\t>>> Progress2: OK")

	updt, err = secondOp.GetUpdate()
	if err != nil {
		return err
	}

	for updt != nil {
		nbmsg++
		fmt.Println("\t>>> Progress2: Update -> ", string(updt.Body))
		updt, err = secondOp.GetUpdate()
		if err != nil {
			return err
		}
	}

	resp, err = secondOp.GetResponse()
	if err != nil {
		return err
	}
	nbmsg++
	fmt.Println("\t>>> Progress2: Response -> ", string(resp.Body))

	if nbmsg != 17 {
		fmt.Printf("Error: Received %d messages, expected %d.\n", nbmsg, 17)
	}

	// Waits for message reception
	time.Sleep(250 * time.Millisecond)

	if provider.nbmsg != len(msg) {
		fmt.Printf("Error: Received %d, expected %d.\n", provider.nbmsg, len(msg))
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

	// Call request method
	errRequest := request(msgRequest1, msgRequest2)
	if errRequest != nil {
		fmt.Println("Error: problem with request function -> ", errRequest)
	}

	// -- INVOKE --
	// Invoke variables
	var msgInvoke1 = []byte("invoke_test_1")
	var msgInvoke2 = []byte("invoke_test_2")

	// Call invoke method
	errInvoke := invoke(msgInvoke1, msgInvoke2)
	if errInvoke != nil {
		fmt.Println("Error: problem with invoke function -> ", errInvoke)
	}

	// -- PROGRESS
	// Progress variables
	var msgProgress1 = []byte("progress_test_1")
	var msgProgress2 = []byte("progress_test_2")

	// Call progress method
	errProgress := progress(msgProgress1, msgProgress2)
	if errProgress != nil {
		fmt.Println("Error: problem with progress function -> ", errProgress)
	}

}
