<h1 align="center">
  <img src="https://www.hackmanit.de/images/beitragsbilder/blog/Web-Cache-Vulnerability-Scanner-Banner.png" width="100%" alt="Web Cache Vulnerability Scanner"/>
</h1>

<p align="center">
  <a href="https://github.com/Spondon16/wcvs/releases/latest"><img src="https://img.shields.io/github/release/Spondon16/wcvs.svg?color=brightgreen" alt="Release"></a>
  <a href="https://goreportcard.com/report/github.com/Spondon16/wcvs"><img src="https://goreportcard.com/badge/github.com/Spondon16/wcvs" alt="Go Report Card"></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/github/go-mod/go-version/Spondon16/wcvs" alt="Go Version"></a>
  <a href="https://www.apache.org/licenses/LICENSE-2.0"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License"></a>
</p>

**Web Cache Vulnerability Scanner (WCVS)** — a fork by [Spondon16](https://github.com/Spondon16) with improved Web Cache Deception coverage and bug fixes.

Based on the original by [Hackmanit](https://hackmanit.de) and [Maximilian Hildebrand](https://www.github.com/m10x), this fork extends the deception technique set and provides a reliable, up-to-date build.

---

## Table of Contents

- [Features](#features)
- [Web Cache Deception Coverage](#web-cache-deception-coverage)
- [Installation](#installation)
- [Usage](#usage)
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

### Option 1: Pre-built Binary (Recommended)
Prebuilt binaries for Linux, macOS, and Windows are available on the [releases page](https://github.com/Spondon16/wcvs/releases).

### Option 2: Install Using Go
```bash
go install -v github.com/Spondon16/wcvs@latest
```

### Option 3: Build from Source
```bash
git clone https://github.com/Spondon16/wcvs.git
cd wcvs
go build -o wcvs .
```

### Option 4: Docker
```bash
git clone https://github.com/Spondon16/wcvs.git
cd wcvs
docker build -t wcvs .
docker run -it wcvs /wcvs --help
```

---

## Usage

```bash
wcvs -u https://example.com
```

Two wordlists are required for header and parameter poisoning techniques — one for headers and one for parameters. Place them in the same directory as WCVS, or specify them with `--headerwordlist/-hw` and `--parameterwordlist/-pw`:

```bash
wcvs -u https://example.com -hw "file:wordlists/header_wordlist.txt" -pw "file:wordlists/parameter_wordlist.txt"
```

### Specify Headers, Parameters, Cookies, and More

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

```bash
wcvs -u https://example.com -sc "PHPSESSID=123"
wcvs -u https://example.com -sh "Referer: localhost"
wcvs -u https://example.com -post -sb "admin=true"
```

### Generate a JSON Report

```bash
wcvs -u https://example.com -gr
wcvs -u https://example.com -gr -gp /home/user/Documents
```

### Crawl for URLs

```bash
wcvs -u https://example.com -r 5
wcvs -u https://example.com -r 5 -rl 2
```

### Use a Proxy

```bash
wcvs -u https://example.com -up
wcvs -u https://example.com -up -purl http://127.0.0.1:8081
```

### Throttle or Accelerate

```bash
wcvs -u https://example.com -rr 10
wcvs -u https://example.com -t 50
```

Run `wcvs -h` for the full list of flags.

---

## Background Information

### Web Cache Poisoning
Web cache poisoning exploits discrepancies between what a cache stores and what it actually serves to users. By injecting malicious content via unkeyed inputs (headers, parameters, etc.), an attacker can poison a cached response that is then delivered to every subsequent visitor.

### Web Cache Deception
Web cache deception tricks a cache into storing a response to a sensitive, authenticated resource by disguising the request URL as a public static file.

### Further Reading
1. [Is Your Application Vulnerable to Web Cache Poisoning?](https://www.hackmanit.de/en/blog-en/142-is-your-application-vulnerable-to-web-cache-poisoning)
2. [Web Cache Vulnerability Scanner (WCVS) - Free, Customizable, Easy-To-Use](https://www.hackmanit.de/en/blog-en/145-web-cache-vulnerability-scanner-wcvs-free-customizable-easy-to-use)

---

## License

WCVS is developed by [Hackmanit](https://hackmanit.de) and [Maximilian Hildebrand](https://www.github.com/m10x), forked by [Spondon16](https://github.com/Spondon16), and licensed under the [Apache License, Version 2.0](LICENSE).
