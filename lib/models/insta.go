package models

//ImageMetaData is a wrapper around the instagram response
type ImageMetaData struct {
	ThumbnailURL string
	Height       int
	Width        int
	MediaID      string
}

type AuthenticationUser struct {
	ID             string `json:"id"`
	UserName       string `json:"user_name"`
	FullName       string `json:"full_name"`
	ProfilePicture string `json:"profile_picture"`
}

type AuthenticationResponse struct {
	AccessToken string             `json:"access_token"`
	User        AuthenticationUser `json:"user"`
}
