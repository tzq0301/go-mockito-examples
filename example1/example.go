package example

var superUsers []string = getSuperUsers()

func getSuperUsers() []string {
	// Get From RPC
	return nil
}

func isSuperUser(user string) bool {
	for _, superUser := range superUsers {
		if superUser == user {
			return true
		}
	}
	return false
}
