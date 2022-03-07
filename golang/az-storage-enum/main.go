package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func main() {
	// Your account name and key can be obtained from the Azure Portal.
	// Ex: "newsta321tierht"
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("Missing required env-var: AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	accountKey, ok := os.LookupEnv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
	if !ok {
		panic("Missing required env-var: AZURE_STORAGE_PRIMARY_ACCOUNT_KEY could not be found")
	}
	// target_container, ok := os.LookupEnv("AZURE_STORAGE_TARGET_CONTAINER")
	// if !ok {
	// 	panic("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY could not be found")
	// }

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	// The service URL for blob endpoints is usually in the form: http(s)://<account>.blob.core.windows.net/
	service, err := azblob.NewServiceClientWithSharedKey(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	// // ===== 1. Create a container =====

	// // First, create a container client, and use the Create method to create a new container in your account
	//  container := service.NewContainerClient("mycontainer")
	// // All functions that make service requests have an options struct as the final parameter.
	// // The options struct allows you to specify optional parameters such as metadata, public access types, etc.
	// // If you want to use the default options, pass in nil.
	// _, err = container.Create(context.TODO(), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // ===== 2. Upload and Download a block blob =====
	// data := "Hello world!"

	// // Create a new BlockBlobClient from the ContainerClient
	// blockBlob := container.NewBlockBlobClient("HelloWorld.txt")

	// // Upload data to the block blob
	// _, err = blockBlob.Upload(context.TODO(), streaming.NopCloser(strings.NewReader(data)), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Download the blob's contents and ensure that the download worked properly
	// get, err := blockBlob.Download(context.TODO(), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Use the bytes.Buffer object to read the downloaded data.
	// downloadedData := &bytes.Buffer{}
	// reader := get.Body(nil) // RetryReaderOptions has a lot of in-depth tuning abilities, but for the sake of simplicity, we'll omit those here.
	// _, err = downloadedData.ReadFrom(reader)
	// if err != nil {
	// 	return
	// }

	// err = reader.Close()
	// if err != nil {
	// 	return
	// }
	// if data != downloadedData.String() {
	// 	log.Fatal("downloaded data does not match uploaded data")
	// } else {
	// 	fmt.Printf("Downloaded data: %s\n", downloadedData.String())
	// }

	// ===== 3. List blobs =====
	// List methods returns a pager object which can be used to iterate over the results of a paging operation.
	// To iterate over a page use the NextPage(context.Context) to fetch the next page of results.
	// PageResponse() can be used to iterate over the results of the specific page.
	// Always check the Err() method after paging to see if an error was returned by the pager. A pager will return either an error or the page of results.

	// Iterate containers
	pagerC := service.ListContainers(nil)
	for pagerC.NextPage(context.TODO()) {
		respC := pagerC.PageResponse()

		// for _, v := range respC.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
		for _, a_containerItem := range respC.ContainerItems {
			fmt.Println("\n==Container: ", *a_containerItem.Name)

			a_containerClient := service.NewContainerClient(*a_containerItem.Name)
			// Iterate blobs
			pager := a_containerClient.ListBlobsFlat(nil)
			for pager.NextPage(context.TODO()) {
				resp := pager.PageResponse()

				for _, a_blobItem := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
					fmt.Println("....Blob: ", *a_blobItem.Name)

					blockBlob := a_containerClient.NewBlockBlobClient(*a_blobItem.Name)

					// Download the blob's contents and ensure that the download worked properly
					get, err := blockBlob.Download(context.TODO(), nil)
					if err != nil {
						log.Fatal(err)
					}

					// Open a buffer, reader, and then download!
					downloadedData := &bytes.Buffer{}
					reader := get.Body(nil) // RetryReaderOptions has a lot of in-depth tuning abilities, but for the sake of simplicity, we'll omit those here.
					_, err = downloadedData.ReadFrom(reader)
					if err != nil {
						log.Fatal(err)
					}
					err = reader.Close()
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("-- blob data start --\n", downloadedData.String(), "\n++ blob data end ++\n")

				}
			}

			if err = pager.Err(); err != nil {
				log.Fatal(err)
			}

		}
	}

	if err = pagerC.Err(); err != nil {
		log.Fatal(err)
	}
	// sleep 1h for pod to be visible
	time.Sleep(3600 * time.Second)

}
