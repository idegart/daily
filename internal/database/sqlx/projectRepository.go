package sqlx

import (
	"bot/internal/external/airtable"
	"bot/internal/model"
	"github.com/jmoiron/sqlx"
)

type ProjectRepository struct {
	db *sqlx.DB
}

func (r *ProjectRepository) Create(project *model.Project) error {
	return r.db.QueryRow(
		"INSERT INTO projects (name, airtable_id, slack_id, is_infographics) VALUES ($1, $2, $3, $4) RETURNING id",
		project.Name,
		project.AirtableId,
		project.SlackId,
		project.IsInfographics,
	).Scan(&project.Id)
}

func (r *ProjectRepository) Update(project *model.Project) error {
	_, err := r.db.NamedExec(
		"UPDATE projects SET name=:name, airtable_id=:airtable_id, slack_id=:slack_id, is_infographics=:is_infographics, updated_at = now() WHERE id=:id",
		project,
	)

	return err
}

func (r *ProjectRepository) GetAll() ([]model.Project, error) {
	var projects []model.Project

	if err := r.db.Select(&projects, "SELECT * FROM projects"); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepository) GenerateFromAirtable(airProjects []airtable.Project) ([]model.Project, error) {
	projects, err := r.GetAll()

	if err != nil {
		return nil, err
	}

LOOP:
	for i := range airProjects {
		for j := range projects {
			if projects[j].AirtableId == airProjects[i].Fields.ID {
				if err := r.Update(&projects[j]); err != nil {
					return nil, err
				}
				continue LOOP
			}
		}

		project := &model.Project{
			Name:           airProjects[i].Fields.Project,
			AirtableId:     airProjects[i].Fields.ID,
			SlackId:        airProjects[i].Fields.SlackID,
			IsInfographics: false,
		}

		if err := r.Create(project); err != nil {
			return nil, err
		}

		projects = append(projects, *project)
	}

	return projects, err
}

func (r *ProjectRepository) GetUsersForProject(project model.Project) ([]model.User, error) {
	var users []model.User

	if err := r.db.Select(
		&users,
		"SELECT u.* FROM project_users pu INNER JOIN users u on pu.user_id = u.id WHERE pu.project_id=$1",
		project.Id,
	); err != nil {
		return nil, err
	}

	return users, nil
}

func (r ProjectRepository) AttachUsers(project model.Project, users []model.User) error {
	for _, user := range users {
		if err := r.AttachUser(project, user); err != nil {
			return err
		}
	}

	return nil
}

func (r *ProjectRepository) AttachUser(project model.Project, user model.User) error {
	_, err := r.db.Exec(
		"INSERT INTO project_users (project_id, user_id) VALUES ($1, $2) RETURNING id",
		project.Id,
		user.Id,
	)

	return err
}

func (r ProjectRepository) DettachUsers(project model.Project, users []model.User) error {
	var usersIds []int

	for _, user := range users {
		usersIds = append(usersIds, user.Id)
	}

	query, args, err := sqlx.In("DELETE FROM project_users WHERE user_id IN (?) AND project_id=?", usersIds, project.Id)

	if err != nil {
		return err
	}

	_, err = r.db.Exec(r.db.Rebind(query), args...)

	return err
}

func (r *ProjectRepository) SyncUsers(project model.Project, users []model.User) error {
	attachedUsers, err := r.GetUsersForProject(project)

	if err != nil {
		return err
	}

	var usersToAttach []model.User
	var usersToDettach []model.User

LoopAttach:
	for i := range users {
		for j := range attachedUsers {
			if users[i].Id == attachedUsers[j].Id {
				continue LoopAttach
			}
		}

		usersToAttach = append(usersToAttach, users[i])
	}

LoopDettach:
	for i := range attachedUsers {
		for j := range users {
			if attachedUsers[i].Id == users[j].Id {
				continue LoopDettach
			}
		}

		usersToDettach = append(usersToDettach, attachedUsers[i])
	}

	if len(usersToAttach) > 0 {
		if err := r.AttachUsers(project, usersToAttach); err != nil {
			return err
		}
	}

	if len(usersToDettach) > 0 {
		if err := r.DettachUsers(project, usersToDettach); err != nil {
			return err
		}
	}

	return nil
}
