package service

import repository "github.com/ecofriends/authentication-backend/repository"

type DatabaseProvider struct {
	Repo *repository.PostGreSQL
}

func (database *DatabaseProvider) New(repo *repository.PostGreSQL) {
	database.Repo = repo
}
