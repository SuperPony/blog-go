/*
 * Copyright 2021 SuperPony <superponyyy@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package options

import (
	"encoding/json"

	genericoptions "blog-api/internal/pkg/options"
	"blog-api/pkg/cli/flag"
)

type Options struct {
	ServerRunOptions       *genericoptions.ServerRunOptions       `json:"server" mapstructure:"server"`
	InsecureServingOptions *genericoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	FeatureOptions         *genericoptions.FeatureOptions         `json:"feature" mapstructure:"feature"`
	MySQLOptions           *genericoptions.MySQLOptions           `json:"db" mapstructure:"db"`
	RedisOptions           *genericoptions.RedisOptions           `json:"redis" mapstructure:"redis"`
}

func NewOptions() *Options {
	return &Options{
		ServerRunOptions:       genericoptions.NewServerRunOptions(),
		InsecureServingOptions: genericoptions.NewInsecureServingOptions(),
		FeatureOptions:         genericoptions.NewFeatureOptions(),
		MySQLOptions:           genericoptions.NewMySQLOptions(),
		RedisOptions:           genericoptions.NewRedisOptions(),
	}
}

func (o *Options) Flags() (fss flag.NamedFlagSets) {
	o.ServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.InsecureServingOptions.AddFlags(fss.FlagSet("insecure serving"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.MySQLOptions.AddFlags(fss.FlagSet("db"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	return fss
}

// Complete 设置需要默认值的选项
func (o Options) Complete() error {
	return nil
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}