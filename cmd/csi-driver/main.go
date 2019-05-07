// Copyright © 2019 The OpenEBS Authors
//
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
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/openebs/maya/pkg/csidriver/v1alpha1/config"
	"github.com/openebs/maya/pkg/csidriver/v1alpha1/driver"
	"github.com/openebs/maya/pkg/version"
	"github.com/spf13/cobra"
)

func main() {
	_ = flag.CommandLine.Parse([]string{})
	var driverConfig = config.NewConfig()

	cmd := &cobra.Command{
		Use:   "openebs-csi-driver",
		Short: "openebs-csi-driver",
		Run: func(cmd *cobra.Command, args []string) {
			handle(driverConfig)
		},
	}

	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	cmd.PersistentFlags().StringVar(&driverConfig.RestURL, "url", "", "url")
	cmd.PersistentFlags().StringVar(&driverConfig.NodeID, "nodeid", "node1", "node id")
	cmd.PersistentFlags().StringVar(&driverConfig.Version, "version", "", "Print the version and exit")
	cmd.PersistentFlags().StringVar(&driverConfig.Endpoint, "endpoint", "unix://csi/csi.sock", "CSI endpoint")
	cmd.PersistentFlags().StringVar(&driverConfig.DriverName, "name",
		"openebs-csi.openebs.io", "name of the driver")
	cmd.PersistentFlags().StringVar(&driverConfig.PluginType,
		"plugin", "csi-plugin", "Plugin type controller/node")

	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}
}

func handle(driverConfig *config.Config) {

	if driverConfig.Version == "" {
		driverConfig.Version = version.GetVersion()
	}

	logrus.Infof("%s - %s", version.GetVersion(),
		version.GetGitCommit())

	logrus.Infof("DriverName: %v Plugin: %v EndPoint: %v URL: %v NodeID: %v",
		driverConfig.DriverName, driverConfig.PluginType, driverConfig.Endpoint,
		driverConfig.RestURL, driverConfig.NodeID)
	drvr := driver.New(driverConfig)

	if err := drvr.Run(); err != nil {
		log.Fatalln(err)
	}
	os.Exit(0)

}
