package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

const TableNameTask = "task"

type TaskStatus string

const (
	TaskStatusNormal    TaskStatus = "Normal"
	TaskStatusPending   TaskStatus = "Pending"
	TaskStatusRunning   TaskStatus = "Running"
	TaskStatusSucceeded TaskStatus = "Succeeded"
	TaskStatusFailed    TaskStatus = "Failed"
	TaskStatusUnknown   TaskStatus = "Unknown"
)

type Task struct {
	ID        int64      `gorm:"column:id" json:"id"`                 // 任务 ID
	Name      string     `gorm:"column:name" json:"name"`             // 任务名称
	Namespace string     `gorm:"column:namespace" json:"namespace"`   // 任务 k8s namespace 名称
	Info      TaskInfo   `gorm:"column:info" json:"info"`             // 任务 k8s 相关信息
	Status    TaskStatus `gorm:"column:status" json:"status"`         // 任务状态
	UserID    int64      `gorm:"column:user_id" json:"user_id"`       // 用户 ID
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"` // 修改时间
}

func (*Task) TableName() string {
	return TableNameTask
}

type TaskInfo struct {
	Image   string   `json:"image"`
	Command []string `json:"command"`
	Args    []string `json:"args"`
}

// Scan implements the [Scanner] interface.
func (ti *TaskInfo) Scan(value any) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), ti)
}

// Value implements the [driver.Valuer] interface.
func (ti TaskInfo) Value() (driver.Value, error) {
	bytes, err := json.Marshal(ti)
	return string(bytes), err
}
