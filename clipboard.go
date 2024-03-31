// Copyright (c) 2024 duggavo.
// Copyright (c) 2013 Ato Araki. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package clipboard read/write on clipboard
package clipboard

import "time"

// ReadAll read string from clipboard
func ReadAll() (string, error) {
	return readAll()
}

// WriteAll write string to clipboard
func WriteAll(text string) error {
	return writeAll(text)
}

// Unsupported might be set true during clipboard init, to help callers decide
// whether or not to offer clipboard options.
var Unsupported bool

// Monitor starts monitoring the clipboard for changes. When the clipboard
// content changes, the content is sent in the changes channel
func Monitor(interval time.Duration, changes chan string, stopit chan struct{}) error {
	defer close(changes)

	currentValue, err := ReadAll()
	if err != nil {
		return err
	}

	for {
		select {
		case <-stopit:
			return nil
		default:
			newValue, err := ReadAll()
			if err != nil {
				return err
			}

			if newValue != currentValue {
				currentValue = newValue
				changes <- currentValue
			}
		}
		time.Sleep(interval)
	}
}
