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
	handlers := []func(w http.ResponseWriter, r *http.Request) error{
		circleCIBuild,
		dockerBuild,
	}

	errs := make([]error, 0)
	w.Header().Set("", "")
	for _, f := range handlers {
		err := f(w, r)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		for _, v := range errs {
			w.Write([]byte(v.Error()))
		}
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func circleCIBuild(w http.ResponseWriter, r *http.Request) error {
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
	return err
}

func dockerBuild(w http.ResponseWriter, r *http.Request) error {
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
	return err
}
