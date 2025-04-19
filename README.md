# About
This is my current notetaking system that integrates with my WM (sway)
# Usage
have the create.sh and search.sh symlinked into your path
e.g have a bin folder in your home directory that is in your path and symlink them into it
then have sway register a keybind for super+n & super+s for the script names
# Setup
```bash
ln -s ...search.sh $HOME/.local/bin/search.sh
ln -s ...create.sh $HOME/.local/bin/create.sh
mkdir $HOME/.local/etc
ln -s ...template.md.tmpl $HOME/.local/etc/template.md.tmpl
```
```lua
vim.keymap.set('n', '<leader>t','/<+++><CR>c5l')
```
# Dependencies
## Env vars
- $EDITOR - your preferd editor
- $NOTEBASEPATH - where the notes are stored at
- $TERM - which terminal to use (designed for alacritty as I use a cli option)
## Software
- wofi - for the menu
# Things for the future
sample.md holds my idea for extra syntax that I would in the future use for an indexer to read.
with {protocol}:{which pc in case of file}//{path} as the syntax for source field.
The requirements file holds some of the requirements I have on this project that I want to have and maybe a rough idea of how

# Why certain design choices
## Why are the filenames uuids not content related?
Because I want to be able to link between notes that are in different contexts have a link to something from my work notes submodule
in my normal notes folder without leaking what the content of that file is. Although maybe context clues of where I link it will leak that
anyway and this way I don't need to think of a name now before writting the note and later when I change it have to change all the links
