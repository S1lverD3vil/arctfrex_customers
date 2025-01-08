package inbox

import "gorm.io/gorm"

type inboxRepository struct {
	db *gorm.DB
}

func NewInboxRepository(db *gorm.DB) InboxRepository {
	return &inboxRepository{db: db}
}

func (ir *inboxRepository) Create(inbox *Inbox) error {
	return ir.db.Create(inbox).Error
}

func (ir *inboxRepository) GetAll(receiver string) ([]Inbox, error) {
	var inboxes []Inbox
	err := ir.db.Where("receiver = ?", receiver).Find(&inboxes).Error
	return inboxes, err
}

func (ir *inboxRepository) GetAllInboxesByUserId(userId string) ([]Inbox, error) {
	var inboxes []Inbox
	err := ir.db.Where("is_active = ? AND user_id = ?", true, userId).Find(&inboxes).Error
	return inboxes, err
}

func (ir *inboxRepository) MarkAsRead(userId string, id uint) error {
	return ir.db.Model(&Inbox{}).Where("is_active = ? AND id = ? AND user_id = ?", true, id, userId).Update("is_read", true).Error
}

func (ir *inboxRepository) Delete(userId string, id uint) error {
	return ir.db.Model(&Inbox{}).Where("id = ? AND user_id = ?", id, userId).Update("is_active", false).Error
	// return ir.db.Where("id = ?", id).Delete(&Inbox{}).Error
}
