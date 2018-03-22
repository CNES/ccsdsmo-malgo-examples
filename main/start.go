package main

import (
	"fmt"
	"os"

	. "github.com/EtienneLndr/MAL_API_Go_Project/archive"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("ERROR: You must use this program like this: go run start.go [provider|consumer]")
		return
	}

	var archiveService *ArchiveService
	// Create the Archive Service
	element := archiveService.CreateService()
	archiveService = element.(*ArchiveService)

	if args[0] == "provider" {
		archiveService.StartProvider()
	} else if args[0] == "consumer" {
		archiveService.StartConsumer()
	} else {
		fmt.Println("ERROR: You must use this program like this: go run start.go [provider|consumer]")
		return
	}
}
