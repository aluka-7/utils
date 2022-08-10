package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type testStringTime struct {
	Time StringTime
}

func Test_StringTime(t *testing.T) {
	now := time.Unix(0, 1582775051014441000)
	st := testStringTime{Time: StringTime(now)}
	v, err := json.Marshal(st)
	if err != nil {
		t.Error("时间转换成字符串格式错误", err)
	}
	expected := "{\"Time\":\"2020-02-27 11:44:11\"}"
	actual := string(v)
	if actual != expected {
		t.Error("生成的结果不匹配\n", "预期:", expected, "|", "实际:", actual)
	}
	var tst testStringTime
	err = json.Unmarshal(v, &tst)
	if err != nil {
		t.Error("字符串转换成时间格式错误", err)
	}
	expected = now.Format("2006-01-02 15:04:05")
	actual = time.Time(tst.Time).Format("2006-01-02 15:04:05")
	if actual != expected {
		t.Error("生成的结果不匹配\n", "预期:", expected, "|", "实际:", actual)
	}
}

type testNumberTime struct {
	Time NumberTime
}

func Test_NumberTime(t *testing.T) {
	now := time.Unix(0, 1582775051014441000)
	nt := testNumberTime{Time: NumberTime(now)}
	v, err := json.Marshal(nt)
	if err != nil {
		t.Error("时间格式转换成UTC错误", err)
	}
	expected := "{\"Time\":1582775051014441000}"
	actual := string(v)
	if actual != expected {
		t.Error("生成的结果不匹配\n", "预期:", expected, "|", "实际:", actual)
	}

	var tnt testNumberTime
	err = json.Unmarshal(v, &tnt)
	if err != nil {
		t.Error("UTC转换成时间格式错误", err)
	}
	expected = now.String()
	actual = time.Time(tnt.Time).String()
	if actual != expected {
		t.Error("生成的结果不匹配\n", "预期:", expected, "|", "实际:", actual)
	}
}

type testDuration struct {
	Duration Duration
}

func Test_Duration(t *testing.T) {
	std := "{\"Duration\":\"1s\"}"
	var td testDuration
	err := json.Unmarshal([]byte(std), &td)
	if err != nil {
		t.Error("转换成Duration错误", err)
	}
	expected := fmt.Sprintf("%d", td.Duration)
	actual := "1000000000"
	if actual != expected {
		t.Error("生成的结果不匹配\n", "预期:", expected, "|", "实际:", actual)
	}
}

func TestShrink(t *testing.T) {
	var d Duration
	err := d.UnmarshalText([]byte("1s"))
	if err != nil {
		t.Fatalf("TestShrink:  d.UnmarshalText failed!err:=%v", err)
	}
	c := context.Background()
	to, ctx, cancel := d.Shrink(c)
	defer cancel()
	if time.Duration(to) != time.Second {
		t.Fatalf("new timeout must be equal 1 second")
	}
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > time.Second || time.Until(deadline) < time.Millisecond*500 {
		t.Fatalf("ctx deadline must be less than 1s and greater than 500ms")
	}
}

func TestShrinkWithTimeout(t *testing.T) {
	var d Duration
	err := d.UnmarshalText([]byte("1s"))
	if err != nil {
		t.Fatalf("TestShrink:  d.UnmarshalText failed!err:=%v", err)
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	to, ctx, cancel := d.Shrink(c)
	defer cancel()
	if time.Duration(to) != time.Second {
		t.Fatalf("new timeout must be equal 1 second")
	}
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > time.Second || time.Until(deadline) < time.Millisecond*500 {
		t.Fatalf("ctx deadline must be less than 1s and greater than 500ms")
	}
}

func TestShrinkWithDeadline(t *testing.T) {
	var d Duration
	err := d.UnmarshalText([]byte("1s"))
	if err != nil {
		t.Fatalf("TestShrink:  d.UnmarshalText failed!err:=%v", err)
	}
	c, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	to, ctx, cancel := d.Shrink(c)
	defer cancel()
	if time.Duration(to) >= time.Millisecond*500 {
		t.Fatalf("new timeout must be less than 500 ms")
	}
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > time.Millisecond*500 || time.Until(deadline) < time.Millisecond*200 {
		t.Fatalf("ctx deadline must be less than 500ms and greater than 200ms")
	}
}
