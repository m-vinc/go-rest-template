package interfaces

type ICryptoService interface {
	// Need to implement a function which take a username and create a hash which identify the user within the service
	GenerateUserHash(username string) (string, error)
}
