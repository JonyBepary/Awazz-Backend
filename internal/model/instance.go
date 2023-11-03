package model

import "github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"

func (p *Community) Create() error {

	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	str := `CREATE TABLE IF NOT EXISTS COMMUNITY (
	Id VARCHAR(255) PRIMARY KEY NOT NULL,
	InstanceId VARCHAR(255),
	Name VARCHAR(255),
	Description VARCHAR(255),
	Icon VARCHAR(255),
	Cover VARCHAR(255),
	CreatedAt BIGINT,
	UpdatedAt BIGINT,
	Members BIGINT,
	Admins VARCHAR(255),
	Moderators VARCHAR(255),
	Post VARCHAR(255)
)`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}
	statement, err := db.Prepare("INSERT INTO COMMUNITY (Id,InstanceId,Name,Description,Icon,Cover,CreatedAt,UpdatedAt,Members,Admins,Moderators,Post) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.Id, p.InstanceId, p.Name, p.Description, p.Icon, p.Cover, p.CreatedAt, p.UpdatedAt, p.Members, p.Admins, p.Moderators, p.Post)
	if err != nil {
		panic(err)
	}
	return nil
}

func (p *Community) GetCommunity(cid string) error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// spew.Dump(rows)
	//! fmt.Println("message id is: ", pid)
	row, err := db.Query("SELECT * FROM COMMUNITY WHERE Id=?", cid)
	if err != nil {
		panic(err)
	}

	row.Next()
	err = row.Scan(&p.Id, &p.InstanceId, &p.Name, &p.Description, &p.Icon, &p.Cover, &p.CreatedAt, &p.UpdatedAt, &p.Members, &p.Admins, &p.Moderators, &p.Post)
	if err != nil {
		panic(err)
	}

	err = row.Err()
	if err != nil {
		panic(err)
	}
	row.Close()

	//! spew.Dump(p.Id)
	return nil
}
func (p *Community) DeleteCommunity(cid string) error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sql_cmd := `DELETE FROM COMMUNITY WHERE Id=?`
	statement, err := db.Prepare(sql_cmd)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(cid)
	if err != nil {
		panic(err)
	}
	return nil
}

func (p *Community) UpdatedCommuninty() error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sql_cmd := `UPDATE COMMUNITY SET Id=?,InstanceId=?,Name=?,Description=?,Icon=?,Cover=?,CreatedAt=?,UpdatedAt=?,Members=?,Admins=?,Moderators=?,Post=?`
	statement, err := db.Prepare(sql_cmd)
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
	CREATE TABLE IF NOT EXISTS INSTANCE (
    Id VARCHAR(255) PRIMARY KEY,
    Name VARCHAR(255),
    Description TEXT,
    Type VARCHAR(255),
    Status VARCHAR(255),
    CreatedAt TIMESTAMP,
    UpdatedBy VARCHAR(255),
    UpdatedAt TIMESTAMP,
    DeletedBy VARCHAR(255),
    DeletedAt TIMESTAMP
)`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}
	statement, err := db.Prepare("INSERT INTO INSTANCE (Id,Name,Description,Type,Status,CreatedAt,UpdatedBy,UpdatedAt,DeletedBy,DeletedAt) VALUES (?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.Id, p.Name, p.Description, p.Type, p.Status, p.CreatedAt, p.UpdatedBy, p.UpdatedAt, p.DeletedBy, p.DeletedAt)
	if err != nil {
		panic(err)
	}
	return nil
}

func (p *Instance) GetInstance(cid string) error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// spew.Dump(rows)
	//! fmt.Println("message id is: ", pid)
	row, err := db.Query("SELECT * FROM INSTANCE WHERE Id=?", cid)
	if err != nil {
		panic(err)
	}

	row.Next()
	err = row.Scan(p.Id, p.Name, p.Description, p.Type, p.Status, p.Owner, p.CreatedBy, p.CommunityIds, p.CreatedAt, p.UpdatedBy, p.UpdatedAt, p.DeletedBy, p.DeletedAt, p.Tags, p.Labels, p.PublicDomain)
	if err != nil {
		panic(err)
	}

	err = row.Err()
	if err != nil {
		panic(err)
	}
	row.Close()

	//! spew.Dump(p.Id)
	return nil
}

func (p *Instance) DeleteInstance(cid string) error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sql_cmd := `DELETE FROM INSTANCE WHERE Id=?`
	statement, err := db.Prepare(sql_cmd)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(cid)
	if err != nil {
		panic(err)
	}
	return nil
}

func (p *Instance) UpdatedInstance() error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sql_cmd := `UPDATE INSTANCE SET Id=?,Name=?,Description=?,Type=?,Status=?,Owner=?,CreatedBy=?,CommunityIds=?,CreatedAt=?,UpdatedBy=?,UpdatedAt=?,DeletedBy=?,DeletedAt=?,Tags=?,Labels=?,PublicDomain=? WHERE Id=?`
	statement, err := db.Prepare(sql_cmd)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.Id, p.Name, p.Description, p.Type, p.Status, p.Owner, p.CreatedBy, p.CommunityIds, p.CreatedAt, p.UpdatedBy, p.UpdatedAt, p.DeletedBy, p.DeletedAt, p.Tags, p.Labels, p.PublicDomain, p.Id)
	if err != nil {
		panic(err)
	}
	return nil
}
