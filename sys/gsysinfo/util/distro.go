package util

/*
  This module is copied from github.com/nir0s/distgo, nothing changed here
*/

/*
  Package distgo implements a simple library for identifying the linux
  distribution you're running on and some of its properies.
*/

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

const unixEtcDir string = "/etc"
const osReleaseFileName string = "os-release"

// LinuxDistributionObject is the base struct
type LinuxDistributionObject struct {
	OsReleaseFile     string
	DistroReleaseFile string
	OsReleaseInfo     map[string]string
	LSBReleaseInfo    map[string]string
	DistroReleaseInfo map[string]string
}

// LinuxDistribution instantiates a LinuxDistributionObject and returns it
// after having parsed all relevant information.
func LinuxDistribution(d *LinuxDistributionObject) *LinuxDistributionObject {
	if d == nil {
		d = &LinuxDistributionObject{
			OsReleaseFile: path.Join(unixEtcDir, osReleaseFileName),
		}
	}
	d.OsReleaseInfo = d.GetOSReleaseFileInfo()
	d.LSBReleaseInfo = d.GetLSBReleaseInfo()
	d.DistroReleaseInfo = d.GetDistroReleaseFileInfo()
	return d
}

// GetOSReleaseFileInfo retrieves parsed information from an
// os-release file and returns a map with its key-value's
func (d *LinuxDistributionObject) GetOSReleaseFileInfo() map[string]string {
	defaultMap := make(map[string]string)

	if _, err := os.Stat(d.OsReleaseFile); err == nil {
		content := readFileContents(d.OsReleaseFile)
		return parseOSReleaseFile(content)
	}
	return defaultMap
}

// GetLSBReleaseInfo retrieves parsed information from an
// `lsb_release -a` command and returns a map with its key-value's
func (d *LinuxDistributionObject) GetLSBReleaseInfo() map[string]string {
	defaultMap := make(map[string]string)
	var (
		cmdOut []byte
		err    error
	)
	cmdName := "/usr/bin/lsb_release"
	cmdArgs := []string{"-a"}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to run lsb_release -a", err)
		return defaultMap
	}
	return parseLSBRelease(string(cmdOut))
}

// GetDistroReleaseFileInfo retrieves parsed information from an
// `lsb_release -a` command and returns a map with its key-value's
func (d *LinuxDistributionObject) GetDistroReleaseFileInfo() map[string]string {
	defaultMap := make(map[string]string)

	ignoredBasenames := []string{
		"debian_version",
		"lsb-release",
		"oem-release",
		"os-release",
		"system-release",
	}

	distroFileNamePattern := `(\w+)[-_](release|version)$`
	compiledPattern := regexp.MustCompile(distroFileNamePattern)

	files, _ := ioutil.ReadDir(unixEtcDir)
	for _, f := range files {
		isReleaseFile := compiledPattern.MatchString(f.Name())
		if isReleaseFile {
			matches := compiledPattern.FindAllStringSubmatch(f.Name(), -1)
			releaseFilePath := path.Join(unixEtcDir, f.Name())
			if !stringInSlice(f.Name(), ignoredBasenames) {
				content := readFileContents(releaseFilePath)
				defaultMap = parseDistroReleaseFile(content)
				if _, ok := defaultMap["name"]; ok {
					defaultMap["id"] = matches[0][1]
				}
			}
		}
	}

	return defaultMap
}

// ParseOSReleaseFile parses `/etc/os-release` files
// and returns a map with its key=value's
func parseOSReleaseFile(content string) map[string]string {
	props := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, element := range lines {
		if strings.Contains(element, "=") {
			kv := strings.Split(element, "=")
			if kv[0] == "VERSION" {
				compiledPattern := regexp.MustCompile(`(\(\D+\))|,(\s+)?\D+`)
				codenameFound := compiledPattern.MatchString(kv[1])
				codename := compiledPattern.FindString(kv[1])
				if codenameFound {
					codename = strings.TrimSpace(codename)
					codename = strings.Trim(codename, "()")
					props["codename"] = codename
				} else {
					props["codename"] = ""
				}
			}
			props[strings.ToLower(kv[0])] = strings.Trim(kv[1], "\"")
		}
	}

	return props
}

