package models

// TableName overrides the table name
func (Category) TableName() string {
	return "category"
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"Name"`
}

func GetCategoryById(id int) (*Category, error) {
	var category *Category
	if err := db.Select(&category).Where("id=?", id).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func GetAllCategories() ([]*Category, error) {
	var categories []*Category
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func AddCategory(name string) error {
	category := Category{Name: name}
	if err := db.Create(&category).Error; err != nil {
		return err
	}
	return nil
}

func EditCategory(id int, name string) error {
	if err := db.Where("id=?", id).Update("name", name).Error; err != nil {
		return err
	}
	return nil
}
