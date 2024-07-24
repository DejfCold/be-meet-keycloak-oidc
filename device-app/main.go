package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code,omitempty"`
	UserCode                string `json:"user_code,omitempty"`
	VerificationUri         string `json:"verification_uri,omitempty"`
	VerificationUriComplete string `json:"verification_uri_complete,omitempty"`
	ExpiresIn               int    `json:"expires_in,omitempty"`
	Interval                int    `json:"interval,omitempty"`
}

type DeviceTokenResponse struct {
	Error            *string `json:"error,omitempty"`
	ErrorDescription *string `json:"error_description,omitempty"`
	AccessToken      string  `json:"access_token,omitempty"`
	ExpiresIn        int     `json:"expires_in,omitempty"`
	RefreshExpiresIn int     `json:"refresh_expires_in,omitempty"`
	RefreshToken     string  `json:"refresh_token,omitempty"`
	TokenType        string  `json:"token_type,omitempty"`
	IdToken          string  `json:"id_token,omitempty"`
	NotBeforePolicy  int     `json:"not-before-policy,omitempty"`
	SessionState     string  `json:"session_state,omitempty"`
	Scope            string  `json:"scope,omitempty"`
}

func main() {
	headerTxt := "Device Login Example"

	a := app.New()

	window := a.NewWindow(headerTxt)
	window.Resize(fyne.Size{
		Width:  500,
		Height: 300,
	})

	loginBtn := widget.NewButton("Login", func() {})

	cont := container.NewVBox(
		loginBtn)

	loginBtn.OnTapped = func() {
		askForLoginOpen(window)
	}

	window.SetContent(cont)

	window.ShowAndRun()
}

func askForLoginOpen(win fyne.Window) {
	fmt.Println("Tapped")
	deviceCodeResponse := getDeviceCode()
	compUrl := deviceCodeResponse.VerificationUriComplete
	loginUrl, err := url.ParseRequestURI(compUrl)
	if err != nil {
		fmt.Println(err.Error())
	}

	hint := widget.NewLabel("Open in browser to login!")
	hint.Alignment = fyne.TextAlignCenter
	link := widget.NewHyperlinkWithStyle(compUrl, loginUrl, fyne.TextAlignCenter, fyne.TextStyle{Underline: true})
	link.OnTapped = func() {
		fyne.Clipboard.SetContent(win.Clipboard(), compUrl)
		waitForCompletion(win, deviceCodeResponse.Interval, deviceCodeResponse.DeviceCode)
	}

	hint2 := widget.NewLabel("Or open this")
	hint2.Alignment = fyne.TextAlignCenter
	partUrl := deviceCodeResponse.VerificationUri
	loginUrlPart, err := url.ParseRequestURI(partUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	link2 := widget.NewHyperlinkWithStyle(partUrl, loginUrlPart, fyne.TextAlignCenter, fyne.TextStyle{Underline: true})
	link2.OnTapped = func() {
		fyne.Clipboard.SetContent(win.Clipboard(), partUrl)
		waitForCompletion(win, deviceCodeResponse.Interval, deviceCodeResponse.DeviceCode)
	}

	hint3 := widget.NewLabel("and then enter this code:")
	hint3.Alignment = fyne.TextAlignCenter
	codeLabel := widget.NewLabel(deviceCodeResponse.UserCode)
	codeLabel.Alignment = fyne.TextAlignCenter

	cont := container.NewVBox(hint, link, hint2, link2, hint3, codeLabel)
	win.SetContent(cont)
}

func waitForCompletion(win fyne.Window, interval int, deviceCode string) {
	pb := widget.NewProgressBarInfinite()
	pb.Start()
	cont := container.NewVBox(pb)
	win.SetContent(cont)

	var deviceTokenResponse DeviceTokenResponse
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		isOk, val := checkIsLoggedIn(deviceCode)
		deviceTokenResponse = val
		if isOk {
			break
		}
	}

	showCompleted(win, deviceTokenResponse)

}

func showCompleted(win fyne.Window, response DeviceTokenResponse) {
	loggedInLabel := widget.NewLabel("Successfully logged in!")
	sep1 := widget.NewSeparator()
	sep2 := widget.NewSeparator()

	accessToken := response.AccessToken
	tokenParts := strings.Split(accessToken, ".")
	payload, err := base64.RawStdEncoding.DecodeString(tokenParts[1])
	if err != nil {
		fmt.Println("Decode err", err)
	}

	fmt.Println()
	fmt.Println(tokenParts[1])
	fmt.Println()

	var prettyPayload bytes.Buffer
	err = json.Indent(&prettyPayload, payload, "", "\t")
	if err != nil {
		fmt.Println("Indent err", err)
	}

	payloadGrid := widget.NewTextGridFromString(prettyPayload.String())
	exitBtn := widget.NewButton("Exit", func() {
		os.Exit(0)
	})

	fmt.Println("payload ", prettyPayload.String())

	cont := container.NewVBox(loggedInLabel, sep1, payloadGrid, sep2, exitBtn)
	win.SetContent(cont)
}

func checkIsLoggedIn(deviceCode string) (bool, DeviceTokenResponse) {
	fmt.Println("Checking if logged in")
	resp, err := http.PostForm("http://localhost:8080/realms/be-meet/protocol/openid-connect/token",
		url.Values{
			"client_id":   {"be-meet-device"},
			"device_code": {deviceCode},
			"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		})
	if err != nil {
		fmt.Println(err.Error())
		return false, DeviceTokenResponse{}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println("check body:", string(body))

	var deviceTokenResponse DeviceTokenResponse
	err = json.Unmarshal(body, &deviceTokenResponse)
	if err != nil {
		fmt.Println("unmarshall error", err.Error())
		return false, DeviceTokenResponse{}
	}

	if deviceTokenResponse.Error != nil {
		fmt.Println("Not yet authed")
		return false, DeviceTokenResponse{}
	}
	return true, deviceTokenResponse
}

func getDeviceCode() DeviceCodeResponse {

	fmt.Println("Getting device code")
	resp, err := http.PostForm("http://localhost:8080/realms/be-meet/protocol/openid-connect/auth/device",
		url.Values{"client_id": {"be-meet-device"}})
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var deviceCodeResponse DeviceCodeResponse
	err = json.Unmarshal(body, &deviceCodeResponse)
	if err != nil {
		fmt.Println("unmarshall error", err.Error())
		return DeviceCodeResponse{}
	}
	fmt.Println("device code", deviceCodeResponse.DeviceCode)
	return deviceCodeResponse
}
