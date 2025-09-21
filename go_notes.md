in go there is one file in which we import package main, than file is the executable file

basic syntax for executable file 
```
package main

import {
	"fmt"
}

func main(){
}
```

go searches in this particular order for packages that is ->

-> Standard library (GOROOT)
-> Local project packages (module)
-> External modules (from Go proxy or GitHub)

Interfaces, slices, maps, channels → already act like references → no * needed for usual operations.

Structs and arrays → copied by value → use * pointer for efficiency or mutability.

Passing pointers → avoids copying and allows modification of original data.

defining struct in go-> 
```
type Person struct {
	Name string 
	Age int 
}
```

defining interface in go-> 
```
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct{
	Width, Height float64
}

func (r Rectangle) Area() float64{
	return r.Width * r.Height
}
```

ok so to call at any page i need to register with a mux, when i write 
http.HandleFunc("/hello", Hello)
it register Hello function with defaultmultipexer , i could also create a custom multiplexer
``` mux := http.NewServeMux()
    mux.HandleFunc("hello", Hello)
	http.ListenAndServe(":8080", mux)
```

this is how i tell what will the json field be 
```
	ID string `json:"id"`
```

this is how i set the header of response 
```
w.Header().Set("Content-Type", "application/json")
```

this is how i take params 
```
params = mux.Vars(r);
```

this is how i define a map
```
var m map[KeyType]ValueType
```

when calling function for [] slice i will pass &slicename because it will point to the header that's what i want 

in graphQL, we have two things one is query second is mutation, query is for get and mutation is for post, put delete update, in this we make a file schemaQL ( something like this ) and the package do the rest of the work

in grpc we need to install 

``` go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	winget install protobuf
```

and than first create file in something.proto in proto folder and run this to generate go files which will have all the functions created automatically
```
protoc --go_out=. --go-grpc_out=. proto/greet.proto
```






