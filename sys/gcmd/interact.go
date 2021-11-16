package gcmd

import (
	"github.com/manifoldco/promptui"
)

/*
"github.com/abiosoft/ishell" 交互
github.com/jroimartin/gocui
*/

// github.com/hypebeast/vspark
// termui
// https://github.com/gosuri/uiprogress
// https://github.com/cheggaaa/pb

// 参考https://github.com/ddo/fast/blob/master/fast.go的滚动显示

type (
	Interact struct {
		pt promptui.Prompt
	}
)

func NewInteract() *Interact {
	return &Interact{pt: promptui.Prompt{}}
}

func (t *Interact) Select(first string, others ...string) (string, error) {
	prompt := promptui.Select{
		Label: "Select Day",
		Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
			"Saturday", "Sunday"},
	}

	_, result, err := prompt.Run()
	return result, err
}