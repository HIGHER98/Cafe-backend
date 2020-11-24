package models

const (
	PENDING_TRANSACTION = 1
	PENDING             = 2
	CONFIRMED           = 3
	COLLECTION          = 4
)

type Status struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by User to `profiles`
func (Status) TableName() string {
	return "status"
}

//INSERT INTO TABLE status (description) VALUES (?);
func AddStatus(desc string) error {
	newStatus := Status{Description: desc}
	err := db.Create(&newStatus).Error
	return err
}

//Get Status
func GetStatus(id int) (*Status, error) {
	var status Status
	if err := db.Where("id").Where("id=? AND is_del=0", id).First(&status).Error; err != nil {
		return nil, err
	}
	return &status, nil
}

//UPDATE TABLE status SET is_del=0 WHERE id=?;
func DeleteStatus(id int) error {
	status, err := GetStatus(id)
	if err != nil {
		return err
	}
	db.Model(&status).Update("is_del", 1)
	return err
}

func GetStatuses() ([]*Status, error) {
	var s []*Status
	if err := db.Select("id", "description").Where("is_del=0").Find(&s).Error; err != nil {
		return nil, err
	}
	return s, nil
}
