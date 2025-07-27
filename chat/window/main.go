package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const Prod_Name = "TermoTalks"

// LineCharPos for tracking position within the content
type LineCharPos struct {
	Line int
	Char int
}

type ChatWindowModel struct {
	Prod_Name     string
	Room_ID       string
	Sender_Name   string
	Receiver_Name string
	Sender_Time   time.Time
	Receiver_Time time.Time

	// New fields for chat content and selection
	content        []string    // The actual chat messages/lines
	selectionMode  bool        // Is the user currently in selection mode?
	selectionStart LineCharPos // Start of the selection
	selectionEnd   LineCharPos // End of the selection (updated as cursor moves)
	cursor         LineCharPos // Current position of the cursor

	selectedText string // Stores the text that was selected and "Ctrl+R"ed

	// Window dimensions for potential scrolling/layout
	width, height int
}

// Styles for rendering selection and cursor
var (
	selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("205")).Foreground(lipgloss.Color("0"))
	cursorStyle   = lipgloss.NewStyle().Reverse(true)

	// New style for the centered header
	headerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center). // Center align text within its own width
			Padding(0, 2)           // Optional: Add some horizontal padding
)

func NewChatModel(RID string, SN string, ST time.Time, RN string, RT time.Time) ChatWindowModel {
	return ChatWindowModel{
		Prod_Name:     Prod_Name, // Using the constant
		Room_ID:       RID,
		Sender_Name:   SN,
		Receiver_Name: RN,
		Sender_Time:   ST,
		Receiver_Time: RT,
		content: []string{ // Example chat content
			"Welcome to TermoTalks!",
			"This is a sample chat message from Alice.",
			"Bob replies: Hello Alice, how are you?",
			"Alice: I'm doing great, thanks for asking!",
			"Bob: What's up for today?",
			"Alice: Just working on some Go projects.",
			"This is a longer line to test wrapping and selection across lines.",
			"Last message in this conversation.",
		},
		cursor: LineCharPos{Line: 0, Char: 0}, // Initial cursor position
	}
}

func (c ChatWindowModel) Init() tea.Cmd { return nil }

func (c ChatWindowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		return c, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c": // Global quit
			return c, tea.Quit

		case "v": // Toggle selection mode (Vim-like visual mode)
			if !c.selectionMode {
				c.selectionMode = true
				c.selectionStart = c.cursor
				c.selectionEnd = c.cursor
			} else {
				c.selectionMode = false
				c.selectedText = "" // Clear selected text on exit
			}
			return c, nil

		case "up", "k":
			if c.selectionMode {
				c.moveCursorUp()
				c.selectionEnd = c.cursor // Update end of selection
			} else {
				// Normal cursor movement (not selection, just moving view or cursor)
				if c.cursor.Line > 0 {
					c.cursor.Line--
				}
				// Adjust char if line becomes shorter
				if c.cursor.Char > len(c.content[c.cursor.Line]) {
					c.cursor.Char = len(c.content[c.cursor.Line])
				}
			}
			return c, nil

		case "down", "j":
			if c.selectionMode {
				c.moveCursorDown()
				c.selectionEnd = c.cursor // Update end of selection
			} else {
				// Normal cursor movement
				if c.cursor.Line < len(c.content)-1 {
					c.cursor.Line++
				}
				// Adjust char if line becomes shorter
				if c.cursor.Char > len(c.content[c.cursor.Line]) {
					c.cursor.Char = len(c.content[c.cursor.Line])
				}
			}
			return c, nil

		case "left", "h":
			if c.selectionMode {
				c.moveCursorLeft()
				c.selectionEnd = c.cursor
			} else {
				if c.cursor.Char > 0 {
					c.cursor.Char--
				} else if c.cursor.Line > 0 {
					c.cursor.Line--
					c.cursor.Char = len(c.content[c.cursor.Line]) // Move to end of previous line
				}
			}
			return c, nil

		case "right", "l":
			if c.selectionMode {
				c.moveCursorRight()
				c.selectionEnd = c.cursor
			} else {
				if c.cursor.Char < len(c.content[c.cursor.Line]) {
					c.cursor.Char++
				} else if c.cursor.Line < len(c.content)-1 {
					c.cursor.Line++
					c.cursor.Char = 0 // Move to start of next line
				}
			}
			return c, nil

		case "home": // Move cursor to start of line
			if c.cursor.Line < len(c.content) { // Ensure line exists
				c.cursor.Char = 0
				if c.selectionMode {
					c.selectionEnd = c.cursor
				}
			}
			return c, nil

		case "end": // Move cursor to end of line
			if c.cursor.Line < len(c.content) { // Ensure line exists
				c.cursor.Char = len(c.content[c.cursor.Line])
				if c.selectionMode {
					c.selectionEnd = c.cursor
				}
			}
			return c, nil

		case "ctrl+d": // Vim-like half page down
			linesToMove := c.height / 2
			for i := 0; i < linesToMove; i++ {
				c.moveCursorDown()
			}
			if c.selectionMode {
				c.selectionEnd = c.cursor
			}
			return c, nil

		case "ctrl+u": // Vim-like half page up
			linesToMove := c.height / 2
			for i := 0; i < linesToMove; i++ {
				c.moveCursorUp()
			}
			if c.selectionMode {
				c.selectionEnd = c.cursor
			}
			return c, nil

		case "ctrl+r": // This is where the 'select+ctrl+r' action happens
			if c.selectionMode {
				// 1. Get the selected text
				c.selectedText = c.getSelectedText()

				// 2. Perform the action with the selected text
				c.selectionMode = false // Exit selection mode after acting
				return c, tea.Batch(
					tea.Printf("Ctrl+R action: Processed selected text: \"%s\"", c.selectedText),
				)
			}
			return c, nil

		case "esc": // Exit selection mode
			if c.selectionMode {
				c.selectionMode = false
				c.selectedText = "" // Clear selected text
			}
			return c, nil
		}
	}
	return c, nil
}

