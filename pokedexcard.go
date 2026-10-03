package main

// pokedexcard.go renders the "inspect" card: Pokedex data and base stats in
// the left column of a framed table, the krabby sprite in the right one.
//
// Nothing here touches the CLI directly. commandInspect calls exactly one
// function, pokemonCard, which returns a string; the card code itself never
// writes to stdout, never reads os.Args, and never shells out by itself --
// the sprite binary sits behind the SpriteSource interface, and the whole
// card sits behind PokedexCard, so both can be swapped or faked in tests.

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/NedjmeEddine/Pokedex_Go/internal/pokeapi"
)

// PokedexCard abstracts turning a caught Pokemon into a terminal card.
type PokedexCard interface {
	Render(p pokeapi.Pokemon) string
}

// SpriteSource abstracts sprite production. The card builder holds one of
// these instead of calling exec itself, so a missing renderer degrades to a
// stats-only card instead of failing the command.
type SpriteSource interface {
	Sprite(name string) (string, error)
}

// krabbySprite renders sprites through the krabby binary. This struct is the
// only place in the file that touches os/exec.
type krabbySprite struct {
	bin string
}

// validPokemonName guards what we hand to exec: an argument slice means no
// shell is involved, but we still only accept what PokeAPI itself emits.
var validPokemonName = regexp.MustCompile(`^[a-z0-9-]+$`)

