package inbox

type InboxUseCase interface {
	CreateInbox(inbox *Inbox) error
	GetUserInboxes(userId string) ([]Inbox, error)
	MarkInboxAsRead(userId string, id uint) error
	DeleteInbox(userId string, id uint) error
}

type inboxUsecase struct {
	inboxRepo InboxRepository
}

func NewInboxUseCase(repo InboxRepository) InboxUseCase {
	return &inboxUsecase{inboxRepo: repo}
}

func (u *inboxUsecase) CreateInbox(inbox *Inbox) error {
	inbox.IsActive = true
	return u.inboxRepo.Create(inbox)
}

func (u *inboxUsecase) GetUserInboxes(userId string) ([]Inbox, error) {
	return u.inboxRepo.GetAllInboxesByUserId(userId)
}

func (u *inboxUsecase) MarkInboxAsRead(userId string, id uint) error {
	return u.inboxRepo.MarkAsRead(userId, id)
}

func (u *inboxUsecase) DeleteInbox(userId string, id uint) error {
	return u.inboxRepo.Delete(userId, id)
}
