// Package fortune 提供算命内容数据与随机抽取逻辑。
package fortune

import (
	"math/rand"
	"time"
)

func init() { rand.Seed(time.Now().UnixNano()) }

// Category 表示运势类别。
type Category string

const (
	CategoryLove   Category = "love"
	CategoryCareer Category = "career"
	CategoryHealth Category = "health"
	CategoryWealth Category = "wealth"
	CategoryGeneral Category = "general"
)

// CategoryLabel 返回类别中文名。
func (c Category) Label() string {
	switch c {
	case CategoryLove:
		return "姻缘"
	case CategoryCareer:
		return "事业"
	case CategoryHealth:
		return "健康"
	case CategoryWealth:
		return "财运"
	default:
		return "综合"
	}
}

// Item 表示一条运势内容。
type Item struct {
	Category Category `json:"category"`
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Hint     string   `json:"hint,omitempty"`
}

var loveFortunes = []Item{
	{CategoryLove, "红鸾星动", "近期桃花运佳，宜主动把握机会，多参与社交。", "宜静心，勿急躁。"},
	{CategoryLove, "缘定三生", "命中良缘将至，留意身边默默关心你的人。", "真诚待人，自有回报。"},
	{CategoryLove, "琴瑟和鸣", "与伴侣关系和谐，可计划共同目标或短途出行。", "沟通为上。"},
}

var careerFortunes = []Item{
	{CategoryCareer, "平步青云", "工作上有贵人相助，适合推进重要项目或争取晋升。", "保持谦逊，多请教。"},
	{CategoryCareer, "厚积薄发", "前期积累将在本月显现，宜主动汇报成果、争取机会。", "稳扎稳打。"},
	{CategoryCareer, "左右逢源", "团队协作顺畅，适合牵头跨部门事项。", "注意分寸与边界。"},
}

var healthFortunes = []Item{
	{CategoryHealth, "神清气爽", "整体状态良好，宜保持作息与适度运动。", "少熬夜。"},
	{CategoryHealth, "养精蓄锐", "适合调养身心，可减少应酬、早睡早起。", "情绪也会影响体质。"},
	{CategoryHealth, "顺其自然", "小恙易愈，不必过度担忧，遵医嘱即可。", "多晒太阳、多走动。"},
}

var wealthFortunes = []Item{
	{CategoryWealth, "小有进账", "正财稳定，偏财有机会但不宜贪多。", "理性理财。"},
	{CategoryWealth, "稳中求进", "适合储蓄与稳健投资，避免冲动消费。", "量入为出。"},
	{CategoryWealth, "贵人带财", "合作或介绍可带来额外收益，宜把握靠谱机会。", "勿轻信高回报承诺。"},
}

var generalFortunes = []Item{
	{CategoryGeneral, "吉星高照", "诸事顺遂，宜做重要决定或开启新计划。", "顺势而为。"},
	{CategoryGeneral, "否极泰来", "之前的不顺将逐渐化解，保持信心。", "多行善事。"},
	{CategoryGeneral, "心想事成", "心念专注之事，有望在近期看到进展。", "行动比空想重要。"},
}

// All 返回按类别分组的全部运势条目（只读）。
func All() map[Category][]Item {
	return map[Category][]Item{
		CategoryLove:    loveFortunes,
		CategoryCareer:  careerFortunes,
		CategoryHealth:  healthFortunes,
		CategoryWealth:  wealthFortunes,
		CategoryGeneral: generalFortunes,
	}
}

// DrawRandom 从所有类别中随机抽一条运势。
func DrawRandom() Item {
	all := All()
	flat := make([]Item, 0, 15)
	for _, items := range all {
		flat = append(flat, items...)
	}
	return flat[rand.Intn(len(flat))]
}

// DrawByCategory 从指定类别中随机抽一条；若类别不存在则从全部中抽。
func DrawByCategory(c Category) Item {
	all := All()
	if items, ok := all[c]; ok && len(items) > 0 {
		return items[rand.Intn(len(items))]
	}
	return DrawRandom()
}
