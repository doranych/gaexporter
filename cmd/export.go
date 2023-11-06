/*
Package cmd
Copyright © 2023 Doranych <dartdoran@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/doranych/gaexporter/export"
)

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.PersistentFlags().StringP("input", "i", "",
		"input file. Might be image or text file with GA URLs")
	exportCmd.PersistentFlags().StringP("output", "o", "",
		"output destination. If not set - stdout")
	exportCmd.PersistentFlags().StringP("format", "f", "qr",
		"output format. Might be txt or qr. Default is qr")
}

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports fata from Google Authenticator QR codes or otpauth URIs",
	Long:  `Exports fata from Google Authenticator QR codes or otpauth URIs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputString := cmd.Flag("input").Value.String()
		outputString := cmd.Flag("output").Value.String()

		dest := export.Output{
			Type: export.OutputTypeStdout,
			Dest: cmd.OutOrStdout(),
		}
		if outputString != "" {
			f, err := os.Create(outputString)
			if err != nil {
				return errors.Wrap(err, "failed to create output file")
			}
			defer f.Close()
			dest = export.Output{
				Type: export.OutputTypeFile,
				Dest: f,
			}
		}

		path, err := getInputPath(inputString)
		if err != nil {
			return err
		}

		var input *os.File
		f, err := os.Open(path)
		if err != nil {
			return errors.Wrap(err, "failed to open input file")
		}
		defer f.Close()
		input = f

		if cmd.Flag("format").Value.String() == "txt" {
			return export.RunText(input, dest)
		}
		return export.RunQR(input, dest)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if cmd.Flag("input").Value.String() == "" {
			return errors.New("input file is required")
		}
		switch cmd.Flag("format").Value.String() {
		case "txt":
		case "qr":
			//do nothing
		default:
			return errors.New("format must be txt or qr")
		}
		return nil
	},
}

func getInputPath(inputString string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "failed to get current working directory")
	}
	path := filepath.Clean(inputString)
	path = filepath.Join(wd, path)
	return path, nil
}
