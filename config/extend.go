package config

var ExtConfig Extend

type Extend struct {
	AMap AMap
}

type AMap struct {
	Key string
}