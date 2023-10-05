# gex! Hex Editor

## Introduction

gex! is a hex editor designed with vi-like keybindings.

**Disclaimer:** gex! is currently in early development, and its features and
keybindings may change as it evolves. Use it with caution, and make sure to back
up your data before editing any files.

## Usage

- Load a file: `gex <filename>`
- See this help file: `gex --help`
- See all avaliable help files: `gex --list-help`
- See a specific help file: `gex --help <help file>`

## Keybindings

### Movement Keys

- `h` / `j` / `k` / `l`: Move the cursor left / down / up / right.
- `0` / `$`: Move the cursor to the beginning / end of the current line.
- `gg` / `G`: Move the cursor to the start / end of the file.
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
- `goto <offset>`: Jump to `<offset>` (hex).
- `set <option> <value>`: Set an option for the current session. See below for
  the list of options.

### Options

- `cols <n>`: Number of colums displayed. Default is 16.
- `inspector.enabled <true|false>`: Enable / disable the inspector
- `inspector.byteOrder <byteOrder>`: Set the byte order of the inspector. Value
  could be `big`, `be`, or `b` for BE, or `little`, `le`, or `l` for LE.
  Defaults to LE.

## Caveats

Note that at the current stage, gex! might behave differently than other text /
hex editors. For instance:

- gex! does not load the whole file into memory. Instead, only sections that are
  visible will be read. This means if the underlying file gets modified outside
  of gex! while editing, and the file offsets change, then the changes you've
  made may no longer correspond to the same point in the original file.
- When saving a file, the undo and redo stack gets cleared. This may be fixed in
  future versions of gex!.
- When saving a file, gex! first writes to the file name suffixed with `~`.
  Then, the temporary file is swapped with the original file. The original file
  will now have a `~` suffix in the filename. This means there will be two
  copies of the file on disk, the modified one and the original one.
