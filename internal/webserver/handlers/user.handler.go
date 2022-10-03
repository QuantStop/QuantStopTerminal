package handlers

/*func UpdatePassword(db *sql.DB, user *models.User, w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	if user.Status != models.UserStatusActive {
		return write.Error(errors.RouteUnauthorized)
	}

	decoder := json.NewDecoder(r.Body)
	var u models.User
	err := decoder.Decode(&u)
	if err != nil || &u == nil {
		return write.Error(errors.NoJSONBody)
	}

	// salt and hash it
	u.Salt = utils.GenerateRandomString(32)
	u.Password, err = utils.HashPassword(u.Password, u.Salt)
	if err != nil {
		return write.Error(err)
	}

	// todo:
	err = env.DB().UpdateUserPassword(r.Context(), db.UpdateUserPasswordParams{
		ID:   user.ID,
		Pass: u.Pass,
		Salt: u.Salt,
	})
	if err != nil {
		return write.Error(err)
	}

	return write.Success()
}*/

/*type WhoamiResponse struct {
	ID       uint32   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

func Whoami(bot internal.IEngine, user *models.User, w http.ResponseWriter, r *http.Request) http.HandlerFunc {

	res := WhoamiResponse{
		ID:       user.ID,
		Username: user.Username,
		Roles:    user.Roles,
	}
	return write.JSON(res)
}

func GetAllUsers(bot internal.IEngine, user *models.User, w http.ResponseWriter, r *http.Request) http.HandlerFunc {

	db, _ := bot.GetSQL("core")
	users, err := user.GetUsers(db)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("cannot get users")
			return write.Error(errors.FailedLogin)
		}
		return write.Error(errors.FailedLogin)
	}
	return write.JSON(users)

}
*/