// ParseLSBRelease parses the contents of the `lsb_release -a` command
// and returns a map with its key=value's
func parseLSBRelease(content string) map[string]string {
	props := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, element := range lines {
		trimmedElement := strings.Trim(element, "\n")
		if strings.Contains(trimmedElement, ":") {
			kv := strings.Split(trimmedElement, ":")
			key := strings.Replace(kv[0], " ", "_", -1)
			key = strings.ToLower(key)
			props[key] = strings.TrimSpace(kv[1])
		}
	}

	return props
}

// ParseDistroReleaseFile parses a distro-specific release/version
// file and returns a map of its data. Not all data is necessarily
// found in each release file and that depends on the distribution
func parseDistroReleaseFile(content string) map[string]string {
	props := make(map[string]string)
	line := strings.Split(content, "\n")[0]

	distroFileContentReversePattern := `(?:[^)]*\)(.*)\()? *(?:STL )?([\d.+\-a-z]*\d) *(?:esaeler *)?(.+)`
	compiledPattern := regexp.MustCompile(distroFileContentReversePattern)
	matches := compiledPattern.FindAllStringSubmatch(reverse(line), -1)
	if len(matches) > 0 {
		groups := matches[0]
		props["name"] = reverse(groups[3])
		props["version_id"] = reverse(groups[2])
		props["codename"] = reverse(groups[1])
	} else if len(line) > 0 {
		props["name"] = strings.TrimSpace(line)
	}

	return props
}

// getOSReleaseAttribute retrives a single attribute from a parsed os-release file
func (d *LinuxDistributionObject) getOSReleaseAttribute(attribute string) string {
	return d.OsReleaseInfo[attribute]
}

// getLSBReleaseAttribute retrives a single attribute from the parsed `lsb_release -a` command
func (d *LinuxDistributionObject) getLSBReleaseAttribute(attribute string) string {
	return d.LSBReleaseInfo[attribute]
}

// getDistroReleaseAttribute retrives a single attribute from a parsed distro release file
func (d *LinuxDistributionObject) getDistroReleaseAttribute(attribute string) string {
	return d.DistroReleaseInfo[attribute]
}

// Name returns the name of the distribution.
// Passing `pretty` as true will return the pretty name
func (d *LinuxDistributionObject) Name(pretty bool) string {
	var name string
	names := []string{
		d.getOSReleaseAttribute("name"),
		d.getLSBReleaseAttribute("distributor_id"),
		d.getDistroReleaseAttribute("name"),
	}

	for _, element := range names {
		if name == "" {
			name = element
		}
	}
	if pretty {
		name = d.getOSReleaseAttribute("pretty_name")
		if name == "" {
			name = d.getLSBReleaseAttribute("description")
		}
		if name == "" {
			name = d.getDistroReleaseAttribute("name")
			version := d.Version(true, false)
			if version != "" {
				name = name + " " + version
			}
		}
	}

	return name
}

// Version returns the version of the distribution
// Passing `pretty` as true will return the pretty name
// Passing `best` as true will return the best and more verbose result found
// between all levels of hierarchy instead of just the first one found.
func (d *LinuxDistributionObject) Version(pretty bool, best bool) string {
	versions := []string{
		d.getOSReleaseAttribute("version_id"),
		d.getLSBReleaseAttribute("release"),
		d.getDistroReleaseAttribute("version_id"),
		parseDistroReleaseFile(d.getOSReleaseAttribute("pretty_name"))["version_id"],
		parseDistroReleaseFile(d.getLSBReleaseAttribute("description"))["version_id"],
	}
	version := ""

	if best {
		for _, element := range versions {
			if strings.Count(element, ".") > strings.Count(version, ".") || version == "" {
				version = element
			}
		}
	} else {
		for _, element := range versions {
			if element != "" {
				version = element
				break
			}
		}
	}
	// && codename
	if pretty && version != "" && d.Codename() != "" {
		version = fmt.Sprintf("%s (%s)", version, d.Codename())
	}

	return version
}

// VersionParts returns three values, one for each part of a distribution's version.
// See `Version` for information on `best`
func (d *LinuxDistributionObject) VersionParts(best bool) (string, string, string) {
	versionStr := d.Version(false, best)
	if versionStr != "" {
		compiledPattern := regexp.MustCompile(`(\d+)\.?(\d+)?\.?(\d+)?`)
		matches := compiledPattern.FindAllStringSubmatch(versionStr, -1)
		if len(matches) > 0 {
			groups := matches[0]
			major, minor, buildNumber := groups[1], groups[2], groups[3]
			return major, minor, buildNumber
		}
	}
	return "", "", ""
}

