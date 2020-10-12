// ORIGINAL: Part of javatest/StringUtilTest.java

package stringutil_test

import (
	"testing"

	"github.com/markusmobius/go-domdistiller/internal/stringutil"
	"github.com/stretchr/testify/assert"
)

func Test_StringUtil_FullWordCounter(t *testing.T) {
	counter := stringutil.FullWordCounter{}
	// One Chinese sentence, or a series of Japanese glyphs should not be treated
	// as a single word.
	assert.True(t, counter.Count("一個中文句子不應該當成一個字") > 1) // zh-Hant
	assert.True(t, counter.Count("中国和马来西亚使用简体字") > 1)   // zh-Hans
	assert.True(t, counter.Count("ファイナルファンタジー") > 1)    // Katakana
	assert.True(t, counter.Count("いってらっしゃい") > 1)       // Hiragana
	assert.True(t, counter.Count("仏仮駅辺") > 1)           // Kanji

	// However, treating each Chinese/Japanese glyph as a word is also wrong.
	assert.True(t, counter.Count("一個中文句子不應該當成一個字") < 14)
	assert.True(t, counter.Count("中国和马来西亚使用简体字") < 12)
	assert.True(t, counter.Count("ファイナルファンタジー") < 11)
	assert.True(t, counter.Count("いってらっしゃい") < 8)
	assert.True(t, counter.Count("仏仮駅辺") < 4)

	// Even if they are separated by spaces.
	assert.True(t, counter.Count("一 個 中 文 句 子 不 應 該 當 成 一 個 字") < 14)
	assert.True(t, counter.Count("中 国 和 马 来 西 亚 使 用 简 体 字") < 12)
	assert.True(t, counter.Count("フ ァ イ ナ ル フ ァ ン タ ジ ー") < 11)
	assert.True(t, counter.Count("い っ て ら っ し ゃ い") < 8)
	assert.True(t, counter.Count("仏 仮 駅 辺") < 4)
	assert.Equal(t, 1, counter.Count("字"))
	assert.Equal(t, 1, counter.Count("が"))

	// Mixing ASCII words and Chinese/Japanese glyphs
	assert.Equal(t, 2, counter.Count("word字"))
	assert.Equal(t, 2, counter.Count("word 字"))
}

func Test_StringUtil_LetterWordCounter(t *testing.T) {
	counters := []stringutil.WordCounter{
		stringutil.LetterWordCounter{},
		stringutil.FullWordCounter{},
	}

	for _, counter := range counters {
		// Hangul uses space as word delimiter like English.
		assert.Equal(t, 1, counter.Count("어"))
		assert.Equal(t, 2, counter.Count("한국어 단어"))
		assert.Equal(t, 5, counter.Count("한 국 어 단 어"))
		assert.Equal(t, 8, counter.Count("예비군 훈련장 총기 난사범 최모씨의 군복에서 발견된 유서."))
	}
}

func Test_StringUtil_FastWordCounter(t *testing.T) {
	counters := []stringutil.WordCounter{
		stringutil.FastWordCounter{},
		stringutil.LetterWordCounter{},
		stringutil.FullWordCounter{},
	}

	for _, counter := range counters {
		assert.Equal(t, 0, counter.Count(""))
		assert.Equal(t, 0, counter.Count("  -@# ';]"))
		assert.Equal(t, 1, counter.Count("word"))
		assert.Equal(t, 1, counter.Count("b'fore"))
		assert.Equal(t, 1, counter.Count(" _word.under_score_ "))
		assert.Equal(t, 2, counter.Count(" \ttwo\nwords"))
		assert.Equal(t, 2, counter.Count(" \ttwo @^@^&(@#$([][;;\nwords"))
		// Norwegian
		assert.Equal(t, 5, counter.Count("dør når på svært dårlig"))
		assert.Equal(t, 5, counter.Count("svært få dør av blåbærsyltetøy"))
		// Greek
		assert.Equal(t, 11, counter.Count("Παρέμβαση των ΗΠΑ για τα τεχνητά νησιά που κατασκευάζει η Κίνα"))
		// Arabic
		assert.Equal(t, 6, counter.Count("زلزال بقوة 8.5 درجات يضرب اليابان"))
		// Tibetan
		assert.Equal(t, 1, counter.Count("༧གོང་ས་མཆོག་གི་ནང་གི་ངོ་སྤྲོད་ཀྱི་གསུང་ཆོས་ལེགས་གྲུབ།"))
		// Thai
		assert.Equal(t, 3, counter.Count("โซลาร์ อิมพัลส์ทู เหินฟ้าข้ามมหาสมุทร"))
	}
}

func Test_StringUtil_SelectWordCounter(t *testing.T) {
	counter := stringutil.SelectWordCounter("abc")
	if _, ok := counter.(stringutil.FastWordCounter); !ok {
		t.Errorf("abc should use FastWordCounter")
	}

	counter = stringutil.SelectWordCounter("어")
	if _, ok := counter.(stringutil.LetterWordCounter); !ok {
		t.Errorf("hangul should use LetterWordCounter")
	}

	counter = stringutil.SelectWordCounter("字")
	if _, ok := counter.(stringutil.FullWordCounter); !ok {
		t.Errorf("zh should use FullWordCounter")
	}
}
