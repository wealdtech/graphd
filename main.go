// Copyright © 2022 Weald Technology Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net/http"

	// #nosec G108
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	zerologger "github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	restdaemon "github.com/wealdtech/graphd/services/daemon/rest"
	"github.com/wealdtech/graphd/services/generator"
	matplotlibgenerator "github.com/wealdtech/graphd/services/generator/matplotlib"
	"github.com/wealdtech/graphd/services/metrics"
	nullmetrics "github.com/wealdtech/graphd/services/metrics/null"
	prometheusmetrics "github.com/wealdtech/graphd/services/metrics/prometheus"
)

// ReleaseVersion is the release version for the code.
var ReleaseVersion = "1.0.0"

func main() {
	os.Exit(main2())
}

func main2() int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := fetchConfig(); err != nil {
		zerologger.Error().Err(err).Msg("Failed to fetch configuration")
		return 1
	}

	// runCommands will not return if a command is run.
	runCommands(ctx)

	if err := initLogging(); err != nil {
		log.Error().Err(err).Msg("Failed to initialise logging")
		return 1
	}

	logModules()
	log.Info().Str("version", ReleaseVersion).Msg("Starting graphd")

	if err := initProfiling(); err != nil {
		log.Error().Err(err).Msg("Failed to initialise profiling")
		return 1
	}

	runtime.GOMAXPROCS(runtime.NumCPU() * 8)

	log.Trace().Msg("Starting metrics service")
	monitor, err := startMonitor(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start metrics service")
		return 1
	}
	if err := registerMetrics(ctx, monitor); err != nil {
		log.Error().Err(err).Msg("Failed to register metrics")
		return 1
	}
	setRelease(ctx, ReleaseVersion)
	setReady(ctx, false)

	if err := startServices(ctx, monitor); err != nil {
		log.Error().Err(err).Msg("Failed to initialise services")
		return 1
	}
	setReady(ctx, true)

	log.Info().Msg("All services operational")

	// Wait for signal.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	for {
		sig := <-sigCh
		if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == os.Interrupt || sig == os.Kill {
			break
		}
	}

	log.Info().Msg("Stopping graphd")
	return 0
}

// fetchConfig fetches configuration from various sources.
func fetchConfig() error {
	pflag.String("base-dir", "", "base directory for configuration files")
	pflag.String("log-level", "info", "minimum level of messsages to log")
	pflag.String("log-file", "", "redirect log output to a file")
	pflag.String("profile-address", "", "Address on which to run Go profile server")
	pflag.String("daemon.listen-address", "0.0.0.0:8080", "Address on which the daemon should listen")
	pflag.String("generator.style", "matplotlib", "style of generator (matplotlib)")
	pflag.String("resources-dir", "resources", "location of resources directory")
	pflag.Bool("version", false, "show graphd version and exit")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return errors.Wrap(err, "failed to bind pflags to viper")
	}

	if viper.GetString("base-dir") != "" {
		// User-defined base directory.
		viper.AddConfigPath(resolvePath(""))
		viper.SetConfigName("graphd")
	} else {
		// Home directory.
		home, err := homedir.Dir()
		if err != nil {
			return errors.Wrap(err, "failed to obtain home directory")
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".graphd")
	}

	// Environment settings.
	viper.SetEnvPrefix("GRAPHD")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	// Defaults.
	viper.Set("process-concurrency", 16)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return errors.Wrap(err, "failed to read configuration file")
		}
	}

	return nil
}

// initProfiling initialises the profiling server.
func initProfiling() error {
	profileAddress := viper.GetString("profile-address")
	if profileAddress != "" {
		go func() {
			log.Info().Str("profile_address", profileAddress).Msg("Starting profile server")
			runtime.SetMutexProfileFraction(1)
			if err := http.ListenAndServe(profileAddress, nil); err != nil {
				log.Warn().Str("profile_address", profileAddress).Err(err).Msg("Failed to run profile server")
			}
		}()
	}
	return nil
}

func startServices(ctx context.Context, monitor metrics.Service) error {
	var generator generator.Service
	var err error
	switch viper.GetString("generator.style") {
	case "matplotlib":
		generator, err = matplotlibgenerator.New(ctx,
			matplotlibgenerator.WithResourcesDir(viper.GetString("resources-dir")),
		)
	default:
		err = fmt.Errorf("unknown generator style %s", viper.GetString("generator.style"))
	}
	if err != nil {
		return errors.Wrap(err, "failed to create generator service")
	}

	_, err = restdaemon.New(ctx,
		restdaemon.WithLogLevel(logLevel(viper.GetString("daemon.log-level"))),
		restdaemon.WithGenerator(generator),
		restdaemon.WithListenAddress(viper.GetString("daemon.listen-address")),
	)
	if err != nil {
		return errors.Wrap(err, "failed to start REST daemon")
	}

	return nil
}

func logModules() {
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		log.Trace().Str("path", buildInfo.Path).Msg("Main package")
		for _, dep := range buildInfo.Deps {
			log := log.Trace()
			if dep.Replace == nil {
				log = log.Str("path", dep.Path).Str("version", dep.Version)
			} else {
				log = log.Str("path", dep.Replace.Path).Str("version", dep.Replace.Version)
			}
			log.Msg("Dependency")
		}
	}
}

// resolvePath resolves a potentially relative path to an absolute path.
func resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	baseDir := viper.GetString("base-dir")
	if baseDir == "" {
		homeDir, err := homedir.Dir()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not determine a home directory")
		}
		baseDir = homeDir
	}
	return filepath.Join(baseDir, path)
}

func startMonitor(ctx context.Context) (metrics.Service, error) {
	var monitor metrics.Service
	if viper.Get("metrics.prometheus.listen-address") != nil {
		var err error
		monitor, err = prometheusmetrics.New(ctx,
			prometheusmetrics.WithLogLevel(logLevel(viper.GetString("metrics.prometheus.log-level"))),
			prometheusmetrics.WithAddress(viper.GetString("metrics.prometheus.listen-address")),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start prometheus metrics service")
		}
		log.Info().Str("listen_address", viper.GetString("metrics.prometheus.listen-address")).Msg("Started prometheus metrics service")
	} else {
		log.Debug().Msg("No metrics service supplied; monitor not starting")
		monitor = &nullmetrics.Service{}
	}
	return monitor, nil
}

func runCommands(ctx context.Context) {
	if viper.GetBool("version") {
		fmt.Printf("%s\n", ReleaseVersion)
		os.Exit(0)
	}
}
