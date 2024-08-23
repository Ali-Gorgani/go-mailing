package kavenegar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	SEND            = "sms/send.json"
	SENDARRAY       = "sms/sendarray.json"
	STATUS          = "sms/status.json"
	STATUSBYLOCALID = "sms/statuslocalmessageid.json"
	SELECT          = "sms/select.json"
	SELECTOUTBOX    = "sms/selectoutbox.json"
	LATESTOUTBOX    = "sms/latestoutbox.json"
	COUNTOUTBOX     = "sms/countoutbox.json"
	CANCEL          = "sms/cancel.json"
	RECEIVE         = "sms/receive.json"
	COUNTINBOX      = "sms/countinbox.json"
	LOOKUP          = "verify/lookup.json"
	TTS             = "call/maketts.json"
	INFO            = "account/info.json"
	CONFIG          = "account/config.json"
)

type Kavenegar struct {
	apiKey string
}

func New(apiKey string) *Kavenegar {
	return &Kavenegar{apiKey: apiKey}
}

type SendInputParams struct {
	Receptor []string  `json:"receptor"`
	Message  string    `json:"message"`
	Sender   string    `json:"sender,omitempty"`  // Optional
	Date     time.Time `json:"date,omitempty"`    // Optional
	Type     string    `json:"type,omitempty"`    // Optional
	LocalID  int64     `json:"localid,omitempty"` // Optional
	Hide     byte      `json:"hide,omitempty"`    // Optional
}

// Output represents the structure of the API response
type OutputParams struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
	Entries []struct {
		MessageID  int64  `json:"messageid"`
		Message    string `json:"message"`
		Status     int    `json:"status"`
		StatusText string `json:"statustext"`
		Sender     string `json:"sender"`
		Receptor   string `json:"receptor"`
		Date       int64  `json:"date"`
		Cost       int    `json:"cost"`
	} `json:"entries"`
}

func (k *Kavenegar) CreateURL(method string) string {
	return fmt.Sprintf("https://api.kavenegar.com/v1/%s/%s", k.apiKey, method)
}

