package careers

import (
	"github.com/pgvector/pgvector-go"
)

type Career struct {
	ID                     int64                `db:"id" json:"id"`
	Title                  string               `db:"title" json:"title"`
	Description            string               `db:"description" json:"description"`
	PersonalityDescription string               `db:"personality_description" json:"personalityDescription"`
	Education              string               `db:"education" json:"education"`
	AverageSalary          string               `db:"average_salary" json:"averageSalary"`
	LowerSalary            string               `db:"lower_salary" json:"lowerSalary"`
	HighestSalary          string               `db:"highest_salary" json:"highestSalary"`
	Embedding              pgvector.Vector      `db:"embedding" json:"embedding"`
	TasksString            []string             `db:"-" json:"tasks"`
	Tasks                  []Task               `db:"-" json:"careerTasks"`
	Knowledge              []KnowledgeCategory  `db:"-" json:"knowledge"`
	Abilities              []AbilityCategory    `db:"-" json:"abilities"`
	SkillCategories        []SkillCategory      `db:"-" json:"skills"`
	TecnologyCategories    []TechnologyCategory `db:"-" json:"technology"`
	Personality            Personality          `db:"-" json:"personality"`
}

type Task struct {
	ID              int64  `db:"id" json:"id"`
	CareerID        int64  `db:"career_id" json:"careerId"`
	TaskDescription string `db:"task_description" json:"taskDescription"`
}

type KnowledgeCategory struct {
	ID       int64    `db:"id" json:"id"`
	CareerID int64    `db:"career_id" json:"careerId"`
	Name     string   `db:"category_name" json:"name"`
	Areas    []string `db:"-" json:"areas"` // Now a slice of strings
}

type KnowledgeArea struct {
	ID         int64  `db:"id" json:"id"`
	CategoryID int64  `db:"category_id" json:"categoryId"`
	AreaName   string `db:"area_name" json:"areaName"`
}

type AbilityCategory struct {
	ID       int64    `db:"id" json:"id"`
	CareerID int64    `db:"career_id" json:"careerId"`
	Name     string   `db:"category_name" json:"name"`
	Areas    []string `db:"-" json:"areas"` // Now a slice of strings
}

type AbilityArea struct {
	ID         int64  `db:"id" json:"id"`
	CategoryID int64  `db:"category_id" json:"categoryId"`
	AreaName   string `db:"area_name" json:"areaName"`
}

type SkillCategory struct {
	ID       int64    `db:"id" json:"id"`
	CareerID int64    `db:"career_id" json:"careerId"`
	Name     string   `db:"category_name" json:"name"`
	Areas    []string `db:"-" json:"areas"` // Now a slice of strings
}

type SkillArea struct {
	ID         int64  `db:"id" json:"id"`
	CategoryID int64  `db:"category_id" json:"categoryId"`
	AreaName   string `db:"area_name" json:"areaName"`
}

type TechnologyCategory struct {
	ID       int64    `db:"id" json:"id"`
	CareerID int64    `db:"career_id" json:"careerId"`
	Name     string   `db:"category_name" json:"name"`
	Areas    []string `db:"-" json:"areas"` // Now a slice of strings
}

type TechnologyArea struct {
	ID         int64  `db:"id" json:"id"`
	CategoryID int64  `db:"category_id" json:"categoryId"`
	AreaName   string `db:"area_name" json:"areaName"`
}

type Personality struct {
	Description string   `db:"-" json:"description"`
	Attributes  []string `db:"-" json:"attributes"`
}

type PersonalityAttribute struct {
	ID            int64  `db:"id" json:"id"`
	CareerID      int64  `db:"career_id" json:"careerId"`
	AttributeName string `db:"attribute_name" json:"attributeName"`
}
