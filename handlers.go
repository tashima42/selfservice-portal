package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/tashima42/selfservice-portal/pangolin"
)

// handleGetHealth returns an [http.HandlerFunc] that responds with the health status of the service.
// It includes the service version, VCS revision, build time, and modified status.
// The service version can be set at build time using the VERSION variable (e.g., 'make build VERSION=v1.0.0').
func handleGetHealth(version string) http.HandlerFunc {
	type responseBody struct {
		Version string `json:"Version"`
		Uptime  string `json:"Uptime"`
	}

	res := responseBody{Version: version}

	up := time.Now()

	return func(w http.ResponseWriter, _ *http.Request) {
		res.Uptime = time.Since(up).String()

		if err := encode[responseBody](w, nil, http.StatusOK, res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleRegisterIP(pangolinClient *pangolin.Pangolin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pangolinRule := pangolin.PangolinRule{
			Action:   pangolin.String("ACCEPT"),
			Match:    pangolin.String("IP"),
			Value:    pangolin.String(r.Header.Get("X-Real-IP")),
			Priority: pangolin.Int(10),
			Enabled:  pangolin.Bool(true),
		}

		resourceID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = pangolinClient.CreateRule(pangolinRule, resourceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Success")
	}
}

func handleHomePage(pangolinClient *pangolin.Pangolin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		resources, err := pangolinClient.GetResources()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		type TemplateVars struct {
			RealIP string
			pangolin.PangolinResources
		}

		templateVars := TemplateVars{
			RealIP:            r.Header.Get("X-Real-IP"),
			PangolinResources: *resources,
		}

		w.Header().Set("Content-Type", "text/html")

		tpl := template.New("home")
		tpl, err = tpl.Parse(homePage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tpl.Execute(w, templateVars); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

const homePage = `
<!DOCTYPE html>
<html lang="en" data-theme="caramellatte">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Selfservice Portal</title>
	<link href="https://cdn.jsdelivr.net/npm/daisyui@5" rel="stylesheet" type="text/css" />
	<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
	<link href="https://cdn.jsdelivr.net/npm/daisyui@5/themes.css" rel="stylesheet" type="text/css" />
	<script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.7/dist/htmx.js" integrity="sha384-yWakaGAFicqusuwOYEmoRjLNOC+6OFsdmwC2lbGQaRELtuVEqNzt11c2J711DeCZ" crossorigin="anonymous"></script>
</head>
<body>
	<div class="navbar bg-neutral text-neutral-content">
		<a class="btn btn-ghost text-xl">Selfservice Portal</a>
		<div class="bg-base-100 w-40 text-base-content flex justify-center rounded-sm">
			<p>{{ .RealIP }}</p>
		</div>
	</div>

	<ul class="list bg-base-100 rounded-box shadow-md"> <li class="p-4 pb-2 text-xs opacity-60 tracking-wide">Resources</li>
		
		{{ range $i, $resource := .Resources }}
		<li class="list-row">
			<div class="text-4xl font-thin opacity-30 tabular-nums">{{ $i }}</div>
			<div class="list-col-grow">
				<div>{{ $resource.Name }}</div>
				<div class="text-xs uppercase font-semibold opacity-60">{{ $resource.FullDomain }}</div>
			</div>
			<button class="btn" id="register-{{ $resource.ResourceID }}" hx-put="/register/{{ $resource.ResourceID }}" hx-trigger="click" hx-swap="this">Register</button>
		</li>
		</li>
		{{ end }}
		
	</ul>

	</body>
</html>
`
