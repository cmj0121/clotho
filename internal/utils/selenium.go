// Get the global utility for the internal packages.
package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/tebeka/selenium"
)

type Selenium struct {
	// The options for the Selenium client.
	Driver        string `group:"selenium" help:"The path to the ChromeDriver." default:"chromedriver"`
	Port          int    `group:"selenium" help:"The bind address for the ChromeDriver." default:"9515"`
	DriverVersion string `group:"selenium" help:"The version of the ChromeDriver." default:""`

	// The Selenium service for the ChromeDriver.
	*selenium.Service `kong:"-"`

	// The Selenium web driver for the ChromeDriver.
	selenium.WebDriver `kong:"-"`

	// The internal HTTP client.
	client Client `kong:"-"`
}

// open the necessary resources for the HTTP client.
func (s *Selenium) Prologue() {
	s.client.Prologue()

	if err := s.setupDriver(); err != nil {
		log.Error().Err(err).Msg("failed to setup the ChromeDriver")
		panic(err)
	}

	opts := []selenium.ServiceOption{}

	service, err := selenium.NewChromeDriverService(s.Driver, s.Port, opts...)
	if err != nil {
		log.Error().Err(err).Msg("failed to start the ChromeDriver service")
		panic(err)
	}

	s.Service = service

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", s.Port))
	if err != nil {
		log.Error().Err(err).Msg("failed to start the ChromeDriver")
		panic(err)
	}

	s.WebDriver = wd
}

// clean up the resources for the HTTP client.
func (s *Selenium) Epilogue() {
	s.client.Epilogue()

	if s.WebDriver != nil {
		s.WebDriver.Quit()
		s.WebDriver = nil
	}

	if s.Service != nil {
		s.Service.Stop()
		s.Service = nil
	}
}

// check the driver exists or not, and download it if necessary.
func (s *Selenium) setupDriver() (err error) {
	if _, err = os.Stat(s.Driver); os.IsNotExist(err) {
		if err = s.downloadDriver(); err != nil {
			log.Error().Err(err).Msg("failed to download the ChromeDriver")
			return
		}
	}

	return
}

// get the latest driver version from the official site.
func (s *Selenium) getLatestDriverVersion() (ver string, err error) {
	log.Debug().Str("driver", s.Driver).Msg("downloading the ChromeDriver")

	// download the driver from the official site
	//   - https://chromedriver.storage.googleapis.com/LATEST_RELEASE_${MAJOR}
	//   - https://chromedriver.storage.googleapis.com/${VERSION}/chromedriver_linux64.zip

	// get the latest chromedriver version from CLI option
	url := "https://chromedriver.storage.googleapis.com/LATEST_RELEASE"
	if s.DriverVersion != "" {
		log.Info().Str("version", s.DriverVersion).Msg("use the specific version")
		url = fmt.Sprintf("%v_%v", url, s.DriverVersion)
	}

	var resp *http.Response

	resp, err = s.client.Get(url)
	if err != nil {
		log.Warn().Err(err).Str("url", url).Msg("failed to get the latest version")
		return
	}

	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)

	// it should be XXXX.XXXX.XXXX.XXXX
	ver = string(data)
	if match, _ := regexp.MatchString(`^\d+\.\d+\.\d+\.\d+$`, ver); !match {
		log.Warn().Str("version", s.DriverVersion).Msg("invalid driver version")
		return
	}

	return
}

// download the driver from the official site.
func (s *Selenium) downloadDriver() (err error) {
	var latest_ver string

	latest_ver, err = s.getLatestDriverVersion()
	if err != nil {
		log.Error().Err(err).Msg("failed to get the latest version")
		return
	}

	var platform string

	switch runtime.GOOS {
	case "linux":
		platform = "linux64"
	case "darwin":
		if platform = "mac64"; runtime.GOARCH != "amd64" {
			log.Info().Str("arch", runtime.GOARCH).Msg("use the arm64 version")
			platform = "mac_arm64"
		}
	default:
		log.Warn().Str("os", runtime.GOOS).Msg("unsupported platform")
		err = errors.New("unsupported platform")
		return
	}

	url := fmt.Sprintf("https://chromedriver.storage.googleapis.com/%v/chromedriver_%v.zip", latest_ver, platform)
	log.Debug().Str("url", url).Msg("downloading the ChromeDriver")

	var resp *http.Response
	resp, err = s.client.Get(url)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("failed to download the ChromeDriver")
		return
	}

	defer resp.Body.Close()

	dir := path.Dir(s.Driver)
	if err = os.MkdirAll(dir, 0755); err != nil {
		log.Error().Err(err).Str("dir", dir).Msg("failed to create the directory")
		return
	}

	var file *os.File
	file, err = os.OpenFile(s.Driver, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Error().Err(err).Str("file", s.Driver).Msg("failed to create the file")
		return
	}

	defer file.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		log.Error().Err(err).Str("file", s.Driver).Msg("failed to write the file")
		return
	}

	return
}
