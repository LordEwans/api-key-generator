package internal

import (
	"log"
)

func Handle(err error) {
	if err != nil {
		log.Panicf("Error Occured: %s\n", err)
	}
}
