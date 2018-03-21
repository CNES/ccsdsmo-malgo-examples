package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/api"

	_ "github.com/ccsdsmo/malgo/mal/transport/tcp"
)

var (
	consumerURL string = "maltcp://127.0.0.1:"
)

const (
	providerURLPort string = "maltcp://127.0.0.1:12400/provider"
)

func main() {
	fmt.Println("#### CONSUMER ####")

	rand.Seed(time.Now().UnixNano())

	var port = 1023 + rand.Intn(15000)
	var consumerPort = port + 1 + rand.Intn(1000)
	consumerURLPort := consumerURL + strconv.Itoa(consumerPort)

	fmt.Println("consumer url and port =", consumerURLPort)

	consumerCtx, err := NewContext(consumerURLPort)
	if err != nil {
		fmt.Println("Error creating context,", err)
		return
	}
	defer consumerCtx.Close()

	consumer, err := NewClientContext(consumerCtx, "consumer")
	if err != nil {
		fmt.Println("Error creating consumer,", err)
		return
	}

	providerURI := NewURI(providerURLPort)

	op := consumer.NewRequestOperation(providerURI, 2, 1, 2, 0)

	retour, err := op.Request([]byte("bidule chouette"))
	if err != nil {
		fmt.Println("Error,", err)
		return
	}
	fmt.Println("\t>>> Request1: OK, ", string(retour.Body))
}
