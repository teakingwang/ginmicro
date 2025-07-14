package generator

import (
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ---------- 中文昵称生成 ----------
var zhAdjectives = []string{
	"快乐", "神秘", "冷静", "狂野", "可爱", "调皮", "淡定", "迷人",
}

var zhNouns = []string{
	"虎", "猫", "龙", "鱼", "鸟", "狐", "熊", "狼",
}

func generateChineseNickname() string {
	adjective := zhAdjectives[rand.Intn(len(zhAdjectives))]
	noun := zhNouns[rand.Intn(len(zhNouns))]

	nickname := adjective + noun
	if utf8.RuneCountInString(nickname) > 8 {
		nickname = string([]rune(nickname)[:8])
	}
	return nickname
}

// ---------- 英文昵称生成 ----------
var enPrefixes = []string{
	"Swift", "Foxy", "Nova", "Echo", "Chill", "Luna", "Jazz", "Byte", "Dusk", "Frost",
}

func randomSuffix(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateEnglishNickname() string {
	prefix := enPrefixes[rand.Intn(len(enPrefixes))]
	maxSuffixLen := 8 - len(prefix)
	if maxSuffixLen <= 0 {
		return prefix[:8]
	}
	suffix := randomSuffix(maxSuffixLen)
	return strings.Title(prefix + suffix)
}

// ---------- 通用接口 ----------
func GenerateNickname(lang string) string {
	switch strings.ToLower(lang) {
	case "zh", "cn", "chinese":
		return generateChineseNickname()
	case "en", "english":
		return generateEnglishNickname()
	default:
		// 默认使用中文
		return generateChineseNickname()
	}
}
