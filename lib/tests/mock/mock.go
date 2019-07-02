package mock

import "github.com/SENERGY-Platform/device-manager/lib/config"

func New(cin config.Config) (publisher *Publisher, cout config.Config, close func()) {
	cout = cin

	publisher = NewPublisher()

	repo := NewDeviceRepo(publisher)
	cout.DeviceRepoUrl = repo.Url()

	semantic := NewSemanticRepo(publisher)
	cout.SemanticRepoUrl = semantic.Url()

	perm := NewPermSearch()
	cout.PermissionsUrl = perm.Url()

	close = func() {
		repo.Stop()
		semantic.Stop()
	}

	return
}
