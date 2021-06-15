@echo off
curl -sK -v http://localhost:%1/debug/pprof/profile/%2 > %3.out
go tool pprof %3.out