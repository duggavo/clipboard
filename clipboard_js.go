// Copyright (c) 2024 duggavo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js
// +build js

package clipboard

import (
	"errors"
	"sync"
	"syscall/js"
)

func readAll() (string, error) {
	var wg sync.WaitGroup
	wg.Add(1)
	var result string
	var err error
	promise := js.Global().Get("navigator").Get("clipboard").Call("readText").Call(
		"then",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			result = args[0].String()
			wg.Done()
			return nil
		}),
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			err = errors.New(args[0].String())
			wg.Done()
			return nil
		}),
	)

	if !promise.Truthy() {
		return "", errors.New("javascript did not return valid promise")
	}

	wg.Wait()

	return result, err
}

func writeAll(text string) error {
	var wg sync.WaitGroup
	wg.Add(1)
	var err error
	promise := js.Global().Get("navigator").Get("clipboard").Call("writeText", text).Call(
		"then",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			wg.Done()
			return nil
		}),
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			err = errors.New(args[0].String())
			wg.Done()
			return nil
		}),
	)

	if !promise.Truthy() {
		return errors.New("javascript did not return valid promise")
	}

	wg.Wait()

	return err
}
