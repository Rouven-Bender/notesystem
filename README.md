# About
This is my current notetaking system that integrates with my WM (sway)
# Usage
have the create.sh and search.sh symlinked into your path
e.g have a bin folder in your home directory that is in your path and symlink them into it
then have sway register a keybind for super+n & super+s for the script names
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
