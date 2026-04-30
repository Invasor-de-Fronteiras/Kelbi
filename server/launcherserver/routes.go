package launcherserver

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/julienschmidt/httprouter"
	"erupe-ce/config"
)

func serverList(s *Server, w http.ResponseWriter, r *http.Request) {
	servers := s.erupeConfig.Launcher.Servers
	if len(servers) == 0 {
		servers = []config.LauncherServer{{
			Name: "Kelbi",
			IP:   s.erupeConfig.Host,
			Port: s.erupeConfig.Sign.Port,
		}}
	}

	var groups string
	for i, srv := range servers {
		ip := srv.IP
		if ip == "" {
			ip = s.erupeConfig.Host
		}
		groups += fmt.Sprintf("<group idx='%d' nam='%s' ip='%s' port='%d'/>", i, srv.Name, ip, srv.Port)
	}

	fmt.Fprintf(w, `<?xml version="1.0"?><server_groups>%s</server_groups>`, groups)
}

func serverUniqueName(w http.ResponseWriter, r *http.Request) {
	// TODO(Andoryuuta): Implement checking for unique character name.
	fmt.Fprintf(w, `<?xml version="1.0" encoding="ISO-8859-1"?><uniq code="200">OK</uniq>`)
}

func jpLogin(w http.ResponseWriter, r *http.Request) {
	// HACK(Andoryuuta): Return the given password back as the `skey` to defer the login logic to the sign server.
	resultJSON := fmt.Sprintf(`{"result": "Ok", "skey": "%s", "code": "000", "msg": ""}`, r.FormValue("pw"))

	fmt.Fprintf(w,
		`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
		<html>
		<body onload="doPost();">
		<script type="text/javascript">
		function doPost(){
			parent.postMessage(document.getElementById("result").getAttribute("value"), "http://cog-members.mhf-z.jp");
		}
		</script>
		<input id="result" value="%s"/>
		</body>
		</html>`, html.EscapeString(resultJSON))

}

func (s *Server) setupServerlistRoutes(r *mux.Router) {
	// TW
	r.HandleFunc("/server/unique.php", serverUniqueName) // Name checking is also done on this host.
	r.Handle("/server/serverlist.xml", ServerHandlerFunc{s, serverList})

	// JP
	r.Handle("/serverlist.xml", ServerHandlerFunc{s, serverList})

}

func (s *Server) setupOriginalLauncherRotues(r *mux.Router) {
	// TW
	twMain := r.Host("mhfg.capcom.com.tw").Subrouter()
	twMain.PathPrefix("/").Handler(http.FileServer(http.Dir("./www/tw/")))

	// JP
	jpMain := r.Host("cog-members.mhf-z.jp").Subrouter()
	jpMain.PathPrefix("/").Handler(http.FileServer(http.Dir("./www/jp/")))

	// JP Launcher does additional auth over HTTP that the TW launcher doesn't.
	jpAuth := r.Host("www.capcom-onlinegames.jp").Subrouter()
	jpAuth.HandleFunc("/auth/launcher/login", jpLogin) //.Methods("POST")
	jpAuth.PathPrefix("/auth/").Handler(http.StripPrefix("/auth/", http.FileServer(http.Dir("./www/jp/auth/"))))

}

func (s *Server) setupCustomLauncherRotues(r *mux.Router) {

	// TW
	r.PathPrefix("/g6_launcher/").Handler(http.StripPrefix("/g6_launcher/", http.FileServer(http.Dir(s.erupeConfig.Launcher.Path))))
	r.PathPrefix("/version").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.erupeConfig.PatchServers.En.PatchServerManifest == "" {
			w.WriteHeader(404)
			return
		}

		result := s.getPatchServerVersion(s.erupeConfig.PatchServers.En.PatchServerManifest)
		// nolint:errcheck
		w.Write(result)
	})

	// JP
	r.PathPrefix("/launcher/").Handler(http.StripPrefix("/launcher/", http.FileServer(http.Dir(s.erupeConfig.Launcher.Path))))
	r.PathPrefix("/version").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.erupeConfig.PatchServers.Jp.PatchServerManifest == "" {
			w.WriteHeader(404)
			return
		}

		result := s.getPatchServerVersion(s.erupeConfig.PatchServers.Jp.PatchServerManifest)
		// nolint:errcheck
		w.Write(result)
	})
}

func (s *Server) getPatchServerVersion(serverUrl string) []byte {
	resp, err := http.Get(fmt.Sprintf("http://%s/version", serverUrl))
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	var payload struct {
		Version string `json:"version"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil
	}

	return []byte(payload.Version)
}
