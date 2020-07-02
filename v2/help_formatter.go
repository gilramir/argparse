// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

package argparse

import (
	//	"fmt"
	"github.com/gilramir/unicodemonowidth"
	"strings"
)

type rowData struct {
	// the various forms of a single option
	lhs []string
	// the description
	rhs []*unicodemonowidth.PrintedWord
}

type helpFormatter struct {
	rows []rowData

	leftPadding      int
	optionWidth      int
	middlePadding    int
	descriptionWidth int
}

func (s *helpFormatter) addOption(lhs []string, rhs string) {
	s.rows = append(s.rows, rowData{
		lhs: lhs,
		rhs: unicodemonowidth.WhitespaceSplit(rhs),
	})
}

func (s *helpFormatter) produceString(width int) string {
	s.analyze(width)

	/*
		fmt.Println("leftPadding", s.leftPadding)
		fmt.Println("optionWidth", s.optionWidth)
		fmt.Println("middlePadding", s.middlePadding)
		fmt.Println("descriptionWidth", s.descriptionWidth)
	*/

	textRows := s.makeOptionRows(width)

	text := ""
	for _, row := range textRows {
		text += row + "\n"
	}
	return text
}

func (s *helpFormatter) makeOptionRows(width int) []string {
	textRows := make([]string, 0, len(s.rows))
	leftPad := strings.Repeat(" ", s.leftPadding)
	middlePad := strings.Repeat(" ", s.middlePadding)
	emptyOption := strings.Repeat(" ", s.optionWidth)

	for _, pair := range s.rows {
		// lhs
		lhsRows := make([]string, 0, 1)
		lhs := ""
		lhsWidth := 0
		for _, lhsItem := range pair.lhs {
			lhsItemWidth := unicodemonowidth.MonoWidth(lhsItem)
			if lhsWidth > 0 {
				if lhsWidth+lhsItemWidth+1 > s.optionWidth {
					lhsRows = append(lhsRows, lhs)
					lhs = lhsItem
					lhsWidth = lhsItemWidth
				} else {
					lhs += "," + lhsItem
					lhsWidth += 1 + lhsItemWidth
				}
			} else {
				lhs = lhsItem
				lhsWidth = lhsItemWidth
			}
		}
		if lhsWidth > 0 {
			lhsRows = append(lhsRows, lhs)
			lhs = ""
			lhsWidth = 0
		}

		// rhs
		rhsRows := unicodemonowidth.WrapPrintedWords(pair.rhs, s.descriptionWidth)

		numRows := len(lhsRows)
		if len(rhsRows) > numRows {
			numRows = len(rhsRows)
		}

		for i := 0; i < numRows; i++ {
			var lhsText string
			var rhsText string
			if i < len(lhsRows) {
				lhsText = lhsRows[i] + strings.Repeat(" ",
					s.optionWidth-unicodemonowidth.MonoWidth(lhsRows[i]))
			} else {
				lhsText = emptyOption
			}

			if i < len(rhsRows) {
				rhsText = rhsRows[i]
			} else {
				rhsText = ""
			}
			textRows = append(textRows, leftPad+lhsText+middlePad+rhsText)
		}
	}

	return textRows
}

func (s *helpFormatter) analyze(width int) {

	var lhsSingleMin int
	var lhsSingleMax int
	var lhsCombinedMax int
	// The biggest single "word" in a description
	var rhsSingleMax int

	for _, pair := range s.rows {
		// LHS
		combinedWidth := 0
		for i, text := range pair.lhs {
			printedWidth := unicodemonowidth.MonoWidth(text)
			// Possibly alter printedWidth, to count a command between
			// versions of an option
			if i < len(text)-1 {
				printedWidth += 1
			}
			combinedWidth += printedWidth
			if lhsSingleMin == 0 || printedWidth < lhsSingleMin {
				lhsSingleMin = printedWidth
			}
			if lhsSingleMax == 0 || printedWidth > lhsSingleMax {
				lhsSingleMax = printedWidth
			}
		}
		if lhsCombinedMax == 0 || combinedWidth > lhsCombinedMax {
			lhsCombinedMax = combinedWidth
		}

		// RHS
		for _, pw := range pair.rhs {
			if pw.Width > rhsSingleMax {
				rhsSingleMax = pw.Width
			}
		}
	}

	/*
		fmt.Println("width", width)
		fmt.Println("lhsSingleMin", lhsSingleMin)
		fmt.Println("lhsSingleMax", lhsSingleMax)
		fmt.Println("lhsCombinedMax", lhsCombinedMax)
		fmt.Println("rhsSingleMax", rhsSingleMax)
	*/

	// Is the width so narrow as to be impossible? The absolute minimum
	// needed is enough for the lhs single min, the rhs single max, and 1 space between
	absMinNeeded := lhsSingleMin + 1 + rhsSingleMax
	if width < absMinNeeded {
		// It just won't fit. Don't do any formatting
		s.leftPadding = 0
		s.optionWidth = lhsCombinedMax
		s.middlePadding = 1
		s.descriptionWidth = -1
		return
	}

	// Start with the ideal
	s.optionWidth = lhsCombinedMax
	s.leftPadding = width / 10
	s.middlePadding = width / 20
	rightPadding := width / 40
	s.descriptionWidth = width - lhsCombinedMax - s.leftPadding - s.middlePadding - rightPadding
}
