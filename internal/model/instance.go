package model

import "github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"

func (p *Community) Create() error {

	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	str := `
	CREATE TABLE IF NOT EXISTS Community (
		id VARCHAR(128) PRIMARY KEY,
		instance_id VARCHAR(128),
		name TEXT,
		description TEXT,
		icon TEXT,
		cover TEXT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		members BIGINT,
		admins TEXT[],
		moderators TEXT[],
		post TEXT[]
	);`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}
	statement, err := db.Prepare("INSERT INTO Community (id, instance_id, name, description, icon, cover, created_at, updated_at, members, admins, moderators, post) VALUES ( ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?,  ?)")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.Id, p.InstanceId, p.Name, p.Description, p.Icon, p.Cover, p.CreatedAt, p.UpdatedAt, p.Members, p.Admins, p.Moderators, p.Post)
	if err != nil {
		panic(err)
	}
	return nil
}
func (p *Instance) Create() error {

	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	str := `
	CREATE TABLE INSTANCE (
    Id VARCHAR(255) PRIMARY KEY,
    Name VARCHAR(255),
    Description TEXT,
    Type VARCHAR(255),
    Status VARCHAR(255),
    Owner VARCHAR(255),
    CreatedBy VARCHAR(255),
    CommunityIds TEXT[],
    CreatedAt TIMESTAMP,
    UpdatedBy VARCHAR(255),
    UpdatedAt TIMESTAMP,
    DeletedBy VARCHAR(255),
    DeletedAt TIMESTAMP,
    Tags TEXT[],
    Labels TEXT[],
    PublicDomain TEXT[]
)`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}
	statement, err := db.Prepare("INSERT INTO INSTANCE (Id,Name,Description,Type,Status,Owner,CreatedBy,CommunityIds,CreatedAt,UpdatedBy,UpdatedAt,DeletedBy,DeletedAt,Tags,Labels,PublicDomain) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.Id, p.Name, p.Description, p.Type, p.Status, p.Owner, p.CreatedBy, p.CommunityIds, p.CreatedAt, p.UpdatedBy, p.UpdatedAt, p.DeletedBy, p.DeletedAt, p.Tags, p.Labels, p.PublicDomain)
	if err != nil {
		panic(err)
	}
	return nil
}
