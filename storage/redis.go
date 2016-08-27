package storage

import "fmt"

func NewStorage() {

}

func HandleChan(c chan string) {
	for {
		key := <-c
		fmt.Printf("Key: %s\n", key)
	}
}
