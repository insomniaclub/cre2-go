package cre2

import (
	"regexp"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	patternUnsupport = `(?=abc)123??????xxx------`
	patternSimple    = `(?:学历|education|文化程度|教育经历|white)\s?=\s?(小学|初中|高中|中专|专科|大专|本科|研究生|博士)`

	textUnmatch    = `这个字符串不会被上面的pattern匹配命中`
	textSimple     = `测试匹配这个字符串学历=小学&education=大专&教育经历=专科&white=本科结束`
	simpleMatch    = []string{"学历=小学", "education=大专", "教育经历=专科", "white=本科"}
	simpleMatchIdx = [][]int{{27, 40}, {41, 57}, {58, 77}, {78, 90}}
	simpleSubmatch = [][]string{{"学历=小学", "小学"}, {"education=大专", "大专"}, {"教育经历=专科", "专科"}, {"white=本科", "本科"}}
)

func TestCompile(t *testing.T) {
	convey.Convey("TestCompile", t, func() {
		var (
			re  *Regexp
			err error
		)

		re, err = Compile(patternSimple)
		convey.So(err, convey.ShouldBeNil)
		convey.So(re, convey.ShouldNotBeNil)

		re, err = Compile(patternUnsupport)
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(re, convey.ShouldBeNil)
	})
}

func TestMatch(t *testing.T) {
	convey.Convey("TestMatch", t, func() {
		re := MustCompile(patternSimple)
		convey.So(re.MatchString(textSimple), convey.ShouldBeTrue)
		convey.So(re.MatchString(textUnmatch), convey.ShouldBeFalse)
	})
}

func TestFindString(t *testing.T) {
	convey.Convey("TestFindAllString", t, func() {
		re := MustCompile(patternSimple)
		convey.So(re.FindString(textUnmatch), convey.ShouldEqual, "")
		convey.So(re.FindString(textSimple), convey.ShouldEqual, simpleMatch[0])
	})
}

func TestFindAllString(t *testing.T) {
	convey.Convey("TestFindAllString", t, func() {
		re := MustCompile(patternSimple)
		convey.So(re.FindAllString(textUnmatch, -1), convey.ShouldEqual, []string(nil))
		convey.So(re.FindAllString(textSimple, 0), convey.ShouldEqual, []string(nil))
		convey.So(re.FindAllString(textSimple, -1), convey.ShouldResemble, simpleMatch)
		convey.So(re.FindAllString(textSimple, 10), convey.ShouldResemble, simpleMatch)
		convey.So(re.FindAllString(textSimple, 2), convey.ShouldResemble, simpleMatch[:2])
	})
}

func TestFindAllStringSubmatch(t *testing.T) {
	convey.Convey("TestFindAllStringSubmatch", t, func() {
		re := MustCompile(patternSimple)
		convey.So(re.FindAllStringSubmatch(textUnmatch, -1), convey.ShouldEqual, [][]string(nil))
		convey.So(re.FindAllStringSubmatch(textSimple, 0), convey.ShouldEqual, [][]string(nil))
		convey.So(re.FindAllStringSubmatch(textSimple, -1), convey.ShouldResemble, simpleSubmatch)
		convey.So(re.FindAllStringSubmatch(textSimple, 10), convey.ShouldResemble, simpleSubmatch)
		convey.So(re.FindAllStringSubmatch(textSimple, 2), convey.ShouldResemble, simpleSubmatch[:2])
	})
}

func TestFindAllStringIndex(t *testing.T) {
	convey.Convey("TestFindAllStringSubmatch", t, func() {
		re := MustCompile(patternSimple)
		convey.So(re.FindAllStringIndex(textUnmatch, -1), convey.ShouldEqual, [][]int(nil))
		convey.So(re.FindAllStringIndex(textSimple, 0), convey.ShouldEqual, [][]int(nil))
		convey.So(re.FindAllStringIndex(textSimple, -1), convey.ShouldResemble, simpleMatchIdx)
		convey.So(re.FindAllStringIndex(textSimple, 10), convey.ShouldResemble, simpleMatchIdx)
		convey.So(re.FindAllStringIndex(textSimple, 2), convey.ShouldResemble, simpleMatchIdx[:2])
	})
}

func BenchmarkCre2Compile_Simple(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MustCompile(patternSimple)
	}
}

func BenchmarkOriginCompile_Simple(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = regexp.MustCompile(patternSimple)
	}
}

func BenchmarkCre2Match_Simple(b *testing.B) {
	re := MustCompile(patternSimple)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.MatchString(textSimple)
	}
}

func BenchmarkOriginMatch_Simple(b *testing.B) {
	re := regexp.MustCompile(patternSimple)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.MatchString(textSimple)
	}
}

func BenchmarkCre2FindString_Simple(b *testing.B) {
	re := MustCompile(patternSimple)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.FindString(textSimple)
	}
}

func BenchmarkOriginFindString_Simple(b *testing.B) {
	re := regexp.MustCompile(patternSimple)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.FindString(textSimple)
	}
}

func BenchmarkCre2FindAllString_Simple(b *testing.B) {
	re := MustCompile(patternSimple)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.FindAllString(textSimple, -1)
	}
}

func BenchmarkOriginFindAllString_Simple(b *testing.B) {
	re := regexp.MustCompile(patternSimple)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.FindAllString(textSimple, -1)
	}
}

func BenchmarkCre2FindAllStringSubmatch_Simple(b *testing.B) {
	re := MustCompile(patternSimple)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.FindAllStringSubmatch(textSimple, -1)
	}
}

func BenchmarkOriginFindAllStringSubmatch_Simple(b *testing.B) {
	re := regexp.MustCompile(patternSimple)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.FindAllStringSubmatch(textSimple, -1)
	}
}

func BenchmarkCre2FindAllStringIndex_Simple(b *testing.B) {
	re := MustCompile(patternSimple)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.FindAllStringIndex(textSimple, -1)
	}
}

func BenchmarkOriginFindAllStringIndex_Simple(b *testing.B) {
	re := regexp.MustCompile(patternSimple)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.FindAllStringIndex(textSimple, -1)
	}
}
