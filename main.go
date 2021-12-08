package main

import (
	"flag"
	"fmt"
	"os"
	"paycheck/box"
	"paycheck/config"
	"paycheck/helpers"
	"paycheck/secure"
	"paycheck/svcloud"
	"sync"
)

type processStatus struct {
	success int
	fails   int
}

func removeTmpAssets(documents []svcloud.Document) {
	helpers.MessageProcessing(
		fmt.Sprintf("Removing %d files from temp folder", len(documents)),
	)

	var wg sync.WaitGroup

	ps := processStatus{
		success: 0,
		fails:   0,
	}

	wg.Add(len(documents))
	ch := make(chan (int), len(documents))

	go func() {
		for v := range ch {
			ps.success = ps.success + v
			ps.fails = len(documents) - ps.success
			wg.Done()
		}
	}()

	for _, document := range documents {
		go func(doc svcloud.Document) {
			var err error
			defer func() {
				if err != nil {
					ch <- 0
				} else {
					ch <- 1
				}
			}()

			filePath := "/tmp/" + doc.Name + ".gpg"
			err = os.Remove(filePath)
			if err != nil {
				helpers.MessageError(err)
			}
		}(document)
	}

	wg.Wait()

	helpers.MessageOK(fmt.Sprintf("Deleted %d files", ps.success))
	helpers.MessageNOK(fmt.Sprintf("%d files failed", ps.fails))
}

func downloadDocuments(documents []svcloud.Document, sv svcloud.ISvcloud, c config.IConfig) processStatus {
	var wg sync.WaitGroup

	ps := processStatus{
		success: 0,
		fails:   0,
	}

	wg.Add(len(documents))
	ch := make(chan (int), len(documents))

	go func() {
		for v := range ch {
			ps.success = ps.success + v
			ps.fails = len(documents) - ps.success
			wg.Done()
		}
	}()

	for _, document := range documents {
		go func(doc svcloud.Document) {
			var err error
			defer func() {
				if err != nil {
					ch <- 0
				} else {
					ch <- 1
				}
			}()

			filePath := c.GetRepo() + "/" + doc.Name

			// Download the document from server
			err = sv.DownloadDocument(doc, filePath)
			if err != nil {
				helpers.MessageError(err)
				return
			}
		}(document)
	}

	wg.Wait()

	return ps
}

func encryptDocuments(documents []svcloud.Document, sv svcloud.ISvcloud, c config.IConfig) processStatus {
	s := secure.New()

	var wg sync.WaitGroup

	ps := processStatus{
		success: 0,
		fails:   0,
	}

	wg.Add(len(documents))
	ch := make(chan (int), len(documents))

	go func() {
		for v := range ch {
			ps.success = ps.success + v
			ps.fails = len(documents) - ps.success
			wg.Done()
		}
	}()

	for _, document := range documents {
		go func(doc svcloud.Document) {
			var err error
			defer func() {
				if err != nil {
					ch <- 0
				} else {
					ch <- 1
				}
			}()

			filePath := c.GetRepo() + "/" + doc.Name
			fileGPGName := "/tmp/" + doc.Name + ".gpg"

			file, err := os.Open(filePath)
			if err != nil {
				return
			}
			defer file.Close()

			fileEncrypted, err := s.Encrypt(fileGPGName, file, c.GetPubKey())
			if err != nil {
				helpers.MessageError(err)
				return
			}
			defer fileEncrypted.Close()
		}(document)
	}

	wg.Wait()

	return ps
}

func uploadDocuments(documents []svcloud.Document, sv svcloud.ISvcloud, c config.IConfig) processStatus {
	b := box.New(c.GetBoxToken())

	var wg sync.WaitGroup

	ps := processStatus{
		success: 0,
		fails:   0,
	}

	wg.Add(len(documents))
	ch := make(chan (int), len(documents))

	go func() {
		for v := range ch {
			ps.success = ps.success + v
			ps.fails = len(documents) - ps.success
			wg.Done()
		}
	}()

	for _, document := range documents {
		go func(doc svcloud.Document) {
			var err error
			defer func() {
				if err != nil {
					ch <- 0
				} else {
					ch <- 1
				}
			}()

			fileGPGName := "/tmp/" + doc.Name + ".gpg"

			file, err := os.Open(fileGPGName)
			if err != nil {
				return
			}
			defer file.Close()

			err = b.Upload(file, c.GetBoxTargetID())
			// helpers.Message(fmt.Sprintf("File %s uploaded to BOX", doc.Name), err)
			if err != nil {
				return
			}
		}(document)
	}

	wg.Wait()

	return ps
}

func main() {
	var withEncrypt bool
	var withUpload bool

	flag.BoolVar(&withEncrypt, "c", false, "Create an encrypted version of the zip")
	flag.BoolVar(&withUpload, "u", false, "Upload the encrypted file to BOX")

	flag.Parse()

	c, err := config.New()
	if err != nil {
		helpers.MessageError(err)
		return
	}

	sv := svcloud.New(c.GetSvCloudBaseURL(), c.GetSvCloudEndpoint())

	// Login into paycheck platform
	_, err = sv.Login(c.GetSvCloudUsername(), c.GetSvCloudPassword())
	helpers.Message("Login to "+c.GetSvCloudBaseURL(), err)
	if err != nil {
		return
	}

	// Retrieve all documents
	documents, err := sv.ListDocuments()
	helpers.Message("Get documents to download", err)
	if err != nil {
		return
	}

	// Download documents
	helpers.MessageProcessing(
		fmt.Sprintf(
			"Start downloading %d documents",
			len(documents),
		),
	)

	ps := downloadDocuments(documents, sv, c)
	helpers.MessageOK(fmt.Sprintf("Downloaded %d files", ps.success))
	helpers.MessageNOK(fmt.Sprintf("%d files failed", ps.fails))

	// If enable encryption
	if withEncrypt {
		helpers.MessageProcessing(
			fmt.Sprintf(
				"Start encrypting %d documents",
				len(documents),
			),
		)

		ps := encryptDocuments(documents, sv, c)
		helpers.MessageOK(fmt.Sprintf("Encrypted %d files", ps.success))
		helpers.MessageNOK(fmt.Sprintf("%d files failed", ps.fails))
	}

	// If enable upload + encryption
	if withUpload && withEncrypt {
		defer removeTmpAssets(documents)
		helpers.MessageProcessing(
			fmt.Sprintf(
				"Start uploading %d documents",
				len(documents),
			),
		)

		ps := uploadDocuments(documents, sv, c)
		helpers.MessageOK(fmt.Sprintf("Uploaded %d files", ps.success))
		helpers.MessageNOK(fmt.Sprintf("%d files failed", ps.fails))
	}
}
