package pkg

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"moul.io/http2curl"
)

func TestWebCacheDeception() reportResult {
	var repResult reportResult
	repResult.Technique = "Cache Deception"

	// cacheable extensions: class, css, jar, js, jpg, jpeg, gif, ico, png, bmp, pict, csv, doc, docx, xls, xlsx, ps, pdf, pls, ppt, pptx, tif, tiff, ttf, otf, webp, woff, woff2, svg, svgz, eot, eps, ejs, swf, torrent, midi, mid

	appendings := []string{
		"/.css",                                       // Path parameter
		"/nonexistent1.css",                           // Path parameter
		"/../nonexistent2.css",                        // Path traversal
		"/%2e%2e/nonexistent3.css",                    // Encoded path traversal
		"%0Anonexistent4.css",                         // Encoded Newline
		"%00nonexistent5.css",                         // Encoded Null Byte
		"%09nonexistent6.css",                         // Encoded Tab
		"%3Bnonexistent7.css",                         // Encoded Semicolon
		"%23nonexistent8.css",                         // Encoded Pound
		"%3Fname=valnonexistent9.css",                 // Encoded Question Mark
		"%26name=valnonexistent10.css",                // Encoded Ampersand
		";nonexistent11.css",                          // Semicolon
		"?nonexistent12.css",                          // Question Mark
		"&nonexistent13.css",                          // Ampersand
		"%0A%2f%2e%2e%2fresources%2fnonexistent1.css", // Encoded Path Traversal to static directory using Encoded Newline
		"%00%2f%2e%2e%2fresources%2fnonexistent2.css", // Encoded Path Traversal to static directory using Encoded Null Byte
		"%09%2f%2e%2e%2fresources%2fnonexistent3.css", // Encoded Path Traversal to static directory using Encoded Tab
		"%3B%2f%2e%2e%2fresources%2fnonexistent4.css", // Encoded Path Traversal to static directoryEncoded using Semicolon
		"%23%2f%2e%2e%2fresources%2fnonexistent5.css", // Encoded Path Traversal to static directory using Encoded Pound
		"%3F%2f%2e%2e%2fresources%2fnonexistent6.css", // Encoded Path Traversal to static directory using Encoded Question Mark
		"%26%2f%2e%2e%2fresources%2fnonexistent7.css", // Encoded Path Traversal to static directory using Encoded Ampersand
		";%2f%2e%2e%2fresources%2fnonexistent8.css",   // Encoded Path Traversal to static directory using Semicolon
		"?%2f%2e%2e%2fresources%2fnonexistent9.css",   // Encoded Path Traversal to static directoy using Question Mark
		"&%2f%2e%2e%2fresources%2fnonexistent10.css",  // Encoded Path Traversal to static directory using Ampersand
		"%0A%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt?a", // Encoded Path Traversal to robots.txt using Encoded Newline
		"%00%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt?b", // Encoded Path Traversal to robots.txt directory using Encoded Null Byte
		"%09%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt?c", // Encoded Path Traversal to robots.txt directory using Encoded Tab
		"%3B%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt?d", // Encoded Path Traversal to robots.txt directoryEncoded using Semicolon
		"%23%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt?e", // Encoded Path Traversal to robots.txt directory using Encoded Pound
		"%3F%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt?f", // Encoded Path Traversal to robots.txt directory using Encoded Question Mark
		"%26%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt?g", // Encoded Path Traversal to robots.txt directory using Encoded Ampersand
		";%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt?h",   // Encoded Path Traversal to robots.txt directory using Semicolon
		"?%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt?i",   // Encoded Path Traversal to robots.txt directoy using Question Mark
		"&%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt?j",   // Encoded Path Traversal to robots.txt directory using Ampersand
	}
	// TODO add "Exploiting normalization by the origin server" cache deception which needs to prepend something before the url path

	if Config.Website.StatusCode != 200 || Config.Website.Body == "" {
		msg := "Skipping Web Cache Deception test, as it requires a valid website configuration with a status code of 200 and a non-empty body.\n"
		Print(msg, Yellow)
		repResult.HasError = true
		repResult.ErrorMessages = append(repResult.ErrorMessages, msg)
		return repResult
	}
	PrintVerbose("Testing for Web Cache Deception\n", NoColor, 1)

	// test each appending one after another
	for _, appendStr := range appendings {
		err := webCacheDeceptionTemplate(&repResult, appendStr)
		if err != nil {
			repResult.HasError = true
			repResult.ErrorMessages = append(repResult.ErrorMessages, err.Error())
		}
	}

	return repResult
}

