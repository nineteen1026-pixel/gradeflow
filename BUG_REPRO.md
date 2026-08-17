# BUG_REPRO

## 现象

多个评分 worker 同时给同一个学生/作业加分时，`tally.Counter` 的总分会少算：
8 个 goroutine 各调用 `Bump("hw-1", 1)` 500 次，最后 `Value("hw-1")` 拿到的不是 4000，
而是 3500、3921、3994 这类偏小且每次都不一样的数。开 `-race` 跑还会直接报
`WARNING: DATA RACE`，读写两端分别落在 `Total` 和 `Bump` 上。单线程顺序调用一切正常。

## 复现步骤

1. `cp _v6/bug1/zz_verify_test.go internal/tally/zz_verify_test.go`
2. `GOTOOLCHAIN=local go test -race ./internal/tally -run TestVerify -count=20`

## 实际输出

```
==================
WARNING: DATA RACE
Read at 0x00c000102218 by goroutine 78:
  gradeflow/internal/tally.(*Counter).Total()
      /root/work/gradeflow/internal/tally/tally.go:64 +0x1aa
  gradeflow/internal/tally_test.TestVerifyBumpAndReadUnderConcurrency.func2()
      /root/work/gradeflow/internal/tally/zz_verify_test.go:52 +0xaf

Previous write at 0x00c000102218 by goroutine 73:
  gradeflow/internal/tally.(*Counter).Bump()
      /root/work/gradeflow/internal/tally/tally.go:40 +0x14d
  gradeflow/internal/tally_test.TestVerifyBumpAndReadUnderConcurrency.func1()
      /root/work/gradeflow/internal/tally/zz_verify_test.go:45 +0xb1

Goroutine 78 (running) created at:
  gradeflow/internal/tally_test.TestVerifyBumpAndReadUnderConcurrency()
      /root/work/gradeflow/internal/tally/zz_verify_test.go:48 +0x104
  testing.tRunner()
      /usr/local/go1.24.7/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go1.24.7/src/testing/testing.go:1851 +0x44

Goroutine 73 (running) created at:
  gradeflow/internal/tally_test.TestVerifyBumpAndReadUnderConcurrency()
      /root/work/gradeflow/internal/tally/zz_verify_test.go:42 +0x1c4
  testing.tRunner()
      /usr/local/go1.24.7/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go1.24.7/src/testing/testing.go:1851 +0x44
==================
--- FAIL: TestVerifyBumpAndReadUnderConcurrency (0.01s)
    testing.go:1490: race detected during execution of test
--- FAIL: TestVerifyBumpUnderConcurrency (0.01s)
    zz_verify_test.go:29: Value(hw-1) after 4000 concurrent Bump calls = 3921, want 4000
--- FAIL: TestVerifyBumpUnderConcurrency (0.01s)
    zz_verify_test.go:29: Value(hw-1) after 4000 concurrent Bump calls = 3994, want 4000
FAIL
FAIL	gradeflow/internal/tally	0.168s
FAIL
```

## 期望行为

`Bump` 对同一个 key 的并发调用必须是原子的：N 次 `Bump(key, 1)` 之后 `Value(key)` 恰好等于 N，且 `-race` 无告警。

## 影响范围

`internal/tally` 的所有并发计分路径；课程总分、单学生总分会随机偏小，成绩汇总不可信。
