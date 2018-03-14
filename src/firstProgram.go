package main

// Imports
import (
	"errors"
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
	consumerURL   string = "maltcp://127.0.0.1:"
	providerURL   string = "maltcp://127.0.0.1:"
	brokerURL     string = "maltcp://127.0.0.1:"
	subscriberURL string = "maltcp://127.0.0.1:"
	publisherURL  string = "maltcp://127.0.0.1:"
)

// Variables globales
var consumerURLPort = consumerURL
var providerURLPort = providerURL
var brokerURLPort = brokerURL
var subscriberURLPort = subscriberURL
var publisherURLPort = publisherURL

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
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

// ProgressProvider :
type ProgressProvider struct {
	ctx   *Context
	cctx  *ClientContext
	uri   *URI
	nbmsg int
}

// PubSubProvider :
type PubSubProvider struct {
	ctx  *Context
	cctx *ClientContext
	subs SubscriberTransaction
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

func newInvokeProvider() (*InvokeProvider, error) {
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

func newProgressProvider() (*ProgressProvider, error) {
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
	progressHandler1 := func(msg *Message, t Transaction) error {
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
	progressHandler2 := func(msg *Message, t Transaction) error {
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

func newPubSubProvider() (*PubSubProvider, error) {
	ctx, err := NewContext(brokerURLPort)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "broker")
	if err != nil {
		return nil, err
	}

	broker := &PubSubProvider{ctx, cctx, nil}

	// Register broker handler
	brokerHandler := func(msg *Message, t Transaction) error {
		if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER {
			broker.OnPublishRegister(msg, t.(PublisherTransaction))
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH {
			broker.OnPublish(msg, t.(PublisherTransaction))
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER {
			broker.OnPublishDeregister(msg, t.(PublisherTransaction))
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_REGISTER {
			broker.OnRegister(msg, t.(SubscriberTransaction))
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_DEREGISTER {
			broker.OnDeregister(msg, t.(SubscriberTransaction))
		} else {
			return errors.New("Bad stage")
		}

		return nil
	}

	// Register Broker handler
	cctx.RegisterBrokerHandler(2, 1, 2, 0, brokerHandler)

	return broker, nil
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

func (provider *PubSubProvider) close() {
	provider.cctx.Close()
}

// Méthodes utiles pour le publisher et le subscriber de la fonction pubSub
func (broker *PubSubProvider) OnRegister(msg *Message, tx SubscriberTransaction) error {
	fmt.Println("\t> OnRegister: ", string(msg.Body))
	broker.subs = tx
	tx.AckRegister(nil, false)
	return nil
}

func (broker *PubSubProvider) OnDeregister(msg *Message, tx SubscriberTransaction) error {
	fmt.Println("\t> OnDeregister:", string(msg.Body))
	broker.subs = nil
	tx.AckDeregister(nil, false)
	return nil
}

func (broker *PubSubProvider) OnPublishRegister(msg *Message, tx PublisherTransaction) error {
	fmt.Println("\t> OnPublishRegister:", string(msg.Body))
	tx.AckRegister(nil, false)
	return nil
}

func (broker *PubSubProvider) OnPublish(msg *Message, tx PublisherTransaction) error {
	fmt.Println("\t> OnPublish:", string(msg.Body))
	if broker.subs != nil {
		broker.subs.Notify(msg.Body, false)
	}
	return nil
}

func (broker *PubSubProvider) OnPublishDeregister(msg *Message, tx PublisherTransaction) error {
	fmt.Println("\t> OnPublishDeregister:", string(msg.Body))
	tx.AckDeregister(nil, false)
	return nil
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

func invoke(msg ...[]byte) error {
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

func progress(msg ...[]byte) error {
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

func pubSub() error {
	// Waiting for the previous socket to close
	time.Sleep(250 * time.Millisecond)

	// -- Broker --
	broker, err := newPubSubProvider()
	if err != nil {
		return err
	}
	defer broker.close()

	// -- Publisher --
	// Publisher Context
	pubCtx, err := NewContext(publisherURLPort)
	if err != nil {
		return err
	}
	defer pubCtx.Close()

	// Publisher
	publisher, err := NewClientContext(pubCtx, "publisher")
	if err != nil {
		return err
	}

	// Create publisher operation
	pubop := publisher.NewPublisherOperation(broker.cctx.Uri, 2, 1, 2, 0)
	// Call register method
	pubop.Register([]byte("register"))

	// -- Subscriber
	// Subscriber Context
	subCtx, err := NewContext(subscriberURLPort)
	if err != nil {
		return err
	}
	defer subCtx.Close()

	// Subscriber
	subscriber, err := NewClientContext(subCtx, "subscriber")
	if err != nil {
		return err
	}

	// Create subscriber operation
	subop := subscriber.NewSubscriberOperation(broker.cctx.Uri, 2, 1, 2, 0)
	// Call register operation
	subop.Register([]byte("register"))

	// Publish messages
	for i := 1; i <= 2; i++ {
		pubop.Publish([]byte(fmt.Sprintf("publish #%d", i)))
	}

	// Try to get notify by first publish
	resp1, err := subop.GetNotify()
	fmt.Println("\t>>> Subscriber notified: OK, ", string(resp1.Body))

	// Try to get notify by second publish
	resp2, err := subop.GetNotify()
	fmt.Println("\t>>> Subscriber notified: OK, ", string(resp2.Body))

	pubop.Deregister([]byte("deregister"))
	subop.Deregister([]byte("deregister"))

	return nil
}

// Fonction principale: lancement des échanges entre différents acteurs à l'aide de différents moyens
func main() {
	rand.Seed(time.Now().UnixNano())

	var port = 1023 + rand.Intn(15000)
	var consumerPort = port + 1 + rand.Intn(1000)
	var providerPort = consumerPort + 1 + rand.Intn(1000)
	var brokerPort = providerPort + 1 + rand.Intn(1000)
	var subscriberPort = brokerPort + 1 + rand.Intn(1000)
	var publisherPort = subscriberPort + 1 + rand.Intn(1000)

	consumerURLPort += strconv.Itoa(consumerPort)
	providerURLPort += strconv.Itoa(providerPort)
	brokerURLPort += strconv.Itoa(brokerPort)
	subscriberURLPort += strconv.Itoa(subscriberPort)
	publisherURLPort += strconv.Itoa(publisherPort)

	fmt.Println("consumer port   = ", consumerPort)
	fmt.Println("provider port   = ", providerPort)
	fmt.Println("broker port     = ", brokerPort)
	fmt.Println("subscriber port = ", subscriberPort)
	fmt.Println("publisher port  = ", publisherPort)

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

	// -- PubSub --
	errPubSub := pubSub()
	if errPubSub != nil {
		fmt.Println("Error: problem with pubSub function -> ", errPubSub)
	}
}
