package main

func main() {
	router := CreateRouter()
	router.Run() //nolint:errcheck
}
