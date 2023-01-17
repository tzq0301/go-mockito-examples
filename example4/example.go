package example

import (
	"log"
	"strings"
)

func refresh() {
	log.Println("refreshing...")
	log.Println("refreshing...")
}

func refreshAsync() {
	go refresh()
}

func isSuperUser(user string) bool {
	refreshAsync()

	return strings.HasPrefix(user, "super")
}
