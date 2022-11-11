package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func request(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

func apiCep(cep string, ch chan<- string) {
	ch <- request("https://cdn.apicep.com/file/apicep/" + cep + ".json")
}

func viaCep(cep string, ch chan<- string) {
	ch <- request("http://viacep.com.br/ws/" + cep + "/json/")
}

func getCep(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cepCh := make(chan string)

	cep := chi.URLParam(r, "cep")

	go viaCep(cep, cepCh)
	go apiCep(cep, cepCh)

	select {
	case cep := <-cepCh:
		io.WriteString(w, cep)
	case <-time.After(time.Second * 1):
		io.WriteString(w, "timeout")
	}

}

func main() {
	r := chi.NewRouter()
	r.Get("/{cep}", getCep)
	log.Fatal(http.ListenAndServe(":8080", r))
}
