package helper

type Change struct {
	Table      string
	Column     string
	Id         string
	OpenDataId string
	OldValue   string
	NewValue   string
	//onderstaande drie zijn standaard leeg, enkel bij gegenereerde data ingevuld
	Station_id string
	User_id    string
	Defect     string
}
