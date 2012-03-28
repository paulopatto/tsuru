package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/timeredbull/tsuru/api/unit"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func Upload(w http.ResponseWriter, r *http.Request) error {
	app := App{Name: r.URL.Query().Get(":name")}
	app.Get()

	if app.Id == 0 {
		http.NotFound(w, r)
	} else {
		f, _, err := r.FormFile("application")
		if err != nil {
			return err
		}

		releaseName := time.Now().Format("20060102150405")
		zipFile := fmt.Sprintf("/tmp/%s.zip", releaseName)
		zipDir := fmt.Sprintf("/tmp/%s", releaseName)

		newFile, err := os.Create(zipFile)
		if err != nil {
			return err
		}
		out, _ := ioutil.ReadAll(f)
		newFile.Write(out)

		cmd := exec.Command("unzip", zipFile, "-d", zipDir)
		output, err := cmd.Output()
		if err != nil {
			return err
		}
		log.Printf(string(output))

		appDir := "/home/application"
		currentDir := appDir + "/releases/current"
		gunicorn := appDir + "/env/bin/gunicorn_django"
		releasesDir := appDir + "/releases"
		releaseDir := releasesDir + "/" + releaseName

		u := unit.Unit{Name: app.Name}
		err = u.SendFile(zipDir, releaseDir)
		if err != nil {
			return err
		}
		//u.Command(fmt.Sprintf("'rm -rf %s'", currentDir))
		output, err = u.Command(fmt.Sprintf("cd %s && ln -nfs %s current", releasesDir, releaseName))
		log.Printf(string(output))
		if err != nil {
			return err
		}
		output, err = u.Command("sudo killall gunicorn_django")
		log.Printf(string(output))
		if err != nil {
			return err
		}
		output, err = u.Command(fmt.Sprintf("cd %s && sudo %s --daemon --workers=3 --bind=127.0.0.1:8888", currentDir, gunicorn))
		log.Printf(string(output))
		if err != nil {
			return err
		}

		fmt.Fprint(w, "success")
	}
	return nil
}

func AppDelete(w http.ResponseWriter, r *http.Request) error {
	app := App{Name: r.URL.Query().Get(":name")}
	app.Destroy()
	fmt.Fprint(w, "success")
	return nil
}

func AppList(w http.ResponseWriter, r *http.Request) error {
	apps, err := AllApps()
	if err != nil {
		return err
	}

	b, err := json.Marshal(apps)
	if err != nil {
		return err
	}
	fmt.Fprint(w, bytes.NewBuffer(b).String())
	return nil
}

func AppInfo(w http.ResponseWriter, r *http.Request) error {
	app := App{Name: r.URL.Query().Get(":name")}
	app.Get()

	if app.Id == 0 {
		http.NotFound(w, r)
	} else {
		b, err := json.Marshal(app)
		if err != nil {
			return err
		}
		fmt.Fprint(w, bytes.NewBuffer(b).String())
	}
	return nil
}

func CreateAppHandler(w http.ResponseWriter, r *http.Request) error {
	var app App

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &app)
	if err != nil {
		return err
	}

	err = app.Create()
	if err != nil {
		return err
	}
	fmt.Fprint(w, "success")
	return nil
}
