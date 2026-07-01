package main

type TreeStyle int

const (
	TreeStyleDefault TreeStyle = iota
	TreeStyleRounded
	TreeStyleHeavy
	TreeStyleDouble
	TreeStyleASCII
	TreeStyleDotted
)

var currentTreeStyle = TreeStyleDefault

func SetTreeStyle(style TreeStyle) {
	currentTreeStyle = style
}

func GetTreeStyle() TreeStyle {
	return currentTreeStyle
}

type TreeConnectors struct {
	Branch      string // ├──
	LastBranch  string // └──
	Vertical    string // │
	Space       string // (empty)
	TopLeft     string // ┌
	TopRight    string // ┐
	BottomLeft  string // └
	BottomRight string // ┘
	Horizontal  string // ─
	VerticalHor string // ├
	VerticalVer string // ┤
	Cross       string // ┼
}

var treeConnectors = map[TreeStyle]TreeConnectors{
	TreeStyleDefault: {
		Branch:      "├── ",
		LastBranch:  "└── ",
		Vertical:    "│   ",
		Space:       "    ",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
		Horizontal:  "─",
		VerticalHor: "├",
		VerticalVer: "┤",
		Cross:       "┼",
	},
	TreeStyleRounded: {
		Branch:      "╭── ",
		LastBranch:  "╰── ",
		Vertical:    "│   ",
		Space:       "    ",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
		Horizontal:  "─",
		VerticalHor: "╞",
		VerticalVer: "╡",
		Cross:       "┼",
	},
	TreeStyleHeavy: {
		Branch:      "┏━━ ",
		LastBranch:  "┗━━ ",
		Vertical:    "┃   ",
		Space:       "    ",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
		Horizontal:  "━",
		VerticalHor: "┣",
		VerticalVer: "┫",
		Cross:       "╋",
	},
	TreeStyleDouble: {
		Branch:      "╠══ ",
		LastBranch:  "╚══ ",
		Vertical:    "║   ",
		Space:       "    ",
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
		Horizontal:  "═",
		VerticalHor: "╠",
		VerticalVer: "╣",
		Cross:       "╬",
	},
	TreeStyleASCII: {
		Branch:      "|-- ",
		LastBranch:  "+-- ",
		Vertical:    "|   ",
		Space:       "    ",
		TopLeft:     "+",
		TopRight:    "+",
		BottomLeft:  "+",
		BottomRight: "+",
		Horizontal:  "-",
		VerticalHor: "+",
		VerticalVer: "+",
		Cross:       "+",
	},
	TreeStyleDotted: {
		Branch:      "┆ ┆ ",
		LastBranch:  "┆ ┆ ",
		Vertical:    "┆   ",
		Space:       "    ",
		TopLeft:     "┆",
		TopRight:    "┆",
		BottomLeft:  "┆",
		BottomRight: "┆",
		Horizontal:  "┄",
		VerticalHor: "┆",
		VerticalVer: "┆",
		Cross:       "┆",
	},
}

func GetConnectors() TreeConnectors {
	return treeConnectors[currentTreeStyle]
}