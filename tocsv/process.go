package tocsv

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/wzshiming/easycel"
)

type Source interface {
	Values(fields []string) []string
}

func ProcessToCSV[T Source](out io.Writer, condition string, fields []string, ch <-chan T) error {
	output := csv.NewWriter(out)
	defer output.Flush()

	var program *cel.Program

	if condition != "" {
		registry := easycel.NewRegistry("ext")

		var t T
		_ = registry.RegisterType(t)
		_ = registry.RegisterVariable("self", t)

		env, err := easycel.NewEnvironment(cel.Lib(registry))
		if err != nil {
			return err
		}

		p, err := env.Program(condition)
		if err != nil {
			return err
		}
		program = &p
	}

	err := output.Write(fields)
	if err != nil {
		return err
	}
	for entry := range ch {
		if program != nil {
			rev, _, err := (*program).Eval(map[string]any{
				"self": entry,
			})
			if err != nil {
				return err
			}

			if match, ok := rev.(types.Bool); !ok {
				return fmt.Errorf("condition is not a bool")
			} else if !match {
				continue
			}
		}

		err = output.Write(entry.Values(fields))
		if err != nil {
			return err
		}
	}
	return nil
}
