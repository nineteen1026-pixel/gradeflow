# BUG_REPRO

## 现象

调用 `rubric.Apply(criteria)` 之后，调用方自己手里那个 `[]rubric.Criterion` 的顺序被改掉了。
传进去的是 `[style, correctness, tests]`，`Apply` 返回后再看这个切片变成了
`[correctness, tests, style]`。调用方拿同一个切片去渲染评分表 UI 或者存回配置，
顺序就全乱了；同一份 rubric 模板复用两次，第二次的展示顺序和第一次不一样。

## 复现步骤

1. `cp _v6/bug2/zz_verify_test.go internal/rubric/zz_verify_test.go`
2. `GOTOOLCHAIN=local go test ./internal/rubric -run TestVerify -count=20`

## 实际输出

```
--- FAIL: TestVerifyApplyKeepsCallerSliceOrder (0.00s)
    zz_verify_test.go:31: Apply reordered the caller's slice: before [style correctness tests], after [correctness tests style]
--- FAIL: TestVerifyApplyLeavesCallerElementsInPlace (0.00s)
    zz_verify_test.go:46: caller slice after Apply = [correctness style], want [style correctness]
FAIL
FAIL	gradeflow/internal/rubric	0.007s
FAIL
```

## 期望行为

`Apply` 只读入参，返回的 `Report` 自己按权重从大到小排序，调用方传入的切片顺序保持原样。

## 影响范围

`internal/rubric.Apply` 的所有调用方；共享的 rubric 模板切片会被评分过程就地改写，
后续展示、持久化和再次评分都会看到被污染的顺序。