func (k *Kavenegar) Send(p SendInputParams) (OutputParams, error) {
	// Base URL
	baseURL := k.CreateURL(SEND)

	// Combine receptor into a comma-separated string
	var receptorStrings []string
	receptorStrings = append(receptorStrings, p.Receptor...)
	receptorParam := strings.Join(receptorStrings, ",")

	// Build query parameters
	params := url.Values{}
	params.Add("receptor", receptorParam)
	params.Add("message", p.Message)
	if p.Sender != "" {
		params.Add("sender", p.Sender)
	}
	if !p.Date.IsZero() {
		params.Add("date", fmt.Sprintf("%d", p.Date.Unix()))
	}
	if p.Type != "" {
		params.Add("type", p.Type)
	}
	if p.LocalID != 0 {
		params.Add("localid", fmt.Sprintf("%d", p.LocalID))
	}
	if p.Hide != 0 {
		params.Add("hide", fmt.Sprintf("%d", p.Hide))
	}

	// Construct the full URL
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response JSON
	var output OutputParams
	if err := json.Unmarshal(body, &output); err != nil {
		return OutputParams{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

type SendArrayInputParams struct {
	Receptor        []string  `json:"receptor"`
	Sender          []string  `json:"sender"`
	Message         []string  `json:"message"`
	Date            time.Time `json:"date,omitempty"`            // Optional
	Type            []int     `json:"type,omitempty"`            // Optional
	LocalMessageIDs []int64   `json:"localmessageids,omitempty"` // Optional
	Hide            byte      `json:"hide,omitempty"`            // Optional
}

func (k *Kavenegar) SendArray(p SendArrayInputParams) (OutputParams, error) {
	// URL
	baseURL := k.CreateURL(SENDARRAY)

	// validation
	if len(p.Receptor) > 200 {
		return OutputParams{}, fmt.Errorf("maximum number of receptors is 200")
	}
	if len(p.Receptor) != len(p.Message) || len(p.Receptor) != len(p.Sender) {
		return OutputParams{}, fmt.Errorf("length of receptor, message, and sender must be equal")
	}
	if len(p.Type) > 0 && len(p.Type) != len(p.Receptor) {
		return OutputParams{}, fmt.Errorf("length of type must be equal to the length of receptor")
	}
	if len(p.LocalMessageIDs) > 0 && len(p.LocalMessageIDs) != len(p.Receptor) {
		return OutputParams{}, fmt.Errorf("length of localmessageids must be equal to the length of receptor")
	}

	// Build query parameters
	params := url.Values{}
	for i := 0; i < len(p.Receptor); i++ {
		params.Add("receptor", p.Receptor[i])
		params.Add("message", p.Message[i])
		params.Add("sender", p.Sender[i])
		if len(p.Type) > 0 {
			params.Add("type", fmt.Sprintf("%d", p.Type[i]))
		}
		if len(p.LocalMessageIDs) > 0 {
			params.Add("localmessageids", fmt.Sprintf("%d", p.LocalMessageIDs[i]))
		}
	}
	if !p.Date.IsZero() {
		params.Add("date", fmt.Sprintf("%d", p.Date.Unix()))
	}
	if p.Hide != 0 {
		params.Add("hide", fmt.Sprintf("%d", p.Hide))
	}

	// Convert input parameters to JSON
	jsonData, err := json.Marshal(params)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to marshal input parameters: %v", err)
	}
	jsonData = bytes.Trim(jsonData, "{}")

	// Create a new HTTP request
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	fmt.Println(req)

	// // Manually build the request body
	// var builder strings.Builder

	// // Helper function to format slices into the required format
	// formatParam := func(key string, values []string) string {
	// 	return fmt.Sprintf("%s=[%s]", key, strings.Join(values, ","))
	// }

	// builder.WriteString(formatParam("receptor", p.Receptor))
	// builder.WriteString(" ")
	// builder.WriteString(formatParam("sender", p.Sender))
	// builder.WriteString(" ")
	// builder.WriteString(formatParam("message", p.Message))

	// // Convert builder content to string
	// nbody := builder.String()
	// fmt.Println(nbody)

	// // Create a new HTTP request
	// req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer([]byte(nbody)))
	// if err != nil {
	// 	return OutputParams{}, fmt.Errorf("failed to create HTTP request: %v", err)
	// }
	// fmt.Println(req.Body)

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response JSON
	var output OutputParams
	if err := json.Unmarshal(body, &output); err != nil {
		return OutputParams{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

type StatusResponse struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
	Entries []struct {
		MessageID  int64  `json:"messageid"`
		Status     int    `json:"status"`
		StatusText string `json:"statustext"`
	} `json:"entries"`
}

func (k *Kavenegar) Status(messageIDs []int64) (StatusResponse, error) {
	baseURL := k.CreateURL(STATUS)

	// Combine message IDs into a comma-separated string
	var messageIDStrings []string
	for _, messageID := range messageIDs {
		messageIDStrings = append(messageIDStrings, fmt.Sprintf("%d", messageID))
	}
	messageIDParam := strings.Join(messageIDStrings, ",")

	params := url.Values{}
	params.Add("messageid", messageIDParam)

	// Construct the full URL with the combined message IDs
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return StatusResponse{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return StatusResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response JSON
	var output StatusResponse
	if err := json.Unmarshal(body, &output); err != nil {
		return StatusResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

type StatusByLocalIDInputResponse struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
	Entries []struct {
		MessageID  int64  `json:"messageid"`
		LocalID    int64  `json:"localid"`
		Status     int    `json:"status"`
		StatusText string `json:"statustext"`
	} `json:"entries"`
}

func (k *Kavenegar) StatusByLocalid(localIDs []int64) (StatusByLocalIDInputResponse, error) {
	baseURL := k.CreateURL(STATUSBYLOCALID)

	// Combine message IDs into a comma-separated string
	var localIDStrings []string
	for _, localID := range localIDs {
		localIDStrings = append(localIDStrings, fmt.Sprintf("%d", localID))
	}
	localIDParam := strings.Join(localIDStrings, ",")

	params := url.Values{}
	params.Add("localid", localIDParam)

	// Construct the full URL with the combined message IDs
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return StatusByLocalIDInputResponse{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return StatusByLocalIDInputResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response JSON
	var output StatusByLocalIDInputResponse
	if err := json.Unmarshal(body, &output); err != nil {
		return StatusByLocalIDInputResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

func (k *Kavenegar) Select(messageIDs []int64) (OutputParams, error) {
	// Base URL
	baseURL := k.CreateURL(SELECT)

	// Combine message IDs into a comma-separated string
	var messageIDStrings []string
	for _, messageID := range messageIDs {
		messageIDStrings = append(messageIDStrings, fmt.Sprintf("%d", messageID))
	}
	messageIDParam := strings.Join(messageIDStrings, ",")

	params := url.Values{}
	params.Add("messageid", messageIDParam)

	// Construct the full URL with the combined message IDs
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response JSON
	var output OutputParams
	if err := json.Unmarshal(body, &output); err != nil {
		return OutputParams{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

type SelectOutboxInputParams struct {
	StartDate int64  `json:"startdate"`
	EndDate   int64  `json:"enddate,omitempty"` // Optional
	Sender    string `json:"sender,omitempty"`  // Optional
}

func (k *Kavenegar) SelectOutbox(p SelectOutboxInputParams) (OutputParams, error) {
	// Base URL
	baseURL := k.CreateURL(SELECTOUTBOX)

	params := url.Values{}
	params.Add("startdate", fmt.Sprintf("%d", p.StartDate))
	if p.EndDate != 0 {
		params.Add("enddate", fmt.Sprintf("%d", p.EndDate))
	}
	if p.Sender != "" {
		params.Add("sender", p.Sender)
	}

	// Construct the full URL with the combined message IDs
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response JSON
	var output OutputParams
	if err := json.Unmarshal(body, &output); err != nil {
		return OutputParams{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

type LatestOutboxInputParams struct {
	PageSize int64  `json:"pagesize,omitempty"` // Optional
	Sender   string `json:"sender,omitempty"`   // Optional
}

func (k *Kavenegar) LatestOutBox(p LatestOutboxInputParams) (OutputParams, error) {
	// Base URL
	baseURL := k.CreateURL(LATESTOUTBOX)

	params := url.Values{}
	if p.PageSize != 0 {
		params.Add("pagesize", fmt.Sprintf("%d", p.PageSize))
	}
	if p.Sender != "" {
		params.Add("sender", p.Sender)
	}

	// Construct the full URL with the combined message IDs
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response JSON
	var output OutputParams
	if err := json.Unmarshal(body, &output); err != nil {
		return OutputParams{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

type CountOutboxInputParams struct {
	StartDate int64  `json:"startdate"`
	EndDate   int64  `json:"enddate,omitempty"` // Optional
	Sender    string `json:"sender,omitempty"`  // Optional
}

type CountOutboxResponse struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
	Entries struct {
		StartDate int64 `json:"startdate"`
		EndDate   int64 `json:"enddate"`
		SumPart   int64 `json:"sumpart"`
		SumCount  int64 `json:"sumcount"`
		Cost      int64 `json:"cost"`
	} `json:"entries"`
}

func (k *Kavenegar) CountOutbox(p CountOutboxInputParams) (CountOutboxResponse, error) {
	// Base URL
	baseURL := k.CreateURL(COUNTOUTBOX)

	params := url.Values{}
	params.Add("startdate", fmt.Sprintf("%d", p.StartDate))
	if p.EndDate != 0 {
		params.Add("enddate", fmt.Sprintf("%d", p.EndDate))
	}
	if p.Sender != "" {
		params.Add("sender", p.Sender)
	}

	// Construct the full URL with the combined message IDs
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return CountOutboxResponse{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CountOutboxResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response JSON
	var output CountOutboxResponse
	if err := json.Unmarshal(body, &output); err != nil {
		return CountOutboxResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

func (k *Kavenegar) Cancel(messageIDs []int64) (StatusResponse, error) {
	baseURL := k.CreateURL(CANCEL)

	// Combine message IDs into a comma-separated string
	var messageIDStrings []string
	for _, messageID := range messageIDs {
		messageIDStrings = append(messageIDStrings, fmt.Sprintf("%d", messageID))
	}
	messageIDParam := strings.Join(messageIDStrings, ",")

	params := url.Values{}
	params.Add("messageid", messageIDParam)

	// Construct the full URL with the combined message IDs
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	fmt.Println(fullURL)

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return StatusResponse{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return StatusResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response JSON
	var output StatusResponse
	if err := json.Unmarshal(body, &output); err != nil {
		return StatusResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

type ReceiveInputParams struct {
	LineNumber string `json:"linenumber"`
	IsRead     int    `json:"isread"`
}

type ReceiveResponse struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
	Entries []struct {
		MessageID int64     `json:"messageid"`
		Message   string    `json:"message"`
		Sender    string    `json:"sender"`
		Receptor  string    `json:"receptor"`
		Date      time.Time `json:"date"`
	} `json:"entries"`
}

func (k *Kavenegar) Receive(p ReceiveInputParams) (ReceiveResponse, error) {
	baseURL := k.CreateURL(RECEIVE)

	params := url.Values{}
	params.Add("linenumber", p.LineNumber)
	params.Add("isread", fmt.Sprintf("%d", p.IsRead))

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return ReceiveResponse{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ReceiveResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var output ReceiveResponse
	if err := json.Unmarshal(body, &output); err != nil {
		return ReceiveResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

type CountInboxInputParams struct {
	StartDate  int64  `json:"startdate"`
	EndDate    int64  `json:"enddate,omitempty"`    // Optional
	LineNumber string `json:"linenumber,omitempty"` // Optional
	IsRead     int    `json:"isread,omitempty"`     // Optional
}

type CountInboxResponse struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
	Entries []struct {
		StartDate int64 `json:"startdate"`
		EndDate   int64 `json:"enddate"`
		SumCount  int64 `json:"sumcount"`
	} `json:"entries"`
}

func (k *Kavenegar) CountInbox(p CountInboxInputParams) (CountInboxResponse, error) {
	baseURL := k.CreateURL(COUNTINBOX)

	params := url.Values{}
	params.Add("startdate", fmt.Sprintf("%d", p.StartDate))
	if p.EndDate != 0 {
		params.Add("enddate", fmt.Sprintf("%d", p.EndDate))
	}
	if p.LineNumber != "" {
		params.Add("linenumber", p.LineNumber)
	}
	params.Add("isread", fmt.Sprintf("%d", p.IsRead))

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return CountInboxResponse{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CountInboxResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var output CountInboxResponse
	if err := json.Unmarshal(body, &output); err != nil {
		return CountInboxResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

type LookupInputParams struct {
	Receptor string `json:"receptor"`
	Token    string `json:"token"`
	Token2   string `json:"token2,omitempty"` // Optional
	Token3   string `json:"token3,omitempty"` // Optional
	Template string `json:"template"`
	Type     string `json:"type,omitempty"` // Optional
}

func (k *Kavenegar) Lookup(p LookupInputParams) (OutputParams, error) {
	baseURL := k.CreateURL(LOOKUP)

	params := url.Values{}
	params.Add("receptor", p.Receptor)
	params.Add("token", p.Token)
	params.Add("template", p.Template)
	if p.Token2 != "" {
		params.Add("token2", p.Token2)
	}
	if p.Token3 != "" {
		params.Add("token3", p.Token3)
	}
	if p.Type != "" {
		params.Add("type", p.Type)
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var output OutputParams
	if err := json.Unmarshal(body, &output); err != nil {
		return OutputParams{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

type TTSInputParams struct {
	Receptor []string  `json:"receptor"`
	Message  string    `json:"message"`
	Date     time.Time `json:"date,omitempty"`    // Optional
	LocalID  string    `json:"localid,omitempty"` // Optional
	Repeat   int       `json:"repeat,omitempty"`  // Optional
}

func (k *Kavenegar) TTS(p TTSInputParams) (OutputParams, error) {
	baseURL := k.CreateURL(TTS)

	var receptorStrings []string
	receptorStrings = append(receptorStrings, p.Receptor...)
	receptorParam := strings.Join(receptorStrings, ",")

	params := url.Values{}
	params.Add("receptor", receptorParam)
	params.Add("message", p.Message)
	if !p.Date.IsZero() {
		params.Add("date", fmt.Sprintf("%d", p.Date.Unix()))
	}
	if p.LocalID != "" {
		params.Add("localid", p.LocalID)
	}
	if p.Repeat != 0 {
		params.Add("repeat", fmt.Sprintf("%d", p.Repeat))
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OutputParams{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var output OutputParams
	if err := json.Unmarshal(body, &output); err != nil {
		return OutputParams{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

type InfoResponse struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
	Entries struct {
		RemainCredit int64  `json:"remaincredit"`
		Expiredate   string `json:"expiredate"`
		Type         string `json:"type"`
	} `json:"entries"`
}

func (k *Kavenegar) Info() (InfoResponse, error) {
	baseURL := k.CreateURL(INFO)

	// Make the GET request
	resp, err := http.Get(baseURL)
	if err != nil {
		return InfoResponse{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return InfoResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var output InfoResponse
	if err := json.Unmarshal(body, &output); err != nil {
		return InfoResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}

type ConfigInputParams struct {
	APILogs        string `json:"apilogs,omitempty"`        // Optional
	DailyReport    string `json:"dailyreport,omitempty"`    // Optional
	DebugMode      string `json:"debugmode,omitempty"`      // Optional
	DefaultSender  string `json:"defaultsender,omitempty"`  // Optional
	MinCreditAlarm int    `json:"mincreditalarm,omitempty"` // Optional
	ResendFailed   string `json:"resendfailed,omitempty"`   // Optional
}

type ConfigResponse struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
	Entries struct {
		APILogs        string `json:"apilogs"`
		DailyReport    string `json:"dailyreport"`
		DebugMode      string `json:"debugmode"`
		DefaultSender  string `json:"defaultsender"`
		MinCreditAlarm int    `json:"mincreditalarm"`
		ResendFailed   string `json:"resendfailed"`
	} `json:"entries"`
}

func (k *Kavenegar) Config(p ConfigInputParams) (ConfigResponse, error) {
	baseURL := k.CreateURL(CONFIG)

	params := url.Values{}
	if p.APILogs != "" {
		params.Add("apilogs", p.APILogs)
	}
	if p.DailyReport != "" {
		params.Add("dailyreport", p.DailyReport)
	}
	if p.DebugMode != "" {
		params.Add("debugmode", p.DebugMode)
	}
	if p.DefaultSender != "" {
		params.Add("defaultsender", p.DefaultSender)
	}
	if p.MinCreditAlarm != 0 {
		params.Add("mincreditalarm", fmt.Sprintf("%d", p.MinCreditAlarm))
	}
	if p.ResendFailed != "" {
		params.Add("resendfailed", p.ResendFailed)
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return ConfigResponse{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ConfigResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var output ConfigResponse
	if err := json.Unmarshal(body, &output); err != nil {
		return ConfigResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return output, nil
}
