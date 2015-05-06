package lib

import "fmt"

func DownloadImages(images []string) (int, bool) {
	//TODO: use a combination of generator pattern and select to create n go routines and synchronize them using select.
	//The core download method would just take a string url, and download it to $baseDownloadpath/tag
	//the advantage is that when index creation is done it can directly use the tag.
	//
	var maxGoRoutines int
	if len(images) > 50 {
		maxGoRoutines = 50
	} else {
		maxGoRoutines = len(images)
	}
	var downloadResult []string
	tasks := make(chan string, maxGoRoutines)
	results := make(chan string, len(images))
	for j := 1; j <= maxGoRoutines; j++ {
		go downloader(tasks, results)
	}

	for _, link := range images {
		tasks <- link
	}

	for i := 0; i < len(images); i++ {
		select {
		case result := <-results:
			downloadResult = append(downloadResult, result)
		}
	}
	return len(downloadResult), true
}

func downloader(links <-chan string, results chan<- string) {
	for link := range links {
		//process and put the result in results
		//TODO: Send a get request and download the file. See what built in options go has for image downloading.
		s := fmt.Sprintf(link)
		results <- s
	}
}
