package model

// ProjectServer project server relationship
type ProjectServer struct {
	ID          uint32 `json:"id"`
	ProjectID   uint32 `json:"projectId"`
	ServerID    uint32 `json:"serverId"`
	ServerName  string `json:"serverName"`
	ServerIP    string `json:"serverIP"`
	ServerOwner string `json:"serverOwner"`
	CreateTime  int64  `json:"createTime"`
	UpdateTime  int64  `json:"updateTime"`
}

// ProjectServers project server relationship
type ProjectServers []ProjectServer

// AddMany add many row to table project_server
func (ps *ProjectServers) AddMany() error {
	db := NewDB()
	sqlStr := "INSERT INTO project_server (project_id, server_id, create_time, update_time) VALUES "
	vals := []interface{}{}

	for _, row := range *ps {
		sqlStr += "(?, ?, ?, ?),"
		vals = append(vals, row.ProjectID, row.ServerID, row.CreateTime, row.UpdateTime)
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	//format all vals at once
	_, err = stmt.Exec(vals...)
	return err
}

// Query server row
func (ps *ProjectServers) Query(projectID uint32) error {
	db := NewDB()
	rows, err := db.Query(
		`SELECT 
			project_server.id,
			project_id,
			server_id,
			server.name,
			server.ip,
			server.owner
		FROM project_server
		LEFT JOIN server 
		ON project_server.server_id = server.id
		WHERE project_id = ?`, projectID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var projectServer ProjectServer

		if err := rows.Scan(&projectServer.ID, &projectServer.ProjectID, &projectServer.ServerID, &projectServer.ServerName, &projectServer.ServerIP, &projectServer.ServerOwner); err != nil {
			return err
		}
		*ps = append(*ps, projectServer)
	}
	return nil
}