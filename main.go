package main

import (
	"github.com/gusanmaz/Synaesthesia_Test/datastore"
	"github.com/gusanmaz/Synaesthesia_Test/server"
	"log"
	"net/http"
	"os"
)

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		// see @andreiavrammsd comment: often 307 > 301
		http.StatusTemporaryRedirect)
}

func main() {
	var mongo datastore.Mongo
	mongo.New("127.0.0.1:27017")
	defer mongo.Close()
	http.Handle("/", server.Router(mongo))

	sslCert := os.Getenv("SSL_CERT")
	sslKey  := os.Getenv("SSL_KEY")

	if sslCert == "" || sslKey == "" {
		http.ListenAndServe(":8087", nil)

	}else{
		err := http.ListenAndServeTLS(":443", sslCert, sslKey, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
		//go http.ListenAndServe(":80", http.HandlerFunc(redirect))

		go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
		}))

	}
}