func (k krabbySprite) Sprite(name string) (string, error) {
	if k.bin == "" {
		return "", fmt.Errorf("no sprite renderer configured")
	}
	if !validPokemonName.MatchString(name) {
		return "", fmt.Errorf("invalid pokemon name %q", name)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// exec.Command resolves the binary through PATH; a missing krabby just
	// comes back as an error here and the caller falls back to no sprite.
	out, err := exec.CommandContext(ctx, k.bin, "name", name, "--no-title").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// statLine is one base stat of the card's view model.
type statLine struct {
	Label string
	Value int
}

// cardData is the plain-data middle layer: everything the card prints, in
// one struct that knows nothing about JSON tags or the terminal.
type cardData struct {
	Name     string
	ID       int
	HeightM  string
	WeightKg string
	Types    []string
	Stats    []statLine
}

// statLabels maps PokeAPI's stat keys onto the names a player expects.
var statLabels = map[string]string{
	"hp":              "HP",
	"attack":          "Attack",
	"defense":         "Defense",
	"special-attack":  "Sp. Atk",
	"special-defense": "Sp. Def",
	"speed":           "Speed",
}

// newCardData converts a pokeapi.Pokemon into the card's view model. Only
// fields that exist on the struct are read; PokeAPI returns Stats in
// hp/attack/defense/special-attack/special-defense/speed order and that order
// is preserved.
func newCardData(p pokeapi.Pokemon) cardData {
	d := cardData{
		Name:     p.Name,
		ID:       p.ID,
		HeightM:  fmt.Sprintf("%.1f m", float64(p.Height)/10),  // decimetres -> metres
		WeightKg: fmt.Sprintf("%.1f kg", float64(p.Weight)/10), // hectograms -> kilograms
	}
	for _, t := range p.Types {
		d.Types = append(d.Types, t.Type.Name)
	}
	for _, s := range p.Stats {
		label := statLabels[s.Stat.Name]
		if label == "" {
			label = s.Stat.Name
		}
		d.Stats = append(d.Stats, statLine{Label: label, Value: s.BaseStat})
	}
	return d
}

// cardBuilder renders cardData into the framed, two-column card.
type cardBuilder struct {
	sprites SpriteSource
	// minLeft is the floor for the data column width, so short cards keep a
	// consistent shape. Longer lines widen it instead of being truncated.
	minLeft int
}

const (
	barWidth    = 16 // stat bar length in cells
	cardMaxCols = 78 // widest card we will draw before stacking the sprite
	cellPad     = 1  // spaces either side of a cell's contents
)

// Render builds the whole card. Sprite failures are not errors: the card is
// still useful without art, so the fallback is the single-column table.
func (b cardBuilder) Render(p pokeapi.Pokemon) string {
	d := newCardData(p)
	var art []string
	if b.sprites != nil {
		if drawn, err := b.sprites.Sprite(d.Name); err == nil {
			art = spriteLines(drawn)
		}
	}
	return frameCard(b.header(d), b.bodyRows(d), art, b.minLeft)
}

// header is the card's title row: the API's own name plus the dex number.
func (b cardBuilder) header(d cardData) string {
	return fmt.Sprintf("%s  #%04d", displayName(d.Name), d.ID)
}

// bodyRows is the data half of the card, one string per table row.
func (b cardBuilder) bodyRows(d cardData) []string {
	rows := []string{
		fmt.Sprintf("%-13s %s", "Height", d.HeightM),
		fmt.Sprintf("%-13s %s", "Weight", d.WeightKg),
		fmt.Sprintf("%-13s %s", "Types", strings.Join(d.Types, ", ")),
		"",
	}
	if len(d.Stats) == 0 {
		return rows
	}
	maxStat := 1
	for _, s := range d.Stats {
		if s.Value > maxStat {
			maxStat = s.Value
		}
	}
	for _, s := range d.Stats {
		rows = append(rows, fmt.Sprintf("%-13s %3d  %s",
			s.Label, s.Value, statBar(s.Value, maxStat)))
	}
	return rows
}

// statBar draws value/maxStat as a filled/empty bar.
func statBar(value, max int) string {
	filled := value * barWidth / max
	if filled < 1 && value > 0 {
		filled = 1
	}
	if filled > barWidth {
		filled = barWidth
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
}

// frameCard draws the card as a table: header row, rule, then the body rows
// with the sprite in the right-hand column. Widths are counted in visible
// columns, so neither the sprite's truecolor escapes nor its block glyphs
// can push the borders out of alignment.
//
//	┌──────────────┬──────────┐
//	│ Caterpie #10 │          │
//	├──────────────┼──────────┤
//	│ Height 0.3 m │  ▄▀▀ ▄▄  │
//	└──────────────┴──────────┘
func frameCard(head string, body, art []string, minLeft int) string {
	dataW := minLeft
	if w := visibleWidth(head); w > dataW {
		dataW = w
	}
	for _, l := range body {
		if w := visibleWidth(l); w > dataW {
			dataW = w
		}
	}
	artW := 0
	for _, a := range art {
		if w := visibleWidth(a); w > artW {
			artW = w
		}
	}
	// Too wide for the terminal: put the sprite under the data instead of
	// beside it, still inside the same frame. The frame itself costs 3
	// borders + 4 pads on top of the two columns.
	if artW > 0 && dataW+artW+3+4*cellPad > cardMaxCols {
		body = append(append(body, ""), art...)
		art = nil
		for _, l := range body {
			if w := visibleWidth(l); w > dataW {
				dataW = w
			}
		}
		artW = 0
	}
	widths := []int{dataW}
	if artW > 0 {
		widths = append(widths, artW)
	}

	rows := len(body)
	if len(art) > rows {
		rows = len(art)
	}

	var sb strings.Builder
	sb.WriteString(hrule("┌", "┬", "┐", widths))
	sb.WriteString("\n")
	sb.WriteString(frameRow([]string{head, strings.Repeat(" ", artW)}, widths))
	sb.WriteString("\n")
	sb.WriteString(hrule("├", "┼", "┤", widths))
	sb.WriteString("\n")
	for i := 0; i < rows; i++ {
		cells := []string{""}
		if i < len(body) {
			cells[0] = body[i]
		}
		if artW > 0 {
			a := ""
			if i < len(art) {
				a = art[i]
			}
			cells = append(cells, a)
		}
		sb.WriteString(frameRow(cells, widths))
		sb.WriteString("\n")
	}
	sb.WriteString(hrule("└", "┴", "┘", widths))
	sb.WriteString("\n")
	return sb.String()
}

// frameRow pads every cell to its column and wraps the row in borders.
func frameRow(cells []string, widths []int) string {
	parts := make([]string, 0, len(widths))
	for i := range widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}
		parts = append(parts, padTo(cell, widths[i]))
	}
	return "│ " + strings.Join(parts, " │ ") + " │"
}

// hrule draws a horizontal border: left corner, one segment per column
// (content width + cellPad on both sides), right corner.
func hrule(left, mid, right string, widths []int) string {
	segments := make([]string, len(widths))
	for i, w := range widths {
		segments[i] = strings.Repeat("─", w+2*cellPad)
	}
	return left + strings.Join(segments, mid) + right
}

// spriteLines splits krabby output into art rows. Blank edges are dropped,
// and every surviving row is closed with a reset -- krabby's lines typically
// end on an open colour escape with no reset of their own, so anything
// printed after the card would otherwise inherit the sprite's last colour
// (the next "Pokedex >" prompt included).
func spriteLines(art string) []string {
	lines := strings.Split(strings.Trim(art, "\n"), "\n")
	for len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && (lines[len(lines)-1] == "" || visibleWidth(lines[len(lines)-1]) == 0) {
		lines = lines[:len(lines)-1]
	}
	for i, line := range lines {
		lines[i] = closeAnsi(line)
	}
	return lines
}

// closeAnsi appends a reset to a line that leaves colour open.
func closeAnsi(line string) string {
	if strings.Contains(line, "\x1b[") && !strings.HasSuffix(line, "\x1b[0m") {
		return line + "\x1b[0m"
	}
	return line
}

// ansiSeq matches terminal escape sequences so they never count as columns.
var ansiSeq = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// visibleWidth is the printed width of a line: escape sequences are zero
// width and the result is counted in runes, not bytes -- the sprite's
// truecolor escapes and block glyphs are multi-byte.
func visibleWidth(s string) int {
	return utf8.RuneCountInString(ansiSeq.ReplaceAllString(s, ""))
}

// padTo right-pads a line with spaces up to w visible columns.
func padTo(line string, w int) string {
	if n := w - visibleWidth(line); n > 0 {
		return line + strings.Repeat(" ", n)
	}
	return line
}

// displayName capitalises the API's lowercase name for the card header.
func displayName(name string) string {
	if name == "" {
		return name
	}
	first, size := utf8.DecodeRuneInString(name)
	return strings.ToUpper(string(first)) + name[size:]
}

// pokemonCard is the single entry point commandInspect uses: it returns the
// finished card instead of printing it, keeping the CLI out of this file's
// hands entirely.
func pokemonCard(p pokeapi.Pokemon) string {
	return cardBuilder{sprites: krabbySprite{bin: "krabby"}, minLeft: 26}.Render(p)
}
