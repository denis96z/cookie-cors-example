package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	//auth.example.com
	http.HandleFunc(
		"/cookie",
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%+v", *r)

			http.SetCookie(w, &http.Cookie{
				Name:     "xmc",
				Value:    "1234567890",
				Domain:   "example.com",
				HttpOnly: true,
				Expires:  time.Now().Add(time.Hour),
				SameSite: http.SameSiteNoneMode,
			})

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Origin", "http://game.example.com")
			w.Header().Set("Access-Control-Allow-Methods", "POST")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		},
	)

	//game.example.com
	http.HandleFunc(
		"/render",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Auth: fetch + cors</title>
</head>
<body>
	<script>
		function authorize() {
			fetch('http://auth.example.com/cookie', {
					method: 'POST', mode: 'cors', credentials: 'include'
			})
				.then(response => {
					if ( response.status !== 200 )
					{
						console.log('ERR');
						return;
					}
					return response.json();
				})
				.then(jv => {
					console.log(JSON.stringify(jv));
				})
				.catch(err => {
					console.log('FAIL');
				});
		}
	</script>
	<button id="btnAuth" onclick="authorize()">AUTHORIZE</button>
</body>
</html>
`,
			))
		},
	)

	http.ListenAndServe(":80", nil)
}
