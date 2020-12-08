package api

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type EndPoint string

const Version = "0.0.1-Alpha0"
const UserAgent = "CloudinaryGo/" + Version

var BaseUrl = "https://api.cloudinary.com/v1_1"

type AssetType string

func (a AssetType) ToString() string {
	if a == "" {
		a = Image
	}
	return string(a)
}

const (
	Image AssetType = "image"
	Video           = "video"
	File            = "raw"
	Auto            = "auto"
	All             = "all"
)

type DeliveryType string

func (d DeliveryType) ToString() string {
	if d == "" {
		d = Upload
	}
	return string(d)
}

const (
	Upload          DeliveryType = "upload"
	Private                      = "private"
	Public                       = "public"
	Authenticated                = "authenticated"
	Fetch                        = "fetch"
	Sprite                       = "sprite"
	Text                         = "text"
	Multi                        = "multi"
	Facebook                     = "facebook"
	Twitter                      = "twitter"
	TwitterName                  = "twitter_name"
	Gravatar                     = "gravatar"
	Youtube                      = "youtube"
	Hulu                         = "hulu"
	Vimeo                        = "vimeo"
	Animoto                      = "animoto"
	Worldstarhiphop              = "worldstarhiphop"
	Dailymotion                  = "dailymotion"
)

type ModerationStatus string

const (
	Pending  ModerationStatus = "pending"
	Approved                  = "approved"
	Rejected                  = "rejected"
)

// Option is the optional parameters custom struct
type Option map[string]interface{}

type Coordinates [][]int
type CldApiArray []string

type CldApiMap map[string]string
type Metadata map[string]interface{}

type BriefAssetResult struct {
	AssetID     string    `json:"asset_id"`
	PublicID    string    `json:"public_id"`
	Format      string    `json:"format"`
	Version     int       `json:"version"`
	AssetType   string    `json:"resource_type"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	Bytes       int       `json:"bytes"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Backup      bool      `json:"backup"`
	AccessMode  string    `json:"access_mode"`
	URL         string    `json:"url"`
	SecureURL   string    `json:"secure_url"`
	Tags        []string  `json:"tags,omitempty"`
	Context     CldApiMap `json:"context,omitempty"`
	Metadata    Metadata  `json:"metadata,omitempty"`
	Placeholder bool      `json:"placeholder,omitempty"`
	Error       string    `json:"error,omitempty"`
}

// MarshalJSON writes a quoted string in the custom format
func (cldApiMap CldApiMap) MarshalJSON() ([]byte, error) {
	// FIXME: handle escaping
	var params []string
	for name, value := range cldApiMap {
		params = append(params, strings.Join([]string{name, value}, "="))
	}

	return []byte(strconv.Quote(strings.Join(params, "|"))), nil
}

// MarshalJSON writes a quoted string in the custom format
func (cldApiArr CldApiArray) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strings.Join(cldApiArr[:], ","))), nil
}

// ErrorResp is the failed api request main struct
type ErrorResp struct {
	Message string `json:"message"`
}

func BuildPath(parts ...interface{}) string {
	var partsSlice []string
	for _, part := range parts {
		if part != "" {
			partsSlice = append(partsSlice, fmt.Sprintf("%v", part))
		}
	}

	return strings.Join(partsSlice, "/")
}

func SignRequest(params url.Values, secret string) (string, error) {
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))

	encodedUnescapedParams, err := url.QueryUnescape(params.Encode())
	if err != nil {
		return "", err
	}

	hash := sha1.New()
	hash.Write([]byte(encodedUnescapedParams + secret))

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func StructToParams(inputStruct interface{}) (url.Values, error) {
	var paramsMap map[string]interface{}
	paramsJsonObj, _ := json.Marshal(inputStruct)
	err := json.Unmarshal(paramsJsonObj, &paramsMap)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	for paramName, value := range paramsMap {
		resBytes, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}

		res := string(resBytes)
		if strings.HasPrefix(res, "\"") { // FIXME: Fix this dirty hack that prevents double quoting of strings
			res, _ = strconv.Unquote(string(res))
		}

		params.Add(paramName, res)
	}

	return params, nil
}

func DeferredClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}
