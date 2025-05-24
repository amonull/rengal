package icon

import (
	"github.com/amonull/rengal/color"
	"github.com/amonull/rengal/style"
)

type Icon int

const (
	Lua Icon = iota + 1
	Go
	Fail
	Success
	Question
	Mark
	Downloaded
	Progress
	Search
	Link
)

var icons = map[Icon]*iconDef{
	Lua: {
		emoji:   "🌙",
		nerd:    style.Fg(color.Blue)("\uE620"),
		plain:   style.Fg(color.Blue)("Lua"),
		kaomoji: style.Fg(color.Blue)("(=^･ω･^=)"),
		squares: style.Fg(color.Blue)("◧"),
	},
	Go: {
		emoji:   "🐹",
		nerd:    style.Fg(color.Cyan)("\uE627"),
		plain:   style.Fg(color.Cyan)("Go"),
		kaomoji: style.Fg(color.Cyan)("ʕ •ᴥ• ʔ"),
		squares: style.Fg(color.Cyan)("◨"),
	},
	Fail: {
		emoji:   "💀",
		nerd:    style.Fg(color.Red)("ﮊ"),
		plain:   style.Fg(color.Red)("X"),
		kaomoji: style.Fg(color.Red)("┐('～`;)┌"),
		squares: style.Fg(color.Red)("▨"),
	},
	Success: {
		emoji:   "🎉",
		nerd:    style.Fg(color.Green)("\uF65F "),
		plain:   style.Fg(color.Green)("✓"),
		kaomoji: style.Fg(color.Green)("(ᵔ◡ᵔ)"),
		squares: style.Fg(color.Green)("▣"),
	},
	Mark: {
		emoji:   "🦐",
		nerd:    style.Fg(color.Green)("\uF6D9"),
		plain:   style.New().Bold(true).Foreground(color.Orange).Render("*"),
		kaomoji: style.New().Bold(true).Foreground(color.Red).Render("炎"),
		squares: style.New().Bold(true).Foreground(color.Orange).Render("■"),
	},
	Question: {
		emoji:   "🤨",
		nerd:    style.Fg(color.Yellow)("\uF128"),
		plain:   style.Fg(color.Yellow)("?"),
		kaomoji: style.Fg(color.Yellow)("(￢ ￢)"),
		squares: style.Fg(color.Yellow)("◲"),
	},
	Progress: {
		emoji:   "👾",
		nerd:    style.Fg(color.Blue)("\uF0ED "),
		plain:   style.Fg(color.Blue)("@"),
		kaomoji: style.Fg(color.Blue)("┌( >_<)┘"),
		squares: style.Fg(color.Blue)("◫"),
	},
	Downloaded: {
		emoji:   "📦",
		nerd:    style.Bold("\uF0C5 "),
		plain:   style.New().Bold(true).Faint(true).Render("D"),
		kaomoji: style.Bold("⊂(◉‿◉)つ"),
		squares: style.Bold("◬"),
	},
	Search: {
		emoji:   "🔍",
		nerd:    style.Fg(color.Blue)("\uF002"),
		plain:   style.Fg(color.Blue)("S"),
		kaomoji: style.Fg(color.Blue)("⌐■-■"),
		squares: style.Fg(color.Blue)("◪"),
	},
	Link: {
		emoji:   "🔗",
		nerd:    style.Fg(color.Blue)("\uF0C1"),
		plain:   style.Fg(color.Blue)("L"),
		kaomoji: style.Fg(color.Blue)("⌐■-■"),
		squares: style.Fg(color.Blue)("◪"),
	},
}
