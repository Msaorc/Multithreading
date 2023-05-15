package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var urlCdnApiCep = "https://cdn.apicep.com/file/apicep/%s.json"
var urlViaCep = "http://viacep.com.br/ws/%s/json/"

type Cep struct {
	Url string
	Cep string
}

func main() {
	cdnApiChannel := make(chan Cep)
	viaCepChannel := make(chan Cep)
	cep := getCepCommandLine()

	go func() {
		cdnApiChannel <- findCEP(cep, urlCdnApiCep)
	}()

	go func() {
		viaCepChannel <- findCEP(cep, urlViaCep)
	}()

	select {
	case cepResult := <-cdnApiChannel:
		printInConsole(&cepResult)
	case cepResult := <-viaCepChannel:
		printInConsole(&cepResult)
	case <-time.After(time.Second * 1):
		println("Error: Timeout!")
	}
}

func findCEP(cep string, url string) Cep {
	r, err := http.Get(fmt.Sprintf(url, cep))
	if err != nil {
		return Cep{
			Url: fmt.Sprintf(url, cep),
			Cep: err.Error(),
		}
	}
	defer r.Body.Close()
	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Cep{
			Url: fmt.Sprintf(url, cep),
			Cep: err.Error(),
		}
	}
	return Cep{
		Url: fmt.Sprintf(url, cep),
		Cep: string(result),
	}
}

func getCepCommandLine() string {
	cep := ""
	for _, cep2 := range os.Args[1:] {
		cep = fmt.Sprint(cep2)
	}
	return cep
}

func printInConsole(cep *Cep) {
	fmt.Printf("Url utilizada: %s\nCep: %s", cep.Url, cep.Cep)
}
