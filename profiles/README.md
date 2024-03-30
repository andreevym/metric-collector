## Profile Diff

Diff between source version 'base' and version after changes for split param uri without chi 'result'

In result parse arg with chi work faster in this block of code, because chi parsed params earlier and put them through
context.

**CPU**

**command**

```shell
go tool pprof -top -diff_base=profiles/base_cpu.pprof profiles/result_cpu.pprof
```

**result**

```
File: ___go_build_github_com_andreevym_metric_collector_profiles
Type: cpu
Time: Mar 31, 2024 at 12:09am (MSK)
Duration: 401.27ms, Total samples = 0 
Showing nodes accounting for 0, 0% of 0 total
      flat  flat%   sum%        cum   cum%
```

**MEM**

**command**

```shell
go tool pprof -top -diff_base=profiles/base_mem.pprof profiles/result_mem.pprof
```

**result**

```
File: ___go_build_github_com_andreevym_metric_collector_profiles
Type: inuse_space
Time: Mar 31, 2024 at 12:09am (MSK)
Showing nodes accounting for 0, 0% of 1.72MB total
      flat  flat%   sum%        cum   cum%

```