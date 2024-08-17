package browse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

func New(req *BrowseRequest, ctx context.Context) (*rod.Page, error) {
	if req == nil {
		req = &BrowseRequest{}
	}

	proxyServer, err := url.Parse(req.Proxy)
	if err != nil {
		return nil, err
	}

	u, err := launcher.New().Headless(req.Headless).Proxy(proxyServer.Host).Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().ControlURL(u)
	go func() {
		if proxyServer.User.Username() == "" {
			return
		}

		password, _ := proxyServer.User.Password()
		errC := browser.HandleAuth(proxyServer.User.Username(), password)
		err = errC()
		if err != nil {
			log.Println(err)
		}
	}()

	err = browser.Connect()
	if err != nil {
		return nil, err
	}

	browser.MustIgnoreCertErrors(true)

	page, err := stealth.Page(browser)
	if err != nil {
		return nil, err
	}

	err = page.SetCookies(req.Cookies)
	if err != nil {
		return nil, err
	}

	return page.Context(ctx), nil
}

func (req *BrowseRequest) Hijack(page *rod.Page) (*BrowseResponse, error) {
	resp := &BrowseResponse{
		Hijacks: make(map[string]any),
	}
	if len(req.HijackRequests) == 0 {
		return resp, nil
	}

	proxyServer, err := url.Parse(req.Proxy)
	if err != nil {
		return nil, err
	}

	proxyClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyServer),
		},
	}

	router := page.HijackRequests()
	for _, pattern := range req.HijackRequests {
		err := router.Add(pattern, "", func(ctx *rod.Hijack) {

			ctx.LoadResponse(proxyClient, true)
			var data any
			body := ctx.Response.Body()
			if strings.Contains(ctx.Response.Headers().Get("Content-Type"), "application/json") {
				err := json.Unmarshal([]byte(body), &data)
				if err != nil {
					data = body
				}
			}

			resp.Hijacks[pattern] = data
		})
		if err != nil {
			return nil, err
		}
	}

	go router.Run()
	return resp, nil
}

func (in *BrowseRequestInstruction) Act(page *rod.Page) (map[string]string, error) {
	element, err := in.GetElement(page)
	if err != nil {
		return nil, err
	}

	var params map[string]string
	var actionErr error

	switch in.Action {
	case Click:
		actionErr = element.Click(proto.InputMouseButtonLeft, 1)

	case Type:
		actionErr = element.Input(in.Input)

	default:
		return nil, fmt.Errorf("unexpected instructed action %s", in.Action)
	}

	return params, actionErr

}

func (in *BrowseRequestInstruction) GetElement(page *rod.Page) (*rod.Element, error) {
	currentFrame := page
	for _, frame := range in.Frames {
		element, err := currentFrame.Element(frame)
		if err != nil {
			return nil, err
		}

		currentFrame, err = element.Frame()
		if err != nil {
			return nil, err
		}
	}

	return currentFrame.Element(in.Selector)
}
