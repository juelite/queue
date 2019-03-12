package structs

import "net/url"

type RedisQueue struct {
	Method 		string
	Params 		url.Values
}
