package collector

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {

	history, err := CreateHistory("..", "4a89114ba35dd28ed81f11ec3eba769a401789a5", "", false)
	//history, err := CreateHistory("..", "master", "", false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(history)

}
