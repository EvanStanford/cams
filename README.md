# cams
Generates cam profiles (.stl) for the mechanical laser show


Check out these for more info:

https://www.youtube.com/watch?v=_dtBUiaAqRE

https://hackaday.io/project/25447-mechanical-laser-show

https://www.thingiverse.com/thing:2383299




# usage

```
#Install go, git
mkdir go/src/github.com/EvanStanford/
cd go/src/github.com/EvanStanford/
#Add environment variable GOPATH to go directory
git clone https://github.com/EvanStanford/cams.git
cd cams/profiler/
go get
go test
#PASS
go install
cd ../main/
go install
cd ../../../../../bin/
main.exe ../src/github.com/EvanStanford/cams/profiler/testfiles/star_path.csv 0.045
ls out/star_path/
```


