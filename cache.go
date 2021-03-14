package main

import (
	"github.com/coocood/freecache"
	"log"
)

var cache = freecache.NewCache(1 * 1024 * 1024)

func getCache(k []byte) []byte {
	var got []byte
	got, err := cache.Get(k)
	if err != nil {
		log.Print("can't get from cache")
	}

	return got
}

func setCache(k []byte, v []byte)  {
	err := cache.Set(k, v, 60)
	if err != nil {
		log.Print("can't save to cache")

		return
	}
}
