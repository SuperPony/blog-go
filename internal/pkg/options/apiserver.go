/*
 * Copyright 2021 SuperPony <superponyyy@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package options

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

type APIServerOptions struct {
	Mode        string   `json:"mode"`
	Middlewares []string `json:"middlewares"`
}

func NewServerOptions() *APIServerOptions {
	return &APIServerOptions{
		Mode:        gin.ReleaseMode,
		Middlewares: []string{},
	}
}

func (o *APIServerOptions) Validate() []error {
	return []error{}
}

func (o *APIServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Mode, "server.mode", o.Mode, "Server mode. supported mode: debug|test|release")

	fs.StringSliceVar(&o.Middlewares, "server.middlewares", o.Middlewares, "Server middlewares")
}
