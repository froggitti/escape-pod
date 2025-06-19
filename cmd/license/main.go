package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/DDLbots/escape-pod/internal/license/format"
	"github.com/DDLbots/escape-pod/internal/license/issuer"
)

func main() {
	email := flag.String("email", "", "The email address of the user")
	robot := flag.String("robot", "", "the robot s/n.  NOTE:  use vic: prefix.")
	flag.Parse()

	if *email == "" || *robot == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	is := issuer.New()

	req := format.License{
		Email: *email,
		Bot:   strings.ToLower(*robot),
	}

	res, err := is.Generate(&req)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	fmt.Println(res)
}
