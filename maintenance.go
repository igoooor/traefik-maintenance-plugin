package traefik_maintenance_plugin

import (
	"context"
	"net/http"
)

type Config struct {
	Enabled bool `yaml:"enabled"`
}

type Maintenance struct {
	name   string
	next   http.Handler
	config *Config
}

// type ResponseWriter struct {
// 	buffer bytes.Buffer
// 
// 	http.ResponseWriter
// }

type responseWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func CreateConfig() *Config {
	return &Config{
		Enabled: false,
	}
}

// Inform if there are hosts in maintenance
// func Inform(config *Config) {
// 	t := time.NewTicker(time.Second * config.InformInterval)
// 	defer t.Stop()
//
// 	for ; true; <-t.C {
// 		client := http.Client{
// 			Timeout: time.Second * config.InformTimeout,
// 		}
//
// 		req, _ := http.NewRequest(http.MethodGet, config.InformUrl, nil)
// 		res, doErr := client.Do(req)
// 		if doErr != nil {
// 			log.Printf("Inform: %v", doErr) // Don't fatal, just go further
// 			continue
// 		}
//
// 		defer res.Body.Close()
//
// 		decoder := json.NewDecoder(res.Body)
// 		decodeErr := decoder.Decode(&hosts)
// 		if decodeErr != nil {
// 			log.Printf("Inform: %v", decodeErr) // Don't fatal, just go further
// 			continue
// 		}
//
// 		log.Printf("Inform response: %v", hosts)
// 	}
// }

// Get all the client's ips
// func GetClientIps(req *http.Request) []string {
// 	var ips []string
//
// 	if req.RemoteAddr != "" {
// 		ip, _, splitErr := net.SplitHostPort(req.RemoteAddr)
// 		if splitErr != nil {
// 			ip = req.RemoteAddr
// 		}
// 		ips = append(ips, ip)
// 	}
//
// 	forwardedFor := req.Header.Get("X-Forwarded-For")
// 	if forwardedFor != "" {
// 		for _, ip := range strings.Split(forwardedFor, ",") {
// 			ips = append(ips, strings.TrimSpace(ip))
// 		}
// 	}
//
// 	return ips
// }
//
// // Check if one of the client ips has access
// func CheckIpAllowed(req *http.Request, host Host) bool {
// 	for _, ip := range GetClientIps(req) {
// 		for _, allowIp := range host.AllowIps {
// 			if ip == allowIp {
// 				return true
// 			}
// 		}
// 	}
//
// 	return false
// }

// Check if the host is under maintenance
//func CheckIfMaintenance(req *http.Request) bool {
//	for _, host := range hosts {
//		if matched, _ := regexp.Match(host.Regex, []byte(req.Host)); matched {
//			return !CheckIpAllowed(req, host)
//		}
//	}
//
//	return false
//}

//func (rw *ResponseWriter) Header() http.Header {
//	return rw.ResponseWriter.Header()
//}
//
//func (rw *ResponseWriter) Write(bytes []byte) (int, error) {
//	return rw.buffer.Write(bytes)
//}
//
//func (rw *ResponseWriter) WriteHeader(statusCode int) {
//	rw.ResponseWriter.Header().Del("Last-Modified")
//	rw.ResponseWriter.Header().Del("Content-Length")
//
//	rw.ResponseWriter.WriteHeader(http.StatusServiceUnavailable)
//}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	//go Inform(config)

	return &Maintenance{
		name:   name,
		next:   next,
		config: config,
	}, nil
}

func (a *Maintenance) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	//rw := &ResponseWriter{ResponseWriter: w}
	if !a.config.Enabled || a.bypassingHeaders(req) {
		rw := &responseWriter{ResponseWriter: w}
		a.next.ServeHTTP(rw, req)

		return
	}

	// rw := &ResponseWriter{ResponseWriter: w}
	//rw := &ResponseWriter{ResponseWriter: w}

	// a.next.ServeHTTP(rw, req)

	//bytes := []byte{}

	bytes := getMaintenanceTemplate()

	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write(bytes)
}

// Maintenance page templates
func getMaintenanceTemplate() []byte {
	return []byte(`<!doctype html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport"
			content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>Domain Transfer in Progress</title>
	<style>
		body {
			text-align: center;
		}

		h1 {
			font-size: 42px;
		}

		body {
			font: 20px Helvetica, sans-serif;
			color: #333;
		}

		article {
			display: block;
			text-align: left;
			margin: auto;
			max-width: 640px;
			min-width: 320px;
			padding: 10% 32px;
		}

		a {
			color: #0047AA;
			text-decoration: none;
		}

		a:hover {
			text-decoration: underline;
		}
	</style>
</head>
<body>
<article>
	<h1>Maintenance Mode</h1>
	<p>We're currently updating and improving our infrastructure. This website will be back soon!</p>
	<p>Wir sind gerade dabei, unsere Infrastruktur zu aktualisieren und zu verbessern. Diese Website wird bald wieder verfügbar sein!</p>
	<p>Nous sommes en train de mettre à jour et d'améliorer notre infrastructure. Ce site sera bientôt de retour !</p>
</article>
</body>
</html>`)
}

func (a *Maintenance) bypassingHeaders(r *http.Request) bool {
	return r.Header.Get("X-Conteo-Maintenance") == "bypass"
}
