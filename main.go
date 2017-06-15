package main

import "github.com/dtop/go.demo.iam/iam"

func main() {

	srv := iam.New()
	srv.Run()
}
