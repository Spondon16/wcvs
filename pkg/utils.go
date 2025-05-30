package pkg

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	NoColor = 0
	Red     = 1
	Yellow  = 2
	Green   = 3
	Cyan    = 4
)

func PrintNewLine() {
	Print("\n", NoColor)
}

func PrintLog(msg string) {
	if !Config.NoLog && Config.Intitialized {
		log.Print(msg)
	}
}

func PrintVerbose(msg string, c int, threshold int) {
	if c == Red {
		PrintLog("[ERR] " + msg)
		msg = color.RedString("[ERR] ") + msg
	} else if c == Yellow {
		PrintLog("[!] " + msg)
		msg = color.YellowString("[!] ") + msg
	} else if c == Green {
		PrintLog("[+] " + msg)
		msg = color.GreenString("[+] ") + msg
	} else if c == Cyan {
		PrintLog("[*] " + msg)
		msg = color.CyanString("[*] ") + msg
	} else {
		PrintLog(msg)
	}

	if Config.Verbosity >= threshold || !Config.Intitialized {
		fmt.Print(msg)
	}
}

func Print(msg string, c int) {
	PrintVerbose(msg, c, 0)
}

func PrintFatal(msg string) {
	Print(msg, Red)
	os.Exit(1)
}

func ReadLocalFile(path string, name string) []string {
	path = strings.TrimPrefix(path, "file:")

	if strings.HasPrefix(strings.ToLower(path), "file:") {
		PrintFatal("Please make sure that path: is lowercase")
	}

	w, err := os.ReadFile(path)
	if err != nil {
		additional := ""
		if name == "header" {
			additional = "Use the flag \"-hw path/to/wordlist\" to specify the path to a header wordlist\n"
		} else if name == "parameter" {
			additional = "Use the flag \"-pw path/to/wordlist\" to specify the path to a parameter wordlist\n"
		}
		PrintFatal("The specified " + name + " file path " + path + " couldn't be found: " + err.Error() + "\n" + additional)
	}

	return strings.Split(string(w), "\n")
}

func setRequest(req *http.Request, doPost bool, cb string, cookie http.Cookie, prependCB bool) {

	cache := Config.Website.Cache
	if cb != "" && cache.CBisParameter {
		var newUrl string
		newUrl, _ = addCachebusterParameter(req.URL.String(), cb, cache.CBName, prependCB)

		var err error
		req.URL, err = url.Parse(newUrl)
		if err != nil {
			msg := "Converting " + newUrl + " to URL:" + err.Error() + "\n"
			Print(msg, Red)
		}
	}

	setRequestHeaders(req, cb)
	//TODO config nötig oder nur config.Website.Cookies?
	setRequestCookies(req, cb, cookie)

	// Content-Type nur hinzufügen, wenn nicht schon vorher geschehen
	if doPost {
		if req.Header.Get("Content-Type") == "" && Config.ContentType != "" {
			req.Header.Add("Content-Type", Config.ContentType)
		}
	}
}

/* TODO wie bei requestCookies nur die erste occurrence eines headers aufnehmen */
func setRequestHeaders(req *http.Request, cb string) {
	cache := Config.Website.Cache

	req.Header.Set("User-Agent", useragent)
	for _, h := range Config.Headers {
		h = strings.TrimSuffix(h, "\r")
		h = strings.TrimSpace(h)
		if h == "" {
			continue
		} else if !strings.Contains(h, ":") {
			msg := "Specified header" + h + "doesn't contain a : and will be skipped"
			Print(msg, NoColor)
			continue
		} else {
			hSplitted := strings.SplitN(h, ":", 2)

			// is this header the cachebuster?
			if cb != "" && cache.CBisHeader && strings.EqualFold(hSplitted[0], cache.CBName) {
				req.Header.Set(cache.CBName, cb)
			}

			req.Header.Set(strings.TrimSpace(hSplitted[0]), strings.TrimSpace(hSplitted[1]))
		}
	}
}

