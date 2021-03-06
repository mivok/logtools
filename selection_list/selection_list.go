// List widget with selection bar based on the termui List widget
package selection_list

// TODO
// - Implement wrapping (Overflow: wrap) - this is especially useful when
// viewing individual log lines which are long

import (
	ui "github.com/gizak/termui"
)

type SelectionList struct {
	ui.List
	SelectedItem        int // Selected item
	ScrollY             int // First visible item
	ScrollX             int
	EnableSelection     bool // Is item selection active?
	SelectedItemBgColor ui.Attribute
	SelectedItemFgColor ui.Attribute
}

func NewSelectionList() *SelectionList {
	l := &SelectionList{List: *ui.NewList()}
	// We don't have theme items for these, so just pick some values here by
	// default
	l.SelectedItemFgColor = ui.ColorBlack
	l.SelectedItemBgColor = ui.ColorCyan
	l.EnableSelection = true
	return l
}

func (l *SelectionList) Scroll(x, y int, absolute bool) {
	if absolute {
		l.ScrollX = x
		l.ScrollY = y
	} else {
		l.ScrollX += x
		l.ScrollY += y
	}
	if l.ScrollY >= (len(l.Items) - l.InnerHeight()) {
		l.ScrollY = (len(l.Items) - l.InnerHeight() - 1)
	}
	if l.ScrollY < 0 {
		l.ScrollY = 0
	}
	if l.ScrollX < 0 {
		l.ScrollX = 0
	}
}

func (l *SelectionList) SelectItem(count int, absolute bool) {
	if !l.EnableSelection {
		// Selection mode is turned off, so just scroll instead
		l.Scroll(0, count, absolute)
	}
	// Select a new item, with bounds checking
	if absolute {
		l.SelectedItem = count
	} else {
		l.SelectedItem += count
	}
	if l.SelectedItem >= len(l.Items) {
		l.SelectedItem = len(l.Items) - 1
	}
	if l.SelectedItem < 0 {
		l.SelectedItem = 0
	}
	// Scroll so that the selected item is always in view
	if l.ScrollY > l.SelectedItem {
		l.ScrollY = l.SelectedItem
	}
	if l.ScrollY < l.SelectedItem+1-l.InnerHeight() {
		l.ScrollY = l.SelectedItem + 1 - l.InnerHeight()
	}
}

func (l *SelectionList) Buffer() ui.Buffer {
	buf := l.Block.Buffer()

	trimItems := l.Items[l.ScrollY:]
	if len(trimItems) > l.Block.InnerHeight() {
		trimItems = trimItems[:l.Block.InnerHeight()]
	}
	for i, v := range trimItems {
		fg := l.ItemFgColor
		bg := l.ItemBgColor
		if l.ScrollX > 0 {
			// Trim the beginning of the line if we are scrolled to the right
			if len(v) > l.ScrollX {
				v = v[l.ScrollX:]
			} else {
				v = ""
			}
		}
		cs := ui.DTrimTxCls(ui.DefaultTxBuilder.Build(v, fg, bg),
			l.Block.InnerWidth())
		j := 0
		for _, vv := range cs {
			w := vv.Width()
			buf.Set(l.Block.InnerX()+j, l.Block.InnerY()+i, vv)
			j += w
		}
		if i+l.ScrollY == l.SelectedItem && l.EnableSelection {
			// Draw selection bar
			for x := 0; x < l.Block.InnerWidth(); x++ {
				cell := buf.At(l.Block.InnerX()+x, l.Block.InnerY()+i)
				cell.Fg = l.SelectedItemFgColor
				cell.Bg = l.SelectedItemBgColor
				buf.Set(l.Block.InnerX()+x, l.Block.InnerY()+i, cell)
			}
		}
	}
	return buf
}
