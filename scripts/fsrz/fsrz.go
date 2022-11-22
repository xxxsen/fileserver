package main

import (
	"context"
	"fileserver/client"
	"flag"
	"fmt"
	"log"
)

var ak = flag.String("ak", "", "access key")
var sk = flag.String("sk", "", "secret key")
var schema = flag.String("schema", "http", "host schema, default http, can use https instead")
var domain = flag.String("domain", "", "host include schema")
var downloadDomain = flag.String("download_domain", "", "download domain, equals to host if empty")
var file = flag.String("file", "", "file to upload")

func main() {
	flag.Parse()
	c, err := client.New(
		client.WithHost(*schema+"://"+*domain),
		client.WithKey(*ak, *sk),
	)
	if err != nil {
		log.Fatalf("init client fail, err:%v", err)
	}
	ctx := context.Background()
	downkey, err := c.UploadFile(ctx, *file)
	if err != nil {
		log.Fatalf("upload file fail, err:%v", err)
	}
	host := *domain
	if len(*downloadDomain) > 0 {
		host = *downloadDomain
	}
	log.Printf("file link: %s", fmt.Sprintf("%s://%s/file?down_key=%s", *schema, host, downkey))
}
