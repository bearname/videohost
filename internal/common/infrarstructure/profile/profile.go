package profile

import (
	"github.com/gorilla/mux"
	"net/http/pprof"
	"os"
	"runtime"
	runtimepprof "runtime/pprof"
)

func ProfileHandler(router *mux.Router) *mux.Router {
	pprofRouter := router.PathPrefix("/debug/pprof").Subrouter()
	pprofRouter.HandleFunc("/", pprof.Index)
	pprofRouter.HandleFunc("/cmdline", pprof.Cmdline)
	pprofRouter.HandleFunc("/symbol", pprof.Symbol)
	pprofRouter.HandleFunc("/trace", pprof.Trace)

	// Debug: pprof.WriteHeapProfile()
	memprofile := "/run/dnsd/memprofile"
	f, _ := os.Create(memprofile)
	runtime.GC() // get up-to-date statistics
	runtimepprof.WriteHeapProfile(f)
	f.Close()

	profile := pprofRouter.PathPrefix("/profile").Subrouter()
	profile.HandleFunc("", pprof.Profile)
	profile.Handle("/goroutine", pprof.Handler("goroutine"))
	profile.Handle("/threadcreate", pprof.Handler("threadcreate"))
	profile.Handle("/heap", pprof.Handler("heap"))
	profile.Handle("/block", pprof.Handler("block"))
	return router
}
