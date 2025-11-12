// SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: MPL-2.0

package bdpmock

import "testing"

func Test1(t *testing.T) {
	var in = BdpMockCalls{}
	var out = BdpMockCalls{}
	CompareBdpMockCalls(t, out, in)
}