// Helper functions for cursor movement
func (c *ChatWindowModel) moveCursorUp() {
	if c.cursor.Line > 0 {
		c.cursor.Line--
		if c.cursor.Char > len(c.content[c.cursor.Line]) {
			c.cursor.Char = len(c.content[c.cursor.Line])
		}
	} else {
		c.cursor.Char = 0 // At the top, move to start of line 0
	}
}

func (c *ChatWindowModel) moveCursorDown() {
	if c.cursor.Line < len(c.content)-1 {
		c.cursor.Line++
		if c.cursor.Char > len(c.content[c.cursor.Line]) {
			c.cursor.Char = len(c.content[c.cursor.Line])
		}
	} else {
		// At the bottom, move to end of last line
		if len(c.content) > 0 { // Avoid panic if content is empty
			c.cursor.Char = len(c.content[c.cursor.Line])
		}
	}
}

func (c *ChatWindowModel) moveCursorLeft() {
	if c.cursor.Char > 0 {
		c.cursor.Char--
	} else if c.cursor.Line > 0 {
		c.cursor.Line--
		c.cursor.Char = len(c.content[c.cursor.Line]) // Move to end of previous line
	}
	if c.cursor.Line < 0 { // Clamp to 0
		c.cursor.Line = 0
		c.cursor.Char = 0
	}
}

func (c *ChatWindowModel) moveCursorRight() {
	if c.cursor.Char < len(c.content[c.cursor.Line]) {
		c.cursor.Char++
	} else if c.cursor.Line < len(c.content)-1 {
		c.cursor.Line++
		c.cursor.Char = 0 // Move to start of next line
	}
	if c.cursor.Line >= len(c.content) { // Clamp to last line
		if len(c.content) > 0 {
			c.cursor.Line = len(c.content) - 1
			c.cursor.Char = len(c.content[c.cursor.Line])
		} else { // No content
			c.cursor.Line = 0
			c.cursor.Char = 0
		}
	}
}

// getSelectedText extracts the text between selectionStart and selectionEnd
func (c ChatWindowModel) getSelectedText() string {
	if !c.selectionMode || (c.selectionStart.Line == c.selectionEnd.Line && c.selectionStart.Char == c.selectionEnd.Char) {
		return ""
	}

	// Normalize selection: ensure 'start' is always before 'end' in document order
	start := c.selectionStart
	end := c.selectionEnd
	if start.Line > end.Line || (start.Line == end.Line && start.Char > end.Char) {
		start, end = end, start
	}

	var builder strings.Builder
	for i := start.Line; i <= end.Line; i++ {
		line := c.content[i]
		lineStartChar := 0
		lineEndChar := len(line)

		if i == start.Line {
			lineStartChar = start.Char
		}
		if i == end.Line {
			lineEndChar = end.Char
		}

		// Ensure indices are within bounds of the current line
		if lineStartChar < 0 {
			lineStartChar = 0
		}
		if lineEndChar > len(line) {
			lineEndChar = len(line)
		}

		if lineStartChar < len(line) && lineEndChar >= 0 && lineStartChar <= lineEndChar {
			builder.WriteString(line[lineStartChar:lineEndChar])
		}

		if i < end.Line {
			builder.WriteString("\n") // Add newline between lines
		}
	}
	return builder.String()
}

