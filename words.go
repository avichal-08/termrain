package main

import "strings"

var wordList = strings.Fields(`
go fmt nil map var int byte rune func chan err ctx api cli git bug fix log
test sync http json yaml port host ping tcp udp sql cron heap stack slice
defer panic select struct const range break build merge rebase commit push
pull fork clone deploy docker server client proxy cache queue mutex token
lambda vector tensor kernel socket buffer stream thread router schema
pointer channel package import return switch method string float error
module vendor binary compile linker runtime context timeout latency
goroutine interface waitgroup benchmark profiling allocator scheduler
middleware container kubernetes terraform microservice serializer
refactor deadlock checksum hashmap iterator generics closure recursion
algorithm database migration callback promise async await unicode
pipeline observability concurrency throughput bottleneck segfault
monorepo dependency changelog hotfix rollback staging production
`)

func wordsWithLen(min, max int) []string {
	var out []string
	for _, w := range wordList {
		if len(w) >= min && len(w) <= max {
			out = append(out, w)
		}
	}
	return out
}
