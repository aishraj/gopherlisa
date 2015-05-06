package lib

func DownloadImages(images []string) (string, bool) {
	//TODO: use a comination of generator pattern and select to create n go routines and synchronize them using select.
	//The core download method would just take a string url, and download it to $baseDownloadpath/tag
	//the advantage is that when index creation is done it can directly use the tag.
	//
}
