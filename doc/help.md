# gex! Hex Editor

## Introduction

gex! is a hex editor designed with vi-like keybindings.

**Disclaimer:** gex! is currently in early development, and its features and
keybindings may change as it evolves. Use it with caution, and make sure to back
up your data before editing any files.

## Keybindings

### Movement Keys

- `h` / `j` / `k` / `l`: Move the cursor left / down / up / right.
- `0` / `$`: Move the cursor to the beginning / end of the current line.
- `g` / `G`: Move the cursor to the start / end of the file.
- `ctrl+d` / `ctrl+u`: Scroll down / up one screen.

### Action Keys

- `x`: Delete the selected byte(s).
- `y`: Yank (copy) the selected byte(s).
- `s`: Substitute (replace) the selected byte(s).
- `P` / `p`: Paste the last yanked byte(s) before / after the cursor.

### Normal Mode Commands

- `tab`: Switch focus between the hex / ascii column.
- `i` / `a`: Enter insert mode before / after the cursor position.
- `v`: Enter visual mode to select a range of bytes.
- `R`: Enter replace mode to overwrite bytes.
- `:`: Enter command mode to execute commands.
- `u` / `ctrl+r`: Undo / redo the last edit.

### Commands

- `w`: Write changes to the file.
- `q`: Quit gex! if there are no unsaved changes.
- `q!`: Quit gex! forcefully, discarding unsaved changes.
