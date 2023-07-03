package main

import (
	"io/ioutil"
	"log"

	"github.com/Icarus-xD/DocRuleVerifier/pkg/verifier"
)

const RULE = "Hello и (Golang или Test) и Индия"
const FILE = "test.txt"

func main() {
	b, err := ioutil.ReadFile(FILE)
	if err != nil {
		log.Fatal(err)
	}
	doc := string(b)

	verifier.Verify(&doc, RULE)
}