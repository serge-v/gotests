ver=`git log -n 1 --format=%h`
[ ! -d version ] && mkdir version
echo "package version" > version/version~.txt
echo "const HEAD = \""${ver}"\"" >> version/version~.txt
echo "const STATUS = \`" >> version/version~.txt
git st -s >> version/version~.txt
echo "\`" >> version/version~.txt
echo "const DIFF = \`" >> version/version~.txt
git diff | sed "s/\`/ __BK__ /g" >> version/version~.txt
echo "\`" >> version/version~.txt

diff version/version~.txt version/version.go 2> /dev/null
if [ $? -ne 0 ] ; then
	cp version/version~.txt version/version.go
	rm version/version~.txt
	echo "version.go generated. New version:" ${ver}
fi
