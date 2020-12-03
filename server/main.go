package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var method string
var path string
var headers map[string]string = map[string]string{}

func main() {
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	// Add timeout - not working
	// go timeout(conn, 10000)

	request(conn)

	respond(conn)
}

func scan(scn *bufio.Scanner) bool {
	return scn.Scan()
}

func request(conn net.Conn) {
	scn := bufio.NewScanner(conn)

	i := 0

	for scan(scn) {
		ln := scn.Text()
		if len(ln) == 0 {
			break
		}

		io.WriteString(os.Stdout, ln+"\r\n")

		if i == 0 {
			method = strings.Fields(ln)[0]
			path = strings.Fields(ln)[1]
			fmt.Fprintf(os.Stdout, "Method: %s\r\n", method)
			fmt.Fprintf(os.Stdout, "Path: %s\r\n", path)
		} else {
			sep := ":"
			headerPart := strings.Split(ln, sep)
			if len(headerPart) >= 2 {
				headers[headerPart[0]] = strings.Join(headerPart[1:], sep)
			}
		}

		i++
	}

}

func respond(conn net.Conn) {
	switch path {
	case "/":
		respondSuccess(conn, "Welcome on home")
	case "/contact":
		respondSuccess(conn, "Contact us now !")
	case "/newsletter-ok":
		respondSuccess(conn, "Your email have been accepted to newsletter")
	case "/api":
		respondApi(conn)
	default:
		respond404(conn)
	}
}

func respondSuccess(conn net.Conn, title string) {

	// for {
	// 	// Make inf loop to test timeout
	// }

	url := fmt.Sprintf("http://%s%s", headers["Host"], path)

	body := `<DOCTYPE html><html lang="fr"><head></head>
	<body>
		<h1>` + title + `</h1>
		<p>Url courante : ` + url + `</p>
		<h3>Menu</h3>
		<ul>
			<li><a href="/">Home</a></li>
			<li><a href="/contact">Contact</a></li>
			<li><a href="/api">Api</a></li>
		</ul>
		<form method="POST" action="/newsletter-ok">
			<input name="email" placeholder="Register to newsletter"/>
			<input type="submit" />
		</form>
	</body></html>`

	io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprintf(conn, "Content-Type: text/html\r\n")
	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, body)

}

func respondApi(conn net.Conn) {
	data := struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}{
		Firstname: "Kévin",
		Lastname:  "José",
	}
	json, _ := json.Marshal(data)
	jsons := string(json)

	io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(jsons))
	fmt.Fprintf(conn, "Content-Type: application/json\r\n")
	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, jsons)
}

func respond404(conn net.Conn) {
	io.WriteString(conn, "HTTP/1.1 404 Not Found\r\n")
	fmt.Fprintf(conn, "\r\n")
}

func timeout(conn net.Conn, ms int) {
	timeToStop := time.Now().Add(time.Duration(ms) * time.Millisecond)
	for {
		if time.Now().After(timeToStop) {
			// fmt.Fprintf(conn, "Connection timed out after %d ms.\r\nCheck that there is no infinite loop in the code.", ms)
			conn.Close()
			break
		}
	}
}
