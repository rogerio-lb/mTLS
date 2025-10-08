package main

import (
	"context"
	"fmt"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
)

func main() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		//URL: "pulsar://localhost:6650",
		URL:            "pulsar+ssl://pc-a5bec094.aws-use2-production-snci-pool-kid.streamnative.aws.snio.cloud",
		Authentication: pulsar.NewAuthenticationToken("eyJhbGciOiJSUzI1NiIsImtpZCI6IjE0NjNhODQ5LTNkNzUtNTlmMi1hMTgyLTVjNzE0ODY4YjBhMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsidXJuOnNuOnB1bHNhcjpvLWU1NWNwOmxiLXBheSJdLCJodHRwczovL3N0cmVhbW5hdGl2ZS5pby9zY29wZSI6WyJhZG1pbiIsImFjY2VzcyJdLCJodHRwczovL3N0cmVhbW5hdGl2ZS5pby91c2VybmFtZSI6ImxiLXN0Z0BvLWU1NWNwLmF1dGguc3RyZWFtbmF0aXZlLmNsb3VkIiwiaWF0IjoxNzUzMTMxMTcxLCJpc3MiOiJodHRwczovL3BjLWE1YmVjMDk0LmF3cy11c2UyLXByb2R1Y3Rpb24tc25jaS1wb29sLWtpZC5zdHJlYW1uYXRpdmUuYXdzLnNuaW8uY2xvdWQvYXBpa2V5cy8iLCJqdGkiOiJhYTk3YTVhN2YwNzA0N2FjYTI3MzQ4ODdlOTI1ZDMyMyIsInBlcm1pc3Npb25zIjpbXSwic3ViIjoiVWNDbGVyaENVRFh6S1NUZ21WbHFPVkF5b1R0aDlUT1lAY2xpZW50cyJ9.UfLKDZNusPHya-xgdWHoSNXbp6nhBMaEyizzULkWCsriY4VKdfkJ6OqnrPXP9xOi0aVKzCL-9ObgzxBklKoFObguZJ1MrIgzeiQfp0FUfmylwWz_jb-zbPZ5cbclvrbMXojKJte1lk9GxmmggBf-zUpuRGDiVGV42ZnU2AVJ-1PXx5frQ5SUbfJkfIRDp566b6PoF9r80gYc594CCo_Z0nUUMjHR_1molD5BDBYoK3O71yy-kEf-_J_nOdfMBQdHZbQGitpo5BzLvE-kdpHg0JZ392IZPeWhoZCEyGfLaNp6aNi8tyCMe--NQFm78h6bQ4L3VHge7BVR7dOjMRiBQQ"),
	})
	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}
	defer client.Close()

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:                       "persistent://lb-core/spi/psti-to-bridge-topic",
		SubscriptionName:            "rsfn-connect-spi-v2-psti-subscription",
		SubscriptionInitialPosition: pulsar.SubscriptionPositionEarliest,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	counter := 0

	for {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		/*fmt.Printf("Received message msgId: %v -- content: '%s'\n",
		msg.ID(), string(msg.Payload()))*/

		counter++

		fmt.Println("counter:", counter)

		consumer.Ack(msg)
	}
}
