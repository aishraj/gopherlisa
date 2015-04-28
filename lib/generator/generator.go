package generator

import (
	"fmt"
)

func Generate(userName string, tag string, fileName string) {
	//should call communicator with the type of data required
	fmt.Sprintf("The values to generate are %v %v %v", userName, tag, fileName)
}
