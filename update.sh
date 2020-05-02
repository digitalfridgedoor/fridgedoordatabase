./govetgotest.sh
if [ $? -ne 0 ]
  then
    echo "Failing test"
    exit $?
fi

git push
go get -u github.com/digitalfridgedoor/fridgedoordatabase
