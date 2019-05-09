# Steps

- Clone this repository to a location outside `GOPATH`.
- Run the `repro` bash script. It should clone `x/tools`, run `repro.go` which consumes `10` packages in `x/tools`, then display the memory profile from `go tool pprof`.

# My results (local)

One commit prior to potential regression:

```
: cat 45dd101d8784.top | head -20
File: repro
Type: alloc_space
Time: May 9, 2019 at 12:03am (UTC)
Showing nodes accounting for 170.73MB, 96.71% of 176.54MB total
Dropped 46 nodes (cum <= 0.88MB)
      flat  flat%   sum%        cum   cum%
         0     0%     0%   159.41MB 90.30%  golang.org/x/tools/go/packages.(*loader).loadPackage
         0     0%     0%   159.41MB 90.30%  golang.org/x/tools/go/packages.(*loader).loadRecursive.func1
         0     0%     0%   159.41MB 90.30%  sync.(*Once).Do
         0     0%     0%   158.91MB 90.01%  golang.org/x/tools/go/packages.(*loader).loadRecursive
         0     0%     0%   153.74MB 87.09%  golang.org/x/tools/go/packages.(*loader).loadRecursive.func1.1
         0     0%     0%   152.67MB 86.48%  golang.org/x/tools/go/gcexportdata.Read
         0     0%     0%   152.67MB 86.48%  golang.org/x/tools/go/packages.(*loader).loadFromExportData
         0     0%     0%   146.98MB 83.26%  bytes.(*Buffer).grow
  146.98MB 83.26% 83.26%   146.98MB 83.26%  bytes.makeSlice
         0     0% 83.26%   146.48MB 82.97%  bytes.(*Buffer).ReadFrom
         0     0% 83.26%   143.37MB 81.21%  io/ioutil.readAll
         0     0% 83.26%   142.86MB 80.92%  io/ioutil.ReadAll
    0.50MB  0.28% 83.54%    11.81MB  6.69%  golang.org/x/tools/go/internal/gcimporter.(*iimporter).doDecl
         0     0% 83.54%    11.81MB  6.69%  golang.org/x/tools/go/internal/gcimporter.(*importReader).obj
```

Potential regression commit based on bisect:

```
: cat dbeab5af4b8d.top | head -20
File: repro
Type: alloc_space
Time: May 9, 2019 at 12:03am (UTC)
Showing nodes accounting for 2647.39MB, 95.13% of 2782.91MB total
Dropped 158 nodes (cum <= 13.91MB)
      flat  flat%   sum%        cum   cum%
         0     0%     0%  1890.71MB 67.94%  go/types.(*Checker).checkFiles
         0     0%     0%  1863.15MB 66.95%  go/types.(*Checker).Files
         0     0%     0%  1842.55MB 66.21%  golang.org/x/tools/go/packages.(*loader).loadPackage
         0     0%     0%  1827.82MB 65.68%  golang.org/x/tools/go/packages.(*loader).loadRecursive.func1
         0     0%     0%  1795.45MB 64.52%  sync.(*Once).Do
         0     0%     0%  1762.59MB 63.34%  golang.org/x/tools/go/packages.(*loader).loadRecursive
         0     0%     0%  1730.81MB 62.19%  golang.org/x/tools/go/packages.(*loader).loadRecursive.func1.1
         0     0%     0%  1410.51MB 50.68%  go/types.(*Checker).rawExpr
   19.50MB   0.7%   0.7%  1302.49MB 46.80%  go/types.(*Checker).stmt
         0     0%   0.7%  1299.35MB 46.69%  go/types.(*Checker).multiExpr
         0     0%   0.7%  1298.87MB 46.67%  go/types.(*Checker).stmtList
    2.16MB 0.077%  0.78%  1275.78MB 45.84%  go/types.(*Checker).exprInternal
         0     0%  0.78%  1227.46MB 44.11%  go/types.(*Checker).funcBody
         0     0%  0.78%  1188.49MB 42.71%  go/types.(*Checker).funcDecl.func1
```

# Current newest

```
: cat cf84161cff3f.top | head -20
File: repro
Type: alloc_space
Time: May 9, 2019 at 12:04am (UTC)
Showing nodes accounting for 2675.72MB, 95.78% of 2793.53MB total
Dropped 165 nodes (cum <= 13.97MB)
      flat  flat%   sum%        cum   cum%
         0     0%     0%  1853.25MB 66.34%  go/types.(*Checker).checkFiles
         0     0%     0%  1832.76MB 65.61%  go/types.(*Checker).Files
         0     0%     0%  1804.23MB 64.59%  golang.org/x/tools/go/packages.(*loader).loadPackage
         0     0%     0%  1792.70MB 64.17%  golang.org/x/tools/go/packages.(*loader).loadRecursive.func1
         0     0%     0%  1762.17MB 63.08%  sync.(*Once).Do
         0     0%     0%  1734.12MB 62.08%  golang.org/x/tools/go/packages.(*loader).loadRecursive
         0     0%     0%  1704.61MB 61.02%  golang.org/x/tools/go/packages.(*loader).loadRecursive.func1.1
         0     0%     0%  1390.45MB 49.77%  go/types.(*Checker).rawExpr
         0     0%     0%  1273.83MB 45.60%  go/types.(*Checker).multiExpr
   14.50MB  0.52%  0.52%  1271.30MB 45.51%  go/types.(*Checker).stmt
         0     0%  0.52%  1268.72MB 45.42%  go/types.(*Checker).stmtList
    2.10MB 0.075%  0.59%  1248.86MB 44.71%  go/types.(*Checker).exprInternal
         0     0%  0.59%  1200.02MB 42.96%  go/types.(*Checker).funcBody
         0     0%  0.59%  1161.72MB 41.59%  go/types.(*Checker).funcDecl.func1

```

# My results (on Travis CI)

https://travis-ci.org/codeactual/repro/builds/530055112#L355

# Travis config

> For one-off build:

```
script:
  - cd github/golang.org/go/x_tools_go_packages_dbeab5a
  - ./repro
  - head -20 45dd101d8784.top
  - head -20 dbeab5af4b8d.top
  - head -20 cf84161cff3f.top
```
