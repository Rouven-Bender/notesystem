#!/usr/bin/env bash
#ENV vars needed
#$NOTEBASEPATH
#$EDITOR
#$TERM

#Function space
function create_rough_note () {
	title=$(echo "" | wofi --dmenu -p title)
	if [[ $title != "" ]]
	then
		title=$(echo $title | sed -e 's/ /-/g')
		touch $NOTEBASEPATH/1-rough-notes/$title.md
		$TERM -e $EDITOR $NOTEBASEPATH/1-rough-notes/$title.md &
	fi
}
function create_source_material () {
	title=$(echo "" | wofi --dmenu -p title)
	if [[ $title != "" ]]
	then
		title=$(echo $title | sed -e 's/ /-/g')
		touch $NOTEBASEPATH/2-source-material/$title.md
		$TERM -e $EDITOR $NOTEBASEPATH/2-source-material/$title.md &
	fi
}
function create_full_note () {
	title=$(echo "" | wofi --dmenu -p title)
	if [[ $title != "" ]]
	then
		title=$(echo $title | sed -e 's/ /-/g')
		touch $NOTEBASEPATH/3-full-notes/$title.md
		$TERM -e $EDITOR $NOTEBASEPATH/3-full-notes/$title.md &
	fi
}
function unknown () {
	echo "I don't know what you want from me"
}
#

type=$(echo -e "rough note\nsource material\nfull note" | wofi --dmenu -p type)

case $type in
		"rough note")
			create_rough_note
			;;
		"source material")
			create_source_material
			;;
		"full note")
			create_full_note
			;;
		*)
			unknown
			;;
esac
