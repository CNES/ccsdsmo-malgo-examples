package main

// Imports
import (
	"fmt"
	"time"
	"strconv"
	"math/rand"
	. "github.com/ccsdsmo/malgo/src/mal"
	. "github.com/ccsdsmo/malgo/src/mal/api"
	_ "github.com/ccsdsmo/malgo/src/mal/transport/tcp"
)

// Constantes
const (
	consumer_url string = "maltcp://127.0.0.1:"
	provider_url string = "maltcp://127.0.0.1:"
)

// Variables globales
var consumer_url_port = consumer_url
var provider_url_port = provider_url

// Définition de nouveaux types
type SendProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

type SubmitProvider struct {
	ctx *Context
	cctx *ClientContext
	nbmsg int
}

type RequestProvider struct {
	ctx *Context
	cctx *ClientContext
	nbmsg int
}

// Création du provider en fonction du moyen de communication
func newSendProvider () (*SendProvider, error) {
	ctx, err := NewContext(provider_url_port);
	if err != nil {
		fmt.Println("context");
		return nil, err;
	}

	cctx, err := NewClientContext(ctx, "provider");
	if err != nil {
		fmt.Println("client context");
		return nil, err;
	}
	provider := &SendProvider{ctx, cctx, 0};

	// Register handler
	sendHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			fmt.Println("\t>>> sendHandler receives: ", string(msg.Body));
			provider.nbmsg++;
		} else {
			fmt.Println("receive: nil");
		}
		return nil;
	}
	cctx.RegisterSendHandler(2, 1, 2, 0, sendHandler);

	return provider, nil;
}

func newSubmitProvider () (*SubmitProvider, error) {
	ctx, err := NewContext(provider_url_port);
	if (err != nil) {
		return nil, err;
	}

	cctx, err := NewClientContext(ctx, "provider");
	if (err != nil) {
		return nil, err;
	}
	provider := &SubmitProvider{ctx, cctx, 0};

	// Register handler
	submitHandler := func(msg *Message, t Transaction) error {
		if (msg != nil) {
			transaction := t.(SubmitTransaction);
			fmt.Println("\t>>> submitHandler receives: ", string(msg.Body));
			provider.nbmsg++;
			transaction.Ack(nil, false);
		} else {
			fmt.Println("receive: nil");
		}
		return nil;
	}

	cctx.RegisterSubmitHandler(2, 1, 2, 0, submitHandler);

	return provider, nil;
}

func newRequestProvider () (*RequestProvider, error) {
	ctx, err := NewContext(provider_url_port);
	if (err != nil) {
		return nil, err;
	}

	cctx, err := NewClientContext(ctx, "provider");
	if (err != nil) {
		return nil, err;
	}
	provider := &RequestProvider{ctx, cctx, 0};

	// Register handler
	requestHandler := func(msg *Message, t Transaction) error {
		if (msg != nil) {
			transaction := t.(RequestTransaction);
			fmt.Println("\t>>> requestHandler receives: ", string(msg.Body));
			provider.nbmsg++;
			transaction.Reply([]byte("reply message"), false);
		} else {
			fmt.Println("receive: nil");
		}

		return nil;
	}

	cctx.RegisterRequestHandler(2, 1, 2, 0, requestHandler);

	return provider, nil;
}

// Méthodes permettant de clôturer un provider
func (provider *SendProvider) close() {
	provider.ctx.Close()
}

func (provider *SubmitProvider) close() {
	provider.ctx.Close();
}

func (provider *RequestProvider) close() {
	provider.ctx.Close();
}

//
func send (msg... []byte) error {
	// Waiting for the previous socket to close 
	time.Sleep(250 * time.Millisecond);

	provider, err := newSendProvider();
	if (err != nil) {
		return err;
	}
	defer provider.close();

	consumer_ctx, err := NewContext(consumer_url_port);
	if (err != nil) {
		return err;
	}
	defer consumer_ctx.Close();

	consumer, err := NewClientContext(consumer_ctx, "consumer");
	if (err != nil) {
		return err;
	}

	// Create submit operation
	first_op := consumer.NewSendOperation(provider.cctx.Uri, 2, 1, 2, 0);
	// Call send method
	first_op.Send(msg[0]);

	// Create submit operation
	second_op := consumer.NewSendOperation(provider.cctx.Uri, 2, 1, 2, 0);
	// Call send method
	second_op.Send(msg[1]);

	// Waits for message reception
	time.Sleep(250 * time.Millisecond);

	if (provider.nbmsg != 2) {
		fmt.Printf("Error: Received %d messages, expected %d.\n", provider.nbmsg, 2);
	}

	return nil;
}