func webCacheDeceptionTemplate(repResult *reportResult, appendStr string) error {
	var msg string
	var repCheck reportCheck
	var req *http.Request
	var resp *http.Response
	var err error
	var body []byte

	rUrl := Config.Website.Url.String()
	// Überprüfen, ob der String genau zwei `//` enthält
	if strings.Count(rUrl, "/") == 2 && !strings.HasPrefix(appendStr, "/") {
		// append `/`, so e.g. https://example%0A does not throw an error when building the request
		rUrl += "/"
	}

	req, err = http.NewRequest("GET", rUrl+appendStr, nil)
	if err != nil {
		msg = fmt.Sprintf("webCacheDeceptionTemplate: %s: http.NewRequest: %s\n", appendStr, err.Error())
		Print(msg, Red)
		return errors.New(msg)
	}

	// Do request
	newClient := http.Client{
		CheckRedirect: http.DefaultClient.CheckRedirect,
		Timeout:       http.DefaultClient.Timeout,
	}
	setRequest(req, false, "", http.Cookie{}, false)
	resp, err = newClient.Do(req)
	if err != nil {
		msg = fmt.Sprintf("webCacheDeceptionTemplate: %s: newClient.Do: %s\n", appendStr, err.Error())
		PrintVerbose(msg, Yellow, 1)
		return errors.New(msg)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		msg = fmt.Sprintf("webCacheDeceptionTemplate: %s: io.ReadAll: %s\n", appendStr, err.Error())
		Print(msg, Red)
		return errors.New(msg)
	}
	if resp.StatusCode != Config.Website.StatusCode || string(body) != Config.Website.Body {
		return nil // no cache deception, as the response is not the same as the original one
	}

	if Config.Website.Cache.NoCache || Config.Website.Cache.Indicator == "age" {
		time.Sleep(1 * time.Second) // wait a second to ensure that age header is not set to 0
	}

	resp, err = newClient.Do(req)
	if err != nil {
		msg = fmt.Sprintf("webCacheDeceptionTemplate: %s: newClient.Do: %s\n", appendStr, err.Error())
		PrintVerbose(msg, Yellow, 1)
		return errors.New(msg)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		msg = fmt.Sprintf("webCacheDeceptionTemplate: %s: io.ReadAll: %s\n", appendStr, err.Error())
		Print(msg, Red)
		return errors.New(msg)
	}

	// Add the request as curl command to the report
	command, err := http2curl.GetCurlCommand(req)
	if err != nil {
		PrintVerbose("Error: http2curl: "+err.Error()+"\n", Yellow, 1)
	}
	repCheck.Request.CurlCommand = command.String()

	var indicValue string
	if Config.Website.Cache.Indicator == "" { // check if now a cache indicator exists
		cache := analyzeCacheIndicator(resp.Header)
		Config.Website.Cache.Indicator = cache.Indicator
	}

	indicValue = strings.TrimSpace(strings.ToLower(resp.Header.Get(Config.Website.Cache.Indicator)))

	// check if there's a cache hit and if the body didn't change (otherwise it could be a cached error page, for example)
	if checkCacheHit(indicValue, "") && string(body) == Config.Website.Body {
		repResult.Vulnerable = true
		repCheck.Reason = "The response got cached due to Web Cache Deception"
		msg = fmt.Sprintf("%s was successfully decepted! appended: %s\n", rUrl, appendStr)
		Print(msg, Green)
		msg = "Curl: " + repCheck.Request.CurlCommand + "\n\n"
		Print(msg, Green)

		repCheck.Identifier = appendStr
		repCheck.URL = req.URL.String()
		// Dump the request without the body
		var dumpReqBytes []byte
		dumpReqBytes, _ = httputil.DumpRequest(req, false)
		repCheck.Request.Request = string(dumpReqBytes)
		// Dump the response
		responseBytes, _ := httputil.DumpResponse(resp, false)
		repCheck.Request.Response = string(responseBytes)

		repResult.Checks = append(repResult.Checks, repCheck)
	} else {
		PrintVerbose("Curl command: "+repCheck.Request.CurlCommand+"\n", NoColor, 2)
	}

	return nil
}
