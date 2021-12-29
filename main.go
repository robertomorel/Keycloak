package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	cliendID     = "myclient"
	clientSecret = "Jruhk14ruaKZCv2AlkhInFTyzXfzZ0rv"
)

func main() {
	// Pacote que nos possibilita em parar uma solicitação no meio dela
	ctx := context.Background()

	/*
		Endpoints do Keycloak: http://localhost:8080/auth/realms/master/.well-known/openid-configuration
	*/
	// Cria providers usando o OpenId conect
	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/auth/realms/master" /* End. do processo de login */)

	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     cliendID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	/*
	  Sempre que iniciamos o processo de oauth, mandamos o state.
	  No recebimento, devemos receber o mesmo state e isto deve ser verificado
	*/
	state := "123"

	// Redirect URL
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		/*
		  Sempre que acessamos "/", vamos ser redirecionados à uma uri de autenticação
		  que será gerada automaticamente com o AuthCodeURL
		*/
		http.Redirect(rw, r, config.AuthCodeURL(state), http.StatusFound)
	})

	// Gerando o Access Token (autenticação) e ID Token (autorização)
	http.HandleFunc("/auth/callback", func(rw http.ResponseWriter, r *http.Request) {
		// Verificando o state
		if r.URL.Query().Get("state") != state {
			http.Error(rw, "State inválido", http.StatusBadRequest)
			return
		}

		// Exchange troca um code que recebemos por um token
		// Pegando o Access Token
		token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(rw, "Falha ao trocar o token", http.StatusInternalServerError)
			return
		}

		// Add o ID Token
		idToken, ok := token.Extra("id_token").(string) // .(string) Transforma em string
		if !ok {
			http.Error(rw, "Falha ao gerar o ID Token", http.StatusInternalServerError)
			return
		}

		userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
		if err != nil {
			http.Error(rw, "Erro ao pegar User Info", http.StatusInternalServerError)
			return
		}

		// A resposta vai ser uma struct que iremos converter para JSON
		resp := struct {
			AccessToken *oauth2.Token
			IDToken     string
			UserInfo    *oidc.UserInfo
		}{
			token,
			idToken,
			userInfo,
		}

		// Criando o json
		data, err := json.Marshal(resp)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Write(data)
	})

	// Subindo o servidor web
	log.Fatal(http.ListenAndServe(":8081", nil))
}
