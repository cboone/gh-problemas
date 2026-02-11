package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// mockView is a minimal View implementation for testing.
type mockView struct {
	name string
}

func (v *mockView) Init() tea.Cmd            { return nil }
func (v *mockView) Update(tea.Msg) (View, tea.Cmd) { return v, nil }
func (v *mockView) View() string             { return v.name }
func (v *mockView) KeyHints() []string       { return []string{v.name} }

func TestPushPopViewStack(t *testing.T) {
	app := NewApp(nil, "owner/repo", nil)

	v1 := &mockView{name: "dashboard"}
	v2 := &mockView{name: "detail"}

	app.PushView(v1)
	if app.ViewStackLen() != 1 {
		t.Fatalf("expected stack len 1, got %d", app.ViewStackLen())
	}
	if app.CurrentView() != v1 {
		t.Fatal("expected current view to be dashboard")
	}

	app.PushView(v2)
	if app.ViewStackLen() != 2 {
		t.Fatalf("expected stack len 2, got %d", app.ViewStackLen())
	}
	if app.CurrentView() != v2 {
		t.Fatal("expected current view to be detail")
	}

	app.PopView()
	if app.ViewStackLen() != 1 {
		t.Fatalf("expected stack len 1 after pop, got %d", app.ViewStackLen())
	}
	if app.CurrentView() != v1 {
		t.Fatal("expected current view to be dashboard after pop")
	}

	// Pop on last view should not remove it
	app.PopView()
	if app.ViewStackLen() != 1 {
		t.Fatalf("expected stack len 1 (cannot pop last), got %d", app.ViewStackLen())
	}
}

func TestNavigateBackMsg_PopsView(t *testing.T) {
	app := NewApp(nil, "owner/repo", nil)
	app.PushView(&mockView{name: "dashboard"})
	app.PushView(&mockView{name: "detail"})

	model, _ := app.Update(NavigateBackMsg{})
	a := model.(*App)
	if a.ViewStackLen() != 1 {
		t.Fatalf("expected stack len 1 after NavigateBackMsg, got %d", a.ViewStackLen())
	}
}

func TestForceQuit_ProducesQuit(t *testing.T) {
	app := NewApp(nil, "owner/repo", nil)
	app.PushView(&mockView{name: "dashboard"})
	app.PushView(&mockView{name: "detail"})

	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected tea.Quit command from ctrl+c")
	}

	// Execute the command and check it produces a QuitMsg
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected QuitMsg, got %T", msg)
	}
}

func TestQuit_OnLastView(t *testing.T) {
	app := NewApp(nil, "owner/repo", nil)
	app.PushView(&mockView{name: "dashboard"})

	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected tea.Quit command from q on last view")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected QuitMsg, got %T", msg)
	}
}

func TestQuit_OnDetailView_DelegatesToView(t *testing.T) {
	app := NewApp(nil, "owner/repo", nil)
	app.PushView(&mockView{name: "dashboard"})
	app.PushView(&mockView{name: "detail"})

	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	// q on detail view should not quit (delegates to view)
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(tea.QuitMsg); ok {
			t.Fatal("q on detail view should not produce QuitMsg")
		}
	}
}
