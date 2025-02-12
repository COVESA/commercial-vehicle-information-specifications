#!/bin/bash

usage() {
	echo "usage: source venv.sh startme|installme" >&2
}

startme() {
	echo "Activating venv..."
	source ~/.venv/bin/activate
	echo "To quit venv use the command 'deactivate'"
}

installme() {
	echo "Installs venv in ~/.venv, activates venv, and then installs vss-tools.\nThis shall only be done if the .venv directory does not exist in the Home directory"
	python3 -m venv ~/.venv
	source ~/.venv/bin/activate
	pip install --pre vss-tools
	echo "To quit venv use the command 'deactivate'"

}

case "$1" in 
	startme)
		startme;;
	installme)
		installme;;
	*)
		usage;;
esac
