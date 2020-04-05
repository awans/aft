package server

import (
	"log"
	"net/http"
	"time"
)

func Middleware(op Operation) http.Handler {
	// this goes inside out
	// invoke the server
	server := op.Server

	// and audit log it
	auditLoggedServer := AuditLog(op, server)

	// wrap it as a handler
	handler := ServerToHandler(auditLoggedServer)

	// before that we log the request
	logged := Logger(handler, op.Name)

	// before that we set CORS
	cors := CORS(logged)

	return cors
}

type AuditLoggedServer struct {
	inner Server
}

func (a AuditLoggedServer) Parse(req *http.Request) interface{} {
	return a.inner.Parse(req)
}

func (a AuditLoggedServer) Serve(w http.ResponseWriter, req interface{}) {
	a.inner.Serve(w, req)
}

func AuditLog(op Operation, inner Server) Server {
	return AuditLoggedServer{inner: inner}
}

func ServerToHandler(server Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsed := server.Parse(r)
		server.Serve(w, parsed)
	})
}

func CORS(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		inner.ServeHTTP(w, r)
	})
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}