// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package target

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateTargetFromHostname(t *testing.T) {

	t.Parallel()

	Convey("A target created", t, func() {
		fl := NewFlags()
		fl.SysInfo = &SysInfo{Hostname: "test-c4", Region: "test-region"}

		Convey("for a device with autogenerated hostname should have autogen: hostname prefix", func() {
			fl.TargetType = DeviceType
			fl.AutoGenHostname = true
			fl.SetDefaultsFromHostname()
			target, err := NewFromFlags(&fl)
			So(err, ShouldBeNil)
			So(target, ShouldHaveSameTypeAs, (*NetworkDevice)(nil))
			nd := target.(*NetworkDevice)
			So(nd.AsProto().GetHostname(), ShouldEqual, "autogen:test-c4")
			So(nd.AsProto().GetHostgroup(), ShouldEqual, "4")
		})
		Convey("for a device with a static hostname should not have a prefix", func() {
			fl.TargetType = DeviceType
			fl.SetDefaultsFromHostname()
			target, err := NewFromFlags(&fl)
			So(err, ShouldBeNil)
			So(target, ShouldHaveSameTypeAs, (*NetworkDevice)(nil))
			nd := target.(*NetworkDevice)
			So(nd.AsProto().GetHostname(), ShouldEqual, "test-c4")
			So(nd.AsProto().GetHostgroup(), ShouldEqual, "4")
		})
		Convey("for a task with autogenerated hostname should have autogen: hostname prefix", func() {
			fl.TargetType = TaskType
			fl.TaskServiceName = "test-service"
			fl.TaskJobName = "test-job"
			fl.AutoGenHostname = true
			fl.SetDefaultsFromHostname()
			target, err := NewFromFlags(&fl)
			So(err, ShouldBeNil)
			So(target, ShouldHaveSameTypeAs, (*Task)(nil))
			task := target.(*Task)
			So(task.AsProto().GetHostName(), ShouldEqual, "autogen:test-c4")
		})
		Convey("for a task with a static hostname should not have a prefix", func() {
			fl.TargetType = TaskType
			fl.TaskServiceName = "test-service"
			fl.TaskJobName = "test-job"
			fl.SetDefaultsFromHostname()
			target, err := NewFromFlags(&fl)
			So(err, ShouldBeNil)
			So(target, ShouldHaveSameTypeAs, (*Task)(nil))
			task := target.(*Task)
			So(task.AsProto().GetHostName(), ShouldEqual, "test-c4")
		})
	})
}