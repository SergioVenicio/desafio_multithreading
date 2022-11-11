package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func request(url string) string {
	log := fmt.Sprintf("[request] %s...", url)
	fmt.Println(log)
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

	viaCh := make(chan string)
	apiCh := make(chan string)

	cep := chi.URLParam(r, "cep")

	go viaCep(cep, viaCh)
	go apiCep(cep, apiCh)

	select {
	case cep := <-viaCh:
		fmt.Println("[getCep] used viaCep...")
		io.WriteString(w, cep)
		break
	case cep := <-apiCh:
		fmt.Println("[getCep] used apicep...")
		io.WriteString(w, cep)
		break
	case <-time.After(time.Second * 1):
		w.WriteHeader(http.StatusRequestTimeout)
		break
	}

}

func main() {
	r := chi.NewRouter()
	r.Get("/{cep}", getCep)
	log.Fatal(http.ListenAndServe(":8080", r))
}
