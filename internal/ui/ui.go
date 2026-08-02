package ui

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"golang.org/x/term"
)

const (
	minBoxWidth = 50
	maxBoxWidth = 100
)

var spinnerFrames = [...]rune{'|', '/', '-', '\\'}

type Style struct {
	text       string
	attributes []color.Attribute
	out        io.Writer
	fixed      bool
}

func new(format string, out io.Writer, args ...any) *Style {
	text := format
	if len(args) > 0 {
		text = fmt.Sprintf(format, args...)
	}
	return &Style{
		text:       text,
		attributes: make([]color.Attribute, 0),
		out:        out,
	}
}

func Text(format string, args ...any) *Style {
	return new(format, os.Stdout, args...)
}

func Success(format string, args ...any) *Style {
	return Text("%s %s", Check(), Text(format, args...))
}

func Info(format string, args ...any) *Style {
	return Text("%s %s", Arrow(), fmt.Sprintf(format, args...))
}

func Warning(format string, args ...any) *Style {
	return Text("%s %s", Warn(), fmt.Sprintf(format, args...))
}

func Spinner(tick int, format string, args ...any) {
	frame := Text(string(spinnerFrames[tick%len(spinnerFrames)])).Blue()
	fmt.Printf("\r%s %s", Text(format, args...), frame)
}

func ClearLine() {
	fmt.Print("\r\x1b[K")
}

func Print(s string) {
	fmt.Print(s)
}

func ClearScreen() {
	fmt.Print("\x1b[H\x1b[2J")
}

func Error(err error) *Style {
	return new("%s %s", os.Stderr, X(), err)
}

func ErrorVerbose(err error) *Style {
	return new("%s %s", os.Stderr, X(), err)
}

func ShowSplash() {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Println()

	lines := []string{
		"  " + green("/") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + "  " + green("/") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + " ",
		" " + green("/") + yellow("♫") + yellow("♫") + green("_____/") + " " + green("/") + yellow("♫") + yellow("♫") + green("__") + "  " + yellow("♫") + yellow("♫"),
		green("|") + "  " + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + " " + green("|") + " " + yellow("♫") + yellow("♫") + "  " + green("\\") + " " + yellow("♫") + yellow("♫"),
		" " + green("\\____") + "  " + yellow("♫") + yellow("♫") + green("|") + " " + yellow("♫") + yellow("♫") + "  " + green("|") + " " + yellow("♫") + yellow("♫"),
		" " + green("/") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + green("/|") + " " + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + yellow("♫") + green("/"),
		green("|_______/") + " " + green("|") + " " + yellow("♫") + yellow("♫") + green("____/") + " ",
		"          " + green("|") + " " + yellow("♫") + yellow("♫") + "      ",
		"          " + green("|") + " " + yellow("♫") + yellow("♫") + "      ",
		"          " + green("|__/"),
	}
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println()
}

func Check() *Style {
	return Text("✓").Green()
}

func Arrow() *Style {
	return Text("→").Blue()
}

func Warn() *Style {
	return Text("!").Yellow()
}

func X() *Style {
	return Text("✘").Red()
}

func Question() *Style {
	return Text("?").HiBlue()
}

func Dot() *Style {
	return Text("•")
}

func Bool(v bool) *Style {
	if v {
		return Text("true").Green()
	}
	return Text("false").Red()
}

func Play() *Style {
	return Text("▶")
}

func Pause() *Style {
	return Text("⏸")
}

func (s *Style) Green() *Style {
	s.attributes = append(s.attributes, color.FgGreen)
	return s
}

func (s *Style) Yellow() *Style {
	s.attributes = append(s.attributes, color.FgYellow)
	return s
}

func (s *Style) Red() *Style {
	s.attributes = append(s.attributes, color.FgRed)
	return s
}

func (s *Style) Blue() *Style {
	s.attributes = append(s.attributes, color.FgBlue)
	return s
}

