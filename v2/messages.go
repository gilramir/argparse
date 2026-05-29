package argparse

// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

// Strings that can be printed out to the user. They can be
// overridden for i18n.
//
// Fields whose name ends in "Fmt" are format strings for fmt.Sprintf or
// fmt.Errorf; the verbs (%s, %v, %w) must be preserved in any translation.
//
// These messages are only the ones shown to the *user* of the program (help
// text and command-line errors). Programmer errors (misuse of the API) are
// reported via panic() and are intentionally not translatable.
type Messages struct {
	// The header for the list of sub-commands in the help output:
	// "Sub-Commands"
	SubCommands string

	// The description for the help options (-h / --help):
	// "See this list of options"
	HelpDescription string

	// ---- value parsing errors ----

	// Error when parsing a boolean.
	// "Cannot convert \"%s\" to a boolean"
	CannotParseBooleanFmt string

	// Error when parsing an int or int64.
	// "Cannot convert \"%s\" to an integer: %w"
	CannotParseIntegerFmt string

	// Error when parsing a float64.
	// "Cannot convert \"%s\" to a float"
	CannotParseFloatFmt string

	// Error when parsing a time.Duration. The first %s is the text, the
	// second %s is the underlying error.
	// "Cannot parse \"%s\" as a time duration: %s"
	CannotParseDurationFmt string

	// The Choices slice is of the wrong type.
	// "Choices should be []%s"
	ChoicesOfWrongTypeFmt string

	// The given value is not a valid choice.
	// "Not a valid choice. Should be one of: %v"
	ShouldBeAValidChoiceFmt string

	// ---- "a value is required" errors ----
	// Shown when a non-boolean switch is given with no value following it.

	NeedStringValue   string // "Need a string value"
	NeedIntValue      string // "Need an int value"
	NeedInt64Value    string // "Need an int64 value"
	NeedFloatValue    string // "Need a float value"
	NeedDurationValue string // "Need a time duration value"
	NeedBoolValue     string // "Need a bool value"
	NeedValue         string // "Need a value" (for custom encoding.TextUnmarshaler types)

	// Error when a custom encoding.TextUnmarshaler type fails to parse. The
	// %s is the text, the %w is the underlying error.
	// "Cannot parse \"%s\": %w"
	CannotParseValueFmt string

	// Reported (via panic) when Choices is set on a custom
	// encoding.TextUnmarshaler argument, which is not supported.
	// "Choices is not supported for a custom (TextUnmarshaler) type"
	ChoicesNotSupportedForCustomType string

	// ---- command-line structure errors ----

	// "Unexpected argument: %s"
	UnexpectedArgumentFmt string

	// "Unexpected positional argument: %s"
	UnexpectedPositionalArgumentFmt string

	// "No such switch: %s"
	NoSuchSwitchFmt string

	// "Expected a value after %s"
	ExpectedValueAfterFmt string

	// "Expected a required '%s' argument"
	ExpectedRequiredArgumentFmt string

	// A required switch was not given on the command-line.
	// "Missing required switch: %s"
	MissingRequiredSwitchFmt string

	// "The %s switch does not take a value"
	SwitchDoesNotTakeValueFmt string

	// Used when a help switch (-h/--help) is given a value.
	// "%s does not accept a value"
	DoesNotAcceptValueFmt string

	// "A switch name cannot begin with '='"
	SwitchNameCannotBeginWithEquals string

	// "'--' is given but there's no positional argument allowed"
	DashDashWithoutPositional string

	// Shown when an empty-string argument is encountered.
	// "<empty string>"
	EmptyArgument string

	// Wraps an error that occurred while parsing a value. The %s is the
	// argument label, the %w is the underlying error.
	// "While parsing value for %s: %w"
	WhileParsingValueForFmt string

	// Wraps an error reported for a specific argument. The %s is the argument
	// label, the %w is the underlying error.
	// "%s argument: %w"
	ArgumentErrorFmt string
}

