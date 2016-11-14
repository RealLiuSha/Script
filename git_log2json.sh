git log --date=iso --pretty=format:'{%n  "commit": "%H",%n  "author": "%aN",%n  "date": "%ad",%n  "message": "%f"%n},' $@ | perl -pe 'BEGIN{print "["}; END{print "]\n"}' | perl -pe 's/},]/}]/'
