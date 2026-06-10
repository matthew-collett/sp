package ui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"golang.org/x/term"
)

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

func StatusRow(label string, values ...fmt.Stringer) {
	parts := make([]string, 0, len(values))
	for _, v := range values {
		parts = append(parts, v.String())
	}
	fmt.Printf("%s  %s\n", Text("%13s", label).Dimmed(), strings.Join(parts, ""))
}

func ProgressBar(progressMS, durationMS int, filledStyle, emptyStyle func(string) *Style) *Style {
	const width = 20
	filled := 0
	if durationMS > 0 {
		filled = min(int(float64(progressMS)/float64(durationMS)*float64(width)), width)
	}
	bar := filledStyle(strings.Repeat("█", filled)).String() + emptyStyle(strings.Repeat("░", width-filled)).String()
	return Text(bar)
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
