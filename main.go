package main

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

var (
	allowedPublicKeys    = map[string]ssh.PublicKey{}
	challengedPublicKeys = map[string]ssh.PublicKey{}
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, fmt.Sprintf("Welcome, %s!\n", s.User()))
		t := term.NewTerminal(s, fmt.Sprintf("(%s) ~ ", s.User()))
		for {
			line, err := t.ReadLine()
			if err != nil {
				break
			}
			if _, err := t.Write([]byte(line + "\n")); err != nil {
				break
			}
		}
	})
	ssh.ListenAndServe(":2222", nil,
		ssh.KeyboardInteractiveAuth(func(ctx ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool {
			if _, ok := allowedPublicKeys[ctx.User()]; ok {
				// Username is already claimed.
				return false
			}
			challengedKey, ok := challengedPublicKeys[ctx.User()]
			if !ok {
				// No public-key challenge is initiated.
				return false
			}
			a, b := rand.Intn(10), rand.Intn(10)
			answers, err := challenger(
				"New User Registration",
				"Solve the following challenge to complete your registration.",
				[]string{fmt.Sprintf("%d + %d: ", a, b)},
				[]bool{true},
			)
			if err != nil {
				return false
			}
			answer, err := strconv.Atoi(answers[0])
			if err != nil || answer != a+b {
				return false
			}
			allowedPublicKeys[ctx.User()] = challengedKey
			delete(challengedPublicKeys, ctx.User())
			return true
		}),
		ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			if allowedKey, ok := allowedPublicKeys[ctx.User()]; ok {
				return ssh.KeysEqual(key, allowedKey)
			}
			challengedPublicKeys[ctx.User()] = key
			return false
		}),
	)
}
