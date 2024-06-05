package api

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<html>
			<head>
				<title>Find Castles</title>
			</head>
			<body>
				<h1>Find castles soon...</h1>
			</body>
		</html>
	`))
	w.WriteHeader(200)
}
