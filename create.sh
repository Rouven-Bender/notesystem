#!/usr/bin/env bash
#ENV vars needed
#$NOTEBASEPATH
#$EDITOR
#$TERMINAL

type=$1

#function space
function private_note() {
		$TERMINAL -e $EDITOR "+0r $HOME/.local/etc/template.md.tmpl" $NOTEBASEPATH/private/$(cat /proc/sys/kernel/random/uuid).md
}
function normal_note() {
		$TERMINAL -e $EDITOR "+0r $HOME/.local/etc/template.md.tmpl" $NOTEBASEPATH/$(cat /proc/sys/kernel/random/uuid).md
}
function work_note() {
		$TERMINAL -e $EDITOR "+0r $HOME/.local/etc/template.md.tmpl" $NOTEBASEPATH/work/$(cat /proc/sys/kernel/random/uuid).md
}
#

case $type in
		"private")
				private_note
				;;
		"normal" | "")
				normal_note
				;;
		"work")
				work_note
				;;
		*)
				;;
esac
