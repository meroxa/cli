package prompt

import (
	"context"
	"errors"

	"github.com/manifoldco/promptui"
)

type Prompt interface {
	Show(ctx context.Context) (interface{}, error)
	IsSkipped() bool
}

func Show(ctx context.Context, prompts []Prompt) error {
	for _, p := range prompts {
		if p.IsSkipped() {
			continue
		}
		_, err := p.Show(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

type StringPrompt struct {
	Label    string
	Default  string
	Validate func(string) error

	Value *string
	Skip  bool
}

func (sp StringPrompt) IsSkipped() bool {
	return sp.Skip
}

func (sp StringPrompt) Show(ctx context.Context) (interface{}, error) {
	p := promptui.Prompt{
		Label:    sp.Label,
		Default:  sp.Default,
		Validate: sp.Validate,
	}

	val, err := p.Run()
	if err != nil {
		return nil, err
	}
	if sp.Value != nil {
		*sp.Value = val
	}
	return val, nil
}

type BoolPrompt struct {
	Label string
	Value *bool
	Skip  bool
}

func (bp BoolPrompt) IsSkipped() bool {
	return bp.Skip
}

func (bp BoolPrompt) Show(ctx context.Context) (interface{}, error) {
	p := promptui.Prompt{
		Label:     bp.Label,
		IsConfirm: true,
	}

	_, err := p.Run()

	var value bool
	switch {
	case errors.Is(err, promptui.ErrAbort):
		value = false
	case err != nil:
		return nil, err
	default:
		value = true
	}

	if bp.Value != nil {
		*bp.Value = value
	}
	return value, nil
}

type ConditionalPrompt struct {
	If   BoolPrompt
	Then Prompt
	Skip bool
}

func (cp ConditionalPrompt) IsSkipped() bool {
	return cp.Skip
}

func (cp ConditionalPrompt) Show(ctx context.Context) (interface{}, error) {
	val, err := cp.If.Show(ctx)
	if err != nil {
		return nil, err
	}
	if !val.(bool) {
		return nil, nil
	}
	return cp.Then.Show(ctx)
}

type MapPrompt struct {
	Label string

	KeyPrompt   Prompt
	ValuePrompt Prompt
	NextLabel   string

	Value map[interface{}]interface{}
	Skip  bool
}

func (mp MapPrompt) Show(ctx context.Context) (interface{}, error) {
	value := make(map[interface{}]interface{})
	hasNext := true

	next := BoolPrompt{
		Label: mp.NextLabel,
	}

	for hasNext {
		key, err := mp.KeyPrompt.Show(ctx)
		if err != nil {
			return nil, err
		}

		val, err := mp.ValuePrompt.Show(ctx)
		if err != nil {
			return nil, err
		}

		value[key] = val

		tmp, err := next.Show(ctx)
		if err != nil {
			return nil, err
		}
		hasNext = tmp.(bool)
	}

	if mp.Value != nil {
		for k, v := range value {
			mp.Value[k] = v
		}
	}
	return value, nil
}

func (mp MapPrompt) IsSkipped() bool {
	return mp.Skip
}
