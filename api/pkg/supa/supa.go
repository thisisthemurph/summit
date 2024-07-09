package supa

import (
	"fmt"
	"github.com/nedpals/supabase-go"
)

const AuthEndpoint = "auth/v1"

type supabaseURLs struct {
	RecoverPassword string
	SignUp          string
}

type Supabase struct {
	Host   string
	Client *supabase.Client
	secret string
	URL    supabaseURLs
}

func New(host, secret string) *Supabase {
	return &Supabase{
		Host:   host,
		Client: supabase.CreateClient(host, secret),
		URL: supabaseURLs{
			RecoverPassword: fmt.Sprintf("%s/%s/recover", host, AuthEndpoint),
			SignUp:          fmt.Sprintf("%s/%s/signup", host, AuthEndpoint),
		},
		secret: secret,
	}
}
