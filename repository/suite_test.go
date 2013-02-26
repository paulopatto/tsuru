// Copyright 2013 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"github.com/globocom/config"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	config.Set("git:host", "mygithost")
	config.Set("git:protocol", "http")
	config.Set("git:port", "8090")
	config.Set("git:unit-repo", "/home/application/current")
}

func (s *S) TestGetGitServerPanicsIfTheConfigFileHasNoServer(c *C) {
	oldConfig, err := config.Get("git")
	c.Assert(err, IsNil)
	err = config.Unset("git")
	c.Assert(err, IsNil)
	defer config.Set("git", oldConfig)
	c.Assert(getGitServer, PanicMatches, `key "git:host" not found`)
}
