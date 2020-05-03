allExitcodes=()

go vet ./recipe
allExitcodes+=$?
go vet ./userview
allExitcodes+=$?

go test ./recipe
allExitcodes+=$?
go test ./userview
allExitcodes+=$?

for t in ${allExitcodes[@]}; do
  if [[ $t != 0 ]]
    then exit $t
  fi
done