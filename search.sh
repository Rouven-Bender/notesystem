#!/usr/bin/env bash
#Needs $NOTEBASEPATH
query=$(echo "" | wofi --dmenu -p search)
if [[ $query != "" ]]
then
	query=$(echo $query | sed -e 's/ /-/g')
	note=$(find $NOTEBASEPATH | grep $query | wofi --dmenu -p which?)
	if [[ $note != "" ]]
	then
		$TERM -e nvim $note
	fi
fi
