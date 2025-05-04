#!/usr/bin/env bash
#Needs $NOTEBASEPATH
#Needs ripgrep
query=$(echo "" | wofi --dmenu -p search)
if [[ $query != "" ]]
then
	note=$(rg -l $query $NOTEBASEPATH | wofi --dmenu -p which?)
	if [[ $note != "" ]]
	then
		$TERMINAL -e nvim $note
	fi
fi