// MajorVersion returns the major version of the distribution.
// See `Version` for information on `best`
func (d *LinuxDistributionObject) MajorVersion(best bool) string {
	major, _, _ := d.VersionParts(best)
	return major
}

// MinorVersion returns the minor version of the distribution.
// See `Version` for information on `best`
func (d *LinuxDistributionObject) MinorVersion(best bool) string {
	_, minor, _ := d.VersionParts(best)
	return minor
}

// BuildNumber returns the build number of the distribution.
// See `Version` for information on `best`
func (d *LinuxDistributionObject) BuildNumber(best bool) string {
	_, _, buildNumber := d.VersionParts(best)
	return buildNumber
}

// normalizeDistroID normalizes a distribution id found and returns it.
// For example, some earlier versions of RHEL based distros might return `red-hat`
// while new ones return `rhel`. For that reason, if `redhat` is found, we normalize
// it and return `rhel` instead so that you always get the same name.
func normalizeDistroID(id string, normalizationTable map[string]string) string {
	distroID := strings.ToLower(id)
	distroID = strings.Replace(distroID, " ", "_", -1)
	normalizedID := normalizationTable[distroID]
	if normalizedID == "" {
		normalizedID = distroID
	}
	return normalizedID
}

// ID returns the id of the distribution.
func (d *LinuxDistributionObject) ID() string {
	var distroID string

	normalizedOSIDTable := map[string]string{}
	normalizedLSBIDTable := map[string]string{
		"enterpriseenterprise":        "oracle",
		"redhatenterpriseworkstation": "rhel",
	}
	normalizedDistroIDTable := map[string]string{
		"redhat": "rhel",
	}

	distroID = d.getOSReleaseAttribute("id")
	if distroID != "" {
		return normalizeDistroID(distroID, normalizedOSIDTable)
	}
	distroID = d.getLSBReleaseAttribute("distributor_id")
	if distroID != "" {
		return normalizeDistroID(distroID, normalizedLSBIDTable)
	}
	distroID = d.getDistroReleaseAttribute("id")
	if distroID != "" {
		return normalizeDistroID(distroID, normalizedDistroIDTable)
	}

	return distroID
}

// Codename returns the distribution's codename.
// Not every distrubtion has a codename, and sometimes, even if a codename is found
// it will not necessarily really be a valid codename but rather an identifier
// of its architecture (like in OpenSUSE's case which returns `x86_64`).
// Codenames should NEVER be used to identify a distribution's version. That's what
// `Version` is for.
func (d *LinuxDistributionObject) Codename() string {
	var codename string

	codenames := []string{
		d.getOSReleaseAttribute("codename"),
		d.getLSBReleaseAttribute("codename"),
		d.getDistroReleaseAttribute("codename"),
	}

	for _, element := range codenames {
		if codename == "" {
			codename = element
		}
	}

	return codename
}

// Like returns the ID_LIKE field of an os-release file if applicable.
// For example, this will return `rhel fedora` for centos or `arch` for Antergos.
func (d *LinuxDistributionObject) Like() string {
	return d.getOSReleaseAttribute("id_like")
}

type version struct {
	Major       string
	Minor       string
	BuildNumber string
}

// Info encompasses most of the information in this library in a single struct.
type Info struct {
	ID           string
	Version      string
	VersionParts version
	Like         string
	Codename     string
}

// Info returns the Info struct.
func (d *LinuxDistributionObject) Info(pretty bool, best bool) *Info {
	infoStruct := &Info{
		ID:      d.ID(),
		Version: d.Version(pretty, best),
		VersionParts: version{
			Major:       d.MajorVersion(best),
			Minor:       d.MinorVersion(best),
			BuildNumber: d.BuildNumber(best),
		},
		Like:     d.Like(),
		Codename: d.Codename(),
	}
	return infoStruct
}

// stringInSlice returns true if str is in list, otherwise returns false.
func stringInSlice(str string, list []string) bool {
	for _, element := range list {
		if element == str {
			return true
		}
	}
	return false
}

// reverse returns its argument string reversed rune-wise left to right.
func reverse(str string) string {
	r := []rune(str)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// readFileContents returns the string representation of the contents of `filepath`.
func readFileContents(filePath string) string {
	contentBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Print(err)
	}
	return string(contentBytes)
}