var DefaultMessages_en = Messages{
	SubCommands:     "Sub-Commands",
	HelpDescription: "See this list of options",

	CannotParseBooleanFmt:  "Cannot convert \"%s\" to a boolean",
	CannotParseIntegerFmt:  "Cannot convert \"%s\" to an integer: %w",
	CannotParseFloatFmt:    "Cannot convert \"%s\" to a float",
	CannotParseDurationFmt: "Cannot parse \"%s\" as a time duration: %s",

	ChoicesOfWrongTypeFmt:   "Choices should be []%s",
	ShouldBeAValidChoiceFmt: "Not a valid choice. Should be one of: %v",

	NeedStringValue:   "Need a string value",
	NeedIntValue:      "Need an int value",
	NeedInt64Value:    "Need an int64 value",
	NeedFloatValue:    "Need a float value",
	NeedDurationValue: "Need a time duration value",
	NeedBoolValue:     "Need a bool value",
	NeedValue:         "Need a value",

	CannotParseValueFmt: "Cannot parse \"%s\": %w",

	ChoicesNotSupportedForCustomType: "Choices is not supported for a custom (TextUnmarshaler) type",

	UnexpectedArgumentFmt:           "Unexpected argument: %s",
	UnexpectedPositionalArgumentFmt: "Unexpected positional argument: %s",
	NoSuchSwitchFmt:                 "No such switch: %s",
	ExpectedValueAfterFmt:           "Expected a value after %s",
	ExpectedRequiredArgumentFmt:     "Expected a required '%s' argument",
	MissingRequiredSwitchFmt:        "Missing required switch: %s",
	SwitchDoesNotTakeValueFmt:       "The %s switch does not take a value",
	DoesNotAcceptValueFmt:           "%s does not accept a value",
	SwitchNameCannotBeginWithEquals: "A switch name cannot begin with '='",
	DashDashWithoutPositional:       "'--' is given but there's no positional argument allowed",
	EmptyArgument:                   "<empty string>",
	WhileParsingValueForFmt:         "While parsing value for %s: %w",
	ArgumentErrorFmt:                "%s argument: %w",
}

var DefaultMessages_ko = Messages{
	SubCommands:     "하위 명령",
	HelpDescription: "이 옵션 목록 보기",

	CannotParseBooleanFmt:  "\"%s\"을(를) 불리언으로 변환할 수 없습니다",
	CannotParseIntegerFmt:  "\"%s\"을(를) 정수로 변환할 수 없습니다: %w",
	CannotParseFloatFmt:    "\"%s\"을(를) 실수로 변환할 수 없습니다",
	CannotParseDurationFmt: "\"%s\"을(를) 시간 기간으로 구문 분석할 수 없습니다: %s",

	ChoicesOfWrongTypeFmt:   "Choices는 []%s 형식이어야 합니다",
	ShouldBeAValidChoiceFmt: "유효한 선택이 아닙니다. 다음 중 하나여야 합니다: %v",

	NeedStringValue:   "문자열 값이 필요합니다",
	NeedIntValue:      "정수 값이 필요합니다",
	NeedInt64Value:    "int64 값이 필요합니다",
	NeedFloatValue:    "실수 값이 필요합니다",
	NeedDurationValue: "시간 기간 값이 필요합니다",
	NeedBoolValue:     "불리언 값이 필요합니다",
	NeedValue:         "값이 필요합니다",

	CannotParseValueFmt: "\"%s\"을(를) 구문 분석할 수 없습니다: %w",

	ChoicesNotSupportedForCustomType: "사용자 정의(TextUnmarshaler) 형식에는 Choices가 지원되지 않습니다",

	UnexpectedArgumentFmt:           "예상치 못한 인수: %s",
	UnexpectedPositionalArgumentFmt: "예상치 못한 위치 인수: %s",
	NoSuchSwitchFmt:                 "그런 스위치가 없습니다: %s",
	ExpectedValueAfterFmt:           "%s 다음에 값이 필요합니다",
	ExpectedRequiredArgumentFmt:     "필수 '%s' 인수가 필요합니다",
	MissingRequiredSwitchFmt:        "필수 스위치가 누락되었습니다: %s",
	SwitchDoesNotTakeValueFmt:       "%s 스위치는 값을 받지 않습니다",
	DoesNotAcceptValueFmt:           "%s은(는) 값을 받지 않습니다",
	SwitchNameCannotBeginWithEquals: "스위치 이름은 '='로 시작할 수 없습니다",
	DashDashWithoutPositional:       "'--'가 주어졌지만 허용되는 위치 인수가 없습니다",
	EmptyArgument:                   "<빈 문자열>",
	WhileParsingValueForFmt:         "%s 값을 구문 분석하는 중: %w",
	ArgumentErrorFmt:                "%s 인수: %w",
}