func (c ChatWindowModel) View() string {
	var s strings.Builder

	// 1. Construct the header string
	headerInfo := fmt.Sprintf("%s - Room ID: %s\nSender: %s (%s) | Receiver: %s (%s)",
		c.Prod_Name,
		c.Room_ID,
		c.Sender_Name, c.Sender_Time.Format("15:04"),
		c.Receiver_Name, c.Receiver_Time.Format("15:04"))

	// 2. Use lipgloss.Place to center the header horizontally
	// We need to ensure c.width is properly set from WindowSizeMsg
	if c.width == 0 { // Fallback if width isn't set yet (e.g., very first render)
		return "Loading..."
	}
	centeredHeader := headerStyle.Width(c.width).Render(headerInfo)
	s.WriteString(centeredHeader)
	s.WriteString("\n\n") // Add some space after the header

	s.WriteString(strings.Repeat("-", c.width) + "\n\n") // Separator line

	// Normalize selection points for rendering, as in getSelectedText
	renderStart := c.selectionStart
	renderEnd := c.selectionEnd
	if renderStart.Line > renderEnd.Line || (renderStart.Line == renderEnd.Line && renderStart.Char > renderEnd.Char) {
		renderStart, renderEnd = renderEnd, renderStart
	}

	// Display chat content
	for lineNum, line := range c.content {
		lineContent := line // Start with the full line content

		if c.selectionMode {
			// Check if this line is within the selected range
			if lineNum >= renderStart.Line && lineNum <= renderEnd.Line {
				effectiveStartChar := 0
				effectiveEndChar := len(line)

				if lineNum == renderStart.Line {
					effectiveStartChar = renderStart.Char
				}
				if lineNum == renderEnd.Line {
					effectiveEndChar = renderEnd.Char
				}

				// Ensure indices are valid for the current line
				if effectiveStartChar < 0 {
					effectiveStartChar = 0
				}
				if effectiveEndChar > len(line) {
					effectiveEndChar = len(line)
				}
				if effectiveStartChar > effectiveEndChar { // Happens for empty selection on a line
					effectiveStartChar = effectiveEndChar
				}

				// Apply selection style
				var parts []string
				if effectiveStartChar > 0 {
					parts = append(parts, lineContent[:effectiveStartChar])
				}
				parts = append(parts, selectedStyle.Render(lineContent[effectiveStartChar:effectiveEndChar]))
				if effectiveEndChar < len(lineContent) {
					parts = append(parts, lineContent[effectiveEndChar:])
				}
				lineContent = strings.Join(parts, "")
			}
		}

		// Add cursor (only in non-selection mode or at the selection end point)
		if c.cursor.Line == lineNum && (!c.selectionMode || c.cursor == c.selectionEnd) {
			if c.cursor.Char < len(lineContent) {
				s.WriteString(lineContent[:c.cursor.Char])
				s.WriteString(cursorStyle.Render(string(lineContent[c.cursor.Char])))
				s.WriteString(lineContent[c.cursor.Char+1:])
			} else {
				s.WriteString(lineContent)
				s.WriteString(cursorStyle.Render(" ")) // Cursor at end of line
			}
		} else {
			s.WriteString(lineContent)
		}
		s.WriteString("\n")
	}

	// Status line/debug info
	s.WriteString(strings.Repeat("-", c.width) + "\n")
	s.WriteString(fmt.Sprintf("Mode: %s | Cursor: L%d C%d",
		func() string {
			if c.selectionMode {
				return "SELECT"
			}
			return "NORMAL"
		}(), c.cursor.Line, c.cursor.Char))

	if c.selectionMode {
		s.WriteString(fmt.Sprintf(" | Sel: (L%d C%d) to (L%d C%d)",
			renderStart.Line, renderStart.Char, renderEnd.Line, renderEnd.Char))
	}
	if c.selectedText != "" {
		s.WriteString(fmt.Sprintf("\nLAST PROCESSED: %q", c.selectedText))
	}

	s.WriteString("\n\n(Press 'v' to toggle select mode, 'h/j/k/l' or arrows to move, 'home/end' to jump.)")
	s.WriteString("\n(Press 'ctrl+r' to process selection, 'q' or 'ctrl+c' to quit.)")

	return s.String()
}

func main() {
	// Create an instance of your chat model
	// Note: You might want to adjust these times to be more relevant to "now"
	// if you're running this as a test.
	model := NewChatModel(
		"Room123",
		"Alice",
		time.Now().Add(-5*time.Minute), // 5 mins ago
		"Bob",
		time.Now(), // now
	)

	// Create a new Bubble Tea program and run it
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}