func submit (msg... []byte) error {
	// Waiting for the previous socket to close
	time.Sleep(250 * time.Millisecond);

	provider, err := newSubmitProvider();
	if (err != nil) {
		return err;
	}
	defer provider.close();

	consumer_ctx, err := NewContext(consumer_url_port);
	if (err != nil) {
		return err;
	}
	defer consumer_ctx.Close();

	consumer, err := NewClientContext(consumer_ctx, "consumer");
	if (err != nil) {
		return err;
	}

	// Create submit operation
	first_op := consumer.NewSubmitOperation(provider.cctx.Uri, 2, 1, 2, 0);
	// Call submit method
	_, err_first_op := first_op.Submit(msg[0]);
	if (err_first_op != nil) {
		return err_first_op;
	}
	fmt.Println("\t>>> First Submit: OK");

	// Create submit operation
	second_op := consumer.NewSubmitOperation(provider.cctx.Uri, 2, 1, 2, 0);
	// Call submit method
	_, err_second_op := second_op.Submit(msg[1]);
	if (err_second_op != nil) {
		return err_second_op;
	}
	fmt.Println("\t>>> Second Submit: OK");

	// Waits for message reception
	time.Sleep(250 * time.Millisecond);

	if (provider.nbmsg != 2) {
		fmt.Printf("Error: Received %d messages, expected %d.\n", provider.nbmsg, len(msg));
	}

	return nil;
}

func request(msg... []byte) error {
	// Waiting for the previous socket to close
	time.Sleep(250 * time.Millisecond);

	provider, err := newRequestProvider();
	if (err != nil) {
		return err;
	}
	defer provider.close();

	consumer_ctx, err := NewContext(consumer_url_port);
	if (err != nil) {
		return err;
	}
	defer consumer_ctx.Close();

	consumer, err := NewClientContext(consumer_ctx, "consumer");
	if (err != nil) {
		return err;
	}

	// Create first request operation
	first_op := consumer.NewRequestOperation(provider.cctx.Uri, 2, 1, 2, 0);
	// Call request method
	ret1, err := first_op.Request(msg[0]);
	if (err != nil) {
		return err;
	}
	fmt.Println("\t>>> Request1: OK, ", string(ret1.Body));

	// Create second request operation
	second_op := consumer.NewRequestOperation(provider.cctx.Uri, 2, 1, 2, 0);
	// Call request method
	ret2, err := second_op.Request(msg[1]);
	if (err != nil) {
		return err;
	}
	fmt.Println("\t>>> Request2: OK, ", string(ret2.Body));

	if (provider.nbmsg != 2) {
		fmt.Printf("Error: Received %d messages, expected %d.\n", provider.nbmsg, len(msg));
	}

	return nil;
}

// Fonction principale: lancement des échanges entre différents acteurs à l'aide de différents moyens
func main () {
	rand.Seed(time.Now().UnixNano());

	var consumer_port int = 1023 + rand.Intn(15000);
	var provider_port int = 1023 + rand.Intn(15000);

	consumer_url_port += strconv.Itoa(consumer_port);
	provider_url_port += strconv.Itoa(provider_port);

	fmt.Println("consumer port = ", consumer_port);
	fmt.Println("provider port = ", provider_port);

	// -- SEND --
	// Send variables
	var msg1_send = []byte("send_test_1");
	var msg2_send = []byte("send_test_2");

	// Call send method
	err_send := send(msg1_send, msg2_send);
	if (err_send != nil) {
		fmt.Println("Error: problem with send function -> ", err_send);
	}

	// -- SUBMIT --
	// Submit variables
	var msg1_submit = []byte("submit_test_1");
	var msg2_submit = []byte("submit_test_2");

	// Call sumbit method
	err_submit := submit(msg1_submit, msg2_submit);
	if (err_submit != nil) {
		fmt.Println("Error: problem with submit function -> ", err_submit);
	}

	// -- REQUEST --
	// Request variables
	var msg1_request = []byte("request_test_1");
	var msg2_request = []byte("request_test_2");

	// Call request mathod
	err_request := request(msg1_request, msg2_request);
	if (err_request != nil) {
		fmt.Println("Error: problem with request function -> ", err_request);
	}

}
