package civogo_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/civo/civogo/v3"
	"github.com/civo/civogo/v3/compute"
)

// List every instance across two regions; All fetches further pages lazily.
func Example() {
	client, err := civogo.New(os.Getenv("CIVO_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	for inst, err := range client.Compute.Instances.All(context.Background(), compute.ListInstancesRequest{},
		civogo.WithRegion("lon1"), civogo.WithRegion("fra1")) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(inst.Region, inst.Hostname, inst.Status)
	}
}

func ExampleNew() {
	client, err := civogo.New(os.Getenv("CIVO_TOKEN"), civogo.WithDefaultRegion("lon1"))
	if err != nil {
		log.Fatal(err)
	}

	inst, err := client.Compute.Instances.Get(context.Background(), "my-instance-id")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(inst.Hostname)
}

// Branch on the error's reason, and read field violations on validation failures.
func ExampleError() {
	client, err := civogo.New(os.Getenv("CIVO_TOKEN"), civogo.WithDefaultRegion("lon1"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Compute.Instances.Get(context.Background(), "my-instance-id")
	switch {
	case err == nil:
	case civogo.HasReason(err, compute.ReasonInstanceNotFound):
		fmt.Println("instance is gone")
	case civogo.HasReason(err, compute.ReasonRequestValidationFailed):
		var apiErr *civogo.Error
		errors.As(err, &apiErr)
		for _, v := range apiErr.FieldViolations {
			fmt.Printf("%s: %s\n", v.Field, v.Description)
		}
	default:
		log.Fatal(err)
	}
}

// Optional fields are pointers: nil leaves them unchanged, Ptr sends a value.
func ExamplePtr() {
	client, err := civogo.New(os.Getenv("CIVO_TOKEN"), civogo.WithDefaultRegion("lon1"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Compute.Instances.Update(context.Background(), "my-instance-id", compute.UpdateInstanceRequest{
		PublicIPRequired: civogo.Ptr(false),
		Notes:            civogo.Ptr(""),
	})
	if err != nil {
		log.Fatal(err)
	}
}
