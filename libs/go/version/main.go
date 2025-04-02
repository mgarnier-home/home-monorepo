package version

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"mgarnier11.fr/go/libs/utils"
)

func SetupVersionRoute(router *mux.Router) {

	router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		// TODO : use the golang function to load file at build time
		versionFilePath := utils.GetEnv("VERSION_FILE_PATH", "./version.txt")

		version, err := os.ReadFile(versionFilePath)

		if err != nil {
			http.Error(w, "Version not found", http.StatusNotFound)
		} else {
			w.Write([]byte(version))
		}
	})
}
