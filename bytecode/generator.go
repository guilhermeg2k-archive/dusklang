package bytecode

import (
	"bufio"
	"os"

	"github.com/guilhermeg2k/dusklang/ast"
)

func generateByteCode(program *ast.Program) error {
	file, err := os.Create("data/program.dskbc")
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	/* var consts vm.Consts
	var labels vm.Labels
	storageCounter := 0
	constCounter := 0
	labelCounter := 0 */
	function := program.Functions[0]
	for _, statement := range function.Statements {
		switch statement.Type {
		case "FullVarDeclaration":
			for _, variable := range statement.Statement.(ast.FullVarDeclaration).Variables {
				if variable.Expression != nil {
					switch variable.Type {
					case "int":

					}
				}
			}
		}
	}
	return nil
}