func (s *Style) HiBlue() *Style {
	s.attributes = append(s.attributes, color.FgHiBlue)
	return s
}

func (s *Style) Magenta() *Style {
	s.attributes = append(s.attributes, color.FgMagenta)
	return s
}

func (s *Style) White() *Style {
	s.attributes = append(s.attributes, color.FgWhite)
	return s
}

func (s *Style) BgGreen() *Style {
	s.attributes = append(s.attributes, color.BgGreen)
	return s
}

func (s *Style) BgBlack() *Style {
	s.attributes = append(s.attributes, color.BgBlack)
	return s
}

func (s *Style) Bold() *Style {
	s.attributes = append(s.attributes, color.Bold)
	return s
}

func (s *Style) Dimmed() *Style {
	s.attributes = append(s.attributes, color.Faint)
	return s
}

func (s *Style) Underline() *Style {
	s.attributes = append(s.attributes, color.Underline)
	return s
}

func (s *Style) Italic() *Style {
	s.attributes = append(s.attributes, color.Italic)
	return s
}

func (s *Style) Fixed() *Style {
	s.fixed = true
	return s
}

func (s *Style) Indent(spaces int) *Style {
	indent := strings.Repeat(" ", spaces)
	s.text = indent + s.text
	return s
}

func (s *Style) Tab() *Style {
	return s.Indent(2)
}

func (s *Style) TabDeep() *Style {
	return s.Indent(4)
}

func (s *Style) Capitalize() *Style {
	if len(s.text) > 0 {
		s.text = strings.ToUpper(s.text[:1]) + s.text[1:]
	}
	return s
}

func (s *Style) Show() {
	fmt.Fprintln(s.out, s) //nolint:errcheck
}

func (s *Style) String() string {
	if len(s.attributes) > 0 {
		return color.New(s.attributes...).Sprint(s.text)
	}
	return s.text
}

func Truncate(s string, width int) string {
	r := []rune(s)
	if width <= 0 {
		return ""
	}
	if len(r) <= width {
		return s
	}
	if width <= 3 {
		return string(r[:width])
	}
	return string(r[:width-3]) + "..."
}

func Fit(main, extra string, width int) (string, string) {
	main = Truncate(main, width)
	room := width - len([]rune(main)) - 3
	if room <= 0 {
		return main, ""
	}
	return main, Truncate(extra, room)
}

func SideBySide(left, right []string, leftWidth int) []string {
	blank := strings.Repeat(" ", leftWidth)
	rows := max(len(left), len(right))
	out := make([]string, rows)
	for i := range out {
		l, r := blank, ""
		if i < len(left) {
			l = left[i]
		}
		if i < len(right) {
			r = right[i]
		}
		out[i] = l + "  " + r
	}
	return out
}

var reANSI = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func visibleWidth(s string) int {
	return utf8.RuneCountInString(reANSI.ReplaceAllString(s, ""))
}

type Box struct {
	title string
	width int
	rows  []string
}

func NewBox(title string) *Box {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = maxBoxWidth
	}
	return &Box{title: title, width: max(minBoxWidth, min(w, maxBoxWidth))}
}

func (b *Box) Width() int {
	return b.width - 4
}

func (b *Box) Row(content string) {
	b.rows = append(b.rows, content)
}

func (b *Box) Render() string {
	var sb strings.Builder

	label := " " + b.title + " "
	dashes := max(b.width-2-visibleWidth(label), 0)
	sb.WriteString("┌" + label + strings.Repeat("─", dashes) + "┐\n")

	for _, content := range b.rows {
		pad := max(b.Width()-visibleWidth(content), 0)
		sb.WriteString("│ " + content + strings.Repeat(" ", pad) + " │\n")
	}

	sb.WriteString("└" + strings.Repeat("─", b.width-2) + "┘\n")
	return sb.String()
}

func RenderImage(ctx context.Context, url string, cols, rows int) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req) //nolint:gosec
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}
	return renderHalfBlocks(img, cols, rows), nil
}

