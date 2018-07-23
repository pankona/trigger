package trigger

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func getEnvVar(varName string) (result string) {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if pair[0] == varName {
			return pair[1]
		}
	}
	return ""
}

func createHttpClient(r *http.Request) *http.Client {
	return urlfetch.Client(appengine.NewContext(r))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	circleCIBuild(w, r)
	dockerBuild(w, r)
}

func circleCIBuild(w http.ResponseWriter, r *http.Request) {
	endpoint := "https://circleci.com/api/v1/project/pankona/gomo-simra/tree/master"
	query := url.Values{"circle-token": {getEnvVar("CIRCLECI_API_KEY")}}

	req, _ := http.NewRequest(
		"POST",
		endpoint+"?"+query.Encode(),
		nil,
	)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := createHttpClient(r)
	_, err := client.Do(req)

	w.Header().Set("", "")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(200)
}

func dockerBuild(w http.ResponseWriter, r *http.Request) {
	endpoint := "https://registry.hub.docker.com/u/pankona/gomo-simra/trigger/" + getEnvVar("DOCKERHUB_TRIGGER_TOKEN") + "/"
	query := url.Values{"docker_tag": {"master"}}

	req, _ := http.NewRequest(
		"POST",
		endpoint+"?"+query.Encode(),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	client := createHttpClient(r)
	_, err := client.Do(req)

	w.Header().Set("", "")
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
	}
	w.WriteHeader(200)
}
