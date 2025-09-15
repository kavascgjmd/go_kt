package main
import (
	"fmt"
	"log"
	"net/http"
)

func formHandler(w http.ResponseWriter, r * http.Request){
  fmt.Fprintf(w, "hi")
  return 
}

func helloHandler(w http.ResponseWriter, r * http.Request){
   fmt.Fprintf(w, "bi");
   return
}

func main(){
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)
	fmt.Printf("starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err);
	}
}