func renderHalfBlocks(img image.Image, cols, rows int) []string {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	lines := make([]string, rows/2)
	for row := range lines {
		var b strings.Builder
		for col := 0; col < cols; col++ {
			tr, tg, tb := sampleAvg(img, bounds, w, h, cols, rows, col, row*2)
			br, bg, bb := sampleAvg(img, bounds, w, h, cols, rows, col, row*2+1)
			fmt.Fprintf(&b, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀\x1b[0m", tr, tg, tb, br, bg, bb)
		}
		lines[row] = b.String()
	}
	return lines
}

func sampleAvg(img image.Image, bounds image.Rectangle, w, h, cols, rows, col, pixelRow int) (r, g, b uint8) {
	x0 := bounds.Min.X + col*w/cols
	x1 := bounds.Min.X + (col+1)*w/cols
	y0 := bounds.Min.Y + pixelRow*h/rows
	y1 := bounds.Min.Y + (pixelRow+1)*h/rows
	if x1 <= x0 {
		x1 = x0 + 1
	}
	if y1 <= y0 {
		y1 = y0 + 1
	}

	var sr, sg, sb, n uint32
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			pr, pg, pb, _ := img.At(x, y).RGBA()
			sr += pr >> 8
			sg += pg >> 8
			sb += pb >> 8
			n++
		}
	}
	return uint8(sr / n), uint8(sg / n), uint8(sb / n)
}

func ProgressBarOverlay(width int, progress float64, label string) string {
	filled := max(min(int(progress*float64(width)), width), 0)

	l := []rune(label)
	start := max((width-len(l))/2, 0)
	end := min(start+len(l), width)

	var sb strings.Builder
	for i := 0; i < width; i++ {
		r := ' '
		if i >= start && i < end {
			r = l[i-start]
		}

		style := Text("%c", r)
		if i < filled {
			style = style.BgGreen()
		} else {
			style = style.BgBlack()
		}
		if i >= start && i < end {
			style = style.White().Bold()
		}
		sb.WriteString(style.String())
	}
	return sb.String()
}

type row struct {
	cols []*Style
}

type Table struct {
	headers []string
	rows    []row
	title   *Style
}

func NewTable() *Table {
	return &Table{}
}

func (t *Table) Title(s *Style) {
	t.title = s
}

func (t *Table) Header(cols ...string) {
	t.headers = cols
}

func (t *Table) Row(cols ...*Style) {
	t.rows = append(t.rows, row{cols: cols})
}

func (t *Table) Empty() bool {
	return len(t.rows) == 0
}

func (t *Table) Render(pager bool) {
	total := len(t.rows)
	shown := total
	height, ok := termHeight()

	if pager && ok {
		headerLines := 1
		if t.title != nil {
			headerLines = 3
		}
		if total+headerLines > height {
			shown = height - headerLines
		}
	}

	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth = 120
	}

	var buf bytes.Buffer
	t.write(&buf, shown, termWidth)

	if pager && ok && total > shown {
		cmd := exec.Command("less", "-R")
		cmd.Stdin = strings.NewReader(buf.String())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Print(buf.String())
		}
		return
	}
	fmt.Print(buf.String())
}

