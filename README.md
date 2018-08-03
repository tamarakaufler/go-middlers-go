# go-middlers-go
Alternative implementations of middleware in Go

## middlerware is not a typo :)
  - middler : one belonging to an intermediate group (Merriam-Webster)
  - middler is someone who will always care for you, be there for you and never betray you (Urban dictionary)

## Implementations

All three approaches define a Middler type, a function with the following signature:

```
  type Middler func(http.Handler) http.Handler
```

All implementations provide two off the shelf Middlers, one for logging, one for tracing.

Example of a logging Middler:

```
func LoggingMiddler(h http.Handler) http.Handler {
  // This is printed when the route is registered
  // h is either the route handler, or another Middler
  // wrapping either the route handler or another Middler.
  // -----------------------------------------------
	log.Println("Initialised LoggingMiddler1")

  // This handler is run when the route is requested
  // -----------------------------------------------
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Logging in LoggingMiddler before")
		log.Printf("%s", r.RequestURI)

		defer log.Println("Logging in LoggingMiddler after")
		h.ServeHTTP(w, r)
	})
}
```

## Implementation 1

  - middlerware/middler1.go
  - examples/example1/main.go

This approach provides a type Middlers:

```
  type Middlers []Middler
```

which has a method Apply. The method is responsible for merging the effect of all provided
middleware components - Middlers - and its output is an http.Handler, which, when a route
is called, runs all the nested middleware functionality. The Apply method call itself recursively.
The same set of Middlers is applied to all routes.

```
func (middlers Middlers) Apply(h http.Handler) http.Handler {
	if len(middlers) == 0 {
		return h
	}
	return middlers[1:].Apply(middlers[0](h))
}
```

Example os a routing setup:
```
	m1 := middler1.Middlers{middler1.TracingMiddler, middler1.LoggingMiddler}
	router.HandleFunc("/hi", hiHandler)
	router.HandleFunc("/bye", byeHandler)

	server := &http.Server{
		Addr:    ":7777",
		Handler: m1.Apply(router),
	}
```

## Implementation 2

  - middlerware/middler2.go
  - examples/example2/main.go

Ihe implementation uses an Apply function, which accepts as an argument

  - a route handler
  - a list of Middlers

```
func Apply(h http.Handler, middlers ...Middler) http.Handler {
	for _, middler := range middlers {
		h = middler(h)
	}
	return h
}
```
The routing can then be set up in the following way:

```
http.Handle("/hi", middler2.Apply(hiHandler,
                                  middler2.LoggingMiddler(logger), middler2.TracingMiddler()))
```

## Implementation 3
  
  - middlerware/middler3.go
  - examples/example3/main.go

The implementation introduces two new types:

```
type middlerware interface {
	Middlerware(handler http.Handler) http.Handler
}
```

and

```
type Middling struct {
	middlers []middleware
}
```

Users can use the built in Middler type or can introduce their own middleware type, which
needs to satisfy the middlerware interface, ie implement the Middleware method.

### Setting up middleware

a) Create a Middling for a group of routes where the same middleware needs to run:

The Middling Use method gathers one ot more Middlers, that should be applied together:

	hiMiddling := &middler3.Middling{}
	hiMiddling.Use(middler3.LoggingMiddler(logger))
	hiMiddling.Use(middler3.TracingMiddler())

	byeMiddling := &middler3.Middling{}
	byeMiddling.Use(middler3.TracingMiddler())

b) Applying the middlerware on routes:

  // applying a sequence of middlers
	http.Handle("/hi", hiMiddling.Apply(hiHandler))

	// applying just one middler
	http.Handle("/bye", byeMiddling.Apply(byeHandler))

An implementation of a custom middlerware (uthentication fot the /piggybank route)
is shown in middlerware.go.

# Related reading

https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702

https://drstearns.github.io/tutorials/gomiddleware/

https://hackernoon.com/simple-http-middleware-with-go-79a4ad62889b

https://gowebexamples.com/advanced-middleware/