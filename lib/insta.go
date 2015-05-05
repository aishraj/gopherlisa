package lib

import "encoding/json"

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

type APIResponse struct {
	Pagination PaginationObject
	Meta       MetaData
	Data       []TagResponse
}

type PaginationObject struct {
	NextURL            string
	NextMaxID          string
	DeprecationWarning string
	NextMaxTagID       string
	NextMinID          string
	MinTagID           string
}

type MetaData struct {
	ErrorType    string
	Code         int32
	ErrorMessage string
}

type TagResponse struct {
	Attribution  json.RawMessage
	Videos       json.RawMessage
	Tags         []string
	MediaType    string
	Location     json.RawMessage
	Comments     json.RawMessage
	Filter       json.RawMessage
	CreatedTime  string
	Link         string
	Images       DigitialMediaInfo
	Likes        json.RawMessage
	UsersInPhoto json.RawMessage
	Caption      json.RawMessage
	UserLinked   json.RawMessage
	ID           string
	UserInfo     json.RawMessage
}

type DigitialMediaInfo struct {
	LowResolution      ImageDetails
	Thumbnail          ImageDetails
	StandardResolution ImageDetails
}

type ImageDetails struct {
	URL    string
	Height int32
	Width  int32
}