func setRequestCookies(req *http.Request, cb string, cookie http.Cookie) {
	cache := Config.Website.Cache

	for _, c := range Config.Website.Cookies {
		// only add first occurrence of a cookie to the request
		_, err := req.Cookie(c.Name)
		if err == http.ErrNoCookie {
			if cb != "" && cache.CBisCookie && c.Name == cache.CBName {
				c.Value = cb

				if c.Name == cookie.Name {
					msg := "Can't test cookie " + c.Name + " for Web Cache Poisoning, as it is used as Cachebuster\n"
					Print(msg, Yellow)
				}
			} else if c.Name == cookie.Name {
				c = &cookie
			}
			req.AddCookie(c)
		}
	}
}

func addCachebusterParameter(strUrl string, cbvalue string, cb string, prepend bool) (string, string) {
	if cbvalue == "" {
		cbvalue = "cb" + randInt()
	}
	if cb == "" {
		cb = Config.Website.Cache.CBName
	}
	if !strings.Contains(strUrl, "?") {
		strUrl += "?" + cb + "=" + cbvalue
	} else {
		if prepend {
			parts := strings.SplitN(strUrl, "?", 2)
			strUrl = parts[0] + "?" + cb + "=" + cbvalue + Config.QuerySeparator + parts[1]
		} else {
			strUrl += Config.QuerySeparator + cb + "=" + cbvalue
		}
	}

	return strUrl, cbvalue
}

// RandomString generates a random string of the specified length
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			Print(err.Error(), Red)
			return "99999999"
		}
		result[i] = charset[n.Int64()]
	}
	return string(result)
}

func removeParam(rawURL string, paramToRemove string) (string, string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}

	// Get current query parameters
	query := parsedURL.Query()

	// Check if the parameter exists and remove it
	if _, exists := query[paramToRemove]; exists {
		query.Del(paramToRemove)
		parsedURL.RawQuery = query.Encode()
	}

	return parsedURL.String(), query.Get(paramToRemove), nil
}

/* Create a random long integer */
func randInt() string {
	min := int64(100000000000)
	max := int64(999999999999)
	// Range size
	rangeSize := max - min + 1

	n, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
	if err != nil {
		Print(err.Error(), Red)
		return "999999999999"
	}

	result := n.Int64() + min
	return strconv.FormatInt(result, 10)
}

func waitLimiter(identifier string) {
	err := Config.Limiter.Wait(context.Background())
	if err != nil {
		msg := identifier + " rate Wait: " + err.Error()
		Print(msg, Red)
	}
}

func searchBodyHeadersForString(cb string, body string, headers http.Header) bool {
	if strings.Contains(body, cb) {
		return true
	}
	for _, h := range headers {
		for _, v := range h {
			if strings.Contains(v, cb) {
				return true
			}
		}
	}
	return false
}

// check if cache was hit
func checkCacheHit(value string, indicator string) bool {
	if indicator == "" {
		indicator = Config.Website.Cache.Indicator
	}
	if strings.EqualFold("age", indicator) {
		value = strings.TrimSpace(value)
		if value != "0" && value != "" {
			return true
		}
	} else if strings.EqualFold("x-iinfo", indicator) {
		// String anhand von Leerzeichen aufteilen
		parts := strings.Split(value, " ")

		// Prüfen, ob der zweite Part existiert
		if len(parts) > 1 {
			secondPart := parts[1]

			// Sicherstellen, dass der zweite Part mindestens zwei Zeichen lang ist
			if len(secondPart) > 1 {
				secondChar := strings.ToUpper(string(secondPart[1]))
				if secondChar == "C" || secondChar == "V" {
					return true
				}
			}
		}
		// Cache Hit may have 0,>0 or >0,0 as value. Both responses are cached
	} else if strings.EqualFold("x-cache-hits", indicator) {
		for _, x := range strings.Split(indicator, ",") {
			x = strings.TrimSpace(x)
			if x != "0" {
				return true
			}
		}
		// Some Headers may have "miss,hit" or "hit,miss" as value. But both are cached responses.
	} else if strings.Contains(strings.ToLower(value), "hit") || strings.Contains(strings.ToLower(value), "cached") {
		return true
	}
	return false
}
