package careers

type Career struct {
	ID                     int64     `db:"id" json:"id"`
	Title                  string    `db:"title" json:"title"`
	Description            string    `db:"description" json:"description"`
	PersonalityDescription string    `db:"personality_description" json:"personalityDescription"`
	Education              string    `db:"education" json:"education"`
	AverageSalary          string    `db:"average_salary" json:"averageSalary"`
	LowerSalary            string    `db:"lower_salary" json:"lowerSalary"`
	HighestSalary          string    `db:"highest_salary" json:"highestSalary"`
	Embedding              []float64 `db:"embedding" json:"embedding"`
	// Tasks                  []Task                 `db:"-" json:"tasks"`
	// Knowledge              []KnowledgeCategory    `db:"-" json:"knowledge"`
	// Abilities              []AbilityCategory      `db:"-" json:"abilities"`
	// SkillCategories        []SkillCategory        `db:"-" json:"skillCategories"`
	// TecnologyCategories    []TechnologyCategory   `db:"-" json:"technologyCategories"`
	// PersonalityAttributes  []PersonalityAttribute `db:"-" json:"personalityAttributes"`
}

type Task struct {
	ID              int64  `db:"id" json:"id"`
	CareerID        int64  `db:"career_id" json:"careerId"`
	TaskDescription string `db:"task_description" json:"taskDescription"`
}

type KnowledgeCategory struct {
	ID             int64           `db:"id" json:"id"`
	CareerID       int64           `db:"career_id" json:"careerId"`
	CategoryName   string          `db:"category_name" json:"categoryName"`
	KnowledgeAreas []KnowledgeArea `db:"-" json:"knowledgeAreas"`
}

type KnowledgeArea struct {
	ID       int64  `db:"id" json:"id"`
	CareerID int64  `db:"career_id" json:"careerId"`
	AreaName string `db:"area_name" json:"areaName"`
}

type AbilityCategory struct {
	ID           int64         `db:"id" json:"id"`
	CareerID     int64         `db:"career_id" json:"careerId"`
	CategoryName string        `db:"category_name" json:"categoryName"`
	AbilityAreas []AbilityArea `db:"-" json:"abilityAreas"`
}

type AbilityArea struct {
	ID       int64  `db:"id" json:"id"`
	CareerID int64  `db:"career_id" json:"careerId"`
	AreaName string `db:"area_name" json:"areaName"`
}

type SkillCategory struct {
	ID           int64       `db:"id" json:"id"`
	CareerID     int64       `db:"career_id" json:"careerId"`
	CategoryName string      `db:"category_name" json:"categoryName"`
	SkillAreas   []SkillArea `db:"-" json:"skillAreas"`
}

type SkillArea struct {
	ID       int64  `db:"id" json:"id"`
	CareerID int64  `db:"career_id" json:"careerId"`
	AreaName string `db:"area_name" json:"areaName"`
}

type TechnologyCategory struct {
	ID              int64            `db:"id" json:"id"`
	CareerID        int64            `db:"career_id" json:"careerId"`
	CategoryName    string           `db:"category_name" json:"categoryName"`
	TechnologyAreas []TechnologyArea `db:"-" json:"technologyAreas"`
}

type TechnologyArea struct {
	ID       int64  `db:"id" json:"id"`
	CareerID int64  `db:"career_id" json:"careerId"`
	AreaName string `db:"area_name" json:"areaName"`
}

type PersonalityAttribute struct {
	ID            int64  `db:"id" json:"id"`
	CareerID      int64  `db:"career_id" json:"careerId"`
	AttributeName string `db:"attribute_name" json:"attributeName"`
}
