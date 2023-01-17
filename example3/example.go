package example

import "fmt"

func getSuperUsers() []string {
	// Get From RPC
	var superUsers []string
	for i := 1; i <= 3; i++ {
		superUsers = append(superUsers, fmt.Sprintf("super%d from rpc", i))
	}
	return superUsers
}

func isSuperUser(user string) bool {
	for _, superUser := range getSuperUsers() {
		if superUser == user {
			return true
		}
	}
	return false
}
