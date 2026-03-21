<h1 align="center">
  <img src="https://www.hackmanit.de/images/beitragsbilder/blog/Web-Cache-Vulnerability-Scanner-Banner.png" width="100%" alt="Web Cache Vulnerability Scanner"/>
</h1>

<p align="center">
  <a href="https://github.com/Hackmanit/Web-Cache-Vulnerability-Scanner/releases/latest"><img src="https://img.shields.io/github/release/Hackmanit/Web-Cache-Vulnerability-Scanner.svg?color=brightgreen" alt="Release"></a>
  <a href="https://goreportcard.com/report/github.com/Hackmanit/Web-Cache-Vulnerability-Scanner"><img src="https://goreportcard.com/badge/github.com/Hackmanit/Web-Cache-Vulnerability-Scanner" alt="Go Report Card"></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/github/go-mod/go-version/Hackmanit/Web-Cache-Vulnerability-Scanner" alt="Go Version"></a>
  <a href="https://www.apache.org/licenses/LICENSE-2.0"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License"></a>
</p>

**Web Cache Vulnerability Scanner (WCVS)** is a fast and versatile CLI tool for detecting [web cache poisoning](#web-cache-poisoning) and [web cache deception](#web-cache-deception) vulnerabilities. Developed by [Hackmanit](https://hackmanit.de) and [Maximilian Hildebrand](https://www.github.com/m10x), WCVS automates complex cache-based attack techniques, includes a built-in crawler, and adapts to specific web cache configurations for more efficient and accurate testing.

Whether you're a penetration tester, a bug bounty hunter, or a security engineer integrating checks into CI/CD pipelines — WCVS covers the techniques you need.

---

## Table of Contents

- [Features](#features)
- [Web Cache Deception Coverage](#web-cache-deception-coverage)
- [Installation](#installation)
  - [Option 1: Pre-built Binary](#option-1-pre-built-binary)
  - [Option 2: Kali Linux / BlackArch Repository](#option-2-kali-linux--blackarch-repository)
  - [Option 3: Install Using Go](#option-3-install-using-go)
  - [Option 4: Docker](#option-4-docker)
- [Usage](#usage)
  - [Specify Headers, Parameters, Cookies, and More](#specify-headers-parameters-cookies-and-more)
  - [Generate a JSON Report](#generate-a-json-report)
  - [Crawl for URLs](#crawl-for-urls)
  - [Use a Proxy](#use-a-proxy)
  - [Throttle or Accelerate](#throttle-or-accelerate)
  - [Further Flags](#further-flags)
- [Background Information](#background-information)
- [License](#license)

---

## Features

### Web Cache Poisoning (10 techniques)
1. Unkeyed header poisoning
2. Unkeyed parameter poisoning
3. Parameter cloaking
4. Fat GET
5. HTTP response splitting
6. HTTP request smuggling
7. HTTP header oversize (HHO)
8. HTTP meta character (HMC)
9. HTTP method override (HMO)
10. Parameter pollution

### Web Cache Deception (extensive coverage)
- **Path parameter injection** — appending static-looking path segments (e.g., `/.css`, `/nonexistent.css`)
- **Path traversal** — using `/../`, `/%2e%2e/`, double-encoded, and Tomcat-style (`/..;/`) traversals targeting `.css`, `.js`, and `/robots.txt`
- **Origin server normalization exploitation** — cache keys a path under a static prefix (e.g., `/static/`), while the origin decodes `%2F..%2F` and serves the sensitive resource
- **Single-level encoded path traversal to `/robots.txt`** — using delimiters like `;`, `?`, `&`, `%0A`, `%09`, `%00`, `%3B`, `%23`, `%3F`, `%26` followed by `%2f%2e%2e%2frobots.txt`
- **Special character delimiters** — both encoded (`%0A`, `%09`, `%00`, `%3B`, `%23`, `%3F`, `%26`) and literal (`;`, `?`, `&`) before static-looking extensions
- **Double URL-encoding** — `%252e%252e%2F` and `%252F..%252F` style traversals
- **Multiple static file extensions** — `.css`, `.js`, `.png`, `.ico`, `.woff2`, `.svg`, `.json`

### Additional Capabilities
- Automatic web cache fingerprinting before testing (adapts strategy per cache type)
- JSON report generation (with optional HTML special-character escaping)
- Built-in URL crawler with configurable depth, domain filtering, and exclusions
- Proxy support (Burp Suite, OWASP ZAP, etc.)
- Rate limiting and multi-threading controls for responsible testing
- CI/CD pipeline-friendly design

---

## Web Cache Deception Coverage

Web cache deception attacks trick a cache into storing a response to a sensitive, authenticated resource by making the request URL look like a public static file. WCVS tests a comprehensive set of techniques including:

| Technique | Example Pattern | Description |
|---|---|---|
| Path parameter | `/.css` | Appends a static-looking path segment |
| Path traversal (unencoded) | `/../nonexistent.css` | Traverses up using standard `..` |
| Path traversal (encoded) | `/%2e%2e/nonexistent.css` | Traverses using percent-encoded dots |
| Encoded delimiter + traversal | `;%2f%2e%2e%2frobots.txt` | Semicolon delimiter with encoded traversal to robots.txt |
| Origin normalization | `/static/..%2Fmy-account` | Cache keys static path; origin resolves traversal |
| Double URL-encoding | `%252e%252e%2Fnonexistent.css` | Bypasses single-decode defences |
| Tomcat-style traversal | `/..;/nonexistent.css` | Uses `..;` path traversal |
| Special char + extension | `%0Anonexistent.css` | Newline or other control chars before extension |
| Query/fragment injection | `?nonexistent.css` | Question mark before static-looking path |

---

## Installation

### Option 1: Pre-built Binary
Prebuilt binaries for Linux, macOS, and Windows are available on the [releases page](https://github.com/Hackmanit/Web-Cache-Vulnerability-Scanner/releases).

### Option 2: Kali Linux / BlackArch Repository
```bash
# Kali Linux
apt install web-cache-vulnerability-scanner

# BlackArch
pacman -S wcvs
```

### Option 3: Install Using Go
Requires Go 1.21 or higher.
```bash
go install -v github.com/Hackmanit/Web-Cache-Vulnerability-Scanner@latest
```

### Option 4: Docker

**1. Clone repository or download the [latest source code release](https://github.com/Hackmanit/Web-Cache-Vulnerability-Scanner/releases/latest)**

**2. Build the Docker image** (the wordlists folder is automatically included):
```bash
docker build -t wcvs .
```

**3. Run WCVS**:
```bash
docker run -it wcvs /wcvs --help
```

---

## Usage

WCVS is highly customizable via flags. Most flags accept either a direct value or a path to a file prefixed with `file:`.

The only required flag is `-u/--url` — the target URL to test. WCVS accepts several URL formats:

```bash
wcvs -u 127.0.0.1
wcvs -u http://127.0.0.1
wcvs -u https://example.com
wcvs -u file:path/to/url_list
```

> **Note:** Two wordlists are required for the first 5 poisoning techniques — one for headers and one for parameters. Place them in the same directory as WCVS, or specify them with `--headerwordlist/-hw` and `--parameterwordlist/-pw`.

```bash
wcvs -u https://example.com -hw "file:/home/user/Documents/wordlist-header.txt"
wcvs -u https://example.com -pw "file:/home/user/Documents/wordlist-parameter.txt"
wcvs -u https://example.com -hw "file:/home/user/wordlist-header.txt" -pw "file:/home/user/wordlist-parameter.txt"
```

---

## Specify Headers, Parameters, Cookies, and More

| Flag | Short | Description |
|---|---|---|
| `--cacheheader` | `-ch` | Custom cache header to detect hits/misses |
| `--setcookies` | `-sc` | Cookies to add to every request |
| `--setheaders` | `-sh` | Headers to add to every request |
| `--setparameters` | `-sp` | URL parameters to add to every request |
| `--post` | `-post` | Use POST instead of GET |
| `--setbody` | `-sb` | Request body to send (used with `-post`) |
| `--contenttype` | `-ct` | Value for the `Content-Type` header |
| `--useragentchrome` | `-uac` | Use a Chrome user-agent string |

> **Tip:** To specify more than one cookie, header, or parameter, use a file. See the [available templates](https://github.com/Hackmanit/Web-Cache-Vulnerability-Scanner/tree/master/templates).

### Examples:
```bash
wcvs -u https://example.com -ch "X-Custom-Header-ABC"

# Cookies
wcvs -u https://example.com -sc "PHPSESSID=123"
wcvs -u https://example.com -sc "file:/home/user/Documents/cookies.txt"

# Headers
wcvs -u https://example.com -sh "Referer: localhost"
wcvs -u https://example.com -sh "file:/home/user/Documents/headers.txt"

# Parameters
wcvs -u https://example.com -sp "admin=true"
wcvs -u https://example.com -sp "file:/home/user/Documents/parameters.txt"

# POST with body
wcvs -u https://example.com -post -sb "admin=true"
wcvs -u https://example.com -post -sb "file:/home/user/Documents/body.txt"
wcvs -u https://example.com -post -sb "{}" -ct "application/json"

# Chrome User-Agent
wcvs -u https://example.com -uac
```

---

## Generate a JSON Report

Use `--generatereport/-gr` to save a JSON report that is updated after each scanned URL.

| Flag | Short | Description |
|---|---|---|
| `--generatereport` | `-gr` | Enable JSON report generation |
| `--generatepath` | `-gp` | Directory where report and log files are written (default: `./`) |
| `--escapejson` | `-ej` | Encode HTML special chars in the report |

### Examples:
```bash
wcvs -u https://example.com -gr
wcvs -u https://example.com -gr -ej
wcvs -u https://example.com -gr -gp /home/user/Documents
wcvs -u https://example.com -gr -gp /home/user/Documents -ej
```

---

## Crawl for URLs

WCVS includes a built-in crawler to discover and test additional pages automatically.

| Flag | Short | Description |
|---|---|---|
| `--recursivity` | `-r` | Crawl depth (number of recursion levels) |
| `--recdomains` | `-red` | Also crawl external/cross-domain URLs |
| `--recinclude` | `-rin` | Only crawl URLs containing a specific string |
| `--reclimit` | `-rl` | Max URLs to crawl per recursion depth |
| `--recexclude` | `-rex` | File with URLs to skip |
| `--generatecompleted` | `-gc` | Save a list of all tested URLs for future exclusion |

### Examples:
```bash
wcvs -u https://example.com -r 5
wcvs -u https://example.com -r 5 -red /home/user/Documents/mydomains.txt
wcvs -u https://example.com -r 5 -rl 2
wcvs -u https://example.com -r 5 -rex /home/user/Documents/donttest.txt
```

---

## Use a Proxy

Use `--useproxy/-up` to route traffic through a proxy such as Burp Suite or OWASP ZAP.

> **Burp Suite note:** Uncheck *"Settings > Network > HTTP > HTTP/2 > Default to HTTP/2 if the server supports it"* — otherwise some non-RFC-compliant techniques will fail.

| Flag | Short | Description |
|---|---|---|
| `--useproxy` | `-up` | Enable proxy (default: `http://127.0.0.1:8080`) |
| `--proxyurl` | `-purl` | Custom proxy URL |

### Examples:
```bash
wcvs -u https://example.com -up
wcvs -u https://example.com -up -purl http://127.0.0.1:8081
```

---

## Throttle or Accelerate

| Flag | Short | Description |
|---|---|---|
| `--reqrate` | `-rr` | Max requests per second (default: unrestricted) |
| `--threads` | `-t` | Number of concurrent threads (default: 20) |

### Examples:
```bash
wcvs -u https://example.com -rr 10
wcvs -u https://example.com -rr 1
wcvs -u https://example.com -rr 0.5
wcvs -u https://example.com -t 50
```

---

## Further Flags

Run `wcvs -h` for the full list of flags, descriptions, and usage:

```bash
wcvs -h
```

---

## Background Information

### Web Cache Poisoning

Web cache poisoning exploits discrepancies between what a cache stores and what it actually serves to users. By injecting malicious content via unkeyed inputs (headers, parameters, etc.), an attacker can poison a cached response that is then delivered to every subsequent visitor.

### Web Cache Deception

Web cache deception tricks a cache into storing a response to a sensitive, authenticated resource by disguising the request URL as a public static file. For example, appending `;%2f%2e%2e%2frobots.txt` to a sensitive URL may cause the cache to key the response as a static resource while the origin still serves the sensitive page.

WCVS tests a wide range of deception patterns — including single-level and multi-level encoded path traversals, delimiter injection, and origin-server normalization attacks — to provide broad coverage against real-world targets.

### Further Reading

A short series of blog posts giving more context about web cache poisoning and WCVS:

1. [Is Your Application Vulnerable to Web Cache Poisoning?](https://www.hackmanit.de/en/blog-en/142-is-your-application-vulnerable-to-web-cache-poisoning)
2. [Web Cache Vulnerability Scanner (WCVS) - Free, Customizable, Easy-To-Use](https://www.hackmanit.de/en/blog-en/145-web-cache-vulnerability-scanner-wcvs-free-customizable-easy-to-use)

The first version of WCVS was developed as part of a [bachelor's thesis by Maximilian Hildebrand](https://hackmanit.de/images/download/thesis/Automated-Scanning-for-Web-Cache-Poisoning-Vulnerabilities.pdf).

---

## License

WCVS is developed by [Hackmanit](https://hackmanit.de) and [Maximilian Hildebrand](https://www.github.com/m10x) and licensed under the [Apache License, Version 2.0](LICENSE).

<a href="https://hackmanit.de"><img src="https://www.hackmanit.de/templates/hackmanit-v2/img/wbm_hackmanit.png" width="30%"></a>
