package version

import (
	"net/http"
	"os"
	"regexp"
	"strconv"

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

var semverRegexp = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)$`)

type SemVer struct {
	Major int
	Minor int
	Patch int
	Raw   string // original tag value
}

func ParseSemver(s string) (SemVer, bool) {
	m := semverRegexp.FindStringSubmatch(s)
	if m == nil {
		return SemVer{}, false
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])
	return SemVer{Major: major, Minor: minor, Patch: patch, Raw: s}, true
}

func (a SemVer) NewerThan(b SemVer) bool {
	if a.Major != b.Major {
		return a.Major > b.Major
	}
	if a.Minor != b.Minor {
		return a.Minor > b.Minor
	}
	return a.Patch > b.Patch
}
