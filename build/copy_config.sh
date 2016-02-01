#!/bin/sh
mkdir ./$BUILDDIR/builds/
for dir in $(find ./$BUILDDIR/*/* -type d)
do
	cp -i ./conf.initial.json "$dir"/conf.json;
	cp -i ./conf.example.json "$dir"/conf.example.json;
	# {{.Arch}}_{{.Platform}} pattern
	tar -czf ./$BUILDDIR/builds/$(basename $(dirname $dir))_$(basename $dir).tar.gz $dir
	# tar -cvzf ./$BUILDDIR/builds/$ARCHIVENAME.tar.gz $dir
	# echo $(basename $(dirname $dir))_$(basename $dir)
done