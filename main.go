package main

func main() {
	a := App{}
	a.Initialize("admin", "123", "go_api")

	a.Run(":8080")
}
