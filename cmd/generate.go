package cmd

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"maple/internal/utilities"
	"maple/pkg/generate"
	"strings"
)

const (
	generateInputFormatErrorCodes = "error_codes"

	generateOutputFormatTypescriptTypes = "typescript:types"
)

var generateCommand = cobra.Command{
	Use:     "generate",
	Short:   "Generate types or json from source Go file",
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleGenerateCommand(cmd, args); err != nil {
			panic(err)
		}
	},
}

func buildGenerateCommand() *cobra.Command {
	return &generateCommand
}

func handleGenerateCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		logrus.Fatalf("Least 1 argument is required")
	}

	sourceFile := args[0]
	if sourceFile == "" {
		logrus.Fatalf("Source file is not provided")
	}

	inputFormat := generateInputFormatErrorCodes
	if len(args) > 1 {
		inputFormat = args[1]
	} else {
		logrus.Infof("Using default input format '%s'", inputFormat)
	}

	switch inputFormat {
	case generateInputFormatErrorCodes:
	default:
		logrus.Fatalf("Unknown input format '%s'", inputFormat)
	}

	outputFormat := generateOutputFormatTypescriptTypes
	if len(args) > 2 {
		outputFormat = args[2]
	} else {
		logrus.Infof("Using default output format '%s'", outputFormat)
	}

	switch outputFormat {
	case generateOutputFormatTypescriptTypes:
	default:
		logrus.Fatalf("Unknown output format '%s'", outputFormat)
	}

	decls := generate.GetLiteralConstantDecls(sourceFile, nil)

	if inputFormat == generateInputFormatErrorCodes {
		decls = utilities.Filter(decls, func(decl generate.LiteralConstantDecl) bool {
			return decl.Type.Name == "Code"
		})
	}

	switch outputFormat {
	case generateOutputFormatTypescriptTypes:
		decls = utilities.Filter(decls, func(decl generate.LiteralConstantDecl) bool {
			return decl.Literal != nil
		})
		values := utilities.Map(decls, func(d generate.LiteralConstantDecl) string {
			return d.Literal.Value
		})

		fmt.Println(strings.Join(values, " | "))
	}

	return nil
}
