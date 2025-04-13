package storage

type Store interface{
	User() UserRepository
	Project() ProjectRepository
	Column() ColumnRepository
	Task_log()Task_log_Repository
	Task() TaskRepository
}