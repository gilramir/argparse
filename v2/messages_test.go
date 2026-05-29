package argparse

// Copyright (c) 2026 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"strings"

	. "gopkg.in/check.v1"
)

// Verify that overriding ArgumentParser.Messages translates the user-facing
// error text, using the built-in Korean messages.
func (s *MySuite) TestKoreanMessages(c *C) {
	// Unknown switch -> NoSuchSwitchFmt
	{
		_, ap := createPTestParser()
		ap.Messages = DefaultMessages_ko
		results := ap.parseArgv([]string{"--nope"})
		c.Assert(results.parseError, NotNil)
		c.Check(results.parseError.Error(), Equals,
			"그런 스위치가 없습니다: --nope")
	}

	// Missing required positional -> ExpectedRequiredArgumentFmt
	{
		_, ap := createPTestParser()
		ap.Messages = DefaultMessages_ko
		ap.Add(&Argument{
			Name:        "PosStringSlice",
			NumArgsGlob: "+",
		})
		results := ap.parseArgv([]string{})
		c.Assert(results.parseError, NotNil)
		c.Check(results.parseError.Error(), Equals,
			"필수 'PosStringSlice' 인수가 필요합니다")
	}

	// Switch given with no following value -> ExpectedValueAfterFmt
	{
		_, ap := createPTestParser()
		ap.Messages = DefaultMessages_ko
		results := ap.parseArgv([]string{"--int1"})
		c.Assert(results.parseError, NotNil)
		c.Check(results.parseError.Error(), Equals,
			"--int1 다음에 값이 필요합니다")
	}

	// Unparseable integer value -> WhileParsingValueForFmt wrapping
	// CannotParseIntegerFmt. The underlying strconv error stays in English,
	// so check that the translated fragments are present.
	{
		_, ap := createPTestParser()
		ap.Messages = DefaultMessages_ko
		results := ap.parseArgv([]string{"--int1", "abc"})
		c.Assert(results.parseError, NotNil)
		msg := results.parseError.Error()
		c.Check(strings.Contains(msg, "값을 구문 분석하는 중"), Equals, true)
		c.Check(strings.Contains(msg, "정수로 변환할 수 없습니다"), Equals, true)
	}
}
