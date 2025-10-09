package q1

//包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。

type Student struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"size:100;not null" json:"name"`
	Age   int    `gorm:"not null" json:"age"`
	Grade string `gorm:"size:10;not null" json:"grade"`
}