func (t *Table) write(w io.Writer, shown, termWidth int) {
	n := len(t.headers)

	// compute column widths and fixed flags in one pass
	widths := make([]int, n)
	fixed := make([]bool, n)
	for i, h := range t.headers {
		widths[i] = utf8.RuneCountInString(h)
	}
	for _, r := range t.rows {
		for i, col := range r.cols {
			if i >= n {
				break
			}
			if w := utf8.RuneCountInString(col.text); w > widths[i] {
				widths[i] = w
			}
			if col.fixed {
				fixed[i] = true
			}
		}
	}

	// proportionally shrink non-fixed columns to fit terminal width
	const gap = 2
	if total := totalWidth(widths) + gap*(n-1); total > termWidth {
		overflow := total - termWidth
		var truncatable []int
		truncatableTotal := 0
		for i, w := range widths {
			if !fixed[i] {
				truncatable = append(truncatable, i)
				truncatableTotal += w
			}
		}
		reduced := 0
		for _, i := range truncatable {
			share := int(float64(widths[i]) / float64(truncatableTotal) * float64(overflow))
			widths[i] = max(widths[i]-share, utf8.RuneCountInString(t.headers[i]))
			reduced += share
		}
		if remainder := overflow - reduced; remainder > 0 {
			widest := truncatable[0]
			for _, i := range truncatable[1:] {
				if widths[i] > widths[widest] {
					widest = i
				}
			}
			widths[widest] = max(widths[widest]-remainder, utf8.RuneCountInString(t.headers[widest]))
		}
	}

	if t.title != nil {
		fmt.Fprintf(w, "Showing %d of %d %s\n\n", shown, len(t.rows), t.title.String()) //nolint:errcheck
	}

	for i, h := range t.headers {
		if i > 0 {
			fmt.Fprint(w, "  ") //nolint:errcheck
		}
		fmt.Fprint(w, Text("%-*s", widths[i], h).Dimmed().Underline()) //nolint:errcheck
	}
	fmt.Fprintln(w) //nolint:errcheck

	for _, r := range t.rows {
		for i, col := range r.cols {
			if i > 0 {
				fmt.Fprint(w, "  ") //nolint:errcheck
			}
			text := col.text
			if !fixed[i] && utf8.RuneCountInString(text) > widths[i] {
				text = string([]rune(text)[:widths[i]-3]) + "..."
			}
			fmt.Fprint(w, &Style{text: fmt.Sprintf("%-*s", widths[i], text), attributes: col.attributes, out: col.out}) //nolint:errcheck
		}
		fmt.Fprintln(w) //nolint:errcheck
	}
}

func totalWidth(widths []int) int {
	total := 0
	for _, w := range widths {
		total += w
	}
	return total
}

func termHeight() (height int, ok bool) {
	_, h, err := term.GetSize(int(os.Stdout.Fd()))
	return h, err == nil
}

type SelectOption struct {
	Label string
	Value any
}

func Prompt(format string, args ...any) (string, error) {
	var res string
	if err := askOne(&survey.Input{Message: Text(format, args...).String()}, &res); err != nil {
		return "", err
	}
	return res, nil
}

func PromptSecret(format string, args ...any) (string, error) {
	var res string
	if err := askOne(&survey.Password{Message: Text(format, args...).String()}, &res); err != nil {
		return "", err
	}
	return res, nil
}

func PromptConfirm(format string, args ...any) (bool, error) {
	var res bool
	if err := askOne(&survey.Confirm{Message: Text(format, args...).String()}, &res); err != nil {
		return false, err
	}
	return res, nil
}

func PromptConfirmYes(format string, args ...any) (bool, error) {
	var res bool
	if err := askOne(&survey.Confirm{Message: Text(format, args...).String(), Default: true}, &res); err != nil {
		return false, err
	}
	return res, nil
}

func PromptSelect(message string, opts []SelectOption) (*SelectOption, error) {
	labels := make([]string, len(opts))
	for i, option := range opts {
		labels[i] = option.Label
	}

	var res int
	prompt := &survey.Select{
		Message: message,
		Options: labels,
	}

	err := askOne(prompt, &res)
	if err != nil {
		return nil, err
	}

	return &opts[res], nil
}

func askOne(prompt survey.Prompt, res any) error {
	if err := survey.AskOne(prompt, res, survey.WithIcons(func(is *survey.IconSet) {
		is.Question.Text = Question().String()
		is.Error.Text = X().String()
	})); err != nil {
		if err == terminal.InterruptErr {
			return nil // silent exit
		}
		return err
	}
	return nil
}
