package main

import (
	"fmt"
	"log"
	"os"
	"flag"
	"../../src/dependents/github.com/sideshow/certificate"
	"../../src/dependents/github.com/sideshow"
	"../../src/dependents/github.com/sideshow/payload"
)

func main() {
	certPath := flag.String("cert", "jxyr-push-prod-Certificates.p12", "Path to .p12 certificate file (Required)")
	password := flag.String("password", "1", "Path to .p12 certificate password (Required)")
	count := flag.Int("count", 1, "Number of pushes to send")
	token := flag.String("token", "39c7ee3bb48381d853a89b89b487e4285e1fae4050d5b2d1d9da543c59b26eaf", "Push token (Required)")
	topic := flag.String("topic", "com.jxyr.ygmobile", "Topic (Required)")
	flag.Parse()

	if *certPath == "" || *token == "" || *topic == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	cert, err := certificate.FromP12File(*certPath, *password)
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	notifications := make(chan *apns2.Notification, 100)
	responses := make(chan *apns2.Response, *count)

	client := apns2.NewClient(cert).Production()

	for i := 0; i < 50; i++ {
		go worker(client, notifications, responses)
	}

	for i := 0; i < *count; i++ {
		n := &apns2.Notification{
			DeviceToken: *token,
			Topic:       *topic,
			Payload:     payload.NewPayload().Alert(fmt.Sprintf("Hello! %v", i)),
		}
		notifications <- n
	}

	for i := 0; i < *count; i++ {
		res := <-responses
		fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	}

	close(notifications)
	close(responses)
}

func worker(client *apns2.Client, notifications <-chan *apns2.Notification, responses chan<- *apns2.Response) {
	for n := range notifications {
		res, err := client.Push(n)
		if err != nil {
			log.Fatal("Push Error:", err)
		}
		responses <- res
	}
}
